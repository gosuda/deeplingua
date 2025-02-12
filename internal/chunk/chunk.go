package chunk

import (
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"cloud.google.com/go/vertexai/genai/tokenizer"
	"github.com/rs/zerolog/log"
)

var tok, _ = tokenizer.New("gemini-1.5-flash")

// ChunkMarkdown splits the input text into smaller chunks, ensuring that each chunk does not exceed the token limit.
// It handles code blocks and tries to split paragraphs at natural breakpoints (e.g., periods) to preserve the original formatting.
// The resulting chunks are returned as a slice of strings.
func ChunkMarkdown(input string) []string {
	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0
	inCodeBlock := false

	// Split the input into paragraphs, keeping the delimiters
	paragraphs := strings.SplitAfter(input, "\n\n")

	for _, paragraph := range paragraphs {
		// Check if this paragraph is a code block
		if strings.HasPrefix(paragraph, "```") || strings.HasSuffix(paragraph, "```") {
			inCodeBlock = !inCodeBlock
		}

		paragraphTokens, err := tok.CountTokens(genai.Text(paragraph))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to count tokens")
		}

		// If adding this paragraph would exceed the token limit or it's a code block
		if currentTokens+int(paragraphTokens.TotalTokens) > 4096 || inCodeBlock {
			// If the current chunk is not empty, add it to chunks
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
				currentTokens = 0
			}

			// If this paragraph itself exceeds 4096 tokens, split it
			if int(paragraphTokens.TotalTokens) > 4096 {
				lines := strings.SplitAfter(paragraph, "\n")
				for _, line := range lines {
					lineTokens, _ := tok.CountTokens(genai.Text(line))
					if int(lineTokens.TotalTokens) > 4096 {
						// Split by rune count or "."
						runes := []rune(line)
						for len(runes) > 0 {
							splitIndex := min(4096, len(runes))
							for idx := range runes {
								if runes[idx] == '.' ||
									runes[idx] == '?' ||
									runes[idx] == '!' ||
									runes[idx] == ';' {
									splitIndex = idx + 1
									break
								}
							}
							chunks = append(chunks, string(runes[:splitIndex]))
							runes = runes[splitIndex:]
						}
					} else {
						if currentTokens+int(lineTokens.TotalTokens) > 4096 {
							chunks = append(chunks, currentChunk.String())
							currentChunk.Reset()
							currentTokens = 0
						}
						currentChunk.WriteString(line)
						currentTokens += int(lineTokens.TotalTokens)
					}
				}
			} else {
				// Add the entire paragraph as a chunk
				chunks = append(chunks, paragraph)
			}
		} else {
			// Add the paragraph to the current chunk
			currentChunk.WriteString(paragraph)
			currentTokens += int(paragraphTokens.TotalTokens)
		}
	}

	// Add the last chunk if it's not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	grouped := groupChunks(chunks, 4096)

	var finalChunks []string
	for _, group := range grouped {
		chunk := strings.Join(group, "")
		finalChunks = append(finalChunks, chunk)
	}

	return finalChunks
}

func groupChunks(chunks []string, maxTokens int) [][]string {
	var groupedChunks [][]string
	var currentGroup []string

	currentTokens := 0

	for _, chunk := range chunks {
		chunkTokens, err := tok.CountTokens(genai.Text(chunk))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to count tokens")
		}
		if currentTokens+int(chunkTokens.TotalTokens) > maxTokens {
			groupedChunks = append(groupedChunks, currentGroup)
			currentGroup = []string{chunk}
			currentTokens = int(chunkTokens.TotalTokens)
		} else {
			currentGroup = append(currentGroup, chunk)
			currentTokens += int(chunkTokens.TotalTokens)
		}
	}

	if len(currentGroup) > 0 {
		groupedChunks = append(groupedChunks, currentGroup)
	}

	return groupedChunks
}
