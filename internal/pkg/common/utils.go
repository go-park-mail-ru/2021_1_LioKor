package common

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	randNumStr := strconv.Itoa(rand.Intn(32000))

	h := sha256.New()
	h.Write([]byte(randNumStr))
	return hex.EncodeToString(h.Sum(nil))
}

