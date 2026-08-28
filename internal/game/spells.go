package game

import "fmt"

// Sort décrit une attaque offensive utilisable en combat.
type Sort struct {
	Nom      string
	Type     string
	Degats   int
	CoutMana int
}

// sortsDisponibles associe chaque nom d'attaque à ses caractéristiques.
var sortsDisponibles = map[string]Sort{
	"Charge":        {Nom: "Charge", Type: "Normal", Degats: 8, CoutMana: 4},
	"Lance-Flammes": {Nom: "Lance-Flammes", Type: "Feu", Degats: 18, CoutMana: 10}, // = "Boule de feu" du sujet
	// Plante
	"Fouet Lianes": {Nom: "Fouet Lianes", Type: "Plante", Degats: 14, CoutMana: 5},
	"Tranch'Herbe": {Nom: "Tranch'Herbe", Type: "Plante", Degats: 22, CoutMana: 9},
	"Lance-Soleil": {Nom: "Lance-Soleil", Type: "Plante", Degats: 40, CoutMana: 20},
	// Feu
	"Flammèche":    {Nom: "Flammèche", Type: "Feu", Degats: 14, CoutMana: 5},
	"Roue de Feu":  {Nom: "Roue de Feu", Type: "Feu", Degats: 22, CoutMana: 9},
	"Déflagration": {Nom: "Déflagration", Type: "Feu", Degats: 40, CoutMana: 20},
	// Eau
	"Pistolet à O": {Nom: "Pistolet à O", Type: "Eau", Degats: 14, CoutMana: 5},
	"Bulles d'O":   {Nom: "Bulles d'O", Type: "Eau", Degats: 22, CoutMana: 9},
	"Hydrocanon":   {Nom: "Hydrocanon", Type: "Eau", Degats: 40, CoutMana: 20},
	// Normal / générique
	"Vive-Attaque": {Nom: "Vive-Attaque", Type: "Normal", Degats: 12, CoutMana: 3},
	"Écras'Face":   {Nom: "Écras'Face", Type: "Normal", Degats: 22, CoutMana: 9},
	"Ultimapoing":  {Nom: "Ultimapoing", Type: "Normal", Degats: 34, CoutMana: 16},
	// Ténèbres
	"Morsure":      {Nom: "Morsure", Type: "Ténèbres", Degats: 14, CoutMana: 5},
	"Coup Bas":     {Nom: "Coup Bas", Type: "Ténèbres", Degats: 22, CoutMana: 9},
	"Tranche-Nuit": {Nom: "Tranche-Nuit", Type: "Ténèbres", Degats: 40, CoutMana: 20},
	// Psy
	"Choc Mental": {Nom: "Choc Mental", Type: "Psy", Degats: 14, CoutMana: 5},
	"Vague Psy":   {Nom: "Vague Psy", Type: "Psy", Degats: 22, CoutMana: 9},
	"Psyko":       {Nom: "Psyko", Type: "Psy", Degats: 40, CoutMana: 20},
	// Insecte
	"Piqûre":    {Nom: "Piqûre", Type: "Insecte", Degats: 14, CoutMana: 5},
	"Dard-Nuée": {Nom: "Dard-Nuée", Type: "Insecte", Degats: 22, CoutMana: 9},
	"Bourdon":   {Nom: "Bourdon", Type: "Insecte", Degats: 40, CoutMana: 20},
	// Sol / Roche
	"Jet-Pierres": {Nom: "Jet-Pierres", Type: "Sol", Degats: 14, CoutMana: 5},
	"Éboulement":  {Nom: "Éboulement", Type: "Sol", Degats: 22, CoutMana: 9},
	"Séisme":      {Nom: "Séisme", Type: "Sol", Degats: 40, CoutMana: 20},
	// Électrik (Zekrom)
	"Éclair":       {Nom: "Éclair", Type: "Électrik", Degats: 16, CoutMana: 5},
	"Tonnerre":     {Nom: "Tonnerre", Type: "Électrik", Degats: 26, CoutMana: 10},
	"Éclair Croix": {Nom: "Éclair Croix", Type: "Électrik", Degats: 48, CoutMana: 22},
	// Glace (Kyurem)
	"Éclats Glace": {Nom: "Éclats Glace", Type: "Glace", Degats: 16, CoutMana: 5},
	"Laser Glace":  {Nom: "Laser Glace", Type: "Glace", Degats: 26, CoutMana: 10},
	"Blizzard":     {Nom: "Blizzard", Type: "Glace", Degats: 48, CoutMana: 22},
}

