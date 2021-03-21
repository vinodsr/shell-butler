package ui

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	ui "github.com/gizak/termui/v3"
	addcli "github.com/vinodsr/shell-butler/lib/cli/add"
	"github.com/vinodsr/shell-butler/lib/config"
	runtime "github.com/vinodsr/shell-butler/lib/runtime"
	lib "github.com/vinodsr/shell-butler/lib/types"

	"github.com/gizak/termui/v3/widgets"
)

// Render the layout
func Render() {

	var selectedContext []string
	rt := runtime.GetRunTime()
	commandStr := ""
	var filteredDataSource []string
	var commandLevel int = 1

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	// Load the widget components

	contextTextBox := widgets.NewParagraph()
	contextTextBox.Title = "Search"
	contextTextBox.Text = formatCommandString(selectedContext, commandStr)
	contextTextBox.SetRect(0, 0, termWidth, 3)
	contextTextBox.TextStyle.Fg = ui.ColorGreen
	contextTextBox.BorderStyle.Fg = ui.ColorCyan

	alertBox := widgets.NewParagraph()
	alertBox.Title = "Alert"
	alertBox.Text = ""
	alertBox.SetRect((termWidth/2)-20, (termHeight/2)-2, (termWidth/2)+20, (termHeight/2)+2)
	alertBox.TextStyle.Fg = ui.ColorRed
	alertBox.BorderStyle.Fg = ui.ColorRed

	alertText := ""

	debugBox := widgets.NewParagraph()
	debugBox.Title = "Debug"
	debugBox.Text = ""
	debugBox.SetRect(0, termHeight-10, termWidth, termHeight)

	contextListBox := widgets.NewList()
	filteredDataSource = displayContextatLevel(updateContextList(rt.GetContextDS(), ""), commandLevel)
	contextListBox.Rows = filteredDataSource
	contextListBox.Title = "Commands"
	contextListBox.TextStyle = ui.NewStyle(ui.ColorWhite)
	contextListBox.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen, ui.ModifierBold)
	contextListBox.WrapText = false
	contextListBox.SetRect(0, 6, termWidth, termHeight)
	contextListBox.BorderStyle.Fg = ui.ColorGreen

	addButton := widgets.NewParagraph()
	addButton.SetRect(0, 3, 20, 6)
	addButton.Text = "Ins - Add Command"
	addButton.BorderStyle.Bg = ui.ColorBlue

	deleteButton := widgets.NewParagraph()
	deleteButton.SetRect(24, 3, 40, 6)
	deleteButton.Text = "Del - Delete Command"
	deleteButton.BorderStyle.Bg = ui.ColorRed

	exitButton := widgets.NewParagraph()
	exitButton.SetRect(44, 3, 60, 6)
	exitButton.Text = "Esc - Quit"
	exitButton.BorderStyle.Bg = ui.ColorCyan

	staicDisplayItems := []ui.Drawable{contextListBox, contextTextBox, addButton, deleteButton, exitButton}

	ui.Render(staicDisplayItems...)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			ui.Close()
			os.Exit(1)
			return
		case "j", "<Down>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollDown()
			}
		case "k", "<Up>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollUp()
			}
		case "<C-d>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollHalfPageDown()
			}
		case "<C-u>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollHalfPageUp()
			}
		case "<C-f>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollPageDown()
			}
		case "<C-b>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollPageUp()
			}
		case "<Home>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollTop()
			}
		case "G", "<End>":
			if len(contextListBox.Rows) > 0 {
				contextListBox.ScrollBottom()
			}
		case "<Backspace>":
			if len(commandStr) > 0 {
				commandStr = commandStr[0 : len(commandStr)-1]
			} else if len(commandStr) == 0 && commandLevel > 1 {
				// clear the previous selection

				commandLevel--
				selectedContext = selectedContext[:len(selectedContext)-1]

			}
		case "<Resize>":
			termWidth, termHeight := ui.TerminalDimensions()
			contextListBox.SetRect(0, 3, termWidth, termHeight)
			contextTextBox.SetRect(0, 0, termWidth, 3)

		case "<Space>":
			commandStr += " "
		case "<Insert>":
			ui.Close()
			addcli.Execute()
			rt.Init()
			ui.Init()
		case "<Delete>":
			toDeleteContext := append(selectedContext, filteredDataSource[contextListBox.SelectedRow])
			_joinedContext := strings.Join(toDeleteContext, ":")

			shouldDelete := ConfirmBox(fmt.Sprintf("Are you sure to delete %s  ?", _joinedContext), uiEvents)
			if shouldDelete {
				if contextListBox.SelectedRow >= 0 && len(filteredDataSource) > 0 {
					_joinedContext += ":"
					debug("JOined context = "+_joinedContext, debugBox)
					var newCommandList []lib.Command
					deleted := false
					for _, configEntry := range rt.GetConfig().Commands {
						debug(configEntry.Context, debugBox)
						if strings.HasPrefix(configEntry.Context, _joinedContext) {
							debug("Found Context", debugBox)
							deleted = true
						} else {
							newCommandList = append(newCommandList, configEntry)
						}
					}
					debug("Should replace ? "+strconv.FormatBool(deleted), debugBox)
					if deleted {
						newConfig := rt.GetConfig()
						newConfig.Commands = newCommandList
						config.WriteConfig(newConfig)
						rt.Init()
					}

				}
				break
			}
		case "<Enter>":
			// Now take the context available in the select box
			if contextListBox.SelectedRow >= 0 && len(filteredDataSource) > 0 {
				selectedContext = append(selectedContext, filteredDataSource[contextListBox.SelectedRow])

				_joinedContext := strings.Join(selectedContext, ":") + ":"
				debug("JOined context = "+_joinedContext, debugBox)
				if rt.GetCommandMap()[_joinedContext] != "" {
					fmt.Println(rt.GetCommandMap()[_joinedContext])
					ui.Clear()
					ui.Close()
					os.Exit(0)
				}
				// Check if this is a terminal command ?

				debugBox.Text = fmt.Sprint()
				commandLevel++
				commandStr = ""
				contextListBox.SelectedRow = 0
			} else {
				alertText = "No command found"
			}

		default:
			if len(e.ID) == 1 {
				commandStr += e.ID
			}
		}

		//fmt.Println(emoji.Squid)
		//+ emoji.Squid.String()

		contextTextBox.Text = formatCommandString(selectedContext, commandStr)
		contextSep := ""
		if len(selectedContext) > 0 {
			contextSep = ":"
		}

		filterText := "(?i)^" + strings.Join(selectedContext, ":") + contextSep + "[^:]*" + commandStr + ".*\\:?"
		//contextTextBox.Text = filterText
		for _, cc := range rt.GetContextDS() {
			debug(" context main DS = "+cc, debugBox)

		}
		filteredDataSource = displayContextatLevel(updateContextList(rt.GetContextDS(), filterText), commandLevel)
		contextListBox.Rows = filteredDataSource
		if alertText == "" && len(filteredDataSource) == 0 {
			alertText = "No match found for : " + commandStr
			if len(selectedContext) > 0 {
				alertText += " in " + strings.Join(selectedContext, " \u2b95 ")
			}
		}
		debug(filterText, debugBox)

		if len(rt.GetContextDS()) == 0 {
			alertText = " No commands found. Please add one"
		}
		displayItems := staicDisplayItems
		if alertText != "" {
			alertBox.Text = alertText
			displayItems = append(displayItems, alertBox)
		}
		//displayItems = append(displayItems, debugBox)

		ui.Clear()
		ui.Render(displayItems...)
		alertText = ""
	}
}
