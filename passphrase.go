package passgen

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// PassphraseCasing represents the casing of each word within a passphrase.
type PassphraseCasing uint8

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
	// Validate the supplied count parameter.
	if count < PassphraseCountMin || count > PassphraseCountMax {
		return nil, fmt.Errorf("count must be at least %d and at most %d", PassphraseCountMin, PassphraseCountMax)
	}

	// Validate the supplied word count parameter.
	if wordCount < PassphraseWordCountMin || wordCount > PassphraseWordCountMax {
		return nil, fmt.Errorf("word count must be at least %d and at most %d", PassphraseWordCountMin, PassphraseWordCountMax)
	}

	// Validate the supplied casing parameter.
	switch casing {
	case PassphraseCasingLower, PassphraseCasingUpper, PassphraseCasingTitle, PassphraseCasingNone:
		break
	default:
		return nil, fmt.Errorf("invalid word casing")
	}

	// Deduplicate the provided word list.
	words := map[string]struct{}{}
	for _, word := range wordList {
		switch casing {
		case PassphraseCasingLower:
			words[strings.ToLower(word)] = struct{}{}
		case PassphraseCasingUpper:
			words[strings.ToUpper(word)] = struct{}{}
		case PassphraseCasingTitle:
			words[strings.Title(word)] = struct{}{}
		case PassphraseCasingNone:
			words[word] = struct{}{}
		}
	}
	var wordSet []string
	for word := range words {
		wordSet = append(wordSet, word)
	}

	// Validate the provided word list.
	if len(wordSet) < WordListLengthMin {
		return nil, fmt.Errorf("word list must contain at least %d unique words", WordListLengthMin)
	}

	// Determine how many bytes are needed to represent a passphrase of the specified word count in
	// the provided word list.
	bitsPerWord := uint(math.Ceil(math.Log2(float64(len(wordSet)))))
	bitsPerPassphrase := bitsPerWord * wordCount
	bytesPerPassphrase := bitsPerPassphrase / 8
	if bitsPerPassphrase%8 > 0 {
		bytesPerPassphrase++
	}

	var (
		i                uint            // Passphrase counter.
		b                strings.Builder // String builder for efficiently constructing passphrases.
		passphraseBuffer []byte          // Byte buffer for random data used as a passphrase source.
	)

	for i = 0; i < count; i++ {
		// Read enough random data to sufficiently produce a passphrase.
		passphraseBuffer = make([]byte, bytesPerPassphrase)
		_, err = io.ReadFull(randSource, passphraseBuffer)
		if err != nil {
			return nil, err
		}

		var (
			j          uint // Passphrase word counter.
			wordIdx    uint // Word index within the provided word list.
			bitIdx     uint // Source buffer bit counter.
			byteIdx    uint // Source buffer byte counter.
			wordBitIdx uint // Passphrase word bit counter.
		)

		for j = 0; j < wordCount; j++ {
			for wordBitIdx = 0; wordBitIdx < bitsPerWord; wordBitIdx++ {
				// Left shift the word index to read the next bit.
				wordIdx <<= 1

				// Set the bit from the passphrase source in the word index.
				wordIdx |= uint((passphraseBuffer[byteIdx] & (0x80 >> (bitIdx % 8))) >> (7 - (bitIdx % 8)))

				// Increment the bit counter and, if necessary, the byte counter.
				bitIdx++
				if bitIdx%8 == 0 {
					byteIdx++
				}
			}

			// Ensure the word index is within the bounds of the word set.
			wordIdx = wordIdx % uint(len(wordSet))

			// Retrieve the word from the word set and write it to the passphrase.
			b.WriteString(wordSet[wordIdx])

			// Write the provided separator if this is not the final word in the passphrase.
			if j < wordCount-1 {
				b.WriteRune(separator)
			}

			// Reset the word index value.
			wordIdx = 0
		}

		// Append the passphrase to the return list.
		passphrases = append(passphrases, b.String())

		// Reset the string builder for the next passphrase.
		b.Reset()
	}

	return
}
