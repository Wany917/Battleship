package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func playMusic() {
	cmd := exec.Command("mpg123", "-q", "music/navy_music.mp3", "-loop", "100000000000")
	err := cmd.Start()
	if err != nil {
		fmt.Println("faill music")
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("faill music")
		log.Fatal(err)
	}
}

func music() {
	fmt.Print("\033[1;1H\033[2J\n") // Clear the screen
	fmt.Print("Music [N/Y] ? ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ifMusic := scanner.Text()

	if strings.ToLower(ifMusic) == "y" {
		go playMusic()
		return
	} else if strings.ToLower(ifMusic) == "n" {
		return
	} else {
		fmt.Println(RED, "y TO PLAY MUSIC AND n FOR NOT\n", RESET)
		music()
	}
}
