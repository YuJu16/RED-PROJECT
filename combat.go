package main

import (
	"fmt"
	"time"
)

// trainingFight lance un combat d'entrainement tour par tour contre un monstre sauvage au choix.
func trainingFight(p *Personnage) {
	ClearScreen()
	fmt.Println("\n=== Zone d'entraînement ===")
	fmt.Println("1. Affronter un Ratentif sauvage")
	fmt.Println("2. Affronter un Zorua sauvage")
	fmt.Println("3. Affronter un Golette sauvage")
	fmt.Println("R. Retour au menu")

	choix := lireLigne("> ")
	var m *Monster
	switch choix {
	case "1":
		m = InitRatentif()
	case "2":
		m = InitZorua()
	case "3":
		m = InitGolette()
	default:
		return
	}

	ClearScreen()
	fmt.Println(getAscii(m.Ascii))
	TypeText(fmt.Sprintf("Un %s sauvage apparaît !", m.Nom))

	joueurCommence := p.Initiative >= m.Initiative
	if joueurCommence {
		fmt.Printf("%s (initiative %d) est plus rapide et commence le combat !\n", p.Nom, p.Initiative)
	} else {
		fmt.Printf("%s (initiative %d) est plus rapide et commence le combat !\n", m.Nom, m.Initiative)
	}

	tour := 1
	for {
		fmt.Printf("\n--- Tour %d ---\n", tour)

		if joueurCommence {
			charTurn(p, m)
			if m.PV <= 0 {
				break
			}
			monsterPattern(m, p, tour)
			if p.PV <= 0 {
				break
			}
		} else {
			monsterPattern(m, p, tour)
			if p.PV <= 0 {
				break
			}
			charTurn(p, m)
			if m.PV <= 0 {
				break
			}
		}

		tour++
	}

	if p.PV <= 0 {
		fmt.Printf("\n%s est tombé au combat !\n", p.Nom)
		dead(p)
	} else {
		orGagne := 5 + (m.Experience / 2)
		fmt.Printf("\n%s est vaincu ! %s remporte le combat et trouve %d pièces d'or.\n", m.Nom, p.Nom, orGagne)
		p.Or += orGagne

		// Drop de ressource aléatoire
		ressources := []string{"Plume de Poichigeon", "Peau de Grotichon", "Cuir de Roitiflam", "Plume de Déflaisan"}
		drop := ressources[time.Now().UnixNano()%int64(len(ressources))]
		fmt.Printf("Vous trouvez un objet : %s !\n", drop)
		addInventory(p, Objet{Nom: drop, Quantite: 1, Type: "Ressource"})

		gainExperience(p, m.Experience)
	}
}

// charTurn simule le tour de jeu du joueur : Attaquer, Sorts ou Inventaire.
func charTurn(p *Personnage, m *Monster) {
	fmt.Printf("\n%s, c'est ton tour !\n", p.Nom)
	fmt.Println("1. Attaquer")
	fmt.Println("2. Sorts")
	fmt.Println("3. Inventaire")
	fmt.Println("4. Capturer")

	choix := ""
	for choix != "1" && choix != "2" && choix != "3" && choix != "4" {
		choix = lireLigne("> ")
		switch choix {
		case "1":
			degats := 5
			m.PV -= degats
			if m.PV < 0 {
				m.PV = 0
			}
			TypeText(fmt.Sprintf("%s utilise Attaque Rapide et inflige %d dégâts à %s !", p.Starter, degats, m.Nom))
			fmt.Printf("%s - PV : %d/%d\n", m.Nom, m.PV, m.PVMax)
		case "2":
			castSpellInCombat(p, m)
		case "3":
			useItemInCombat(p)
		case "4":
			if possede(p, "Pokéball", 1) {
				removeInventory(p, "Pokéball", 1)
				TypeText(fmt.Sprintf("%s lance une Pokéball !", p.Nom))
				// Chance de capture basée sur les PV restants
				chance := (float64(m.PVMax-m.PV) / float64(m.PVMax)) * 100.0
				if chance > 40.0 { // 40% de dégâts infligés pour avoir une chance
					TypeText(fmt.Sprintf("1... 2... 3... et hop ! %s est capturé !", m.Nom))
					addInventory(p, Objet{Nom: m.Nom, Quantite: 1, Type: "Pokemon"})
					m.PV = 0 // Fin du combat
				} else {
					TypeText("Oh non ! Le Pokémon s'est libéré !")
				}
			} else {
				fmt.Println("Tu n'as pas de Pokéball !")
				choix = "" // On annule pour re-choisir
			}
		default:
			fmt.Println("Choix invalide.")
			choix = ""
		}
	}
}

// useItemInCombat affiche l'inventaire pendant le combat et applique l'objet choisi.
func useItemInCombat(p *Personnage) {
	if len(p.Inventaire) == 0 {
		fmt.Println("Ton inventaire est vide.")
		return
	}
	fmt.Println("=== Inventaire ===")
	for i, o := range p.Inventaire {
		fmt.Printf("%d. %s x%d\n", i+1, o.Nom, o.Quantite)
	}
	fmt.Println("R. Retour")

	choix := lireLigne("> ")
	if choix == "R" || choix == "r" {
		return
	}
	index := indexFromChoice(choix, len(p.Inventaire))
	if index == -1 {
		fmt.Println("Choix invalide.")
		return
	}
	useItem(p, index)
}

// monsterPattern fait attaquer le monstre m : 100% de son attaque, 200% tous les 3 tours.
func monsterPattern(m *Monster, p *Personnage, tour int) {
	degats := m.Attaque
	pourcentage := "100%"
	if tour%3 == 0 {
		degats = m.Attaque * 2
		pourcentage = "200%"
	}
	p.PV -= degats
	if p.PV < 0 {
		p.PV = 0
	}
	fmt.Printf("%s inflige à %s %d de dégâts (%s de son attaque) !\n", m.Nom, p.Nom, degats, pourcentage)
	fmt.Printf("%s - PV : %d/%d\n", p.Nom, p.PV, p.PVMax)
}
