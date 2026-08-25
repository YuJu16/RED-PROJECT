package main

import (
	"fmt"
	"time"
)

// ClearScreen efface le terminal pour rendre l'affichage plus agréable.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// TypeText affiche du texte avec un effet d'animation (machine à écrire).
func TypeText(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(15 * time.Millisecond)
	}
	fmt.Println()
}
