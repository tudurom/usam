# usam

## micro-sam

An experiment in blending sam with the shell.

## Mode of operation

Each sam command is implemented as a stand-alone command line tool.

The commands are chained together via pipes.

## Tools

### `e <filename>`

Opens the file for editing. Dot is set to `0`. It is the start of the pipe chain.

Example:

```bash
e file.txt | further editing
```

### `ca <new_dot>`

Sets the dot to a new address.

Example:

```bash
e file.txt | ca 2,3 | p
```

### `p`

Prints the content of the dot.

### `c <text>`

Changes the dot's value with `text`.

Example:

```bash
<<EOF > file.txt
Happy 2018!
EOF

e file.txt | ca /2018/ | c 2019 | ca , | p
# Prints: Happy 2019!
```

### `i <text>`

Like `c`, but inserts `text` right before the dot.

### `a <text>`

Like `c`, but inserts `text` right after the dot.

### `d`

Deletes the dot's content.

Example:

```bash
<<EOF > file.txt
We live in a different society.
EOF

e file.txt | ca '/ different/' | d | ca , | p
# Prints: We live in a society.
```

### `s <regexp> <text> [n|g]`

Substitute `text` for the first match to the `regexp` in the dot. Set dot to the modified range. 
In `text`, `$` signs are interpreted as in Go's [`Expand`](https://golang.org/pkg/regexp/#Regexp.Expand).

Example:

```bash
<<EOF > file.txt
y y a x y y
EOF

e file.txt | ca 1 | s '(a) (x)' '$2 $1' | ca , | p
# Prints: y y x a y y
```