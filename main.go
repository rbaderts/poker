package main

import (
	//"fmt"
	//"github.com/gobuffalo/packr"

	"flag"
	"fmt"
	"github.com/rbaderts/poker/adapters"
	"github.com/rbaderts/poker/app"
	"github.com/rbaderts/poker/cmd/server"
	"github.com/rbaderts/poker/db"
	"github.com/rbaderts/poker/domain"
	_ "github.com/rbaderts/pokerlib"
	"log"
	"net/http"
	"os"
	_ "os"
	"runtime/pprof"
	_ "strings"
	//	"github.com/gobuffalo/packr"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//"io/ioutil"
	//"log"
	//"os"
	//"path/filepath"
)

var Client http.Client

func main() {
	/*
		var deck *pokerlib.Deck
		if len(os.Args) > 1 {
			argsWithoutProg := os.Args[1:]
			for _, arg := range argsWithoutProg {
				if strings.Contains(arg, "-deck=") {
					parts := strings.Split(arg, "=")
					deckFile := parts[1]
					fmt.Printf("deckFile = %s\n", deckFile)
					deck = pokerlib.ReadDeck(deckFile)

				}

			}
		}

	*/
	var gamescript = flag.String("gamescript", "", "a game automation file")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()
	if *cpuprofile != "" {
		fmt.Printf("CPU Profile = %s\n", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *gamescript != "" {

	}

	db := db.SetupDB()
	adapters.MigrateDB()

	application := app.NewApplication(db)

	recovered := domain.RecoverTables()

	for u, amt := range recovered {
		fmt.Printf("Recovered %d dollards for user: %d\n", amt, u)
	}

	table := app.GblTableService.NewTable(6)

	application.SetTable(table)
	application.DB = db

	/*
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		//signal.Notify registers the given channel to receive notifications of the specified signals.

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		//This goroutine executes a blocking receive for signals. When it gets one it’ll print it out and then notify the program that it can finish.
		go func() {
			sig := <-sigs
			_ = sig
			table.ReturnPlayerCash()
			os.Exit(-1)
			done <- true
		}()
		//The program will wait here until it gets the expected signal (as indicated by the goroutine above sending a value on done) and then exit.

	*/
	done := make(chan bool, 1)
	go server.Server(done)
	//	setupKeycloak()

	//fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}

/*
func setupKeycloak() {
	server := "http://localhost:3000"
	keycloakserver := "http://localhost:8080"
	auth.Init(keycloakserver, server,"", "logout", "Pokr")

}

*/
