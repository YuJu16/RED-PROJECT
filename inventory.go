package main

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
		fmt.Println("\n=== Inventaire ===")
		fmt.Printf("(%d/%d emplacements utilisés)\n", len(p.Inventaire), p.InventaireMax)
		if len(p.Inventaire) == 0 {
			fmt.Println("(vide)")
		}
		for i, o := range p.Inventaire {
			fmt.Printf("%d. %s x%d\n", i+1, o.Nom, o.Quantite)
		}
		fmt.Println("R. Retour au menu")

		choix := lireLigne("Choisis un objet à utiliser (numéro) ou R pour revenir > ")
		if choix == "R" || choix == "r" {
			return
		}

		index := indexFromChoice(choix, len(p.Inventaire))
		if index == -1 {
			fmt.Println("Choix invalide.")
			continue
		}
		useItem(p, index)
	}
}

// useItem applique l'effet de l'objet situé à l'index i de l'inventaire.
func useItem(p *Personnage, i int) {
	if i < 0 || i >= len(p.Inventaire) {
		fmt.Println("Choix invalide.")
		return
	}
	item := p.Inventaire[i]
	switch item.Type {
	case "Potion":
		takePot(p)
	case "PotionPoison":
		poisonPot(p)
	case "PotionMana":
		takeManaPot(p)
	case "LivreSort":
		spellBook(p)
		removeInventory(p, item.Nom, 1)
	case "Equipement":
		equip(p, item)
	default:
		fmt.Println("Cet objet ne peut pas être utilisé directement.")
	}
}

// consumeOne retire une unité de l'objet situé à l'index i (le supprime si la quantité tombe à 0).
func consumeOne(p *Personnage, i int) {
	p.Inventaire[i].Quantite--
	if p.Inventaire[i].Quantite <= 0 {
		p.Inventaire = append(p.Inventaire[:i], p.Inventaire[i+1:]...)
	}
}

// takePot utilise une Potion de vie (soin de 50 PV, cap à PVMax).
func takePot(p *Personnage) {
	for i, o := range p.Inventaire {
		if o.Type == "Potion" {
			if p.PV >= p.PVMax {
				fmt.Println("Tes PV sont déjà au maximum !")
				return
			}
			p.PV += 50
			if p.PV > p.PVMax {
				p.PV = p.PVMax
			}
			consumeOne(p, i)
			fmt.Printf("Tu utilises une Potion de vie. PV : %d/%d\n", p.PV, p.PVMax)
			return
		}
	}
	fmt.Println("Erreur : tu n'as pas de Potion de vie !")
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
