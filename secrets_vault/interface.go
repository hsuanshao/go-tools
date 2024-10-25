package secretsvault

import "github.com/hsuanshao/go-tools/ctx"

// VaultHandler is the secret/credential or any sensitive, internal service setup handler
type VaultHandler interface {
	// SetSecret to set a new secret, which require input major service name, and sub service name, with credential with []byte data type
	SetSecret(ctx ctx.CTX)
}
