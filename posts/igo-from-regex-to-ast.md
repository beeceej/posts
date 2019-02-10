<!--id: 7-->
<!--title: Writing a Go Interpreter in Go, pt2 -->
<!--author: Brian Jones-->
<!--visible: true-->

This is part 2 in a mini series about writing an interpreter in Go. The first post can be found here [here](https://blog.beeceej.com/blog/4).

As a quick recap, in part 1, I began parsing out snippets of go code into a data structure with regexes and evaluating them on the fly. I quickly realized regex weren't going to be a scalable solution (though it was fun to build them), the code wasn't as robust as it needed to be, and also i was unable to parse out what I needed at the lowest granularity. I knew the [go ast library](https://godoc.org/go/ast) existed, but I hadn't used it yet. This project is the perfect excuse to use it. And I'll tell ya, after a little while working with it, the library is wonderful. I'm not going to go into too much detail, if you want to see what the code looked like before you can reference the previous post, but here's what the regex looked like. nice and juicy.

```go
 `(?m)(func (.+?)(\(.+?\,|\(.+?)*\)((\(.+?\,|\(.+?)+\)|.*)\{([\s\Sa-zA-Z1-9]*?(^}|(\}$\}\s\S))))`
```


This is what the parsing looks like now:

```go
package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
)

var hasPackageStatementRegexp = regexp.MustCompile("^package.*")

// ASTParse is
type ASTParse struct {
	Raw       string
	Functions []*Function

	fset *token.FileSet
	root ast.Node
}

func (a *ASTParse) Parse() {
	a.Setup()
	ast.Inspect(a.root, func(n ast.Node) bool {
		err := ifFunctionDeclaration(
			compose(
				a.print,
				a.ParseFn,
			), n)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		return true
	})
}

func (a *ASTParse) Setup() {
	if !hasPackageStatementRegexp.MatchString(a.Raw) {
		const withPkg = `package _
		%s`
		a.Raw = fmt.Sprintf(withPkg, a.Raw)
	}
	a.Functions = []*Function{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "/tmp/tmp.go", a.Raw, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	a.fset = fset
	a.root = node
}

func (a *ASTParse) ParseFn(n *ast.FuncDecl) error {
	var (
		identifier      string
		params          string
		returnSignature string
		body            string
		buffer          bytes.Buffer
	)
	identifier = n.Name.Name
	params = a.getFunctionParameters(n)

	returnSignature = a.getFunctionReturnSignature(n)

	printer.Fprint(&buffer, a.fset, n.Body)
	body = buffer.String()

	a.Functions = append(a.Functions, &Function{
		Identifier: identifier,
		Params:     params,
		// This is actually an ast.Node in-itself, we could parse sub functions recursively,
		// not sure that's completely necessary just yet though
		Body:   body,
		Return: returnSignature,
	})
	return nil
}

func (a *ASTParse) getFunctionParameters(n *ast.FuncDecl) string {
	var (
		params []string
		buffer bytes.Buffer
	)
	if n.Type.Results == nil {
		return ""
	}

	offset := n.Pos()
	printer.Fprint(&buffer, a.fset, n)
	rawFn := buffer.String()
	for _, pGroup := range n.Type.Params.List {
		p := pGroup.Type.Pos() - offset
		e := pGroup.Type.End() - offset
		groupType := rawFn[p:e]
		params = append(params, paramGroup(pGroup.Names, groupType))
	}

	return strings.Join(params, ", ")
}

func (a *ASTParse) getFunctionReturnSignature(n *ast.FuncDecl) string {
	var (
		returns []string
		buffer  bytes.Buffer
	)
	if n.Type.Results == nil {
		return ""
	}
	offset := n.Pos()
	printer.Fprint(&buffer, a.fset, n)
	fnRaw := buffer.String()

	for _, rGroup := range n.Type.Results.List {
		p := rGroup.Type.Pos() - offset
		e := rGroup.Type.End() - offset
		returns = append(returns, fnRaw[p:e])
	}

	return strings.Join(returns, ", ")

}

func paramGroup(idents []*ast.Ident, t string) string {
	if len(idents) < 1 {
		return ""
	}
	paramNames := []string{}
	for _, i := range idents {
		if i.Name == "" {
			continue
		}
		paramNames = append(paramNames, i.Name)
	}

	return fmt.Sprint(strings.Join(paramNames, ", "), " ", t)
}

func (a *ASTParse) print(n *ast.FuncDecl) error {
	printer.Fprint(os.Stdout, a.fset, n)
	return nil
}

```

Using go's ast package gives me all the power I need to traverse through the code. It turns out at this point I don't need the granularity that it provides, so in some cases I'm throwing away data that I just don't need. Think comment data, or recursing into the body of a function.

Some key concepts when working with the ast package are:

* [func Inspect(node Node, f func(Node) bool](https://godoc.org/go/ast#Inspect)
  * Akin to `filepath.Walk` allows you to recurse through the tree of a source file
* to pretty Print a node to std out:
```go
fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, "/tmp/tmp.go", yourProgramAsString, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
  printer.Fprint(os.Stdout, fset, root)
```

* to pretty Print a node to a string:
```go
fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, "/tmp/tmp.go", yourProgramAsAString, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
  }
  var b bytes.Buffer
  printer.Fprint(&buffer, fset, root)
  programAsString := b.String()
```

* The ast.Node interface, This one is awesome!
  * The AST is made up of many different types of nodes in the tree, but there are generic operations (like in the printing above) that can be done on an ast.Node. `root` is an ast.Node

* While `ast.Inspect'ing` you can type assert the ast.Node and access properties specific to the current node; or example, if you're looking for Function Declarations you can do the following:
```go
	if fnDecl, ok := n.(*ast.FuncDecl); ok {
		// Do Some Stuff wit the *ast.FuncDecl struct
	}
```

* I decided to abstract away the explicit type assertions into functions, like:
```go 
func ifFunctionDeclaration(h funcDeclarationHandler, n ast.Node) error {
	if fnDecl, ok := n.(*ast.FuncDecl); ok {
		return h(fnDecl)
	}
	return nil
}

func usage() {
  err := ifFunctionDeclaration(func(n *ast.FuncDecl) error {
    // do stuff
  }, node)
}
```

You can even compose these functions...

```go
// compose declaration handlers into one function, generics who?
func compose(fns ...funcDeclarationHandler) funcDeclarationHandler {
	return func(n *ast.FuncDecl) error {
		for _, fn := range fns {
			if err := fn(n); err != nil {
				return err
			}
		}
		return nil
	}
}

```

In a nutshell, don't be afraid of AST's. They are awesome, code introspection is a lot of fun. And remember, IT'S JUST DATA.
As always, [you can find iGo here](https://github.com/beeceej/iGo)