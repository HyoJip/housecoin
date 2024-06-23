package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"strings"
	"testing"
)

const (
	pk      = "307702010104207bca7f55150530322cbdb0e44779614da6ff962f055ba5a7564adf8c39c13c0fa00a06082a8648ce3d030107a144034200040c5c936a37d5abbb66d34ce12d2c8924b4d59863757b86e3d01671675c9ccb45697f0fb829442fe25ea06ad2d2430daf2d6e242b91f66559748b5a7af93bed2a"
	payload = "002a0b6849c412926b4b4ce7efd0508f0a52500261836c0c20f4f2b792abf503"
	sign    = "9c11c2fad5954b9d32c2d86fbb7f045f62f4b2246ce0b2ec614c271f3ab44d81a5d298b31d3cf2462e32dab745abbc27e4c9a7dcac0763127e776d2528a2d8bf"
)

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{payload, true},
		{strings.Replace(payload, "a", "b", 1), false},
	}
	w := makeWallet()
	for _, tc := range tests {
		ok := Verify(sign, tc.input, w.Address)
		if ok != tc.ok {
			t.Errorf("Verify(%q, %q) = %v, want %v", tc.input, w.Address, ok, tc.ok)
		}
	}

	t.Run("invalid hex decoded string returns error", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Errorf("Verify(%q, %q, %q) dosen't throw panic", sign[1:], payload, w.Address)
			}
		}()
		Verify(sign[1:], payload, w.Address)
	})

}

func TestSign(t *testing.T) {
	s := Sign(makeWallet(), payload)
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return hex encoded string. got %s. error: %v", s, err)
	}
}

func makeWallet() *wallet {
	b, _ := hex.DecodeString(pk)
	key, _ := x509.ParseECPrivateKey(b)
	return &wallet{
		privateKey: key,
		Address:    toStringFrom(key.X, key.Y),
	}
}
