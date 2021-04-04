// -----------------------------------------------------------------------------
// Go Language Experiments                          go-experiments/[aes_demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

// This file demonstrates how to encrypt and
// decrypt data using AES-256 symmetric cipher.
//
// Most of the comments are taken from the Go
// language standard library's documentation.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

var _ = encryptAES
var _ = decryptAES
var _ = aesDemo

// encryptAES encrypts plaintext using secretKey and returns
// the encrypted cipherthext, using AES-256 symmetric cipher.
func encryptAES(plaintext, secretKey []byte) (ciphertext []byte, err error) {
	var gcm cipher.AEAD
	{
		// NewCipher creates and returns a new cipher.Block
		// The key argument should be the AES key, either 16, 24,
		// or 32 bytes to select AES-128, AES-192, or AES-256.
		cip, err := aes.NewCipher(secretKey)
		if err != nil {
			fmt.Println(err)
		}
		// NewGCM returns the given 128-bit, block cipher wrapped
		// in Galois Counter Mode with the standard nonce length.
		//
		// In general, the GHASH operation performed by this implementation
		// of GCM is not constant-time. An exception is when the underlying
		// Block was created by aes.NewCipher on systems with hardware support
		// for AES. See the crypto/aes package documentation for details.
		//
		// See also https://en.wikipedia.org/wiki/Galois/Counter_Mode
		//
		// func cipher.NewGCM(cipher cipher.Block) (cipher.AEAD, error)
		gcm, err = cipher.NewGCM(cip)
		if err != nil {
			return nil, err
		}
	}
	// creates a byte array the size of the nonce which must be passed to Seal
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	//
	// fill the nonce with cryptographically secure random bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		// rand.Reader is a global, shared instance of a cryptographically
		// secure random number generator.
		//
		// On Linux and FreeBSD, Reader uses getrandom(2) if available,
		// /dev/urandom otherwise. On OpenBSD, Reader uses getentropy(2).
		// On other Unix-like systems, Reader reads from /dev/urandom.
		// On Windows systems, Reader uses the RtlGenRandom API.
		// On Wasm, Reader uses the Web Crypto API.
		//
		return nil, err
	}
	// Seal encrypts and authenticates plaintext, authenticates the additional
	// data and appends the result to dst, returning the updated slice.
	// The nonce must be NonceSize() bytes long
	// and unique for all time, for a given key.
	//
	// To reuse plaintext's storage for the encrypted output,
	// use plaintext[:0] as dst. Otherwise, the remaining
	// capacity of dst must not overlap plaintext.
	//
	ciphertext = gcm.Seal(
		nonce,     // dst []byte,
		nonce,     // nonce []byte,
		plaintext, // plaintext []byte,
		nil,       // additionalData []byte) []byte
	)
	return ciphertext, nil
} //                                                                  encryptAES

// decryptAES decrypts cipherthext using secretKey and returns
// the decrypted plaintext, using AES-256 symmetric cipher.
func decryptAES(ciphertext, secretKey []byte) (plaintext []byte, err error) {
	//
	// NewCipher creates and returns a new cipher.Block.
	// The key argument should be the AES key, either 16, 24,
	// or 32 bytes to select AES-128, AES-192, or AES-256.
	chp, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	// NewGCM returns the given 128-bit, block cipher wrapped
	// in Galois Counter Mode with the standard nonce length.
	//
	// In general, the GHASH operation performed by this implementation
	// of GCM is not constant-time. An exception is when the underlying
	// Block was created by aes.NewCipher on systems with hardware support
	// for AES. See the crypto/aes package documentation for details.
	//
	gcm, err := cipher.NewGCM(chp)
	if err != nil {
		return nil, err
	}
	// NonceSize returns the size of the nonce
	// that must be passed to Seal and Open.
	//
	n := gcm.NonceSize()
	if len(ciphertext) < n {
		return nil, err
	}
	nonce := ciphertext[:n]
	ciphertext = ciphertext[n:]
	//
	// Open decrypts and authenticates ciphertext, authenticates
	// the additional data and, if successful, appends the
	// resulting plaintext to dst, returning the updated slice.
	// The nonce must be NonceSize() bytes long and both it and
	// the additional  data must match the value passed to Seal.
	//
	// To reuse ciphertext's storage for the decrypted output,
	// use ciphertext[:0] as dst. Otherwise, the remaining
	// capacity of dst must not overlap plaintext.
	//
	// Even if the function fails, the contents of dst,
	// up to its capacity, may be overwritten.
	//
	plaintext, err = gcm.Open(
		nil,        // dst []byte
		nonce,      // nonce []byte
		ciphertext, // ciphertext []byte
		nil,        // additionalData []byte
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                  decryptAES

func aesDemo() {
	fmt.Println(div)
	fmt.Println("Running aesDemo")
	var (
		msg    = "The quick brown fox\n"            // 20 bytes
		input  = strings.Repeat(msg, 1024)          // 20K of data
		aesKey = "abcdefghijklmnopqrstuvwxyz789012" // should be 32 bytes
	)
	ciphertext, err := encryptAES([]byte(input), []byte(aesKey))
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	b, err := decryptAES(ciphertext, []byte(aesKey))
	plaintext := string(b)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	n := len(msg) * 10
	if len(plaintext) > n {
		plaintext = plaintext[:n]
	}
	if plaintext == input[:n] {
		fmt.Println("AES encryption and decryption successful")
	}
	fmt.Print("Sample of decrypted plaintext:\n" + plaintext)
} //                                                                     aesDemo

// end
