# Projet RED — Pokémon Noir 3

Mini-jeu de rôle en ligne de commande, écrit en Go, dans l'univers de *Pokémon
Version Noire*. Le cadre imposé par le sujet (personnage, inventaire, marchand,
forgeron, combat tour par tour) est respecté intégralement ; l'habillage
« Pokémon » et une trame en chapitres viennent par-dessus (Bonus 1).

## Lancer le jeu

```
go run .
```
ou
```
go build -o red . && ./red        # Windows : red.exe
```

Variables d'environnement optionnelles :
- `NO_COLOR=1` : désactive les couleurs ANSI.
- `RED_FAST=1` : désactive l'animation « machine à écrire » (utile pour tester).
- `RED_SOUND=1` / `RED_SOUND=0` : force le son activé / désactivé (sinon le jeu
  le demande au lancement si le dossier `audio/` contient des fichiers).
- `RED_VOL=0..100` : volume du son (35 par défaut).

Le son (musiques `.mp3`/`.wav` dans `audio/`) est **optionnel et Windows
uniquement** ; sans fichiers, sans Windows ou en répondant « n », le jeu tourne
en silence, sans aucun impact. Voir `audio/README.txt`.

Le fichier `save.json` est créé à côté de l'exécutable quand on sauvegarde au
Centre Pokémon. Le supprimer revient à repartir de zéro.

## Commandes

- **Menu principal** : on tape le numéro du choix puis Entrée.
- **Carte (choix 6)** : flèches directionnelles **ou** `Z` `Q` `S` `D` pour se
  déplacer, `R` pour revenir au menu. On entre dans un lieu ou on parle à un
  personnage en se dirigeant dessus. Si l'entrée standard n'est pas un vrai
  terminal, la carte bascule en lecture ligne par ligne (une lettre + Entrée).

## Déroulé

Création du personnage (nom, genre, starter) → prologue → exploration de la
région. La trame se joue via le choix **6. Explorer la région** :

| Ch. | Objectif | Lieu débloqué |
|----|----------|---------------|
| 1 | Rejoindre la Route 1 puis Arabelle | Route 1, Arabelle |
| 2 | S'équiper, se soigner, battre le Maître d'Arène | — |
| 3 | Traverser la Forêt d'Ombreflore | Forêt |
| 4 | Infiltrer le QG de la Neo-Plasma, battre Nikolai | QG |
| 5 | Affronter N et son Zekrom | — |
| 6 | Épilogue (Ghetsis & Kyurem), exploration libre | — |

Le menu classique (infos, inventaire, marchand, forgeron, entraînement) reste
accessible à tout moment.

## Organisation des fichiers

```
RED-PROJECT/
├── main.go                point d'entrée (appelle game.Run())
├── go.mod / go.sum
├── README.md
├── internal/
│   └── game/              TOUT le jeu (paquet "game")
│       ├── *.go           un domaine par fichier (tableau ci-dessous)
│       └── ascii/*.txt    dessins ASCII, embarqués via //go:embed
├── audio/                 musiques/bruitages optionnels + README.txt
└── plan/                  énoncé (PDF) et fiche de route
```

