package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"crypto/rand"
	"io"
	"encoding/base64"
	"strings"
 	"encoding/hex"
)


// AES encryption using GCM mode (widely adopted because of its efficiency and performance)
// The key argument should be the AES key,
// either 16, 24, or 32 bytes corresponding to the AES-128, AES-192 or AES-256 algorithms, respectively
func encrypt(plaintext []byte, key []byte) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.URLEncoding.EncodeToString(ciphertext),nil
}

func decrypt(ciphertext []byte, key []byte) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	decryptedBytes,err:= gcm.Open(nil, nonce, ciphertext, nil)
	return string(decryptedBytes),err
}

// format encrypted text , keyId is first 6 chars of hashed(sha256) signing key
// {{key-id}}${{algo}}${{encrypted-value}}
// example: e67gef$aes256$adsad321424324fdsdfs3Rddi90oP34xV
func Encrypt(text,key,keyId string) (string,error){
	// base 64 decoding of key
	KeyByte, err := hex.DecodeString(key)
	cipherText,err:=encrypt([]byte(text),KeyByte)
	if err != nil{
		return "",err
	}
	return keyId + "$" + "aes256" + "$" + cipherText,nil
}


func Decrypt(FormattedCipherText,key,keyId string) (string,error){
	formatEncryption := keyId + "$" + "aes256" + "$"
	if strings.Index(FormattedCipherText, formatEncryption) == -1{
		return "",errors.New("Cipher text is not well formatted")
	}
	//keep cipher text only
 	cipherText:=strings.Replace(FormattedCipherText,formatEncryption,"",-1)
	// base 64 decoding of key and text
	KeyByte, err := hex.DecodeString(key)
	cipherTextByte, err := base64.URLEncoding.DecodeString(cipherText)
	text,err:=decrypt(cipherTextByte,KeyByte)
	if err != nil{
		return "",err
	}
	return text,nil
}

func IsTextEncrypted(FormattedCipherText,key,keyId string) (bool,error){
	formatEncryption := keyId + "$" + "aes256" + "$"
	if strings.Index(FormattedCipherText, formatEncryption) == -1{
		return false,errors.New("Cipher text is not well formatted")
	}
	//keep cipher text only
	cipherText:=strings.Replace(FormattedCipherText,formatEncryption,"",-1)
	// base 64 decoding of key and text
	KeyByte, err := base64.URLEncoding.DecodeString(key)
	cipherTextByte, err := base64.URLEncoding.DecodeString(cipherText)
	_,err=decrypt(cipherTextByte,KeyByte)
	if err != nil{
		return false,err
	}
	return true,nil
}