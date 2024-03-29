package bls

import (
	"crypto/rand"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/llifezou/hdwallet"
)

func BenchmarkBLSSign(b *testing.B) {
	signer := blsSigner{}
	mnemonic, err := hdwallet.NewMnemonic()
	if err != nil {
		b.Fatal(err)
	}
	seed, err := hdwallet.GenerateSeedFromMnemonic(mnemonic, "")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pk, _ := signer.GenPrivate(seed)
		randMsg := make([]byte, 32)
		_, _ = rand.Read(randMsg)
		b.StartTimer()

		_, _ = signer.Sign(pk, randMsg)
	}
}

func BenchmarkBLSVerify(b *testing.B) {
	signer := blsSigner{}
	mnemonic, err := hdwallet.NewMnemonic()
	if err != nil {
		b.Fatal(err)
	}
	seed, err := hdwallet.GenerateSeedFromMnemonic(mnemonic, "")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		randMsg := make([]byte, 32)
		_, _ = rand.Read(randMsg)

		priv, _ := signer.GenPrivate(seed)
		pk, _ := signer.ToPublic(priv)
		addr, _ := address.NewBLSAddress(pk)
		sig, _ := signer.Sign(priv, randMsg)

		b.StartTimer()

		_ = signer.Verify(sig, addr, randMsg)
	}
}
