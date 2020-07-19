# passgen

![Builds](https://github.com/schultz-is/passgen/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/schultz-is/passgen)](https://goreportcard.com/report/github.com/schultz-is/passgen)
[![GoDoc](https://godoc.org/github.com/schultz-is/passgen?status.svg)](https://pkg.go.dev/github.com/schultz-is/passgen)
[![License](https://img.shields.io/github/license/schultz-is/passgen)](./LICENSE)

`passgen` is an API and command-line utility for generating passwords and passphrases.

## Installation

To install the command-line utility:
```console
git clone -b stable --depth 1 https://github.com/schultz-is/passgen.git
cd passgen
make install
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
> passgen pw 5 10 --alphabet "ACGT"
ATAAG
CATTC
TTGAT
CGGAT
TGTAG
GCGAC
ACATG
TTATT
ACTAT
CCGTA
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
> passgen pp 5 6 -uw words.txt
TANNING TRICKLE PRECOOK KEEP ARMHOLE
MARITIME LADYLIKE ELM UNDRAFTED BONANZA
EGOTISM MANTIS BANNER MUNICIPAL AMUSING
EVOLUTION WIRING TRACK BLURT GREYHOUND
UNTITLED RURAL SHAKINESS GEOMETRIC ARMREST
WHY OUTCLASS RIVETING OVERLORD UNFIXED
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

```console
> make test
> make cover
```

## Benchmarks

```console
> make benchmark
> go tool pprof prof/cpu.prof
> go tool pprof prof/mem.prof
```

## Build

```console
> make build
> ./dist/passgen --version
```
