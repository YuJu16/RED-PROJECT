package game

import (
	"fmt"
	"strings"
)

// Couleurs ANSI. Désactivées si la variable d'environnement NO_COLOR est présente.
var couleurActive = true

const (
	cReset = "\033[0m"
	cGras  = "\033[1m"
	cRouge = "\033[31m"
	cVert  = "\033[32m"
	cJaune = "\033[33m"
	cCyan  = "\033[36m"
)

// col colore s si les couleurs sont actives.
func col(s, c string) string {
	if !couleurActive {
		return s
	}
	return c + s + cReset
}

// titre encadre un intitulé de section.
func titre(s string) string {
	return col(cGras+"=== "+s+" ===", cCyan)
}

// afficherAscii imprime le dessin ascii/<name>.txt. Si le fichier est absent ou
// vide, un cadre "dessin à venir" est affiché à la place (aucun plantage).
func afficherAscii(name string) {
	art := getAscii(name)
	if strings.TrimSpace(art) != "" {
		fmt.Println(art)
		return
	}
	titre := " " + strings.ToUpper(name) + " "
	largeur := 44
	pad := (largeur - len(titre)) / 2
	fmt.Println("+" + strings.Repeat("-", largeur) + "+")
	fmt.Println("|" + strings.Repeat(" ", largeur) + "|")
	fmt.Println("|" + strings.Repeat(" ", pad) + titre + strings.Repeat(" ", largeur-pad-len(titre)) + "|")
	fmt.Println("|" + center("(dessin ascii/"+name+".txt à venir)", largeur) + "|")
	fmt.Println("|" + strings.Repeat(" ", largeur) + "|")
	fmt.Println("+" + strings.Repeat("-", largeur) + "+")
}

func center(s string, largeur int) string {
	if len(s) >= largeur {
		return s[:largeur]
	}
	pad := (largeur - len(s)) / 2
	return strings.Repeat(" ", pad) + s + strings.Repeat(" ", largeur-pad-len(s))
}

// barre rend une jauge texte : [#####-----] cur/max.
func barre(cur, max, largeur int) string {
	if max <= 0 {
		max = 1
	}
	if cur < 0 {
		cur = 0
	}
	rempli := cur * largeur / max
	if rempli > largeur {
		rempli = largeur
	}
	jauge := strings.Repeat("#", rempli) + strings.Repeat("-", largeur-rempli)
	ratio := float64(cur) / float64(max)
	switch {
	case ratio > 0.5:
		jauge = col(jauge, cVert)
	case ratio > 0.25:
		jauge = col(jauge, cJaune)
	default:
		jauge = col(jauge, cRouge)
	}
	return "[" + jauge + "]"
}

// barresPV affiche les jauges de PV du joueur et du monstre au début d'un tour.
func barresPV(p *Personnage, m *Monster) {
	fmt.Printf("  %-16s Nv.%-3d %s %d/%d\n", m.Nom, m.Niveau, barre(m.PV, m.PVMax, 20), m.PV, m.PVMax)
	fmt.Printf("  %-16s Nv.%-3d %s %d/%d   Mana %d/%d\n", p.Starter, p.Niveau, barre(p.PV, p.PVMax, 20), p.PV, p.PVMax, p.Mana, p.ManaMax)
}
