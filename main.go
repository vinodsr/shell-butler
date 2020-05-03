package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/enescakir/emoji"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Command structure
type Command struct {
	Context string
	Program string
}

// ConfigData structure
type ConfigData struct {
	Commands []Command
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func updateContextList(contextList []string, contextFilter string) []string {
	filteredList := filter(contextList, func(str string) bool {
		// return strings.Contains(str, sel)

		match, _ := regexp.MatchString(contextFilter, str)
		return match
	})
	return filteredList
}

func displayContextatLevel(contextList []string, level int) []string {
	result := []string{}
	contextMap := make(map[string]bool)
	for _, s := range contextList {
		splits := strings.Split(s, ":")
		if len(splits) >= level {
			if contextMap[splits[level-1]] == false {
				result = append(result, splits[level-1])
				contextMap[splits[level-1]] = true
			}

		}
	}
	return result
}

func main() {

	var contextDataSource []string
	var filteredDataSource []string
	var commandLevel int = 1
	var selectedContext []string
	var commandMap = make(map[string]string)
	var commands ConfigData

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	dir, _ = os.UserHomeDir()
	configFile := dir + "/.butler/settings.json"
	dat, err := ioutil.ReadFile(configFile)
	if err != nil {
		initStruct := ConfigData{
			Commands: []Command{{
				Context: "Help",
				Program: "echo Edit " + configFile + " to add more commands",
			}},
		}
		//fmt.Printf("%+v\n", initStruct)
		//os.Exit(0)
		initStructJSON, _ := json.Marshal(initStruct)
		_, err := os.Stat(dir + "/.butler")

		if os.IsNotExist(err) {
			errDir := os.MkdirAll(dir+"/.butler", 0755)
			check(errDir)

		}
		err = ioutil.WriteFile(configFile, initStructJSON, 0644)
		check(err)
		commands = initStruct
	}

	json.Unmarshal([]byte(string(dat)), &commands)
	//fmt.Printf("Data %+v\n", commands)
	//fmt.Print(commands.Commands[0].Context)

	// initialise the list datasource .
	for _, command := range commands.Commands {
		// splits := strings.Split(command.Context, ":")
		contextDataSource = append(contextDataSource, command.Context)
		commandMap[command.Context] = command.Program
	}
	if len(os.Args) > 1 {
		argCommand := os.Args[1]
		if argCommand == "add" {
			fmt.Print("adding")
		}

	}
	//os.Exit(1)

	commandStr := ""

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

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
	debugBox.SetRect(0, 50, 50, 53)

	contextListBox := widgets.NewList()
	filteredDataSource = displayContextatLevel(updateContextList(contextDataSource, ""), commandLevel)
	contextListBox.Rows = filteredDataSource
	contextListBox.TextStyle = ui.NewStyle(ui.ColorWhite)
	contextListBox.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen, ui.ModifierBold)
	contextListBox.WrapText = false
	contextListBox.SetRect(0, 3, termWidth, termHeight)
	contextListBox.BorderStyle.Fg = ui.ColorGreen

	ui.Render(contextTextBox, contextListBox)

	previousKey := ""
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
		case "g":
			if previousKey == "g" {
				contextListBox.ScrollTop()
			}
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

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		//fmt.Println(emoji.Squid)
		//+ emoji.Squid.String()

		contextTextBox.Text = formatCommandString(selectedContext, commandStr)
		contextSep := ""
		if len(selectedContext) > 0 {
			contextSep = ":"
		}

		filterText := "(?i)^" + strings.Join(selectedContext, ":") + contextSep + ".*" + commandStr + ".*\\:?"
		filteredDataSource = displayContextatLevel(updateContextList(contextDataSource, filterText), commandLevel)
		contextListBox.Rows = filteredDataSource
		if alertText == "" && len(filteredDataSource) == 0 {
			alertText = "No match found for : " + commandStr
			if len(selectedContext) > 0 {
				alertText += " in " + strings.Join(selectedContext, " \u2b95 ")
			}
		}
		debug(filterText, debugBox)

		displayItems := []ui.Drawable{contextListBox, contextTextBox}
		if alertText != "" {
			alertBox.Text = alertText
			displayItems = append(displayItems, alertBox)
		}
		ui.Render(displayItems...)
		alertText = ""
	}
}

func debug(str string, d *widgets.Paragraph) {
	d.Text += str + "\n"
}

func formatCommandString(selectedContext []string, commandStr string) string {
	formattedCommandStr := " "
	for _, s := range selectedContext {
		formattedCommandStr += "[" + s + "] "
	}
	return formattedCommandStr + emoji.RightArrow.String() + " " + commandStr
}
