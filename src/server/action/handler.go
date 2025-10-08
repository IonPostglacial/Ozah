package action

import (
	"context"
	"fmt"
	"mime/multipart"
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

func NewActionWithFileUpload(buttonName, fileFieldName string, cb func(ctx context.Context, file multipart.File, header *multipart.FileHeader) error) Action {
	return func(ctx context.Context, r *http.Request) error {
		if r.PostFormValue(buttonName) == "" {
			return nil
		}
		if r.MultipartForm == nil {
			if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
				return fmt.Errorf("could not parse multipart form: %w", err)
			}
		}
		file, header, err := r.FormFile(fileFieldName)
		if err != nil {
			if err == http.ErrMissingFile {
				return nil // No file provided, skip this action
			}
			return fmt.Errorf("could not get file from form: %w", err)
		}
		defer file.Close()

		return cb(ctx, file, header)
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
