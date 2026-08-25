package main

import "fmt"

// Equipment regroupe les emplacements d'équipement du personnage.
type Equipment struct {
	Tete  string
	Torse string
	Pieds string
}

// Recette décrit un équipement fabricable par le forgeron.
type Recette struct {
	Nom        string
	Slot       string // "Tete", "Torse" ou "Pieds"
	Ressources map[string]int
	BonusPV    int
}

// recettesForgeron liste les équipements fabricables et leur coût en ressources.
var recettesForgeron = []Recette{
	{
		Nom:        "Veste de Combat",
		Slot:       "Torse",
		Ressources: map[string]int{"Plume de Déflaisan": 1, "Cuir de Roitiflam": 1},
		BonusPV:    10,
	},
	{
		Nom:        "Restes",
		Slot:       "Tete", // On met 'Tete' juste pour l'emplacement d'équipement dans notre struct générique
		Ressources: map[string]int{"Plume de Poichigeon": 2, "Peau de Grotichon": 1},
		BonusPV:    25,
	},
	{
		Nom:        "Grelot Coque",
		Slot:       "Pieds", // Idem
		Ressources: map[string]int{"Plume de Poichigeon": 1, "Cuir de Roitiflam": 1},
		BonusPV:    15,
	},
}

// coutForge est le coût en or de toute fabrication chez le forgeron.
const coutForge = 5

// forgeronMenu affiche le menu du forgeron et gère la fabrication d'équipement.
func forgeronMenu(p *Personnage) {
	for {
		ClearScreen()
		fmt.Println("\n=== Forgeron ===")
		fmt.Printf("Or : %d\n", p.Or)
		for i, r := range recettesForgeron {
			fmt.Printf("%d. %s (%d or)\n", i+1, r.Nom, coutForge)
		}
		fmt.Println("R. Retour au menu")

		choix := lireLigne("> ")
		if choix == "R" || choix == "r" {
			return
		}
		index := indexFromChoice(choix, len(recettesForgeron))
		if index == -1 {
			fmt.Println("Choix invalide.")
			continue
		}
		fabriquer(p, recettesForgeron[index])
	}
}

// fabriquer vérifie l'or et les ressources nécessaires puis fabrique l'équipement demandé.
func fabriquer(p *Personnage, r Recette) {
	if p.Or < coutForge {
		fmt.Printf("Pas assez d'or ! Il te faut %d or pour fabriquer %s.\n", coutForge, r.Nom)
		return
	}
	for nom, qte := range r.Ressources {
		if !possede(p, nom, qte) {
			fmt.Printf("Pas les ressources nécessaires : il te faut %d %s pour fabriquer %s.\n", qte, nom, r.Nom)
			return
		}
	}
	if checkInventoryFull(p) {
		fmt.Println("Inventaire plein ! Impossible de fabriquer " + r.Nom + ".")
		return
	}

	for nom, qte := range r.Ressources {
		removeInventory(p, nom, qte)
	}
	p.Or -= coutForge
	addInventory(p, Objet{Nom: r.Nom, Quantite: 1, Type: "Equipement", Slot: r.Slot})
	fmt.Println("Tu as fabriqué : " + r.Nom)
}

// possede indique si le personnage possède au moins qte exemplaires de l'objet nom.
func possede(p *Personnage, nom string, qte int) bool {
	for _, item := range p.Inventaire {
		if item.Nom == nom {
			return item.Quantite >= qte
		}
	}
	return false
}

// bonusPVEquipement renvoie le bonus de PV max associé à un équipement fabricable.
func bonusPVEquipement(nom string) int {
	for _, r := range recettesForgeron {
		if r.Nom == nom {
			return r.BonusPV
		}
	}
	return 0
}

// equip équipe l'objet item (retiré de l'inventaire), applique son bonus de PV max,
// et renvoie dans l'inventaire l'équipement précédemment porté sur le même emplacement.
func equip(p *Personnage, item Objet) {
	var ancien string
	switch item.Slot {
	case "Tete":
		ancien = p.Equipement.Tete
		p.Equipement.Tete = item.Nom
	case "Torse":
		ancien = p.Equipement.Torse
		p.Equipement.Torse = item.Nom
	case "Pieds":
		ancien = p.Equipement.Pieds
		p.Equipement.Pieds = item.Nom
	default:
		fmt.Println("Cet objet ne peut pas être équipé.")
		return
	}

	removeInventory(p, item.Nom, 1)
	bonus := bonusPVEquipement(item.Nom)
	p.PVMax += bonus
	if p.PV > p.PVMax {
		p.PV = p.PVMax
	}
	fmt.Printf("Tu équipes : %s (+%d PV max)\n", item.Nom, bonus)

	if ancien != "" && ancien != item.Nom {
		p.PVMax -= bonusPVEquipement(ancien)
		if p.PV > p.PVMax {
			p.PV = p.PVMax
		}
		addInventory(p, Objet{Nom: ancien, Quantite: 1, Type: "Equipement", Slot: item.Slot})
		fmt.Println(ancien + " retourne dans ton inventaire.")
	}
}
