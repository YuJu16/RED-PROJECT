package game

import "strings"

// ---------------------------------------------------------------------------
// Cartes du monde. Les grilles ne contiennent QUE de la structure :
//   '#' mur/bordure   ' ' sol   ',' hautes herbes   '^' arbre (bloquant)
// Les PNJ et les sorties sont posés par coordonnées et dessinés par-dessus,
// donc leurs cases doivent être du sol ' '.
// ---------------------------------------------------------------------------

// salle crée une pièce rectangulaire vide entourée de murs.
func salle(w, h int) []string {
	g := make([]string, h)
	for y := 0; y < h; y++ {
		if y == 0 || y == h-1 {
			g[y] = strings.Repeat("#", w)
		} else {
			g[y] = "#" + strings.Repeat(" ", w-2) + "#"
		}
	}
	return g
}

// stamp écrit s dans la grille à partir de (x,y) sans changer sa taille.
func stamp(g []string, x, y int, s string) {
	if y < 0 || y >= len(g) {
		return
	}
	r := []rune(g[y])
	for i, c := range s {
		if x+i > 0 && x+i < len(r)-1 { // ne pas écraser les murs de bord
			r[x+i] = c
		}
	}
	g[y] = string(r)
}

// herbe pose un rectangle de hautes herbes.
func herbe(g []string, x, y, w, h int) {
	for j := 0; j < h; j++ {
		stamp(g, x, y+j, strings.Repeat(",", w))
	}
}

func zoneRenouet() *Zone {
	g := salle(52, 15)
	stamp(g, 3, 2, "^^^^^        ^^^^^        ^^^^^        ^^^^")
	stamp(g, 3, 12, "^^^     ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^     ^^")
	return &Zone{
		Cle: "renouet", Nom: "RENOUET", Banniere: "Le village du vent — ton point de départ",
		Grille: g, Depart: [2]int{26, 11},
		PNJs: []*PNJ{
			{X: 12, Y: 5, Glyph: "L", Nom: "Laboratoire de la Professeure", Service: "labo"},
			{X: 39, Y: 5, Glyph: "H", Nom: "Hugo (rival)", Portrait: "hugo", Dialogue: []string{
				"Qu'est-ce que tu attends ? Arabelle ne viendra pas à toi !",
			}},
			{X: 20, Y: 7, Glyph: "P", Nom: "Maman", Portrait: "mere", Service: "maison", Dialogue: []string{
				"Tu pars déjà ? Viens là, je soigne tes Pokémon.",
				"Sois prudent avec cette Neo-Plasma... et écris-moi !",
			}},
			{X: 12, Y: 9, Glyph: "M", Nom: "Boutique", Service: "boutique"},
			{X: 39, Y: 9, Glyph: "F", Nom: "Forge", Service: "forge"},
		},
		Sorties: []*Sortie{
			{X: 26, Y: 1, ZoneCible: "route1", SpawnX: 26, SpawnY: 13, ChapitreMin: 1},
		},
	}
}

