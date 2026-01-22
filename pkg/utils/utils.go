// Copyright (C) 2021-2023 YumeMichi
//
// SPDX-License-Identifier: Apache-2.0
package utils

import (
	"encoding/hex"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	rwMutex sync.RWMutex
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadAllText(path string) string {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func WriteAllText(path, text string) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	_ = os.WriteFile(path, []byte(text), 0644)
}

func XOR(s1, s2 []byte) []byte {
	n := min(len(s1), len(s2))

	res := make([]byte, n)
	for i, b := range s1[:n] {
		res[i] = b ^ s2[i]
	}
	return res
}

func Sub16(str []byte) []byte {
	return str[16:]
}

func RandomStr(len int) string {
	rand.Seed(time.Now().UnixNano())
	mRand := make([]byte, len)
	rand.Read(mRand)
	mRandStr := hex.EncodeToString(mRand)[0:len]

	return mRandStr
}
