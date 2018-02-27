package editor

const (
	CharCtrlAt = iota
	CharCtrlA
	CharCtrlB
	CharCtrlC
	CharCtrlD
	CharCtrlE
	CharCtrlF
	CharCtrlG
	CharCtrlH
	CharCtrlI
	CharCtrlJ
	CharCtrlK
	CharCtrlL
	CharCtrlM
	CharCtrlN
	CharCtrlO
	CharCtrlP
	CharCtrlQ
	CharCtrlR
	CharCtrlS
	CharCtrlT
	CharCtrlU
	CharCtrlV
	CharCtrlW
	CharCtrlX
	CharCtrlY
	CharCtrlZ
	CharEscape

	// In ASCII, the first 32 codes are used for control characters, which are not printable.
	// EndOfControlCharacters represents the last of these characters, "^_".
	EndOfControlCharacters = 31

	CharBackspace = 127
)
