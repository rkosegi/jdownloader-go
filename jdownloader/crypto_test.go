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
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	assert.Equal(t, [32]byte{
		0xee, 0x06, 0xc8, 0x31, 0x92, 0xf3, 0xc0, 0x33, 0x27, 0xfd, 0x93, 0x0, 0x10,
		0xe8, 0x8f, 0x21, 0x3c, 0x21, 0x10, 0x51, 0x44, 0xe8, 0x6a, 0x64, 0x15, 0xfc,
		0x3c, 0xff, 0x2d, 0xf0, 0x1e, 0x57,
	}, createSecret("test@acme.tld", "123456", "something"))
}

func TestEncryptDecryptLong(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	key := createSecret("test@acme.tld", "123456", "something")
	size := 500
	data := make([]byte, size)
	_, err := r.Read(data)
	assert.Nil(t, err)
	encrypted, _ := encrypt(data, key)
	// t.Log(base64.StdEncoding.EncodeToString(encrypted))
	decrypted, _ := decrypt(encrypted, key)
	assert.Equal(t, decrypted, data)
}

func TestEncryptDecryptShort(t *testing.T) {
	key := createSecret("test@acme.tld", "123456", "something")
	ciphertext, _ := encrypt([]byte("test"), key)
	plaintext, _ := decrypt(ciphertext, key)
	assert.Equal(t, []byte("test"), plaintext)
}
