package passgen

import "crypto/rand"

const (
	PasswordCountMin     = 1    // Fewest allowed passwords to generate.
	PasswordCountMax     = 1024 // Most allowed passwords to generate.
	PasswordCountDefault = 1    // Default number of passwords to generate.

	PasswordLengthMin     = 5    // Shortest allowed password to generate.
	PasswordLengthMax     = 1024 // Longest allowed password to generate.
	PasswordLengthDefault = 16   // Default length of password to generate.

	AlphabetLengthMin        = 2                                               // Smallest allowed alphabet.
	AlphabetLower            = "abcdefghijkmnopqrstuvwxyz"                     // Lowercase English letters, ambiguous characters removed.
	AlphabetLowerAmbiguous   = "l" + AlphabetLower                             // Lowercase English letters.
	AlphabetUpper            = "ABCDEFGHJKLMNPQRSTUVWXYZ"                      // Uppercase English letters, ambiguous characters removed.
	AlphabetUpperAmbiguous   = "IO" + AlphabetUpper                            // Uppercase English letters.
	AlphabetNumeric          = "23456789"                                      // Arabic numerals, ambiguous characters removed.
	AlphabetNumericAmbiguous = "01" + AlphabetNumeric                          // Arabic numerals.
	AlphabetSpecial          = "!@#$%^&*_-+="                                  // Selection of special characters.
	AlphabetDefault          = AlphabetLower + AlphabetUpper + AlphabetNumeric // Alphanumeric English characters, ambiguous characters removed.
)

var (
	// By default, the generators will use the random source provided by crypto/rand. This
	// package-level variable is only included to aid test coverage.
	randSource = rand.Reader
)
