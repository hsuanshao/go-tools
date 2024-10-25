package vaultmodels

// SecretInfo describe a secret information that storage in
type SecretInfo struct {
	Content     []byte
	UpdateEpoch int64
}
