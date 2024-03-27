package sizer

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

func Run(targetYbVersion string, inputs map[string]int) {
	// read required inputs: may change from version to version
	tables := inputs["tables"]
	//requiredSelectThroughput := inputs["requiredSelectThroughput"]
	//requiredInsertThroughput := inputs["requiredInsertThroughput"]

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
	//printRows()
	checkTableLimits(tables)
	//getThroughputData(2, requiredInsertThroughput, requiredSelectThroughput)
}

func printRows() {
	rows, err := DB.Query("SELECT * from sizing limit 10")
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
	remotePath := "https://raw.githubusercontent.com/shaharuk-yb/sizing-calc/maps/" + fileName
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

func checkTableLimits(req_tables int) {
	rows, err := DB.Query("select num_cores from sizing where num_tables > ? and dimension like '%TableLimits-3nodeRF=3%' order by num_cores", req_tables)
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

func getThroughputData(minCoresReq int, requiredInsertThroughput int, requiredSelectThroughput int) {
	//rows, err := DB.Query("select foo.* from (select id, (cast(?/inserts_per_core) as int + ((?/inserts_per_core) > cast(?/inserts_per_core) as int)) insert_total_cores, (cast(?/selects_per_core) as int + ((?/selects_per_core) > cast(?/selects_per_core) as int) select_total_cores, num_cores, num_nodes from sizing where dimension='MaxThroughput' and num_cores>=?) as foo order by select_total_cores + insert_total_cores, num_cores", requiredInsertThroughput, requiredInsertThroughput, requiredInsertThroughput, requiredSelectThroughput, requiredSelectThroughput, requiredSelectThroughput, minCoresReq)
	rows, err := DB.Query("select foo.* from (select id, ? , ?, num_cores, num_nodes from sizing where dimension='MaxThroughput' and num_cores>=?) as foo order by select_total_cores, insert_total_cores, num_cores", requiredInsertThroughput, requiredSelectThroughput, minCoresReq)
	if err != nil {
		fmt.Println("no records found")
	}
	defer rows.Close()
	err = rows.Err()

	allMaps := convertToMap(rows)
	printMap(allMaps)
	if err != nil {
		fmt.Println("error occurred")
	}
}
