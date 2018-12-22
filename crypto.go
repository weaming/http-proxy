package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func AESStream(secret string, outWriter io.Writer, inReader io.Reader) {
	stream := getStream(secret)
	reader := &cipher.StreamReader{S: stream, R: inReader}
	io.Copy(outWriter, reader)
}

func getStream(secret string) cipher.Stream {
	key, _ := hex.DecodeString(Md5(secret))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// If the key is unique for each ciphertext, then it's ok to use a zero IV.
	var iv [aes.BlockSize]byte
	return cipher.NewOFB(block, iv[:])
}

func Md5(secret string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(secret)))
}

func AES(origin []byte, secret string) []byte {
	stream := getStream(secret)
	rv := make([]byte, len(origin))
	stream.XORKeyStream(rv, origin)
	return rv
}
