package parser

import (
	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/token"
)

func ParseSrc(src []byte) (f *ast.File, err error) {
	return ParseFile(token.NewFileSet(), "command line", src)
}

func ParseFile(fset *token.FileSet, filename string, src []byte) (f *ast.File, err error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}
	var p parser
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
		p.errors.Sort()
		err = p.errors.Err()
	}()

	// parse source
	p.init(fset, filename, src)
	f = p.parseFile()

	return
}
