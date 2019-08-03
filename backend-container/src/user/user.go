package user

import (
	"app/database"
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id        string `json:"id" gorm:"unique"`
	Username  string `json:"username" gorm:"size:255" valid:"required,length(3|20)" `
	Firstname string `json:"firstname" gorm:"size:255" valid:"required,length(2|20)" `
	Lastname  string `json:"lastname" gorm:"size:255" valid:"required,length(2|20)" `
	creds     []Credential
}

type Credential struct {
	gorm.Model
	AAGUID     []byte `gorm:"size:255"`
	Details    []byte `gorm:"size:2048"`
	SignCount  uint32
	FKUsername string
}

//Add a user to database
func NewUser(user User) (string, error) {
	var existuser User
	if database.DBCon.Where("username = ?", user.Username).First(&existuser).RecordNotFound() {
		user.Id = randomID()
		errors := database.DBCon.Create(&user).GetErrors()
		for _, err := range errors {
			return "", err
		}
		return user.Id, nil
	}
	return "", errors.New("Der Benutzername " + user.Username + "  ist bereits vergeben.")
}

//Find user by username
func FindUser(username string) (User, error) {
	var user User
	if result := database.DBCon.Where("username = ?", username).First(&user); result.Error != nil {
		return user, errors.New("Der Benutzer existiert nicht.")
	}
	return user, nil
}

//Find Credential by AAGUID
func FindCred(aaguid []byte) (Credential, error) {
	var credential Credential

	if result := database.DBCon.Where("aa_guid = ?", aaguid).First(&credential); result.Error != nil {
		return credential, errors.New("credential does not exist.")
	}
	return credential, nil
}

//Update Counter in database
func UpdateCred(aaguid []byte, counter uint32) {
	var credential Credential
	database.DBCon.Model(&credential).Where("aa_guid = ?", aaguid).Update("sign_count", counter)

}

//Generate an ID
func randomID() string {
	buf := make([]byte, 8)
	rand.Read(buf)
	return string(buf[:])
}

/*
The following functions implement WebAuthn interface
*/

// get WebAuthnID
func (u User) WebAuthnID() []byte {
	return []byte(u.Id)
}

// get WebAuthnName
func (u User) WebAuthnName() string {
	return u.Username
}

// get WebAuthnDisplayName
func (u User) WebAuthnDisplayName() string {
	return u.Firstname + " " + u.Lastname
}

// get WebAuthnIcon
func (u User) WebAuthnIcon() string {
	return ""
}

// set AddCredential
func (u *User) AddCredential(credInfo webauthn.Credential) {

	credJson, err := json.Marshal(credInfo)
	if err != nil {
	}

	databaseCred := Credential{
		Details:    credJson,
		FKUsername: u.Username,
		AAGUID:     credInfo.Authenticator.AAGUID,
		SignCount:  credInfo.Authenticator.SignCount,
	}
	database.DBCon.Save(&databaseCred)
}

// WebAuthnCredentials
func (u User) WebAuthnCredentials() []webauthn.Credential {
	var creds []Credential
	credentialList := []webauthn.Credential{}

	database.DBCon.Where("fk_username = ?", u.Username).Find(&creds)

	for _, cred := range creds {
		oneCred := webauthn.Credential{}
		json.Unmarshal(cred.Details, &oneCred)
		credentialList = append(credentialList, oneCred)
	}

	return credentialList
}

// CredentialExcludeList
//This list is used to prevent users from registering the same authenticator twice
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	var databaseCreds []Credential

	database.DBCon.Where("fk_username = ?", u.Username).Find(&databaseCreds)

	for _, cred := range databaseCreds {
		oneCred := webauthn.Credential{}
		json.Unmarshal(cred.Details, &oneCred)
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: oneCred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}
