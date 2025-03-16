package action

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Action func(ctx context.Context, r *http.Request) error

func NewActionWithIntArgument(argName string, cb func(ctx context.Context, value int) error) Action {
	return func(ctx context.Context, r *http.Request) error {
		if rawValue := r.PostFormValue(argName); rawValue != "" {
			value, err := strconv.Atoi(rawValue)
			if err != nil {
				return fmt.Errorf("'%s' is not a valid argument for '%s': %w", rawValue, argName, err)
			}
			return cb(ctx, value)
		}
		return nil
	}
}

func NewActionWithStringArgument(argName string, cb func(ctx context.Context, value string) error) Action {
	return func(ctx context.Context, r *http.Request) error {
		if value := r.PostFormValue(argName); value != "" {
			return cb(ctx, value)
		}
		return nil
	}
}

type Registry struct {
	actions []Action
}

func NewRegistry() *Registry {
	return &Registry{
		actions: make([]Action, 0, 8),
	}
}

type Registrable interface {
	Register(ar *Registry)
}

func (ar *Registry) AppendAction(a Action) {
	ar.actions = append(ar.actions, a)
}

func (ar *Registry) Register(r Registrable) {
	r.Register(ar)
}

func (ar *Registry) ExecuteActions(ctx context.Context, r *http.Request) error {
	for _, action := range ar.actions {
		err := action(ctx, r)
		if err != nil {
			return err
		}
	}
	return nil
}
