package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

type Node struct {
	name   string
	text   string
	parent *Node
	child  []*Node
}

func (n *Node) Start(name string) *Node {
	tmp := new(Node)
	tmp.parent = n
	tmp.name = name

	return tmp
}

func (n *Node) AddText(text string) {
	n.text = text
}

func (n *Node) Stop() *Node {
	n.parent.child = append(n.parent.child, n)
	return n.parent
}

func (n *Node) Print(level int) {
	tab(level)
	level++
	fmt.Printf("<%s>\n", n.name)
	if len(n.text) > 0 {
		tab(level + 1)
		fmt.Printf("%s\n", n.text)
	}
	if len(n.child) > 0 {
		for i := 0; i < len(n.child); i++ {
			n.child[i].Print(level)
		}
	}
	level--
	tab(level)
	fmt.Printf("</%s>\n", n.name)
}

var pageFontSize float64 = 12
var pageFontFamily string = "Helvetica"
var pageFontColor string = "black"
var pageFontStyle string = "" //B - bold or I - italic or U - underscore

func (n *Node) PrintSelf(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.SetFontSize(16)
	fmt.Println(pdf.GetXY())
	pdf.SetXY(0, 0)
	fmt.Println(pdf.GetXY())
	n.pdf(pdf)

	return pdf
}

var newline bool

func setNewLine(pdf *gofpdf.Fpdf) {
	if !newline {
		y := pdf.GetY()
		_, fontSize := pdf.GetFontSize()
		pdf.SetY(y + fontSize)
		newline = true
	}
}

func (n *Node) pdf(pdf *gofpdf.Fpdf) {
	parseChilder := true
	pdf.SetFont(pageFontFamily, pageFontStyle, pageFontSize)
	switch n.name {
	case "div":
		setNewLine(pdf)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		drawHX(pdf, n)
	case "b":
		pdf.SetFont(pageFontFamily, "B", pageFontSize)
		pdf.Text(pdf.GetX(), pdf.GetY(), n.text)
		x := pdf.GetX()
		pdf.SetX(x + pdf.GetStringWidth(n.text) + 1)
		newline = false
	case "table":
		setNewLine(pdf)
		drawTable(pdf, n)
		parseChilder = false
	}
	if parseChilder {
		for i := 0; i < len(n.child); i++ {
			n.child[i].pdf(pdf)
		}
	}
}

var lastMargin float64

func drawHX(pdf *gofpdf.Fpdf, n *Node) {
	var fontSize float64
	var marginTop float64
	var marginBottom float64
	marginTop = 6
	marginBottom = 6

	switch n.name {
	case "h1":
		fontSize = 18
	case "h2":
		fontSize = 16
	case "h3":
		fontSize = 14
	case "h4", "h5", "h6":
		fontSize = 12
	}
	if marginTop > lastMargin {
		marginTop -= lastMargin
	} else {
		marginTop = 0
	}
	pdf.SetFontSize(fontSize)
	y := pdf.GetY()
	_, fontSizeMM := pdf.GetFontSize()
	pdf.SetY(y + fontSizeMM + marginTop)

	pdf.Text(pdf.GetX(), pdf.GetY(), n.text)

	y = pdf.GetY()
	pdf.SetY(y + marginBottom)
	lastMargin = marginBottom
}

func drawTable(pdf *gofpdf.Fpdf, n *Node) {
	t := make(map[int]float64)
	for i := 0; i < len(n.child); i++ {
		tmp := n.child[i]
		for j := 0; j < len(tmp.child); j++ {
			stringSize := pdf.GetStringWidth(tmp.child[j].text)
			if t[i] < stringSize {
				t[i] = stringSize
			}
		}
	}
	fmt.Println(t)
	for i := 0; i < len(n.child); i++ {
		drawTr(pdf, n.child[i], t)
	}
	y := pdf.GetY()
	_, fontSize := pdf.GetFontSize()
	pdf.SetY(y + fontSize)
}

func drawTr(pdf *gofpdf.Fpdf, n *Node, t map[int]float64) {
	y := pdf.GetY()
	_, fontSize := pdf.GetFontSize()
	pdf.SetY(y + fontSize)

	for i := 0; i < len(n.child); i++ {
		pdf.Text(pdf.GetX(), pdf.GetY(), n.child[i].text)
		x := pdf.GetX()
		pdf.SetX(x + t[i] + 1)
	}
}