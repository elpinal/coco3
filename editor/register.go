package main

const (
	// Numbered registers "0 to "9 are represented as 'x'.
	// 26 named registers "a to "z or "A to "Z are represented as 'a'.

	registerUnnamed        = '"'
	registerSmallDelete    = '-'
	registerRecentExecuted = ':'
	registerLastInserted   = '.'
	registerCurrentFile    = '%'
	registerExpression     = '='
	registerClipboard      = '*'
	registerBlackHole      = '_'
	registerLastSearch     = '/'
	// registerAlternate     = '#'
	// registerClipboardPlus = '+'
	// registerDropped       = '~'
)

type registers struct {
	numbered [10][]rune
	named    map[rune][]rune

	unnamed,
	smallDelete,
	recentExecuted,
	lastInserted,
	currentFile,
	expression,
	clipboard,
	lastSearch []rune
}

func (r *registers) init() {
	if r.named == nil {
		r.named = make(map[rune][]rune)
	}
}

func (r *registers) register(where rune, s []rune) {
	switch {
	case '0' <= where && where <= '9':
		r.numbered[where-'0'] = s
		return
	case 'a' <= where && where <= 'z':
		r.named[where] = s
		return
	case 'A' <= where && where <= 'Z':
		i := where - 'A' + 'a'
		r.named[i] = append(r.named[i], s...)
		return
	}
	switch where {
	case registerUnnamed:
		r.unnamed = s
	case registerSmallDelete:
		r.smallDelete = s
	case registerRecentExecuted:
		r.recentExecuted = s
	case registerLastInserted:
		r.lastInserted = s
	case registerCurrentFile:
		r.currentFile = s
	case registerExpression:
		r.expression = s
	case registerClipboard:
		r.clipboard = s
	case registerBlackHole:
		// no-op
	case registerLastSearch:
		r.lastSearch = s
	}
}

func (r *registers) read(where rune) []rune {
	switch {
	case '0' <= where && where <= '9':
		return r.numbered[where-'0']
	case 'a' <= where && where <= 'z':
		return r.named[where]
	case 'A' <= where && where <= 'Z':
		return r.named[where-'A'+'a']
	}
	switch where {
	case registerUnnamed:
		return r.unnamed
	case registerSmallDelete:
		return r.smallDelete
	case registerRecentExecuted:
		return r.recentExecuted
	case registerLastInserted:
		return r.lastInserted
	case registerCurrentFile:
		return r.currentFile
	case registerExpression:
		return r.expression
	case registerClipboard:
		return r.clipboard
	case registerBlackHole:
		// no-op
	case registerLastSearch:
		return r.lastSearch
	}
	return nil
}
