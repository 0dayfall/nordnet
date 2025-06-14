// Package util contains utility functions
package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

const Stockholm_Timezone = "Europe/Stockholm"

func SafeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("goroutine", name).
					Interface("panic", r).
					Msg("panic recovered")
			}
		}()
		fn()
	}()
}

func GenerateCredentials(username, password, rawPem []byte) (cred string, err error) {
	var loc *time.Location
	if loc, err = time.LoadLocation(Stockholm_Timezone); err != nil {
		return
	}
	ms := time.Now().In(loc).Unix() * 1000

	unixStr := strconv.FormatInt(ms, 10)

	userBase64 := base64.StdEncoding.EncodeToString(username)
	passBase64 := base64.StdEncoding.EncodeToString(password)
	timeBase64 := base64.StdEncoding.EncodeToString([]byte(unixStr))

	formated := fmt.Sprintf("%s:%s:%s", userBase64, passBase64, timeBase64)

	block, _ := pem.Decode(rawPem)
	if block == nil {
		err = errors.New("Could not Decode PEM")
		return
	}

	pubKeyVal, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}

	rsaPubKey, ok := pubKeyVal.(*rsa.PublicKey)
	if !ok {
		err = errors.New("Could not make DER to an RSA PublicKey")
		return
	}

	encr, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(formated))
	if err != nil {
		return
	}

	cred = base64.StdEncoding.EncodeToString(encr)
	return
}
