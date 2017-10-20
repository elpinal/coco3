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

type TestScreen struct {
}

func (_ *TestScreen) Start(_ *config.Config, _ bool, _ []rune, _ int, _ *Hi) {
}

func (_ *TestScreen) Refresh(_ *config.Config, _ bool, _ []rune, _ int, _ *Hi) {
}

func (_ *TestScreen) SetLastLine(_ string) {
}
