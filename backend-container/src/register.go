package main

import (
	"app/user"
	"encoding/json"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func RegisterBegin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser user.User
	err := decoder.Decode(&newUser)
	if err != nil {
		jsonResponse(w, "Die Registrierung war nicht vollständig.", http.StatusBadRequest)
		return
	}
	newUser.Username = strings.ToLower(newUser.Username)
	id, err := user.NewUser(newUser)
	newUser.Id = id

	if err != nil {
		jsonResponse(w, "Es existiert bereits ein Benutzer mit diesem Namen.", http.StatusBadRequest)
		return
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = newUser.CredentialExcludeList()
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := webAuthn.BeginRegistration(
		newUser,
		registerOptions,
	)

	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionStore.SaveWebauthnSession("registration", sessionData, r, w)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, options, http.StatusOK)
}

func RegisterFinish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := strings.ToLower(vars["username"])
	var databaseUser user.User
	databaseUser, err := user.FindUser(username)

	if err != nil {
		jsonResponse(w, "Der Benutzer existiert nicht!", http.StatusBadRequest)
		return
	}
	// read registration information from session
	sessionData, err := sessionStore.GetWebauthnSession("registration", r)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	credential, err := webAuthn.FinishRegistration(databaseUser, sessionData, r)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	databaseUser.AddCredential(*credential)
	jsonResponse(w, "Die Registrierung wurde abgeschlossen.", http.StatusOK)
}
