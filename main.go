package main

import (
	"os"

	addcli "github.com/vinodsr/shell-butler/lib/cli/add"
	"github.com/vinodsr/shell-butler/lib/cli/help"
	"github.com/vinodsr/shell-butler/lib/runtime"
	layout "github.com/vinodsr/shell-butler/lib/ui"
)

//Main program
func main() {
	// Initialize the runtime
	var rt *runtime.Runtime = runtime.GetRunTime()

	rt.Init()

	if len(os.Args) > 1 {
		argCommand := os.Args[1]
		if argCommand == "add" {
			addcli.Execute()
			os.Exit(1)
		} else if argCommand == "help" || argCommand == "--help" {
			help.Execute()
			os.Exit(1)
		}

	}

	// if len(rt.GetConfig().Commands) == 0 {
	// 	fmt.Println("No commands to load. Please add commands.")
	// 	os.Exit(1)
	// }
	//fmt.Printf("Data %+v\n", commands)
	//fmt.Print(commands.Commands[0].Context)

	// initialise the list datasource .

	layout.Render()
	//os.Exit(1)

}
