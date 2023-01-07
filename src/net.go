package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var port = ":9000"

func getBoatsNumber(name string) {
	usrIndex := 0
	ip := "NULL"
	for usrIndex = 0; usrIndex != len(playerList); usrIndex++ {
		if playerList[usrIndex].playerName == name {
			ip = playerList[usrIndex].playerIp
			break
		}
	}
	if ip == "NULL" {
		fmt.Println(RED, name, "ERREUR PLAYER_NAME", RESET)
		return
	}
	response, err := http.Get("http://" + ip + "/boats") //Envoie une requete GET à l'IP entré par l'utilisateur
	if err != nil {                                      //Si il y a une erreur
		fmt.Println(RED, name, "ERREUR ADRESSE INVALIDE OR USER_NAME", RESET)
		return
	}
	body, err := ioutil.ReadAll(response.Body) // Si la requete fonctionne récupère le contenu de la page (dans la page il n'y a que l'ip et le nom du joueur adverse)
	if err != nil {                            //Si il y a une erreur, affiche l'erreur et relance la boucle
		fmt.Println(RED, "ERREUR LECTURE DE LA PAGE", RESET)
		return
	}
	response.Body.Close() //Obliger de fermer le body
	boats := string(body)
	fmt.Println("Il lui rest [", boats, "] navire(s) encore debut.\n")
}

func getPlayerMap(name string) int {
	usrIndex := 0
	ip := "NULL"
	for usrIndex = 0; usrIndex != len(playerList); usrIndex++ {
		if playerList[usrIndex].playerName == name {
			ip = playerList[usrIndex].playerIp
			break
		}
	}
	if ip == "NULL" {
		fmt.Println(RED, name, "ERREUR USER_NAME\n", RESET)
		return -1
	}
	response, err := http.Get("http://" + ip + "/board") //Envoie une requete GET à l'IP entré par l'utilisateur
	if err != nil {                                      //Si il y a une erreur
		fmt.Println(RED, "ERREUR ADRESSE INVALIDE OR USER_NAME", RESET)
		return -1
	}
	body, err := ioutil.ReadAll(response.Body) // Si la requete fonctionne récupère le contenu de la page (dans la page il n'y a que l'ip et le nom du joueur adverse)
	if err != nil {                            //Si il y a une erreur, affiche l'erreur et relance la boucle
		fmt.Println(RED, "ERREUR LECTURE DE LA PAGE", RESET)
		return -1
	}
	response.Body.Close() //Obliger de fermer le body
	playermap := string(body)
	playerList[usrIndex].playermap.playermap = make([][]rune, rows)
	for i := range playerList[usrIndex].playermap.playermap {
		playerList[usrIndex].playermap.playermap[i] = make([]rune, cols)
		for j := 0; j != cols; j++ {
			playerList[usrIndex].playermap.playermap[i][j] = '~'
		}
	}
	for count, i := 0, 0; i != cols; i++ {
		for j := 0; j != rows; j++ {
			playerList[usrIndex].playermap.playermap[i][j] = rune(playermap[count])
			count++
		}
		count++
	}
	return usrIndex
}

func addNewOpponents(targetList []Target) {
	usrInput := targetList[0].targetName
	fmt.Printf("\n%s\n", usrInput)
	response, err := http.Get("http://" + usrInput + port + "/added") //Envoie une requete GET à l'IP entré par l'utilisateur
	if err != nil {                                                   //Si il y a une erreur
		fmt.Println(RED, "ERREUR ADRESSE INVALIDE", RESET)
		return
	}
	body, err := ioutil.ReadAll(response.Body) // Si la requete fonctionne récupère le contenu de la page (dans la page il n'y a que l'ip et le nom du joueur adverse)
	if err != nil {                            //Si il y a une erreur, affiche l'erreur et relance la boucle
		fmt.Println(RED, "ERREUR LECTURE DE LA PAGE", RESET)
		return
	}
	response.Body.Close()                   //Obliger de fermer le body
	ToSplit := string(body)                 //stocke la réponse récupérer dans le body
	Splitted := strings.Split(ToSplit, "/") //Sépare le résultat en 2 slice pour avoir le nom et l'ip séparé

	fmt.Println(Splitted[1])
	opponent := player{playerIp: Splitted[0], playerName: Splitted[1]} // crée un jouer en assignant le nom et l'ip récup plus tôt
	if !checkPlayersList(opponent.playerIp) {                          //Vérifie que le joueur n'a pas été ajouté,si il n'a pas été ajouté :
		addPlayerToPlayerList(&opponent) //ajoute le joueur à la liste de joueurs
		fmt.Println("Ennemis repéré mon capitaine ! Nous l'avons ajouter au radar ! (/list pour le radar) ")
	} else { //	 sinon relance la boucle
		fmt.Println(RED, "Impossible d'ajouter cet ennemi au radar !", RESET)
		fmt.Println(RED, "Il se peut qu'il y soit déja ! (/list pour le radar)", RESET)
	}
}

