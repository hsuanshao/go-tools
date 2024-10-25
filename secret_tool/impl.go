package secrettool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"
	"github.com/hsuanshao/go-tools/randm"
)

func NewSecretTool() (u Utility) {
	randDice := randm.NewDice(5)
	return &impl{Dice: randDice}
}

type impl struct {
	Dice randm.Method
}

var (
	// ErrEmptyMessageForEncrypt ...
	ErrEmptyMessageForEncrypt = errors.New("input empty string for encrypt is not supported")
	// ErrNewCipherFailed ...
	ErrNewCipherFailed = errors.New("aes new cipher get error")

	// ErrNewGCM ...
	ErrNewGCM = errors.New("cipher based on c to new a gcm failed")

	// ErrAesGcmDecrypt ...
	ErrAesGcmDecrypt = errors.New("aes.gcm decrypt failed")

	// ErrHexDecodeRawMessage ....
	ErrHexDecodeRawMessage = errors.New("hex decode raw message failed")

	// ErrBase64StdDecodeStr ...
	ErrBase64StdDecodeStr = errors.New("base64 decode encrypt message back to origin but failed")

	// ErrRandReadIV ....
	ErrRandReadIV = errors.New("rand.read iv failed")
)

// Encrypt handles encrypt raw message and provide result and public key for decrypt
func (im *impl) Encrypt(ctx ctx.CTX, message string) (encryptedMessage, publicKey string, err error) {
	if message == "" || strings.TrimSpace(message) == "" {
		ctx.WithField("message", message).Error("input message for encrypt should not as empty string")
		return "", "", ErrEmptyMessageForEncrypt
	}
	pubKey := im.Dice.GenRandomString(32)

	c, err := aes.NewCipher([]byte(pubKey))
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "generate pk": pubKey}).Error("AES New Cipher failed")
		return "", "", ErrNewCipherFailed
	}

	blockSize := len(pubKey)
	paddingLen := blockSize - len(message)%blockSize

	message += string(bytes.Repeat([]byte{byte(paddingLen)}, paddingLen))
	cipherText := make([]byte, aes.BlockSize+len(message))
	iv := cipherText[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		ctx.WithField("err", err).Error("rand read iv get error")
		return "", "", ErrRandReadIV
	}

	m := cipher.NewCBCEncrypter(c, iv)
	m.CryptBlocks(cipherText[aes.BlockSize:], []byte(message))

	encrypted := base64.StdEncoding.EncodeToString(cipherText)

	return encrypted, pubKey, nil

}

// Decrypt for get raw message
func (im *impl) Decrypt(ctx ctx.CTX, encryptedMessage, publicKey string) (message string, err error) {
	c, err := aes.NewCipher([]byte(publicKey))
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "public key": publicKey}).Error("decrypt based on input public key but failed")
		return "", ErrNewCipherFailed
	}

	ciphercode, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		ctx.WithField("err", err).Error("decode base64 encrypt string failed")
		return "", ErrBase64StdDecodeStr
	}

	iv := ciphercode[:aes.BlockSize]
	ciphercode = ciphercode[aes.BlockSize:]

	m := cipher.NewCBCDecrypter(c, iv)
	m.CryptBlocks(ciphercode, ciphercode)

	plainText := string(ciphercode)
	res := plainText[:len(plainText)-int(plainText[len(plainText)-1])]

	return res, nil
}
