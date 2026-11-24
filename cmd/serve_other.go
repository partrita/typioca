//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !solaris

package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

// sessionProgramOptions connects the TUI to the session channel, the portable
// option for platforms without a POSIX PTY.
func sessionProgramOptions(session ssh.Session) []tea.ProgramOption {
	return []tea.ProgramOption{
		tea.WithInput(session),
		tea.WithOutput(session),
	}
}
