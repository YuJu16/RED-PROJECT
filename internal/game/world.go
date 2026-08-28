package game

import (
	"fmt"
	"math/rand"
	"strings"
)

// ---------------------------------------------------------------------------
// Moteur de monde : plusieurs zones reliées, explorées aux flèches (ou ZQSD).
// ---------------------------------------------------------------------------

// PNJ est un personnage non-joueur posé sur une carte. Selon ses champs il peut
// tenir un service (labo, centre, boutique, forge), déclencher une scène
// d'histoire, ou proposer un combat de dresseur.
type PNJ struct {
	X, Y        int
	Glyph       string
	Nom         string
	Portrait    string   // dessin ascii/<Portrait>.txt affiché au contact (optionnel)
	Theme       string   // musique jouée au contact / avant le combat (ex. "plasma")
	Dialogue    []string // lignes dites au contact (ou avant le combat)
	Apres       []string // lignes dites juste après avoir été battu
	ApresVaincu []string // lignes dites si on lui reparle plus tard (optionnel)

	Service       string   // "labo", "centre", "boutique", "forge"
	Declenche     string   // id de scène d'histoire jouée une seule fois
	DeclencheFait bool     // la scène a déjà été jouée
	Equipe        []string // espèces envoyées si c'est un dresseur
	Butin         int      // Poké-Dollars gagnés en le battant
	MonteChapitre int      // fait passer p.Chapitre à cette valeur après victoire
	Vaincu        bool
}

// Sortie relie une case d'une zone à une autre zone.
type Sortie struct {
	X, Y           int
	ZoneCible      string
	SpawnX, SpawnY int
	ChapitreMin    int    // le joueur doit avoir atteint ce chapitre
	MsgBloque      string // message affiché si le passage est verrouillé
}

// Zone est un écran de carte.
type Zone struct {
	Cle         string
	Nom         string
	Banniere    string
	Grille      []string
	Depart      [2]int
	PNJs        []*PNJ
	Sorties     []*Sortie
	Sauvages    []string // espèces croisées dans les hautes herbes ','
	OnEnter     string   // scène jouée la 1re fois qu'on entre dans la zone
	onEnterFait bool
}

var zones = map[string]*Zone{}

// initMonde construit toutes les zones. Appelé une fois au démarrage.
func initMonde() {
	zones["renouet"] = zoneRenouet()
	zones["route1"] = zoneRoute1()
	zones["arabelle"] = zoneArabelle()
	zones["foret"] = zoneForet()
	zones["qg"] = zoneQG()
}

// ---------------------------------------------------------------------------
// Boucle d'exploration
// ---------------------------------------------------------------------------

func explorerMonde(p *Personnage) {
	if p.ZoneActuelle == "" {
		p.ZoneActuelle = "renouet"
	}
	z := zones[p.ZoneActuelle]
	if z == nil {
		z = zones["renouet"]
		p.ZoneActuelle = "renouet"
	}
	x, y := p.PosX, p.PosY
	if x <= 0 || y <= 0 {
		x, y = z.Depart[0], z.Depart[1]
	}

	for {
		musique(z.Cle) // (re)lance la musique de la zone, ex. après un combat
		ClearScreen()
		afficherZone(z, x, y, p)

		nx, ny := x, y
		switch ReadKey() {
		case "Up", "z", "Z":
			ny--
		case "Down", "s", "S":
			ny++
		case "Left", "q", "Q":
			nx--
		case "Right", "d", "D":
			nx++
		case "i", "I", "p", "P":
			menuPause(p)
			continue
		case "r", "R", "Quit":
			p.PosX, p.PosY = x, y
			ClearScreen()
			if r := strings.ToLower(lireLigne("Sauvegarder la partie avant de revenir au menu ? [O/n] ")); r == "" || r == "o" || r == "oui" {
				if err := sauvegarder(p); err == nil {
					fmt.Println("Partie sauvegardée.")
				} else {
					fmt.Println("Échec de la sauvegarde : " + err.Error())
				}
				lireLigne("\nAppuyez sur Entrée...")
			}
			return
		default:
			continue
		}

		if ny < 0 || ny >= len(z.Grille) || nx < 0 || nx >= len([]rune(z.Grille[ny])) {
			continue
		}

		// Un PNJ occupe la case ciblée -> on lui parle, on ne bouge pas.
		if n := pnjA(z, nx, ny); n != nil {
			interagirPNJ(p, z, n)
			if zc := zones[p.ZoneActuelle]; zc != z { // une scène a pu nous déplacer
				z, x, y = zc, p.PosX, p.PosY
			}
			continue
		}

		// Une sortie occupe la case ciblée -> changement de zone.
		if s := sortieA(z, nx, ny); s != nil {
			if p.Chapitre < s.ChapitreMin {
				ClearScreen()
				msg := s.MsgBloque
				if msg == "" {
					msg = "Le passage est bloqué pour le moment."
				}
				TypeText(msg)
				lireLigne("\nAppuyez sur Entrée...")
				continue
			}
			p.ZoneActuelle = s.ZoneCible
			z = zones[s.ZoneCible]
			x, y = s.SpawnX, s.SpawnY
			p.PosX, p.PosY = x, y
			banniereZone(z)
			if z.OnEnter != "" && !z.onEnterFait {
				if jouerScene(p, z.OnEnter) {
					z.onEnterFait = true
				}
			}
			continue
		}

		if !marchable(z, nx, ny) {
			continue
		}

		x, y = nx, ny
		p.PosX, p.PosY = x, y

		if caseRune(z, x, y) == ',' {
			tenterRencontre(p, z)
		}
	}
}

