package importalias

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"path"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

func Run(pkgPaths []string, verbose bool, w io.Writer) error {
	for _, pkgPath := range pkgPaths {
		loadedPkg, _ := build.ImportDir(pkgPath, 0)

		var goFilesForPkg []string
		goFilesForPkg = append(goFilesForPkg, loadedPkg.GoFiles...)
		goFilesForPkg = append(goFilesForPkg, loadedPkg.TestGoFiles...)
		goFilesForPkg = append(goFilesForPkg, loadedPkg.XTestGoFiles...)
		sort.Strings(goFilesForPkg)

		hasError := false
		for _, currGoFileName := range goFilesForPkg {
			currFile := path.Join(pkgPath, currGoFileName)
			failedChecks, err := checkFile(currFile)
			if err != nil {
				return errors.Wrapf(err, "file %v failed k8s importalias check", currFile)
			}
			for _, failedCheck := range failedChecks {
				fmt.Fprintln(w, failedCheck.message)
			}
			if len(failedChecks) > 0 {
				hasError = true
			}
		}
		if hasError {
			return fmt.Errorf("")
		}
	}
	return nil
}

func checkFile(filename string) ([]failedCheck, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %v", filename, err)
	}

	var visitor visitFn
	var failedChecks []failedCheck
	visitor = visitFn(func(node ast.Node) ast.Visitor {
		if node == nil {
			return visitor
		}
		switch v := node.(type) {
		case *ast.ImportSpec:
			failedCheck := checkImportAlias(filename, v.Name, v.Path.Value, fset.Position(v.Pos()))
			if failedCheck != nil {
				failedChecks = append(failedChecks, *failedCheck)
			}
			break
		}
		return visitor
	})
	ast.Walk(visitor, file)
	return failedChecks, nil
}

func checkImportAlias(filename string, alias *ast.Ident, path string, pos token.Position) *failedCheck {
	for _, aliasDeriver := range aliasDerivers {
		expected := aliasDeriver.Alias(path)
		if expected != "" {
			if alias == nil {
				return &failedCheck{fmt.Sprintf("%s:%d:%d: %s must declare import alias %q", filename, pos.Line, pos.Column, path, expected)}
			}
			if expected != alias.Name {
				return &failedCheck{fmt.Sprintf("%s:%d:%d: expected %s to declare import alias %q but was %q", filename, pos.Line, pos.Column, path, expected, alias.Name)}
			}
			return nil
		}
	}
	return nil
}

type failedCheck struct {
	message string
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
	match := k8sioApimachineryRegex.MatchString(filename)
	if match {
		return "meta"
	}
	return ""
}

var aliasDerivers = []aliasDeriver{
	aliasDeriverFn(k8sioApiAlias),
	aliasDeriverFn(k8sioApimachineryAlias),
}
