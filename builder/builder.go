package builder

import (
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Builder hold the build retails
type Builder struct {
	app         string
	dir         string
	latestBuild time.Time
}

// New returns new builder object
func New(name, path string) *Builder {
	return &Builder{
		app: name,
		dir: path,
	}
}

// GetLatestBuild returns last build time
func (b *Builder) GetLatestBuild() time.Time {
	return b.latestBuild
}

// SetLatestBuild update lastest build time
func (b *Builder) SetLatestBuild(time time.Time) {
	b.latestBuild = time
}

// Build builds our application
func (b *Builder) Build() bool {

	commands := []string{"go", "build", "-o", b.app}

	log.Infof("Build: %s", strings.Join(commands, " "))

	// Execute build commands
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Dir = b.dir
	out, err := cmd.CombinedOutput()
	// Update last build time
	b.SetLatestBuild(time.Now())
	if err != nil {
		log.Errorf("Build: %s : %s\n", err.Error(), string(out))
		return false
	}

	return true
}
