/*
Copyright © 2024 shaharuk-yb <sshaikh@yugabyte.com>
*/
package cmd

import (
	"fmt"
	"github.com/shaharuk-yb/sizing-calc/src"
	"github.com/spf13/cobra"
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
		src.Switching(targetYbVersion)
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
