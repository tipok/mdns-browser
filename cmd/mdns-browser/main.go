package main

import (
	"context"
	"fmt"
	"log/slog"
	"mdns-browser/internal/data"
	"mdns-browser/internal/discovery"
	"mdns-browser/internal/tui"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	addCh := make(chan data.ListItem, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		slog.Info("received termination signal, shutting down gracefully")
		cancel()
	}()

	go func() {
		err := discovery.ListAllServices(ctx, addCh)
		if err != nil {
			slog.Error("error discovering services", "error", err)
			os.Exit(1)
		}
	}()

	m := tui.Tui(tui.ListOpts{
		Title: "Found Services",
		AddCh: addCh,
	})

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
