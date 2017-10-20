package screen

import "github.com/elpinal/coco3/config"

type Screen interface {
	Start(*config.Config, bool, []rune, int, *Hi)
	Refresh(*config.Config, bool, []rune, int, *Hi)
	SetLastLine(string)
}

// Hi represents a range for highlight.
type Hi struct {
	Left, Right int
}

// TestScreen is one of the simplest implementation for Screen. It does not
// anything. It is useful for testing.
type TestScreen struct{}

func (_ *TestScreen) Start(_ *config.Config, _ bool, _ []rune, _ int, _ *Hi) {
}

func (_ *TestScreen) Refresh(_ *config.Config, _ bool, _ []rune, _ int, _ *Hi) {
}

func (_ *TestScreen) SetLastLine(_ string) {
}
