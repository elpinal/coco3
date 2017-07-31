package editor

import (
	"fmt"
	"io"
)

func help(w io.Writer) error {
	for k, v := range quickref {
		_, err := w.Write([]byte(fmt.Sprintf("%s   %s\n", k, v)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Reference: Vim's quickref.txt.
var quickref = map[string]string{
	"h": "left",
	"l": "right",
	"0": "to first character in the line",
	"^": "to first non-blank character in the line",
	"$": "to the last character in the line",
	"|": "to column N",
	"f": "to the Nth occurrence of {char} to the right",
	"F": "to the Nth occurrence of {char} to the left",

	"k": "go back history",
	"j": "go forward history",

	"-": "decrement the number at or after the cursor",
	"+": "increment the number at or after the cursor",

	"w": "N words forward",
	"W": "N blank-separated WORDs forward",
	"e": "forward to the end of the Nth word",
	"E": "forward to the end of the Nth blank-separated WORD",
	"b": "N words backward",
	"B": "N blank-separated WORDs backward",

	"/": "search forward",
	"?": "search backward",

	"n": "repeat last search",
	"N": "repeat last search, in opposite direction",

	"a":  "append text after the cursor",
	"A":  "append text at the end of the line",
	"i":  "insert text before the cursor",
	"I":  "insert text before the first non-blank in the line",
	"gI": "insert text in column 1",

	// insert mode...

	"i_<Esc>":  "end Insert mode, back to Normal mode",
	"i_CTRL-C": "like <Esc>",

	"i_CTRL-R": "insert the contents of a register",
	"i_CTRL-Y": "complete the word before the cursor in various ways",
	"i_<BS":    "delete the character before the cursor",
	"i_CTRL-W": "delete word before the cursor",
	"i_CTRL-U": "delete all entered characters in the current line",
}
