package ui

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func ConfirmBox(msg string, uiEvents <-chan ui.Event) bool {

	alertBox := widgets.NewParagraph()
	alertBox.Title = "Alert"
	alertBox.Text = msg + " [y/n] "
	alertBox.SetRect(5, 5, 100-5, 8)
	alertBox.TextStyle.Fg = ui.ColorRed
	alertBox.BorderStyle.Fg = ui.ColorRed
	ret := true
	fmt.Println(msg)
	ui.Render(alertBox)
	loop := true
	for loop {
		e := <-uiEvents

		switch e.ID {
		case "n", "N":
			loop = false
			ret = false
			break
		case "y", "Y":
			loop = false
			ret = true
			break
		}
		ui.Render(alertBox)

	}
	return ret
}
