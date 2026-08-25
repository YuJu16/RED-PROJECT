package main

import "fmt"

// Sort décrit un sort offensif utilisable en combat.
type Sort struct {
	Nom      string
	Degats   int
	CoutMana int
}

// sortsDisponibles associe chaque nom de sort connu à ses caractéristiques.
var sortsDisponibles = map[string]Sort{
	"Charge":        {Nom: "Charge", Degats: 8, CoutMana: 4},
	"Lance-Flammes": {Nom: "Lance-Flammes", Degats: 18, CoutMana: 10},
}

// spellBook ajoute le sort "Boule de Feu" à la liste de sorts du personnage.
// Un même sort ne peut être appris qu'une seule fois.
func spellBook(p *Personnage) {
	for _, s := range p.Skills {
		if s == "Lance-Flammes" {
			fmt.Println("Ce sort est déjà appris : tu connais déjà Lance-Flammes !")
			return
		}
	}
	p.Skills = append(p.Skills, "Lance-Flammes")
	fmt.Println("Tu as appris l'attaque Lance-Flammes !")
}

// castSpellInCombat propose au joueur de lancer un des sorts connus contre le monstre m.
func castSpellInCombat(p *Personnage, m *Monster) {
	if len(p.Skills) == 0 {
		fmt.Println("Tu ne connais aucun sort.")
		return
	}
	fmt.Println("=== Sorts connus ===")
	for i, s := range p.Skills {
		sort := sortsDisponibles[s]
		fmt.Printf("%d. %s (%d dégâts, %d mana)\n", i+1, sort.Nom, sort.Degats, sort.CoutMana)
	}
	fmt.Println("R. Retour")

	choix := lireLigne("> ")
	if choix == "R" || choix == "r" {
		return
	}
	index := indexFromChoice(choix, len(p.Skills))
	if index == -1 {
		fmt.Println("Choix invalide.")
		return
	}

	sort := sortsDisponibles[p.Skills[index]]
	if p.Mana < sort.CoutMana {
		fmt.Println("Mana insuffisant pour lancer " + sort.Nom + " !")
		return
	}

	p.Mana -= sort.CoutMana
	m.PV -= sort.Degats
	if m.PV < 0 {
		m.PV = 0
	}
	fmt.Printf("%s lance %s et inflige %d dégâts à %s !\n", p.Starter, sort.Nom, sort.Degats, m.Nom)
	fmt.Printf("%s - PV : %d/%d\n", m.Nom, m.PV, m.PVMax)
	fmt.Printf("Mana restant : %d/%d\n", p.Mana, p.ManaMax)
}