func addOpponents() { //permet d'ajouter des adversaires
	var usrInput string
	fmt.Println("Entrez l'adresse de vos adversaires (STOP pour arreter): ")
	for strings.Compare(usrInput, "STOP") != 0 { //Tant que l'user n'entre pas STOP, continue de demander des IP d'adversaires
		fmt.Scanf("%s", &usrInput)
		if strings.Compare(usrInput, "STOP") == 0 {
			break
		}
		response, err := http.Get("http://" + usrInput + port + "/added") //Envoie une requete GET à l'IP entré par l'utilisateur
		fmt.Println(usrInput)
		if err != nil { //Si il y a une erreur, relance la boucle
			fmt.Println(RED, "ERREUR ADRESSE INVALIDE", RESET)
			fmt.Println("Entrez l'adresse de vos adversaires (STOP pour arreter): ")
			continue
		}
		body, err := ioutil.ReadAll(response.Body) // Si la requete fonctionne récupère le contenu de la page (dans la page il n'y a que l'ip et le nom du joueur adverse)
		if err != nil {                            //Si il y a une erreur, affiche l'erreur et relance la boucle
			fmt.Println(RED, "ERREUR LECTURE DE LA PAGE", RESET)
			fmt.Println("Entrez l'adresse de vos adversaires (STOP pour arreter): ")
			continue
		}
		response.Body.Close()                   //Obliger de fermer le body
		ToSplit := string(body)                 //stocke la réponse récupérer dans le body
		Splitted := strings.Split(ToSplit, "/") //Sépare le résultat en 2 slice pour avoir le nom et l'ip séparé

		fmt.Println(Splitted[1])
		opponent := player{playerIp: Splitted[0], playerName: Splitted[1]} // crée un jouer en assignant le nom et l'ip récup plus tôt
		if !checkPlayersList(opponent.playerIp) {                          //Vérifie que le joueur n'a pas été ajouté,si il n'a pas été ajouté :
			addPlayerToPlayerList(&opponent) //ajoute le joueur à la liste de joueurs
			fmt.Println("Ennemis repéré mon capitaine ! Nous l'avons ajouter au radar ! (/list pour le radar) ")
			fmt.Println("Entrez l'adresse de vos adversaires (STOP pour arreter): ")
			continue //relance la boucle
		} else { // sinon relance la boucle
			fmt.Println(RED, "Impossible d'ajouter cet ennemi au radar !", RESET)
			fmt.Println(RED, "Il se peut qu'il y soit déja ! (/list pour le radar)", RESET)
			fmt.Println("Entrez l'adresse de vos adversaires (STOP pour arreter): ")
			continue
		}
	}
}
func beAdded(w http.ResponseWriter, r *http.Request) { //renvoie l'IP et le nom du joueur afin qu'il puisse être ajouté
	if r.URL.Path != "/added" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, playerList[0].playerIp+"/"+playerList[0].playerName) //Affiche le joueur 0 dans le tableau des joueurs car il s'agit
	} //toujours du propriétaire du serveur

}

func handlboats(w http.ResponseWriter, r *http.Request) {
	boats := 0
	if r.URL.Path != "/boats" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		for i := 0; i != len(playerList[0].playermap.Ships); i++ {
			if playerList[0].playermap.Ships[i].isSunken == false {
				boats++
			}
		}
		fmt.Fprintf(w, strconv.Itoa(boats))
	}
}

