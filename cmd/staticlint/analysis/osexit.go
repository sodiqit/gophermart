package analysis

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if c, ok := node.(*ast.CallExpr); ok {
				if s, ok := c.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						if _, ok := c.Args[0].(*ast.BasicLit); ok {
							pass.Reportf(s.Pos(), "declaration os.Exit shouldn't be used")
						}
					}
				}
			}
			return true
		})
	}

	return nil, nil
}
