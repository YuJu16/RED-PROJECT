package main

import (
	"fmt"
	"strings"
	"time"
)

var worldMap = []string{
	"#################################################################################",
	"#                                                                               #",
	"#  [ ARABELLE ]                                                                 #",
	"#  ________________________                                                     #",
	"# |  [Centre Pokémon]      |                                                    #",
	"# |         (M)            |     ________                 [ ROUTE 1 ]           #",
	"# |                        |    | BOSS N |                                      #",
	"# |  [Atelier de Forge]    |    |  (N)   |              ,,,,,,,,,,,,,           #",
	"# |         (F)            |    |________|              ,           ,           #",
	"# |                        |                            ,           ,           #",
	"# |________________________|----------------------------,           ,           #",
	"#                                                       ,           ,           #",
	"#                                                       ,,,,,,,,,,,,,           #",
	"#                                                             |                 #",
	"#                                                             |                 #",
	"#                                                             |                 #",
	"#                                                       ______|______           #",
	"#                                                      |             |          #",
	"#                                                      |  [RENOUET]  |          #",
	"#                                                      |   Maison    |          #",
	"#                                                      |    (P)      |          #",
	"#                                                      |_____________|          #",
	"#                                                                               #",
	"#################################################################################",
}

func explorerMap(perso *Personnage) {
	// Coordonnées de départ du joueur (devant sa maison)
	playerX, playerY := 63, 19

	for {
		ClearScreen()
		fmt.Println("=== EXPLORATION DE LA RÉGION D'UNYS ===")
		fmt.Println("Utilisez Z (Haut), Q (Gauche), S (Bas), D (Droite) pour vous déplacer.")
		fmt.Println("Appuyez sur 'R' pour retourner au menu classique.")
		fmt.Println("Lieux : (P) Chez vous, (M) Marchand, (F) Forgeron, (N) Boss N, (,) Hautes Herbes")
		fmt.Println()

		// Affichage de la carte avec le joueur
		for y, ligne := range worldMap {
			if y == playerY {
				runes := []rune(ligne)
				runes[playerX] = 'X' // 'X' représente le joueur
				fmt.Println(string(runes))
			} else {
				fmt.Println(ligne)
			}
		}

		choix := strings.ToLower(lireLigne("Action > "))

		newX, newY := playerX, playerY

		switch choix {
		case "z":
			newY--
		case "s":
			newY++
		case "q":
			newX--
		case "d":
			newX++
		case "r":
			return // Retour au menu principal
		}

		// Collisions et interactions
		if newY >= 0 && newY < len(worldMap) && newX >= 0 && newX < len(worldMap[newY]) {
			cell := worldMap[newY][newX]
			if cell == ' ' || cell == ',' || cell == '|' || cell == '_' || cell == '-' {
				playerX, playerY = newX, newY
			} else if cell == 'M' {
				merchantMenu(perso)
			} else if cell == 'F' {
				forgeronMenu(perso)
			} else if cell == 'P' {
				displayInfo(perso)
				lireLigne("\nAppuyez sur Entrée pour fermer vos informations...")
			} else if cell == 'N' {
				combatBossN(perso)
				playerY++ // Recule après le combat
			}
		}

		// Zone des hautes herbes (Route 1)
		if playerY >= 7 && playerY <= 12 && playerX >= 56 && playerX <= 68 {
			// 20% de chance de lancer un combat à chaque pas dans l'herbe
			if time.Now().UnixNano()%5 == 0 {
				ClearScreen()
				TypeText("Vous marchez dans les hautes herbes... Un Pokémon sauvage attaque !")
				trainingFight(perso)
				lireLigne("\nAppuyez sur Entrée pour continuer l'exploration...")
			}
		}
	}
}

func combatBossN(p *Personnage) {
	ClearScreen()
	TypeText("N : Mon nom est N. Je suis le roi de la Team Plasma.")
	TypeText("N : Tu penses pouvoir enfermer les Pokémon dans des Pokéballs ?")
	TypeText("N : Montre-moi la force de tes convictions !")
	time.Sleep(2 * time.Second)

	// Création du boss
	boss := &Monster{
		Nom:        "Zekrom (Boss)",
		PVMax:      150,
		PV:         150,
		Attaque:    12,
		Initiative: 8,
		Experience: 500,
		Ascii:      "zekrom", // Utilise l'ascii de Zekrom
	}

	ClearScreen()
	fmt.Println(getAscii(boss.Ascii))
	TypeText("Le Boss N envoie Zekrom !")

	// Boucle de combat similaire à trainingFight
	tour := 1
	joueurCommence := p.Initiative >= boss.Initiative

	for {
		fmt.Printf("\n--- Tour %d ---\n", tour)

		if joueurCommence {
			charTurn(p, boss)
			if boss.PV <= 0 {
				break
			}
			monsterPattern(boss, p, tour)
			if p.PV <= 0 {
				break
			}
		} else {
			monsterPattern(boss, p, tour)
			if p.PV <= 0 {
				break
			}
			charTurn(p, boss)
			if boss.PV <= 0 {
				break
			}
		}
		tour++
	}

	if p.PV <= 0 {
		fmt.Printf("\n%s est tombé au combat face à N...\n", p.Nom)
		dead(p)
	} else {
		ClearScreen()
		fmt.Println(getAscii(boss.Ascii))
		TypeText("N : Je vois... Tes sentiments pour tes Pokémon sont réels.")
		TypeText("N : Tu as gagné. Je vais dissoudre la Team Plasma.")
		TypeText("FÉLICITATIONS ! VOUS AVEZ BATTU LE JEU !")
		p.Or += 500
		gainExperience(p, boss.Experience)
		lireLigne("\nAppuyez sur Entrée pour retourner à l'exploration...")
	}
}
