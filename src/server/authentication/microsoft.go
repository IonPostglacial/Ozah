package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"nicolas.galipot.net/hazo/storage/appdb"
)

type MSConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type MSTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type MSUserInfo struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
}

func LoadMSConfig() *MSConfig {
	clientID := os.Getenv("MICROSOFT_CLIENT_ID")
	clientSecret := os.Getenv("MICROSOFT_CLIENT_SECRET")
	redirectURI := os.Getenv("MICROSOFT_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil
	}

	return &MSConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

func (c *MSConfig) GenerateAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", c.RedirectURI)
	params.Set("scope", "openid profile email User.Read")
	params.Set("state", state)

	return "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?" + params.Encode()
}

func (c *MSConfig) ExchangeCode(code string) (*MSTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", c.RedirectURI)

	resp, err := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get Microsoft OAuth token: %s", string(body))
	}

	var tokenResp MSTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func GetMSUserInfo(accessToken string) (*MSUserInfo, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get Microsoft Graph user info: %s", string(body))
	}

	var userInfo MSUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HandleMSLogin(w http.ResponseWriter, r *http.Request, config *MSConfig) {
	if config == nil {
		http.Error(w, "Microsoft authentication not configured", http.StatusInternalServerError)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "ms_auth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	authURL := config.GenerateAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func HandleMSCallback(w http.ResponseWriter, r *http.Request, config *MSConfig, queries *appdb.Queries) {
	if config == nil {
		http.Error(w, "Microsoft authentication not configured", http.StatusInternalServerError)
		return
	}

	stateCookie, err := r.Cookie("ms_auth_state")
	if err != nil {
		http.Error(w, "Missing state cookie", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" || state != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "ms_auth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		errorDesc := r.URL.Query().Get("error_description")
		if errorDesc == "" {
			errorDesc = "No authorization code received"
		}
		http.Error(w, errorDesc, http.StatusBadRequest)
		return
	}

	tokenResp, err := config.ExchangeCode(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange code: %v", err), http.StatusInternalServerError)
		return
	}

	userInfo, err := GetMSUserInfo(tokenResp.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	creds, err := queries.GetCredentialsByMSAccountId(ctx, sql.NullString{
		String: userInfo.ID,
		Valid:  true,
	})

	if err == nil {
		sessionToken, err := generateToken()
		if err != nil {
			http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(2 * time.Hour)
		_, err = queries.InsertSession(ctx, appdb.InsertSessionParams{
			Token:      sessionToken,
			Login:      creds.Login,
			ExpiryDate: expiresAt.Format(dateFormat),
		})
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	_, err = queries.GetMSAccountRequest(ctx, userInfo.ID)
	if err == nil {
		http.Error(w, "Your access request is pending admin approval", http.StatusForbidden)
		return
	}

	email := userInfo.Mail
	if email == "" {
		email = userInfo.UserPrincipalName
	}

	_, err = queries.CreateMSAccountRequest(ctx, appdb.CreateMSAccountRequestParams{
		MsAccountID:   userInfo.ID,
		Email:         email,
		FullName:      userInfo.DisplayName,
		RequestedDate: time.Now().Format(time.RFC3339),
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create access request: %v", err), http.StatusInternalServerError)
		return
	}

	http.Error(w, "Access request submitted. Please wait for admin approval.", http.StatusForbidden)
}
