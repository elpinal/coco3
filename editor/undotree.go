package main

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

func (u *undoTree) undo() []rune {
	if u.current.parent == nil {
		return nil
	}
	u.current = u.current.parent
	return u.current.data
}

func (u *undoTree) redo() []rune {
	if u.current == nil {
		return nil
	}
	u.current = u.current.child
	return u.current.data
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
	n := &node{data: s, parent: u.current, when: time.Now()}
	u.current.child = n
	u.current = n
	u.nodes = append(u.nodes, n)
}
