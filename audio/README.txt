DOSSIER AUDIO (optionnel)
=========================

Le son est DESACTIVE par defaut. Pour l'activer (PowerShell) :

    $env:RED_SOUND = 1
    $env:RED_VOL   = 20        # volume 0..100 (35 par defaut), propre au jeu
    .\red.exe

Sans RED_SOUND, hors Windows, ou si un fichier manque : aucun son, aucun impact
sur le jeu ni sur `go build`.

FORMAT : .mp3 OU .wav (le .mp3 est essaye en premier). AUCUNE conversion
necessaire, la lecture passe par MCI (winmm.dll). Place les fichiers ICI avec
EXACTEMENT ces noms (sans accent, sans espace) :

  MUSIQUES DE FOND (jouees en boucle)
  ----------------------------------
  title            ecran-titre                    [PRESENT]
  menu             menu principal                            (repli -> title)
  renouet          village de Renouet             [PRESENT]
  route1           Route 1                         [PRESENT]
  arabelle         ville d'Arabelle               [PRESENT]
  foret            Foret d'Ombreflore             [PRESENT]
  qg               QG de la Neo-Plasma            [PRESENT]
  battle_wild      combat contre un Pokemon sauvage [PRESENT]
  battle_trainer   combat contre un dresseur       [PRESENT]
  battle_boss      combat de boss (Ghetsis / Kyurem) [PRESENT]
  N_battle         combat contre N (Zekrom)                  (repli: battle_boss)
  low_hp           PV du Pokemon actif sous 25%    [PRESENT]
  n                scenes dialoguees avec N        [PRESENT]
  plasma           scenes Neo-Plasma (ambiance)    [PRESENT]
  plasma_battle    combats contre la Neo-Plasma    [PRESENT] (repli: battle_trainer)
  ghetsis          confrontation finale Ghetsis    [PRESENT]
  arene            Maitre d'Arene Cedric           [PRESENT]

  Repli auto : "menu" utilise "title" si menu.mp3 est absent.

  BRUITAGES (joues une fois)
  --------------------------
  victory          victoire en combat              [PRESENT]
  level            passage de niveau               [PRESENT]
  ko               equipe mise a terre                       (optionnel, court)

Tout fichier absent est simplement ignore.

NOTE : MCI (le lecteur de Windows utilise ici) refuse certains MP3 avec un
gros tag ID3v2.4 / pochette embarquee (code d'erreur 277 -> plus de son).
Si un morceau ne se lance pas : re-exporte-le sans metadonnees / sans image
(ou en .wav), ou retire le tag ID3.
