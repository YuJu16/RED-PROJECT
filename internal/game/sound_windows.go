//go:build windows

package game

import (
	"os/exec"
	"syscall"
)

// processusCache empêche PowerShell de créer une fenêtre console (donc pas de
// flash ni de vol de focus). CREATE_NO_WINDOW = 0x08000000.
func processusCache(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}
