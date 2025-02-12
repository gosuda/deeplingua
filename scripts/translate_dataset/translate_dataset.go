package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/lemon-mint/coord/llm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fastjson"
	"gosuda.org/deeplingua/internal/translate"
	"gosuda.org/deeplingua/jsonl"
	"gosuda.org/deeplingua/normalize"
)

var (
	translationModel  llm.Model                                     // required
	customPrompt      string                                        // required
	evaluationModel   llm.Model                             = nil   // optional
	doEvaluation      bool                                  = false // optional
	customPipelinePre func(index int, v *jsonl.Value) error = func(index int, v *jsonl.Value) error {
		if string(v.GetStringBytes("custom_id")) != "" {
			return nil
		}

		v.Set("custom_id", fastjson.MustParse(fmt.Sprintf(`"%020d"`, index)))
		return nil
	} // optional (default: add a custom_id field with the index)
	customPipelinePost func(index int, v *jsonl.Value) error     // optional
	startIndex         int                                   = 0 // optional
)

var (
	lastReadIndex          atomic.Int64
	successfulTranslations atomic.Int64
	failedTranslations     atomic.Int64
)

type Job struct {
	Index int
	Value *jsonl.Value
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	log.Logger = log.Logger.Level(zerolog.DebugLevel)

	var inFile string
	var outFile string
	var inLang string
	var outLang string
	var workers int

	flag.StringVar(&inFile, "in", "", "Input file")
	flag.StringVar(&outFile, "out", "", "Output file")
	flag.StringVar(&inLang, "src", "", "Source language")
	flag.StringVar(&outLang, "dst", "", "Target language")
	flag.IntVar(&workers, "workers", 64, "Workers")
	flag.Parse()
	if inFile == "" || outFile == "" || inLang == "" || outLang == "" {
		panic("Usage: translate_dataset -in <input.jsonl> -out <output.jsonl> -src <source_lang> -dst <target_lang>")
	}

	// Load configuration
	var err error
	var config Configs
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		config_path = "config.json"
	}
	config_path, err = filepath.Abs(config_path)
	log.Debug().Str("path", config_path).Msg("loading config file")
	config_data, err := os.ReadFile(config_path)
	if err != nil {
		log.Fatal().Str("path", config_path).Err(err).Msg("failed to read config file")
	}
	err = json.Unmarshal(config_data, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config file")
	}
	ApplyConfig(&config)

	log.Info().Str("in", inFile).Str("out", outFile).Str("src", inLang).Str("dst", outLang).Int("workers", workers).Msg("starting")

	f, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := jsonl.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	wf, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	w, err := jsonl.NewWriter(wf)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	wfailf, err := os.Create(outFile + ".failed")
	if err != nil {
		panic(err)
	}
	defer wfailf.Close()

	wfail, err := jsonl.NewWriter(wfailf)
	if err != nil {
		panic(err)
	}
	defer wfail.Close()

	var wgWorkers sync.WaitGroup // Wait for all workers to finish
	var wgWriter sync.WaitGroup  // Wait for writer to finish

	stopSignal := make(chan struct{})
	jobQueue := make(chan Job, 1)
	errorQueue := make(chan *jsonl.Value, workers*2)
	completionQueue := make(chan *jsonl.Value, workers*2)

	// Start Translation Workers
	wgWorkers.Add(workers)
	for i := 0; i < workers; i++ {
		go translationWorker(i, inLang, outLang, jobQueue, completionQueue, errorQueue, &wgWorkers)
	}

	// Start Writer Worker
	wgWriter.Add(1)
	go writerWorker(completionQueue, errorQueue, w, wfail, &wgWriter)

	// Start Reader
	go reader(r, jobQueue, stopSignal)

	// Listen for SIGINT or SIGTERM to stop gracefully
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		close(stopSignal)
		signal.Stop(sigChan)
		log.Info().Msg("received stop signal, gracefully stopping remaining workers... (interrupt again to force stop)")
	}()

	wgWorkers.Wait() // Wait for all workers to finish translation
	log.Debug().Msg("all workers stopped")
	close(completionQueue) // Signal writer that no more successful jobs are coming
	close(errorQueue)      // Signal writer that no more error jobs are coming
	log.Debug().Msg("stopping writer worker")
	wgWriter.Wait() // Wait for writer to finish writing all results
	log.Debug().Msg("writer worker stopped")

	// Print statistics
	log.Info().
		Int64("Successful Translations", successfulTranslations.Load()).
		Int64("Failed Translations", failedTranslations.Load()).
		Int64("Last Read Index", lastReadIndex.Load()).
		Msg("finished")

}

