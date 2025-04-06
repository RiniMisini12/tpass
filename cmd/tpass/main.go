package tpass

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rinimisini112/tpass/internal/store"
	"github.com/rinimisini112/tpass/internal/tui"
	"log"
	"os"
)

func Run() {
	masterPassword := store.GetMasterPassword()

	passwordStore, err := store.LoadStore(masterPassword)
	if err != nil {
		log.Fatalf("Error loading password store: %v", err)
	}

	p := tea.NewProgram(tui.InitialModel(passwordStore), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}
}
