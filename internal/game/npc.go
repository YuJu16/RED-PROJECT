package game

import (
	"fmt"
	"strings"
)

// attendreEntree met le jeu en pause jusqu'à ce que le joueur appuie sur Entrée.
func attendreEntree() {
	lireLigne(col("   ▶ [Entrée]", cCyan))
}

// jouerDialogue affiche les répliques une par une, façon boîte de dialogue Pokémon :
// on appuie sur Entrée après CHAQUE réplique. Les lignes vides sont ignorées.
// Si une ligne contient déjà " : " (autre interlocuteur), elle est affichée telle quelle.
func jouerDialogue(nom string, lignes []string) {
	for _, l := range lignes {
		if strings.TrimSpace(l) == "" {
			continue
		}
		switch {
		case nom == "" || strings.Contains(l, " : "):
			TypeText(l)
		default:
			TypeText(nom + " : " + l)
		}
		attendreEntree()
	}
}

// combatDresseur enchaîne les Pokémon de l'équipe d'un dresseur (pas de capture).
// Renvoie true si le joueur les a tous vaincus.
func combatDresseur(p *Personnage, n *PNJ) bool {
	musiqueCombat = "battle_trainer"
	if n.Theme == "plasma" {
		musiqueCombat = "plasma_battle" // repli auto sur battle_trainer si le fichier manque
	}
	musique(musiqueCombat)
	ClearScreen()
	TypeText(fmt.Sprintf("%s veut se battre !", n.Nom))
	TypeText(fmt.Sprintf("Équipe adverse : %s Pokémon.", nombreEnLettres(len(n.Equipe))))
	attendreEntree()

	for i, espece := range n.Equipe {
		m := nouveauAuNiveau(espece, p.Niveau) // dresseur : Pokémon au niveau du joueur
		ClearScreen()
		TypeText(fmt.Sprintf("%s envoie %s ! (%d / %d)", n.Nom, espece, i+1, len(n.Equipe)))
		attendreEntree()
		if !deroulerCombat(p, m, false) {
			fmt.Printf("\n%s a gagné le duel...\n", n.Nom)
			dead(p)
			return false
		}
		distribuerXP(p, m.Experience)
		if i < len(n.Equipe)-1 {
			attendreEntree()
		}
	}
	jingle("victory")
	return true
}

// nombreEnLettres : petit confort d'affichage pour 1..6.
func nombreEnLettres(n int) string {
	noms := []string{"zéro", "un", "deux", "trois", "quatre", "cinq", "six"}
	if n >= 0 && n < len(noms) {
		return noms[n]
	}
	return fmt.Sprintf("%d", n)
}
