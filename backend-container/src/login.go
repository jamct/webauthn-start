package main

import (
	"app/settings"
	"app/user"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Struct for incomoing JSON data
type LoginUser struct {
	Username string `json:"username"`
}

//Struct for JWT response
type Claims struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	jwt.StandardClaims
}

//JWT response for successfull login
type LoginSuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

var jwtKey = []byte(settings.JWTKey)

func LoginBegin(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var loginUser LoginUser
	err := decoder.Decode(&loginUser)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		fmt.Println("")
		return
	}
	loginUser.Username = strings.ToLower(loginUser.Username)

	var databaseUser user.User

	databaseUser, err = user.FindUser(loginUser.Username)
	if err != nil {
		jsonResponse(w, "Der Benutzer existiert nicht.", http.StatusNotFound)
		return
	}

	// prepare session
	options, sessionData, err := webAuthn.BeginLogin(databaseUser)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store session data JSON
	err = sessionStore.SaveWebauthnSession("authentication", sessionData, r, w)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Login Challenge is ready. Send it to browser:
	jsonResponse(w, options, http.StatusOK)

}

func LoginFinish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := strings.ToLower(vars["username"])
	var databaseUser user.User
	databaseUser, err := user.FindUser(username)

	if err != nil {
		jsonResponse(w, "Der übergebene Benutzer existiert nicht!", http.StatusNotFound)
		return
	}

	// load the session data
	sessionData, err := sessionStore.GetWebauthnSession("authentication", r)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cred, err := webAuthn.FinishLogin(databaseUser, sessionData, r)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if cred.Authenticator.CloneWarning {
		jsonResponse(w, "Der Authenticator könnte manipuliert sein.", http.StatusForbidden)
		return
	}

	databaseCred, _ := user.FindCred(cred.Authenticator.AAGUID)

	if cred.Authenticator.SignCount < databaseCred.SignCount {
		jsonResponse(w, "Der Authenticator könnte manipuliert sein.", http.StatusForbidden)
		return
	}
	user.UpdateCred(cred.Authenticator.AAGUID, cred.Authenticator.SignCount)

	// Login was successful. Return JWT
	fmt.Println("successful login.")
	response := LoginSuccess{
		Message: "Login erfolgreich.",
		Token:   createJWT(databaseUser),
		Success: true,
	}

	jsonResponse(w, response, http.StatusOK)
}

func createJWT(user user.User) string {
	t, err := strconv.ParseUint(settings.JWTTime, 10, 32)
	expirationTime := time.Now().Add(time.Duration(t) * time.Minute)

	claims := &Claims{
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    settings.JWTIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err.Error()
	}
	return tokenString
}
