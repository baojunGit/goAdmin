package biz

import (
	"crypto/md5"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"testing"
)

func TestGetMd5(t *testing.T) {
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, hashedPwd := password.Encode("happy", &options)
	fmt.Println(salt)
	fmt.Println(hashedPwd)
	check := password.Verify("happy", salt, hashedPwd, nil)
	fmt.Println(check)

}
