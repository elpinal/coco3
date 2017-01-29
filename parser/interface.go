package parser

import "github.com/elpinal/coco3/ast"

func ParseSrc(src []byte) (f *ast.File, err error) {
	var p parser
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
	}()

	// parse source
	p.init(src)
	f = p.parseFile()

	return
}
