package utils

import "unicode"

// MixCase чередует нижний и верхний регистр букв сохраняя остальные символы
func MixCase(text string) string {
	result := []rune(text)
	upper := false
	for i, letter := range result {
		if !unicode.IsLetter(letter) {
			continue
		}
		if upper {
			result[i] = unicode.ToUpper(letter)
		} else {
			result[i] = unicode.ToLower(letter)
		}
		upper = !upper
	}
	return string(result)
}
