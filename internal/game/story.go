package game

import "fmt"

// ---------------------------------------------------------------------------
// POKÉMON NOIR 3 — trame en chapitres.
//
// Des années après la dissolution de la Team Plasma, un groupe dissident, la
// « Neo-Plasma », vole les Pokémon des dresseurs au nom d'une fausse libération.
// La Professeure Keteleeria envoie le joueur enquêter depuis Renouet.
// ---------------------------------------------------------------------------

// objectifChapitre renvoie la ligne d'objectif affichée sur la carte.
func objectifChapitre(n int) string {
	switch n {
	case 1:
		return "Rejoins la Route 1 au nord, puis la ville d'Arabelle."
	case 2:
		return "À Arabelle : équipe-toi, soigne-toi, puis bats le Maître d'Arène."
	case 3:
		return "Traverse la Forêt d'Ombreflore jusqu'au repaire de la Neo-Plasma."
	case 4:
		return "Infiltre le QG de la Neo-Plasma et affronte Nikolai."
	case 5:
		return "Affronte N et son Zekrom pour clore cette histoire."
	default:
		return "L'aventure est accomplie. Explore librement la région."
	}
}

// configurerRival fixe le rival (sexe opposé au joueur) et met à jour son PNJ
// sur la carte de Renouet. À appeler après la création OU le chargement.
func configurerRival(p *Personnage) {
	if p.Genre == "Garçon" {
		p.Rival = "Bianca"
	} else {
		p.Rival = "Hugo"
	}
	if z := zones["renouet"]; z != nil {
		for _, n := range z.PNJs {
			if n.Glyph == "H" {
				n.Nom = p.Rival + " (rival)"
				n.Portrait = rivalAscii(p)
			}
		}
	}
}

// rivalAscii renvoie le nom de fichier ascii du rival.
func rivalAscii(p *Personnage) string {
	if p.Rival == "Bianca" {
		return "bianca"
	}
	return "hugo"
}

// soigner remet les PV et le mana au maximum (et soigne toute l'équipe).
func soigner(p *Personnage) {
	p.PV = p.PVMax
	p.Mana = p.ManaMax
	soignerEquipe(p)
}

// prologue : scène d'introduction jouée juste après la création du personnage.
func prologue(p *Personnage) {
	configurerRival(p)

	ClearScreen()
	afficherAscii("professor")
	jouerDialogue("Professeure Keteleeria", []string{
		fmt.Sprintf("%s, la région d'Unys n'a jamais vraiment retrouvé la paix.", p.Nom),
		"Un groupe se fait appeler la « Neo-Plasma ». Il « libère » les Pokémon...",
		"...en les arrachant de force à leurs dresseurs.",
		"Prends ce Pokédex et ces 3 Potions. Va voir ce qui se trame à Arabelle.",
	})

	ClearScreen()
	afficherAscii(rivalAscii(p))
	jouerDialogue(p.Rival+" (rival)", []string{
		fmt.Sprintf("Alors %s, toi aussi la Prof t'envoie enquêter ?", p.Nom),
		"On se retrouve sur la route. Le premier à Arabelle a gagné !",
	})

	soigner(p)
	ClearScreen()
	fmt.Println("La Professeure soigne ton équipe avant le grand départ.")
	fmt.Println(">>> Choisis « Explorer la région » dans le menu pour commencer.")
	fmt.Printf(">>> Objectif : %s\n", objectifChapitre(p.Chapitre))
	lireLigne("\nAppuyez sur Entrée...")
}

// sceneLabo : le laboratoire de la Professeure à Renouet (soin + Pokéballs).
func sceneLabo(p *Personnage) {
	ClearScreen()
	afficherAscii("professor")
	jouerDialogue("Professeure Keteleeria", []string{
		"Repose-toi ici quand tu veux, ton équipe sera soignée.",
	})
	soigner(p)
	if !p.LaboPokeballs {
		addInventory(p, Objet{Nom: "Pokéball", Quantite: 5, Type: "Pokeball"})
		p.LaboPokeballs = true
		jouerDialogue("Professeure Keteleeria", []string{"Tiens, 5 Pokéballs pour la route. Sois prudent !"})
	}
	fmt.Printf("\n%s : PV %d/%d - Mana %d/%d\n", p.Nom, p.PV, p.PVMax, p.Mana, p.ManaMax)
	lireLigne("\nAppuyez sur Entrée...")
}

