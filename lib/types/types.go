package lib

// Command structure
type Command struct {
	Context string
	Program string
}

// ConfigData structure
type ConfigData struct {
	Commands []Command
}
