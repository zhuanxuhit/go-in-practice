package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/robertkrimen/otto"
)

func testMarkdown() {
	extensions := parser.CommonExtensions | parser.MathJax
	myparser := parser.NewWithExtensions(extensions)

	md := []byte(`
	$$
	\left[ \begin{array}{a} a^l_1 \\ ⋮ \\ a^l_{d_l} \end{array}\right]
	= \sigma(
	\left[ \begin{matrix}
	w^l_{1,1} & ⋯  & w^l_{1,d_{l-1}} \\
	⋮ & ⋱  & ⋮  \\
	w^l_{d_l,1} & ⋯  & w^l_{d_l,d_{l-1}} \\
	\end{matrix}\right]  ·
	\left[ \begin{array}{x} a^{l-1}_1 \\ ⋮ \\ ⋮ \\ a^{l-1}_{d_{l-1}} \end{array}\right] +
	\left[ \begin{array}{b} b^l_1 \\ ⋮ \\ b^l_{d_l} \end{array}\right])
	$$
		`)
	output := markdown.ToHTML(md, myparser, nil)
	fmt.Printf("%q", output)
}

func main() {
	vm := otto.New()
	vm.Run(`
    abc = 2 + 2;
    console.log("The value of abc is " + abc); // 4
`)
}
