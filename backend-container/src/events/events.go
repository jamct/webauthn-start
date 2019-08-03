package events

//Log Events to database. Not implemented at the moment.

import (
	/*
		"app/database"
		"encoding/json"
		"errors"
		"fmt"
	*/
	"github.com/jinzhu/gorm"
)

type Login struct {
	gorm.Model
	Username []byte `gorm:"size:255"`
}

func AddLogin(username string) {
	//TODO
}
