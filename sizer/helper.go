package sizer

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func convertToMap(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()

	// for each database row / record, a map with the column names and row values is added to the allMaps slice
	var allMaps []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		checkErr(err)
		resultMap := make(map[string]interface{})
		for i, val := range values {
			//fmt.Printf("Adding key=%s val=%v\n", columns[i], val)
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}
	return allMaps
}

func printMap(allMaps []map[string]interface{}) {
	for _, v := range allMaps {
		for k1, v1 := range v {
			fmt.Printf("%v : %v\n", k1, v1)
		}
		fmt.Println()
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkInternetAccess() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func checkLocalFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
