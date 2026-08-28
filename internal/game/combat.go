package game

import (
	"fmt"
	"math/rand"
)

// ---------------------------------------------------------------------------
// Boucle de combat tour par tour (Tâches 17 à 20 + Missions 3, 5).
// ---------------------------------------------------------------------------

var (
	combatCapture bool   // le combat s'est terminé par une capture
	combatFuite   bool   // le joueur a pris la fuite
	musiqueCombat string // morceau de combat courant (pour revenir après "low_hp")
)

// verifierLowHP bascule sur la musique de PV bas quand le Pokémon actif est
// sous 25 % de ses PV, et revient à la musique de combat quand il remonte.
func verifierLowHP(p *Personnage) {
	bas := p.PV > 0 && p.PV*4 <= p.PVMax
	switch {
	case bas && musiqueCourante != "low_hp":
		musique("low_hp")
	case !bas && musiqueCourante == "low_hp" && musiqueCombat != "":
		musique(musiqueCombat)
	}
}

// deroulerCombat enchaîne les tours entre le joueur et le monstre m.
// L'initiative décide qui commence (Mission 3).
// Renvoie true si le joueur a gagné (monstre K.O., joueur debout, sans fuite).
func deroulerCombat(p *Personnage, m *Monster, capturable bool) bool {
	combatCapture = false
	combatFuite = false
	ClearScreen()
	afficherAscii(m.Ascii)
	TypeText(fmt.Sprintf("%s (Nv. %d) entre en scène ! (PV %d - ATQ %d)", m.Nom, m.Niveau, m.PVMax, m.Attaque))

	joueurCommence := p.Initiative >= m.Initiative
	if joueurCommence {
		fmt.Printf("%s (initiative %d) est le plus rapide et commence le combat !\n", p.Nom, p.Initiative)
	} else {
		fmt.Printf("%s (initiative %d) est le plus rapide et commence le combat !\n", m.Nom, m.Initiative)
	}

	tour := 1
	for {
		fmt.Println("\n" + col(fmt.Sprintf("--- Tour %d ---", tour), cCyan))
		barresPV(p, m)
		verifierLowHP(p) // bascule la musique dès l'affichage du tour (même au tour 1)

		// tourJoueur / tourMonstre renvoient false si le combat doit s'arrêter.
		tourJoueur := func() bool {
			charTurn(p, m, capturable)
			if combatFuite {
				return false
			}
			return m.PV > 0
		}
		tourMonstre := func() bool {
			monsterPattern(m, p, tour)
			if p.PV <= 0 {
				return gererKO(p, m.Type) // true = un autre Pokémon prend le relais
			}
			return true
		}

		if joueurCommence {
			if !tourJoueur() || !tourMonstre() {
				break
			}
		} else {
			if !tourMonstre() || !tourJoueur() {
				break
			}
		}
		verifierLowHP(p)
		tour++
	}

	syncPersoVersActif(p)
	return p.PV > 0 && !combatFuite
}

// combatSauvage : rencontre dans les hautes herbes. Capture et fuite possibles.
func combatSauvage(p *Personnage, m *Monster) {
	musiqueCombat = "battle_wild"
	musique("battle_wild")
	TypeText(fmt.Sprintf("Un %s sauvage (Nv. %d) apparaît !", m.Nom, m.Niveau))
	if flagOK(p, "capture:"+m.Nom) {
		TypeText(col("(Pokédex : "+m.Nom+" est déjà enregistré.)", cCyan))
	} else {
		TypeText(col("(Pokédex : "+m.Nom+" est une nouvelle espèce !)", cJaune))
	}
	gagne := deroulerCombat(p, m, true)
	if combatFuite {
		fmt.Println("\nVous avez pris la fuite.")
		return
	}
	if gagne {
		jingle("victory")
		recompenseSauvage(p, m)
	} else {
		dead(p)
	}
}

// recompenseSauvage : or, ressource et expérience après une victoire sauvage.
func recompenseSauvage(p *Personnage, m *Monster) {
	if combatCapture {
		distribuerXP(p, m.Experience/2)
		return
	}
	orGagne := 5 + m.Experience/2
	p.Or += orGagne
	fmt.Printf("\n%s est vaincu ! %s trouve %d Poké-Dollars.\n", m.Nom, p.Nom, orGagne)

	ressources := []string{"Plume de Poichigeon", "Peau de Grotichon", "Cuir de Roitiflam", "Plume de Déflaisan"}
	drop := ressources[rand.Intn(len(ressources))]
	if err := addInventory(p, Objet{Nom: drop, Quantite: 1, Type: "Ressource"}); err == nil {
		fmt.Printf("Vous ramassez : %s !\n", drop)
	}
	distribuerXP(p, m.Experience)
}

// trainingFight : combat d'entraînement du sujet (Tâche 17), accessible au menu.
func trainingFight(p *Personnage) {
	ClearScreen()
	fmt.Println("\n" + titre("Zone d'entraînement"))
	fmt.Println("1. Affronter un Ratentif sauvage")
	fmt.Println("2. Affronter un Zorua sauvage")
	fmt.Println("3. Affronter un Golette sauvage")
	fmt.Println("R. Retour au menu")

	var m *Monster
	switch lireLigne("> ") {
	case "1":
		m = InitRatentif()
	case "2":
		m = InitZorua()
	case "3":
		m = InitGolette()
	default:
		return
	}
	combatSauvage(p, m)
}

