package brain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/mouuff/SmartCut/pkg/brain"
	"github.com/mouuff/SmartCut/pkg/types"
)

func TestOllamaBrainGenerateString(t *testing.T) {
	ctx := context.Background()
	gen, err := brain.NewOllamaBrain("")

	if err != nil {
		t.Fatal(err)
	}

	r := &types.PromptRequest{
		Model:        "llama3.2",
		SystemPrompt: "You are a translation assistant. Your only task is to translate any input text into clear and natural French. Do not add explanations, comments, or extra details—only provide the translation.",
		PropertyName: "translated_text",
		Prompt:       "Translate this to french: 'hello'",
	}

	result, err := gen.GenerateString(ctx, r)
	lowerResult := strings.ToLower(result)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(lowerResult, "bonjour") {
		t.Fatal("Did not get expected translation, got: " + lowerResult)
	}
}
