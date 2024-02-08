# gophercraft_mpq_set

## Building

To get started, you'll need:

- [the Go programming language](https://go.dev/dl/) 
- git

```
git clone https://github.com/Gophercraft/mpq gophercraft_mpq
cd gophercraft_mpq
go build github.com/Gophercraft/mpq/cmd/gophercraft_mpq_set
```

## Usage

First you need MPQs. Find a game folder ending with "Data"

Second, you need a list of MPQ glob paths represented as an array of strings in a JSON file. These MPQ paths can be relative or absolute (refer to [docs/wotlk-chain.json](./wotlk-chain.json)).

To generate a SHA-256 manifest of the MPQ set:

```bash
./gophercraft_mpq_set list \
    --chain-json docs/wotlk-chain.json \
    --working-directory /wotlk-game-dir/Data
```

To export the contents of an MPQ set into a directory:

```bash
./gophercraft_mpq_set export \
    --chain-json docs/wotlk-chain.json \
    --working-directory /wotlk-game-dir/Data \
    --export-directory /wotlk-game-dir/Export/
```