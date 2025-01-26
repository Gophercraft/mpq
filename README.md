![Logo](./docs/img/icon.png)

# Gophercraft/mpq

[![Go Reference](https://pkg.go.dev/badge/github.com/Gophercraft/mpq.svg)](https://pkg.go.dev/github.com/Gophercraft/mpq)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Chat on discord](https://img.shields.io/discord/556039662997733391.svg)](https://discord.gg/xPtuEjt)

This package provides an MPQ decoder compatible with MPQ versions 1-4. 

# `mopaq` command-line utility

This repository includes a command-line tool that you can use to inspect an MPQ archive's contents, extract a file, or a large set of files.

See [mopaq.md](docs/mopaq.md) to read about how to install and use the tool.

# Usage

Start by opening an Archive with `mpq.Open()`:

```go
// open an Archive
archive, err := mpq.Open("sample/test.SC2Replay")
// handle err
```

Now you can read files like this:

```go
// File implements io.ReadCloser
file, err := archive.Open("Path\\to\\file")
// handle err

// read all bytes from File
file_bytes, err := io.ReadAll(file)
// handle err

// close 
err = file.Close()
// handle err
```

You can make use of the listfile that most MPQs contain with the `List()` method:

```go
list, err := archive.List()
// handle err

for list.Next() {
    fmt.Println(list.Path())
}

list.Close()
```

You can also load multiple Archives on top of eachother, using the `mpq.Set` structure.

```go
set, err := mpq.GlobSet(
    "common.MPQ",
    "*/locale-*.MPQ",
    "*/speech-*.MPQ")
// handle err

// Open a file from a Set, favoring the highest-order Archives
file, err := set.Open("Path\\to\\file")

// List files from a Set
list, err := set.List()
```

# Thanks

Documentation and logic inspired by [StormLib](https://github.com/ladislav-zezula/StormLib) by Ladislav Zezula, the gold standard for MPQ parsing

Cryptographic routines taken from [go.Zamara](https://github.com/aphistic/go.Zamara) by Kristin Davidson

All glory to Mike O'Brien.