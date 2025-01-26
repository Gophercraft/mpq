# mopaq

## Building

To get started, you'll need:

- [the Go programming language](https://go.dev/dl/) 
- git

```sh
git clone https://github.com/Gophercraft/mpq mpq && cd mpq
go install github.com/Gophercraft/mpq/cmd/mopaq
```

## Usage

### export

The `export` command mass-exports one or more MPQs into a directory, specified by the `--export-directory` flag.

```bash
mopaq export sample/test.SC2Replay --export-directory /tmp/test-sc2-replay
```

If desired, rather than specifying a list of MPQs manually, you can load a JSON list of MPQs from a file.

```bash
mopaq export wotlk-chain.json --export-directory /tmp/wotlk/
```

### list

The `list` command lists the contents of one or more MPQs.

```bash
mopaq list sample/test.SC2Replay
```

If desired, the `-a` / `--hash-algorithm` flag will display the filenames along with hash checksums.

```bash
mopaq list sample/test.SC2Replay -a sha256
```

### info

The `info` command outputs diagnostic information about a single MPQ archive.

```bash
mopaq info sample/test.SC2Replay
```