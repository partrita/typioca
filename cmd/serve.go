package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/muesli/termenv"
	cryptossh "golang.org/x/crypto/ssh"
)

// serveSession runs the typioca TUI for a single SSH session.
func serveSession(session ssh.Session) {
	pty, windowChanges, active := session.Pty()

	if !active {
		fmt.Fprintln(session.Stderr(), "not a tty")
		_ = session.Exit(1)
		_ = session.Close()
		return
	}

	log.Printf("session started: %s", sessionRemoteAddr(session))

	program := tea.NewProgram(
		initialModel(
			termenv.ANSI256,
			termenv.ANSIWhite,
			pty.Window.Width,
			pty.Window.Height,
		),
		append(sessionProgramOptions(session), tea.WithAltScreen())...,
	)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case window := <-windowChanges:
				program.Send(tea.WindowSizeMsg{Width: window.Width, Height: window.Height})
			}
		}
	}()

	if _, err := program.Run(); err != nil {
		log.Printf("session %s ended with error: %v", sessionRemoteAddr(session), err)
	}

	program.Kill()
	close(done)

	log.Printf("session ended: %s", sessionRemoteAddr(session))
}

func sessionRemoteAddr(session ssh.Session) string {
	if addr := session.RemoteAddr(); addr != nil {
		return addr.String()
	}
	return "unknown"
}

// ensureHostKey writes a fresh ed25519 host key to path when no key exists yet,
// so `typioca serve` works on a new machine without extra setup.
func ensureHostKey(path string) error {
	if path == "" {
		return fmt.Errorf("server key path must not be empty")
	}

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	block, err := cryptossh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0o600); err != nil {
		return err
	}

	sshPublicKey, err := cryptossh.NewPublicKey(publicKey)
	if err != nil {
		return err
	}

	return os.WriteFile(path+".pub", cryptossh.MarshalAuthorizedKey(sshPublicKey), 0o644)
}
