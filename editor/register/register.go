package register

const (
	// Numbered registers "0 to "9 are represented as 'x'.
	// 26 named registers "a to "z or "A to "Z are represented as 'a'.

	Unnamed        = '"'
	SmallDelete    = '-'
	RecentExecuted = ':'
	LastInserted   = '.'
	CurrentFile    = '%'
	Expression     = '='
	Clipboard      = '*'
	BlackHole      = '_'
	LastSearch     = '/'
	// Alternate     = '#'
	// ClipboardPlus = '+'
	// Dropped       = '~'
)

func IsValid(r rune) bool {
	switch r {
	case Unnamed, SmallDelete, RecentExecuted, LastInserted,
		CurrentFile, Expression, Clipboard, BlackHole, LastSearch:
		return true
	}
	if '0' <= r && r <= '9' || 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' {
		return true
	}
	return false
}

type Registers struct {
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

func (r *Registers) Init() {
	if r.named == nil {
		r.named = make(map[rune][]rune)
	}
}

func (r *Registers) Register(where rune, s []rune) {
	switch where {
	case Unnamed:
		r.unnamed = s
		return
	case SmallDelete:
		r.smallDelete = s
	case RecentExecuted:
		r.recentExecuted = s
	case LastInserted:
		r.lastInserted = s
	case CurrentFile:
		r.currentFile = s
	case Expression:
		r.expression = s
	case Clipboard:
		r.clipboard = s
	case BlackHole:
		// no-op
		return
	case LastSearch:
		r.lastSearch = s
	}
	switch {
	case '0' <= where && where <= '9':
		r.numbered[where-'0'] = s
	case 'a' <= where && where <= 'z':
		r.named[where] = s
	case 'A' <= where && where <= 'Z':
		i := where - 'A' + 'a'
		r.named[i] = append(r.named[i], s...)
	}

	// The unnamed register is pointing to the last used register.
	r.unnamed = s
}

func (r *Registers) Read(where rune) []rune {
	switch {
	case '0' <= where && where <= '9':
		return r.numbered[where-'0']
	case 'a' <= where && where <= 'z':
		return r.named[where]
	case 'A' <= where && where <= 'Z':
		return r.named[where-'A'+'a']
	}
	switch where {
	case Unnamed:
		return r.unnamed
	case SmallDelete:
		return r.smallDelete
	case RecentExecuted:
		return r.recentExecuted
	case LastInserted:
		return r.lastInserted
	case CurrentFile:
		return r.currentFile
	case Expression:
		return r.expression
	case Clipboard:
		return r.clipboard
	case BlackHole:
		// no-op
	case LastSearch:
		return r.lastSearch
	}
	return nil
}