// centrePokemon : soin + sauvegarde (voir save.go).
func centrePokemon(p *Personnage) {
	for {
		ClearScreen()
		afficherAscii("infirmiere")
		fmt.Println(titre("Centre Pokémon"))
		afficherEtatEquipe(p)
		fmt.Println()
		fmt.Println("1. Soigner toute l'équipe")
		fmt.Println("2. Sauvegarder la partie")
		fmt.Println("3. Charger la dernière sauvegarde")
		fmt.Println("R. Sortir")
		switch lireLigne("> ") {
		case "1":
			soigner(p)
			TypeText("L'infirmière soigne votre équipe. Vos Pokémon débordent d'énergie !")
			lireLigne("\nAppuyez sur Entrée...")
		case "2":
			if err := sauvegarder(p); err != nil {
				fmt.Println("Échec de la sauvegarde : " + err.Error())
			} else {
				fmt.Println("Partie sauvegardée.")
			}
			lireLigne("\nAppuyez sur Entrée...")
		case "3":
			if err := charger(p); err != nil {
				fmt.Println("Aucune sauvegarde valide : " + err.Error())
			} else {
				fmt.Println("Sauvegarde chargée.")
			}
			lireLigne("\nAppuyez sur Entrée...")
		case "r", "R":
			return
		}
	}
}

// jouerScene exécute une scène scriptée. Renvoie false si la scène a refusé de se
// jouer (conditions non remplies) : le déclencheur reste alors actif.
func jouerScene(p *Personnage, id string) bool {
	switch id {
	case "ch2_grunt":
		return sceneGruntArabelle(p)
	case "final_n":
		return sceneFinaleN(p)
	default:
		jouerDialogue("", []string{"(...)"})
		return true
	}
}

// sceneGruntArabelle : à l'entrée d'Arabelle, un sbire vole le Pokémon d'un enfant.
func sceneGruntArabelle(p *Personnage) bool {
	if p.Chapitre >= 2 {
		return true // déjà fait
	}
	musique("plasma")
	ClearScreen()
	afficherAscii("team_plasme")
	jouerDialogue("Sbire Neo-Plasma", []string{
		"Ce Pokémon souffre entre tes mains, gamin ! La Neo-Plasma le libère !",
	})
	ClearScreen()
	afficherAscii("enfant")
	jouerDialogue("Enfant", []string{"Au secours ! Il veut prendre mon Pokémon !"})

	jouerDialogue(p.Nom, []string{"Hé ! Lâche ce Pokémon, tout de suite !"})

	ClearScreen()
	afficherAscii(rivalAscii(p))
	jouerDialogue(p.Rival+" (rival)", []string{
		fmt.Sprintf("%s ? Toujours à jouer les héros...", p.Nom),
		"Vas-y, occupe-toi de ce sbire. Je surveille l'enfant.",
	})

	grunt := &PNJ{Nom: "Sbire Neo-Plasma", Theme: "plasma", Equipe: []string{"Ponchiot"}}
	if combatDresseur(p, grunt) {
		ClearScreen()
		afficherAscii("team_plasme")
		jouerDialogue("Sbire Neo-Plasma", []string{"Grr... Tu le regretteras. La Neo-Plasma est partout !"})

		ClearScreen()
		afficherAscii("enfant")
		jouerDialogue("Enfant", []string{
			"Merci ! Tiens, prends ça pour te remercier.",
			"Le Maître d'Arène te cherchait justement.",
		})
		p.Or += 200
		addInventory(p, Objet{Nom: "Pokéball", Quantite: 3, Type: "Pokeball"})
		addInventory(p, Objet{Nom: "Potion (Restaure 50 PV)", Quantite: 2, Type: "Potion"})
		fmt.Println("\nVous recevez : 200 Poké-Dollars, 3 Pokéballs, 2 Potions !")

		jouerDialogue(p.Rival+" (rival)", []string{"Pas mal. On se retrouve à l'Arène, ne traîne pas."})

		p.Chapitre = 2
		fmt.Printf("\n>>> Nouvel objectif : %s\n", objectifChapitre(p.Chapitre))
		lireLigne("\nAppuyez sur Entrée...")
	} else {
		jouerDialogue(p.Rival+" (rival)", []string{"Pathétique... Va te soigner au Centre et reviens."})
	}
	return p.Chapitre >= 2
}

// sceneFinaleN : dans le QG, N apparaît. Selon l'avancement : pas prêt / combat /
// épilogue si l'histoire est déjà finie.
func sceneFinaleN(p *Personnage) bool {
	musique("n")
	ClearScreen()
	afficherAscii("N")

	if p.Chapitre >= 6 || flagOK(p, "fin") {
		jouerDialogue("N", []string{
			"Unys respire de nouveau. Prends soin de Zekrom.",
			"Nos routes se recroiseront, j'en suis sûr.",
		})
		return true
	}
	if flagOK(p, "n_battu") {
		// N déjà battu (Zekrom obtenu) mais Ghetsis pas encore vaincu.
		jouerDialogue("N", []string{"Ghetsis est toujours là. Finissons-en ensemble."})
		combatFinalGhetsis(p)
		return true
	}
	if p.Chapitre < 5 {
		jouerDialogue("N", []string{
			"Tu n'es pas encore prêt.",
			"Nikolai t'attend au fond du repaire. Bats-le, et reviens me voir.",
		})
		return false
	}
	combatBossN(p)
	return true
}

