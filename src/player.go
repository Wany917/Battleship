package main

import (
	"fmt"
	"strings"
)

type player struct {
	playerIp     string
	playerName   string
	playerHealth int
	playermap    Playermap_s
	munition     int
}

type Target struct {
	targetName string
	isTrue     bool
}

var playerList []player

func VerifAttack(attack string) bool {
	isAttack := false
	verif := [100]string{
		"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9",
		"B0", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9",
		"C0", "C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9",
		"D0", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9",
		"E0", "E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9",
		"F0", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9",
		"G0", "G1", "G2", "G3", "G4", "G5", "G6", "G7", "G8", "G9",
		"H0", "H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9",
		"I0", "I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8", "I9",
		"J0", "J1", "J2", "J3", "J4", "J5", "J6", "J7", "J8", "J9"}
	for i := 0; i != len(verif); i++ {
		if attack == verif[i] {
			isAttack = true
			break
		}
	}
	return (isAttack)
}

func IsGameOver(player player) bool {

	GameOver := true
	for i := 0; i != len(player.playermap.Ships); i++ {
		if player.playermap.Ships[i].isSunken == false {
			GameOver = false
		}
	}
	return GameOver
}

func applyAttack(playermap *Playermap_s, attack Coordinate) {
	for i, ship := range playermap.Ships {
		if attack.row >= ship.start.row && attack.row <= ship.end.row && attack.col >= ship.start.col && attack.col <= ship.end.col {
			// Attack was successful, mark the attack on the ship's boat slice
			rowIndex := attack.row - ship.start.row
			if ship.boat[rowIndex] == '0' {
				ship.boat[rowIndex] = '*'
				playermap.Ships[i] = ship
			}
			break
		}
	}
}

func attackPlayer(attack string, attackedPlayer player) (player, bool) {
	var hit bool
	if VerifAttack(attack) == false { // error menagement if user attack a wrong position
		fmt.Printf("position incorrecte\n")
		return attackedPlayer, false
	}
	tab := make(map[rune]int)
	tab['A'] = 0
	tab['B'] = 1
	tab['C'] = 2
	tab['D'] = 3
	tab['E'] = 4
	tab['F'] = 5
	tab['G'] = 6
	tab['H'] = 7
	tab['I'] = 8
	tab['J'] = 9

	j := tab[rune(attack[0])]
	i := int(attack[1] - '0')
	updateMap(&attackedPlayer.playermap)
	if attackedPlayer.playermap.playermap[i][j] == '0' {
		//attackedPlayer.playermap.playermap[i][j] = '*' // player hit
		fmt.Printf("\n%s%s: Touché%s\n", GREEN, attack, RESET)
		applyAttack(&attackedPlayer.playermap, Coordinate{i, j})
		moveShips(&attackedPlayer.playermap)
		hit = true
	} else {
		attackedPlayer.playermap.playermap[i][j] = ' ' // player missed
		fmt.Printf("\n%s%s: Raté%s\n", RED, attack, RESET)
		hit = false
	}
	return attackedPlayer, hit
}

func createPlayer() player { //crée un joueur
	playerIp, err := LocalIP() //récupère l'ip du joueur

	if err != nil {
		fmt.Println(RED, "Adresse IP Invalide", RESET) //Si un problème au niveau de l'ip Revoie l'utilisateur avec comme nom ERROR pour couper la suite du programme
		user := player{playerIp: "ERROR", playerName: "ERROR"}
		return user
	} else { //sinon lui demande son nom et renvoie l'utilisateur
		user := player{playerIp: playerIp.String() + port, playerName: "", playerHealth: 10, munition: 5}
		fmt.Printf("Entrez votre nom capitaine : ")
		fmt.Scanf("%s", &user.playerName)
		user.playermap = generateMap() // generate a random map
		return user
	}
}

func isAlive(user player) bool {
	if user.playerHealth > 0 {
		return true
	} else {
		return false
	}
}

// Tue le joueur en local (l'enleve de sa liste de joueur)
func killPlayer(y int) {
	i := y
	x := player{playerIp: "", playerName: "", playerHealth: 0}
	playerList[i] = playerList[len(playerList)-1]
	playerList[len(playerList)-1] = x
	playerList = playerList[:len(playerList)-1]
}

func getPosition() int {
	x := len(playerList)
	return x + 1
}

func ranking(x int) {
	switch x {
	case 1:
		fmt.Println(GREEN, "Bravo ! Tu est le grand vainqueur de cette bataille !", RESET)
	case 2:
		fmt.Println(GREEN, "C'est fini, tout nos navires sont tombés.", RESET)
		fmt.Println(GREEN, "Bien joué ! Tu as fini 2eme", RESET)
	case 3:
		fmt.Println(GREEN, "C'est fini, tout nos navires sont tombés.", RESET)
		fmt.Println(GREEN, "Bien joué ! Tu as fini 3eme", RESET)
	default:
		fmt.Println(RED, "C'est fini, tout nos navires sont tombés.", RESET)
		fmt.Println(RED, "Dommage, tu as fini %deme\n", x, RESET)
	}
}

func printPlayerList() {
	for i := 0; i < len(playerList); i++ {
		fmt.Println(playerList[i])
	}
}

func checkPlayersList(playerIp string) bool { //Renvoie si le joueur a déjà été ajouté à la liste des joueurs
	for i := 0; i < len(playerList); i++ {
		if strings.Compare(playerList[i].playerIp, playerIp) == 0 {
			return true
		}
	}

	return false
}

// vérifie les nom/IP entré par l'utilisateur et renvoie un tableau de structure TargetList[{nom, true}{nom,false}] pour savoir si les cibles sont bonnes
func checkUserNames(nameArray []string) []Target {
	tmpTargetList := [6]Target{}
	var size = 0
	for i := 0; i < len(nameArray)-1; i++ { // -1 car le dernier élément du tableau donner est l'action donc sert a rien

		tmpTargetList[i].targetName = nameArray[i] //copie les nom dans la liste des cibles
		if strings.Compare(tmpTargetList[i].targetName, "") != 0 {
			size++
		}
		if isPlayer(nameArray[i]) { //vérifie chaque nom, si il est bon
			tmpTargetList[i].isTrue = true //le isTrue passe a vraie
		} else {
			tmpTargetList[i].isTrue = false //sinon il passe a faux
		}
	}

	targetList := make([]Target, size)

	for j := 0; j < len(nameArray); j++ {
		if strings.Compare(tmpTargetList[j].targetName, "") == 0 {
			continue
		} else {
			targetList[j].targetName = tmpTargetList[j].targetName
			targetList[j].isTrue = tmpTargetList[j].isTrue
		}
	}

	return targetList
}

// vérifie que le joueur entré par l'utilisateur éxiste bien que ce soit par nom ou IP
func isPlayer(name string) bool {
	ip := name + port
	for i := 1; i < len(playerList); i++ {
		if strings.Compare(playerList[i].playerName, name) == 0 || strings.Compare(playerList[i].playerIp, ip) == 0 { //Si nomDansListe = Entré de l'user ou IP = Entré de l'user alors TRUE
			return true
		}
	}
	return false //sinon FALSE
}

func getTargetFromName(name string) player {
	ip := name + port
	var target int
	for i := 0; i < len(playerList); i++ {
		if strings.Compare(name, playerList[i].playerName) == 0 || strings.Compare(ip, playerList[i].playerIp) == 0 {
			target = i
		}
	}
	return playerList[target]
}
