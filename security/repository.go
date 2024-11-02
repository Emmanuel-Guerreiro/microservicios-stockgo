package security

import (
	"emmanuel-guerreiro/stockgo/lib"
	"emmanuel-guerreiro/stockgo/lib/log"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

func getRemoteToken(token string, ctx ...interface{}) (*User, error) {
	// Buscamos el usuario remoto
	req, err := http.NewRequest("GET", lib.GetEnv().SecurityServerURL+"/v1/users/current", nil)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, lib.UnauthorizedError
	}
	req.Header.Add("Authorization", "bearer "+token)
	if corrId, ok := log.Get(ctx...).Data()[log.LOG_FIELD_CORRELATION_ID].(string); ok {
		req.Header.Add(log.LOG_FIELD_CORRELATION_ID, corrId)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Get(ctx...).Error(err)
		return nil, lib.UnauthorizedError
	}
	defer resp.Body.Close()

	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}
	return user, nil
}