func handlBoard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/board" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		updateMap(&playerList[0].playermap)
		for i := 0; i != rows; i++ {
			fmt.Fprintf(w, string(playerList[0].playermap.playermap[i]))
			fmt.Fprintf(w, "\n")
		}
	}
}

func addPlayerToPlayerList(user *player) { //Ajoute un joueur à la liste de joueurs
	playerList = append(playerList, *user)
}

func listServers(w http.ResponseWriter, r *http.Request) { // Liste les joueur mais n'a servi qua des test et peut être enlevé
	if r.URL.Path != "/list" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "User / Ip")
		printPlayerList()
	}
}

func LocalIP() (net.IP, error) { //Renvoie l'ip local du joueur
	ifaces, err := net.Interfaces() //Récupère les interfaces de l'utilisateurs (ex:eth1,ens0...)
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces { //Pour chaque interface
		addrs, err := i.Addrs() //Récupère les ip des interfaces
		if err != nil {         //Si erreur STOP
			return nil, err
		}

		for _, addr := range addrs { //Pour chaque IP récupéré
			var ip net.IP
			switch v := addr.(type) { //En fonction du type de l'addresse
			case *net.IPNet: //si c'est une adresse Réseau la récupère
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP //Si c'est une IP la récupère
			}

			if isPrivateIP(ip) { //Vérifie que c'est une ip Privé, si c'est le cas la renvoie sinon
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP") //Renvoie PAS D'IP
}

func isPrivateIP(ip net.IP) bool { //Vérifie qu'il s'agit d'une IP privé
	var privateIPBlocks []*net.IPNet //Tableau d'IP Réseau
	for _, cidr := range []string{   //Vérifie l'IP, Si elle fait partie d'un réseau privé
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)               //Récupère l'IP et l'IP réseau quel implique
		privateIPBlocks = append(privateIPBlocks, block) //l'ajoute au tableau des IP
	}

	for _, block := range privateIPBlocks { //Si il y a une ip sa veut dire qu'elle est privé donc renvoie vraie
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

func xDied(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil { // Parsing des paramètres envoyés
			fmt.Println(RED, "Something went bad", RESET) // par le client et gestion // d’erreurs
			fmt.Fprintln(w, "Something went bad")
			return
		}
		name := r.FormValue("name")
		ip := r.FormValue("ip")
		var userToKill int
		for i := 0; i < len(playerList); i++ {
			if strings.Compare(name, playerList[i].playerName) == 0 && strings.Compare(ip, playerList[i].playerIp) == 0 {
				userToKill = i
				break
			}
		}
		killPlayer(userToKill)
		fmt.Fprintf(w, "%s à été tué ! Il reste %d ennemis", name, len(playerList))
	}
}

func sendDeathInfo(user player) {
	data := url.Values{}
	data.Add("name", user.playerName)
	data.Add("ip", user.playerIp)
	for i := 0; i < len(playerList); i++ {
		http.PostForm(playerList[i].playerIp+port+"/death", data)
	}
}

func sendAttackCoords(target player, coords string) {
	data := url.Values{}
	data.Add("coords", coords)
	http.PostForm("http://"+target.playerIp+"/hit", data)
}

func handllogs(w http.ResponseWriter, r *http.Request) {
	var entry = "test"
	switch r.Method {
	case http.MethodPost:
		data := url.Values{}
		data.Set("entry", entry)
		req, err := http.NewRequest("POST", LogIp+"/logs", bytes.NewBufferString(data.Encode()))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Envoi de la requête et récupération de la réponse
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		// Lecture de la réponse
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		return
	}
}

func attack(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		if err := r.ParseForm(); err != nil { // Parsing des paramètres envoyés
			fmt.Println(RED, "Something went bad", RESET) // par le client et gestion // d’erreurs
			fmt.Fprintln(w, "Something went bad")
			return
		}
		coords := r.FormValue("coords")
		playerList[0], _ = attackPlayer(coords, playerList[0])
		shikShip(&playerList[0].playermap)
	}
}

//func test(w http.ResponseWriter, r *http.Request, player *player) player {
//	return *player
//}
