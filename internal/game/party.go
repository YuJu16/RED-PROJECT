package game

import "fmt"

// ---------------------------------------------------------------------------
// Équipe de Pokémon.
//
// Le sujet modélise le joueur avec UN couple PV/PVMax : dans ce jeu, ce couple
// est toujours celui du Pokémon ACTIF, c'est-à-dire Equipe[0]. Changer de
// Pokémon = échanger un membre avec Equipe[0] et recharger PV/PVMax.
// ---------------------------------------------------------------------------

const equipeMax = 6

// Pokemon est un membre de l'équipe du joueur. Chacun a sa propre progression
// et son propre moveset (adapté à son type).
type Pokemon struct {
	Nom           string
	Type          string
	Niveau        int
	PVMax         int
	PV            int
	Experience    int
	ExperienceMax int
	Skills        []string
}

func copieSkills(s []string) []string {
	c := make([]string, len(s))
	copy(c, s)
	return c
}

func xpMaxDepart(p *Personnage) int {
	if p.ExperienceMax > 0 {
		return p.ExperienceMax
	}
	return 100
}

// initEquipe crée l'équipe de départ à partir du starter choisi.
func initEquipe(p *Personnage) {
	p.Equipe = []Pokemon{{
		Nom:           p.Starter,
		Type:          p.Type,
		Niveau:        p.Niveau,
		PVMax:         p.PVMax,
		PV:            p.PV,
		Experience:    p.Experience,
		ExperienceMax: xpMaxDepart(p),
		Skills:        copieSkills(p.Skills),
	}}
}

// syncPersoVersActif recopie l'état du joueur dans le Pokémon actif (Equipe[0]).
func syncPersoVersActif(p *Personnage) {
	if len(p.Equipe) == 0 {
		initEquipe(p)
		return
	}
	p.Equipe[0].PV = p.PV
	p.Equipe[0].PVMax = p.PVMax
	p.Equipe[0].Niveau = p.Niveau
	p.Equipe[0].Experience = p.Experience
	p.Equipe[0].ExperienceMax = p.ExperienceMax
	p.Equipe[0].Type = p.Type
	p.Equipe[0].Skills = copieSkills(p.Skills)
}

// chargerActifVersPerso recharge l'état du joueur depuis le Pokémon actif.
func chargerActifVersPerso(p *Personnage) {
	e := p.Equipe[0]
	p.Starter = e.Nom
	p.Type = e.Type
	p.PVMax = e.PVMax
	p.PV = e.PV
	p.Niveau = e.Niveau
	if e.ExperienceMax > 0 {
		p.Experience = e.Experience
		p.ExperienceMax = e.ExperienceMax
	}
	if len(e.Skills) > 0 {
		p.Skills = copieSkills(e.Skills)
	}
}

// distribuerXP : le Pokémon actif reçoit l'XP complète (Mission 4, inchangé) et
// les autres membres de l'équipe en reçoivent la moitié (« Multi-Exp »).
func distribuerXP(p *Personnage, xp int) {
	gainExperience(p, xp)
	syncPersoVersActif(p)
	if len(p.Equipe) < 2 {
		return
	}
	partage := xp / 2
	if partage < 1 {
		partage = 1
	}
	for i := 1; i < len(p.Equipe); i++ {
		gagnerXPMembre(&p.Equipe[i], partage)
	}
}

// gagnerXPMembre fait progresser un Pokémon non actif de l'équipe.
func gagnerXPMembre(pk *Pokemon, xp int) {
	if pk.ExperienceMax <= 0 {
		pk.ExperienceMax = 100
	}
	pk.Experience += xp
	fmt.Printf("%s gagne %d points d'expérience (Multi-Exp).\n", pk.Nom, xp)
	for pk.Experience >= pk.ExperienceMax {
		pk.Experience -= pk.ExperienceMax
		pk.Niveau++
		pk.ExperienceMax = int(float64(pk.ExperienceMax) * 1.5)
		pk.PVMax += 8
		pk.PV = pk.PVMax
		bruitage("level")
		fmt.Printf(">>> %s passe niveau %d !\n", pk.Nom, pk.Niveau)
		mettreAJourSortsPokemon(pk)
	}
}

// mettreAJourEquipe recale le moveset de tous les membres (rattrapage au chargement).
func mettreAJourEquipe(p *Personnage) {
	mettreAJourSorts(p)
	for i := range p.Equipe {
		mettreAJourSortsPokemon(&p.Equipe[i])
	}
}

// soignerEquipe remet tous les Pokémon de l'équipe à PV max.
func soignerEquipe(p *Personnage) {
	if len(p.Equipe) == 0 {
		initEquipe(p)
	}
	for i := range p.Equipe {
		p.Equipe[i].PV = p.Equipe[i].PVMax
	}
	p.Equipe[0].PV = p.PV // l'actif suit p.PV, remis à max par soigner()
}

