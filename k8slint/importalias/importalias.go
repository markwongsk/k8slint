package importalias

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

func Run() error {
	for _, pkgPath := range pkgPaths {
		loadedPkg, _ := build.ImportDir(pkgPath, 0)

		var goFilesForPkg []string
		goFilesForPkg = append(goFilesForPkg, loadedPkg.GoFiles...)
		goFilesForPkg = append(goFilesForPkg, loadedPkg.TestGoFiles...)
		goFilesForPkg = append(goFilesForPkg, loadedPkg.XTestGoFiles...)
		sort.Strings(goFilesForPkg)

		for _, currGoFileName := range goFilesForPkg {
			currFile := path.Join(pkgPath, currGoFileName)
			if err := checkFile(currFile); err != nil {
				return errors.Wrapf(err, "file %v failed k8s importalias check", currFile)
			}
		}
	}
}

func checkFile(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %v", filename, err)
	}

	hasError = false
	visitor = visitFn(func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.ImportSpec:
			err := checkImportAlias(filename, v.Name, v.Path, fset.Position(v.Pos()))
			if err != nil {
				hasError = true
			}
			break
		}
		return visitor
	})
	ast.Walk(visitor, file)
	if hasError {
		return fmt.Errorf("")
	}
	return nil
}

func checkImportAlias(filename string, fset *token.FileSet, alias *ast.Ident, path *ast.BasicLit, pos token.Position) error {
	for _, aliasDeriver := range aliasDerivers {
		expected := aliasDeriver.Alias(path)
		if expected != nil {
			if alias == nil {
				fmt.Fprintf("%s:%d:%d: %q must declare import alias %q", filename, pos.Line, pos.Column, path, expected)
				return fmt.Errorf("")
			}
			if expected != alias.Name.Name {
				fmt.Fprintf("%s:%d:%d: expected %q to declare import alias %q but was %q", filename, pos.Line, pos.Column, path, expected, alias.Name)
				return fmt.Errorf("")
			}
			return nil
		}
	}
}

type visitFn func(node ast.Node) ast.Visitor

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}

type aliasDeriver interface {
	Alias(filename string) string
}

type aliasDeriverFn func(filename string) string

func (fn aliasDeriverFn) Alias(filename string) string {
	return fn(filename)
}

var k8sioApiRegex = regexp.MustCompile(`k8s\.io\/api\/([^\/]*)\/v.*`)

func k8sioApiAlias(filename string) string {
	sm := k8sioApiRegex.FindStringSubmatch(filename)
	if sm == nil {
		return ""
	}

	return sm[1]
}

var k8sioApimachineryRegex = regexp.MustCompile(`k8s\.io\/apimachinery\/pkg\/apis\/meta\/v.*`)

func k8sioApimachineryAlias(filename string) string {
	bool := k8sioApimachineryRegex.MatchString(filename)
	if match {
		return "meta"
	}
	return ""
}

var aliasDerivers = []aliasDeriver{
	aliasDeriverFn(k8sioApiAlias),
	aliasDeriverFn(k8sioApimachineryAlias),
}
