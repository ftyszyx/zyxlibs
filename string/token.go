package string

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type TokenData struct {
	Expire int64  `json:"expire"`
	Value  string `json:"value"`
}

func NewToken(value string, expiretime int64) *TokenData {
	var data = new(TokenData)
	if expiretime > 0 {
		data.Expire = time.Now().Unix() + expiretime
	}

	data.Value = value
	return data
}

func (data *TokenData) GetToken(password string) (string, error) {
	text, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	res, err := AesEncrypt([]byte(text), []byte(password))
	if err != nil {
		return "", err
	}

	return string(res), err
}

func ParseToken(tokenstr string, key string) (*TokenData, error) {

	origData, err := AesDecrypt([]byte(tokenstr), []byte(key))
	if err != nil {
		return nil, err
	}

	var tokendata = new(TokenData)
	err = json.Unmarshal(origData, tokendata)
	if err != nil {
		return nil, err
	}

	if tokendata.Expire > 0 && tokendata.Expire < time.Now().Unix() {
		return nil, errors.New("token expire")
	}
	return tokendata, nil
}
