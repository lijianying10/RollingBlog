title: Metaprogramming and code generation/replace in action
date: 2021-03-14 22:01:54
categories: 技术
tags: [metaprogramming]
---

## What is Metaprogramming

Metaprogramming, for short, is a program that can treat other programs as data. It can be designed to read, generate, analyze or transform other programs. [Detail](https://en.wikipedia.org/wiki/Metaprogramming)

In this article, we concentrate on static code not in the process running phase.

## Macro

Macro means how a code block should be mapped to a replacement output. [Detail](https://en.wikipedia.org/wiki/Macro_(computer_science)). Macro is a commonly used tool for the implementation of MetaProgramming.

### Macro in C

Suppose you are new to the C language. You may don't know you are already use Metaprogramming in action. Such as the following example:

some examples:

e.g 1:

```
#include <stdio.h>
int main() {
   printf("Hello, World!");
   return 0;
}
```

Here we got `#include` is a macro to import function definition from the header file.

e.g 2:

```
#define PI 3.14159
```

This code fragment will cause the string "PI" to be replaced with "3.14159", we call it parameterized macro.

e.g 3:

```
#if WIN32
	#include	<winsock.h>
#elif LINUX
	#include	<sys/socket.h>
#endif
```

A macro could read variables. Different target platforms use specific herder files.

e.g 4:

```
#define pred(x)  ((x)-1)
```

We can easily take advantage of the `inline function` definition through the definition of a Macro. We don't care about the type of x cause the code needs to expand before compile. Moreover, define Macro as a function without any type definition. It can provide you lots of flexibility to make your project more clean and readable.

the C language toolchain compiles your source code through the following workflow in detail:

![](https://user-images.githubusercontent.com/3077762/111071428-1d73d200-8511-11eb-9abd-23d41a801d61.png)

### summary

As the previous example showing that Metaprogramming is not a strange thing. We often use it in daily programming.

## Code replace in action

### AST parser 

For actually and safety code replace, we need to study two new concepts: `AST` and `AST parser`.

AST is the abbreviation of Abstract Syntax Tree, a kind of tree data structure. Each node of the tree denotes a construct occurring in the source code.
[Detail](https://en.wikipedia.org/wiki/Abstract_syntax_tree)

AST parser is a library or package for parsing a source code file or fragment to a tree data structure tool. Programming language specification too complex to replace by simple string replacement. It can not cover most usage cases. We use the `AST parser` tool to search or operation a tree structure to meet the requirements.

example parser for jsx

``` javascript
babelParser = require('@babel/parser');
const res = babelParser.parse(`
ReactDOM.render(
  <h1>Hello, world!</h1>,
  document.getElementById('root')
);`,{plugins:[ 'jsx', 'flow' ]});
console.log(JSON.stringify(res));
```

Then we got a very large JSON to tell us every detail of this code. [pastebin: detail](https://pastebin.com/AAch75nV)

``` json
{
  "type": "CallExpression",
  "start": 1,
  "end": 79,
  "loc": {
    "start": {
      "line": 2,
      "column": 0
    },
    "end": {
      "line": 5,
      "column": 1
    }
  },
  "callee": {
    "blah":"blah"
  }
}
```

Here's part of the result. As an example, we can get the information that the expression is call function expression and get the code's position.
Even the location line and column information and tell us every `callee` (every member of function call parameter)

The parser is a part of the language compiler (or part of the compiler frontend).
The parser is often used as an infrastructure for code static analysis and code automation complete in an editor.

More example:

1. Markdown-parser https://www.npmjs.com/package/markdown-parser
2. yaml-parser https://www.npmjs.com/package/js-yaml

### code scanning for ci

After parsing the source code to the syntax tree, we can write a deep search first algorithm for specific coding behavior.
Then deploy the algorithm to CI workflow before compile source code.

Here's an example for checkout whether any code behavior returns an error but doesn't return the None 200 status code.
The function implementation based on if we check error is not `nil` there must call some function `w.Writeheader(xxx)`, and the parameter may not be 2xx

We got this requirement from the middleware team, who develop a web load balance and monitor the website's error status.

For shorter and straightforward we made a small example only check a function is HTTP handler:

``` golang
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// src is the input for which we want to print the AST.
	src := `
package main
func (rt *Runtime) FileReceive(w http.ResponseWriter, r *http.Request) {
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

  // Print the AST.
  ast.Print(fset, f)

	ast.Inspect(f, func(n ast.Node) bool { // deep first node visitor
		switch x := n.(type) {
		case *ast.FuncType: // check node is a function
			if len(x.Params.List) != 2 { // check function parameter length
				return false
			}
			p1, p2 := false, false
			if val, ok := x.Params.List[0].Type.(*ast.SelectorExpr); ok { // check parameter type is http.ResponseWriter
				if valExpr, ok := val.X.(*ast.Ident); ok {
					if valExpr.Name == "http" && val.Sel.Name == "ResponseWriter" {
						p1 = true
					}
				}
			}
			if val, ok := x.Params.List[1].Type.(*ast.StarExpr); ok { // check parameter type is *http.Request
				if starX, ok := val.X.(*ast.SelectorExpr); ok {
					if starXX, ok := starX.X.(*ast.Ident); ok {
						if starXX.Name == "http" && starX.Sel.Name == "Request" {
							p2 = true
						}
					}
				}
			}
			if p1 && p2 {
				fmt.Println("handler found")
			}
		}
		return true
	})
}
```

Tips: type assertation following the printed AST node tree.

### example of replacing code

Follow the path got safe and accurately replace code result:

1. reuse previously mentioned search algorithm methodology
2. direct operation to the syntax tree
3. compile the syntax tree to source code
4. format code makes it looks pretty cool.

Use case: When you got a GO project to upgrade the original log module (fmt.Println) to a newly designed log module (log.Log())

and example code for golang (by golang.org/x/tools/go/ast/astutil) https://play.golang.org/p/jDSpmV_Kxnt

### example of code gen

For the example of code generation, we have already have lots of widely used tools:

#### protobuf

GRPC is a common tool for [RPC](https://en.wikipedia.org/wiki/Remote_procedure_call). The message encoding uses the standard named protobuf.

here is an example for code generation use following protobuf file example fill the left blank: https://protogen.marcgravell.com/

example protobuf file:

``` protobuf
syntax = "proto3";

message ExampleMessage {
    int32 foo = 1;
    string bar = 3;
}
```

#### sqlc

To fully control the SQL query and performance, we don't use ORM but code generation.
This code generation tool can help us a lot for those repeat work on translate SQL results to the local data structure.

Example convert SQL file to golang DAO layer: https://play.sqlc.dev/

## In summary

Metaprogramming is a commonly used method for programming. It can reduce repeat work and improve engineering quality.

Learning, sharing, and improve together.