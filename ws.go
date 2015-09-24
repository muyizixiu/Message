package M

import (
	"errors"
)

func decode(b []byte) ([]byte, error) {
	length := len(b)
	if length < 6 || length > 1023 {
		return nil, nil
	}
	var start, end int
	if b[1] > 253 {
		start = 8
		end = int(b[2])*256 + int(b[3]) + start
	} else {
		start = 6
		end = int(b[1]) - 122
	}
	mask_key := b[start-4 : start]
	body := b[start:end]
	var re []byte
	for i, v := range body {
		re = append(re, (mask_key[i%4] ^ v))
	}
	return re, nil
}
func encode(body []byte) ([]byte, error) {
	re := []byte{129}
	length := len(body)
	if length > 1023 {
		return nil, errors.New("out of 1024,can not send data!")
	}
	if length > 125 {
		re = append(re, uint8(126))
		re = append(re, uint8(length>>8), uint8(length%256))
	} else {
		re = append(re, uint8(length))
	}
	re = append(re, body...)
	return re, nil
}
