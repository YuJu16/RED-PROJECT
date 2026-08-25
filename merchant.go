package main

import "fmt"

// ArticleMarchand décrit un article vendu par le marchand.
type ArticleMarchand struct {
	Nom   string
	Prix  int
	Objet Objet
}

// nomAugmentationInventaire est le nom spécial déclenchant upgradeInventorySlot à l'achat.
const nomAugmentationInventaire = "Sac de Dresseur (+10 places)"

// catalogueMarchand liste tous les articles disponibles chez le marchand et leur prix.
var catalogueMarchand = []ArticleMarchand{
	{"Potion (Restaure 50 PV)", 3, Objet{Nom: "Potion (Restaure 50 PV)", Quantite: 1, Type: "Potion"}},
	{"Baie Pecha (Poison)", 6, Objet{Nom: "Baie Pecha (Poison)", Quantite: 1, Type: "PotionPoison"}},
	{"Huile (Restaure PP/Mana)", 5, Objet{Nom: "Huile (Restaure PP/Mana)", Quantite: 1, Type: "PotionMana"}},
	{"CT35 - Lance-Flammes", 25, Objet{Nom: "CT35 - Lance-Flammes", Quantite: 1, Type: "LivreSort"}},
	{"Plume de Poichigeon", 4, Objet{Nom: "Plume de Poichigeon", Quantite: 1, Type: "Ressource"}},
	{"Peau de Grotichon", 7, Objet{Nom: "Peau de Grotichon", Quantite: 1, Type: "Ressource"}},
	{"Cuir de Roitiflam", 3, Objet{Nom: "Cuir de Roitiflam", Quantite: 1, Type: "Ressource"}},
	{"Plume de Déflaisan", 1, Objet{Nom: "Plume de Déflaisan", Quantite: 1, Type: "Ressource"}},
	{"Pokéball", 10, Objet{Nom: "Pokéball", Quantite: 1, Type: "Pokeball"}},
	{nomAugmentationInventaire, 30, Objet{}},
}

// merchantMenu affiche le menu du marchand et gère les achats.
func merchantMenu(p *Personnage) {
	for {
		ClearScreen()
		fmt.Println("\n=== Marchand ===")
		fmt.Printf("Or : %d\n", p.Or)
		for i, a := range catalogueMarchand {
			fmt.Printf("%d. %s - %d or\n", i+1, a.Nom, a.Prix)
		}
		fmt.Println("R. Retour au menu")

		choix := lireLigne("> ")
		if choix == "R" || choix == "r" {
			return
		}
		index := indexFromChoice(choix, len(catalogueMarchand))
		if index == -1 {
			fmt.Println("Choix invalide.")
			continue
		}
		acheter(p, catalogueMarchand[index])
	}
}

// acheter vérifie l'or du joueur puis effectue l'achat de l'article a.
func acheter(p *Personnage, a ArticleMarchand) {
	if p.Or < a.Prix {
		fmt.Printf("Pas assez d'or ! Il te faut %d or, tu en as %d.\n", a.Prix, p.Or)
		return
	}

	if a.Nom == nomAugmentationInventaire {
		if err := upgradeInventorySlot(p); err != nil {
			fmt.Println("Erreur : " + err.Error())
			return
		}
		p.Or -= a.Prix
		fmt.Printf("Ton inventaire passe à %d emplacements !\n", p.InventaireMax)
		return
	}

	if err := addInventory(p, a.Objet); err != nil {
		fmt.Println("Erreur : " + err.Error())
		return
	}
	p.Or -= a.Prix
	fmt.Println("Tu as acheté : " + a.Nom)
}
