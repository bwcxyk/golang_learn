package aes

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAesECB(t *testing.T) {
	str := "{\"hello\", \"world\"}"
	key := "12345678"
	en, err := JavaEncryptECB(str, key)
	if err != nil {
		t.Log(err)
	}
	t.Log(en)
	de, err := JavaDecryptECB(en, key)
	if err != nil {
		t.Log(err)
	}
	t.Log(de)
	assert.Equal(t, de, str, "ecb fail.")
}

func TestAesCBC(t *testing.T) {
	str := "{\"hello\", \"world\"}"
	key := "12345678"
	en, err := JavaEncryptCBC(str, key)
	if err != nil {
		t.Log(err)
	}
	t.Log(en)
	de, err := JavaDecryptCBC(en, key)
	if err != nil {
		t.Log(err)
	}
	t.Log(de)
	assert.Equal(t, de, str, "cbc fail.")
}
