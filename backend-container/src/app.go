package main

import (
	"app/database"
	"app/events"
	"app/session"
	"app/settings"
	"app/user"
	"fmt"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/qor/validations"
	"log"
	"net/http"
	"os"
)

var webAuthn *webauthn.WebAuthn
var sessionStore *session.Store

func main() {

	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: settings.RPDisplayName,
		RPID:          settings.RPID,
	})

	if err != nil {
		log.Fatal("unable to create webauthn:", err)
	}

	sessionStore, err = session.NewStore([]byte("test"))
	if err != nil {
		log.Fatal("failed to create session store:", err)
	}

	//check cli command
	if len(os.Args) == 1 {
		fmt.Println("no command found.")
		os.Exit(1)
	}

	//connect to database and run migrations
	database.DBCon, err = gorm.Open("sqlite3", settings.SqliteStorage)
	if err != nil {
		panic("failed to connect database")
	}
	database.DBCon.AutoMigrate(user.User{}, user.Credential{}, events.Login{})
	validations.RegisterCallbacks(database.DBCon)

	switch command := os.Args[1]; command {
	case "serve":
		port := settings.ApiPort
		serveHttp(port)
	case "version":
		fmt.Println("1.0")
	}
}

func serveHttp(p string) {

	r := mux.NewRouter()
	r.HandleFunc("/register/begin/", RegisterBegin).Methods("POST")
	r.HandleFunc("/register/finish/{username}", RegisterFinish).Methods("POST")

	r.HandleFunc("/login/begin/", LoginBegin).Methods("POST")
	r.HandleFunc("/login/finish/{username}", LoginFinish).Methods("POST")

	fmt.Println("serving on port " + p)
	if err := http.ListenAndServe(":"+p, r); err != nil {
		panic(err)
	}
}