// ---------------------------------------------------------------------------
// Rendu
// ---------------------------------------------------------------------------

func afficherZone(z *Zone, px, py int, p *Personnage) {
	fmt.Println(titre(z.Nom))
	if z.Banniere != "" {
		fmt.Println(z.Banniere)
	}
	fmt.Printf("%s %s\n", col("Objectif :", cJaune), objectifChapitre(p.Chapitre))
	fmt.Println("Flèches ou Z/Q/S/D : bouger   |   I : sac & équipe   |   R : menu principal")
	fmt.Println()

	grille := make([][]rune, len(z.Grille))
	for i, ligne := range z.Grille {
		grille[i] = []rune(ligne)
	}
	place := func(x, y int, r rune) {
		if y >= 0 && y < len(grille) && x >= 0 && x < len(grille[y]) {
			grille[y][x] = r
		}
	}
	for _, s := range z.Sorties {
		place(s.X, s.Y, '>')
	}
	for _, n := range z.PNJs {
		g := 'o'
		if n.Glyph != "" {
			g = []rune(n.Glyph)[0]
		}
		place(n.X, n.Y, g)
	}
	place(px, py, '@')

	for _, ligne := range grille {
		s := string(ligne)
		if strings.ContainsRune(s, '@') {
			s = strings.Replace(s, "@", col("@", cVert+cGras), 1)
		}
		fmt.Println(s)
	}
	fmt.Println()
	fmt.Println("Légende : " + col("@", cVert) + " vous   " + col(">", cCyan) + " sortie   , hautes herbes   lettres = personnes")
}

// menuPause : écran de pause sur la carte (état de l'équipe + sac).
func menuPause(p *Personnage) {
	for {
		ClearScreen()
		fmt.Println(titre("Pause"))
		afficherEtatEquipe(p)
		fmt.Printf("\nMana : %d/%d    Or : %d\n\n", p.Mana, p.ManaMax, p.Or)
		fmt.Println("1. Utiliser un objet du sac")
		fmt.Println("2. Changer de Pokémon actif")
		fmt.Println("R. Reprendre l'exploration")
		switch lireLigne("> ") {
		case "1":
			accessInventory(p)
		case "2":
			syncPersoVersActif(p)
			if changerPokemon(p, "") {
				chargerActifVersPerso(p)
			}
			lireLigne("\nAppuyez sur Entrée...")
		case "r", "R":
			return
		}
	}
}

func banniereZone(z *Zone) {
	ClearScreen()
	fmt.Println("========================================")
	fmt.Printf("   %s\n", z.Nom)
	if z.Banniere != "" {
		fmt.Printf("   %s\n", z.Banniere)
	}
	fmt.Println("========================================")
	lireLigne("\nAppuyez sur Entrée...")
}

// ---------------------------------------------------------------------------
// Helpers de carte
// ---------------------------------------------------------------------------

func caseRune(z *Zone, x, y int) rune {
	if y < 0 || y >= len(z.Grille) {
		return '#'
	}
	r := []rune(z.Grille[y])
	if x < 0 || x >= len(r) {
		return '#'
	}
	return r[x]
}

// marchable : on ne marche que sur le sol ' ' et les hautes herbes ','.
func marchable(z *Zone, x, y int) bool {
	switch caseRune(z, x, y) {
	case ' ', ',':
		return true
	}
	return false
}

func pnjA(z *Zone, x, y int) *PNJ {
	for _, n := range z.PNJs {
		if n.X == x && n.Y == y {
			return n
		}
	}
	return nil
}

