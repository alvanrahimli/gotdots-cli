package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"io"
	"net/http"
	"os"
	"path"
)

func login() {
	var username, password string

	fmt.Print("Enter your username: ")
	_, scanErr := fmt.Scanln(&username)
	if scanErr != nil {
		return
	}

	fmt.Print("Enter your password: ")
	_, scanErr = fmt.Scanln(&password)
	if scanErr != nil {
		return
	}

	backendUrl := os.Getenv("BACKEND_URL")
	requestBody, marshallErr := json.Marshal(models.LoginDto{
		Username: username,
		Password: password,
	})
	if marshallErr != nil {
		fmt.Printf("ERROR: %s\n", marshallErr.Error())
		return
	}

	response, loginErr := http.Post(backendUrl, "application/json", bytes.NewReader(requestBody))
	if loginErr != nil {
		fmt.Printf("ERROR: %s\n", loginErr.Error())
		panic(loginErr)
		// return
	}
	//goland:noinspection ALL
	defer response.Body.Close()

	// Handle failed response
	if response.StatusCode == http.StatusUnauthorized {
		fmt.Printf("Could not login with given credentials. Hint: %s\n", response.Status)
		return
	}

	body, bodyReadErr := io.ReadAll(response.Body)
	if bodyReadErr != nil {
		fmt.Printf("ERROR: %s\n", bodyReadErr.Error())
		return
	}

	var loginResponse models.LoginResponse
	unmarshallErr := json.Unmarshal(body, &loginResponse)
	if unmarshallErr != nil {
		fmt.Printf("ERROR: %s\n", unmarshallErr.Error())
		return
	}

	// Write token to file as: "Bearer <token>"
	archivesFolder, _ := getArchivesFolder()
	tokenFile := path.Join(archivesFolder, ".token")
	authHeaderValue := fmt.Sprintf("Bearer %s", loginResponse.Token)
	writeErr := utils.WriteToFile(tokenFile, authHeaderValue)
	if writeErr != nil {
		fmt.Println("Could not save token")
		fmt.Printf("ERROR: %s\n", writeErr.Error())
		return
	}

	// Get userinfo & write to $HOME/.dots-archives/.userinfo
	infoUrl := os.Getenv("GET_USERINFO_URL")
	client := http.Client{}
	req, _ := http.NewRequest("GET", infoUrl, nil)
	req.Header.Set("Authorization", authHeaderValue)
	res, resErr := client.Do(req)
	if resErr != nil {
		handleError(resErr, true)
	}

	if res.StatusCode != http.StatusOK {
		// We should not reach here
		fmt.Println("Wrong credentials")
		os.Exit(1)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		handleError(err, true)
	}

	infoFile := path.Join(archivesFolder, ".userinfo")
	writeErr = utils.WriteToFile(infoFile, string(bodyBytes))
	if writeErr != nil {
		fmt.Println("Could not save userinfo")
		handleError(writeErr, true)
	}

	fmt.Println("Successfully logged in and saved userinfo")
}
