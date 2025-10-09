package user

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"nicolas.galipot.net/hazo/storage/app"
	"nicolas.galipot.net/hazo/storage/appdb"
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

func (u *T) CanAccessDataset(dsName string) (dataset.Private, error) {
	dsPath := getUserPrivateDatasetPath(u.privateDirectory, dsName)
	inUserDir, err := filepath.Match(path.Join(u.privateDirectory, "*.sq3"), dsPath)
	if err != nil {
		return dataset.InvalidPrivate, fmt.Errorf("could not find dataset '%s': %w", dsName, err)
	}
	if inUserDir {
		if _, err := os.Stat(dsPath); err == nil {
			return dataset.Private(dsPath), nil
		}
	}

	_, queries, err := app.OpenDb()
	if err != nil {
		return dataset.InvalidPrivate, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()

	readableDatasets, err := queries.GetReadableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return dataset.InvalidPrivate, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	for _, d := range readableDatasets {
		if d.Name == dsName {
			path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
			return dataset.Private(path), nil
		}
	}

	writableDatasets, err := queries.GetWritableDatasetSharedWithUser(ctx, u.Login)
	if err != nil {
		return dataset.InvalidPrivate, fmt.Errorf("could not retrieve shared datasets for user '%s': %w", u.Login, err)
	}
	for _, d := range writableDatasets {
		if d.Name == dsName {
			path := getUserPrivateDatasetPath(d.PrivateDirectory, d.Name)
			return dataset.Private(path), nil
		}
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

func (u *T) GetCapabilities() ([]string, error) {
	_, queries, err := app.OpenDb()
	if err != nil {
		return nil, fmt.Errorf("could not open the users database: %w", err)
	}
	ctx := context.Background()
	capabilities, err := queries.GetUserCapabilities(ctx, u.Login)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve capabilities for user '%s': %w", u.Login, err)
	}
	names := make([]string, len(capabilities))
	for i, cap := range capabilities {
		names[i] = cap.CapabilityName
	}
	return names, nil
}

func (u *T) HasCapability(capabilityName string) (bool, error) {
	capabilities, err := u.GetCapabilities()
	if err != nil {
		return false, err
	}
	for _, cap := range capabilities {
		if cap == capabilityName {
			return true, nil
		}
	}
	return false, nil
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

type CreateUserParams struct {
	Login            string
	Password         string
	PrivateDirectory string
	Capabilities     []string
	GrantedBy        string
}

func Create(ctx context.Context, params CreateUserParams) error {
	if params.Login == "" || params.Password == "" || params.PrivateDirectory == "" {
		return fmt.Errorf("login, password, and private directory are required")
	}

	if err := os.MkdirAll(params.PrivateDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory '%s': %w", params.PrivateDirectory, err)
	}

	_, queries, err := app.OpenDb()
	if err != nil {
		return fmt.Errorf("could not open users database: %w", err)
	}

	const bcryptCost = 11
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	_, err = queries.InsertCredentials(ctx, appdb.InsertCredentialsParams{
		Login:      params.Login,
		Encryption: "bcrypt",
		Password:   string(hash),
	})
	if err != nil {
		return fmt.Errorf("could not insert credentials of user '%s': %w", params.Login, err)
	}

	_, err = queries.InsertUserConfiguration(ctx, appdb.InsertUserConfigurationParams{
		Login:            params.Login,
		PrivateDirectory: params.PrivateDirectory,
	})
	if err != nil {
		return fmt.Errorf("could not insert configuration of user '%s': %w", params.Login, err)
	}

	if len(params.Capabilities) > 0 {
		grantedDate := time.Now().Format("2006-01-02 15:04:05")
		grantedBy := params.GrantedBy
		if grantedBy == "" {
			grantedBy = params.Login
		}
		for _, cap := range params.Capabilities {
			_, err = queries.GrantUserCapability(ctx, appdb.GrantUserCapabilityParams{
				UserLogin:      params.Login,
				CapabilityName: cap,
				GrantedDate:    grantedDate,
				GrantedBy:      grantedBy,
			})
			if err != nil {
				return fmt.Errorf("could not grant capability '%s' to user '%s': %w", cap, params.Login, err)
			}
		}
	}

	return nil
}
