package passgen

import (
	"bytes"
	"io"
	"testing"

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
					require.Len(t, password, PasswordLengthDefault)
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
					require.Len(t, password, PasswordLengthDefault+1)
					for _, char := range password {
						require.Contains(t, AlphabetDefault+"-", string(char))
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
