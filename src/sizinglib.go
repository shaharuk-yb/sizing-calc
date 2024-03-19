package src

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"net/http"
	"os"
	"strings"
)

var baseDownloadPath = "resources/remote/"
var DB *sql.DB

func Switching(targetYbVersion string) {
	filePath := "resources/yb_" + strings.ReplaceAll(targetYbVersion, ".", "_") + ".db"
	if checkInternetAccess() {
		remoteFileExists := checkFileExistsOnRemoteRepo(filePath)
		if remoteFileExists {
			// print the contents of the file
			//fmt.Println(contents)
			filePath = strings.ReplaceAll(filePath, "resources/", baseDownloadPath)
			fmt.Println("connect to downloaded data")
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
			panic("file doesn't exist locally")
		}
	}
	err := ConnectDatabase(filePath)
	checkErr(err)
	printRows()
}

func printRows() {
	rows, err := DB.Query("SELECT id, dimension from sizing limit 10")
	if err != nil {
		fmt.Println("no records found")
	}
	defer rows.Close()
	allMaps := convertToMap(rows)
	printMap(allMaps)

	err = rows.Err()
	if err != nil {
		fmt.Println("error occurred")
	}
}

func checkFileExistsOnRemoteRepo(fileName string) bool {
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
		return false
	} else {
		//body, err := io.ReadAll(resp.Body)
		downloadPath := strings.ReplaceAll(fileName, "resources/", baseDownloadPath)
		out, err := os.Create(downloadPath)
		defer out.Close()
		_, err = io.Copy(out, resp.Body)

		if err != nil {
			panic(err)
		}
		return true
	}
}

func ConnectDatabase(file string) error {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	DB = db
	return nil
}
