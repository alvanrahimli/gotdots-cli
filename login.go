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

	fmt.Printf("Backend url: %s\n", backendUrl)
	response, loginErr := http.Post(backendUrl, "application/json", bytes.NewReader(requestBody))
	if loginErr != nil {
		fmt.Printf("ERROR: %s\n", loginErr.Error())
		panic(loginErr)
		// return
	}
	//goland:noinspection ALL
	defer response.Body.Close()

	// Handle failed response
	if response.StatusCode == 401 {
		fmt.Printf("Could not login with given credentials. Hint: %s\n", response.Status)
		return
	}

	body, bodyReadErr := io.ReadAll(response.Body)
	if bodyReadErr != nil {
		fmt.Printf("ERROR: %s\n", bodyReadErr.Error())
		// panic(bodyReadErr)
		return
	}

	var loginResponse models.LoginResponse
	unmarshallErr := json.Unmarshal(body, &loginResponse)
	if unmarshallErr != nil {
		fmt.Printf("ERROR: %s\n", unmarshallErr.Error())
		// panic(unmarshallErr)
		return
	}

	// Write token to file as: "Bearer <token>"
	folder, _ := getArchivesFolder()
	tokenFile := fmt.Sprintf("%s/token", folder)
	writeErr := utils.WriteToFile(tokenFile, fmt.Sprintf("Bearer %s", loginResponse.Token))
	if writeErr != nil {
		fmt.Println("Could not save token")
		fmt.Printf("ERROR: %s\n", writeErr.Error())
		return
	}

	fmt.Println("Successfully logged in")
}
