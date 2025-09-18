package types

import "context"

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}
