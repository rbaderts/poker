package ports

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rbaderts/poker/app"
	"github.com/rbaderts/poker/domain"
	"github.com/rbaderts/poker/util"
	"html/template"
	"time"

	_ "mime"

	"github.com/go-chi/chi"
	_ "github.com/go-chi/cors"
	_ "github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	_ "github.com/oceanicdev/chi-param"
	"io/ioutil"
	_ "os"
	"strconv"

	"fmt"
	//_	"github.com/gobuffalo/packr"
	"github.com/rbaderts/poker/auth"
	_ "github.com/rbaderts/pokerlib"
	//"github.com/gorilla/sessions"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}

	//CentClient    *gocent.Client
	tableTemplate *template.Template
)

const (
	TIME_FORMAT = "02/06/2002 3:04PM"
)

//var Assets http.FileSystem = http.Dir("assets")

var fmap = template.FuncMap{
	"FormatAsDate": FormatAsDate,
	"eq": func(a, b interface{}) bool {
		return a == b
	},
}

func FormatAsDate(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func GetTableTemplate() *template.Template {
	if tableTemplate == nil {
		tableTemplate = template.New("table").Funcs(fmap)
		tableTemplate = tableTemplate.Delims("[[", "]]")
		_, err := tableTemplate.ParseFiles(
			"web/table.html",
			"web/nav.html")

		if err != nil {
			fmt.Errorf("load template error: %v\n", err)
		}
	}
	return tableTemplate
}

func GetTableId(r *http.Request) (int64, error) {
	if tId := chi.URLParam(r, "tableID"); tId != "" {
		id, err := strconv.ParseInt(tId, 10, 64)
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

func AddRestRoutes(r *chi.Router, env *util.Env) {

	//mime.AddExtensionType(".vue", "application/javascript")

	(*r).Route("/tables/{tableID}", func(r chi.Router) {
		r.Use(TableCtx)
		r.Get("/", util.Handler{env, TableRenderHandler}.ServeHTTP)
		//			r.Get("/ws", Handler{env, WebsocketHandler}.ServeHTTP)
	})

	(*r).Route("/api", func(r chi.Router) {

		r.Route("/users/{userID}", func(r chi.Router) {
			r.Get("/", util.Handler{env, GetUserHandler}.ServeHTTP)
			r.Post("/stack/{amount:[0-9]+}", util.Handler{env, AddMoneyHandler}.ServeHTTP)

		})
		//			r.Get("/players/{playerID}", Handler{env, GetPlayerHandler}.ServeHTTP)
		//			r.Delete("/players/{playerID}", Handler{env, DeletePlayerHandler}.ServeHTTP)

		r.Route("/table/{tableID}", func(r chi.Router) {
			r.Get("/", util.Handler{env, GetTableHandler}.ServeHTTP)
			r.Route("/activeGame", func(r chi.Router) {
				r.Get("/", util.Handler{env, GetTableHandler}.ServeHTTP)
			})
			r.Route("/game/{gameId}", func(r chi.Router) {
			})

		})

		r.Route("/tables", func(r chi.Router) {

			r.Get("/", util.Handler{env, GetTablesHandler}.ServeHTTP)
			r.Post("/", util.Handler{env, CreateTableHandler}.ServeHTTP)
			r.Route("/{tableID}", func(r chi.Router) {
				r.Use(TableCtx)
				r.Get("/", util.Handler{env, GetTableHandler}.ServeHTTP)
				r.Post("/bet", util.Handler{env, PlayerTurnHandler}.ServeHTTP)
				r.Get("/join", util.Handler{env, JoinTableHandler}.ServeHTTP)
				r.Post("/startgame", util.Handler{env, StartGameHandler}.ServeHTTP)

				r.Route("/seat/{seatNum}", func(r chi.Router) {
					r.Get("/", util.Handler{env, GetTableSeatHandler}.ServeHTTP)
					r.Post("/message", util.Handler{env, PostMessageHandler}.ServeHTTP)
					r.Post("/unjoin", util.Handler{env, LeaveTableHandler}.ServeHTTP)
				})
			})
		})

	})

	/*
		fmt.Printf("launching server on 3000\n")

		if err := http.ListenAndServe(":3000", r); err != nil {
			fmt.Printf("ListenAndServe error = %v\n", err)

		}

	*/
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
	fmt.Printf("Render table for user: %v with seat: %v\n", uid, usersSeat)
	data := struct {
		UserName   string
		TableID    int
		OurSeatNum int
	}{
		user.GivenName,
		table.Id,
		usersSeat,
	}

	if err := GetTableTemplate().Execute(w, data); err != nil {
		return util.StatusError{500, err}
	}
	return nil
}

func GetUserHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	userIdStr := chi.URLParam(r, "userID")
	var targetUid int = 0
	if userIdStr == "current" {
		uid := r.Context().Value("uid").(int)
		targetUid = uid
	} else {
		targetUid, _ = strconv.Atoi(userIdStr)
	}

	fmt.Printf("targetUid = %v\n", targetUid)
	user, err := domain.GblUserRepository.LoadUserById(r.Context(), targetUid)
	fmt.Printf("user = %v\n", user)

	if err != nil {
		return util.StatusError{500, err}
	}
	render.JSON(w, r, user)

	return nil
}

func PlayerTurnHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	//	table := r.Context().Value("table").(*domain.Table)

	var msg domain.Message
	responseStr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	err = json.Unmarshal(responseStr, &msg)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	err = domain.GblMessageDispatcher.RespondToRequest(msg.ResponseTo, &msg)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	return nil
}

func StartGameHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	table := r.Context().Value("table").(*domain.Table)
	table.SetupGame()

	//table.SetupGame()

	return nil
}

func AddMoneyHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	var userId int
	var amount int
	userIdStr := chi.URLParam(r, "userID")

	if userIdStr == "current" {
		uid := r.Context().Value("uid").(int)
		userId = uid
	} else {
		userId, _ = strconv.Atoi(userIdStr)
	}

	//	user, err := env.UserService.LoadUserById( userId)

	amountStr := chi.URLParam(r, "amount")
	amount, _ = strconv.Atoi(amountStr)
	err := env.UserService.AddMoney(userId, amount)

	if err != nil {
		return util.StatusError{500, err}
	}

	return nil

}

func PostMessageHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	table := r.Context().Value("table").(*domain.Table)
	_ = table
	//userId := r.Context().Value("uid").(int)
	//fmt.Printf("userId = %d\n", userId)

	seatNumStr := chi.URLParam(r, "seatNum")
	seatNum := 0
	var err error
	if seatNum, err = strconv.Atoi(seatNumStr); err != nil {
		return util.StatusError{http.StatusBadRequest, err}
	}

	_ = seatNum
	/*userId := 0
	user := table.Players[seatNum].User
	if user != nil {
		userId = user.Id
	}*/
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}
	fmt.Printf("message: %v\n", message)

	//	table.SendPlayerMessage(seatNum, string(message))

	return nil

}

func GetTablesHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	//userId := r.Context().Value("uid").(int)

	var tRecords []*domain.Table

	var err error
	tRecords, err = domain.GblTableRepository.GetAllTables(r.Context())

	if err != nil {
		return err
	}
	//bytes, err := json.Marshal(tRecords)
	//if err != nil {
	///	return util.StatusError{http.StatusInternalServerError, err}
	//}
	render.JSON(w, r, tRecords)

	return nil

}

func GetTableHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	table := r.Context().Value("table").(*domain.Table)

	tableResource := domain.CalculateUpdate(table)
	render.JSON(w, r, tableResource)
	return nil
}

func GetTableSeatHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {

	table := r.Context().Value("table").(*domain.Table)
	seatNumStr := chi.URLParam(r, "seatNum")
	seatNum := 0
	var err error
	if seatNum, err = strconv.Atoi(seatNumStr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	seat := table.GetSeat(seatNum)
	render.JSON(w, r, seat)
	return nil
}

func CreateTableHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	//	accountId := r.Context().Value("accountId").(int)
	//	userId := r.Context().Value("uid").(int)

	deckfile := chi.URLParam(r, "deckfile")

	var table *domain.Table

	if deckfile != "" {
		//		table = MainApp.NewTableWithFixedDec(deckfile, 6)
	} else {
		table = app.GblTableService.NewTable(6)
	}

	//t := NewTable(env.DB, 6)

	render.JSON(w, r, table)

	return nil
}

func JoinTableHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value("uid").(int)
	table := r.Context().Value("table").(*domain.Table)
	sessionId, err := auth.GetJWTToken(r)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	ctx := context.Background()
	user, err := domain.GblUserRepository.LoadUserById(ctx, userId)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	_, err = table.AssignSeat(ctx, user, user.Stack, sessionId)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	//	client := GetClient(sessionId)
	//	user, err := LoadUserById(env.DB, userId)

	token := auth.GenerateJWTToken(user.Subject, user.Id)

	resp := struct {
		table *domain.Table `json:"table"`
		JWT   string        `json:"jwt"`
	}{
		table,
		token,
	}

	render.JSON(w, r, resp)

	return nil
}

func LeaveTableHandler(env *util.Env, w http.ResponseWriter, r *http.Request) error {
	table := r.Context().Value("table").(*domain.Table)
	sessionId, err := auth.GetJWTToken(r)
	if err != nil {
		return util.StatusError{http.StatusInternalServerError, err}
	}

	seatNumStr := chi.URLParam(r, "seatNum")
	seatNum := 0
	if seatNum, err = strconv.Atoi(seatNumStr); err != nil {
		return util.StatusError{http.StatusBadRequest, err}
	}

	/*
		user, err := LoadUserById(env.DB, userId)
		if err != nil {
			return util.StatusError{http.StatusInternalServerError, err}
		}
	*/

	table.RemovePlayer(seatNum, sessionId)

	return nil
}

func TableCtx(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var t *domain.Table

		if tId := chi.URLParam(r, "tableID"); tId != "" {
			//			id, err = strconv.Atoi(tId)
			//t, err = GetTableById(DB, id)
			tableId, err := strconv.Atoi(tId)
			if err != nil {
				fmt.Errorf("Error with tableId: %v\n", tId)
				render.Render(w, r, util.ErrNotFound)
			}

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

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
/*
type Handler struct {
	*util.Env
	H func(e *util.Env, w http.ResponseWriter, r *http.Request) error
}
*/

// ServeHTTP allows our Handler type to satisfy http.Handler.
/*
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Env, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			//http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

*/

/*
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
*/
