package admin

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"nicolas.galipot.net/hazo/server/action"
	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/storage/appdb"
	"nicolas.galipot.net/hazo/user"
)

type actions struct {
	cc *common.Context
}

func NewActions(cc *common.Context) *actions {
	return &actions{cc}
}

func (h *actions) addUser(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-add-user") == "" {
		return nil
	}

	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	folderPath := r.PostFormValue("folder")

	return user.Create(ctx, user.CreateUserParams{
		Login:            login,
		Password:         password,
		PrivateDirectory: folderPath,
		Capabilities:     nil,
		GrantedBy:        h.cc.User.Login,
	})
}

func (h *actions) deleteUser(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-delete-user") == "" {
		return nil
	}

	login := r.PostFormValue("delete-login")

	if login == "" {
		return fmt.Errorf("login is required")
	}

	if login == h.cc.User.Login {
		return fmt.Errorf("you cannot delete your own account")
	}

	return h.cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		capabilities, err := qtx.GetUserCapabilities(ctx, login)
		if err == nil {
			for _, cap := range capabilities {
				_, err = qtx.RevokeUserCapability(ctx, appdb.RevokeUserCapabilityParams{
					UserLogin:      login,
					CapabilityName: cap.CapabilityName,
				})
				if err != nil {
					return err
				}
			}
		}
		_, err = qtx.DeleteUserSessions(ctx, login)
		if err != nil {
			return err
		}
		_, err = qtx.DeleteCredentials(ctx, login)
		return err
	})
}

func (h *actions) grantCapability(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-grant-capability") == "" {
		return nil
	}

	login := r.PostFormValue("grant-login")
	capability := r.PostFormValue("grant-capability")

	if login == "" || capability == "" {
		return fmt.Errorf("login and capability are required")
	}

	return h.cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		_, err := qtx.GrantUserCapability(ctx, appdb.GrantUserCapabilityParams{
			UserLogin:      login,
			CapabilityName: capability,
			GrantedDate:    time.Now().Format("2006-01-02 15:04:05"),
			GrantedBy:      h.cc.User.Login,
		})
		return err
	})
}

func (h *actions) revokeCapability(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-revoke-capability") == "" {
		return nil
	}

	login := r.PostFormValue("revoke-login")
	capability := r.PostFormValue("revoke-capability")

	if login == "" || capability == "" {
		return fmt.Errorf("login and capability are required")
	}

	return h.cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		_, err := qtx.RevokeUserCapability(ctx, appdb.RevokeUserCapabilityParams{
			UserLogin:      login,
			CapabilityName: capability,
		})
		return err
	})
}

func (h *actions) approveMSRequest(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-approve-ms-request") == "" {
		return nil
	}

	msAccountId := r.PostFormValue("approve-ms-account-id")
	login := r.PostFormValue("approve-login")
	folderPath := r.PostFormValue("approve-folder")

	if msAccountId == "" || login == "" || folderPath == "" {
		return fmt.Errorf("MS account ID, login, and folder path are required")
	}

	return h.cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		_, err := qtx.GetCredentials(ctx, login)
		if err != nil {
			if err := user.CreateDirectory(folderPath); err != nil {
				return fmt.Errorf("could not create directory '%s': %w", folderPath, err)
			}

			_, err = qtx.InsertCredentials(ctx, appdb.InsertCredentialsParams{
				Login:      login,
				Encryption: "ms-oauth",
				Password:   "", // No password for MS auth users
			})
			if err != nil {
				return fmt.Errorf("could not insert credentials for user '%s': %w", login, err)
			}

			_, err = qtx.InsertUserConfiguration(ctx, appdb.InsertUserConfigurationParams{
				Login:            login,
				PrivateDirectory: folderPath,
			})
			if err != nil {
				return fmt.Errorf("could not insert configuration for user '%s': %w", login, err)
			}
		}

		_, err = qtx.LinkMSAccountToCredentials(ctx, appdb.LinkMSAccountToCredentialsParams{
			MsAccountID:  sql.NullString{String: msAccountId, Valid: true},
			LastModified: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
			Login:        login,
		})
		if err != nil {
			return err
		}

		_, err = qtx.ApproveMSAccountRequest(ctx, appdb.ApproveMSAccountRequestParams{
			ProcessedDate: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
			ProcessedBy:   sql.NullString{String: h.cc.User.Login, Valid: true},
			LinkedLogin:   sql.NullString{String: login, Valid: true},
			MsAccountID:   msAccountId,
		})
		return err
	})
}

func (h *actions) rejectMSRequest(ctx context.Context, r *http.Request) error {
	if r.PostFormValue("admin-reject-ms-request") == "" {
		return nil
	}

	msAccountId := r.PostFormValue("reject-ms-account-id")

	if msAccountId == "" {
		return fmt.Errorf("MS account ID is required")
	}

	return h.cc.AppQueriesTx(func(qtx *appdb.Queries) error {
		_, err := qtx.RejectMSAccountRequest(ctx, appdb.RejectMSAccountRequestParams{
			ProcessedDate: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
			ProcessedBy:   sql.NullString{String: h.cc.User.Login, Valid: true},
			MsAccountID:   msAccountId,
		})
		return err
	})
}

func (h *actions) Register(reg *action.Registry) {
	reg.AppendAction(action.Action(h.addUser))
	reg.AppendAction(action.Action(h.deleteUser))
	reg.AppendAction(action.Action(h.grantCapability))
	reg.AppendAction(action.Action(h.revokeCapability))
	reg.AppendAction(action.Action(h.approveMSRequest))
	reg.AppendAction(action.Action(h.rejectMSRequest))
}
