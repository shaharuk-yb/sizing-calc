/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
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
	sizerCmd.Flags().IntP("tables", "t", 1, "number of tables")
	sizerCmd.Flags().IntP("select-throughput", "s", 1, "desired select throughput")
	sizerCmd.Flags().IntP("insert-throughput", "i", 1, "desired isnert throughput")
	sizerCmd.Flags().StringP("target-yb-version", "y", "2.20", "target yugabyte db version")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sizerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sizerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func switching(targetYbVersion string) {
	filePath := "resources/yb_" + strings.ReplaceAll(targetYbVersion, ".", "_") + ".txt"

	fmt.Println(filePath)

	isFileExist := checkLocalFileExists(filePath)

	if isFileExist {
		fmt.Println("file exist")
	} else {

		fmt.Println("file doesn't exist")
	}

}

func checkLocalFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

/*func checkFileExistsOnRepo(fileName string) bool {

}*/