// ajouterAEquipe intègre un Pokémon capturé. Renvoie false si l'équipe est pleine.
func ajouterAEquipe(p *Personnage, m *Monster, niveau int) bool {
	if len(p.Equipe) >= equipeMax {
		return false
	}
	p.Equipe = append(p.Equipe, Pokemon{
		Nom:           m.Nom,
		Type:          m.Type,
		Niveau:        niveau,
		PVMax:         m.PVMax,
		PV:            m.PVMax,
		ExperienceMax: 100,
		Skills:        []string{"Charge"},
	})
	mettreAJourSortsPokemon(&p.Equipe[len(p.Equipe)-1]) // moveset selon son niveau/type
	return true
}

// afficherEtatEquipe montre la santé de chaque Pokémon (menu, Centre, pause carte).
func afficherEtatEquipe(p *Personnage) {
	if len(p.Equipe) == 0 {
		initEquipe(p)
	}
	syncPersoVersActif(p)
	fmt.Println("Équipe :")
	for i, e := range p.Equipe {
		marque := "  "
		if i == 0 {
			marque = col("> ", cVert)
		}
		etat := ""
		if e.PV <= 0 {
			etat = col(" K.O.", cRouge)
		}
		fmt.Printf("%s%-12s Nv.%-2d %s %d/%d%s\n",
			marque, e.Nom, e.Niveau, barre(e.PV, e.PVMax, 14), e.PV, e.PVMax, etat)
	}
}

// hintMatchup : indication de type d'un Pokémon face à un adversaire.
// Rien si c'est neutre (ni efficace, ni peu efficace).
func hintMatchup(typePokemon, typeAdverse string) string {
	if typeAdverse == "" || typePokemon == "" {
		return ""
	}
	switch mult, _ := efficacite(typePokemon, typeAdverse); {
	case mult >= 2:
		return col(" ★ attaques efficaces", cVert)
	case mult <= 0.5:
		return col(" ✗ attaques peu efficaces", cRouge)
	}
	return ""
}

// changerPokemon : le joueur choisit un autre Pokémon actif. typeAdverse (peut
// être "") sert à afficher quels Pokémon sont efficaces contre l'adversaire.
// Renvoie true si le changement a eu lieu (le tour est alors consommé).
func changerPokemon(p *Personnage, typeAdverse string) bool {
	if len(p.Equipe) < 2 {
		fmt.Println("Tu n'as qu'un seul Pokémon.")
		return false
	}
	syncPersoVersActif(p)
	fmt.Println("=== Changer de Pokémon ===")
	for i, e := range p.Equipe {
		statut := ""
		if i == 0 {
			statut = " (actif)"
		} else if e.PV <= 0 {
			statut = col(" K.O.", cRouge)
		}
		fmt.Printf("%d. %-12s %s %d/%d%s%s\n", i+1, e.Nom, barre(e.PV, e.PVMax, 12), e.PV, e.PVMax, statut, hintMatchup(e.Type, typeAdverse))
	}
	fmt.Println("R. Retour")

	choix := lireLigne("> ")
	if choix == "R" || choix == "r" {
		return false
	}
	idx := indexFromChoice(choix, len(p.Equipe))
	if idx <= 0 {
		if idx == 0 {
			fmt.Println("Ce Pokémon est déjà au combat.")
		} else {
			fmt.Println("Choix invalide.")
		}
		return false
	}
	if p.Equipe[idx].PV <= 0 {
		fmt.Println(p.Equipe[idx].Nom + " est K.O. et ne peut pas se battre.")
		return false
	}
	ancien := p.Equipe[0].Nom
	p.Equipe[0], p.Equipe[idx] = p.Equipe[idx], p.Equipe[0]
	chargerActifVersPerso(p)
	TypeText(fmt.Sprintf("Reviens, %s ! En avant, %s !", ancien, p.Equipe[0].Nom))
	return true
}

// gererKO est appelé quand le Pokémon actif tombe à 0 PV. Si un autre Pokémon
// est valide, on force un changement et le combat continue (renvoie true).
// Sinon toute l'équipe est à terre (renvoie false).
func gererKO(p *Personnage, typeAdverse string) bool {
	syncPersoVersActif(p)
	p.Equipe[0].PV = 0

	vivants := []int{}
	for i := 1; i < len(p.Equipe); i++ {
		if p.Equipe[i].PV > 0 {
			vivants = append(vivants, i)
		}
	}
	if len(vivants) == 0 {
		return false
	}

	TypeText(fmt.Sprintf("%s est K.O. !", p.Equipe[0].Nom))
	fmt.Println("Choisis ton prochain Pokémon :")
	for _, i := range vivants {
		e := p.Equipe[i]
		fmt.Printf("%d. %-12s %s %d/%d%s\n", i+1, e.Nom, barre(e.PV, e.PVMax, 12), e.PV, e.PVMax, hintMatchup(e.Type, typeAdverse))
	}
	var idx int
	for {
		idx = indexFromChoice(lireLigne("> "), len(p.Equipe))
		if idx > 0 && p.Equipe[idx].PV > 0 {
			break
		}
		fmt.Println("Choix invalide.")
	}
	p.Equipe[0], p.Equipe[idx] = p.Equipe[idx], p.Equipe[0]
	chargerActifVersPerso(p)
	TypeText(fmt.Sprintf("En avant, %s !", p.Equipe[0].Nom))
	return true
}
