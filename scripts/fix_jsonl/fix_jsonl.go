package main

import (
	"encoding/json"
	"os"

	"gosuda.org/deeplingua/jsonl"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rf, err := jsonl.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer rf.Close()

	wf, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer wf.Close()

	type Messages struct {
		Content           string `json:"content"`
		Role              string `json:"role"`
		TranslatedContent string `json:"translated_content"`
	}

	type Row struct {
		Messages []Messages `json:"messages"`
		CustomID string     `json:"custom_id"`
	}

	for {
		r, err := rf.Scan()
		if err != nil {
			break
		}

		var m Row

		arr := r.GetArray("messages")
		for i := range arr {
			m.Messages = append(m.Messages, Messages{
				Content:           string(arr[i].GetStringBytes("content")),
				Role:              string(arr[i].GetStringBytes("role")),
				TranslatedContent: string(arr[i].GetStringBytes("translated_content")),
			})
		}
		m.CustomID = string(r.GetStringBytes("custom_id"))

		data, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		wf.Write(data)
		wf.Write([]byte("\n"))
	}
}