func sortieA(z *Zone, x, y int) *Sortie {
	for _, s := range z.Sorties {
		if s.X == x && s.Y == y {
			return s
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Interaction avec un PNJ
// ---------------------------------------------------------------------------

// flagOK / flagSet : petits helpers sur p.Flags (progression persistée).
func flagOK(p *Personnage, k string) bool { return p.Flags != nil && p.Flags[k] }
func flagSet(p *Personnage, k string) {
	if p.Flags == nil {
		p.Flags = map[string]bool{}
	}
	p.Flags[k] = true
}

// idDresseur : clé unique d'un dresseur (zone + glyphe).
func idDresseur(z *Zone, n *PNJ) string { return "vaincu:" + z.Cle + ":" + n.Glyph }

// appliquerProgression réapplique aux zones (reconstruites à chaque lancement)
// l'état sauvegardé : dresseurs déjà battus.
func appliquerProgression(p *Personnage) {
	// Toute espèce déjà dans l'équipe compte comme enregistrée au Pokédex.
	for _, e := range p.Equipe {
		flagSet(p, "capture:"+e.Nom)
	}
	for _, z := range zones {
		for _, n := range z.PNJs {
			if len(n.Equipe) == 0 {
				continue
			}
			// dresseur explicitement mémorisé, OU dresseur qui débloque un
			// chapitre déjà dépassé (rattrapage pour les sauvegardes d'avant ce système).
			if flagOK(p, idDresseur(z, n)) || (n.MonteChapitre > 0 && p.Chapitre >= n.MonteChapitre) {
				n.Vaincu = true
				// un dresseur battu qui gardait un chapitre le débloque.
				if n.MonteChapitre > p.Chapitre {
					p.Chapitre = n.MonteChapitre
				}
			}
		}
	}
}

func interagirPNJ(p *Personnage, z *Zone, n *PNJ) {
	// Scène d'histoire (une seule fois — sauf si la scène refuse de se jouer).
	if n.Declenche != "" && !n.DeclencheFait {
		if jouerScene(p, n.Declenche) {
			n.DeclencheFait = true
		}
		return
	}

	// Portrait du PNJ, affiché UNE seule fois ici pour tout le reste de l'échange.
	// (Le Centre et le Labo affichent déjà leur propre dessin.)
	if n.Portrait != "" && n.Service != "centre" && n.Service != "labo" {
		ClearScreen()
		afficherAscii(n.Portrait)
	}

	switch n.Service {
	case "maison":
		jouerDialogue(n.Nom, n.Dialogue)
		soigner(p)
		fmt.Printf("%s : PV %d/%d - Mana %d/%d. Toute l'équipe est soignée.\n", p.Nom, p.PV, p.PVMax, p.Mana, p.ManaMax)
		lireLigne("\nAppuyez sur Entrée...")
		return
	case "labo":
		sceneLabo(p)
		return
	case "centre":
		centrePokemon(p)
		return
	case "boutique":
		merchantMenu(p)
		return
	case "forge":
		forgeronMenu(p)
		return
	}

	// Dresseur pas encore battu.
	if len(n.Equipe) > 0 && !n.Vaincu {
		if n.Theme != "" {
			musique(n.Theme)
		}
		jouerDialogue(n.Nom, n.Dialogue)
		if combatDresseur(p, n) {
			n.Vaincu = true
			flagSet(p, idDresseur(z, n)) // mémorisé pour les prochains lancements
			if len(n.Apres) > 0 {
				jouerDialogue(n.Nom, n.Apres)
			}
			if n.Butin > 0 {
				p.Or += n.Butin
				fmt.Printf("Vous recevez %d Poké-Dollars !\n", n.Butin)
			}
			if n.MonteChapitre > p.Chapitre {
				p.Chapitre = n.MonteChapitre
				fmt.Printf("\n>>> Nouvel objectif : %s\n", objectifChapitre(p.Chapitre))
			}
			lireLigne("\nAppuyez sur Entrée...")
		}
		return
	}

	// PNJ bavard. Un dresseur déjà battu répète sa réplique d'après-combat
	// (ou une phrase générique) plutôt que son défi. (Le portrait est déjà affiché.)
	lignes := n.Dialogue
	if n.Vaincu {
		if len(n.ApresVaincu) > 0 {
			lignes = n.ApresVaincu
		} else {
			lignes = []string{"On s'est déjà battus. Bonne route, dresseur !"}
		}
	}
	jouerDialogue(n.Nom, lignes)
	lireLigne("\nAppuyez sur Entrée...")
}

// tenterRencontre déclenche parfois un combat sauvage dans les hautes herbes.
func tenterRencontre(p *Personnage, z *Zone) {
	if len(z.Sauvages) == 0 {
		return
	}
	if rand.Intn(100) >= 22 { // ~22 % de chance par pas dans l'herbe
		return
	}
	espece := z.Sauvages[rand.Intn(len(z.Sauvages))]
	m := nouveauSauvage(espece, p.Niveau)
	ClearScreen()
	TypeText("Les hautes herbes s'agitent...")
	combatSauvage(p, m)
	lireLigne("\nAppuyez sur Entrée pour continuer l'exploration...")
}
