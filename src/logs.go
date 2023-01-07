package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

var LogIp = "192.168.0.0:9000"

func addLog(entry string) error {
	// Création de la requête POST
	data := url.Values{}
	data.Set("entry", entry)
	req, err := http.NewRequest("POST", LogIp+"/logs", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Envoi de la requête et récupération de la réponse
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Lecture de la réponse
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
