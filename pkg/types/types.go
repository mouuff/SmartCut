package types

import (
	"context"
)

// InputText represents the result of a user input
type InputText struct {
	IsExplicit bool
	Text       string
}

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}
