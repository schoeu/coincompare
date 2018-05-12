package utils


import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type UserInfo struct {
	OpenID    string `json:"openId"`
	UnionId   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

type WXBizDataCrypt struct {
	appID, sessionKey string
}

func NewWXBizDataCrypt(appID, sessionKey string) *WXBizDataCrypt {
	return &WXBizDataCrypt{
		appID:      appID,
		sessionKey: sessionKey,
	}
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

func (w *WXBizDataCrypt) Decrypt(encryptedData, iv string) (UserInfo, error) {
	var userInfo UserInfo
	aesKey, err := base64.StdEncoding.DecodeString(w.sessionKey)
	if err != nil {
		return userInfo, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return userInfo, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return userInfo, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return userInfo, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return userInfo, err
	}
	
	fmt.Println("cipherText", string(cipherText))
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return userInfo, err
	}
	if userInfo.Watermark.AppID != w.appID {
		return userInfo, ErrAppIDNotMatch
	}
	return userInfo, nil
}