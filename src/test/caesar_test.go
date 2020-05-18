package test

import (
	caesar "caesar/cipher"
	"testing"
)

func TestCipher(t *testing.T) {
	cipherTest := "zy esp tyepcype, yzmzoj vyzhd jzf lcp l ozr. apepc deptypc"
	plainTest := "on the internet, nobody knows you are a dog. peter steiner"
	rotate := 11
	cipher := caesar.Encode(plainTest, rotate)
	plain := caesar.Decode(cipherTest, rotate)
	if cipher != cipherTest {
		t.Errorf("Cipher was incorrect, got: %s, want: %s.", cipher, cipherTest)
	}
	if plain != plainTest {
		t.Errorf("Plain was incorrect, got: %s, want: %s.", plain, plainTest)
	}
}
