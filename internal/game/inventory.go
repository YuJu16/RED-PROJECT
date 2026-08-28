package game

import (
	"fmt"
	"time"
)

// Objet représente un élément possédé par le personnage (potion, ressource, équipement...).
type Objet struct {
	Nom      string
	Quantite int
	Type     string // "Potion", "PotionPoison", "PotionMana", "LivreSort", "Ressource", "Equipement"
	Slot     string // "Tete", "Torse", "Pieds" pour un objet de Type "Equipement"
}

// accessInventory affiche l'inventaire et permet d'utiliser un objet, ou de revenir au menu.
func accessInventory(p *Personnage) {
	for {
		ClearScreen()
		fmt.Println("\n" + titre("Sac"))
		fmt.Printf("(%d/%d emplacements utilisés)\n\n", len(p.Inventaire), p.InventaireMax)
		if len(p.Inventaire) == 0 {
			fmt.Println("(vide)")
		}
		for i, o := range p.Inventaire {
			fmt.Printf("%2d. %-28s x%-2d  %s\n", i+1, o.Nom, o.Quantite, col(descriptionObjet(o), cJaune))
		}
		fmt.Println("R. Retour au menu")

		choix := lireLigne("Objet à utiliser (numéro) ou R pour revenir > ")
		if choix == "R" || choix == "r" {
			return
		}

		index := indexFromChoice(choix, len(p.Inventaire))
		if index == -1 {
			fmt.Println("Choix invalide.")
			lireLigne("\nAppuyez sur Entrée...")
			continue
		}
		useItem(p, index)
		lireLigne("\nAppuyez sur Entrée...")
	}
}

// descriptionObjet renvoie une courte explication selon le type d'objet.
func descriptionObjet(o Objet) string {
	switch o.Type {
	case "Potion":
		return "Rend 50 PV au Pokémon actif."
	case "PotionPoison":
		return "Poison : -10 PV/s pendant 3 s."
	case "PotionMana":
		return "Rend 30 points de Mana."
	case "LivreSort":
		return "Apprend un sort (une seule fois)."
	case "Equipement":
		return "À équiper (bonus de PV max)."
	case "Pokeball":
		return "Capture un Pokémon sauvage affaibli."
	case "Ressource":
		return "Matériau de forge."
	case "Pokemon":
		return "Pokémon en pension au Centre."
	default:
		return ""
	}
}

// totalQuantite : somme des quantités de tous les objets (pour détecter une conso).
func totalQuantite(p *Personnage) int {
	n := 0
	for _, o := range p.Inventaire {
		n += o.Quantite
	}
	return n
}

// useItem applique l'effet de l'objet situé à l'index i de l'inventaire.
// Renvoie true si un objet a effectivement été consommé/équipé (sinon annulé).
func useItem(p *Personnage, i int) bool {
	if i < 0 || i >= len(p.Inventaire) {
		fmt.Println("Choix invalide.")
		return false
	}
	avant := totalQuantite(p)
	item := p.Inventaire[i]
	switch item.Type {
	case "Potion":
		takePot(p)
	case "PotionPoison":
		poisonPot(p)
	case "PotionMana":
		takeManaPot(p)
	case "LivreSort":
		if spellBook(p) {
			removeInventory(p, item.Nom, 1)
		}
	case "Equipement":
		equip(p, item)
	default:
		fmt.Println("Cet objet ne peut pas être utilisé directement.")
	}
	return totalQuantite(p) != avant
}

// consumeOne retire une unité de l'objet situé à l'index i (le supprime si la quantité tombe à 0).
func consumeOne(p *Personnage, i int) {
	p.Inventaire[i].Quantite--
	if p.Inventaire[i].Quantite <= 0 {
		p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
	}
}

// choisirCibleSoin : à qui appliquer un soin. Renvoie l'index dans p.Equipe,
// ou -1 si annulé. Pas de question si l'équipe n'a qu'un Pokémon.
func choisirCibleSoin(p *Personnage) int {
	syncPersoVersActif(p)
	if len(p.Equipe) <= 1 {
		return 0
	}
	fmt.Println("Soigner quel Pokémon ?")
	for i, e := range p.Equipe {
		actif := ""
		if i == 0 {
			actif = " (actif)"
		}
		fmt.Printf("%d. %-12s %s %d/%d%s\n", i+1, e.Nom, barre(e.PV, e.PVMax, 12), e.PV, e.PVMax, actif)
	}
	fmt.Println("R. Annuler")
	idx := indexFromChoice(lireLigne("> "), len(p.Equipe))
	if idx == -1 {
		return -1
	}
	return idx
}