type palierSort struct {
	Niv int
	Nom string
}

// lignesAttaque : progression d'attaques selon le TYPE du Pokémon.
var lignesAttaque = map[string][]palierSort{
	"Plante":   {{2, "Fouet Lianes"}, {5, "Tranch'Herbe"}, {8, "Lance-Soleil"}},
	"Feu":      {{2, "Flammèche"}, {5, "Roue de Feu"}, {8, "Déflagration"}},
	"Eau":      {{2, "Pistolet à O"}, {5, "Bulles d'O"}, {8, "Hydrocanon"}},
	"Normal":   {{2, "Vive-Attaque"}, {5, "Écras'Face"}, {8, "Ultimapoing"}},
	"Ténèbres": {{2, "Morsure"}, {5, "Coup Bas"}, {8, "Tranche-Nuit"}},
	"Psy":      {{2, "Choc Mental"}, {5, "Vague Psy"}, {8, "Psyko"}},
	"Insecte":  {{2, "Piqûre"}, {5, "Dard-Nuée"}, {8, "Bourdon"}},
	"Sol":      {{2, "Jet-Pierres"}, {5, "Éboulement"}, {8, "Séisme"}},
	"Électrik": {{2, "Éclair"}, {12, "Tonnerre"}, {22, "Éclair Croix"}}, // Zekrom
	"Glace":    {{2, "Éclats Glace"}, {12, "Laser Glace"}, {22, "Blizzard"}},
}

// ligneGenerique : pour un type sans progression dédiée.
var ligneGenerique = []palierSort{{2, "Vive-Attaque"}, {5, "Écras'Face"}, {8, "Ultimapoing"}}

func lignePourType(typ string) []palierSort {
	if l, ok := lignesAttaque[typ]; ok {
		return l
	}
	return ligneGenerique
}

// mettreAJourSorts : donne au Pokémon ACTIF (représenté par p) toutes les
// attaques que son niveau et son type débloquent.
func mettreAJourSorts(p *Personnage) {
	for _, pa := range lignePourType(p.Type) {
		if p.Niveau >= pa.Niv && !connaitSort(p, pa.Nom) {
			p.Skills = append(p.Skills, pa.Nom)
			fmt.Printf(">>> %s apprend une nouvelle attaque : %s !\n", p.Starter, pa.Nom)
		}
	}
}

// mettreAJourSortsPokemon : idem pour un membre de l'équipe (non actif).
func mettreAJourSortsPokemon(pk *Pokemon) {
	if len(pk.Skills) == 0 {
		pk.Skills = []string{"Charge"}
	}
	for _, pa := range lignePourType(pk.Type) {
		if pk.Niveau >= pa.Niv && !contientChaine(pk.Skills, pa.Nom) {
			pk.Skills = append(pk.Skills, pa.Nom)
			fmt.Printf(">>> %s apprend une nouvelle attaque : %s !\n", pk.Nom, pa.Nom)
		}
	}
}

