package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/schultz-is/passgen"
	"github.com/spf13/cobra"
)

// parseSeparator converts a string to a single rune for use as a separator.
func parseSeparator(separatorString string) (rune, error) {
	if separatorString == "" {
		return 0, nil // No separator provided
	}
	if utf8.RuneCountInString(separatorString) > 1 {
		return 0, errors.New("separator must be a single character")
	}
	return []rune(separatorString)[0], nil
}

// passphraseConfig holds the configuration for passphrase generation from CLI.
type passphraseConfig struct {
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
	casingNone  bool // Generate passphrases without applying any casing transformation.

	wordListFilename string // Filename of a newline-delimited word list to use in passphrases.
}

// runPassphraseCmd executes the passphrase generation command.
func runPassphraseCmd(cfg *passphraseConfig, output io.Writer) error {
	// Parse separator if provided.
	if cfg.separatorString != "" {
		sep, err := parseSeparator(cfg.separatorString)
		if err != nil {
			return err
		}
		cfg.separator = sep
	}

	// Set the casing based on user flags, ensuring no more than one flag is set.
	casing, err := selectCasing(cfg.casingLower, cfg.casingUpper, cfg.casingTitle, cfg.casingNone)
	if err != nil {
		return err
	}
	cfg.casing = casing

	// Read the word list file if provided.
	if cfg.wordListFilename != "" {
		wordList, err := readWordListFile(cfg.wordListFilename)
		if err != nil {
			return err
		}
		cfg.wordList = wordList
	}

	// Generate passphrases based on the command invocation.
	passphrases, err := passgen.GeneratePassphrases(
		cfg.count,
		cfg.wordCount,
		cfg.separator,
		cfg.casing,
		cfg.wordList,
	)
	if err != nil {
		return err
	}

	// Print out a single passphrase per line.
	for _, passphrase := range passphrases {
		if _, err := fmt.Fprintln(output, passphrase); err != nil {
			return err
		}
	}

	return nil
}

// selectCasing determines which casing to use, ensuring mutual exclusion of casing flags.
func selectCasing(lower, upper, title, none bool) (passgen.PassphraseCasing, error) {
	casings := []struct {
		enabled bool
		value   passgen.PassphraseCasing
	}{
		{lower, passgen.PassphraseCasingLower},
		{upper, passgen.PassphraseCasingUpper},
		{title, passgen.PassphraseCasingTitle},
		{none, passgen.PassphraseCasingNone},
	}

	var selected *passgen.PassphraseCasing
	for _, c := range casings {
		if c.enabled {
			if selected != nil {
				return 0, errors.New("at most one casing method is allowed")
			}
			val := c.value
			selected = &val
		}
	}

	if selected == nil {
		return passgen.PassphraseCasingDefault, nil
	}
	return *selected, nil
}

// readWordListFile reads a newline-delimited word list from a file.
func readWordListFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		words = append(words, strings.TrimSpace(scanner.Text()))
	}
	return words, nil
}

// buildPassphraseCmd constructs the passphrase subcommand responsible for generating passphrases.
func buildPassphraseCmd() *cobra.Command {
	// Build a configuration struct for converting commandline input into parameters for a passgen
	// GeneratePassphrases function call.
	cfg := &passphraseConfig{
		count:     passgen.PassphraseCountDefault,
		wordCount: passgen.PassphraseWordCountDefault,
		separator: passgen.PassphraseSeparatorDefault,
		casing:    passgen.PassphraseCasingDefault,
		wordList:  passgen.WordListDefault,
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

			// Parse the first argument (passphrase word count) if provided.
			if len(args) > 0 {
				wordCount, err := parseUintArg(args, 0, "word count")
				if err != nil {
					return err
				}
				cfg.wordCount = wordCount
			}

			// Parse the second argument (passphrase count) if provided.
			if len(args) > 1 {
				count, err := parseUintArg(args, 1, "count")
				if err != nil {
					return err
				}
				cfg.count = count
			}

			return nil
		},

		// Define what the passphrase subcommand does when invoked.
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPassphraseCmd(cfg, cmd.OutOrStdout())
		},
	}

	// Define the flag for the word separator.
	passphraseCmd.Flags().StringVarP(
		&cfg.separatorString,
		"separator",
		"s",
		"",
		"passphrase word separator",
	)

	// Define the flag for generating lowercase passphrases.
	passphraseCmd.Flags().BoolVarP(
		&cfg.casingLower,
		"lowercase",
		"l",
		false,
		"generate lowercase passphrases",
	)

	// Define the flag for generating uppercase passphrases.
	passphraseCmd.Flags().BoolVarP(
		&cfg.casingUpper,
		"uppercase",
		"u",
		false,
		"generate uppercase passphrases",
	)

	// Define the flag for generating title-case passphrases.
	passphraseCmd.Flags().BoolVarP(
		&cfg.casingTitle,
		"title-case",
		"t",
		false,
		"generate title-case passphrases",
	)

	// Define the flag for generating uncased passphrases.
	passphraseCmd.Flags().BoolVarP(
		&cfg.casingNone,
		"no-casing",
		"n",
		false,
		"generate passphrases without applying any case transformation",
	)

	// Define the flag for a word list filename.
	passphraseCmd.Flags().StringVarP(
		&cfg.wordListFilename,
		"word-list",
		"w",
		"",
		"file containing a newline-delimited word list for use in passphrases",
	)

	return passphraseCmd
}
