package passgen

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePassphrases(t *testing.T) {
	type testReqs func(t *testing.T, passphrases []string, err error)

	type testDef struct {
		name      string
		count     uint
		wordCount uint
		separator rune
		casing    PassphraseCasing
		wordList  []string

		requirements testReqs
		setup        func() interface{}
		teardown     func(interface{})
	}

	var tests = []testDef{
		{
			"rational defaults",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.NoError(t, err)
				require.Len(t, passphrases, PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(PassphraseSeparatorDefault))
					require.Len(t, words, PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, WordListDefault, word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"count too small",
			PassphraseCountMin - 1,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"count too large",
			PassphraseCountMax + 1,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word count too small",
			PassphraseCountDefault,
			PassphraseWordCountMin - 1,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word count too large",
			PassphraseCountDefault,
			PassphraseWordCountMax + 1,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"invalid casing",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasing(^uint8(0)),
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"empty word list",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			[]string{},

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word list too small",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			[]string{"x"},

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"random source EOF",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,

			func(t *testing.T, passphrases []string, err error) {
				require.Empty(t, passphrases)
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

				passphrases, err := GeneratePassphrases(
					test.count,
					test.wordCount,
					test.separator,
					test.casing,
					test.wordList,
				)
				test.requirements(t, passphrases, err)

				if test.teardown != nil {
					test.teardown(setupContext)
				}
			},
		)
	}
}
