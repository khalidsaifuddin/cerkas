package helper

import "strings"

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ReplaceUnderscoreWithSpace(str string) string {
	return strings.ReplaceAll(str, "_", " ")
}

func CapitalizeWords(input string) string {
	// Split string menjadi slice berdasarkan spasi
	words := strings.Fields(input)

	for i, word := range words {
		if len(word) > 0 {
			// Ubah huruf pertama menjadi kapital dan gabungkan dengan sisanya
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}

	// Gabungkan kembali slice menjadi string
	return strings.Join(words, " ")
}
