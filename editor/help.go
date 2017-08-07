package editor

import (
	"fmt"
	"io"
)

func Help(w io.Writer) error {
	for _, x := range quickref {
		_, err := w.Write([]byte(fmt.Sprintf("%-12s%s\n", x.k, x.v)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Reference: Vim's quickref.txt.
var quickref = []struct{ k, v string }{
	{"h", "left"},
	{"l", "right"},
	{"0", "to first character in the line"},
	{"^", "to first non-blank character in the line"},
	{"$", "to the last character in the line"},
	{"|", "to column N"},
	{"f", "to the Nth occurrence of {char} to the right"},
	{"F", "to the Nth occurrence of {char} to the left"},
	{"t", "till before the Nth occurrence of {char} to the right"},
	{"T", "till bl before the Nth occurrence of {char} to the left"},

	{"k", "go back history"},
	{"j", "go forward history"},

	{"-", "decrement the number at or after the cursor"},
	{"+", "increment the number at or after the cursor"},

	{"w", "N words forward"},
	{"W", "N blank-separated WORDs forward"},
	{"e", "forward to the end of the Nth word"},
	{"E", "forward to the end of the Nth blank-separated WORD"},
	{"b", "N words backward"},
	{"B", "N blank-separated WORDs backward"},
	{"ge", "backward to the end of the Nth word"},
	{"gE", "backward to the end of the Nth blank-separated WORD"},

	{"[(", "N times back to unclosed '('"},
	{"[{", "N times back to unclosed '{'"},
	{"])", "N times forward to unclosed ')'"},
	{"]}", "N times forward to unclosed '}'"},

	{"/", "search forward"},
	{"?", "search backward"},

	{"n", "repeat last search"},
	{"N", "repeat last search, in opposite direction"},

	{"a", "append text after the cursor"},
	{"A", "append text at the end of the line"},
	{"i", "insert text before the cursor"},
	{"I", "insert text before the first non-blank in the line"},
	{"gI", "insert text in column 1"},

	// insert mode...

	{"i_<Esc>", "end Insert mode, back to Normal mode"},
	{"i_CTRL-C", "like <Esc>"},

	{"i_CTRL-R", "insert the contents of a register"},
	{"i_CTRL-X", "complete the word before the cursor in various ways"},
	{"i_<BS>", "delete the character before the cursor"},
	{"i_CTRL-W", "delete word before the cursor"},
	{"i_CTRL-U", "delete all entered characters in the current line"},

	{"x", "delete N characters under and after the cursor"},
	{"<Del>", "delete N characters under and after the cursor"},
	{"X", "delete N characters before the cursor"},
	{"d", "delete the text that is moved over with {motion}"},
	{"v_d", "delete the highlighted text"},
	{"dd", "delete N lines"},
	{"D", "delete to the end of the line"},

	{"\"", "use register {char} for the next delete, yank, or put"},
	{"y", "yank the text moved over with {motion} into a register"},
	{"v_y", "yank the highlighted text into a register"},
	{"yy", "yank N lines into a register"},
	{"Y", "yank to the end of line into a register"},
	{"p", "put a register after the cursor position (N times)"},
	{"P", "put a register before the cursor position (N times)"},

	{"r", "replace N characters with {char}"},
	{"R", "enter Replace mode"},
	{"c", "change the text that is moved over with {motion}"},
	{"v_c", "change the highlighted text"},
	{"cc", "change N lines"},
	{"C", "change to the end of the line"},
	{"~", "switch case for N characters and advance cursor"},
	{"v_~", "switch case for highlighted text"},
	{"v_u", "make highlighted text lowercase"},
	{"v_U", "make highlighted text uppercase"},
	{"g~", "switch case for the text that is moved over with {motion}"},
	{"gu", "make the text that is moved over with {motion} lowercase"},
	{"gU", "make the text that is moved over with {motion} uppercase"},

	{"v", "start highlighting characters  }  move cursor and use"},
	{"V", "start highlighting linewise    }  operator to affect"},
	{"o", "exchange cursor position with start of highlighting"},
	{"v_v", "highlight characters or stop highlighting"},

	{"aw", `Select "a word"`},
	{"iw", `Select "inner word"`},
	{"aW", `Select "a |WORD|"`},
	{"iW", `Select "inner |WORD|"`},
	{"as", `Select "a sentence"`},
	{"is", `Select "inner sentence"`},
	{"ap", `Select "a paragraph"`},
	{"ip", `Select "inner paragraph"`},
	{"ab", `Select "a block" (from "[(" to "])")`},
	{"ib", `Select "inner block" (from "[(" to "])")`},
	{"aB", `Select "a Block" (from "[{" to "]}")`},
	{"iB", `Select "inner Block" (from "[{" to "]}")`},
	{"a>", `Select "a <> block"`},
	{"i>", `Select "inner <> block"`},
	{"at", `Select "a tag block" (from <aaa> to </aaa>)`},
	{"it", `Select "inner tag block" (from <aaa> to </aaa>)`},
	{"a'", `Select "a single quoted string"`},
	{"i'", `Select "inner single quoted string"`},
	{"a\"", `Select "a double quoted string"`},
	{"i\"", `Select "inner double quoted string"`},
	{"a`", `Select "a backward quoted string"`},
	{"i`", `Select "inner backward quoted string"`},
}
