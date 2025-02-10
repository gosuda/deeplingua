package main

import (
	"os"
	"strings"

	"gosuda.org/deeplingua/jsonl"
	"gosuda.org/deeplingua/normalize"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: normalize_sharegpt_main <input.jsonl>")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := jsonl.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	wf, err := os.Create(strings.TrimSuffix(os.Args[1], ".jsonl") + ".normalized.jsonl")
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	w, err := jsonl.NewWriter(wf)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	for {
		rec, err := r.Scan()
		if err != nil {
			break
		}
		normalize.NormalizeShareGPT(rec)

		if err := w.Write(rec); err != nil {
			panic(err)
		}
	}
}