func zoneRoute1() *Zone {
	g := salle(52, 16)
	herbe(g, 4, 3, 12, 3)
	herbe(g, 36, 3, 12, 3)
	herbe(g, 4, 10, 12, 3)
	herbe(g, 30, 10, 14, 3)
	stamp(g, 20, 2, "^^^^")
	return &Zone{
		Cle: "route1", Nom: "ROUTE 1", Banniere: "Hautes herbes et jeunes dresseurs",
		Grille: g, Depart: [2]int{26, 13},
		Sauvages: []string{"Ratentif", "Poichigeon", "Ponchiot"},
		PNJs: []*PNJ{
			{X: 24, Y: 7, Glyph: "T", Nom: "Gamin Timéo", Portrait: "gamin", Butin: 60,
				Equipe:      []string{"Ratentif"},
				Dialogue:    []string{"Un vrai dresseur ! On se bat, dis, on se bat ?"},
				Apres:       []string{"T'es trop fort... file à Arabelle."},
				ApresVaincu: []string{"Un jour je te battrai ! ...enfin, peut-être."}},
			{X: 22, Y: 13, Glyph: "U", Nom: "Exploratrice Lya", Portrait: "exploratrice", Butin: 90,
				Equipe:      []string{"Ponchiot", "Ratentif"},
				Dialogue:    []string{"La Neo-Plasma rôde jusqu'ici. Montre-moi que tu tiens le coup."},
				Apres:       []string{"Bien joué. Arabelle est au nord. Méfie-toi des capes grises."},
				ApresVaincu: []string{"La Neo-Plasma se rapproche du repaire. Sois prête... euh, prêt."}},
			{X: 27, Y: 5, Glyph: "K", Nom: "Dresseur Kenji", Portrait: "dresseur", Butin: 70,
				Equipe:      []string{"Poichigeon"},
				Dialogue:    []string{"Personne ne passe sans un petit duel !"},
				Apres:       []string{"Pas mal du tout."},
				ApresVaincu: []string{"Reviens t'entraîner quand tu veux."}},
		},
		Sorties: []*Sortie{
			{X: 26, Y: 14, ZoneCible: "renouet", SpawnX: 26, SpawnY: 2, ChapitreMin: 1},
			{X: 26, Y: 1, ZoneCible: "arabelle", SpawnX: 26, SpawnY: 12, ChapitreMin: 1},
		},
	}
}

func zoneArabelle() *Zone {
	g := salle(52, 15)
	stamp(g, 3, 2, "^^^                                       ^^^")
	return &Zone{
		Cle: "arabelle", Nom: "ARABELLE", Banniere: "La ville de l'harmonie",
		Grille: g, Depart: [2]int{26, 12}, OnEnter: "ch2_grunt",
		PNJs: []*PNJ{
			{X: 9, Y: 4, Glyph: "C", Nom: "Centre Pokémon", Portrait: "infirmiere", Service: "centre"},
			{X: 9, Y: 8, Glyph: "B", Nom: "Boutique", Service: "boutique"},
			{X: 9, Y: 11, Glyph: "S", Nom: "Forge", Service: "forge"},
			{X: 34, Y: 6, Glyph: "G", Nom: "Maître d'Arène Cédric", Portrait: "cedric", Theme: "arene", Butin: 300, MonteChapitre: 3,
				Equipe: []string{"Ponchiot", "Munna", "Golette"},
				Dialogue: []string{
					"C'est toi qui as chassé ce sbire ? Voyons ton niveau.",
					"Arène d'Arabelle : trois Pokémon, aucun quartier !",
				},
				Apres: []string{
					"Superbe ! Tu as l'étoffe d'un héros.",
					"La Neo-Plasma s'est retranchée derrière la Forêt d'Ombreflore. Va.",
				}},
			{X: 40, Y: 10, Glyph: "A", Nom: "Habitante", Portrait: "habitant", Dialogue: []string{
				"On raconte que N lui-même est revenu à Unys...",
			}},
		},
		Sorties: []*Sortie{
			{X: 26, Y: 13, ZoneCible: "route1", SpawnX: 26, SpawnY: 2, ChapitreMin: 1},
			{X: 50, Y: 7, ZoneCible: "foret", SpawnX: 26, SpawnY: 13, ChapitreMin: 3,
				MsgBloque: "Des sbires de la Neo-Plasma bloquent la route. Bats le Maître d'Arène pour ouvrir la voie."},
		},
	}
}

