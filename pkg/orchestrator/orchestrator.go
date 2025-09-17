package orchestrator

import (
	"context"
	"log"

	"golang.design/x/clipboard"
)

type Brain interface {
	GenerateString(ctx context.Context, propertyName, prompt string) (string, error)
}

type GenerationResult struct {
	Text string
}

type GenerationRequest struct {
	Prompt       string
	PropertyName string
}

type Orchestrator struct {
	Context context.Context
	Brain   Brain
	In      chan GenerationRequest
	Out     chan GenerationResult
}

func NewOrchestrator(context context.Context, brain Brain) *Orchestrator {
	return &Orchestrator{
		Context: context,
		Brain:   brain,
		In:      make(chan GenerationRequest),
		Out:     make(chan GenerationResult),
	}
}

func (o *Orchestrator) Start() {
	for req := range o.In {
		// Process the request (placeholder logic)
		result, err := o.Brain.GenerateString(o.Context, req.PropertyName, req.Prompt)

		if err != nil {
			log.Println("Error generating:", err)
			continue
		} else {
			log.Println("Generated:", result)
			o.Out <- GenerationResult{
				Text: result,
			}
		}
	}
}

func (o *Orchestrator) StartFeedFromClipboard() {
	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {
		// print out clipboard data whenever it is changed
		println(string(data))

		o.In <- GenerationRequest{
			Prompt:       string(data),
			PropertyName: "result",
		}
	}

	log.Println("Stopped feeding from clipboard.")

}
