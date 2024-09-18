package repo

import (
	"errors"
	"math/rand"
	"strings"

	jwt "github.com/golang-jwt/jwt"
)

var tokenSecret = []byte("MZeyamYJDjHYPFTy47OenqS9NdVu1nA4") // TODO: get from env

type Session struct {
	Repo
}

func randomWord(n int) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}
	return sb.String()
}

func (r *Session) create(user *tUserRec) (string, error) {
	r.sessionID++

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": r.sessionID,
		"sa": randomWord(11),
	})

	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}

	r.sessions[r.sessionID] = &tSession{
		id:   r.sessionID,
		user: user,
	}

	return tokenString, nil
}

func getIDFromToken(tokenString string) (uint64, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, errors.New("bad sign method")
		}
		return tokenSecret, nil
	}

	token, err := jwt.Parse(tokenString, hashSecretGetter)
	if err != nil || !token.Valid {
		return 0, errors.New("bad token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("no payload")
	}

	v, ok := payload["id"].(float64)
	if !ok {
		return 0, errors.New("no identifier in token")
	}

	return uint64(v), nil
}

func (r *Session) getByToken(token string) (*tSession, error) {

	id, err := getIDFromToken(token)
	if err != nil {
		return nil, err
	}

	rec := r.sessions[id]
	if rec == nil {
		return nil, errors.New("session not found")
	}

	return rec, nil
}

func (r *Session) delete(token string) error {

	id, err := getIDFromToken(token)
	if err != nil {
		return err
	}

	delete(r.sessions, id)

	return nil
}
