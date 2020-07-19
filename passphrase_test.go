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

	var alternateWordList = []string{
		"alfa", "bravo", "charlie", "delta", "echo",
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
			"non-default parameters",
			PassphraseCountDefault + 1,
			PassphraseWordCountDefault + 1,
			PassphraseSeparatorDash,
			PassphraseCasingLower,
			alternateWordList,

			func(t *testing.T, passphrases []string, err error) {
				require.NoError(t, err)
				require.Len(t, passphrases, PassphraseCountDefault+1)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(PassphraseSeparatorDash))
					require.Len(t, words, PassphraseWordCountDefault+1)
					for _, word := range words {
						require.Contains(t, alternateWordList, strings.ToLower(word))
					}
				}
			},

			nil,
			nil,
		},
		{
			"non-default parameters",
			PassphraseCountDefault + 1,
			PassphraseWordCountDefault + 1,
			PassphraseSeparatorUnderscore,
			PassphraseCasingUpper,
			alternateWordList,

			func(t *testing.T, passphrases []string, err error) {
				require.NoError(t, err)
				require.Len(t, passphrases, PassphraseCountDefault+1)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(PassphraseSeparatorUnderscore))
					require.Len(t, words, PassphraseWordCountDefault+1)
					for _, word := range words {
						require.Contains(t, alternateWordList, strings.ToLower(word))
					}
				}
			},

			nil,
			nil,
		},
		{
			"non-default parameters",
			PassphraseCountDefault + 1,
			PassphraseWordCountDefault + 1,
			PassphraseSeparatorDot,
			PassphraseCasingTitle,
			alternateWordList,

			func(t *testing.T, passphrases []string, err error) {
				require.NoError(t, err)
				require.Len(t, passphrases, PassphraseCountDefault+1)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(PassphraseSeparatorDot))
					require.Len(t, words, PassphraseWordCountDefault+1)
					for _, word := range words {
						require.Contains(t, alternateWordList, strings.ToLower(word))
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

func BenchmarkGeneratePassphrases(b *testing.B) {
	type benchmarkDef struct {
		name      string
		count     uint
		wordCount uint
		separator rune
		casing    PassphraseCasing
		wordList  []string
	}

	var benchmarks = []benchmarkDef{
		// Preset value benchmarks.
		{
			"defaults",
			PassphraseCountDefault,
			PassphraseWordCountDefault,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"minimums",
			PassphraseCountMin,
			PassphraseWordCountMin,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			[]string{"a", "b"},
		},
		{
			"maximums",
			PassphraseCountMax,
			PassphraseWordCountMax,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},

		// Increasing count benchmarks.
		{
			"25_6_default",
			25,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"50_6_default",
			50,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"100_6_default",
			100,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"250_6_default",
			250,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"500_6_default",
			500,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},

		// Increasing word count benchmarks.
		{
			"16_8_default",
			16,
			8,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"16_16_default",
			16,
			16,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"16_32_default",
			16,
			32,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},
		{
			"16_64_default",
			16,
			64,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault,
		},

		// Increasing word list size benchmarks.
		{
			"16_6_25",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:25],
		},
		{
			"16_6_50",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:50],
		},
		{
			"16_6_100",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:100],
		},
		{
			"16_6_250",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:250],
		},
		{
			"16_6_500",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:500],
		},
		{
			"16_6_1000",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:1000],
		},
		{
			"16_6_2500",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:2500],
		},
		{
			"16_6_5000",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			WordListDefault[:5000],
		},

		// Special cases.
		{
			"once-duplicated word list",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			append(WordListDefault, WordListDefault...),
		},
		{
			"twice-duplicated word list",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			append(WordListDefault, append(WordListDefault, WordListDefault...)...),
		},
		{
			"thrice-duplicated word list",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingDefault,
			append(WordListDefault, append(WordListDefault, append(WordListDefault, WordListDefault...)...)...),
		},
		{
			"lowercase output",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingLower,
			WordListDefault,
		},
		{
			"uppercase output",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingUpper,
			WordListDefault,
		},
		{
			"title-case output",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingTitle,
			WordListDefault,
		},
		{
			"uncased output",
			16,
			6,
			PassphraseSeparatorDefault,
			PassphraseCasingNone,
			WordListDefault,
		},
	}

	for _, benchmark := range benchmarks {
		b.Run(
			benchmark.name,
			func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					_, _ = GeneratePassphrases(
						benchmark.count,
						benchmark.wordCount,
						benchmark.separator,
						benchmark.casing,
						benchmark.wordList,
					)
				}
			},
		)
	}
}
