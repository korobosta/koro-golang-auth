package korologin

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func currentTimeStamp() string {
	now := time.Now()
	unix := now.Unix()
	timestamp := strconv.FormatInt(unix, 10)

	return timestamp
}

func rolesContains(s []string, e string) int {
	i := 0
	for _, a := range s {
		if a == e {
			return i
		}
		i++
	}
	return -1
}

func EncNoEncrypt(password string) string {

	return password
}

func EncMD5(password string) string {

	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
