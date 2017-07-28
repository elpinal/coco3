package editor

import (
	"fmt"
	"io"
)

func help(w io.Writer) error {
	for k, v := range helpMap {
		_, err := w.Write([]byte(fmt.Sprintf("%s   %s\n", k, v)))
		if err != nil {
			return err
		}
	}
	return nil
}

var helpMap = map[string]string{
	"h": "left",
	"l": "right",
	"0": "to first character in the line",
	"$": "to the last character in the line",
}
