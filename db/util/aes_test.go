package util_test

import (
	"testing"

	"github.com/yaitoo/sparrow/db/util"
)

func TestAesSecretTooShort(t *testing.T) {
	secret := "fz4r2g2deb"
	plainText := "hello"
	_, err := util.AesEncrypt(plainText, secret)
	if err != util.ErrSecretTooShort {
		t.Fatal(err)
	}

	_, err = util.AesDecrypt("abasda", secret)
	if err != util.ErrSecretTooShort {
		t.Fatal(err)
	}
}

func TestAESEncryptAndDecrypt(t *testing.T) {
	secret := "fz4r2g2debh9gofc"
	plainText := "hello"
	cipherText, err := util.AesEncrypt(plainText, secret)
	if err != nil {
		t.Fatal(err)
	}

	decryptResult, err := util.AesDecrypt(cipherText, secret)
	if err != nil {
		t.Fatal(err)
	}
	if decryptResult != plainText {
		t.Fatal("")
	}
}
