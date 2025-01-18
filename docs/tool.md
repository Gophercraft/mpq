# mopaq

## Building

To get started, you'll need:

- [the Go programming language](https://go.dev/dl/) 
- git

```sh
git clone https://github.com/Gophercraft/mpq gophercraft_mpq
cd gophercraft_mpq
go install github.com/Gophercraft/mpq/cmd/mopaq
```

## Usage

First you need MPQs. Find a game folder ending with "Data"

The tool takes as an argument either

1) a single MPQ file

2) a chain file.

A chain file is a JSON array of MPQ glob paths represented as an array of strings. These MPQ paths can be relative or absolute (refer to [docs/wotlk-chain.json](./wotlk-chain.json)).

To generate a SHA-256 manifest of the MPQ set:

```sh
mopaq list docs/wotlk-chain.json \
    --working-directory /wotlk-game-dir/Data
```

To export the contents of an MPQ set into a directory:

```sh
mopaq export docs/wotlk-chain.json \
    --working-directory /wotlk-game-dir/Data \
    --export-directory /wotlk-game-dir/Export/
```