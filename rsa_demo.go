// -----------------------------------------------------------------------------
// Go Language Experiments                          go-experiments/[rsa_demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

// This file demonstrates how to encrypt and decrypt short lengths of
// data using RSA (Rivest-Shamir-Adleman) public-key cryptosystem.
//
// Most of the comments are taken from the Go
// language standard library's documentation.

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"strings"
)

const RSA_DEMO_LABEL = "RSA_DEMO_LABEL"

var _ = rsaDemo

// sha256.New returns a new hash.Hash computing the SHA256 checksum.
// The Hash also implements encoding.BinaryMarshaler and
// encoding.BinaryUnmarshaler to marshal and unmarshal
// the internal state of the hash.
var rsaHash = sha256.New()

// encodeAsPEM encodes an rsa.PrivateKey, rsa.PublicKey or a byte
// array as a string in PEM format (PKCS #1, ASN.1 DER form).
// PEM stands for Privacy-Enhanced Mail.
func encodeAsPEM(keyOrMessage interface{}) string {
	switch input := keyOrMessage.(type) {
	case rsa.PrivateKey:
		return encodeAsPEM(&input)
	case *rsa.PrivateKey:
		pemStr := string(pem.EncodeToMemory(&pem.Block{
			Type: "RSA PRIVATE KEY",
			// MarshalPKCS1PrivateKey converts an RSA private key to PKCS #1,
			// ASN.1 DER form. This kind of key is commonly encoded in PEM
			// blocks of type "RSA PRIVATE KEY". For a more flexible key
			// format which is not RSA specific, use MarshalPKCS8PrivateKey.
			Bytes: x509.MarshalPKCS1PrivateKey(
				input, // key *rsa.PrivateKey
			),
		}))
		return pemStr
	case rsa.PublicKey:
		return encodeAsPEM(&input)
	case *rsa.PublicKey:
		pemStr := string(pem.EncodeToMemory(&pem.Block{
			Type: "RSA PUBLIC KEY",
			// MarshalPKCS1PublicKey converts an RSA public key to PKCS #1,
			// ASN.1 DER form. This kind of key is commonly encoded in PEM
			// blocks of type "RSA PUBLIC KEY".
			// See also x509.MarshalPKIXPublicKey(pub)
			Bytes: x509.MarshalPKCS1PublicKey(
				input, // key *rsa.PublicKey
			),
		}))
		return pemStr
	case []byte: // message
		block := &pem.Block{Type: "MESSAGE", Bytes: input}
		// EncodeToMemory returns the PEM encoding of b. If b has
		// invalid headers and cannot be encoded, EncodeToMemory
		// returns nil. If it is important to report details
		// about this error case, use Encode instead.
		pemStr := string(pem.EncodeToMemory(
			block, // b *pem.Block
		))
		return pemStr
	default:
		log.Println("Wrong data type passed to encodeAsPEM")
		return ""
	}
} //                                                                 encodeAsPEM

// rsaCreateKeys generates a new key pair
func rsaCreateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	//
	// rsa.GenerateKey generates an RSA keypair of the
	// given bit size using the random source 'random'
	//
	var privateKey *rsa.PrivateKey
	var err error
	privateKey, err = rsa.GenerateKey(
		rand.Reader, // random io.Reader
		bits,        // bits int
	)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
} //                                                               rsaCreateKeys

// encryptRSA encrypts plaintext with publicKey and returns
// the cihertext, using RSA public-key cryptosystem.
//
func encryptRSA(plaintext []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	// EncryptOAEP encrypts the given message with RSA-OAEP.
	//
	// OAEP is parameterised by a hash function that is used as a random oracle.
	// Encryption and decryption of a given message must use the same
	// hash function and sha256.New() is a reasonable choice.
	//
	//
	// The random parameter is used as a source of entropy to ensure that
	// encrypting the same message twice doesn't result in the same ciphertext.
	//
	// The label parameter may contain arbitrary data that will not be
	// encrypted, but which gives important context to the message.
	// For example, if a given public key is used to decrypt two types of
	// messages then distinct label values could be used to ensure that
	// a ciphertext for one purpose cannot be used for another by an attacker.
	// If not required it can be empty.
	//
	// The message must be no longer than the length of the public
	// modulus minus twice the hash length, minus a further 2.
	var ciphertext []byte
	var err error
	ciphertext, err = rsa.EncryptOAEP(
		rsaHash,                // hash hash.Hash
		rand.Reader,            // random io.Reader
		publicKey,              // pub *rsa.PublicKey
		plaintext,              // msg []byte
		[]byte(RSA_DEMO_LABEL), // label []byte
	)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
} //                                                                  encryptRSA

// decryptRSA decrypts cihertext with privateKey and returns
// the plaintext, using RSA public-key cryptosystem.
//
func decryptRSA(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// DecryptOAEP decrypts ciphertext using RSA-OAEP.
	//
	// OAEP is parameterised by a hash function that is used as a random oracle.
	// Encryption and decryption of a given message must use the
	// same hash function and sha256.New() is a reasonable choice.
	//
	// The random parameter, if not nil, is used to blind the
	// private-key  operation and avoid timing side-channel attacks.
	// Blinding is purely internal to this function - the random
	// data need not match that used when encrypting.
	//
	// The label parameter must match the value given when encrypting.
	// See EncryptOAEP for details.
	var plaintext []byte
	var err error
	plaintext, err = rsa.DecryptOAEP(
		rsaHash,                // hash hash.Hash
		rand.Reader,            // random io.Reader
		privateKey,             // priv *rsa.PrivateKey
		ciphertext,             // ciphertext []byte
		[]byte(RSA_DEMO_LABEL), // label []byte
	)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} //                                                                  decryptRSA

func rsaDemo() {
	fmt.Println(div)
	fmt.Println("Running rsaDemo")
	var (
		plaintext = strings.Repeat("A", 62) // 62 is the greatest message length
		err       error
	)
	// Generate private and public keys.
	// The public key is used to encrypt.
	// The secret private key is used to decrypt.
	privateKey, publicKey, err := rsaCreateKeys(1024)
	if err != nil {
		fmt.Print(err)
	}
	// encrypt the message
	ciphertext, err := encryptRSA([]byte(plaintext), publicKey)
	if err != nil {
		fmt.Print(err)
	}
	// decrypt the message
	decrypted, err := decryptRSA(ciphertext, privateKey)
	if err != nil {
		fmt.Print(err)
	}
	privateKeyPEM := encodeAsPEM(privateKey)
	fmt.Println(privateKeyPEM)
	//
	publicKeyPEM := encodeAsPEM(publicKey)
	fmt.Println(publicKeyPEM)
	//
	ciphertextPEM := encodeAsPEM(ciphertext)
	fmt.Println(ciphertextPEM)
	fmt.Printf("RSA decrypted message to:\n'%s'\n", decrypted)
} //                                                                     rsaDemo

// end
