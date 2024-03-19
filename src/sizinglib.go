package src

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Switching(targetYbVersion string) {
	filePath := "resources/yb_" + strings.ReplaceAll(targetYbVersion, ".", "_") + ".db"

	if checkInternetAccess() {
		remoteFileExists, contents := checkFileExistsOnRemoteRepo(filePath)
		if remoteFileExists {
			// print the contents of the file
			fmt.Println(contents)
		} else {
			// check if local file exists
			isFileExist := checkLocalFileExists(filePath)
			if isFileExist {
				fmt.Println("file exist locally")
				// read the file
				cont, err := os.ReadFile(filePath)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(cont))
			} else {
				fmt.Println("file doesn't exist locally")
			}
		}
	} else {
		// no network access
		fmt.Println("No network access. Checking file locally...")
		// check if local file exists
		isFileExist := checkLocalFileExists(filePath)

		if isFileExist {
			fmt.Println("file exist locally")
			// read the file
			cont, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(cont))
		} else {
			fmt.Println("file doesn't exist locally")
		}
	}
}

func checkLocalFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func checkInternetAccess() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}
func checkFileExistsOnRemoteRepo(fileName string) (bool, string) {
	remotePath := "https://raw.githubusercontent.com/shaharuk-yb/sizing-calc/init/" + fileName
	resp, _ := http.Get(remotePath)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("File does not exist on remote location")
		return false, ""
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return true, string(body)
	}
}
