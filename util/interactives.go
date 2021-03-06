package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/bcvery1/Deadication/bindata"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	interactivePlacementPath string = "assets/interactiveplacement.csv"
)

var (
	titleV = pixel.V(25, 155)
	menuV  = pixel.V(40, 155)
)

var field1 = field{
	Interactive: Interactive{"Top field", false},
	crop:        NewCrop(""),
	planted:     true,
	amountLeft:  0,
}
var field2 = field{
	Interactive: Interactive{"Mid field", false},
	crop:        NewCrop(""),
	planted:     false,
	amountLeft:  0,
}
var field3 = field{
	Interactive: Interactive{"Bottom field", false},
	crop:        NewCrop(""),
	planted:     false,
	amountLeft:  0,
}

var river = riverInter{
	Interactive: Interactive{"River", false},
}

var pen1 = pen{
	Interactive: Interactive{"Top pen", false},
	humans:      []*Human{},
	food:        0,
}
var pen2 = pen{
	Interactive: Interactive{"Mid pen", false},
	humans:      []*Human{},
	food:        0,
}
var pen3 = pen{
	Interactive: Interactive{"Bottom pen", false},
	humans:      []*Human{},
	food:        2,
}

var house = NewHouse()

var (
	// Pens shares the pens to the gameloop
	Pens = map[string]*pen{
		"Top pen":    &pen1,
		"Mid pen":    &pen2,
		"Bottom pen": &pen3,
	}
	// Fields shares the fields to the gameloop
	Fields = map[string]*field{
		"Top field":    &field1,
		"Mid field":    &field2,
		"Bottom field": &field3,
	}
)

type option struct {
	text string
}

func (o *option) Text() string {
	return o.text
}
func (o *option) Action(i InteractiveI, carrying string) {
	log.Printf("Acting with %s", carrying)
}

type optionI interface {
	Text() string
	Action(InteractiveI, string)
}

// Interactive action for in game objects
type Interactive struct {
	title  string
	active bool
}

// InteractiveI interface for any Interactive element
type InteractiveI interface {
	opts(string) []optionI
	Activate(string)
	Deactivate()
	Title() string
	IsActive() bool
	Update(*pixelgl.Window, string)
}

// Update updates the interactive in the game world
func (i *Interactive) Update(win *pixelgl.Window, carrying string) {
	if !i.IsActive() {
		return
	}

	// Draw box
	imd := getBox()
	imd.Draw(win)

	// Draw title
	title, scale := getText(-1, i.Title(), 1.4, titleV)
	title.Draw(win, scale)

	shiftV := pixel.V(0, -20)
	for j, opt := range i.opts(carrying) {
		v := menuV.Sub(shiftV.Scaled(float64(j)))
		optTxt, scale := getText(j+1, opt.Text(), 1.1, v)
		optTxt.Draw(win, scale)
	}
}

func getBox() *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = colornames.Whitesmoke
	imd.Push(pixel.V(0, 0), pixel.V(300, 175))
	imd.Rectangle(0)

	return imd
}

func getAtlas() *text.Atlas {
	return text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func getText(i int, txt string, scale float64, v pixel.Vec) (*text.Text, pixel.Matrix) {
	atlas := getAtlas()
	outputText := text.New(v, atlas)
	outputText.Color = colornames.Black
	if i < 0 {
		fmt.Fprintf(outputText, txt)
	} else {
		fmt.Fprintf(outputText, "%d - %s", i, txt)
	}

	return outputText, pixel.IM.Scaled(outputText.Orig, scale)

}

// IsActive returns whether this interactive is currently active
func (i *Interactive) IsActive() bool {
	return i.active
}

// Title returns the interactives' title
func (i *Interactive) Title() string {
	return i.title
}

func (i *Interactive) opts(c string) []optionI {
	return []optionI{}
}

// Activate the structs function
// Takes what the player is currently carrying
func (i *Interactive) Activate(carrying string) {
	i.active = true
}

// Deactivate stops the interactives behaivour
func (i *Interactive) Deactivate() {
	i.active = false
}

// AllInteractives gets all interactive entities and collision areas
func AllInteractives() (map[string]InteractiveI, map[pixel.Rect]string) {
	r := make(map[pixel.Rect]string)
	interactiveF, err := bindata.Asset(interactivePlacementPath)
	if err != nil {
		log.Fatal(err)
	}

	csvFile := csv.NewReader(bytes.NewReader(interactiveF))
	for {
		i, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		xMin, _ := strconv.ParseFloat(i[0], 64)
		yMin, _ := strconv.ParseFloat(i[1], 64)
		xMax, _ := strconv.ParseFloat(i[2], 64)
		yMax, _ := strconv.ParseFloat(i[3], 64)
		interName := i[4]
		rect := pixel.R(xMin, yMin, xMax, yMax)
		r[rect] = interName

		// Set the pen rects
		if interName == "pen1" {
			pen1.rect = rect
		} else if interName == "pen2" {
			pen2.rect = rect
		} else if interName == "pen3" {
			pen3.rect = rect
		}
		// Set the field rects
		if interName == "field1" {
			field1.rect = rect
		} else if interName == "field2" {
			field2.rect = rect
		} else if interName == "field3" {
			field3.rect = rect
		}
	}

	m := make(map[string]InteractiveI)
	m["field1"] = &field1
	m["field2"] = &field2
	m["field3"] = &field3
	m["pen1"] = &pen1
	m["pen2"] = &pen2
	m["pen3"] = &pen3
	m["river"] = &river
	m["house"] = house

	return m, r
}

func doOptions(win *pixelgl.Window, optionlist []optionI, carrying string, i InteractiveI) {
	if win.JustPressed(pixelgl.Key1) {
		if len(optionlist) > 0 {
			optionlist[0].Action(i, carrying)
		}
	}
	if win.JustPressed(pixelgl.Key2) {
		if len(optionlist) > 1 {
			optionlist[1].Action(i, carrying)
		}
	}
	if win.JustPressed(pixelgl.Key3) {
		if len(optionlist) > 2 {
			optionlist[2].Action(i, carrying)
		}
	}
	if win.JustPressed(pixelgl.Key4) {
		if len(optionlist) > 3 {
			optionlist[3].Action(i, carrying)
		}
	}
	if win.JustPressed(pixelgl.Key5) {
		if len(optionlist) > 4 {
			optionlist[4].Action(i, carrying)
		}
	}
}
