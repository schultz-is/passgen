package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/schultz-is/passgen"
	"github.com/spf13/cobra"
)

// buildPassphraseCmd constructs the passphrase subcommand responsible for generating passphrases.
func buildPassphraseCmd() *cobra.Command {
	// Build a configuration struct for converting commandline input into paramters for a passgen
	// GeneratePassphrases function call.
	passphraseConfig := struct {
		count     uint                     // Number of passphrases to generate.
		wordCount uint                     // Length, in words, of passphrases to generate.
		separator rune                     // Passphrase word separator.
		casing    passgen.PassphraseCasing // Passphrase word casing.
		wordList  []string                 // List of words to pull passphrase words from.

		separatorString string // Intermediate storage for separator flag before translation to rune.

		// At most one of the following values is allowed to be set.
		casingLower bool // Generate lowercase passphrases.
		casingUpper bool // Generate uppercase passphrases.
		casingTitle bool // Generate title-case passphrases.

		wordListFilename string // Filename of a newline-delimited word list to use in passphrases.
	}{
		passgen.PassphraseCountDefault,
		passgen.PassphraseWordCountDefault,
		passgen.PassphraseSeparatorDefault,
		passgen.PassphraseCasingDefault,
		passgen.WordListDefault,

		"",

		false,
		false,
		false,

		"",
	}

	// Construct the command.
	passphraseCmd := &cobra.Command{
		Use:   "passphrase [word count] [count]",
		Short: "Generate passphrases",

		Aliases: []string{
			"pp",
			"phrase",
		},

		Args: func(cmd *cobra.Command, args []string) error {
			// Don't allow more than two positional arguments (word count and count.)
			if len(args) > 2 {
				return errors.New("too many args provided")
			}

			// The first argument is the passphrase word count.
			if len(args) > 0 {
				wordCount, err := strconv.ParseUint(
					args[0],
					10,
					64,
				)
				if err != nil {
					return errors.New("invalid word count provided")
				}

				// Bounds check the word count for the platform.
				if wordCount > uint64(uintMax) {
					return errors.New("invalid word count provided")
				}

				// Update the configuration with parsed information.
				passphraseConfig.wordCount = uint(wordCount)
			}

			// The second argument is the passphrase count.
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

				// Update the configuration with parsed information.
				passphraseConfig.count = uint(count)
			}

			return nil
		},

		// Define what the passphrase subcommand does when invoked.
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Attempt to convert the provided separator string (if it exists) to a single rune.
			if passphraseConfig.separatorString != "" {
				if utf8.RuneCountInString(passphraseConfig.separatorString) > 1 {
					return errors.New("separator must be a single character")
				}
				passphraseConfig.separator = []rune(passphraseConfig.separatorString)[0]
			}

			// Set the casing based on user flags, ensuring no more than one flag is set.
			casingSet := false
			if passphraseConfig.casingLower {
				passphraseConfig.casing = passgen.PassphraseCasingLower
				casingSet = true
			}
			if passphraseConfig.casingUpper {
				if casingSet {
					return errors.New("at most one casing method is allowed")
				}
				passphraseConfig.casing = passgen.PassphraseCasingUpper
				casingSet = true
			}
			if passphraseConfig.casingTitle {
				if casingSet {
					return errors.New("at most one casing method is allowed")
				}
				passphraseConfig.casing = passgen.PassphraseCasingTitle
			}

			// Attempt to open and read the word list file if provided.
			if passphraseConfig.wordListFilename != "" {
				wordListFile, err := os.Open(passphraseConfig.wordListFilename)
				if err != nil {
					return err
				}

				// Clean up after the file has been read.
				defer wordListFile.Close()

				// Step through each line and append it to the word list.
				var wordList []string
				scanner := bufio.NewScanner(wordListFile)
				for scanner.Scan() {
					wordList = append(wordList, strings.TrimSpace(scanner.Text()))
				}

				passphraseConfig.wordList = wordList
			}

			// Generate passphrases based on the command invocation.
			passphrases, err := passgen.GeneratePassphrases(
				passphraseConfig.count,
				passphraseConfig.wordCount,
				passphraseConfig.separator,
				passphraseConfig.casing,
				passphraseConfig.wordList,
			)
			if err != nil {
				return err
			}

			// Print out a single passphrase per line.
			for _, passphrase := range passphrases {
				cmd.Println(passphrase)
			}

			return
		},
	}

	// Define the flag for the word separator.
	passphraseCmd.Flags().StringVarP(
		&passphraseConfig.separatorString,
		"separator",
		"s",
		"",
		"passphrase word separator",
	)

	// Define the flag for generating lowercase passphrases.
	passphraseCmd.Flags().BoolVarP(
		&passphraseConfig.casingLower,
		"lowercase",
		"l",
		false,
		"generate lowercase passphrases",
	)

	// Define the flag for generating uppercase passphrases.
	passphraseCmd.Flags().BoolVarP(
		&passphraseConfig.casingUpper,
		"uppercase",
		"u",
		false,
		"generate uppercase passphrases",
	)

	// Define the flag for generating title-case passphrases.
	passphraseCmd.Flags().BoolVarP(
		&passphraseConfig.casingTitle,
		"title-case",
		"t",
		false,
		"generate title-case passphrases",
	)

	// Define the flag for a word list filename.
	passphraseCmd.Flags().StringVarP(
		&passphraseConfig.wordListFilename,
		"word-list",
		"w",
		"",
		"file containing a newline-delimited word list for use in passphrases",
	)

	return passphraseCmd
}
