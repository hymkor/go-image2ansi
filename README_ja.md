go-image2ansi
=============
( [English](README.md) / Japanese )

指定した画像をターミナルにテキストとして出力します。

- Sixel Graphics を使わず、Unicode の U+2580 ( &#x2580; ) と ANSIカラーシーケンスを用いて、テキストの１キャラクターに２ピクセルを表示させます。
- Windows Terminal, Command Prompt, ConEmu, WezTerm, GNOME Terminal で動作確認しています。Unicode のブロック文字と、ANSI カラーシーケンスをサポートする端末であれば、理論上動作します。

Install
-------

```
go install github.com/hymkor/go-image2ansi/cmd/image2ansi@latest
```

Usage
-----

```
image2ansi {オプション...} 画像ファイル名...
```

### オプション

- `-24` シーケンス `ESC[38;2;{r1};{g1};{b1}m` `ESC[48;2;{r2};{g2};{b2}m`を使って、
   24bit color で出力します。
   省略時は `ESC[38;5;{id1}m` `ESC[48;5;{id2}m`を使って 256色で出力します。
- `-w WIDTH` 出力する幅を指定します。
   省略時は現在のターミナルの桁数を使用します。
   リダイレクトされている時は幅を制限せず、画像そのものの横幅をそのまま使用します。
- `-gif-from-stdin` 標準入力から GIF 画像を読み込みます。
- `-jpg-from-stdin` 標準入力から JPEG 画像を読み込みます。
- `-png-from-stdin` 標準入力から PNG 画像を読み込みます。

### 画像ファイル名

.GIF、.PNG、.JPG(.JPEG) が利用可能です。