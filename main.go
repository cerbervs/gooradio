package main

import (
	"fmt"
	"gooradio/views/startscreen"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// TODO: Implement
	m, err := startscreen.NewStartScreen()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}
