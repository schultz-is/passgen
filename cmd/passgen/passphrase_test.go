package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/schultz-is/passgen"
	"github.com/stretchr/testify/require"
)

func TestPassphraseCommand(t *testing.T) {
	type testReqs func(t *testing.T, output string, err error)

	type testDef struct {
		name  string
		args  []string
		flags map[string]string

		requirements testReqs
		setup        func() interface{}
		teardown     func(interface{})
	}

	// Construct a custom word list and write it to a temp file.
	var alternateWordList = []string{
		"alfa", "bravo", "charlie", "delta", "echo",
	}

	wordListFile, err := ioutil.TempFile("", "")

	defer func() {
		_ = wordListFile.Close()
	}()

	require.NoError(t, err)
	wordListFilename := wordListFile.Name()
	for _, word := range alternateWordList {
		_, err = wordListFile.WriteString(word + "\n")
		require.NoError(t, err)
	}

	var tests = []testDef{
		{
			"rational defaults",
			nil,
			nil,

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"too many arguments",
			[]string{"1", "2", "3"},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"mistyped word count argument",
			[]string{"x"},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word count argument too small",
			[]string{strconv.FormatUint(passgen.PassphraseWordCountMin-1, 10)},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word count argument too large",
			[]string{strconv.FormatUint(passgen.PassphraseWordCountMax+1, 10)},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"word count argument too large for platform",
			[]string{"16"},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			func() interface{} {
				originalUintMax := uintMax
				uintMax = 15
				return originalUintMax
			},

			func(setupContext interface{}) {
				uintMax = setupContext.(uint)
			},
		},
		{
			"mistyped count argument",
			[]string{
				strconv.FormatUint(passgen.PassphraseWordCountDefault, 10),
				"x",
			},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"count argument too small",
			[]string{
				strconv.FormatUint(passgen.PassphraseWordCountDefault, 10),
				strconv.FormatUint(passgen.PassphraseCountMin-1, 10),
			},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"count argument too large",
			[]string{
				strconv.FormatUint(passgen.PassphraseWordCountDefault, 10),
				strconv.FormatUint(passgen.PassphraseCountMax+1, 10),
			},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"count argument too large for platform",
			[]string{
				"1",
				"16",
			},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			func() interface{} {
				originalUintMax := uintMax
				uintMax = 15
				return originalUintMax
			},

			func(setupContext interface{}) {
				uintMax = setupContext.(uint)
			},
		},
		{
			"invalid separator",
			nil,
			map[string]string{
				"separator": "==",
			},

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"custom separator",
			nil,
			map[string]string{
				"separator": string(passgen.PassphraseSeparatorDash),
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDash))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"too many casing flags",
			nil,
			map[string]string{
				"lowercase": "true",
				"uppercase": "true",
			},

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"too many casing flags",
			nil,
			map[string]string{
				"uppercase":  "true",
				"title-case": "true",
			},

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"too many casing flags",
			nil,
			map[string]string{
				"title-case": "true",
				"no-casing":  "true",
			},

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"lowercase flag",
			nil,
			map[string]string{
				"lowercase": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, strings.ToLower(word))
						require.Equal(t, strings.ToLower(word), word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"uppercase flag",
			nil,
			map[string]string{
				"uppercase": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, strings.ToLower(word))
						require.Equal(t, strings.ToUpper(word), word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"title-case flag",
			nil,
			map[string]string{
				"title-case": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, strings.ToLower(word))
						require.Equal(t, strings.Title(word), word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"no-casing flag",
			nil,
			map[string]string{
				"no-casing": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, passgen.WordListDefault, word)
					}
				}
			},

			nil,
			nil,
		},
		{
			"nonexistent word list file",
			nil,
			map[string]string{
				"word-list": "fake.txt",
			},

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"custom word list",
			nil,
			map[string]string{
				"word-list": wordListFilename,
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passphrases := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passphrases, passgen.PassphraseCountDefault)
				for _, passphrase := range passphrases {
					words := strings.Split(passphrase, string(passgen.PassphraseSeparatorDefault))
					require.Len(t, words, passgen.PassphraseWordCountDefault)
					for _, word := range words {
						require.Contains(t, alternateWordList, word)
					}
				}
			},

			nil,
			nil,
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

				passphraseCmd := buildPassphraseCmd()
				var outputBuffer strings.Builder
				passphraseCmd.SetOut(&outputBuffer)

				passphraseCmd.SetArgs(test.args)
				for flag, value := range test.flags {
					err := passphraseCmd.Flags().Set(flag, value)
					require.NoError(t, err)
				}

				err := passphraseCmd.Execute()
				test.requirements(t, outputBuffer.String(), err)

				if test.teardown != nil {
					test.teardown(setupContext)
				}
			},
		)
	}
}
