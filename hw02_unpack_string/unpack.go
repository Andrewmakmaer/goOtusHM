package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func modifyString(sourceSlice []byte, index, multiplier int) []byte {
	var multiItemSlice []byte
	for item := multiplier; item > 0; item-- {
		multiItemSlice = append(multiItemSlice, sourceSlice[index])
	}
	multiItemSlice = append(multiItemSlice, sourceSlice[index+2:]...)
	resultSlice := sourceSlice[:index]
	resultSlice = append(resultSlice, multiItemSlice...)

	return resultSlice
}

func isNumber(symbolByte byte) bool {
	if int(symbolByte) >= 48 && int(symbolByte) <= 57 {
		return true
	}

	return false
}

func Unpack(str string) (string, error) {
	strBytes := []byte(str)
	for i := 0; i < len(strBytes); i++ {
		if int(strBytes[i]) == 92 {
			if i < len(strBytes)-2 && (isNumber(strBytes[i+1]) || int(strBytes[i+1]) == 92) {
				strBytes = append(strBytes[:i], strBytes[i+1:]...)
				continue
			} else if i == len(strBytes)-2 {
				strBytes = append(strBytes[:i], strBytes[i+1:]...)
				break
			}

			return "", ErrInvalidString
		}
		if isNumber(strBytes[i]) {
			if i == 0 || (i < len(strBytes)-1 && isNumber(strBytes[i+1])) {
				return "", ErrInvalidString
			}
			multiper := int(strBytes[i]) - 48
			strBytes = modifyString(strBytes, i-1, multiper)
			i = i + multiper - 2
			continue
		}

		continue
	}
	return string(strBytes), nil
}
