package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod/lib/utils"
)

func main() {
	log.Println("npx eslint --ext .js,.html .")
	utils.Exec("npx", "eslint", "--ext", ".js,.html", ".")

	log.Println("npx prettier --loglevel error --write .")
	utils.Exec("npx", "prettier", "--loglevel", "error", "--write", ".")

	log.Println("godev lint")
	utils.Exec("godev", "lint")

	lintMustPrefix()
}

func lintMustPrefix() {
	log.Println("lint 'Must' prefix")

	paths, err := filepath.Glob("*.go")
	utils.E(err)
	lintErr := false

	for _, p := range paths {
		name := filepath.Base(p)
		if name == "must.go" || strings.HasSuffix(name, "_test.go") {
			continue
		}

		src, err := utils.ReadString(p)
		utils.E(err)

		list := token.NewFileSet()
		f, err := parser.ParseFile(list, p, src, 0)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if ok && strings.HasPrefix(fd.Name.Name, "Must") {
				log.Printf("%s %s\n", list.Position(fd.Name.Pos()), fd.Name.Name)
				lintErr = true
			}
		}
		break
	}

	if lintErr {
		log.Fatalln("'Must' prefixed function should be declared in file 'must.go'")
	}
}
