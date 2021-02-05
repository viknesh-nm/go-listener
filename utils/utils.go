package utils

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// GracefulShutdown shutdown system cleanly if any interrupts happen.
func GracefulShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for sig := range signalChan {
		if sig == syscall.SIGINT {
			log.Warn("Exiting...")
			os.Exit(0)
		}
	}
}

// GetDefaultPath gets default file path.
func GetDefaultPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Errorf("An error occurred while getting the current working directory: %v\n", err)
		os.Exit(1)
	}
	path, err := filepath.Abs(dir)
	if err != nil {
		log.Errorf("An error occurred while finding an absolute working path: %v\n", err)
		os.Exit(1)
	}

	return path
}

// ValidateCommands split commands and remove unwanted space between these.
func ValidateCommands(c string) []string {

	cmd := strings.Split(strings.TrimSpace(c), " ")
	final := make([]string, 0, len(cmd))
	keys := make(map[string]bool, 0)

	// Eliminate duplicate commands and empty spaces.
	for _, v := range cmd {
		if got := keys[v]; !got && v != "" {
			keys[v] = true
			final = append(final, v)
		}
	}
	return final
}
