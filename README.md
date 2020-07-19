# passgen

![Builds](https://github.com/schultz-is/passgen/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/schultz-is/passgen)](https://goreportcard.com/report/github.com/schultz-is/passgen)
[![GoDoc](https://godoc.org/github.com/schultz-is/passgen?status.svg)](https://pkg.go.dev/github.com/schultz-is/passgen)
[![License](https://img.shields.io/github/license/schultz-is/passgen)](./LICENSE)

`passgen` is an API and command-line utility for generating passwords and passphrases.

## Installation

To install the command-line utility:
```console
go install github.com/schultz-is/passgen/cmd/passgen
```

To install the API for use in other projects:
```console
go get github.com/schultz-is/passgen
```

## Examples

### Using the command-line utility
```console
> passgen password
vt7tStRf3SfLV3V3
```

```console
> passgen pw -alnsu
Bc!Eyca9pHmWuRJr
```

```console
> passgen pw --alphabet "ACGT"
CTGCAGTCAAGGTTTG
```

```console
> passgen passphrase
faceless navigate scabby return snorkel cough
```

```console
> passgen pp -ts.
Cranberry.Deskwork.Ramble.Energize.Gloss.Tranquil
```

```console
> passgen pp -uw words.txt
MUSHROOM CHAMPION POD CHAFE SUITABLY EMPLOYER
```

### Using the API
```go
package main

import (
	"fmt"

	"github.com/schultz-is/passgen"
)

func main() {
	passwords, err := passgen.GeneratePasswords(
		passgen.PasswordCountDefault,
		passgen.PasswordLengthDefault,
		passgen.AlphabetDefault,
	)
	if err != nil {
		panic(err)
	}

	for _, password := range passwords {
		fmt.Println(password)
	}
}
```
[Open in Go Playground](https://play.golang.org/p/H45Sord6t0v)

```go
package main

import (
	"fmt"

	"github.com/schultz-is/passgen"
)

func main() {
	passphrases, err := passgen.GeneratePassphrases(
		passgen.PassphraseCountDefault,
		passgen.PassphraseWordCountDefault,
		passgen.PassphraseSeparatorDefault,
		passgen.PassphraseCasingDefault,
		passgen.WordListDefault,
	)
	if err != nil {
		panic(err)
	}

	for _, passphrase := range passphrases {
		fmt.Println(passphrase)
	}
}
```
[Open in Go Playground](https://play.golang.org/p/I-t1GM0QjUy)

## Tests

### Running tests and generating a coverage report
```console
> make test
```

### Viewing unit test coverage
```console
> make cover
```

## Benchmarks

### Running benchmarks and generating CPU and memory profiles
```console
make benchmark
```

### Viewing CPU and memory profiles
```console
> go tool pprof prof/cpu.prof
> go tool pprof prof/mem.prof
```
