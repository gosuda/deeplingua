package translation

import (
	"context"

	"github.com/lemon-mint/coord/llm"
	"gosuda.org/deeplingua/internal/translate"
)

var ErrFailedToTranslate = translate.ErrFailedToTranslate

func TranslateText(ctx context.Context, l llm.Model, input, targetLanguage string) (string, error) {
	return translate.Translate(ctx, l, input, targetLanguage)
}
