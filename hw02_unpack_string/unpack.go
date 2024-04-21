package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func modifyString(sourceSlice []rune, index, multiplier int) []rune {
	var multiItemSlice []rune
	for item := multiplier; item > 0; item-- {
		multiItemSlice = append(multiItemSlice, sourceSlice[index])
	}
	multiItemSlice = append(multiItemSlice, sourceSlice[index+2:]...)
	resultSlice := sourceSlice[:index]
	resultSlice = append(resultSlice, multiItemSlice...)

	return resultSlice
}

func isNumber(symbolByte rune) bool {
	if int(symbolByte) >= 48 && int(symbolByte) <= 57 {
		return true
	}

	return false
}

func Unpack(str string) (string, error) {
	strRunes := []rune(str)
	for i := 0; i < len(strRunes); i++ {
		if int(strRunes[i]) == 92 {
			if i < len(strRunes)-2 && (isNumber(strRunes[i+1]) || int(strRunes[i+1]) == 92) {
				strRunes = append(strRunes[:i], strRunes[i+1:]...)
				continue
			} else if i == len(strRunes)-2 {
				strRunes = append(strRunes[:i], strRunes[i+1:]...)
				break
			}

			return "", ErrInvalidString
		}
		if isNumber(strRunes[i]) {
			if i == 0 || (i < len(strRunes)-1 && isNumber(strRunes[i+1])) {
				return "", ErrInvalidString
			}
			multiper := int(strRunes[i]) - 48
			strRunes = modifyString(strRunes, i-1, multiper)
			i = i + multiper - 2
			continue
		}

		continue
	}
	return string(strRunes), nil
}
