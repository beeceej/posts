<!--id: 4-->
<!--title: Writing a Go interpreter in Go -->
<!--author: Brian Jones-->
<!--visible: true-->

## Why Interpreters

The main reason to use an interpreter is an improved workflow. If you're like me then you're always searching for ways to improve working efficiency. Interpreters are one surefire way to improve productivity, and it's all about shortening feedback loops. Shorter feedback loops usually equate to quicker turn around times. Even with Go, which compiles extremely fast, seeing results as you type like in Python or Lisp can be extremely beneficial. Python and lisp are both interpreted and have special tooling built out for this purpose. The draw of LSIP is the interpreter and the rapid feedback loop that evaluating code on the fly gives you. Hell, even Java has a REPL now, and if Java has something like an interpreter there's no reason Go can't have one. At the core of Python and LISP's respective design, is the interpreter which makes an interpreted workflow first class within each languages respective ecosystem. Unfortunately, in general, the compilation process of languages like Java (JVM Bytecode) or Go (Binary) doesn't lend itself naturally to evaluating expressions on the fly. Historically, to get a quick feedback while developing in a compiled language you'd need to execute some combination of the following steps

1.  Write Tests
2.  Run Tests, and observe behavior
3.  Make changes to code
4.  Repeat

This is a good workflow, but we can do better. This is why interpreters are awesome. Interpreters _abstract_ away the boring parts about testing our code and give us rapid feedback loops as we work. In EMACs and LISP land, that's just a couple key strokes away, my goal is to replicate this in Go.

So... with software anything is possible? Just because Go is compiled doesn't mean we can't conjure up an interpreter for a compiled language, developers _are wizards_ aren't we?

## Humble Beginnings

So with the project in mind, I started writing some code.

I have a very loose understanding of what an interpreter is. So as a result my implementation of an interpreter is very loose is. At it's most simple form, I broke down each piece into these components:

- **Interpreter will Accept Some Raw Text:**
  - Classify the text into functions, expressions, etc.
  - Store the classified text in memory for later use
  - Store history of user defined input
- **Generate code and compile it**
  - Use previous commands to generate code
  - Compile & Execute the code
  - Given some result (Error or success), report back to the user
- **The interpreter should accept multiple forms of Interaction**
  - Client-Server
  - REPL/CLI (enter some text in, press enter, view results)
  - Communication is decoupled from the core implementation

## An initial REPL

Turns out, building a usable CLI interface to an interpreter is not easy. my initial strategy looked something like this

```go
const mainTmpl = `package main
func main() {
    {{.}}
    {{ .M }}
}
 {{ .F }}`

func main() {
    t := template.Must(template.New("mainTmpl").Parse(mainTmpl))
    t.Execute(os.Stdout, `fmt.Println("iGoIsCool")`)
    history := make(map[string]string)
    var instruct In
    for {
        fmt.Print("\n$ ")
        f, err := os.Create("exe.go")
        if err != nil {
            panic(err.Error())
        }
        defer f.Close()
        reader := bufio.NewReader(os.Stdin)
        text, _ := reader.ReadString('\n')
        if history["a"] != "" {
            instruct.F = history["a"]
        }
        instruct.M = text
        if err := t.Execute(f, instruct); err != nil {
            fmt.Println(err.Error())
        }
        if f, err := iGo.NewFunction(text); err != nil {
            fmt.Println(err.Error())
        } else {
            history[f.Identifier] = f.Raw
            for k, v := range history {
                fmt.Printf("%s: %s\n", k, v)
            }
        }
        cmd := exec.Command("goimports", "-w", "exe.go")
        b, err := cmd.Output()
        if err != nil {
            fmt.Println("Error calling goimports", err.Error())
            continue
        }
        fmt.Println(string(b))

        cmd = exec.Command("go", "run", "exe.go")
        b, err = cmd.Output()
        if err != nil {
            fmt.Println("Error calling go run", err.Error())
            continue
        }
        fmt.Println(fmt.Sprintf(">> %s", string(b)))

        os.Remove(path)
    }
}
```

This is what I call some ugly Code! That's ok! it worked.... kind of

