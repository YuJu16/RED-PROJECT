package game

import "fmt"

// ArticleMarchand décrit un article vendu par le marchand.
type ArticleMarchand struct {
	Nom   string
	Prix  int
	Desc  string
	Objet Objet
}

// nomAugmentationInventaire est le nom spécial déclenchant upgradeInventorySlot à l'achat.
const nomAugmentationInventaire = "Sac de Dresseur (+10 places)"

// catalogueMarchand liste tous les articles disponibles chez le marchand.
var catalogueMarchand = []ArticleMarchand{
	{"Potion (Restaure 50 PV)", 3, "Rend 50 PV au Pokémon actif (jamais au-dessus du max).",
		Objet{Nom: "Potion (Restaure 50 PV)", Quantite: 1, Type: "Potion"}},
	{"Baie Pecha (Poison)", 6, "À manger... très mauvaise idée : -10 PV/s pendant 3 s.",
		Objet{Nom: "Baie Pecha (Poison)", Quantite: 1, Type: "PotionPoison"}},
	{"Huile (Restaure PP/Mana)", 5, "Rend 30 points de Mana, pour relancer des sorts.",
		Objet{Nom: "Huile (Restaure PP/Mana)", Quantite: 1, Type: "PotionMana"}},
	{"CT35 - Lance-Flammes", 25, "Apprend le sort Lance-Flammes (18 dégâts, 10 mana). Une seule fois.",
		Objet{Nom: "CT35 - Lance-Flammes", Quantite: 1, Type: "LivreSort"}},
	{"Plume de Poichigeon", 4, "Ressource de forge (Veste de Combat, Grelot Coque).",
		Objet{Nom: "Plume de Poichigeon", Quantite: 1, Type: "Ressource"}},
	{"Peau de Grotichon", 7, "Ressource de forge (Veste de Combat).",
		Objet{Nom: "Peau de Grotichon", Quantite: 1, Type: "Ressource"}},
	{"Cuir de Roitiflam", 3, "Ressource de forge (Restes, Grelot Coque).",
		Objet{Nom: "Cuir de Roitiflam", Quantite: 1, Type: "Ressource"}},
	{"Plume de Déflaisan", 1, "Ressource de forge (Restes).",
		Objet{Nom: "Plume de Déflaisan", Quantite: 1, Type: "Ressource"}},
	{"Pokéball", 10, "Sert à capturer un Pokémon sauvage affaibli (>40 % de PV en moins).",
		Objet{Nom: "Pokéball", Quantite: 1, Type: "Pokeball"}},
	{nomAugmentationInventaire, 30, "Ajoute 10 emplacements d'inventaire (3 fois maximum).",
		Objet{}},
}

// quantiteObjet renvoie combien d'exemplaires de nom le joueur possède.
func quantiteObjet(p *Personnage, nom string) int {
	for _, o := range p.Inventaire {
		if o.Nom == nom {
			return o.Quantite
		}
	}
	return 0
}

// merchantMenu affiche le marchand et gère les achats (avec quantité).
func merchantMenu(p *Personnage) {
	for {
		ClearScreen()
		fmt.Println("\n" + titre("Marchand"))
		fmt.Printf("Or : %s\n\n", col(fmt.Sprintf("%d", p.Or), cJaune))
		if !p.PotionGratuitePrise {
			fmt.Println(col("0.", cVert) + " Échantillon gratuit : Potion (Restaure 50 PV) — 0 or (1 fois)")
		}
		for i, a := range catalogueMarchand {
			stock := ""
			if q := quantiteObjet(p, a.Objet.Nom); q > 0 {
				stock = col(fmt.Sprintf("  [x%d en stock]", q), cCyan)
			}
			fmt.Printf("%2d. %-26s %4d or%s\n", i+1, a.Nom, a.Prix, stock)
			fmt.Printf("      %s\n", col(a.Desc, cJaune))
		}
		fmt.Println("R. Retour au menu")

		choix := lireLigne("> ")
		if choix == "R" || choix == "r" {
			return
		}

		if choix == "0" && !p.PotionGratuitePrise {
			if err := addInventory(p, Objet{Nom: "Potion (Restaure 50 PV)", Quantite: 1, Type: "Potion"}); err != nil {
				fmt.Println("Erreur : " + err.Error())
			} else {
				p.PotionGratuitePrise = true
				fmt.Println("Tu as reçu ton échantillon gratuit : Potion (Restaure 50 PV)")
			}
			lireLigne("\nAppuyez sur Entrée...")
			continue
		}

		index := indexFromChoice(choix, len(catalogueMarchand))
		if index == -1 {
			fmt.Println("Choix invalide.")
			lireLigne("\nAppuyez sur Entrée...")
			continue
		}
		acheter(p, catalogueMarchand[index])
		lireLigne("\nAppuyez sur Entrée...")
	}
}

// acheter gère l'achat d'un article, avec choix de la quantité pour les objets
// empilables (Pokéball, ressources, potions...).
func acheter(p *Personnage, a ArticleMarchand) {
	// Cas particulier : l'agrandissement d'inventaire s'achète à l'unité.
	if a.Nom == nomAugmentationInventaire {
		if p.Or < a.Prix {
			fmt.Printf("Pas assez d'or ! Il te faut %d or, tu en as %d.\n", a.Prix, p.Or)
			return
		}
		if err := upgradeInventorySlot(p); err != nil {
			fmt.Println("Erreur : " + err.Error())
			return
		}
		p.Or -= a.Prix
		fmt.Printf("Ton inventaire passe à %d emplacements !\n", p.InventaireMax)
		return
	}

	qte := demanderQuantite(p, a)
	if qte <= 0 {
		fmt.Println("Achat annulé.")
		return
	}

	achetes := 0
	for i := 0; i < qte; i++ {
		if p.Or < a.Prix {
			break
		}
		obj := a.Objet
		obj.Quantite = 1
		if err := addInventory(p, obj); err != nil {
			fmt.Println("Erreur : " + err.Error())
			break
		}
		p.Or -= a.Prix
		achetes++
	}

	switch {
	case achetes == 0:
		fmt.Printf("Pas assez d'or ! Il te faut %d or par %s, tu en as %d.\n", a.Prix, a.Nom, p.Or)
	case achetes < qte:
		fmt.Printf("Tu as acheté %d x %s (or ou place insuffisants pour le reste). Or : %d\n", achetes, a.Nom, p.Or)
	default:
		fmt.Printf("Tu as acheté %d x %s. Or restant : %d\n", achetes, a.Nom, p.Or)
	}
}

// demanderQuantite propose une quantité par défaut (1) et le maximum abordable.
func demanderQuantite(p *Personnage, a ArticleMarchand) int {
	maxAbordable := 0
	if a.Prix > 0 {
		maxAbordable = p.Or / a.Prix
	}
	if maxAbordable < 1 {
		return 1 // on laisse acheter() afficher le message "pas assez d'or"
	}
	saisie := lireLigne(fmt.Sprintf("Combien ? (Entrée = 1, max %d) > ", maxAbordable))
	if saisie == "" {
		return 1
	}
	n := 0
	if _, err := fmt.Sscanf(saisie, "%d", &n); err != nil || n < 1 {
		fmt.Println("Quantité invalide, on prend 1.")
		return 1
	}
	if n > maxAbordable {
		n = maxAbordable
	}
	return n
}
