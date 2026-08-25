package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Personnage représente le dresseur contrôlé par le joueur.
type Personnage struct {
	Nom               string
	Genre             string
	Starter           string
	Type              string // type élémentaire du starter, sert de "classe" (impacte les PV)
	Niveau            int
	PVMax             int
	PV                int
	ManaMax           int
	Mana              int
	Experience        int
	ExperienceMax     int
	Or                int
	Initiative        int
	Inventaire        []Objet
	InventaireMax     int
	InventoryUpgrades int
	Skills            []string
	Equipement        Equipment
}

var reader = bufio.NewReader(os.Stdin)

// lireLigne affiche un prompt et lit une ligne saisie par le joueur.
// Si l'entrée standard est fermée (EOF), le programme s'arrête proprement.
func lireLigne(prompt string) string {
	fmt.Print(prompt)
	ligne, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\nEntrée fermée, fin du programme.")
		os.Exit(0)
	}
	return strings.TrimSpace(ligne)
}

// Init crée le personnage (dialogue de création) et initialise ses stats de départ.
func Init() *Personnage {
	p := charCreation()
	p.Niveau = 1
	p.Experience = 0
	p.ExperienceMax = 100
	p.Or = 100
	p.ManaMax = 20
	p.Mana = p.ManaMax
	p.Initiative = 5
	p.Skills = []string{"Charge"}
	p.InventaireMax = 10
	p.InventoryUpgrades = 0
	p.Inventaire = []Objet{
		{Nom: "Potion (Restaure 50 PV)", Quantite: 2, Type: "Potion"},
	}
	return p
}

// charCreation fait dialoguer le joueur avec le Professeur : nom, genre, starter.
func charCreation() *Personnage {
	ClearScreen()
	fmt.Println(getAscii("professor"))
	TypeText("Professeure Keteleeria : Bonjour ! Bienvenue dans le monde des Pokémon !")
	TypeText("Professeure Keteleeria : Je suis la Professeure Keteleeria. Ce monde est peuplé de créatures appelées Pokémon.")
	fmt.Println()

	nom := ""
	for nom == "" {
		nom = lireLigne("Comment t'appelles-tu ? ")
	}
	nom = formaterNom(nom)
	ClearScreen()

	genre := ""
	for genre != "1" && genre != "2" {
		TypeText("Es-tu un garçon ou une fille ?")
		fmt.Println("1. Garçon")
		fmt.Println("2. Fille")
		genre = lireLigne("> ")
	}
	ClearScreen()
	if genre == "1" {
		fmt.Println(getAscii("genre_M"))
		genre = "Garçon"
		TypeText("Tu es donc un garçon !")
	} else {
		fmt.Println(getAscii("genre_F"))
		genre = "Fille"
		TypeText("Tu es donc une fille !")
	}
	lireLigne("\nAppuyez sur Entrée pour continuer...")

	ClearScreen()
	fmt.Println(getAscii("professor"))
	TypeText(fmt.Sprintf("Professeure Keteleeria : Très bien, %s ! Il est temps de choisir ton premier partenaire :", nom))
	fmt.Println()
	fmt.Println("1. Vipélierre (Plante)")
	fmt.Println(getAscii("vipelierre"))
	fmt.Println("2. Gruikui (Feu)")
	fmt.Println(getAscii("gruikui"))
	fmt.Println("3. Moustillon (Eau)")
	fmt.Println(getAscii("moustillon"))

	var starter, typePoke string
	var pvMax int
	choix := ""
	for choix != "1" && choix != "2" && choix != "3" {
		choix = lireLigne("> ")
		switch choix {
		case "1":
			starter, typePoke, pvMax = "Vipélierre", "Plante", 120
		case "2":
			starter, typePoke, pvMax = "Gruikui", "Feu", 100
		case "3":
			starter, typePoke, pvMax = "Moustillon", "Eau", 110
		default:
			fmt.Println("Choix invalide, réessaie.")
		}
	}

	ClearScreen()
	fmt.Println(getAscii("professor"))
	TypeText(fmt.Sprintf("Professeure Keteleeria : %s t'accompagnera dans ton aventure !", starter))
	lireLigne("\nAppuyez sur Entrée pour continuer...")

	return &Personnage{
		Nom:     nom,
		Genre:   genre,
		Starter: starter,
		Type:    typePoke,
		PVMax:   pvMax,
		PV:      pvMax / 2, // PV actuels de départ = 50% des PV max (Mission 1)
	}
}

