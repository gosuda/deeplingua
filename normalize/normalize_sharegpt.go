package normalize

import (
	"github.com/valyala/fastjson"
	"gosuda.org/deeplingua/jsonl"
)

var role_values = []*fastjson.Value{
	fastjson.MustParse(`"system"`),
	fastjson.MustParse(`"user"`),
	fastjson.MustParse(`"assistant"`),
}

func NormalizeShareGPT(v *jsonl.Value) {
	// check if there is a "conversations", "conversation", "message", or "messages" field
	// if so, replace it with a "messages" field

	messages_name := ""
	if v.Get("conversations") != nil {
		messages_name = "conversations"
	} else if v.Get("conversation") != nil {
		messages_name = "conversation"
	} else if v.Get("message") != nil {
		messages_name = "message"
	} else if v.Get("messages") != nil {
		messages_name = "messages"
	} else {
		return
	}

	messages := v.Get(messages_name)
	if messages_name != "messages" {
		v.Set("messages", messages)
		v.Del(messages_name)
	}

	for _, message := range messages.GetArray() {
		// check if there is a "content", "text", "value", or "message" field

		content_name := ""
		if message.Get("content") != nil {
			content_name = "content"
		} else if message.Get("text") != nil {
			content_name = "text"
		} else if message.Get("value") != nil {
			content_name = "value"
		} else if message.Get("message") != nil {
			content_name = "message"
		} else {
			continue
		}

		if content_name != "content" {
			message.Set("content", message.Get(content_name))
			message.Del(content_name)
		}

		// check if there is a "role", or "from" field
		// if so, replace it with a "role" field
		role_name := ""
		if message.Get("role") != nil {
			role_name = "role"
		} else if message.Get("from") != nil {
			role_name = "from"
		} else {
			continue
		}

		if role_name != "role" {
			message.Set("role", message.Get(role_name))
			message.Del(role_name)
		}

		var role_val *fastjson.Value

		switch string(message.Get("role").GetStringBytes()) {
		case "system":
			role_val = role_values[0]
		case "user", "human":
			role_val = role_values[1]
		case "gpt", "assistant", "model", "completion", "bot", "ai", "gpt-4", "gpt-3.5-turbo":
			role_val = role_values[2]
		}

		message.Set("role", role_val)
	}
}
