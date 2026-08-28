package game

import "math/rand"

// Monster représente un adversaire rencontré en combat.
type Monster struct {
	Nom        string
	Type       string
	Niveau     int
	PVMax      int
	PV         int
	Attaque    int
	Initiative int
	Experience int
	Ascii      string
}

// Espece décrit les statistiques de base (niveau 1) d'une espèce de Pokémon.
type Espece struct {
	Nom        string
	Type       string // sert au moveset appris en montant de niveau (voir spells.go)
	PVBase     int
	AttBase    int
	Initiative int
	XP         int // expérience donnée à la défaite (niveau 1)
	Ascii      string
}

// especes est le bestiaire du jeu. Le fichier ascii/<Ascii>.txt fournit le dessin
// (un fichier absent ou vide affiche simplement un cadre "dessin à venir").
var especes = map[string]Espece{
	// Route 1
	"Ratentif":   {"Ratentif", "Normal", 40, 5, 4, 25, "ratentif"}, // = "Gobelin d'entraînement" du sujet
	"Poichigeon": {"Poichigeon", "Normal", 34, 6, 9, 22, "poichigeon"},
	"Ponchiot":   {"Ponchiot", "Normal", 42, 6, 6, 26, "ponchiot"},
	// Forêt d'Ombreflore
	"Zorua":     {"Zorua", "Ténèbres", 55, 7, 8, 40, "zorua"},
	"Chacripan": {"Chacripan", "Ténèbres", 48, 8, 10, 38, "chacripan"},
	"Munna":     {"Munna", "Psy", 60, 5, 3, 44, "munna"},
	"Cheniti":   {"Cheniti", "Insecte", 45, 5, 5, 30, "cheniti"},
	// QG Neo-Plasma / grottes
	"Golette":   {"Golette", "Sol", 70, 6, 2, 55, "golette"},
	"Nodulithe": {"Nodulithe", "Sol", 66, 7, 3, 48, "nodulithe"},
	// Boss / as de dresseurs
	"Zoroark": {"Zoroark", "Ténèbres", 120, 14, 12, 220, "zoroark"},
	"Zekrom":  {"Zekrom", "Électrik", 210, 16, 8, 600, "zekrom"},
	"Kyurem":  {"Kyurem", "Glace", 300, 22, 9, 1200, "kyurem"},
}

// nouveauMonstre crée un Monster niveau 1 aux stats de base de l'espèce.
func nouveauMonstre(nom string) *Monster {
	e, ok := especes[nom]
	if !ok {
		e = especes["Ratentif"]
	}
	return &Monster{
		Nom:        e.Nom,
		Type:       e.Type,
		Niveau:     1,
		PVMax:      e.PVBase,
		PV:         e.PVBase,
		Attaque:    e.AttBase,
		Initiative: e.Initiative,
		Experience: e.XP,
		Ascii:      e.Ascii,
	}
}

// nouveauAuNiveau crée un Pokémon de l'espèce demandée mis à l'échelle du niveau.
// Progression volontairement douce (le jeu ne doit pas être punitif).
func nouveauAuNiveau(nom string, niveau int) *Monster {
	if niveau < 1 {
		niveau = 1
	}
	m := nouveauMonstre(nom)
	m.Niveau = niveau
	bonus := niveau - 1
	m.PVMax += bonus * 5
	m.PV = m.PVMax
	m.Attaque += bonus // +1 attaque par niveau
	m.Experience += bonus * 4
	return m
}

// niveauSauvageAleatoire : dans la nature, on croise des Pokémon de niveau varié,
// mais jamais au-dessus du joueur (0 à 3 niveaux en dessous).
func niveauSauvageAleatoire(niveauJoueur int) int {
	n := niveauJoueur - rand.Intn(4)
	if n < 1 {
		n = 1
	}
	return n
}

// nouveauSauvage : un Pokémon sauvage de niveau aléatoire adapté au joueur.
func nouveauSauvage(nom string, niveauJoueur int) *Monster {
	return nouveauAuNiveau(nom, niveauSauvageAleatoire(niveauJoueur))
}

// InitRatentif garde la signature du sujet (Tâche 16). C'est le "Gobelin d'entraînement".
func InitRatentif() *Monster { return nouveauMonstre("Ratentif") }

// InitZorua initialise un Zorua sauvage, adversaire vif et fragile.
func InitZorua() *Monster { return nouveauMonstre("Zorua") }

// InitGolette initialise un Golette sauvage, résistant mais lent.
func InitGolette() *Monster { return nouveauMonstre("Golette") }
