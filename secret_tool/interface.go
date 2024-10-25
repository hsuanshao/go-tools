package secrettool

import "github.com/hsuanshao/go-tools/ctx"

// Utility defines functions for help process encrypted data
type Utility interface {
	// Encrypt handles encrypt raw message and provide result and public key for decrypt
	Encrypt(ctx ctx.CTX, message string) (encryptedMessage, publicKey string, err error)
	// Decrypt for get raw message
	Decrypt(ctx ctx.CTX, encryptedMessage, publicKey string) (message string, err error)
}
