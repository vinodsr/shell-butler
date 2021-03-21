package add

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/vinodsr/shell-butler/lib/config"
	runtime "github.com/vinodsr/shell-butler/lib/runtime"
	types "github.com/vinodsr/shell-butler/lib/types"
)

//Execute the operations for adding a comamnd
func Execute() {
	var contextInput, commandInput string
	rt := runtime.GetRunTime()
	commandMap := rt.GetCommandMap()
	configData := rt.GetConfig()
	commandAccepted := false
	for commandAccepted == false {
		color.Info.Println("Input the command context seperated by :")
		fmt.Println("  eg : (server:nginx:start)")
		fmt.Print("> ")
		fmt.Scan(&contextInput)
		if _, found := commandMap[contextInput]; found {
			fmt.Println()
			color.Error.Println("Duplicate context !")
			fmt.Println()
		} else {
			commandAccepted = true

		}
	}

	fmt.Printf("Input the command you need to run for %s ", contextInput)
	fmt.Println()
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	commandInput = scanner.Text()

	fmt.Println()
	fmt.Println()
	color.LightCyan.Println("Context : " + contextInput)
	fmt.Println()
	color.LightCyan.Println("Command : " + commandInput)
	fmt.Println()
	fmt.Println()
	color.Info.Print(" > Are you sure to add this [Y/n] ? ")
	var proceed = "y"
	fmt.Scan(&proceed)
	if strings.ToLower(proceed) != "y" {
		commandAccepted = false
	}

	if commandAccepted {
		newCommand := types.Command{
			Context: contextInput + ":",
			Program: commandInput,
		}
		configData.Commands = append(configData.Commands, newCommand)
		config.WriteConfig(configData)
		fmt.Println()
		fmt.Println()
		color.Info.Tips("Added new command")
	} else {
		color.Magenta.Println("No changes have been made")
	}

}
