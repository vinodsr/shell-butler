package ui

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/vinodsr/shell-butler/lib/config"
	runtime "github.com/vinodsr/shell-butler/lib/runtime"
	lib "github.com/vinodsr/shell-butler/lib/types"

	"github.com/gizak/termui/v3/widgets"
)

// Render the layout
func Render() {

	var selectedContext []string
	rt := runtime.GetRunTime()
	commandMap := rt.GetCommandMap()
	contextDataSource := rt.GetContextDS()
	configData := rt.GetConfig()
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
	contextTextBox.Title = "Command"
	contextTextBox.Text = formatCommandString(selectedContext, commandStr)
	contextTextBox.SetRect(0, 0, termWidth, 3)
	contextTextBox.TextStyle.Fg = ui.ColorGreen
	contextTextBox.BorderStyle.Fg = ui.ColorCyan

	alertBox := widgets.NewParagraph()
	alertBox.Title = "Alert"
	alertBox.Text = ""
	alertBox.SetRect(5, 5, termWidth-5, 8)
	alertBox.TextStyle.Fg = ui.ColorRed
	alertBox.BorderStyle.Fg = ui.ColorRed

	alertText := ""

	debugBox := widgets.NewParagraph()
	debugBox.Title = "Debug"
	debugBox.Text = ""
	debugBox.SetRect(0, termHeight-10, termWidth, termHeight)

	contextListBox := widgets.NewList()
	filteredDataSource = displayContextatLevel(updateContextList(contextDataSource, ""), commandLevel)
	contextListBox.Rows = filteredDataSource
	contextListBox.TextStyle = ui.NewStyle(ui.ColorWhite)
	contextListBox.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen, ui.ModifierBold)
	contextListBox.WrapText = false
	contextListBox.SetRect(0, 3, termWidth, termHeight)
	contextListBox.BorderStyle.Fg = ui.ColorGreen

	ui.Render(contextTextBox, contextListBox)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		switch e.ID {
		case "q", "<C-c>":
			ui.Close()
			os.Exit(1)
			return
		case "j", "<Down>":
			contextListBox.ScrollDown()
		case "k", "<Up>":
			contextListBox.ScrollUp()
		case "<C-d>":
			contextListBox.ScrollHalfPageDown()
		case "<C-u>":
			contextListBox.ScrollHalfPageUp()
		case "<C-f>":
			contextListBox.ScrollPageDown()
		case "<C-b>":
			contextListBox.ScrollPageUp()
		case "<Home>":
			contextListBox.ScrollTop()
		case "G", "<End>":
			contextListBox.ScrollBottom()
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
		case "<Delete>":
			shouldDelete := ConfirmBox("Are you to delete the entry ?", uiEvents)
			if shouldDelete {
				if contextListBox.SelectedRow >= 0 && len(filteredDataSource) > 0 {
					toDeleteContext := append(selectedContext, filteredDataSource[contextListBox.SelectedRow])

					_joinedContext := strings.Join(toDeleteContext, ":")
					debug("JOined context = "+_joinedContext, debugBox)
					var newCommandList []lib.Command
					deleted := false
					for _, configEntry := range configData.Commands {
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
						configData.Commands = newCommandList
						config.WriteConfig(configData)
						rt.Init(configData)
						commandMap = rt.GetCommandMap()
						contextDataSource = rt.GetContextDS()
						for _, cc := range contextDataSource {
							debug(" context = "+cc, debugBox)

						}
						configData = rt.GetConfig()
					}

				}
				break
			}
		case "<Enter>":
			// Now take the context available in the select box
			if contextListBox.SelectedRow >= 0 && len(filteredDataSource) > 0 {
				selectedContext = append(selectedContext, filteredDataSource[contextListBox.SelectedRow])

				_joinedContext := strings.Join(selectedContext, ":")
				debug("JOined context = "+_joinedContext, debugBox)
				if commandMap[_joinedContext] != "" {
					fmt.Println(commandMap[_joinedContext])
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
		for _, cc := range contextDataSource {
			debug(" context main DS = "+cc, debugBox)

		}
		filteredDataSource = displayContextatLevel(updateContextList(contextDataSource, filterText), commandLevel)
		contextListBox.Rows = filteredDataSource
		if alertText == "" && len(filteredDataSource) == 0 {
			alertText = "No match found for : " + commandStr
			if len(selectedContext) > 0 {
				alertText += " in " + strings.Join(selectedContext, " \u2b95 ")
			}
		}
		debug(filterText, debugBox)

		displayItems := []ui.Drawable{contextListBox, contextTextBox, debugBox}
		if alertText != "" {
			alertBox.Text = alertText
			displayItems = append(displayItems, alertBox)
		}

		ui.Render(displayItems...)
		alertText = ""
	}
}
