package evaluate

import (
	"context"

	"github.com/lemon-mint/coord/llm"
	"gosuda.org/deeplingua/internal/judge"
)

func EvaluateTranslationLLMJudge(
	ctx context.Context,
	original, translated string,
	originalLanguage, translatedLanguage string,
	llm llm.Model,
) (float64, error) {
	return judge.EvaluateTranslation(
		ctx,
		llm,
		originalLanguage,
		translatedLanguage,
		original,
		translated,
	)
}
