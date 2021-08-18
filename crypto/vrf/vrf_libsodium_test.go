// +build libsodium

package vrf

import (
	"bytes"
	"crypto/ed25519"
	r2ishiguro "github.com/r2ishiguro/vrf/go/vrf_ed25519"
	"github.com/stretchr/testify/require"
	"testing"

	coniks "github.com/coniks-sys/coniks-go/crypto/vrf"
	libsodium "github.com/line/ostracon/crypto/vrf/internal/vrf"
)

func TestKeyPairCompatibilityLibsodium(t *testing.T) {
	secret := [SEEDBYTES]byte{}
	publicKey, privateKey := libsodium.KeyPairFromSeed(&secret)

	privateKey2 := ed25519.PrivateKey(make([]byte, 64))
	copy(privateKey2, privateKey[:])
	publicKey2 := privateKey2.Public().(ed25519.PublicKey)

	if !bytes.Equal(publicKey[:], publicKey2[:]) {
		t.Error("public key is not matched: using same private key which is generated by libsodium",
			"libsodium.Public", enc(publicKey[:]), "ed25519.Public", enc(publicKey2[:]))
	}
}

func TestProveAndVerify_LibsodiumByCryptoED25519(t *testing.T) {
	secret := [SEEDBYTES]byte{}
	privateKey := ed25519.NewKeyFromSeed(secret[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	verified, err := proveAndVerify(t, privateKey, publicKey)
	//
	// verified when using crypto ED25519
	//
	require.Nil(t, err)
	require.True(t, verified)
}

func TestProveAndVerify_LibsodiumByConiksED25519(t *testing.T) {
	secret := [SEEDBYTES]byte{}
	privateKey, _ := coniks.GenerateKey(bytes.NewReader(secret[:]))
	publicKey, _ := privateKey.Public()

	verified, err := proveAndVerify(t, privateKey, publicKey)
	//
	// "un-verified" when using coniks ED25519
	// If you want to use libsodium, you should use crypto/libsodium ED25519
	//
	require.NotNil(t, err)
	require.False(t, verified)
}

func TestProveAndVerify_LibsodiumByLibsodiumED25519(t *testing.T) {
	secret := [SEEDBYTES]byte{}
	publicKey, privateKey := libsodium.KeyPairFromSeed(&secret)

	verified, err := proveAndVerify(t, privateKey[:], publicKey[:])
	//
	// verified when using libsodium ED25519
	//
	require.Nil(t, err)
	require.True(t, verified)
}

func TestProveAndVerifyCompatibilityLibsodium(t *testing.T) {
	secret := [SEEDBYTES]byte{}
	message := []byte("hello, world")
	privateKey := ed25519.NewKeyFromSeed(secret[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	libsodiumImpl := newVrfEd25519libsodium()

	{
		proof, err := libsodiumImpl.Prove(privateKey, message)
		require.Nil(t, err)
		require.NotNil(t, proof)

		output, err := r2ishiguro.ECVRF_verify(publicKey, proof, message)
		//
		// No compatibility between libsodium.Prove and r2ishiguro.Verify
		//
		require.NotNil(t, err)
		require.NotNil(t, output)
	}
	{
		proof, err := r2ishiguro.ECVRF_prove(publicKey, privateKey, message)
		require.Nil(t, err)
		require.NotNil(t, proof)

		output, err := libsodiumImpl.Verify(publicKey, proof, message)
		//
		// No compatibility between r2ishiguro.Prove and libsodium.Verify
		//
		require.NotNil(t, err)
		require.NotNil(t, output)
	}
}