package server

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/rbaderts/poker/adapters"
	"github.com/rbaderts/poker/app"
	"github.com/rbaderts/poker/domain"
	"github.com/rbaderts/poker/ports"
	"github.com/rbaderts/poker/util"

	//	"github.com/centrifugal/centrifuge-go"
	_ "github.com/centrifugal/gocent"
	"mime"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	_ "github.com/oceanicdev/chi-param"
	"io/ioutil"
	"os"
	"strconv"

	"strings"

	"fmt"
	//_	"github.com/gobuffalo/packr"
	"github.com/rbaderts/poker/auth"
	_ "github.com/rbaderts/pokerlib"
	//"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"time"
)

var (
	//	bracketLiveTemplate *template.Template
	homeTemplate  *template.Template
	tableTemplate *template.Template
	//controlTemplate *template.Template
	//	playersTemplate      *template.Template
	upgrader = websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}

	//CentClient *gocent.Client
	//NatsClient *nats.Conn
)

const (
	TIME_FORMAT = "02/06/2002 3:04PM"
)

func init() {
	fmt.Printf("Server.init\n")
}

var fmap = template.FuncMap{
	"FormatAsDate": FormatAsDate,
	"eq": func(a, b interface{}) bool {
		return a == b
	},
}

//var Assets http.FileSystem = http.Dir("assets")

func GetTableId(r *http.Request) (int, error) {
	if tId := chi.URLParam(r, "tableID"); tId != "" {
		id, err := strconv.Atoi(tId)
		//id, err := strconv.ParseInt(tId, 10, 64)
		if err != nil {
			return 0, util.StatusError{500, err}
		}
		return id, nil
	}
	return 0, util.StatusError{http.StatusBadRequest, errors.New("Unable to get tableID from request")}
}

func GetUserId(ctx context.Context) *int64 {
	userId, ok := ctx.Value("uid").(int64)
	fmt.Printf("userId = %d\n", userId)
	if !ok {
		// Log this issue
		return nil
	}
	return &userId
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	fmt.Printf("origin = %v\n", origin)
	return true
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	w.Write([]byte("foo"))
}
func Server(done chan bool) {

	/*
		centrifugoKey := os.Getenv("CENTRIFUGO_KEY")
		centrifugoService := os.Getenv("CENTRIFUGO_SERVICE")
		CentClient = gocent.New(gocent.Config{
			Addr: centrifugoService,
			Key:  centrifugoKey,
		})
	*/

	var err error
	fmt.Printf("Attempting nats connection...\n")
	natsUrl := os.Getenv("NATS_URL")
	adapters.NatsClient, err = nats.Connect(natsUrl)
	if err != nil {
		fmt.Printf("Unable to connect to nats: %v\n", err)
	}

	env := &util.Env{

		DB:   app.MainApp.DB,
		Port: os.Getenv("PORT"),
		Host: os.Getenv("HOST"),
		// We might also have a custom log.Logger, our
		// template instance, and a config struct as fields
		// in our Env struct.
	}

	mime.AddExtensionType(".vue", "application/javascript")

	loadTemplates()

	r := chi.NewRouter()
	//	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:  []string{"http://*"},
		AllowOriginFunc: AllowOriginFunc,

		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           3000, // Maximum value not ignored by any of major browsers
	}))
	FileServer(r, "/static", http.Dir("./web"))

	r.Get("/callback", util.Handler{env, auth.CallbackHandler}.ServeHTTP)
	r.Get("/login", util.Handler{env, auth.LoginHandler}.ServeHTTP)
	r.Get("/devlogin", util.Handler{env, auth.DevLogin}.ServeHTTP)
	r.Get("/logout", util.Handler{env, auth.LogoutHandler}.ServeHTTP)

	r.Group(func(r chi.Router) {
		r.Use(auth.VerifyAuth)
		r.Get("/", util.Handler{env, HomeRenderHandler}.ServeHTTP)
		r.Get("/home", util.Handler{env, HomeRenderHandler}.ServeHTTP)
		r.Get("/lobby", util.Handler{env, HomeRenderHandler}.ServeHTTP)
		//		r.Get("/ws", Handler{env, WebsocketHandler}.ServeHTTP)

		ports.AddRestRoutes(&r, env)
	})

	fmt.Printf("launching server on 3000\n")

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("ListenAndServe error = %v\n", err)

	}

	fmt.Printf("exiting Server\n")
	done <- true

}

func TableRenderHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("TableRenderHander\n")
	table := r.Context().Value("table").(*domain.Table)

	uid := r.Context().Value("uid").(int)
	ctx := context.Background()
	user, err := domain.GblUserRepository.LoadUserById(ctx, uid)
	if err != nil {
		return util.StatusError{500, err}
	}

	usersSeat := table.FindUsersSeat(uid)
	data := struct {
		UserName   string
		TableID    int
		OurSeatNum int
	}{
		user.GivenName,
		table.Id,
		usersSeat,
	}

	if err := tableTemplate.Execute(w, data); err != nil {
		return util.StatusError{500, err}
	}
	return nil
}

func HomeRenderHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("HomeRenderHandler\n")

	uid := r.Context().Value("uid").(int)
	user, err := app.GblUserService.LoadUserById(uid)
	if err != nil {
		return util.StatusError{500, err}
	}

	data := struct {
		UserName string
		TableID  string
		UserId   int
	}{
		user.GivenName,
		"",
		uid,
	}

	//url := "/static/home.html"
	//http.Redirect(w, r, url, http.StatusFound)

	if err := homeTemplate.Execute(w, data); err != nil {
		return util.StatusError{500, err}
	}
	return nil
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

}

func TableCtx(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t *domain.Table

		if tId := chi.URLParam(r, "tableID"); tId != "" {
			//			id, err = strconv.Atoi(tId)
			fmt.Printf("tId = %s\n", tId)
			tableId, err := strconv.Atoi(tId)
			if err != nil {
				fmt.Errorf("Error with tableId: %v\n", tId)
				render.Render(w, r, util.ErrNotFound)
			}

			//t, err = GetTableById(DB, id)
			t = app.GblTableService.GetTable(tableId)
			if t == nil {
				fmt.Errorf("table with Id = %s not found\v", tId)
				render.Render(w, r, util.ErrNotFound)
				return
			}
		} else {
			fmt.Errorf("table with Id = %s not found\v", tId)
			render.Render(w, r, util.ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "table", t)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FormatAsDate(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func loadTemplates() {

	//	controlTemplate = template.Must(template.New("main").ParseFiles(
	//		"./web/index.html", "./web/tmpl/nav.tmpl")).Funcs(fmap)
	fmt.Printf("Loading templates\n")

	homeTemplate = template.New("home").Funcs(fmap)
	homeTemplate = homeTemplate.Delims("[[", "]]")
	var err error
	_, err = homeTemplate.ParseFiles(
		"web/index.html",
		"web/nav.html")

	if err != nil {
		fmt.Errorf("load template error: %v\n", err)
	}

	tableTemplate = template.New("table").Funcs(fmap)
	tableTemplate = tableTemplate.Delims("[[", "]]")
	_, err = tableTemplate.ParseFiles(
		"web/table.html",
		"web/nav.html")
	if err != nil {
		fmt.Errorf("load template error: %v\n", err)
	}

}

type Username struct {
	Name string `json:"name"`
}

type Id struct {
	Id int64 `json:"id"`
}

func (this Username) String() string {
	return fmt.Sprintf(`{"name": "%s"}`, this.Name)
}

func uploadFile(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return util.StatusError{500, err}
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return util.StatusError{500, err}
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return util.StatusError{500, err}
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return nil
}
