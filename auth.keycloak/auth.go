package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	SessionKey    = ""
	Authenticator *OIDCAuthenticator
	Oauth2Config  oauth2.Config
)

type OIDCAuthenticator struct {
	Provider        *oidc.Provider
	Config          oauth2.Config
	IdTokenVerifier *oidc.IDTokenVerifier
}

const ExampleAppState = "I wish to wash my irish wristwatch"

func init() {
	gob.Register(map[string]interface{}{})
	Authenticator, _ = NewAuthenticator()
}

func NewAuthenticator() (*OIDCAuthenticator, error) {

	ctx := context.Background()

	brokerClientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	brokerClientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	authURL := os.Getenv("AUTH_PROVIDER_DOMAIN")
	redirectURL := os.Getenv("POKER_CALLBACK_URL")

	providerURL := authURL

	// Initialize a provider by specifying dex's issuer URL.
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		// handle error
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	// Configure the OAuth2 config with the client values.
	Oauth2Config = oauth2.Config{
		// client_id and client_secret of the client.
		ClientID:     brokerClientId,
		ClientSecret: brokerClientSecret,

		// The redirectURL.
		RedirectURL: redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		//
		// Other scopes, such as "groups" can be requested.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	tokenVerifier := provider.Verifier(&oidc.Config{ClientID: brokerClientId})

	return &OIDCAuthenticator{
		Provider:        provider,
		Config:          Oauth2Config,
		IdTokenVerifier: tokenVerifier,
		//	Ctx:      ctx,
	}, nil

}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthToken ...
// This is what is retured to the user
type AuthToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

type AuthTokenClaim struct {
	*jwt.StandardClaims
	UserId int `json:"userId"`
}

func GenerateJWTToken(sub string, userId int) string {
	expiresAt := time.Now().Add(time.Hour * 24 * 1).Unix()

	//	token := jwt.New(jwt.SigningMethodHS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    sub,
		"exp":    expiresAt,
		"userId": userId,
	})

	secret := os.Getenv("CENTRIFUGO_HMAC_SECRET")
	tokenString, error := token.SignedString([]byte(secret))
	if error != nil {
		fmt.Println(error)
	}

	return tokenString
}

func GetJWTToken(r *http.Request) (string, error) {

	//remoteAddr := r.RemoteAddr;
	session, err := GetAuthSession(r)
	if err != nil {
		fmt.Errorf("Error getting session: %v\n", err)
		return "", err
	}

	var rawIDToken interface{}
	var ok bool
	if rawIDToken, ok = session.Values["rawIDToken"]; !ok {
		fmt.Errorf("Error getting raw ID Token\n")
		return "", err
	}
	return rawIDToken.(string), nil
}

var clients = make(map[string]int)
var clientCounter = 1

func GetCookieName(r *http.Request) string {
	return "PokerCookie"
}

func enableCorsR(r *http.Request) {
	r.Header["Access-Control-Allow-Origin"] = []string{"*"}

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetAuthSession(r *http.Request) (*sessions.Session, error) {
	key := os.Getenv("POKER_SESSION_KEY")
	AuthStore := sessions.NewCookieStore([]byte(key))

	session, err := AuthStore.Get(r, GetCookieName(r))
	if err != nil {
		return nil, err
	}
	return session, nil

}

var StoreLock sync.Mutex

func VerifyAuth(next http.Handler) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		session, err := GetAuthSession(r)

		if err != nil {
			fmt.Printf("Redirecting to auth url1\n")
			http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		jwt, err := GetJWTToken(r)
		if err != nil {
			fmt.Printf("Redirecting to auth url1\n")
			http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", "123")
		_, err = Authenticator.IdTokenVerifier.Verify(ctx, jwt)

		if err != nil {
			fmt.Printf("Redirecting to auth url4\n")
			http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		uid, ok := session.Values["uid"]

		if uid == nil || !ok {
			fmt.Printf("Redirecting to auth url5\n")
			http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		var accountId interface{}
		accountId, ok = session.Values["accountId"]

		if accountId == nil || !ok {
			fmt.Printf("Redirecting to auth url6\n")
			http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		ctx = context.WithValue(r.Context(), "uid", uid)
		ctx = context.WithValue(ctx, "accountId", accountId)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
	return http.HandlerFunc(h)
}

func RedirectToAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, Authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusTemporaryRedirect)
}

func newState() string {
	return "dummystate"
}

/*
func GetSessionToken(r *http.Request) (string, error) {

	session, err := GetAuthSession(r)

	if err != nil {
		return "", err
	}
	var rawIDToken interface{}
	var ok bool
	if rawIDToken, ok = session.Values["rawIDToken"]; !ok {
		return "", err
	}
	return rawIDToken.(string), nil
}

*/

func containsString(strings []string, checkFor string) bool {
	for _, s := range strings {
		if s == checkFor {
			return true
		}
	}
	return false
}
