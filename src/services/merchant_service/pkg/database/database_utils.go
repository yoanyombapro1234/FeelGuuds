package database

import (
	"context"
	"errors"
	"os"
	"time"

	core_database "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-database"
	core_logging "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-logging/json"
	core_tracing "github.com/yoanyombapro1234/FeelGuuds/src/libraries/core/core-tracing"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/gen/github.com/yoanyombapro1234/FeelGuuds/src/merchant_service/proto/merchant_service_proto_v1"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/constants"
	svcErrors "github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/errors"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/saga"
	"github.com/yoanyombapro1234/FeelGuuds/src/services/merchant_service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)



// ValidateAndHashPassword validates, hashes and salts a password
func (db *Db) ValidateAndHashPassword(password string) (string, error) {
	// check if confirmed password is not empty
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	//  hash and salt password
	hashedPassword, err := db.hashAndSalt([]byte(password))
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// hashAndSalt hashes and salts a password
func (db *Db) hashAndSalt(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords compares a hashed password and a plaintext password and returns
// a boolean stating wether they are equal or not
func (db *Db) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
