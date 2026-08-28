package game

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// ClearScreen efface le terminal (écran + historique de défilement) et remet le
// curseur en haut, pour un rendu qui se rafraîchit sur place au lieu de défiler.
//
//	\033[3J : efface le buffer de défilement (évite l'empilement des écrans)
//	\033[H  : curseur en haut à gauche
//	\033[2J : efface la zone visible
func ClearScreen() {
	fmt.Print("\033[3J\033[H\033[2J")
}

// vitesseTexte est le délai entre deux caractères de TypeText.
// Mettre la variable d'environnement RED_FAST=1 la désactive (tests).
var vitesseTexte = 15 * time.Millisecond

// TypeText affiche du texte avec un effet d'animation (machine à écrire).
func TypeText(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		if vitesseTexte > 0 {
			time.Sleep(vitesseTexte)
		}
	}
	fmt.Println()
}

// ReadKey lit une touche pressée sans nécessiter la touche Entrée.
// Si l'entrée n'est pas un vrai terminal (redirigée), on retombe sur une lecture
// ligne par ligne : le jeu reste alors jouable/scriptable (une lettre + Entrée).
func ReadKey() string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		ligne, e := reader.ReadString('\n')
		if e != nil {
			soundQuit()
			fmt.Println("\nEntrée fermée, fin du programme.")
			os.Exit(0)
		}
		return strings.TrimSpace(ligne)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 3)
	os.Stdin.Read(b)

	// Ctrl+C
	if b[0] == 3 {
		term.Restore(int(os.Stdin.Fd()), oldState)
		soundQuit()
		os.Exit(0)
	}
	// Flèches directionnelles
	if b[0] == 27 && b[1] == '[' {
		switch b[2] {
		case 'A':
			return "Up"
		case 'B':
			return "Down"
		case 'C':
			return "Right"
		case 'D':
			return "Left"
		}
	}
	// Autres touches
	return string(b[0:1])
}
