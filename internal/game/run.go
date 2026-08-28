package game

import (
	"fmt"
	"os"
)

// Run est le point d'entrée du jeu (appelé par le main.go racine).
func Run() {
	if os.Getenv("RED_FAST") != "" {
		vitesseTexte = 0
	}
	if os.Getenv("NO_COLOR") != "" {
		couleurActive = false
	}
	ClearScreen()
	soundInit() // peut poser une question -> avant l'écran-titre
	musique("title")
	ClearScreen()
	if titleText := getAscii("title"); titleText != "" {
		// Le dessin ascii/title.txt contient déjà son invite « Appuyez sur Entrée ».
		fmt.Println(titleText)
	} else {
		fmt.Print("\n\n\n\n\n\n\n\n")
		fmt.Println("                    ========================================")
		fmt.Println("                              POKÉMON NOIR 3")
		fmt.Println("                    ========================================")
		fmt.Println("\n              Appuyez sur Entrée pour commencer l'aventure...")
	}
	lireLigne("")

	initMonde()

	perso := demarrerPartie()
	appliquerProgression(perso) // dresseurs déjà battus -> restaurés

	for {
		musique("menu")
		ClearScreen()
		fmt.Println("\n" + titre("Menu principal"))
		fmt.Printf("%s %s\n\n", col(fmt.Sprintf("Chapitre %d —", perso.Chapitre), cJaune), objectifChapitre(perso.Chapitre))
		afficherEtatEquipe(perso)
		fmt.Println()
		fmt.Println("1. Voir mes informations")
		fmt.Println("2. Sac (inventaire)")
		fmt.Println("3. Aller au Marchand")
		fmt.Println("4. Aller au Forgeron")
		fmt.Println("5. Zone d'entraînement (combat rapide)")
		fmt.Println("6. Explorer la région  <<< L'AVENTURE EST ICI")
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
			explorerMonde(perso)
		case "7":
			whoAreThey()
			lireLigne("\nAppuyez sur Entrée pour continuer...")
		case "8":
			if err := sauvegarder(perso); err == nil {
				fmt.Println("Partie sauvegardée automatiquement.")
			}
			soundQuit()
			fmt.Println("À bientôt, " + perso.Nom + " !")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// whoAreThey répond au Bonus 2 : deux artistes sont cachés dans les parties 2 et 3
// du sujet, via les titres des tâches.
func whoAreThey() {
	ClearScreen()
	fmt.Println("\n=== Qui sont-ils ? ===")
	TypeText("Deux artistes se cachent dans les titres des tâches du sujet :")
	fmt.Println()
	TypeText("  Partie 2  ->  ABBA")
	TypeText("    \"Money, Money, Money\", \"Gimme! Gimme! Gimme!\", \"Mamma Mia\",")
	TypeText("    \"On and On and On\", \"Two for the Price of One\" (Take a Chance on Me).")
	fmt.Println()
	TypeText("  Partie 3  ->  Steven Spielberg")
	TypeText("    \"Duel\", \"A.I. Intelligence Artificielle\", \"Ready Player One\".")
}
