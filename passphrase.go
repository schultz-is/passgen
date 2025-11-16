package passgen

import (
	"fmt"
	"io"
	"math"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PassphraseCasing represents the casing of each word within a passphrase.
type PassphraseCasing uint8

// applyCasing applies the specified casing transformation to a word.
func applyCasing(word string, casing PassphraseCasing) string {
	switch casing {
	case PassphraseCasingLower:
		return strings.ToLower(word)
	case PassphraseCasingUpper:
		return strings.ToUpper(word)
	case PassphraseCasingTitle:
		return cases.Title(language.Und).String(word)
	case PassphraseCasingNone:
		return word
	default:
		return word
	}
}

// validatePassphraseParams validates the count, word count, and casing parameters.
func validatePassphraseParams(count, wordCount uint, casing PassphraseCasing) error {
	if count < PassphraseCountMin || count > PassphraseCountMax {
		return fmt.Errorf("count must be at least %d and at most %d", PassphraseCountMin, PassphraseCountMax)
	}
	if wordCount < PassphraseWordCountMin || wordCount > PassphraseWordCountMax {
		return fmt.Errorf("word count must be at least %d and at most %d", PassphraseWordCountMin, PassphraseWordCountMax)
	}
	switch casing {
	case PassphraseCasingLower, PassphraseCasingUpper, PassphraseCasingTitle, PassphraseCasingNone:
		return nil
	default:
		return fmt.Errorf("invalid word casing")
	}
}

// deduplicateWordList removes duplicate words from a word list and returns a word set.
func deduplicateWordList(wordList []string, casing PassphraseCasing) ([]string, error) {
	words := map[string]struct{}{}
	for _, word := range wordList {
		words[applyCasing(word, casing)] = struct{}{}
	}

	var wordSet []string
	for word := range words {
		wordSet = append(wordSet, word)
	}

	if len(wordSet) < WordListLengthMin {
		return nil, fmt.Errorf("word list must contain at least %d unique words", WordListLengthMin)
	}

	return wordSet, nil
}

// sampleFromStringSet samples items from a string set with separators using random bytes.
func sampleFromStringSet(set []string, itemCount uint, separator rune, randomBytes []byte) string {
	bitsPerItem := uint(math.Ceil(math.Log2(float64(len(set)))))
	var b strings.Builder
	var bitIdx, byteIdx uint

	for i := range itemCount {
		var itemIdx uint
		for range bitsPerItem {
			itemIdx <<= 1
			itemIdx |= uint((randomBytes[byteIdx] & (0x80 >> (bitIdx % 8))) >> (7 - (bitIdx % 8)))
			bitIdx++
			if bitIdx%8 == 0 {
				byteIdx++
			}
		}
		itemIdx %= uint(len(set))
		b.WriteString(set[itemIdx])

		if i < itemCount-1 {
			b.WriteRune(separator)
		}
	}

	return b.String()
}

// GeneratePassphrases generates random passphrases based on the configuration provided by the user.
func GeneratePassphrases(
	count uint, // Number of passphrases to generate.
	wordCount uint, // Length, in words, of each generated passphrase.
	separator rune, // Passphrase word separator.
	casing PassphraseCasing, // Passphrase word casing.
	wordList []string, // List of words to pull passphrase words from.
) (
	passphrases []string, // Generated passphrases.
	err error, // Possible error encountered during passphrase generation.
) {
	if err := validatePassphraseParams(count, wordCount, casing); err != nil {
		return nil, err
	}

	wordSet, err := deduplicateWordList(wordList, casing)
	if err != nil {
		return nil, err
	}

	bytesPerPassphrase := calculateBytesNeeded(wordCount, uint(len(wordSet)))

	for range count {
		randomBytes := make([]byte, bytesPerPassphrase)
		if _, err = io.ReadFull(randSource, randomBytes); err != nil {
			return nil, err
		}

		passphrases = append(passphrases, sampleFromStringSet(wordSet, wordCount, separator, randomBytes))
	}

	return
}
