package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/schultz-is/passgen"
	"github.com/spf13/cobra"
)

// buildAlphabet constructs an alphabet string based on character type flags.
func buildAlphabet(allowLowercase, allowUppercase, allowNumeric, allowSpecial, allowAmbiguous bool) string {
	type charSet struct {
		enabled           bool
		normalAlphabet    string
		ambiguousAlphabet string
	}

	sets := []charSet{
		{allowLowercase, passgen.AlphabetLower, passgen.AlphabetLowerAmbiguous},
		{allowUppercase, passgen.AlphabetUpper, passgen.AlphabetUpperAmbiguous},
		{allowNumeric, passgen.AlphabetNumeric, passgen.AlphabetNumericAmbiguous},
	}

	var b strings.Builder
	for _, set := range sets {
		if set.enabled {
			if allowAmbiguous && set.ambiguousAlphabet != "" {
				b.WriteString(set.ambiguousAlphabet)
			} else {
				b.WriteString(set.normalAlphabet)
			}
		}
	}

	if allowSpecial {
		b.WriteString(passgen.AlphabetSpecial)
	}

	alphabet := b.String()
	if alphabet == "" {
		if allowAmbiguous {
			return passgen.AlphabetDefaultAmbiguous
		}
		return passgen.AlphabetDefault
	}

	return alphabet
}

// buildPasswordCmd constructs the password subcommand responsible for generating passwords.
func buildPasswordCmd() *cobra.Command {
	// Build a configuration struct for converting commandline input into parameters for a passgen
	// GeneratePasswords function call.
	passwordConfig := struct {
		count    uint   // Number of passwords to generate.
		length   uint   // Length of passwords to generate.
		alphabet string // Alphabet to use when generating passwords.

		allowUppercase bool // Allow uppercase characters in passwords.
		allowLowercase bool // Allow lowercase characters in passwords.
		allowNumeric   bool // Allow numeric characters in passwords.
		allowSpecial   bool // Allow special characters in passwords.
		allowAmbiguous bool // Allow ambiguous characters in passwords.
	}{
		passgen.PasswordCountDefault,
		passgen.PasswordLengthDefault,
		passgen.AlphabetDefault,

		false,
		false,
		false,
		false,
		false,
	}

	// Construct the command.
	passwordCmd := &cobra.Command{
		Use:   "password [length] [count]",
		Short: "Generate passwords",

		Aliases: []string{
			"pw",
			"word",
		},

		Args: func(cmd *cobra.Command, args []string) error {
			// Don't allow more than two positional arguments (length and count.)
			if len(args) > 2 {
				return errors.New("too many args provided")
			}

			// Parse the first argument (password length) if provided.
			if len(args) > 0 {
				length, err := parseUintArg(args, 0, "length")
				if err != nil {
					return err
				}
				passwordConfig.length = length
			}

			// Parse the second argument (password count) if provided.
			if len(args) > 1 {
				count, err := parseUintArg(args, 1, "count")
				if err != nil {
					return err
				}
				passwordConfig.count = count
			}

			return nil
		},

		// Define what the password subcommand does when invoked.
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Determine the alphabet to use.
			if passwordConfig.alphabet == "" {
				passwordConfig.alphabet = buildAlphabet(
					passwordConfig.allowLowercase,
					passwordConfig.allowUppercase,
					passwordConfig.allowNumeric,
					passwordConfig.allowSpecial,
					passwordConfig.allowAmbiguous,
				)
			}

			// Generate passwords based on the command invocation.
			passwords, err := passgen.GeneratePasswords(
				passwordConfig.count,
				passwordConfig.length,
				passwordConfig.alphabet,
			)
			if err != nil {
				return err
			}

			// Print out a single password per line.
			for _, password := range passwords {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), password)
				if err != nil {
					return err
				}
			}

			return
		},
	}

	// Define the flag for allowance of lowercase characters in generated passwords.
	passwordCmd.Flags().BoolVarP(
		&passwordConfig.allowLowercase,
		"lowercase",
		"l",
		false,
		"allow lowercase letters in passwords",
	)

	// Define the flag for allowance of uppercase characters in generated passwords.
	passwordCmd.Flags().BoolVarP(
		&passwordConfig.allowUppercase,
		"uppercase",
		"u",
		false,
		"allow uppercase letters in passwords",
	)

	// Define the flag for allowance of numeric characters in generated passwords.
	passwordCmd.Flags().BoolVarP(
		&passwordConfig.allowNumeric,
		"numeric",
		"n",
		false,
		"allow numeric characters in passwords",
	)

	// Define the flag for allowance of special characters in generated passwords.
	passwordCmd.Flags().BoolVarP(
		&passwordConfig.allowSpecial,
		"special",
		"s",
		false,
		"allow special characters in passwords",
	)

	// Define the flag for allowance of ambiguous characters in generated passwords.
	passwordCmd.Flags().BoolVarP(
		&passwordConfig.allowAmbiguous,
		"ambiguous",
		"a",
		false,
		"allow ambiguous characters in passwords",
	)

	// Define the flag for specification of a custom alphabet used by generated passwords.
	passwordCmd.Flags().StringVar(
		&passwordConfig.alphabet,
		"alphabet",
		"",
		"alphabet to use for password generation (supersedes other flags)",
	)

	return passwordCmd
}
