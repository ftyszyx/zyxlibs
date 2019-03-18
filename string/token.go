package string

import (
	"encoding/base32"
	"encoding/json"
	"time"

	"github.com/ftyszyx/libs/beego/logs"

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
		return "", errors.WithStack(err)
	}

	res, err := AesEncrypt(text, []byte(password))
	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(res), err
}

func ParseToken(tokenstr string, key string) (*TokenData, error) {

	aesstr, err := base32.StdEncoding.DecodeString(tokenstr)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	origData, err := AesDecrypt(aesstr, []byte(key))
	logs.Info("origin:%s", string(origData))
	if err != nil {
		return nil, err
	}

	var tokendata = new(TokenData)
	err = json.Unmarshal(origData, tokendata)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if tokendata.Expire > 0 && tokendata.Expire < time.Now().Unix() {
		return nil, errors.New("token expire")
	}
	return tokendata, nil
}
