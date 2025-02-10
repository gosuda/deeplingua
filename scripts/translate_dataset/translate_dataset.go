package main

import (
	"context"
	"encoding/json"
	"flag"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/lemon-mint/coord/llm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fastjson"
	"gosuda.org/deeplingua/internal/translate"
	"gosuda.org/deeplingua/jsonl"
	"gosuda.org/deeplingua/normalize"
)

var (
	translationModel   llm.Model                                     // required
	customPrompt       string                                        // required
	evaluationModel    llm.Model                             = nil   // optional
	doEvaluation       bool                                  = false // optional
	customPipelinePre  func(index int, v *jsonl.Value) error         // optional
	customPipelinePost func(index int, v *jsonl.Value) error         // optional
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

	wfailf, err := os.Create(outFile + ".fail")
	if err != nil {
		panic(err)
	}
	defer wfailf.Close()

	wfail, err := jsonl.NewWriter(wfailf)
	if err != nil {
		panic(err)
	}
	defer wfail.Close()

	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	type Job struct {
		Index int
		Value *jsonl.Value
	}

	jobQueue := make(chan Job, workers*2)
	errorQueue := make(chan *jsonl.Value, workers*2)
	completionQueue := make(chan *jsonl.Value, workers*2)

	wg.Add(workers)
	// Start workers
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
		L:
			for v := range jobQueue {
				if customPipelinePre != nil {
					if err := customPipelinePre(v.Index, v.Value); err != nil {
						log.Error().Int("worker", id).Int("Index", v.Index).Err(err).Msg("custom pipeline pre failed")
						continue
					}
				}

				credits := 3
			RL:
				for {
					if credits <= 0 {
						errorQueue <- v.Value
						log.Error().Int("worker", id).Int("Index", v.Index).Msg("repeated translation fail, skipping")
						break
					}
					credits--

					normalize.NormalizeShareGPT(v.Value)
					messages := v.Value.GetArray("messages")

					for i := range messages {
						original := string(messages[i].GetStringBytes("content"))
						if string(messages[i].GetStringBytes("translated_content")) != "" {
							continue
						}
						//original = normalize.Normalize(original)
						translated := ""

						translated, err := translate.TranslateCustomPrompt(context.Background(), translationModel, original, outLang, customPrompt)
						if err != nil {
							log.Error().Int("worker", id).Int("Index", v.Index).Err(err).Int("tokens", credits).Msg("translate failed")
							time.Sleep(time.Duration(float64(10) * rand.Float64() * float64(time.Second)))
							continue RL
						}
						//translated = normalize.Normalize(translated)

						data, err := json.Marshal(translated)
						if err != nil {
							log.Error().Int("worker", id).Int("Index", v.Index).Err(err).Msg("marshal failed")
							continue L
						}
						messages[i].Set("translated_content", fastjson.MustParseBytes(data))
						v.Value.Get("messages").SetArrayItem(i, messages[i])
						if doEvaluation {
							// TODO: evaluate
						}
					}

					if customPipelinePost != nil {
						if err := customPipelinePost(v.Index, v.Value); err != nil {
							log.Error().Int("worker", id).Int("Index", v.Index).Err(err).Msg("custom pipeline post failed")
							continue
						}
					}

					log.Info().Int("worker", id).Int("Index", v.Index).Msg("translated successfully")
					completionQueue <- v.Value
					continue L
				}
			}
		}(i)
	}

	wg2.Add(1)
	// start writer
	go func() {
		defer wg2.Done()
	L:
		for {
			select {
			case v, ok := <-errorQueue:
				if !ok {
					continue L
				}
				if err := wfail.Write(v); err != nil {
					panic(err)
				}
			case v, ok := <-completionQueue:
				if !ok {
					break L
				}

				if err := w.Write(v); err != nil {
					panic(err)
				}
			}
		}

		// drain error queue
		for v := range errorQueue {
			if err := wfail.Write(v); err != nil {
				panic(err)
			}
		}
	}()

	// start reader
	go func() {
		defer close(jobQueue)

		var index int
		for {
			v, err := r.Scan()
			if err != nil {
				break
			}
			if v == nil {
				continue
			}

			jobQueue <- Job{
				Index: index,
				Value: v,
			}
			log.Info().Int("Index", index).Msg("read successfully, queued")

			index++
		}
	}()

	wg.Wait()
	close(completionQueue)
	close(errorQueue)

	time.Sleep(time.Second * 20)
}
