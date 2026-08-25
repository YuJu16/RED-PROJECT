package main

// Monster représente un adversaire rencontré en combat.
type Monster struct {
	Nom        string
	PVMax      int
	PV         int
	Attaque    int
	Initiative int
	Experience int
	Ascii      string
}

// InitRatentif initialise un Ratentif sauvage (équivalent du Gobelin d'entrainement du sujet).
func InitRatentif() *Monster {
	return &Monster{
		Nom:        "Ratentif sauvage",
		PVMax:      40,
		PV:         40,
		Attaque:    5,
		Initiative: 4,
		Experience: 25,
		Ascii:      "ratentif",
	}
}

// InitZorua initialise un Zorua sauvage, adversaire un peu plus vif et fragile.
func InitZorua() *Monster {
	return &Monster{
		Nom:        "Zorua sauvage",
		PVMax:      55,
		PV:         55,
		Attaque:    7,
		Initiative: 8,
		Experience: 40,
		Ascii:      "zorua",
	}
}

// InitGolette initialise un Golette sauvage, adversaire résistant mais lent.
func InitGolette() *Monster {
	return &Monster{
		Nom:        "Golette sauvage",
		PVMax:      70,
		PV:         70,
		Attaque:    6,
		Initiative: 2,
		Experience: 55,
		Ascii:      "golette",
	}
}
