package auth

import (
	"context"
	"github.com/gorilla/sessions"
	"log"
	"os"

	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"net/http"
	"time"

	oidc "github.com/coreos/go-oidc"
)

//var (
//session, err = app.Store.Get(r, "auth-session")
//SessionKey = ""
//}

func init() {
	//gob.Register(map[string]interface{}{})
}

func GetCookieName(r *http.Request) string {
	return "PokerCookie"
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
type Authenticator struct {
	Provider      *oidc.Provider
	Config        oauth2.Config
	Ctx           context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	domain := os.Getenv("AUTH0_DOMAIN")
	fmt.Printf("domain = %v\n", domain)
	provider, err := oidc.NewProvider(ctx, domain)
//	fmt.Printf("provider = %v\n", provider)
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")

	conf := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:3000/callback",
		Endpoint:     provider.Endpoint(),
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Authenticator{
		Provider:      provider,
		Config:        conf,
		Ctx:           ctx,
	}, nil
}

func GetJWTToken(r *http.Request) (string, error) {

	//remoteAddr := r.RemoteAddr;
	//session, err := GetAuthSession(r)
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

func VerifyAuth(next http.Handler) http.Handler {

	h := func(w http.ResponseWriter, r *http.Request) {

		authenticator, err := NewAuthenticator()
		session, err := GetAuthSession(r)

		if err != nil {
			fmt.Printf("Redirecting to auth url1\n")
			http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}



		fmt.Printf("VerifyAuth1\n")
		if _, ok := session.Values["profile"]; !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			uid, ok := session.Values["uid"]
			if uid == nil || !ok {
			//	fmt.Printf("Redirecting to auth url5\n")
				//http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			ctx := context.WithValue(r.Context(), "uid", uid)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
		fmt.Printf("VerifyAuth2\n")

		/*
		jwt, err := GetJWTToken(r)
		if err != nil {
			fmt.Printf("Redirecting to auth url1\n")
			http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", "123")

		clientId := os.Getenv("AUTH0_CLIENT_ID")
		oidcConfig := &oidc.Config{
			ClientID: clientId,
		}

		_, err = authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(),jwt)


		if err != nil {
			fmt.Printf("Redirecting to auth url4\n")
			http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		uid, ok := session.Values["uid"]

		if uid == nil || !ok {
			fmt.Printf("Redirecting to auth url5\n")
			http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		var accountId interface{}
		accountId, ok = session.Values["accountId"]

		if accountId == nil || !ok {
			fmt.Printf("Redirecting to auth url6\n")
			http.Redirect(w, r, authenticator.Config.AuthCodeURL(ExampleAppState), http.StatusFound)
			return
		}

		ctx = context.WithValue(r.Context(), "uid", uid)
		ctx = context.WithValue(ctx, "accountId", accountId)
		next.ServeHTTP(w, r.WithContext(ctx))

		 */

	}
	return http.HandlerFunc(h)
}
