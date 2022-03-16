package secp

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
	crypto2 "github.com/filecoin-project/go-state-types/crypto"
	"github.com/llifezou/fil-sdk/sigs"
	"github.com/minio/blake2b-simd"
)

type secpSigner struct{}

func (secpSigner) GenPrivate(seed []byte) ([]byte, error) {
	seedSha256 := sha256.Sum256(seed)
	ecdsaSeed := append(seed, seedSha256[:]...)

	// ecdsaSeed require length of 40
	if len(ecdsaSeed) < 40 {
		return nil, errors.New("seed is too short")
	}

	priv, err := crypto.GenerateKeyFromSeed(bytes.NewReader(ecdsaSeed))

	if err != nil {
		return nil, err
	}
	return priv, nil
}

func (secpSigner) ToPublic(pk []byte) ([]byte, error) {
	return crypto.PublicKey(pk), nil
}

func (secpSigner) Sign(pk []byte, msg []byte) ([]byte, error) {
	b2sum := blake2b.Sum256(msg)
	sig, err := crypto.Sign(pk, b2sum[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (secpSigner) Verify(sig []byte, a address.Address, msg []byte) error {
	b2sum := blake2b.Sum256(msg)
	pubk, err := crypto.EcRecover(b2sum[:], sig)
	if err != nil {
		return err
	}

	maybeaddr, err := address.NewSecp256k1Address(pubk)
	if err != nil {
		return err
	}

	if a != maybeaddr {
		return fmt.Errorf("signature did not match")
	}

	return nil
}

func init() {
	sigs.RegisterSignature(crypto2.SigTypeSecp256k1, secpSigner{})
}