// charTurn simule le tour du joueur. Le menu s'adapte au contexte (équipe, sauvage).
func charTurn(p *Personnage, m *Monster, capturable bool) {
	for {
		fmt.Printf("\n%s, à toi de jouer (%s) !\n", p.Nom, p.Starter)
		opts := []string{"Attaquer (choisir l'attaque)", "Sac (inventaire)"}
		if len(p.Equipe) > 1 {
			opts = append(opts, "Changer de Pokémon")
		}
		if capturable {
			opts = append(opts, "Capturer (Pokéball)", "Fuir")
		}
		for i, o := range opts {
			fmt.Printf("%d. %s\n", i+1, o)
		}

		choix := indexFromChoice(lireLigne("> "), len(opts))
		if choix == -1 {
			fmt.Println("Choix invalide.")
			continue
		}
		switch opts[choix] {
		case "Attaquer (choisir l'attaque)":
			if choisirAttaque(p, m) {
				return // une attaque a été lancée -> tour consommé
			}
		case "Sac (inventaire)":
			if useItemInCombat(p) {
				return // un objet a été utilisé -> tour consommé, le monstre joue ensuite
			}
		case "Changer de Pokémon":
			if changerPokemon(p, m.Type) {
				return // le changement consomme le tour
			}
		case "Capturer (Pokéball)":
			tenterCapture(p, m)
			return
		case "Fuir":
			if rand.Intn(100) < 70 {
				TypeText("Vous prenez vos jambes à votre cou... Fuite réussie !")
				combatFuite = true
			} else {
				TypeText("Impossible de fuir !")
			}
			return
		}
	}
}

// tenterCapture : lance une Pokéball. La chance dépend des PV entamés de l'adversaire.
func tenterCapture(p *Personnage, m *Monster) {
	if !possede(p, "Pokéball", 1) {
		fmt.Println("Tu n'as pas de Pokéball ! (tour perdu)")
		return
	}
	removeInventory(p, "Pokéball", 1)
	ClearScreen()
	afficherAscii("pokeball")
	TypeText(fmt.Sprintf("%s lance une Pokéball sur %s !", p.Nom, m.Nom))
	chance := float64(m.PVMax-m.PV) / float64(m.PVMax) * 100.0
	if chance <= 40.0 {
		TypeText("Oh non ! Le Pokémon s'est libéré ! (il faut l'affaiblir davantage)")
		return
	}
	TypeText(fmt.Sprintf("1... 2... 3... Clic ! %s est capturé !", m.Nom))
	m.PV = 0
	combatCapture = true
	flagSet(p, "capture:"+m.Nom) // enregistré au Pokédex
	if ajouterAEquipe(p, m, p.Niveau) {
		TypeText(fmt.Sprintf("%s rejoint ton équipe !", m.Nom))
	} else {
		addInventory(p, Objet{Nom: m.Nom, Quantite: 1, Type: "Pokemon"})
		TypeText(fmt.Sprintf("Ton équipe est pleine : %s est envoyé au Centre Pokémon.", m.Nom))
	}
}

// useItemInCombat affiche l'inventaire pendant le combat et applique l'objet choisi.
// Renvoie true si un objet a réellement été utilisé (le tour est alors consommé).
func useItemInCombat(p *Personnage) bool {
	if len(p.Inventaire) == 0 {
		fmt.Println("Ton sac est vide.")
		return false
	}
	fmt.Println("=== Sac ===")
	for i, o := range p.Inventaire {
		fmt.Printf("%d. %s x%d\n", i+1, o.Nom, o.Quantite)
	}
	fmt.Println("R. Retour")

	choix := lireLigne("> ")
	if choix == "R" || choix == "r" {
		return false
	}
	index := indexFromChoice(choix, len(p.Inventaire))
	if index == -1 {
		fmt.Println("Choix invalide.")
		return false
	}
	return useItem(p, index) // true seulement si un objet a été réellement utilisé
}

// monsterPattern : le monstre inflige 100 % de son attaque, 200 % tous les 3 tours
// (Tâche 18). La table des types module ensuite les dégâts (Bonus 1).
func monsterPattern(m *Monster, p *Personnage, tour int) {
	degats := m.Attaque
	pourcentage := "100%"
	if tour%3 == 0 {
		degats = m.Attaque * 2
		pourcentage = "200%"
	}
	degats, label := appliquerType(degats, m.Type, p.Type)
	p.PV -= degats
	if p.PV < 0 {
		p.PV = 0
	}
	fmt.Printf("%s inflige à %s %d de dégâts (%s de son attaque) !\n", m.Nom, p.Nom, degats, pourcentage)
	if label != "" {
		fmt.Println(label)
	}
	fmt.Printf("%s - PV : %d/%d\n", p.Nom, p.PV, p.PVMax)
}
