package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

func init() {
	register(fix{
		name: "log.Fatal",
		desc: `Replaces error checks that call log.Fatal with exit.on`,
		date: "2017-09-10",
		f:    fixLogFatal,
	})
}

func fixLogFatal(file *ast.File) bool {
	fixed := false

	// Only scan files that import "log"
	if !astutil.UsesImport(file, "log") {
		return false
	}

	fixFunc := func(parent ast.Node, name string, index int, n ast.Node) bool {

		// only look at if-statements
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		// ignore if-statements that have an else block
		if ifStmt.Else != nil {
			return true
		}

		// ensure that the condition is "{var type} != nil"
		ifCondExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		if ifCondExpr.Op != token.NEQ {
			return true
		}
		ifCondExprYIdent, ok := ifCondExpr.Y.(*ast.Ident)
		if !ok {
			return true
		}
		if ifCondExprYIdent.Name != "nil" {
			return true
		}
		ifCondExprXIdent, ok := ifCondExpr.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ifCondExprXIdent.Obj.Kind != ast.Var {
			return true
		}

		// ensure that the block of the if-statement only contains a call to log.Fatal
		if len(ifStmt.Body.List) != 1 {
			return true
		}
		ifBodyExprStmt, ok := ifStmt.Body.List[0].(*ast.ExprStmt)
		if !ok {
			return true
		}
		ifBodyStmtCallExpr, ok := ifBodyExprStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}
		ifBodyStmtCallSelectorExpr, ok := ifBodyStmtCallExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ifBodyStmtCallSelectorIdent, ok := ifBodyStmtCallSelectorExpr.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ifBodyStmtCallSelectorIdent.Name != "log" {
			return true
		}
		if ifBodyStmtCallSelectorExpr.Sel.Name != "Fatal" {
			return true
		}
		if len(ifBodyStmtCallExpr.Args) != 1 {
			return true
		}

		// ensure that log.Fatal is called with the variable used in the condition of the if-statement
		ifBodyStmtCallExprArgIdent, ok := ifBodyStmtCallExpr.Args[0].(*ast.Ident)
		if !ok {
			return true
		}
		if ifCondExprXIdent.Name != ifBodyStmtCallExprArgIdent.Name {
			return true
		}

		// replace the if-expression with a call to "exit.On({var})"
		newExprStmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "exit",
					},
					Sel: &ast.Ident{
						Name: "On",
					},
				},
				Args: []ast.Expr{
					ifBodyStmtCallExprArgIdent,
				},
			},
		}
		SetField(parent, name, index, newExprStmt)

		// Add the import
		astutil.AddImport(fset, file, "github.com/Originate/exit")

		fixed = true
		return false
	}

	Apply(file, fixFunc, nil)
	return fixed
}
