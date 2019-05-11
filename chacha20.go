// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-05-11 14:27:11 2A9A26                   go-experiments/[chacha20.go]
// -----------------------------------------------------------------------------

package main

import (
	"bytes"
	"fmt"
	"math/rand"

	chacha20 "golang.org/x/crypto/chacha20poly1305"
)

// chacha20EncryptionDemo demonstrates how you can use the
// ChaCha20-Poly1305 AEAD symmetric encryption algorithm
// (AEAD means Authenticated Encryption with Associated Data)
//
func chacha20EncryptionDemo() {
	info := func(label string, stringFormat bool, ar []byte) {
		var s string
		if stringFormat {
			s = fmt.Sprintf("%q", string(ar))
		} else {
			s = fmt.Sprintf("%#v", ar)
		}
		fmt.Printf("%s [%d] %s\n", label, len(ar), s)
	}
	for i, t := range chacha20Samples {
		{
			fmt.Println()
			info("plaintext:", true, t.plaintext)
			info("additionalData:", true, t.additionalData)
			info("key:", true, t.key)
			info("nonce:", true, t.nonce)
			info("plaintext:", true, t.plaintext)
			info("cyphertext:", false, t.cyphertext)
		}
		// create the encrypter / decrypter
		algo, err := chacha20.New(t.key)
		if err != nil {
			fmt.Printf("failed creating new key")
		}
		// Seal encrypts and authenticates plaintext, authenticates
		// the additional data and appends the result to dst,
		// returning the updated slice. The nonce must be NonceSize()
		// bytes long and unique for all time, for a given key.
		//
		// To reuse plaintext's storage for the encrypted output, use
		// plaintext[:0] as dst. Otherwise, the remaining capacity of
		// dst must not overlap plaintext.
		//
		// Seal(dst, nonce, plaintext, additionalData []byte) []byte
		//
		cyphertext := algo.Seal(
			nil,              // dst            []byte
			t.nonce,          // nonce          []byte
			t.plaintext,      // plaintext      []byte
			t.additionalData, // additionalData []byte
		) // [] byte
		if !bytes.Equal(cyphertext, t.cyphertext) {
			fmt.Printf("#%d: GOT %#v, WANT %#v", i, cyphertext, t.cyphertext)
			continue
		}
		// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
		// Decrypt:
		plaintext, err := algo.Open(nil, t.nonce, cyphertext, t.additionalData)
		if err != nil {
			fmt.Printf("#%d: Open() FAILED", i)
			continue
		}
		if !bytes.Equal(t.plaintext, plaintext) {
			fmt.Printf("#%d: PLAINTEXT DOES NOT MATCH: GOT: %x WANT: %x",
				i, plaintext, t.plaintext)
			continue
		}
		// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
		// Tamper with ciphertext, nonce and additonal data.
		// This should cause Open() to fail:
		//
		const Failed = "#%d: Open() FAILED TO ERROR, GIVEN TAMPERED %s."
		{
			i := rand.Intn(len(cyphertext))
			cyphertext[i] ^= 0x01
			_, err := algo.Open(nil, t.nonce, cyphertext, t.additionalData)
			if err == nil {
				fmt.Printf(Failed, i, "CIPHERTEXT")
			}
			cyphertext[i] ^= 0x01
		}
		{
			i := rand.Intn(algo.NonceSize())
			t.nonce[i] ^= 0x01
			_, err := algo.Open(nil, t.nonce, cyphertext, t.additionalData)
			if err == nil {
				fmt.Printf(Failed, i, "NONCE")
			}
			t.nonce[i] ^= 0x01
		}
		if len(t.additionalData) > 0 {
			i := rand.Intn(len(t.additionalData))
			t.additionalData[i] ^= 0x01
			_, err := algo.Open(nil, t.nonce, cyphertext, t.additionalData)
			if err == nil {
				fmt.Printf(Failed, 1, "ADDITIONAL DATA")
			}
			t.additionalData[i] ^= 0x01
		}
	}
} //                                                        chacha20EncryptionDemo

var chacha20Samples = []struct {
	plaintext, additionalData, key, nonce, cyphertext []byte
}{
	{
		[]byte("X"), //                                plaintext (12 bytes)
		[]byte(""),  //                                additionalData
		[]byte("a6d9x99u4B6c52t01fE150j3729Em4D9"), // 32-byte key
		[]byte("jabberwocky!"),                     // 12-byte nonce
		[]byte{ // cyphertext: 17 bytes
			0x41, 0x0F, 0x2D, 0x23, 0x45, 0x00, 0x75, 0x98, 0xE2, 0xCF,
			0x4E, 0x1C, 0x8D, 0x02, 0x72, 0xBC, 0xBD},
	},
	{
		[]byte("Hello world!"),                     // plaintext (12 bytes)
		[]byte("1234567"),                          // additionalData
		[]byte("012345678901234567890123456789AB"), // 32-byte key
		[]byte("0123456789AB"),                     // 12-byte nonce
		[]byte{ // cyphertext: 28 bytes
			0xF5, 0x8D, 0xAD, 0xA2, 0x11, 0x97, 0x86, 0xA6, 0x00, 0xCB,
			0xE7, 0x93, 0x14, 0x42, 0xBB, 0x10, 0xDF, 0x89, 0x64, 0xEA,
			0x7B, 0x03, 0x3C, 0x84, 0x46, 0xEC, 0x39, 0x94,
		},
	},
}

//end
