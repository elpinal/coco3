package editor

import "time"

type node struct {
	data   []rune
	parent *node
	child  *node
	when   time.Time
}

type undoTree struct {
	current *node
	nodes   []*node // in chronological order
}

func newUndoTree() undoTree {
	return undoTree{current: &node{}}
}

func (u *undoTree) undo() ([]rune, bool) {
	if u.current.parent == nil {
		return nil, false
	}
	u.current = u.current.parent
	return u.current.data, true
}

func (u *undoTree) redo() ([]rune, bool) {
	if u.current.child == nil {
		return nil, false
	}
	u.current = u.current.child
	return u.current.data, true
}

func (u *undoTree) earlier() []rune {
	for i, n := range u.nodes {
		if n == u.current {
			if i == 0 {
				return nil
			}
			u.current = u.nodes[i-1]
			return u.current.data
		}
	}
	panic("unexpected loss of undo node")
}

func (u *undoTree) later() []rune {
	for i, n := range u.nodes {
		if n == u.current {
			if i == len(u.nodes)-1 {
				return nil
			}
			u.current = u.nodes[i+1]
			return u.current.data
		}
	}
	panic("unexpected loss of undo node")
}

func (u *undoTree) add(s []rune) {
	data := make([]rune, len(s))
	copy(data, s)
	n := &node{data: data, parent: u.current, when: time.Now()}
	u.current = n
	u.current.parent.child = n
	u.nodes = append(u.nodes, n)
}