- There's a loop, processing some user input
- A template which generates the boiler plate of an executable go program
- Some exec.Command().Outputs() which:
  - Run Goimports on the generated Go code (Thank god that's part of the Go toolchain)
  - run the generated code, and reports the output back

This quick 30 min hack was enough validation to continue working. So I decided to begin work on parsing out functions from raw user input

## Parsing Functions

At the moment of writing this, I'm parsing the function by regular expressions.

The in memory representation of a function is:

```go
type Function struct {
    // Raw, the raw input which was determined to be a function
    Raw string

    // Identifier, the identifier of the function.
    // For example:
    // func a() {
    //
    // }
    // Identifier would = a
    Identifier string

    // Params is a raw string which identifies the parameters of the function
    Params string

    // Return is the return signature of the function
    Return string
}
```

And I have some hacked together regexes to parse each piece of a function:

```go
const (
    // holds double duty for classifying, and extracting functions from raw text
    isFunctionExpr regexpType = iota
    identifierExpr
    argsExpr
    returnExpr
)

var expressions = map[regexpType]*regexp.Regexp{
    isFunctionExpr: regexp.MustCompile(`func \(?.*\)?\{\n?(.*|\s|\S)*?(\})`),
    identifierExpr: regexp.MustCompile(`(func .* \(|func .*?)\(`),
    argsExpr:       regexp.MustCompile(`\((.*?)\)`),
    returnExpr:     regexp.MustCompile(`\) .* {`),
}
```

Using the 4 above expression I'm able to piece together each part of a function and store it in memory as a Map

For example,

```go
isFunctionExpr: regexp.MustCompile(`func \(?.*\)?\{\n?(.*|\s|\S)*?(\})`),
```

will classify some text as **_is a function_** or **_is not a function_**

for instance:

```go
func hello() int {

}
```

**would** parse, while

```go
function() string {
}
```

**wouldn't** parse

Similarly, the other 3 expressions are able to parse different parts of a function. I'm admittedly not the best at regular expressions so I'm sure there are cleaner ways if coming to the same solution. So this likely won't stay the same forever. But it works pretty well for now.

After a couple days of hacking on iGo I was able to classify functions, store them in memory, look them up, and generate them on the fly. And as a bonus, since Go already has a great ecosystem Goimports is able to look up references to 3rd party packages; that's awesome!

## Read, Interpret, Eval

The Interpreter struct is very basic. It holds a map of references to functions, and a history of text the interpreter has seen.

```go
// Interpreter houses the function references and input history
type Interpreter struct {
    Functions map[string]*parse.Function
    History   []string
}
```

Every time a function is recognized as such, it will be placed in Interpreter.Functions (keyed by its ID), that looks like this:

```go
// Interpret will take some text and
// Classify it as either an expression or a function
// If it is a function it will store the reference of the Function in a map
// If the text is classified as an expression, it will evaluate the expression,
// using the function reference map if needed
func (i *Interpreter) Interpret(text string) {
    if i.Functions == nil {
        i.Functions = make(map[string]*parse.Function)
    }
    i.History = append(i.History, text)
    t := i.classify(text)
    for _, tv := range t {
        switch v := tv.(type) {
        case *parse.Function:
            i.Functions[v.Identifier] = v
            fmt.Printf("# %s\n", v.String())
            break
        case *parse.Expression:
            fmt.Printf(">> %s\n", text)
            i.Eval(text)
            break
        }
    }
}
```

By default, we don't eval function declarations right away, they just sit there until they're needed

Once the user invokes the interpreter with a function they declared earlier, the interpreter will then attempt to eval the function.

```go
case *parse.Expression:
            fmt.Printf(">> %s\n", text)
            i.Eval(text)
            break
        }
```

And the eval code is the same in my 30 min hacked together program, just refactored to use the core interpreter implementation

## Decoupling the Interpreter From Input

When a started this project I began with a simple Terminal based input system. That worked, relatively well but wasn't much fun. I quickly outgrew that and the quickest way to stress test my function parsing was to allow a more robust transfer mechanism from _User_ to _Interpreter_. I wrote a simple HTTP Server which accepts requests and forwards them to the interpreter. Then spits out the result of the code which was interpreted.

That meant I needed to refactor; I began splitting every piece into a module and I ended up with a package structure like

- pkg
  - parse
  - interpreter

parse handles the text parsing and classification, and interpreter utilizes the parsed data to implement the interpreter. Once that was done it was easy to accept Client Server Communication.

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/beeceej/iGo/pkg/interpreter"
    "github.com/davecgh/go-spew/spew"
)

func main() {
    i := interpreter.Interpreter{}
    http.HandleFunc("/interpret", func(w http.ResponseWriter, r *http.Request) {
        b, _ := ioutil.ReadAll(r.Body)
        defer r.Body.Close()
        var m struct {
            Text string `json:"text"`
        }
        json.Unmarshal(b, &m)
        i.Interpret(m.Text)
    })

    if err := http.ListenAndServe(":9999", nil); err != nil {
        fmt.Println(err.Error())
        spew.Dump(i)
    }
}
```

Notice the last line, if the server fails for any reason, we are able to print out the exact state of the interpreter at the time of crash... (`spew.Dump(i)`). This is a very powerful concept and lends itself very nicely to programmatic use! For instance one could read from the file system to hydrate state. One could declare a set of default functions in a config file.

## Integrations

Since we're able to run the interpreter as a server, we can build many unique clients for it. Just write a plugin for your editor of choice and you can have much the same functionality that interpreted language users use and love on a daily basis.

Keep in mind this project is under active development, What's written here will likely change! to follow the development keep an eye out for further posts! Or, if you're interested in iGo you can follow the development [here](https://github.com/beeceej/iGo).

Part 1 of **_Writing a Go Interpreter in Go_**
