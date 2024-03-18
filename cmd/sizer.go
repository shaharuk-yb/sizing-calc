/*
Copyright © 2024 shaharuk-yb <sshaikh@yugabyte.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strings"
)

// sizerCmd represents the sizer command
var sizerCmd = &cobra.Command{
	Use:   "sizer",
	Short: "Generate sizing recommendations",
	Long: `Generate sizing recommendation 
based on inputs provided by the user.
NOTE: the recommendations can change based on the target yugabytedb version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sizer called")
		tables, _ := cmd.Flags().GetInt("tables")
		selectThroughput, _ := cmd.Flags().GetInt("select-throughput")
		insertThroughput, _ := cmd.Flags().GetInt("insert-throughput")
		targetYbVersion, _ := cmd.Flags().GetString("target-yb-version")
		fmt.Printf("user inputs:\n\ttables: %v\n\tselect_throughput: %v\n\tinsert_throughput: %v\n\ttarget yb version: %v\n", tables, selectThroughput, insertThroughput, targetYbVersion)
		switching(targetYbVersion)
	},
}

func init() {
	rootCmd.AddCommand(sizerCmd)

	// Here you will define your flags and configuration settings.
	sizerCmd.Flags().IntP("tables", "t", 1, "number of tables")
	sizerCmd.Flags().IntP("select-throughput", "s", 1, "desired select throughput")
	sizerCmd.Flags().IntP("insert-throughput", "i", 1, "desired isnert throughput")
	sizerCmd.Flags().StringP("target-yb-version", "y", "2.20", "target yugabyte db version")
}

func switching(targetYbVersion string) {
	filePath := "resources/yb_" + strings.ReplaceAll(targetYbVersion, ".", "_") + ".txt"

	remoteFileExists, contents := checkFileExistsOnRemoteRepo(filePath)

	if remoteFileExists {
		// print the contents of the file
		fmt.Println(contents)
	} else {
		//check if local file exists
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
