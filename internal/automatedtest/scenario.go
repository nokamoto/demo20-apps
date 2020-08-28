package automatedtest

import (
	"go.uber.org/zap"
)

// Scenario represents an automated test case.
type Scenario struct {
	Name string
	Run  func(State, *zap.Logger) (State, error)
}

// State represents a state of automated test cases.
type State map[string]string

// Scenarios represents a series of automated test cases.
type Scenarios []Scenario

func (xs Scenarios) run(logger *zap.Logger) {
	logger.Info("test start", zap.Int("size", len(xs)))
	state := make(State)

	for i, x := range xs {
		scoped := logger.With(zap.String("name", x.Name), zap.Int("number", i))

		scoped.Info("run")

		updated, err := x.Run(state, scoped)
		if err != nil {
			scoped.Fatal("failed", zap.Error(err))
		}

		scoped.Info("succeeded")

		state = updated
	}

	logger.Info("test end")
}
