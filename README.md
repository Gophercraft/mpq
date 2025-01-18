![Logo](./docs/img/icon.png)

# Gophercraft/mpq

[![Go Reference](https://pkg.go.dev/badge/github.com/Gophercraft/mpq.svg)](https://pkg.go.dev/github.com/Gophercraft/mpq)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Chat on discord](https://img.shields.io/discord/556039662997733391.svg)](https://discord.gg/xPtuEjt)

This is an MPQ decoder compatible with MPQ v1 and v2. Support for later versions will be added over time.

# Tool

See [tool.md](docs/tool.md) to read about how to use the included extraction tool.

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