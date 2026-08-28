//go:build !windows

package game

import "os/exec"

// processusCache : rien à faire hors Windows (le son est déjà désactivé).
func processusCache(cmd *exec.Cmd) {}
