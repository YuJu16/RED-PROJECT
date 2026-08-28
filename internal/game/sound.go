package game

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf16"
)

// ---------------------------------------------------------------------------
// Son (optionnel, Windows).
//
// UN SEUL processus PowerShell caché ("démon audio") est lancé au démarrage. On
// lui envoie des commandes texte sur son entrée standard :
//   music <chemin>   -> musique de fond en boucle (remplace la précédente)
//   sfx <chemin>     -> bruitage court (remplace le précédent)
//   stop             -> coupe tout
//   vol <0..1000>    -> volume
//   quit             -> fin
// Garantie : au plus 1 musique + 1 bruitage simultanés, jamais de superposition.
//
// RED_SOUND=1/0 force activé/désactivé ; sinon on demande au lancement.
// RED_VOL=0..100 règle le volume (35 % par défaut). Aucun effet hors Windows.
// ---------------------------------------------------------------------------

var (
	soundActif      bool
	volumeMCI       = 350
	musiqueCourante string // nom logique du morceau de fond en cours
	musiquePath     string // chemin réel en cours (évite de relancer le même fichier)
	audioCmd        *exec.Cmd
	audioIn         io.WriteCloser
)

// soundInit décide si le son est actif puis démarre le démon audio.
func soundInit() {
	if v, err := strconv.Atoi(os.Getenv("RED_VOL")); err == nil {
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		volumeMCI = v * 10
	}

	if runtime.GOOS != "windows" {
		return
	}
	switch os.Getenv("RED_SOUND") {
	case "1", "on", "true":
		soundActif = true
	case "0", "off", "false":
		return
	default:
		if !dossierAudioRempli() {
			return
		}
		rep := strings.ToLower(lireLigne("Activer la musique et les sons ? [O/n] "))
		soundActif = rep == "" || rep == "o" || rep == "oui" || rep == "y"
	}
	if soundActif {
		demarrerDemonAudio()
	}
}

func dossierAudioRempli() bool {
	for _, ext := range []string{"*.mp3", "*.wav"} {
		if m, _ := filepath.Glob(filepath.Join("audio", ext)); len(m) > 0 {
			return true
		}
	}
	return false
}

// aliasAudio : repli sur un autre nom si le fichier demandé est absent.
var aliasAudio = map[string]string{
	"menu":          "title",
	"plasma_battle": "battle_trainer",
	"N_battle":      "battle_boss", // repli tant que l'OST dédiée de N n'est pas là
}

// cheminAudio renvoie le chemin absolu du .mp3/.wav, "" si introuvable.
func cheminAudio(nom string) string {
	for _, ext := range []string{".mp3", ".wav"} {
		p := filepath.Join("audio", nom+ext)
		if _, err := os.Stat(p); err == nil {
			if abs, err := filepath.Abs(p); err == nil {
				return abs
			}
		}
	}
	if a, ok := aliasAudio[nom]; ok {
		return cheminAudio(a)
	}
	return ""
}

const pidAudio = "audio/.redaudio.pid"

// demonScript : format = fmt.Sprintf(demonScript, volumeMCI, pidDuJeu).
// Le démon lit ses commandes sur stdin ET surveille le PID du jeu : dès que
// le jeu disparaît (fermeture propre, crash, terminal fermé), il coupe le
// son et se termine — jamais d'orphelin qui continue de jouer.
const demonScript = `$ErrorActionPreference='SilentlyContinue'
try { $PID | Out-File -Encoding ascii -Force 'audio\.redaudio.pid' } catch {}
Add-Type -MemberDefinition '[DllImport("winmm.dll",CharSet=CharSet.Auto)] public static extern int mciSendString(string c,System.Text.StringBuilder r,int l,System.IntPtr h);' -Name M -Namespace R
function mci($s){ [R.M]::mciSendString($s,$null,0,0) | Out-Null }
function fin(){ mci('close redm'); mci('close reds'); Remove-Item 'audio\.redaudio.pid' -ErrorAction SilentlyContinue }
$vol = %d
$jeu = %d
# StreamReader sur le flux brut : ReadLineAsync y est reellement asynchrone
# (celui de [Console]::In est synchrone et bloque, ce qui neutralise le watchdog).
$sr = New-Object System.IO.StreamReader([Console]::OpenStandardInput())
while ($true) {
  $t = $sr.ReadLineAsync()
  while (-not $t.Wait(1200)) {
    if (-not (Get-Process -Id $jeu -ErrorAction SilentlyContinue)) { fin; [Environment]::Exit(0) }
  }
  $line = $t.Result
  if ($null -eq $line) { break }
  $i = $line.IndexOf(' ')
  if ($i -lt 0) { $v = $line; $a = '' } else { $v = $line.Substring(0,$i); $a = $line.Substring($i+1) }
  switch ($v) {
    'music'  { mci('close redm'); mci('open "' + $a + '" alias redm'); mci('setaudio redm volume to ' + $vol); mci('play redm repeat') }
    'jingle' { mci('close reds'); mci('close redm'); mci('open "' + $a + '" alias redm'); mci('setaudio redm volume to ' + $vol); mci('play redm') }
    'sfx'    { mci('close reds'); mci('open "' + $a + '" alias reds'); mci('setaudio reds volume to ' + $vol); mci('play reds from 0') }
    'stop'   { mci('close redm'); mci('close reds') }
    'vol'    { $vol = [int]$a }
    'quit'   { break }
  }
}
fin`

