package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/schultz-is/passgen"
	"github.com/spf13/cobra"
)

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

			// The first argument is the password length.
			if len(args) > 0 {
				length, err := strconv.ParseUint(
					args[0],
					10,
					64,
				)
				if err != nil {
					return errors.New("invalid length provided")
				}

				// Bounds check the length for the platform.
				if length > uint64(uintMax) {
					return errors.New("invalid length provided")
				}

				// Update the configuration with the parsed information.
				passwordConfig.length = uint(length)
			}

			// The second argument is the password count.
			if len(args) > 1 {
				count, err := strconv.ParseUint(
					args[1],
					10,
					64,
				)
				if err != nil {
					return errors.New("invalid count provided")
				}

				// Bounds check the count for the platform.
				if count > uint64(uintMax) {
					return errors.New("invalid count provided")
				}

				// Update the configuration with the parsed information.
				passwordConfig.count = uint(count)
			}

			return nil
		},

		// Define what the password subcommand does when invoked.
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Determine the alphabet to use.
			if passwordConfig.alphabet == "" {
				// Instantiate a string builder for efficient alphabet construction.
				var b strings.Builder

				// Determine if the user wants passwords which include lowercase characters.
				if passwordConfig.allowLowercase {
					// Determine whether or not ambiguous characters should be allowed in the output.
					if passwordConfig.allowAmbiguous {
						_, err = b.WriteString(passgen.AlphabetLowerAmbiguous)
						if err != nil {
							return err
						}
					} else {
						_, err = b.WriteString(passgen.AlphabetLower)
						if err != nil {
							return err
						}
					}
				}

				// Determine if the user wants passwords which include uppercase characters.
				if passwordConfig.allowUppercase {
					// Determine whether or not ambiguous characters should be allowed in the output.
					if passwordConfig.allowAmbiguous {
						_, err = b.WriteString(passgen.AlphabetUpperAmbiguous)
						if err != nil {
							return err
						}
					} else {
						_, err = b.WriteString(passgen.AlphabetUpper)
						if err != nil {
							return err
						}
					}
				}

				// Determine if the user wants passwords which include numeric characters.
				if passwordConfig.allowNumeric {
					// Determine whether or not ambiguous characters should be allowed in the output.
					if passwordConfig.allowAmbiguous {
						_, err = b.WriteString(passgen.AlphabetNumericAmbiguous)
						if err != nil {
							return err
						}
					} else {
						_, err = b.WriteString(passgen.AlphabetNumeric)
						if err != nil {
							return err
						}
					}
				}

				// Determine if the user wants passwords which include numeric characters.
				if passwordConfig.allowSpecial {
					_, err := b.WriteString(passgen.AlphabetSpecial)
					if err != nil {
						return err
					}
				}

				// Build the alphabet based on user input.
				passwordConfig.alphabet = b.String()

				// If the user did not specify allowance of any of lowercase, uppercase, numeric, or
				// special characters, rely on the default alphabet.
				if passwordConfig.alphabet == "" {
					// Determine whether or not ambiguous characters should be allowed in the output.
					if passwordConfig.allowAmbiguous {
						passwordConfig.alphabet = passgen.AlphabetDefaultAmbiguous
					} else {
						passwordConfig.alphabet = passgen.AlphabetDefault
					}
				}
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
				fmt.Println(password)
			}

			return nil
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
