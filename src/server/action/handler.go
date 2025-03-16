package action

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Handler func(ctx context.Context, r *http.Request) error

func NewHandlerWithIntArgument(argName string, cb func(ctx context.Context, value int) error) Handler {
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

func NewHandlerWithStringArgument(argName string, cb func(ctx context.Context, value string) error) Handler {
	return func(ctx context.Context, r *http.Request) error {
		if value := r.PostFormValue(argName); value != "" {
			return cb(ctx, value)
		}
		return nil
	}
}
