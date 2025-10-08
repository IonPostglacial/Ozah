package user

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/dataset"
)

type T struct {
	Login            string
	privateDirectory string
}

var ErrForbiddenAccess = fmt.Errorf("current user cannot access this location")

func Register(login string) (*T, error) {
	_, queries, err := app.OpenDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	conf, err := queries.GetUserConfiguration(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve configuration of user '%s': %w", login, err)
	}
	return &T{
		Login:            login,
		privateDirectory: conf.PrivateDirectory,
	}, nil
}

func getUserPrivateDatasetPath(privateDirectory, dsName string) string {
	return path.Clean(path.Join(privateDirectory, fmt.Sprintf("%s.sq3", dsName)))
}

func (u *T) GetDataset(dsName string) (dataset.Private, error) {
	dsPath := getUserPrivateDatasetPath(u.privateDirectory, dsName)
	inUserDir, err := filepath.Match(path.Join(u.privateDirectory, "*.sq3"), dsPath)
	if err != nil {
		return dataset.InvalidPrivate, fmt.Errorf("could not find private dataset '%s': %w", dsName, err)
	}
	if inUserDir {
		return dataset.Private(dsPath), nil
	}
	return dataset.InvalidPrivate, ErrForbiddenAccess
}

func (u *T) ListDatasets() ([]dataset.T, error) {
	files, err := filepath.Glob(path.Join(u.privateDirectory, "*.sq3"))
	if err != nil {
		return nil, fmt.Errorf("could not read dataset directory of user '%s': %w", u.Login, err)
	}
	ds := make([]dataset.T, len(files))
	for i, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		ds[i].Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ds[i].Path = path
		ds[i].LastModified = info.ModTime().Format("2006-01-02 15:04:05")
	}
	return ds, nil
}

func (u *T) GetReadableSharedDatasets() ([]dataset.Shared, error) {
	_, queries, err := app.OpenDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	datasets, err := queries.GetReadableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	ds := make([]dataset.Shared, len(datasets))
	for i, d := range datasets {
		path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		ds[i].T.Name = d.Name
		ds[i].T.Path = path
		ds[i].T.LastModified = info.ModTime().Format("2006-01-02 15:04:05")
		ds[i].Creator = d.CreatorUserLogin
		ds[i].Mode = "read"
	}
	return ds, nil
}

func (u *T) GetWritableSharedDatasets() ([]dataset.Shared, error) {
	_, queries, err := app.OpenDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	datasets, err := queries.GetWritableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	ds := make([]dataset.Shared, len(datasets))
	for i, d := range datasets {
		path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		ds[i].T.Name = d.Name
		ds[i].T.Path = path
		ds[i].T.LastModified = info.ModTime().Format("2006-01-02 15:04:05")
		ds[i].Creator = d.CreatorUserLogin
		ds[i].Mode = "write"
	}
	return ds, nil
}
