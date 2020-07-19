package passgen

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestGeneratePasswords(t *testing.T) {
	type testReqs func(t *testing.T, passwords []string, err error)

	type testDef struct {
		name     string
		count    uint
		length   uint
		alphabet string

		requirements testReqs
		setup        func() interface{}
		teardown     func(interface{})
	}

	var tests = []testDef{
		{
			"rational defaults",
			PasswordCountDefault,
			PasswordLengthDefault,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.NoError(t, err)
				require.Len(t, passwords, PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(t, AlphabetDefault, string(char))
					}
				}
			},

			nil,
			nil,
		},
		{
			"non-default parameters",
			PasswordCountDefault + 1,
			PasswordLengthDefault + 1,
			AlphabetDefault + "-",

			func(t *testing.T, passwords []string, err error) {
				require.NoError(t, err)
				require.Len(t, passwords, PasswordCountDefault+1)
				for _, password := range passwords {
					require.Equal(t, PasswordLengthDefault+1, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(t, AlphabetDefault+"-", string(char))
					}
				}
			},

			nil,
			nil,
		},
		{
			"multibyte alphabet characters",
			PasswordCountDefault,
			PasswordLengthDefault,
			"🐶🐱",

			func(t *testing.T, passwords []string, err error) {
				require.NoError(t, err)
				require.Len(t, passwords, PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(t, "🐶🐱", string(char))
					}
				}
			},

			nil,
			nil,
		},
		{
			"count too small",
			PasswordCountMin - 1,
			PasswordLengthDefault,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"count too large",
			PasswordCountMax + 1,
			PasswordLengthDefault,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"length too small",
			PasswordCountDefault,
			PasswordLengthMin - 1,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"length too large",
			PasswordCountDefault,
			PasswordLengthMax + 1,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"empty alphabet",
			PasswordCountDefault,
			PasswordLengthDefault,
			"",

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"alphabet too small",
			PasswordCountDefault,
			PasswordLengthDefault,
			"x",

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"random source EOF",
			PasswordCountDefault,
			PasswordLengthDefault,
			AlphabetDefault,

			func(t *testing.T, passwords []string, err error) {
				require.Empty(t, passwords)
				require.Error(t, err)
			},

			func() interface{} {
				originalRandSource := randSource
				randSource = new(bytes.Reader)
				return originalRandSource
			},
			func(setupContext interface{}) {
				randSource = setupContext.(io.Reader)
			},
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var setupContext interface{}
				if test.setup != nil {
					setupContext = test.setup()
				}

				passwords, err := GeneratePasswords(
					test.count,
					test.length,
					test.alphabet,
				)
				test.requirements(t, passwords, err)

				if test.teardown != nil {
					test.teardown(setupContext)
				}
			},
		)
	}
}

func BenchmarkGeneratePasswords(b *testing.B) {
	type benchmarkDef struct {
		name     string
		count    uint
		length   uint
		alphabet string
	}

	// Generate CJK unified ideographs Unicode block for large alphabet benchmarks.
	var (
		sb strings.Builder
		i  rune
	)
	for i = '\u4e00'; i < '\u9fff'; i++ {
		sb.WriteRune(i)
	}
	cjkUnifiedAlphabet := sb.String()

	var benchmarks = []benchmarkDef{
		// Preset value benchmarks.
		{
			"defaults",
			PasswordCountDefault,
			PasswordLengthDefault,
			AlphabetDefault,
		},
		{
			"minimums",
			PasswordCountMin,
			PasswordLengthMin,
			"ab",
		},
		{
			"maximums",
			PasswordCountMax,
			PasswordLengthMax,
			cjkUnifiedAlphabet,
		},

		// Increasing count benchmarks.
		{
			"25_16_default",
			25,
			16,
			AlphabetDefault,
		},
		{
			"50_16_default",
			50,
			16,
			AlphabetDefault,
		},
		{
			"100_16_default",
			100,
			16,
			AlphabetDefault,
		},
		{
			"250_16_default",
			250,
			16,
			AlphabetDefault,
		},
		{
			"500_16_default",
			500,
			16,
			AlphabetDefault,
		},

		// Increasing length benchmarks.
		{
			"16_25_default",
			16,
			25,
			AlphabetDefault,
		},
		{
			"16_50_default",
			16,
			50,
			AlphabetDefault,
		},
		{
			"16_100_default",
			16,
			100,
			AlphabetDefault,
		},
		{
			"16_250_default",
			16,
			250,
			AlphabetDefault,
		},
		{
			"16_500_default",
			16,
			500,
			AlphabetDefault,
		},

		// Increasing alphabet size benchmarks.
		{
			"16_16_26",
			16,
			16,
			AlphabetLowerAmbiguous,
		},
		{
			"16_16_52",
			16,
			16,
			AlphabetLowerAmbiguous + AlphabetUpperAmbiguous,
		},
		{
			"16_16_62",
			16,
			16,
			AlphabetLowerAmbiguous + AlphabetUpperAmbiguous + AlphabetNumericAmbiguous,
		},
		{
			"16_16_74",
			16,
			16,
			AlphabetLowerAmbiguous + AlphabetUpperAmbiguous + AlphabetNumericAmbiguous + AlphabetSpecial,
		},
		{
			"16_16_20991",
			16,
			16,
			cjkUnifiedAlphabet,
		},

		// Special cases.
		{
			"once-duplicated alphabet",
			16,
			16,
			AlphabetDefault + AlphabetDefault,
		},
		{
			"twice-duplicated alphabet",
			16,
			16,
			AlphabetDefault + AlphabetDefault + AlphabetDefault,
		},
		{
			"thrice-duplicated alphabet",
			16,
			16,
			AlphabetDefault + AlphabetDefault + AlphabetDefault + AlphabetDefault,
		},
	}

	for _, benchmark := range benchmarks {
		b.Run(
			benchmark.name,
			func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					_, _ = GeneratePasswords(
						benchmark.count,
						benchmark.length,
						benchmark.alphabet,
					)
				}
			},
		)
	}
}
