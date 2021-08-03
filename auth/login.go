package auth

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	_ "encoding/json"
	"errors"
	_ "errors"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/rbaderts/poker/domain"
	"github.com/rbaderts/poker/util"
	_ "golang.org/x/oauth2"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

const ExampleAppState = "I wish to wash my irish wristwatch"

type handler struct{}

var LoggedInUsers map[int]bool

func init() {
	LoggedInUsers = make(map[int]bool)
}
func isUserLoggedIn(id int) bool {
	_, ok := LoggedInUsers[id]

	if ok {
		return true
	}
	return false
}

func DevLogin(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	//session, err := GetAuthSession(r)
	session, err := GetAuthSession(r)

	if err != nil {
		fmt.Printf("err1 = %v\n", err)
		return err
	}

	// Generate random state
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		fmt.Printf("err2 = %v\n", err)
		return err
	}

	//data := []byte("string of data")
	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encodedData, b)
	state := string(encodedData)

	session.Values["state"] = state

	session.Values["given_name"] = "DevUser"
	session.Values["profile"] = "Dev"
	//	session.Values["user_id"] = user.Id

	//	session.Values["access_token"] = token.AccessToken
	//session.Values["profile"] = profile
	session.Values["uid"] = 1
	session.Values["accountId"] = 1
	err = session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)

	return nil
}

//func LoginHandler(c *gin.Context) {

func LoginHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	//func LoginHandler(w http.ResponseWriter, r *http.Request)

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return util.StatusError{500, err}
	}

	//data := []byte("string of data")
	encodedData := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encodedData, b)
	state := string(encodedData)

	//session, err := GetAuthSession(r)
	session, err := GetAuthSession(r)

	if err != nil {
		return util.StatusError{500, err}
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		return util.StatusError{500, err}
	}

	authenticator, err := NewAuthenticator()
	if err != nil {
		fmt.Printf("err == %v\n", err)
		return util.StatusError{500, err}
	}

	//enableCors(&w)
	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
	return nil
}

func LogoutHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	session, err := GetAuthSession(r)
	if err != nil {
		return util.StatusError{500, err}
	}
	uid, ok := session.Values["uid"]
	if !ok {
		return util.StatusError{500, err}
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		return util.StatusError{500, err}
	}

	domain := os.Getenv("AUTH_PROVIDER_DOMAIN")

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		return util.StatusError{500, err}
	}

	logoutUrl.Path += "v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host + "/home")
	fmt.Printf("returnTo: %v\n", returnTo)

	if err != nil {
		return util.StatusError{500, err}
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("POKER_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()
	//enableCors(&w)

	id := uid.(int)
	delete(LoggedInUsers, id)

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)

	return nil

	//	return nil
}

//func AuthCallbackHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
func CallbackHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	session, err := GetAuthSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	fmt.Printf("CallbackHandler 1\n")

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return err
	}

	fmt.Printf("CallbackHandler 2\n")
	authenticator, err := NewAuthenticator()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	fmt.Printf("CallbackHandler 3\n")
	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		fmt.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	fmt.Printf("CallbackHandler 4\n")
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return errors.New("no id_token field")
	}

	clientID := os.Getenv("AUTH0_CLIENT_ID")
	fmt.Printf("CallbackHandler 4\n")
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	// Getting now the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	fmt.Printf("CallbackHandler 5\n")


	var user *domain.User
	var account *domain.Account

	email := profile["email"]
	ctx := r.Context()
	user, err = domain.GblUserRepository.LoadUserByEmail(ctx, email.(string))

	fmt.Printf("user = %v\n", user)

	if err != nil || user == nil {
		//user = &User{Subject: subject, Name: profile["name"].(string), Provider: issuer}
		ctx := context.Background()
		account, err = domain.GblAccountRepository.LoadAccountByName(ctx, email.(string))
		if err != nil || account == nil {
			ctx := context.Background()
			account, err = domain.GblAccountRepository.CreateAccount(ctx, email.(string))
			if err != nil {
				return util.StatusError{http.StatusInternalServerError, err}
			}
		}


		user, err = domain.GblUserRepository.AddUserFromProfile(ctx, account.Id, profile)

		if err != nil {
			return util.StatusError{http.StatusInternalServerError, err}
		}
	}
	id := user.Id
	_, ok = LoggedInUsers[id]
	if ok {
		// user already logged in
		msg := fmt.Sprintf("User with id %d already logged in", id)
		error := errors.New(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, error}
	}


	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	session.Values["uid"] = user.Id
	session.Values["accountId"] = user.AccountId
	session.Values["given_name"] = profile["name"].(string)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	return nil
}

