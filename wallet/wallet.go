package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/Hyojip/housecoin/utils"
	"io/fs"
	"math/big"
	"os"
)

const (
	filename = "house.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet
var fileLayer = osFileLayer{}

func Wallet() *wallet {
	if w != nil {
		return w
	}

	w = loadWallet()
	return w
}

func Sign(w *wallet, payload string) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleError(err)

	return toStringFrom(r, s)
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleError(err)
	publicKey := parsePublicKey(address)
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	return ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
}

func parsePublicKey(address string) ecdsa.PublicKey {
	x, y, err := restoreBigInts(address)
	utils.HandleError(err)
	return ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
}

func restoreBigInts(signature string) (*big.Int, *big.Int, error) {
	sigAsBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}
	halfIdx := len(sigAsBytes) / 2
	firstOfBytes := sigAsBytes[:halfIdx]
	remainOfBytes := sigAsBytes[halfIdx:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstOfBytes)
	bigB.SetBytes(remainOfBytes)
	return &bigA, &bigB, nil
}

func loadWallet() *wallet {
	w := &wallet{}
	if fileLayer.hasWalletFile() {
		w.privateKey = restoreKey()
	} else {
		key := createPrivateKey()
		persistKey(key)
		w.privateKey = key
	}

	w.Address = extractAddress(w)
	return w
}

func extractAddress(w *wallet) string {
	return toStringFrom(w.privateKey.X, w.privateKey.Y)
}

func toStringFrom(a, b *big.Int) string {
	bytes := append(a.Bytes(), b.Bytes()...)
	return fmt.Sprintf("%x", bytes)
}

func createPrivateKey() *ecdsa.PrivateKey {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return pk
}

func persistKey(key *ecdsa.PrivateKey) {
	keyAsBytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	err = fileLayer.WriteFile(filename, keyAsBytes, os.FileMode(0644))
	utils.HandleError(err)
}
func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := fileLayer.ReadFile(filename)
	utils.HandleError(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)
	return key
}

type walletFileAccessible interface {
	hasWalletFile() bool
	WriteFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type osFileLayer struct{}

func (s *osFileLayer) hasWalletFile() bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (s *osFileLayer) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (s *osFileLayer) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