// formaterNom met une majuscule au premier caractère et le reste en minuscules.
func formaterNom(nom string) string {
	runes := []rune(nom)
	return strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
}

// displayInfo affiche la fiche personnage.
func displayInfo(p *Personnage) {
	fmt.Println("---------------------------")
	fmt.Printf("Nom        : %s (%s)\n", p.Nom, p.Genre)
	fmt.Printf("Starter    : %s [%s]\n", p.Starter, p.Type)
	fmt.Printf("Niveau     : %d (XP : %d/%d)\n", p.Niveau, p.Experience, p.ExperienceMax)
	fmt.Printf("PV         : %d/%d\n", p.PV, p.PVMax)
	fmt.Printf("Mana       : %d/%d\n", p.Mana, p.ManaMax)
	fmt.Printf("Or         : %d\n", p.Or)
	fmt.Printf("Initiative : %d\n", p.Initiative)
	fmt.Printf("Sorts      : %s\n", strings.Join(p.Skills, ", "))
	fmt.Printf("Équipement : Tête=%s | Torse=%s | Pieds=%s\n",
		venOuVide(p.Equipement.Tete), venOuVide(p.Equipement.Torse), venOuVide(p.Equipement.Pieds))
	fmt.Printf("Inventaire : %d/%d emplacements\n", len(p.Inventaire), p.InventaireMax)
	fmt.Println("---------------------------")
}

// venOuVide renvoie "-" si la chaîne est vide, sinon la chaîne elle-même.
func venOuVide(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// dead vérifie si le joueur est K.O. (0 PV) et le ressuscite avec 50% de ses PV max.
func dead(p *Personnage) {
	if p.PV <= 0 {
		fmt.Println(p.Nom + " est K.O. ! Transport en urgence vers le Centre Pokémon...")
		p.PV = p.PVMax / 2
		fmt.Printf("%s se réveille avec %d/%d PV.\n", p.Nom, p.PV, p.PVMax)
	}
}

// gainExperience ajoute de l'expérience au personnage et gère la montée de niveau,
// en reportant le surplus d'XP sur le niveau suivant.
func gainExperience(p *Personnage, xp int) {
	fmt.Printf("%s gagne %d points d'expérience !\n", p.Nom, xp)
	p.Experience += xp
	for p.Experience >= p.ExperienceMax {
		p.Experience -= p.ExperienceMax
		p.Niveau++
		p.ExperienceMax = int(float64(p.ExperienceMax) * 1.5)

		bonusPV := 10
		bonusMana := 5
		p.PVMax += bonusPV
		p.PV = p.PVMax
		p.ManaMax += bonusMana
		p.Mana = p.ManaMax

		fmt.Printf("Niveau supérieur ! %s passe niveau %d (+%d PV max, +%d Mana max, PV/Mana restaurés).\n",
			p.Nom, p.Niveau, bonusPV, bonusMana)
	}
	fmt.Printf("Expérience : %d/%d\n", p.Experience, p.ExperienceMax)
}

// indexFromChoice convertit une saisie utilisateur "1".."max" en index 0-based, -1 si invalide.
func indexFromChoice(choix string, max int) int {
	n := 0
	_, err := fmt.Sscanf(choix, "%d", &n)
	if err != nil || n < 1 || n > max {
		return -1
	}
	return n - 1
}
