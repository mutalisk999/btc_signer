package main

import (
	"encoding/hex"
	"fmt"
	"github.com/mutalisk999/bitcoin-lib/src/base58"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	privKeyB58 := "KzutU4gAuMqf9qFayh57Xb6JkCZv6o4jKkZmuEpk5tpVfvvzKqUT"
	privKeyBytes, _ := base58.Decode(privKeyB58)
	privKeyBytes = privKeyBytes[1:33]
	privKeyHex := hex.EncodeToString(privKeyBytes)
	fmt.Println("privKeyHex:", privKeyHex)

	cryptedBytes := AesEncrypt(privKeyHex, SecurityPassStr)
	cryptedHex := hex.EncodeToString(cryptedBytes)
	fmt.Println("cryptedHex:", cryptedHex)
}

func TestAesDecrypt(t *testing.T) {
	cryptedHex := "157dc942fbfd8b23c796d9d6fb0e4337a7ed79c01dae71e249dc27c3ffa74560f8677faeb5ea1259a61fb98fbf231897bbc83e0ea385a1876d531d44dfd1aa832d2fede6a8001a1a2a5abb2e46426445"
	cryptedBytes, _ := hex.DecodeString(cryptedHex)
	privKeyHexStr := string(AesDecrypt(cryptedBytes, []byte(SecurityPassStr)))
	fmt.Println("privKeyHex:", privKeyHexStr)
}
