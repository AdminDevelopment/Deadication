package util

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type pen struct {
	Interactive
	humans []*human
}

func (p *pen) Update(win *pixelgl.Window, carrying string) {
	if !p.IsActive() {
		return
	}

	// Draw box
	imd := getBox()
	imd.Draw(win)

	// Draw title
	title, scale := getText(-1, p.Title(), 1.4, titleV)
	title.Draw(win, scale)

	shiftV := pixel.V(0, 30)
	penoptions := p.opts(carrying)
	for j, opt := range penoptions {
		v := menuV.Sub(shiftV.Scaled(float64(j + 1)))
		optTxt, scale := getText(j+1, opt.Text(), 1.1, v)
		optTxt.Draw(win, scale)
	}

	// Check if the user presses a number key to select an option
	doOptions(win, penoptions, carrying, p)
}

func (p *pen) opts(c string) []optionI {
	o := observePen{option{"Observe pen"}}
	opts := []optionI{&o}

	if c == "food" {
		o := feedHumans{option{"Feed humans"}}
		opts = append(opts, &o)
	}

	if c == "cloth" {
		o := clotheHumans{option{"Give humans cloth for warmth"}}
		opts = append(opts, &o)
	}

	if len(p.humans) > 0 {
		o := eatBrain{option{"Eat a human"}}
		opts = append(opts, &o)
	}

	return opts
}

type observePen struct {
	option
}

func (o *observePen) Action(p InteractiveI, carrying string) {
	PopupChan <- &Popup{"This pen holds humans for eating!"}
}

type feedHumans struct {
	option
}

func (f *feedHumans) Action(p InteractiveI, carrying string) {
	PickupChan <- ""
	s := fmt.Sprintf("You gave food to the humans in %s", p.Title())
	PopupChan <- &Popup{s}
}

type clotheHumans struct {
	option
}

func (c *clotheHumans) Action(p InteractiveI, carrying string) {
	s := fmt.Sprintf("You gave clothes to the humans in %s", p.Title())
	PopupChan <- &Popup{s}
	PickupChan <- ""
}

type eatBrain struct {
	option
}

func (e *eatBrain) Action(p InteractiveI, carrying string) {
	PopupChan <- &Popup{"You ate some brains!  Yum!"}
	EatChan <- 50
}