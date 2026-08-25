package main

import "fmt"

func main() {
	ClearScreen()
	titleText := getAscii("title")
	if titleText != "" {
		fmt.Println(titleText)
	} else {
		fmt.Println("\n\n\n\n\n\n\n\n")
		fmt.Println("                              ========================================")
		fmt.Println("                                        POKÉMON VERSION NOIRE         ")
		fmt.Println("                              ========================================")
	}
	fmt.Println("\n                               Appuyez sur Entrée pour commencer l'aventure...")
	lireLigne("")

	perso := Init()

	for {
		ClearScreen()
		fmt.Println("\n=== Menu principal ===")
		fmt.Println("1. Voir mes informations")
		fmt.Println("2. Inventaire")
		fmt.Println("3. Aller au Marchand")
		fmt.Println("4. Aller au Forgeron")
		fmt.Println("5. Aller dans les Hautes Herbes (Entrainement)")
		fmt.Println("6. Explorer la carte de Renouet (MAP INTERACTIVE)")
		fmt.Println("7. Qui sont-ils ?")
		fmt.Println("8. Quitter")

		choix := lireLigne("> ")
		switch choix {
		case "1":
			displayInfo(perso)
			lireLigne("\nAppuyez sur Entrée pour continuer...")
		case "2":
			accessInventory(perso)
		case "3":
			merchantMenu(perso)
		case "4":
			forgeronMenu(perso)
		case "5":
			trainingFight(perso)
			lireLigne("\nAppuyez sur Entrée pour continuer...")
		case "6":
			explorerMap(perso)
		case "7":
			whoAreThey()
			lireLigne("\nAppuyez sur Entrée pour continuer...")
		case "8":
			fmt.Println("À bientôt, " + perso.Nom + " !")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// whoAreThey affiche le nom des deux artistes cachés dans les diapositives 2 et 3 du sujet (Bonus 2).
func whoAreThey() {
	ClearScreen()
	fmt.Println("\n=== Qui sont-ils ? ===")
	TypeText("Ce sont des easter eggs cachés dans le PDF du projet (Bonus 2).")
	TypeText("On y retrouve peut-être des références à des artistes bien connus (ex: Daft Punk, PNL...) !")
}