func reader(r *jsonl.Reader, jobQueue chan<- Job, stopSignal <-chan struct{}) {
	log.Debug().Msg("reader started")
	defer close(jobQueue) // Close jobQueue when reader finishes

	var index int
L:
	for {
		v, err := r.Scan()
		if err != nil {
			if err != io.EOF {
				log.Err(err).Msg("error reading input file")
			}
			break // Assume io.EOF is the expected error when reaching end of file
		}
		if v == nil {
			continue
		}

		if index < startIndex {
			index++
			continue
		}
		select {
		case jobQueue <- Job{
			Index: index,
			Value: v,
		}:
			lastReadIndex.Store(int64(index + 1))
			log.Info().Int("Index", index).Msg("read successfully, queued")
		case <-stopSignal:
			log.Info().Msg("stop signal received, stopping reader")
			break L
		}
		index++
	}

	log.Debug().Msg("reader stopped")
}

func translationWorker(id int, inLang string, outLang string, jobQueue <-chan Job, completionQueue chan<- *jsonl.Value, errorQueue chan<- *jsonl.Value, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug().Int("ID", id).Msg("translation worker started")
	_ = inLang
L:
	for job := range jobQueue {
		v := job.Value
		index := job.Index

		if customPipelinePre != nil {
			if err := customPipelinePre(index, v); err != nil {
				log.Error().Int("workerID", id).Int("Index", index).Err(err).Msg("custom pipeline pre failed")
				continue
			}
		}

		credits := 3
	RL:
		for {
			if credits <= 0 {
				errorQueue <- v
				log.Error().Int("workerID", id).Int("Index", index).Msg("repeated translation fail, skipping")
				break
			}
			credits--

			normalize.NormalizeShareGPT(v)
			messages := v.GetArray("messages")

			for i := range messages {
				original := string(messages[i].GetStringBytes("content"))
				if string(messages[i].GetStringBytes("translated_content")) != "" {
					continue
				}
				translated := ""

				if !utf8.ValidString(original) {
					continue
				}
				original = normalize.Normalize(original)

				translated, err := translate.TranslateCustomPrompt(context.Background(), translationModel, original, outLang, customPrompt)
				if err != nil {
					log.Error().
						Int("workerID", id).
						Int("Index", index).
						Err(err).
						Int("tokens", credits).
						Msg("translate failed")
					time.Sleep(time.Duration(float64(10) * rand.Float64() * float64(time.Second)))
					continue RL
				}

				if !utf8.ValidString(translated) {
					log.Error().
						Int("workerID", id).
						Int("Index", index).
						Err(fmt.Errorf("deeplingua: invalid utf8 string")).
						Int("tokens", credits).
						Msg("translate failed")
					time.Sleep(time.Duration(float64(10) * rand.Float64() * float64(time.Second)))
					continue
				}
				translated = normalize.Normalize(translated)

				data, err := json.Marshal(translated)
				if err != nil {
					log.Error().Int("workerID", id).Int("Index", index).Err(err).Msg("marshal failed")
					continue L
				}
				messages[i].Set("translated_content", fastjson.MustParseBytes(data))
				v.Value.Get("messages").SetArrayItem(i, messages[i])
				if doEvaluation {
					// TODO: evaluate - this part would be moved to evaluation worker if you have one
				}
			}

			if customPipelinePost != nil {
				if err := customPipelinePost(index, v); err != nil {
					log.Error().Int("workerID", id).Int("Index", index).Err(err).Msg("custom pipeline post failed")
					continue
				}
			}

			log.Info().Int("workerID", id).Int("Index", index).Msg("translated successfully")
			completionQueue <- v
			continue L
		}
	}

	log.Debug().Int("ID", id).Msg("translation worker stopped")
}

// writerWorker handles writing both successful and failed jobs to their respective files.
func writerWorker(completionQueue <-chan *jsonl.Value, errorQueue <-chan *jsonl.Value, w *jsonl.Writer, wfail *jsonl.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug().Msg("writer worker started")

	for {
		select {
		case v, ok := <-errorQueue:
			if !ok {
				errorQueue = nil // Set to nil to stop selecting this case when channel is closed
			} else {
				failedTranslations.Add(1)
				if err := wfail.Write(v); err != nil {
					panic(err) // Writer error is critical
				}
			}
		case v, ok := <-completionQueue:
			if !ok {
				completionQueue = nil // Set to nil to stop selecting this case when channel is closed
			} else {
				successfulTranslations.Add(1)
				if err := w.Write(v); err != nil {
					panic(err) // Writer error is critical
				}
			}
		}
		if completionQueue == nil && errorQueue == nil {
			break // Exit loop when both channels are closed and drained
		}
	}

	log.Debug().Msg("writer worker stopped")
}