func contientChaine(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func connaitSort(p *Personnage, nom string) bool { return contientChaine(p.Skills, nom) }

// spellBook ajoute le sort "Lance-Flammes" (la "Boule de feu" du sujet) à la liste
// de sorts du personnage. Un même sort ne peut être appris qu'une seule fois.
// Renvoie true si le sort a effectivement été appris (sinon le livre n'est pas consommé).
func spellBook(p *Personnage) bool {
	if connaitSort(p, "Lance-Flammes") {
		fmt.Println("Ce sort est déjà appris : tu connais déjà Lance-Flammes !")
		return false
	}
	p.Skills = append(p.Skills, "Lance-Flammes")
	fmt.Println("Tu as appris l'attaque Lance-Flammes !")
	return true
}

// degatsAttaqueRapide : l'attaque de base vaut 5 (Tâche 19), +2 par niveau
// au-delà (gain de stats compté en bonus, Mission 4).
func degatsAttaqueRapide(p *Personnage) int {
	return 5 + 2*(p.Niveau-1)
}

// choisirAttaque affiche le menu des attaques (attaque de base + sorts connus),
// applique celle choisie, et renvoie true si une attaque a bien été lancée
// (le tour est alors consommé).
func choisirAttaque(p *Personnage, m *Monster) bool {
	base := degatsAttaqueRapide(p)
	for {
		fmt.Println("Quelle attaque ? (dégâts affichés = après effet de type)")
		fmt.Printf(" 1. %-16s [%s] %s\n", "Attaque Rapide", "Normal", ligneAttaque(base, "Normal", m.Type, 0))
		for i, s := range p.Skills {
			sc := sortsDisponibles[s]
			fmt.Printf(" %d. %-16s [%s] %s\n", i+2, sc.Nom, sc.Type, ligneAttaque(sc.Degats, sc.Type, m.Type, sc.CoutMana))
		}
		fmt.Println(" R. Retour")

		choix := lireLigne("> ")
		if choix == "R" || choix == "r" {
			return false
		}
		n := indexFromChoice(choix, len(p.Skills)+1)
		if n == -1 {
			fmt.Println("Choix invalide.")
			continue
		}
		if n == 0 { // Attaque Rapide (type Normal, neutre)
			deg, label := appliquerType(base, "Normal", m.Type)
			infligerAuMonstre(m, deg)
			TypeText(fmt.Sprintf("%s utilise Attaque Rapide et inflige %d dégâts à %s !", p.Starter, deg, m.Nom))
			if label != "" {
				TypeText(label)
			}
			fmt.Printf("%s - PV : %d/%d\n", m.Nom, m.PV, m.PVMax)
			return true
		}
		sc := sortsDisponibles[p.Skills[n-1]]
		if p.Mana < sc.CoutMana {
			fmt.Println("Mana insuffisant pour lancer " + sc.Nom + " ! Choisis autre chose.")
			continue
		}
		p.Mana -= sc.CoutMana
		deg, label := appliquerType(sc.Degats, sc.Type, m.Type)
		infligerAuMonstre(m, deg)
		TypeText(fmt.Sprintf("%s lance %s [%s] et inflige %d dégâts à %s !", p.Starter, sc.Nom, sc.Type, deg, m.Nom))
		if label != "" {
			TypeText(label)
		}
		fmt.Printf("%s - PV : %d/%d\n", m.Nom, m.PV, m.PVMax)
		fmt.Printf("Mana restant : %d/%d\n", p.Mana, p.ManaMax)
		return true
	}
}

// ligneAttaque formate la partie droite d'une entrée du menu d'attaque :
// dégâts réels (après type) + coût mana + mention super/peu efficace.
func ligneAttaque(base int, typeAtt, typeDef string, cout int) string {
	deg, label := appliquerType(base, typeAtt, typeDef)
	s := fmt.Sprintf("%d dégâts", deg)
	if deg != base {
		s = fmt.Sprintf("%d dégâts (base %d)", deg, base)
	}
	s += fmt.Sprintf(", %d mana", cout)
	switch {
	case label != "" && deg > base:
		s += col("  ★ super efficace", cVert)
	case label != "" && deg < base:
		s += col("  ✗ peu efficace", cRouge)
	}
	return s
}

// infligerAuMonstre applique des dégâts au monstre (borné à 0).
func infligerAuMonstre(m *Monster, degats int) {
	m.PV -= degats
	if m.PV < 0 {
		m.PV = 0
	}
}
