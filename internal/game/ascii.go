package game

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

//go:embed ascii/*.txt
var asciiFS embed.FS

// getAscii renvoie le contenu de <name>.txt.
// On lit d'abord le fichier SUR LE DISQUE (pratique : une retouche d'un dessin
// est visible sans recompiler), puis on retombe sur la copie embarquée dans le
// binaire (utile si le jeu est distribué seul). Chaîne vide si rien n'est trouvé.
func getAscii(name string) string {
	for _, p := range []string{
		filepath.Join("internal", "game", "ascii", name+".txt"), // lancé depuis la racine
		filepath.Join("ascii", name+".txt"),                     // lancé depuis internal/game
	} {
		if data, err := os.ReadFile(p); err == nil {
			return strings.TrimRight(string(data), "\n")
		}
	}
	if data, err := asciiFS.ReadFile("ascii/" + name + ".txt"); err == nil {
		return strings.TrimRight(string(data), "\n")
	}
	return ""
}
