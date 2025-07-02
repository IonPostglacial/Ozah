package user

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"nicolas.galipot.net/hazo/storage"
)

type T struct {
	Login            string
	privateDirectory string
}

var ErrForbiddenAccess = fmt.Errorf("current user cannot access this location")

func Register(login string) (*T, error) {
	_, queries, err := storage.OpenAppDb()
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

func (u *T) GetDataset(dsName string) (storage.PrivateDataset, error) {
	dsPath := getUserPrivateDatasetPath(u.privateDirectory, dsName)
	inUserDir, err := filepath.Match(path.Join(u.privateDirectory, "*.sq3"), dsPath)
	if err != nil {
		return storage.InvalidPrivateDataset, fmt.Errorf("could not find private dataset '%s': %w", dsName, err)
	}
	if inUserDir {
		return storage.PrivateDataset(dsPath), nil
	}
	return storage.InvalidPrivateDataset, ErrForbiddenAccess
}

func (u *T) ListDatasets() ([]storage.Dataset, error) {
	files, err := filepath.Glob(path.Join(u.privateDirectory, "*.sq3"))
	if err != nil {
		return nil, fmt.Errorf("could not read dataset directory of user '%s': %w", u.Login, err)
	}
	ds := make([]storage.Dataset, len(files))
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

func (u *T) GetReadableSharedDatasets() ([]storage.SharedDataset, error) {
	_, queries, err := storage.OpenAppDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	datasets, err := queries.GetReadableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	ds := make([]storage.SharedDataset, len(datasets))
	for i, d := range datasets {
		path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		ds[i].Dataset.Name = d.Name
		ds[i].Dataset.Path = path
		ds[i].Dataset.LastModified = info.ModTime().Format("2006-01-02 15:04:05")
		ds[i].Creator = d.CreatorUserLogin
		ds[i].Mode = "read"
	}
	return ds, nil
}

func (u *T) GetWritableSharedDatasets() ([]storage.SharedDataset, error) {
	_, queries, err := storage.OpenAppDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	datasets, err := queries.GetWritableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	ds := make([]storage.SharedDataset, len(datasets))
	for i, d := range datasets {
		path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		ds[i].Dataset.Name = d.Name
		ds[i].Dataset.Path = path
		ds[i].Dataset.LastModified = info.ModTime().Format("2006-01-02 15:04:05")
		ds[i].Creator = d.CreatorUserLogin
		ds[i].Mode = "write"
	}
	return ds, nil
}