/*

func AuthCallbackHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("Logged in users: %v\n", LoggedInUsers)
	var (
		err error
		//token *oauth2.Token
	)

	//enableCors(&w)

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != ExampleAppState {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return util.StatusError{http.StatusBadRequest, err}
	}

	ctx := oidc.ClientContext(r.Context(), HttpClient)

	oauth2Token, err := Authenticator.Config.Exchange(ctx, r.URL.Query().Get("code"))

	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, err}
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, err}
	}

	fmt.Printf("Verifying Token\n")
	idToken, err := Authenticator.IdTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, err}
	}

	var claims map[string]interface{}

	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to read claims: "+err.Error(), http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, err}

		// handle error
	}
	email := claims["email"]

	fmt.Printf("claims = %v\n", claims)
	fmt.Printf("email = %v\n", email)
	var user *domain.User
	var account *domain.Account
	ctx = r.Context()
	user, err = domain.GblUserRepository.LoadUserByEmail(ctx, email.(string))

	fmt.Printf("user = %v\n", user)

	if err != nil || user == nil {
		//user = &User{Subject: subject, Name: profile["name"].(string), Provider: issuer}
		ctx := context.Background()
		account, err = domain.GblAccountRepository.LoadAccountByName(ctx, email.(string))
		if err != nil || account == nil {
			ctx := context.Background()
			account, err = domain.GblAccountRepository.CreateAccount(ctx, email.(string))
			if err != nil {
				return util.StatusError{http.StatusInternalServerError, err}
			}
		}

		user, err = domain.GblUserRepository.AddUserFromProfile(ctx, account.Id, claims)

		if err != nil {
			return util.StatusError{http.StatusInternalServerError, err}
		}
	}
	id := user.Id
	_, ok = LoggedInUsers[id]
	if ok {
		// user already logged in
		msg := fmt.Sprintf("User with id %d already logged in", id)
		error := errors.New(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, error}
	}

	session, err := GetAuthSession(r)

	session.Values["rawIDToken"] = rawIDToken
	session.Values["uid"] = user.Id
	session.Values["accountId"] = user.AccountId
	session.Values["given_name"] = claims["name"].(string)

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return util.StatusError{http.StatusInternalServerError, err}
	}

	err = session.Save(r, w)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	LoggedInUsers[user.Id] = true

	//enableCors(&w)
	//enableRCors(r)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	//	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
	return nil

}
*/

func enableRCors(r *http.Request) {
	r.Header["Access-Control-Allow-Origin"] = []string{"*"}

}

var HttpClient *http.Client

func init() {

	if false {
		var err error
		HttpClient, err = HttpClientForRootCAs("test")
		if err != nil {
			fmt.Printf("Error setting up HttpClient: %v\n", err)
		}
	}
}

func HttpClientForRootCAs(rootCAs string) (*http.Client, error) {
	tlsConfig := tls.Config{RootCAs: x509.NewCertPool()}
	rootCABytes, err := ioutil.ReadFile(rootCAs)
	if err != nil {
		return nil, fmt.Errorf("failed to read root-ca: %v", err)
	}
	if !tlsConfig.RootCAs.AppendCertsFromPEM(rootCABytes) {
		return nil, fmt.Errorf("no certs found in root CA file %q", rootCAs)
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}, nil
}
