package user

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"nicolas.galipot.net/hazo/db"
)

type T struct {
	Login            string
	privateDirectory string
}

var ErrForbiddenAccess = fmt.Errorf("current user cannot access this location")

func Register(login string) (*T, error) {
	_, queries, err := db.OpenCommon()
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

func (u *T) GetDataset(dsName string) (db.PrivateDataset, error) {
	dsPath := path.Clean(path.Join(u.privateDirectory, fmt.Sprintf("%s.sq3", dsName)))
	inUserDir, err := filepath.Match(path.Join(u.privateDirectory, "*.sq3"), dsPath)
	if err != nil {
		return db.InvalidPrivateDataset, fmt.Errorf("could not find private dataset '%s': %w", dsName, err)
	}
	if inUserDir {
		return db.PrivateDataset(dsPath), nil
	}
	return db.InvalidPrivateDataset, ErrForbiddenAccess
}

func (u *T) ListDatasets() ([]db.Dataset, error) {
	files, err := filepath.Glob(path.Join(u.privateDirectory, "*.sq3"))
	if err != nil {
		return nil, fmt.Errorf("could not read dataset directory of user '%s': %w", u.Login, err)
	}
	ds := make([]db.Dataset, len(files))
	for i, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve file information about '%s: %w'", path, err)
		}
		info.ModTime()
		ds[i].Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		ds[i].Path = path
		ds[i].LastModified = info.ModTime().Format("2006-01-02 15:04:05")
	}
	return ds, nil
}
