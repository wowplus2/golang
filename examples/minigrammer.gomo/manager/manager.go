package manager

import (
	"runtime"
	"os"
	"bytes"
	"fmt"
	"path/filepath"
)


// Manager manage the memo file and command list.
type Manager struct {
	file	string
	cmds	map[string]Command
}

// Command has necessary information for each command.
// `Run` is function to execute the specific command
type Command struct {
	Name	string
	Usage	string
	Run	func(string, []string) error
}

// New create default Manager instance.
func New() *Manager {
	homeDir := GetHomeDir()

	return &Manager{
		file:	filepath.Join(homeDir, ".gomo", "memo.json"),
		cmds:	make(map[string]Command),
	}
}

// Usage creates usage message string of all available commands.
func (m *Manager) Usage() string {
	buf := bytes.NewBufferString("\n")

	for _, c := range m.cmds {
		fmt.Fprintf(buf, c.Usage)
	}

	return buf.String()
}

// AddCommand registers the specific command to command list of Manager.
func (m *Manager) AddCommand(cmd Command) {
	m.cmds[cmd.Name] = cmd
}

// Execute parses the command line arguments and
// runs the `Run` function of command with that parsed arguments.
func (m *Manager) Execute(args []string) error {
	var cmdName string
	var cmdArgs []string

	if len(args) > 1 {
		cmdArgs = args[1:]
	}

	cmdName = args[0]
	cmd, ok := m.cmds[cmdName]
	if !ok {
		return fmt.Errorf("%s is not defined...", cmdName)
	}

	if err := cmd.Run(m.file, cmdArgs); err != nil {
		return err
	}

	return nil
}

// GetHomeDir gets gome directory corresponding to each OS.
func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		return home
	}

	return os.Getenv("HOME")
}