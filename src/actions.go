package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// liste des actions possible (a update pour chaque action que l'utilisateur peut faire si on en rajoute)
var actionList = [5]string{"/hit", "/board", "/boats", "/list", "/add"}

// fonction qui vas récupérer l'action souhaiter par l'utilisateur et la traiter
func userCommands() {
	var userCommand string
	fmt.Println("\nQue voulez-vous faire mon capitaine ? Voici la liste des actions possibles :")
	printActions() //Affiche à l'utilisateur les actions possibles
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userCommand = scanner.Text()
	userCommandArray := strings.Split(userCommand, " ") //Split ici pour récupérer chaque entrée séparément auu cas ou il entrent plusieurs adversaires lors du /hit
	if len(userCommandArray) > 6 || len(userCommandArray) < 2 && userCommand != "/list" {
		fmt.Println(RED, "Erreur, 5 cible maximum et 1 minimum !", RESET)
		return
	}
	action := userCommandArray[len(userCommandArray)-1] //stocke l'action ici
	if !isValidAction(action) {                         // si l'action est invalid (l'action est toujours le dernier truc entrer par l'user dans la commande)
		fmt.Println(RED, "ERREUR Action Invalide", RESET)
		return //affiche msg d'erreur et relance la fonction
	} //sinon

	target := checkUserNames(userCommandArray) //stocke la ou les cibles dans un tableaux de structure avec [{nom, valide ou non}] exemple : [{gil,true}{marceau,true}{mbappe,false}]

	executeUserAction(action, target) //execute l'action de l'utilisateur en fonction de l'action, envoie le tableau des cible pour savoir sur qui l'execute
}

func printActions() {
	fmt.Println("Pour effectuer une action veuillez suivre la structure suivante : xxx /kill ")
	fmt.Println("xxx = nom de l'adversaire ou IP , /xxx = action voulu")
	fmt.Println("xxx /board   : Affiche le plateau de l'adversaire ([0] = case intouché, [ ] = coup dans l'eau , [*] = bateau touché)")
	fmt.Println("xxx /boats   : Renvoie le nombre de bateau restant pour l'adversaire xxx")
	fmt.Println("xxx(5) /hit  : Tir sur une case choisi de l'adversaire    (5) = jusqu'à 5 utilisateur peuvent être entré en même temps  Cooldown de 6sc tous les 5 tirs")
	fmt.Println("xxx /list    : Affiche la liste des adversaires")
	fmt.Println("xxx /add     : Ajoute un ou plusieurs adversaire au radar, xxx = IP")
}

//Vérifie si la commande entrée est valide : dernier élément est bien une action valide

// vérifie si c'est bien une action
func isValidAction(userCommand string) bool {
	for i := 0; i < len(actionList); i++ {
		if strings.Compare(userCommand, actionList[i]) == 0 {
			return true
		}
	}
	return false
}

func executeUserAction(action string, targetList []Target) {
	switch action {
	case "/hit":
		for i := 0; i < len(targetList); i++ {
			if targetList[i].isTrue {
				target := getTargetFromName(targetList[i].targetName)
				if playerList[0].munition != 0 {
					launchAttack(target, playerList[0])
				} else {
					println("Chargement des canons en cour ! Veuillez patienter...")
					time.Sleep(1 * time.Second)
					i = i - 1
					continue
				}

			} else {
				println(targetList[i].targetName + " Cible invalide")
			}
		}
	case "/board":

		for i := 0; i != len(targetList); i++ {
			if targetList[i].targetName == playerList[0].playerName {
				addLog(playerList[0].playerName + " /board " + playerList[i].playerName)
				DisplayPlayerMap(playerList[0], true)
				continue
			}
			addLog(playerList[0].playerName + " /board " + playerList[i].playerName)
			index := getPlayerMap(targetList[i].targetName)
			DisplayPlayerMap(playerList[index], false)
		}
	case "/boats":
		getBoatsNumber(targetList[0].targetName)
	case "/list":
		fmt.Println()
		for i := 0; i != len(playerList); i++ {
			addLog(playerList[0].playerName + " /list " + playerList[i].playerName)
			ip := playerList[i].playerIp[:len(playerList[i].playerIp)-5]
			fmt.Println(GREEN, playerList[i].playerName, RESET, "[", CYAN, ip, RESET, "]")
		}
		fmt.Println()
	case "/add":
		addNewOpponents(targetList)
	}
}

func convertCoords(coords string) (int, int) {
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
	x := tab[rune(coords[0])]
	y := int(coords[1] - '0')
	return x, y
}

func launchAttack(target player, shooter player) {
	var attackCoords string
	for {
		fmt.Println("Ou souhaitez-vous tirer mon capitaiane ?")
		fmt.Scanf("%s", &attackCoords)
		if !VerifAttack(attackCoords) {
			fmt.Println(RED, "Coordonnés invalide !", RESET)
		} else {
			index := getPlayerMap(target.playerName)
			sendAttackCoords(target, attackCoords)
			j, i := convertCoords(attackCoords)
			if playerList[index].playermap.playermap[i][j] == '0' {
				fmt.Println(GREEN, string(attackCoords), ": Touché ", RESET)
			} else {
				fmt.Println(RED, string(attackCoords), ": Raté ", RESET)
			}
			break
		}
	}
}
