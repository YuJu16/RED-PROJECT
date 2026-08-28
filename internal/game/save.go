package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Fichiers de sauvegarde, à côté de l'exécutable.
const (
	cheminSauvegarde = "save.json"
	cheminSecours    = "save.json.bak" // copie de la sauvegarde précédente
	cheminAncienne   = "save.json.old" // sauvegarde mise de côté par "Recommencer"
)

// copierFichier copie src vers dst (best effort, erreurs ignorées).
func copierFichier(src, dst string) {
	if data, err := os.ReadFile(src); err == nil {
		_ = os.WriteFile(dst, data, 0644)
	}
}

// sauvegarder écrit l'état du personnage en JSON, après avoir mis l'ancienne
// sauvegarde de côté dans save.json.bak (filet de sécurité).
func sauvegarder(p *Personnage) error {
	if sauvegardeExiste() {
		copierFichier(cheminSauvegarde, cheminSecours)
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cheminSauvegarde, data, 0644)
}

// chargerDepuis relit un fichier de sauvegarde et remplace le contenu de *p.
func chargerDepuis(chemin string, p *Personnage) error {
	data, err := os.ReadFile(chemin)
	if err != nil {
		return err
	}
	var charge Personnage
	if err := json.Unmarshal(data, &charge); err != nil {
		return err
	}
	if charge.Nom == "" {
		return errors.New("sauvegarde corrompue")
	}
	*p = charge
	syncPersoVersActif(p)
	return nil
}

// charger relit la sauvegarde principale.
func charger(p *Personnage) error { return chargerDepuis(cheminSauvegarde, p) }

func fichierExiste(chemin string) bool {
	_, err := os.Stat(chemin)
	return err == nil
}

// sauvegardeExiste indique si la sauvegarde principale est présente.
func sauvegardeExiste() bool { return fichierExiste(cheminSauvegarde) }

// demarrerPartie : au lancement, propose de continuer, de charger le secours,
// ou de recommencer (l'ancienne partie est mise de côté, pas supprimée).
func demarrerPartie() *Personnage {
	p := &Personnage{}

	if sauvegardeExiste() || fichierExiste(cheminSecours) {
		ClearScreen()
		fmt.Println(titre("Sauvegarde trouvée"))
		if sauvegardeExiste() {
			fmt.Println("1. Continuer la partie")
		}
		if fichierExiste(cheminSecours) {
			fmt.Println("2. Charger la sauvegarde de secours (avant-dernière)")
		}
		fmt.Println("3. Nouvelle partie (l'ancienne est renommée save.json.old)")
		switch lireLigne("> ") {
		case "1":
			if sauvegardeExiste() && charger(p) == nil {
				configurerRival(p)
				mettreAJourEquipe(p)
				fmt.Println("Bon retour, " + p.Nom + " !")
				lireLigne("\nAppuyez sur Entrée...")
				return p
			}
			fmt.Println("Impossible de lire la sauvegarde principale.")
			lireLigne("\nAppuyez sur Entrée...")
		case "2":
			if chargerDepuis(cheminSecours, p) == nil {
				configurerRival(p)
				mettreAJourEquipe(p)
				fmt.Println("Sauvegarde de secours chargée. Bon retour, " + p.Nom + " !")
				lireLigne("\nAppuyez sur Entrée...")
				return p
			}
			fmt.Println("Secours illisible.")
			lireLigne("\nAppuyez sur Entrée...")
		}
		// Nouvelle partie : on met l'ancienne de côté sans la détruire.
		if sauvegardeExiste() {
			_ = os.Rename(cheminSauvegarde, cheminAncienne)
		}
	}

	p = Init()
	prologue(p)
	return p
}
