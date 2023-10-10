package utils_test

import (
	"reflect"
	"testing"

	"github.com/yzaimoglu/yzgo/utils"
)

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword("test")
	if err != nil || reflect.ValueOf(hash).IsZero() {
		t.Errorf("HashPassword() failed: %v", err)
		t.FailNow()
	}
	t.Logf("HashPassword() success: %v", hash)
}

func TestGenerateRandomBytes(t *testing.T) {
	randomBytes, err := utils.GenerateRandomBytes(32)
	if err != nil || reflect.ValueOf(randomBytes).IsZero() || len(randomBytes) != 32 {
		t.Errorf("GenerateRandomBytes() failed: %v", err)
		t.FailNow()
	}
	t.Logf("GenerateRandomBytes() success: %v", randomBytes)
}

func TestHashSHA512(t *testing.T) {
	hash := utils.HashSHA512("test")
	if hash != "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff" {
		t.Errorf("HashSHA512() failed: %v", hash)
		t.FailNow()
	}
	t.Logf("HashSHA512(\"test\") pass: %v", hash)
	hash = utils.HashSHA512("hello")
	if hash != "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043" {
		t.Errorf("HashSHA512() failed: %v", hash)
		t.FailNow()
	}
	t.Logf("HashSHA512(\"hello\") pass: %v", hash)
	t.Logf("HashSHA512() success")
}

func TestHashAndCheckPassword(t *testing.T) {
	encodedHash, err := utils.HashPassword("test")
	if err != nil || reflect.ValueOf(encodedHash).IsZero() {
		t.Errorf("HashPassword() failed: %v", err)
		t.FailNow()
	}
	t.Logf("HashPassword() pass: %v", encodedHash)
	match, err := utils.CheckPassword("test", encodedHash)
	if err != nil || !match {
		t.Errorf("CheckPassword() failed: %v", err)
		t.FailNow()
	}
	t.Logf("CheckPassword() right password pass: %v", match)

	match, err = utils.CheckPassword("test1", encodedHash)
	if err != nil || match {
		t.Errorf("CheckPassword() failed: %v", err)
		t.FailNow()
	}
	t.Logf("CheckPassword() wrong password pass: %v", match)
	t.Logf("HashAndCheckPassword() success")
}

func TestHashMD5(t *testing.T) {
	hash := utils.HashMD5("test")
	if hash != "098f6bcd4621d373cade4e832627b4f6" {
		t.Errorf("HashMD5() failed: %v", hash)
		t.FailNow()
	}
	t.Logf("HashMD5(\"test\") pass: %v", hash)
	hash = utils.HashMD5("hello")
	if hash != "5d41402abc4b2a76b9719d911017c592" {
		t.Errorf("HashMD5() failed: %v", hash)
		t.FailNow()
	}
	t.Logf("HashMD5(\"hello\") pass: %v", hash)
	t.Logf("HashMD5() success")
}

func TestRandomAlphanumeric(t *testing.T) {
	randomAlphanumeric, err := utils.RandomAlphanumeric(32)
	if err != nil || reflect.ValueOf(randomAlphanumeric).IsZero() || len(randomAlphanumeric) != 32 {
		t.Errorf("RandomAlphanumeric() failed: %v", randomAlphanumeric)
		t.FailNow()
	}
	t.Logf("RandomAlphanumeric() success: %v", randomAlphanumeric)
}
