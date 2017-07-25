package screen

import "github.com/elpinal/coco3/config"

type Screen interface {
	Start(*config.Config, bool, []rune, int)
	Refresh(*config.Config, bool, []rune, int)
	SetLastLine(string)
}
