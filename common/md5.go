package common

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5String(text string) string {
	sum := md5.Sum([]byte(text)) // 返回 [16]byte
	return hex.EncodeToString(sum[:])
}
