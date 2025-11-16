package passgen

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// validatePasswordParams validates the count and length parameters.
func validatePasswordParams(count, length uint) error {
	if count < PasswordCountMin || count > PasswordCountMax {
		return fmt.Errorf("count must be at least %d and at most %d", PasswordCountMin, PasswordCountMax)
	}
	if length < PasswordLengthMin || length > PasswordLengthMax {
		return fmt.Errorf("length must be at least %d and at most %d", PasswordLengthMin, PasswordLengthMax)
	}
	return nil
}

// deduplicateAlphabet removes duplicate characters from an alphabet and returns a character set.
func deduplicateAlphabet(alphabet string) ([]rune, error) {
	chars := map[rune]struct{}{}
	for _, char := range alphabet {
		chars[char] = struct{}{}
	}

	var charSet []rune
	for char := range chars {
		charSet = append(charSet, char)
	}

	if len(charSet) < AlphabetLengthMin {
		return nil, fmt.Errorf("alphabet must contain at least %d unique characters", AlphabetLengthMin)
	}

	return charSet, nil
}

// calculateBytesNeeded determines how many random bytes are needed for generation.
func calculateBytesNeeded(itemCount, setSize uint) uint {
	bitsPerItem := uint(math.Ceil(math.Log2(float64(setSize))))
	totalBits := bitsPerItem * itemCount
	bytes := totalBits / 8
	if totalBits%8 > 0 {
		bytes++
	}
	return bytes
}

// sampleFromRuneSet samples items from a rune set using random bytes and bit manipulation.
func sampleFromRuneSet(set []rune, itemCount uint, randomBytes []byte) string {
	bitsPerItem := uint(math.Ceil(math.Log2(float64(len(set)))))
	var b strings.Builder
	var bitIdx, byteIdx uint

	for range itemCount {
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
		b.WriteRune(set[itemIdx])
	}

	return b.String()
}

// GeneratePasswords generates random passwords based on the configuration provided by the user.
func GeneratePasswords(
	count uint, // Number of passwords to generate.
	length uint, // Length of each generated password.
	alphabet string, // Alphabet to pull password characters from.
) (
	passwords []string, // Generated passwords.
	err error, // Possible error encountered during password generation.
) {
	if err := validatePasswordParams(count, length); err != nil {
		return nil, err
	}

	charSet, err := deduplicateAlphabet(alphabet)
	if err != nil {
		return nil, err
	}

	bytesPerPassword := calculateBytesNeeded(length, uint(len(charSet)))

	for range count {
		randomBytes := make([]byte, bytesPerPassword)
		if _, err = io.ReadFull(randSource, randomBytes); err != nil {
			return nil, err
		}

		passwords = append(passwords, sampleFromRuneSet(charSet, length, randomBytes))
	}

	return
}
