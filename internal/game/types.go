package game

// ---------------------------------------------------------------------------
// Table des types (Bonus 1, suggéré par la fiche).
// Multiplicateur appliqué aux dégâts : 2 = super efficace, 0.5 = peu efficace,
// 1 = neutre (les dégâts de base). Volontairement simple : pas d'immunités.
// Types présents dans le jeu : Normal, Plante, Feu, Eau, Ténèbres, Psy, Insecte, Sol.
// ---------------------------------------------------------------------------

// forces[attaque] = liste des types sur lesquels cette attaque est SUPER efficace.
var forces = map[string][]string{
	"Feu":      {"Plante", "Insecte"},
	"Eau":      {"Feu", "Sol"},
	"Plante":   {"Eau", "Sol"},
	"Ténèbres": {"Psy"},
	"Psy":      {"Insecte"}, // léger clin d'œil : le Psy craint l'Insecte mais domine ici Insecte pour rester lisible
	"Insecte":  {"Plante", "Ténèbres"},
	"Sol":      {"Feu", "Électrik"},
	"Électrik": {"Eau"},
	"Glace":    {"Plante", "Sol"},
}

// faiblesses[attaque] = types sur lesquels cette attaque est PEU efficace.
var faiblesses = map[string][]string{
	"Feu":      {"Eau", "Feu"},
	"Eau":      {"Plante", "Eau"},
	"Plante":   {"Feu", "Plante", "Insecte"},
	"Ténèbres": {"Ténèbres"},
	"Psy":      {"Psy", "Ténèbres"},
	"Insecte":  {"Feu"},
	"Sol":      {"Plante", "Électrik"},
	"Électrik": {"Plante", "Sol", "Électrik"},
	"Glace":    {"Feu", "Eau", "Glace"},
}

func contientType(l []string, t string) bool {
	for _, x := range l {
		if x == t {
			return true
		}
	}
	return false
}

// efficacite renvoie le multiplicateur de dégâts et un court libellé à afficher.
func efficacite(typeAttaque, typeDefenseur string) (float64, string) {
	if typeAttaque == "" || typeDefenseur == "" || typeAttaque == "Normal" {
		return 1.0, ""
	}
	if contientType(forces[typeAttaque], typeDefenseur) {
		return 2.0, "C'est super efficace !"
	}
	if contientType(faiblesses[typeAttaque], typeDefenseur) {
		return 0.5, "Ce n'est pas très efficace..."
	}
	return 1.0, ""
}

// appliquerType ajuste des dégâts selon les types et renvoie (dégâts, libellé).
func appliquerType(degats int, typeAttaque, typeDefenseur string) (int, string) {
	mult, label := efficacite(typeAttaque, typeDefenseur)
	d := int(float64(degats)*mult + 0.5)
	if d < 1 {
		d = 1
	}
	return d, label
}
