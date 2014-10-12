package model

import (
	"bytes"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	AuthTokenTable = "authtokens"
	AuthRequestLen = 5
	TokenLen       = 32
)

var TokenDigitMult = int(math.Pow(10, float64(AuthRequestLen)))
var TokenBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var TokenBytesLen = len(TokenBytes)

type AuthToken struct {
	Id        string `gorethink:"id,omitempty"`
	UserId    string
	Token     string
	CreatedAt time.Time
}

func GenerateAuthRequestToken() string {
	return strconv.Itoa(
		(rand.New(rand.NewSource(time.Now().UnixNano())).Int() %
			(9 * TokenDigitMult)) + TokenDigitMult)
}

func GenerateAuthToken() (string, error) {
	var token string
	ra := rand.New(rand.NewSource(time.Now().UnixNano()))
	session, err := Connect()
	if err != nil {
		return "", err
	}
	defer session.Close()
	for {
		buf := bytes.NewBuffer([]byte{})
		for i := 0; i < TokenLen; i++ {
			buf.WriteByte(TokenBytes[ra.Int()%TokenBytesLen])
		}
		token = buf.String()
		cur, err := Db().Table(AuthTokenTable).GetAllByIndex("Token", token).Run(session)
		if err != nil {
			return "", err
		}
		if cur.IsNil() {
			break
		}
		if err := cur.Close(); err != nil {
			return "", err
		}
	}
	return token, nil
}

func FindAuthToken(token string) (*AuthToken, bool, error) {
	session, err := Connect()
	if err != nil {
		return nil, false, err
	}
	defer session.Close()
	cur, err := Db().Table(AuthTokenTable).GetAllByIndex("Token", token).
		Limit(1).Run(session)
	if err != nil || cur.IsNil() {
		return nil, false, err
	}
	t := &AuthToken{}
	err = cur.One(t)
	return t, err == nil, err
}

func NewAuthToken(userId string) (*AuthToken, error) {
	token, err := GenerateAuthToken()
	if err != nil {
		return nil, err
	}
	return &AuthToken{
		UserId:    userId,
		Token:     token,
		CreatedAt: time.Now(),
	}, nil
}

func (token *AuthToken) Save() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	if token.Id == "" {
		w, err := Db().Table(AuthTokenTable).Insert(token).RunWrite(session)
		if err != nil {
			return err
		}
		token.Id = w.GeneratedKeys[0]
		return nil
	}
	_, err = Db().Table(AuthTokenTable).Get(token.Id).Update(token).RunWrite(session)
	return err
}