func zoneForet() *Zone {
	g := salle(52, 16)
	stamp(g, 1, 1, strings.Repeat("^", 50))
	stamp(g, 1, 14, strings.Repeat("^", 50))
	stamp(g, 24, 1, "      ")  // passage nord
	stamp(g, 24, 14, "      ") // passage sud
	herbe(g, 6, 3, 10, 3)
	herbe(g, 34, 3, 10, 3)
	herbe(g, 6, 9, 10, 3)
	herbe(g, 34, 9, 10, 3)
	stamp(g, 4, 6, "^^^^^^^^^^^^^^") // barrières d'arbres, gap au centre
	stamp(g, 34, 6, "^^^^^^^^^^^^^^")
	return &Zone{
		Cle: "foret", Nom: "FORÊT D'OMBREFLORE", Banniere: "Sombre, dense, surveillée",
		Grille: g, Depart: [2]int{26, 13},
		Sauvages: []string{"Zorua", "Chacripan", "Munna", "Cheniti"},
		PNJs: []*PNJ{
			{X: 24, Y: 4, Glyph: "V", Nom: "Sbire Neo-Plasma", Portrait: "team_plasme", Theme: "plasma", Butin: 120,
				Equipe:      []string{"Chacripan", "Cheniti", "Zorua"},
				Dialogue:    []string{"Tu ne trouveras jamais la sortie, gamin."},
				Apres:       []string{"L'As t'attend plus loin..."},
				ApresVaincu: []string{"La Neo-Plasma ne pardonnera pas cet affront."}},
			{X: 26, Y: 8, Glyph: "W", Nom: "As de la Neo-Plasma", Portrait: "team_plasme", Theme: "plasma",
				Butin: 500, MonteChapitre: 4,
				Equipe: []string{"Chacripan", "Munna", "Zoroark"},
				Dialogue: []string{
					"Je suis l'As de la Neo-Plasma. Derrière moi : le repaire.",
					"Mon Zoroark va te montrer ce qu'est une véritable illusion.",
				},
				Apres: []string{
					"Comment... ? Personne n'avait battu mon Zoroark.",
					"",
					"N : Impressionnant. Le chef t'attend. Je vous rejoindrai.",
				}},
		},
		Sorties: []*Sortie{
			{X: 26, Y: 14, ZoneCible: "arabelle", SpawnX: 50, SpawnY: 7, ChapitreMin: 1},
			{X: 26, Y: 1, ZoneCible: "qg", SpawnX: 26, SpawnY: 11, ChapitreMin: 4,
				MsgBloque: "L'entrée du repaire est verrouillée. Bats l'As de la Neo-Plasma dans la forêt."},
		},
	}
}

func zoneQG() *Zone {
	g := salle(52, 14)
	stamp(g, 1, 1, strings.Repeat("=", 50))
	return &Zone{
		Cle: "qg", Nom: "QG DE LA NEO-PLASMA", Banniere: "Le repaire du mensonge",
		Grille: g, Depart: [2]int{26, 11},
		PNJs: []*PNJ{
			{X: 12, Y: 4, Glyph: "X", Nom: "Sbire", Portrait: "team_plasme", Theme: "plasma", Butin: 100,
				Equipe:   []string{"Nodulithe", "Chacripan"},
				Dialogue: []string{"Le chef ne recevra pas d'intrus !"}},
			{X: 33, Y: 4, Glyph: "Y", Nom: "Sbire", Portrait: "team_plasme", Theme: "plasma", Butin: 100,
				Equipe:   []string{"Golette", "Cheniti"},
				Dialogue: []string{"Demi-tour, dresseur."}},
			{X: 26, Y: 6, Glyph: "Z", Nom: "Nikolai (Chef de la Neo-Plasma)", Portrait: "plasma_leader", Theme: "plasma",
				Butin: 800, MonteChapitre: 5,
				Equipe: []string{"Golette", "Nodulithe", "Zoroark"},
				Dialogue: []string{
					"Je suis Nikolai. Je dirige la Neo-Plasma, et notre cause est juste.",
					"Les dresseurs sont des geôliers. Nous, nous ouvrons les cages.",
					"Assez parlé. Que le plus convaincu l'emporte !",
				},
				Apres: []string{
					"Nous n'étions... qu'un instrument. Il t'attend. Lui.",
				}},
			{X: 33, Y: 8, Glyph: "N", Nom: "N", Portrait: "N", Declenche: "final_n"},
		},
		Sorties: []*Sortie{
			{X: 26, Y: 12, ZoneCible: "foret", SpawnX: 26, SpawnY: 2, ChapitreMin: 1},
		},
	}
}
