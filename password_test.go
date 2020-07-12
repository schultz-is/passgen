package passgen

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type generatePasswordRequirements func(t *testing.T, pws []string, err error)

type generatePasswordsTest struct {
	name         string
	count        uint
	length       uint
	alphabet     string
	requirements generatePasswordRequirements
	setup        func() interface{}
	teardown     func(interface{})
}

var tests []generatePasswordsTest = []generatePasswordsTest{
	{
		"rational defaults",
		PasswordCountDefault,
		PasswordLengthDefault,
		AlphabetDefault,
		func(t *testing.T, pws []string, err error) {
			require.NoError(t, err)
			require.Len(t, pws, PasswordCountDefault)
			for _, pw := range pws {
				require.Len(t, pw, PasswordLengthDefault)
				for _, c := range pw {
					require.Contains(t, AlphabetDefault, string(c))
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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
		func(t *testing.T, pws []string, err error) {
			require.Empty(t, pws)
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

func TestGeneratePasswords(t *testing.T) {
	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				var setupContext interface{}
				if test.setup != nil {
					setupContext = test.setup()
				}

				pws, err := GeneratePasswords(
					test.count,
					test.length,
					test.alphabet,
				)
				test.requirements(t, pws, err)

				if test.teardown != nil {
					test.teardown(setupContext)
				}
			},
		)
	}
}
