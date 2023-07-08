/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
)

const EncodeKeyEnvStr = "ENCRYPTION_SECRET"

// TODO: maybe move encryption/decryption into helper?
// AES + Base64 encryption using ENCRYPTION_SECRET in .env as key
func Encrypt(encryptionSecret, plainText string) (string, errors.Error) {
	// add suffix to the data part
	inputBytes := append([]byte(plainText), 123, 110, 100, 100, 116, 102, 125)
	// perform encryption
	output, err := AesEncrypt(inputBytes, []byte(encryptionSecret))
	if err != nil {
		return plainText, err
	}
	// Return the result after Base64 processing
	return base64.StdEncoding.EncodeToString(output), nil
}

// Base64 + AES decryption using ENCRYPTION_SECRET in .env as key
func Decrypt(encryptionSecret, encryptedText string) (string, errors.Error) {
	// when encryption key is not set
	if encryptionSecret == "" {
		// return error message
		return encryptedText, errors.Default.New("encryptionSecret is required")
	}

	// Decode Base64
	decodingFromBase64, err1 := base64.StdEncoding.DecodeString(encryptedText)
	if err1 != nil {
		return encryptedText, errors.Convert(err1)
	}
	// perform AES decryption
	output, err2 := AesDecrypt(decodingFromBase64, []byte(encryptionSecret))
	if err2 != nil {
		return encryptedText, err2
	}

	// Verify and remove suffix
	oSize := len(output)
	if oSize >= 7 {
		check := output[oSize-7 : oSize]
		backEnd := []byte{123, 110, 100, 100, 116, 102, 125}
		if string(check) == string(backEnd) {
			output = output[0 : oSize-7]
			// return result
			return string(output), nil
		}
	}
	return "", errors.Default.New("invalid encryptionSecret")
}

// PKCS7Padding PKCS7 padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7 unPadding
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unpadding := int(origData[length-1])
	if unpadding >= length {
		return nil
	}
	return origData[:(length - unpadding)]
}

// AesEncrypt AES encryption, CBC
func AesEncrypt(origData, key []byte) ([]byte, errors.Error) {
	// data alignment fill and encryption
	sha256Key := sha256.Sum256(key)
	key = sha256Key[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// data alignment fill and encryption
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt AES decryption
func AesDecrypt(crypted, key []byte) ([]byte, errors.Error) {
	// Uniformly use sha256 to process as 32-bit Byte (256-bit bit)
	sha256Key := sha256.Sum256(key)
	key = sha256Key[:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// Get the block size and check whether the ciphertext length is legal
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, errors.Default.New(fmt.Sprintf("The length of the data to be decrypted is [%d], so cannot match the required block size [%d]", len(crypted), blockSize))
	}

	// Decrypt and unalign data
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// RandomEncryptionSecret will return a random string of length 128
func RandomEncryptionSecret() (string, errors.Error) {
	return utils.RandLetterBytes(128)
}
