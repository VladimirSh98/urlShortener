package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/gostaticanalysis/forcetypeassert"
	"github.com/orijtech/httperroryzer"

	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		errorsas.Analyzer,
		findcall.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	)

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer != nil && len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer == nil {
			continue
		}
		switch v.Analyzer.Name {
		case "S1000", "ST1000":
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	analyzers = append(analyzers, forcetypeassert.Analyzer)
	analyzers = append(analyzers, httperroryzer.Analyzer)

	analyzers = append(analyzers, noExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}

var noExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "disallow direct calls to os.Exit in main.main",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			if pass.Pkg.Name() != "main" {
				continue
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name.Name != "main" {
					continue
				}
				ast.Inspect(fn.Body, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					sel, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}
					xIdent, ok := sel.X.(*ast.Ident)
					if !ok {
						return true
					}
					if xIdent.Name == "os" && sel.Sel.Name == "Exit" {
						pass.Reportf(call.Pos(), "direct call to os.Exit is forbidden in main")
					}
					return true
				})
			}
		}
		return nil, nil
	},
}
