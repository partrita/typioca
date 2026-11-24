//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
)

// sessionProgramOptions connects the TUI to the session. When the server
// allocated a real PTY the program talks to the slave side, otherwise it reads
// and writes the session channel directly.
func sessionProgramOptions(session ssh.Session) []tea.ProgramOption {
	pty, _, ok := session.Pty()
	if !ok || session.EmulatedPty() {
		return []tea.ProgramOption{
			tea.WithInput(session),
			tea.WithOutput(session),
		}
	}

	return []tea.ProgramOption{
		tea.WithInput(pty.Slave),
		tea.WithOutput(pty.Slave),
	}
}
