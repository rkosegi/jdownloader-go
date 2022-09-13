/*
Copyright 2022 Richard Kosegi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jdownloader

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func decode(body []byte, key [32]byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, errors.Wrap(err, "can't decode base64 string")
	}
	return decrypt(decoded, key)
}

func decrypt(ciphertext []byte, key [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[16:])
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(block, key[:16])
	plaintext := make([]byte, len(ciphertext))
	decrypter.CryptBlocks(plaintext, ciphertext)
	plaintext = trimPKCS5(plaintext)
	return plaintext, nil
}

//encrypt Encrypts plaintext using AES-128 with provided key
//first half of key is IV and second half is actual key
func encrypt(plaintext []byte, key [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[16:])
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:16])
	paddedPlainText := padPKCS7(plaintext, encrypter.BlockSize())
	ciphertext := make([]byte, len(paddedPlainText))
	encrypter.CryptBlocks(ciphertext, paddedPlainText)
	return ciphertext, nil
}

//sign Signs uri with key using Sha256HMAC
func sign(uri string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(uri))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createSecret(email string, password string, domain string) [32]byte {
	return sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", strings.ToLower(email), password, domain)))
}

func hashConfigKeys(email string, password string) string {
	h := sha256.New()
	h.Write([]byte(email))
	h.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func updateToken(newToken []byte, existing [32]byte) [32]byte {
	var buff bytes.Buffer
	buff.Write(existing[:])
	buff.Write(newToken)
	return sha256.Sum256(buff.Bytes())
}

func padPKCS7(topad []byte, blockSize int) []byte {
	padding := blockSize - len(topad)%blockSize
	padded := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(topad, padded...)
}

func trimPKCS5(text []byte) []byte {
	padding := text[len(text)-1]
	return text[:len(text)-int(padding)]
}
