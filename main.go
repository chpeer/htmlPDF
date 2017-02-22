package htmlPDF

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
)

//global pointer to pdf
var pdf *gofpdf.Fpdf

func Generate(in string, out string) {
	fmt.Println(in, out)
	xmlFile, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	//parse html to Node tree
	n := ParseHtml(string(xmlFile))
	n.print(0)

	cssFile, err := ioutil.ReadFile("style.css")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	cssStyle := string(cssFile)
	fmt.Println(cssStyle)
	//todo change NewParser to ParseCSS with return *Rules
	p2 := NewParser(cssStyle)
	stylesheet := p2.parseRules()

	styletree := styleTree(n, &stylesheet)
	fmt.Println("=============")
	fmt.Println(styletree)
	fmt.Println("=============")

	viewport := Dimensions{}
	viewport.content.width = 800
	viewport.content.height = 600

	layoutTree := layoutTree(styletree, viewport)
	fmt.Printf("%+v\n", layoutTree)
	list := buildDisplayList(layoutTree)
	fmt.Println("==========DisplayCommand============")
	fmt.Println(list)
	//list := map[int]DisplayCommand{}

	rect := Rect{10, 10, 50, 50}
	color := Color{255, 0, 0, 0}
	//list[0] = DisplayCommand{SolidColor{color: color, rect: rect}}

	rect.x += 10
	rect.y += 10
	color.r = 0
	color.g = 255
	list[1] = DisplayCommand{SolidColor{color: color, rect: rect}}

	fmt.Println("list")
	fmt.Println(list)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	for i := 0; i < len(list); i++ {
		list[i].draw(pdf)
	}
	err = pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Println("Error pdf", err)
	}
	pdf.Close()
}
