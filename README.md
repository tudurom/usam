# usam

## micro-sam

An experiment in blending sam with the shell.

After the experiment is done, the plan is to build a shell around structural
regular expressions.

## Mode of operation

Each sam command is implemented as a stand-alone command line tool.

The commands are chained together via pipes.

Many commands accept a new dot as an argument. The argument dot is evaluated relative to the current dot.
Works for loops (`x` commands) too.

## Tools

### `e <filename>`

Opens the file for editing. Dot is set to `0`. It is the start of the pipe chain.

Example:

```bash
e file.txt | further editing
```

### `po` and `pc`

If you want to manipulate text from a pipe,
you can use `po` (pipe open) and `pc` (pipe close).

`pc` writes the resulting text to stdout.

`po` reads all the pipe's content before continuing the command chain.

Example:

```bash
dmesg | po | further editing | pc > new_dmesg.txt
```

### `el <new_dot>`

Sets the dot to a new address. When used in a loop, it ends the loop and sets the dot relative to the last dot of the loop.

Example:

```bash
e file.txt | el 2,3 | p
```

### `p [dot]`

Prints the content of the dot.

Example:

```bash
e file.txt | p 2,3 # does the exact thing as the above example
```

### `c <text> [dot]`

Changes the dot's value with `text`. Sets dot to the changed text.

Example:

```bash
<<EOF > file.txt
Happy 2018!
EOF

e file.txt | c 2019 '/2018/' | p ,
# Prints: Happy 2019!
```

### `i <text> [dot]`

Like `c`, but inserts `text` right before the dot.

### `a <text> [dot]`

Like `c`, but inserts `text` right after the dot.

### `d [dot]`

Deletes the dot's content.

Example:

```bash
<<EOF > file.txt
We live in a different society.
EOF

e file.txt | d '/ different/' | p ,
# Prints: We live in a society.
```

### `s <regexp> <text> [n|g] [dot]`

Substitute `text` for the first match to the `regexp` in the dot. Set dot to the modified range.
In `text`, `$` signs are interpreted as in Go's [`Expand`](https://golang.org/pkg/regexp/#Regexp.Expand).
If you want to change dot and substitute the first match, you must call it like so:

```bash
s <regexp> <text> 1 <dot>
```

Example:

```bash
<<EOF > file.txt
y y a x y y
EOF

e file.txt | s '(a) (x)' '$2 $1' 1 1 | p ,
# Prints: y y x a y y
```

### `x [regexp]`

For each `regexp` match in the dot, sets dot to that and execute the next command on that dot.

Because `x` doesn't accept a dot argument, you must use `el` first. Ironically, `el` comes from "end loop".

You can also compose `x` commands.

Example:

```bash
e vim.1 | el , | x 'vim' | p +-
# Prints all lines that contain the word 'vim'.
# If a line has 'vim' in it more than once, the line will pe printed each time.
```