// combatBossN : N & Zekrom. En cas de victoire, Zekrom rejoint l'équipe du joueur.
func combatBossN(p *Personnage) {
	jouerDialogue("N", []string{
		"Mon nom est N. Autrefois, j'ai voulu séparer les hommes et les Pokémon.",
		"J'avais tort. Mais la Neo-Plasma a détourné mon rêve.",
		fmt.Sprintf("%s... montre-moi la force de tes convictions ! Zekrom, à nous !", p.Nom),
	})

	musiqueCombat = "N_battle"
	musique("N_battle")
	boss := nouveauMonstre("Zekrom")
	boss.Nom = "Zekrom de N"
	boss.Niveau = 25
	if !deroulerCombat(p, boss, false) {
		fmt.Printf("\n%s s'effondre face à la puissance de Zekrom...\n", p.Nom)
		dead(p)
		lireLigne("\nAppuyez sur Entrée...")
		return
	}
	flagSet(p, "n_battu")

	musique("n")
	ClearScreen()
	afficherAscii("N")
	jouerDialogue("N", []string{
		"Je vois... Tes Pokémon se battent pour toi, pas par peur.",
		"Zekrom a senti la pureté de ton cœur. Il te choisit.",
	})

	// Zekrom rejoint l'équipe (fort, mais gagné — pas cheaté).
	don := nouveauMonstre("Zekrom")
	don.Nom = "Zekrom"
	if ajouterAEquipe(p, don, 25) {
		mettreAJourSortsPokemon(&p.Equipe[len(p.Equipe)-1])
		bruitage("level") // son court : ne se superpose pas au thème de N
		TypeText("Zekrom rejoint ton équipe !")
	} else {
		addInventory(p, Objet{Nom: "Zekrom", Quantite: 1, Type: "Pokemon"})
		TypeText("Ton équipe est pleine : Zekrom t'attendra au Centre Pokémon.")
	}
	jouerDialogue("N", []string{"Mais tout n'est pas fini. Il est là. Derrière toi."})
	lireLigne("\nAppuyez sur Entrée...")

	combatFinalGhetsis(p)
}

// combatFinalGhetsis : le vrai commanditaire, avec Kyurem.
func combatFinalGhetsis(p *Personnage) {
	musique("ghetsis")
	ClearScreen()
	afficherAscii("ghetsis")
	jouerDialogue("Ghetsis", []string{
		"La Neo-Plasma, ce petit théâtre... tout cela était mon oeuvre.",
		"Cette fois, aucun idéaliste ne viendra t'aider. KYUREM, anéantis-le !",
	})

	musiqueCombat = "battle_boss"
	musique("battle_boss")
	kyurem := nouveauMonstre("Kyurem")
	kyurem.Nom = "Kyurem de Ghetsis"
	kyurem.Niveau = 30
	if !deroulerCombat(p, kyurem, false) {
		fmt.Printf("\n%s est vaincu... mais Ghetsis a sous-estimé ta détermination.\n", p.Nom)
		dead(p)
		lireLigne("\nAppuyez sur Entrée pour retenter...")
		combatFinalGhetsis(p)
		return
	}

	p.Or += 2000
	distribuerXP(p, kyurem.Experience)

	ClearScreen()
	afficherAscii("kyurem")
	jouerDialogue("", []string{
		"Libéré de l'emprise de Ghetsis, Kyurem pousse un rugissement...",
		"...et disparaît vers les montagnes gelées du nord.",
	})

	ClearScreen()
	afficherAscii("ghetsis")
	jouerDialogue("Ghetsis", []string{"Impossible... vaincu par un gamin et ses... « amis » ?!"})
	jouerDialogue("N", []string{
		"C'est terminé, Ghetsis. La police d'Unys t'attend.",
		fmt.Sprintf("Merci, %s. Grâce à toi, hommes et Pokémon peuvent grandir ensemble.", p.Nom),
	})

	ClearScreen()
	fmt.Println("===========================================================")
	fmt.Println("            FIN  —  POKÉMON NOIR 3")
	fmt.Println("   Zekrom à tes côtés, Kyurem quelque part dans les glaces...")
	fmt.Println("      Merci d'avoir joué. L'exploration reste ouverte.")
	fmt.Println("===========================================================")
	if p.Chapitre < 6 {
		p.Chapitre = 6
	}
	flagSet(p, "fin")
	lireLigne("\nAppuyez sur Entrée...")
}