Le jeu tient dans **un seul paquet `game`** (les fonctions se référencent trop
mutuellement pour être éclatées en sous-paquets sans cycles d'import). La racine
ne contient qu'un `main.go` minimal, ce qui garde le projet lisible.
Lancer / compiler se fait toujours depuis la racine (`go run .`,
`go build -o red .`) : `audio/` et `save.json` y sont cherchés à l'exécution.

**Cœur du sujet** (tout dans `internal/game/`)

| Fichier | Rôle |
|---------|------|
| `run.go` | `Run()` : menu principal (switch/case), lancement / reprise de partie |
| `character.go` | struct `Personnage`, `Init`, `charCreation`, `displayInfo`, `dead`, `gainExperience` |
| `inventory.go` | struct `Objet`, `accessInventory`, `takePot`/`poisonPot`/`takeManaPot`, `addInventory`/`removeInventory`, limite 10, `upgradeInventorySlot` |
| `merchant.go` | struct `ArticleMarchand`, `merchantMenu`, `acheter` (avec quantité) |
| `forge.go` | struct `Equipment`, struct `Recette`, `forgeronMenu`, `fabriquer`, `equip` |
| `spells.go` | struct `Sort`, `spellBook`, `choisirAttaque` (sorts en combat, coût de mana) |
| `monster.go` | struct `Monster`, bestiaire `especes`, `InitRatentif` (= InitGoblin), `InitZorua`, `InitGolette` |
| `combat.go` | `deroulerCombat`, `trainingFight`, `charTurn`, `monsterPattern` (= goblinPattern), capture, fuite |
| `utils.go` | `ClearScreen`, `TypeText`, `ReadKey` |

**Extensions (Bonus 1)**

| Fichier | Rôle |
|---------|------|
| `party.go` | équipe de Pokémon (capture, changement, K.O., état de santé) |
| `world.go` / `worldzones.go` | moteur de carte multi-zones + définition des 5 zones |
| `story.go` | trame « Pokémon Noir 3 » en 6 chapitres, scènes, boss N & Ghetsis |
| `npc.go` | dialogues (`jouerDialogue`), combats de dresseurs |
| `save.go` | sauvegarde / chargement JSON, écran « continuer / recommencer » |
| `ui.go` | couleurs ANSI, jauges de PV, cadre « dessin à venir » |
| `sound.go` | lecteur audio optionnel (MCI, Windows), désactivé par défaut |
| `ascii.go` | chargement des dessins `ascii/*.txt` (`//go:embed`) |

## Correspondance avec l'énoncé

Tout ce que le sujet impose est présent. Certains **noms** ont été adaptés à
l'univers Pokémon ; les **valeurs chiffrées** sont celles de l'énoncé.

| Sujet | Dans le jeu |
|-------|-------------|
| T1 Menu `switch/case` + saisie | `run.go` (`Run()`) ; `main.go` (racine) appelle `game.Run()` |
| T2 Struct `Personnage` (nom, classe, niveau, PV max/actuels, inventaire) | `character.go` — la « classe » est le **type du starter** (Plante 120 / Feu 100 / Eau 110 PV) |
| T3 `Init` | `Init` dans `character.go` (remplacé par `charCreation`, cf. Mission 1) |
| T4 `displayInfo` | `character.go` |
| T5 `accessInventory` (+ retour) | `inventory.go` |
| T6 `takePot` (+50 PV, cap PV max) | `inventory.go` |
| T7 Marchand + 1 potion gratuite + `addInventory`/`removeInventory` | `merchant.go` (choix `0` = échantillon gratuit, 1 seule fois) |
| T8 `dead` (résurrection à 50 % PV max) | `character.go` |
| T9 `poisonPot` (−10 PV/s pendant 3 s, `time.Sleep`) | `inventory.go` — « Baie Pecha » chez le marchand |
| T10 `skill`, sort de base, `spellBook`, « déjà appris » | `spells.go` — sort de base **« Charge »** (= Coup de poing, 8 dégâts), livre **« CT35 - Lance-Flammes »** (= Boule de feu, 18 dégâts) |
| Mission 1 `charCreation` (nom lettres uniquement + formaté, classe → PV, PV départ = 50 %, niveau 1) | `character.go` (`estAlphabetique`, `formaterNom`) |
| Mission 2 Limite d'inventaire à 10 (check à l'ajout) | `inventory.go` (`checkInventoryFull` dans `addInventory`) |
| T11 Argent (100 or au départ) | `character.go` |
| T12 Prix (3 / 6 / 25) + 4 ressources (4 / 7 / 3 / 1) | `merchant.go` — ressources renommées (Plume de Poichigeon = Plume de Corbeau, etc.) |
| T13 Forgeron : 3 équipements, −5 or, coût ressources, messages d'erreur | `forge.go` — Restes = Chapeau, Veste de Combat = Tunique, Grelot Coque = Bottes |
| T14 Struct `Equipment` (tête/torse/pieds) sur `Personnage` | `forge.go` / `character.go` |
| T15 Équiper (retrait inventaire, +10 / +25 / +15 PV max, remplacement) | `forge.go` (`equip`) |
| Mission 2.2 `upgradeInventorySlot` (+10, ×3 max, 30 or) | `inventory.go` + `merchant.go` |
| T16 Struct `Monster` + `InitGoblin` (PV 40, ATQ 5) | `monster.go` — le « Gobelin d'entraînement » est le **Ratentif** |
| T17 `trainingFight` + variable de tour | `combat.go` |
| T18 `goblinPattern` (100 %, 200 % tous les 3 tours, affichages) | `combat.go` (`monsterPattern`) |
| T19 `charTurn` (Attaquer 5 dégâts / Inventaire) | `combat.go` (+ Sorts, + Capturer) |
| T20 Boucle de combat, tour affiché, option menu, fin à 0 PV → menu | `combat.go` + `run.go` |
| Mission 3 Initiative (le plus rapide commence) | `combat.go` (`deroulerCombat`) |
| Mission 4 XP + niveau + report du surplus + XP max croissant | `character.go` (`gainExperience`) |
| Mission 5.1 Sorts offensifs en combat (8 / 18 dégâts) | `spells.go` (`choisirAttaque`) |
| Mission 5.2 Mana + coût + blocage si insuffisant + potion de mana | `spells.go` + `inventory.go` — potion = « Huile » |
| Bonus 1 Enrichir le jeu | trame Pokémon Noir 3, carte multi-zones, capture, dresseurs, boss, sauvegarde |
| Bonus 2 « Qui sont-ils ? » | `whoAreThey` : **ABBA** (partie 2) et **Steven Spielberg** (partie 3) |

### Note sur l'attaque de base

Le sujet fixe l'attaque basique à **5 dégâts** (T19). Elle vaut exactement 5 au
niveau 1 ; elle gagne `+2` par niveau au-delà, ce qui relève du « gain de
statistiques compté en bonus » explicitement autorisé par la Mission 4.

## Choix de conception

- **Un seul fichier par domaine**, noms de fonctions du sujet conservés
  (`takePot`, `spellBook`, `charTurn`, `monsterPattern`…).
- **Dessins ASCII** chargés via `//go:embed` : le binaire est autonome. Un
  fichier `ascii/<nom>.txt` absent ou vide affiche un cadre « dessin à venir »
  au lieu de planter.
- **Sauvegarde** JSON simple (hors sujet, confort pour une partie en chapitres).
- Messages d'erreur personnalisés partout (or insuffisant, inventaire plein,
  ressources manquantes, mana insuffisant, sort déjà appris).
