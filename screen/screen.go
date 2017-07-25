package screen

import "github.com/elpinal/coco3/config"

type Screen interface {
	Start(*config.Config, bool, []rune, int, *Hi)
	Refresh(*config.Config, bool, []rune, int, *Hi)
	SetLastLine(string)
}

type Hi struct {
	Left, Right int
}
