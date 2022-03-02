package blockchain

import "strconv"

func IntToHex(in int64) []byte {
	return []byte(strconv.FormatInt(in, 16))
}
