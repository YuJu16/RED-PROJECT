package main

import (
	"embed"
	"strings"
)

//go:embed ascii/*.txt
var asciiFS embed.FS

// getAscii renvoie le contenu du fichier ascii/<name>.txt (chaîne vide si absent).
func getAscii(name string) string {
	data, err := asciiFS.ReadFile("ascii/" + name + ".txt")
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(data), "\n")
}
