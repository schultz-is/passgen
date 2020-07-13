package passgen

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// GeneratePasswords generates random passwords based on the configuration provided by the user.
func GeneratePasswords(
	count uint, // Number of passwords to generate.
	length uint, // Length of each generated password.
	alphabet string, // Alphabet to pull password characters from.
) (
	passwords []string, // Generated passwords.
	err error, // Possible error encountered during password generation.
) {
	// Validate the supplied count parameter.
	if count < PasswordCountMin || count > PasswordCountMax {
		return nil, fmt.Errorf("count must be at least %d and at most %d", PasswordCountMin, PasswordCountMax)
	}

	// Validate the supplied length parameter.
	if length < PasswordLengthMin || length > PasswordLengthMax {
		return nil, fmt.Errorf("length must be at least %d and at most %d", PasswordLengthMin, PasswordLengthMax)
	}

	// Deduplicate the provided alphabet.
	chars := map[rune]struct{}{}
	for _, char := range alphabet {
		chars[char] = struct{}{}
	}
	var charSet []rune
	for char := range chars {
		charSet = append(charSet, char)
	}

	// Validate the provided alphabet.
	if len(charSet) < AlphabetLengthMin {
		return nil, fmt.Errorf("alphabet must contain at least %d unique characters", AlphabetLengthMin)
	}

	// Determine how many bytes are needed to represent a password of the specified length in the
	// provided alphabet.
	bitsPerChar := uint(math.Ceil(math.Log2(float64(len(charSet)))))
	bitsPerPassword := bitsPerChar * length
	bytesPerPassword := bitsPerPassword / 8
	if bitsPerPassword%8 > 0 {
		bytesPerPassword++
	}

	var (
		i              uint            // Password counter.
		b              strings.Builder // String builder for efficiently constructing passwords.
		passwordBuffer []byte          // Byte buffer for random data used as password source.
	)

	for i = 0; i < count; i++ {
		// Read enough random data to sufficiently produce a password.
		passwordBuffer = make([]byte, bytesPerPassword)
		_, err = io.ReadFull(randSource, passwordBuffer)
		if err != nil {
			return nil, err
		}

		var (
			j          uint // Password character counter.
			charIdx    uint // Character index within the provided alphabet.
			bitIdx     uint // Source buffer bit counter.
			byteIdx    uint // Source buffer byte counter.
			charBitIdx uint // Password character bit counter.
		)

		for j = 0; j < length; j++ {
			for charBitIdx = 0; charBitIdx < bitsPerChar; charBitIdx++ {
				// Left shift the character index to read the next bit.
				charIdx <<= 1

				// Set the bit from the password source in the character index.
				charIdx |= uint((passwordBuffer[byteIdx] & (0x80 >> (bitIdx % 8))) >> (7 - (bitIdx % 8)))

				// Increment the bit counter and, if necessary, the byte counter.
				bitIdx++
				if bitIdx%8 == 0 {
					byteIdx++
				}
			}

			// Ensure the character index is within the bounds of the alphabet.
			charIdx = charIdx % uint(len(charSet))

			// Retrieve the character from the alphabet and write it to the password.
			b.WriteRune(charSet[charIdx])

			// Reset character index value.
			charIdx = 0
		}

		// Append the password to the return list.
		passwords = append(passwords, b.String())

		// Reset the string builder for the next password.
		b.Reset()
	}

	return
}