// takePot utilise une Potion de vie (soin de 50 PV, cap à PVMax) sur le Pokémon
// choisi par le joueur (actif par défaut si un seul Pokémon).
func takePot(p *Personnage) {
	pos := -1
	for i, o := range p.Inventaire {
		if o.Type == "Potion" {
			pos = i
			break
		}
	}
	if pos == -1 {
		fmt.Println("Erreur : tu n'as pas de Potion de vie !")
		return
	}

	cible := choisirCibleSoin(p)
	if cible == -1 {
		fmt.Println("Soin annulé.")
		return
	}

	soigné := ""
	if cible == 0 {
		if p.PV >= p.PVMax {
			fmt.Println(p.Starter + " a déjà tous ses PV !")
			return
		}
		p.PV += 50
		if p.PV > p.PVMax {
			p.PV = p.PVMax
		}
		syncPersoVersActif(p)
		soigné = fmt.Sprintf("%s : %d/%d PV", p.Starter, p.PV, p.PVMax)
	} else {
		e := &p.Equipe[cible]
		if e.PV >= e.PVMax {
			fmt.Println(e.Nom + " a déjà tous ses PV !")
			return
		}
		e.PV += 50
		if e.PV > e.PVMax {
			e.PV = e.PVMax
		}
		soigné = fmt.Sprintf("%s : %d/%d PV", e.Nom, e.PV, e.PVMax)
	}
	consumeOne(p, pos)
	fmt.Printf("Tu utilises une Potion de vie. %s\n", soigné)
}

// poisonPot inflige 10 PV de dégâts par seconde pendant 3 secondes.
func poisonPot(p *Personnage) {
	for i, o := range p.Inventaire {
		if o.Type == "PotionPoison" {
			consumeOne(p, i)
			fmt.Println("Tu bois la Potion de poison... quelle mauvaise idée !")
			for s := 0; s < 3; s++ {
				time.Sleep(1 * time.Second)
				p.PV -= 10
				if p.PV < 0 {
					p.PV = 0
				}
				fmt.Printf("PV : %d/%d\n", p.PV, p.PVMax)
			}
			dead(p)
			return
		}
	}
	fmt.Println("Erreur : tu n'as pas de Potion de poison !")
}

// takeManaPot utilise une Potion de mana (restaure 30 Mana, cap à ManaMax).
func takeManaPot(p *Personnage) {
	for i, o := range p.Inventaire {
		if o.Type == "PotionMana" {
			if p.Mana >= p.ManaMax {
				fmt.Println("Ton mana est déjà au maximum !")
				return
			}
			p.Mana += 30
			if p.Mana > p.ManaMax {
				p.Mana = p.ManaMax
			}
			consumeOne(p, i)
			fmt.Printf("Tu utilises une Potion de mana. Mana : %d/%d\n", p.Mana, p.ManaMax)
			return
		}
	}
	fmt.Println("Erreur : tu n'as pas de Potion de mana !")
}

// checkInventoryFull indique si l'inventaire a atteint sa capacité maximale.
func checkInventoryFull(p *Personnage) bool {
	return len(p.Inventaire) >= p.InventaireMax
}

// addInventory ajoute un objet à l'inventaire (fusionne les quantités si l'objet existe déjà).
func addInventory(p *Personnage, o Objet) error {
	for i, item := range p.Inventaire {
		if item.Nom == o.Nom && item.Type == o.Type {
			p.Inventaire[i].Quantite += o.Quantite
			return nil
		}
	}
	if checkInventoryFull(p) {
		return fmt.Errorf("Inventaire plein (%d/%d) : impossible d'ajouter %s", len(p.Inventaire), p.InventaireMax, o.Nom)
	}
	p.Inventaire = append(p.Inventaire, o)
	return nil
}

// removeInventory retire qte unités de l'objet nom de l'inventaire.
func removeInventory(p *Personnage, nom string, qte int) error {
	for i, item := range p.Inventaire {
		if item.Nom == nom {
			if item.Quantite < qte {
				return fmt.Errorf("quantité insuffisante de %s", nom)
			}
			p.Inventaire[i].Quantite -= qte
			if p.Inventaire[i].Quantite <= 0 {
				p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
			}
			return nil
		}
	}
	return fmt.Errorf("tu ne possèdes pas %s", nom)
}

// maxUpgrades est le nombre maximal d'augmentations d'inventaire autorisées.
const maxUpgrades = 3

// upgradeInventorySlot augmente la capacité de l'inventaire de +10 (3 fois maximum).
func upgradeInventorySlot(p *Personnage) error {
	if p.InventoryUpgrades >= maxUpgrades {
		return fmt.Errorf("limite d'améliorations d'inventaire atteinte (%d/%d)", p.InventoryUpgrades, maxUpgrades)
	}
	p.InventaireMax += 10
	p.InventoryUpgrades++
	return nil
}