// tuerDemonOrphelin tue un éventuel lecteur audio laissé par une session
// précédente (fenêtre du terminal fermée sans quitter proprement).
func tuerDemonOrphelin() {
	data, err := os.ReadFile(pidAudio)
	if err != nil {
		return
	}
	pid := strings.TrimSpace(string(data))
	if pid != "" {
		c := exec.Command("taskkill", "/F", "/PID", pid)
		processusCache(c)
		_ = c.Run()
	}
	_ = os.Remove(pidAudio)
}

// demarrerDemonAudio lance l'unique processus PowerShell caché.
func demarrerDemonAudio() {
	tuerDemonOrphelin()
	script := fmt.Sprintf(demonScript, volumeMCI, os.Getpid())
	u16 := utf16.Encode([]rune(script))
	b := make([]byte, len(u16)*2)
	for i, r := range u16 {
		b[i*2] = byte(r)
		b[i*2+1] = byte(r >> 8)
	}
	enc := base64.StdEncoding.EncodeToString(b)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-EncodedCommand", enc)
	processusCache(cmd) // CREATE_NO_WINDOW -> aucune fenêtre, aucun vol de focus
	in, err := cmd.StdinPipe()
	if err != nil {
		soundActif = false
		return
	}
	if err := cmd.Start(); err != nil {
		soundActif = false
		return
	}
	audioCmd, audioIn = cmd, in
}

func envoyerAudio(ligne string) {
	if audioIn == nil {
		return
	}
	_, _ = io.WriteString(audioIn, ligne+"\n")
}

// musique lance une musique de fond en boucle. Sans effet si déjà en cours.
// Si le fichier n'existe pas, on ne change rien (musiqueCourante inchangée) :
// la musique en place continue et la détection pourra réessayer plus tard.
func musique(nom string) {
	if !soundActif || nom == musiqueCourante {
		return
	}
	chemin := cheminAudio(nom)
	if chemin == "" {
		return
	}
	musiqueCourante = nom
	if chemin == musiquePath {
		return // même fichier déjà en cours (ex. title -> menu)
	}
	musiquePath = chemin
	envoyerAudio("music " + chemin)
}

// bruitage joue un son COURT en surimpression (level up, K.O.). À réserver aux
// sons brefs : un fichier long resterait par-dessus la musique.
func bruitage(nom string) {
	if !soundActif {
		return
	}
	if chemin := cheminAudio(nom); chemin != "" {
		envoyerAudio("sfx " + chemin)
	}
}

// jingle joue un morceau UNE fois à la place de la musique de fond (fanfare de
// victoire...). La musique de zone reprend dès l'appel suivant à musique().
func jingle(nom string) {
	if !soundActif {
		return
	}
	if chemin := cheminAudio(nom); chemin != "" {
		musiqueCourante = "" // force le prochain musique() à relancer la zone
		envoyerAudio("jingle " + chemin)
	}
}

// soundStop coupe la musique et les bruitages en cours.
func soundStop() {
	musiqueCourante = ""
	musiquePath = ""
	envoyerAudio("stop")
}

// soundQuit ferme proprement le démon audio (à la sortie du jeu).
func soundQuit() {
	envoyerAudio("quit")
	if audioIn != nil {
		_ = audioIn.Close()
	}
	if audioCmd != nil && audioCmd.Process != nil {
		pid := audioCmd.Process.Pid
		_ = audioCmd.Process.Kill()
		// coup de grâce : tuer l'arbre du process au cas où
		c := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		processusCache(c)
		_ = c.Run()
	}
	_ = os.Remove(pidAudio)
	audioIn, audioCmd = nil, nil
}
