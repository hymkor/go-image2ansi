go-image2ansi
=============
( English / [Japanese](./README_ja.md) )

Display images as text in a terminal.

- Uses Unicode U+2580 ( &#x2580; ) and ANSI color sequences instead of Sixel Graphics. Each character cell represents two vertical pixels.
- Tested with Windows Terminal, Command Prompt, ConEmu,  WezTerm, and GNOME Terminal. In theory, it should work on any terminal that supports Unicode block characters and ANSI color sequences.

Install
-------

```sh
go install github.com/hymkor/go-image2ansi/cmd/image2ansi@latest
```

## Usage

```text
image2ansi [options...] image-file...
```

### Options

- `-24`  
  Use 24-bit color with the following sequences:
  - `ESC[38;2;{r1};{g1};{b1}m`
  - `ESC[48;2;{r2};{g2};{b2}m`

  By default, 256 colors are used:
  - `ESC[38;5;{id1}m`
  - `ESC[48;5;{id2}m`
- `-w WIDTH`  
  Specify the output width.
  By default, the current terminal width is used. When output is redirected, no width limit is applied, and the original image width is used.
- `-gif-from-stdin` Read a GIF image from standard input
- `-jpg-from-stdin` Read a JPEG image from standard input
- `-png-from-stdin` Read a PNG image from standard input

### Supported Image Formats

- GIF
- PNG
- JPEG