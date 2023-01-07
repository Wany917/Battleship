package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func main() {
	music()
	fmt.Print("\033[1;1H\033[2J\n") // Clear the screen
	user := createPlayer()
	DisplayPlayerMap(user, true)
	if strings.Compare(user.playerName, "ERROR") != 0 && strings.Compare(user.playerIp, "ERROR") != 0 { //Si pas d'erreur lors de la création
		addPlayerToPlayerList(&user)        //Ajoute le joueur à la liste des joueur
		fmt.Println(playerList[0].playerIp) //Affiche son ip afin que ce soit plus simple de la connaitre
		http.HandleFunc("/list", listServers)
		http.HandleFunc("/added", beAdded)
		http.HandleFunc("/board", handlBoard)
		http.HandleFunc("/boats", handlboats)
		http.HandleFunc("/death", xDied)
		http.HandleFunc("/hit", attack)
		//http.HandleFunc("/test", test)

		go http.ListenAndServe(user.playerIp, nil)
		addOpponents() //Lance l'ajout d'adversaire via ip

		for {
			if !isAlive(user) {
				killPlayer(0)
				sendDeathInfo(user)
				break
			}
			/*if len(playerList) == 1 {
				killPlayer(0)
				break
			}*/
			userCommands()
			time.Sleep(2 * time.Second)
		}
		ranking(getPosition())
	} else {
		fmt.Println(RED, "ERREUR LORS DE LA CRÉATION DU JOUEUR", RESET)
	}
	Display_All_Board(playerList)
}
