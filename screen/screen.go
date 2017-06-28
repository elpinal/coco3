package screen

import "github.com/elpinal/coco3/config"

type Screen interface {
	Start(*config.Config, []rune, int)
	Refresh(*config.Config, []rune, int)
	SetLastLine(string)
}
