package admin

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"nicolas.galipot.net/hazo/server/common"
	"nicolas.galipot.net/hazo/server/components"
	"nicolas.galipot.net/hazo/storage/appdb"
)

//go:embed admin.html
var adminPage string

type UserWithCapabilities struct {
	Login            string
	CreatedOn        string
	PrivateDirectory string
	Capabilities     []CapabilityInfo
}

type CapabilityInfo struct {
	Name        string
	Description string
	GrantedDate string
	GrantedBy   string
}

type MSAccountRequest struct {
	MSAccountId   string
	Email         string
	FullName      string
	RequestedDate string
}

type ViewModel struct {
	PageTitle       string
	Users           []UserWithCapabilities
	AllCapabilities []appdb.Capability
	MSRequests      []MSAccountRequest
	Debug           bool
	Error           string
	Success         string
}

func Handler(w http.ResponseWriter, r *http.Request, cc *common.Context) error {
	ctx := context.Background()

	cc.RegisterActions(NewActions(cc))
	err := cc.ExecuteActions(ctx, r)

	var errorMsg, successMsg string
	if err != nil {
		errorMsg = err.Error()
	} else if r.Method == "POST" {
		if r.PostFormValue("admin-add-user") != "" {
			successMsg = fmt.Sprintf("Successfully created user '%s'", r.PostFormValue("login"))
		} else if r.PostFormValue("admin-grant-capability") != "" {
			successMsg = fmt.Sprintf("Successfully granted '%s' capability to user '%s'",
				r.PostFormValue("grant-capability"), r.PostFormValue("grant-login"))
		} else if r.PostFormValue("admin-revoke-capability") != "" {
			successMsg = fmt.Sprintf("Successfully revoked '%s' capability from user '%s'",
				r.PostFormValue("revoke-capability"), r.PostFormValue("revoke-login"))
		} else if r.PostFormValue("admin-delete-user") != "" {
			successMsg = fmt.Sprintf("Successfully deleted user '%s'", r.PostFormValue("delete-login"))
		} else if r.PostFormValue("admin-approve-ms-request") != "" {
			successMsg = fmt.Sprintf("Successfully approved Microsoft account request for '%s'", r.PostFormValue("approve-email"))
		} else if r.PostFormValue("admin-reject-ms-request") != "" {
			successMsg = fmt.Sprintf("Successfully rejected Microsoft account request for '%s'", r.PostFormValue("reject-email"))
		}
	}

	usersData, err := cc.AppQueries().GetAllUsersWithCapabilities(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users with capabilities: %w", err)
	}

	capabilities, err := cc.AppQueries().GetAllCapabilities(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all capabilities: %w", err)
	}

	userMap := make(map[string]*UserWithCapabilities)
	for _, row := range usersData {
		user, exists := userMap[row.Login]
		if !exists {
			createdOn := ""
			if row.CreatedOn.Valid {
				createdOn = row.CreatedOn.String
			}
			user = &UserWithCapabilities{
				Login:            row.Login,
				CreatedOn:        createdOn,
				PrivateDirectory: row.PrivateDirectory,
				Capabilities:     []CapabilityInfo{},
			}
			userMap[row.Login] = user
		}
		if row.CapabilityName.Valid {
			user.Capabilities = append(user.Capabilities, CapabilityInfo{
				Name:        row.CapabilityName.String,
				Description: row.CapabilityDescription.String,
				GrantedDate: row.GrantedDate.String,
				GrantedBy:   row.GrantedBy.String,
			})
		}
	}

	users := make([]UserWithCapabilities, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, *user)
	}

	msRequestsData, err := cc.AppQueries().GetAllPendingMSAccountRequests(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MS account requests: %w", err)
	}

	msRequests := make([]MSAccountRequest, 0, len(msRequestsData))
	for _, req := range msRequestsData {
		msRequests = append(msRequests, MSAccountRequest{
			MSAccountId:   req.MsAccountID,
			Email:         req.Email,
			FullName:      req.FullName,
			RequestedDate: req.RequestedDate,
		})
	}

	tmpl := components.NewTemplate()
	tmpl = template.Must(tmpl.Parse(adminPage))
	w.Header().Add("Content-Type", "text/html")
	err = tmpl.Execute(w, ViewModel{
		PageTitle:       "User Administration",
		Users:           users,
		AllCapabilities: capabilities,
		MSRequests:      msRequests,
		Debug:           cc.Config.Debug,
		Error:           errorMsg,
		Success:         successMsg,
	})
	if err != nil {
		return fmt.Errorf("template rendering of the admin page failed: %w", err)
	}
	return nil
}
