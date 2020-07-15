package main

import (
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/schultz-is/passgen"
	"github.com/stretchr/testify/require"
)

func TestPasswordCommand(t *testing.T) {
	type testReqs func(t *testing.T, output string, err error)

	type testDef struct {
		name  string
		args  []string
		flags map[string]string

		requirements testReqs
		setup        func() interface{}
		teardown     func(interface{})
	}

	var tests = []testDef{
		{
			"rational defaults",
			nil,
			nil,

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passwords := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passwords, passgen.PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, passgen.PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(t, passgen.AlphabetDefault, string(char))
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
			"mistyped length argument",
			[]string{"x"},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"length argument too small",
			[]string{strconv.FormatUint(passgen.PasswordLengthMin-1, 10)},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"length argument too large",
			[]string{strconv.FormatUint(passgen.PasswordLengthMax+1, 10)},
			nil,

			func(t *testing.T, output string, err error) {
				require.Error(t, err)
			},

			nil,
			nil,
		},
		{
			"length argument too large for platform",
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
				strconv.FormatUint(passgen.PasswordLengthDefault, 10),
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
				strconv.FormatUint(passgen.PasswordLengthDefault, 10),
				strconv.FormatUint(passgen.PasswordCountMin-1, 10),
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
				strconv.FormatUint(passgen.PasswordLengthDefault, 10),
				strconv.FormatUint(passgen.PasswordCountMax+1, 10),
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
			"all characters allowed in output, excluding ambiguous",
			nil,
			map[string]string{
				"lowercase": "true",
				"uppercase": "true",
				"numeric":   "true",
				"special":   "true",
				"ambiguous": "false",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passwords := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passwords, passgen.PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, passgen.PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(
							t,
							passgen.AlphabetLower+passgen.AlphabetUpper+passgen.AlphabetNumeric+passgen.AlphabetSpecial,
							string(char),
						)
					}
				}
			},

			nil,
			nil,
		},
		{
			"all characters allowed in output, including ambiguous",
			nil,
			map[string]string{
				"lowercase": "true",
				"uppercase": "true",
				"numeric":   "true",
				"special":   "true",
				"ambiguous": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passwords := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passwords, passgen.PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, passgen.PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(
							t,
							passgen.AlphabetLowerAmbiguous+passgen.AlphabetUpperAmbiguous+passgen.AlphabetNumericAmbiguous+passgen.AlphabetSpecial,
							string(char),
						)
					}
				}
			},

			nil,
			nil,
		},
		{
			"default alphabet allowed in output, including ambiguous",
			nil,
			map[string]string{
				"ambiguous": "true",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passwords := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passwords, passgen.PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, passgen.PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.Contains(
							t,
							passgen.AlphabetDefaultAmbiguous,
							string(char),
						)
					}
				}
			},

			nil,
			nil,
		},
		{
			"alphabet flag supersedes all other flags",
			nil,
			map[string]string{
				"lowercase": "true",
				"uppercase": "true",
				"numeric":   "true",
				"special":   "true",
				"ambiguous": "true",
				"alphabet":  "🐶🐱",
			},

			func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				passwords := strings.Split(strings.TrimSpace(output), "\n")
				require.Len(t, passwords, passgen.PasswordCountDefault)
				for _, password := range passwords {
					require.Equal(t, passgen.PasswordLengthDefault, utf8.RuneCountInString(password))
					for _, char := range password {
						require.NotContains(
							t,
							passgen.AlphabetLowerAmbiguous+passgen.AlphabetUpperAmbiguous+passgen.AlphabetNumericAmbiguous+passgen.AlphabetSpecial,
							string(char),
						)
						require.Contains(
							t,
							"🐶🐱",
							string(char),
						)
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

				passwordCmd := buildPasswordCmd()
				var outputBuffer strings.Builder
				passwordCmd.SetOut(&outputBuffer)

				passwordCmd.SetArgs(test.args)
				for flag, value := range test.flags {
					passwordCmd.Flags().Set(flag, value)
				}

				err := passwordCmd.Execute()
				test.requirements(t, outputBuffer.String(), err)

				if test.teardown != nil {
					test.teardown(setupContext)
				}
			},
		)
	}
}
