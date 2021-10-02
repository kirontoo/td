package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to your lask list",
	Run:   runAddCmd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAddCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Missing task argument")
		return
	}

	task := strings.Join(args, " ")
	fmt.Printf("Added \"%s\" to your task list. \n", task)

	// if err != nil {
	// 	// fmt.Errorf("Could not save new task %s", task)
	// 	fmt.Printf("Could not save new task:\n %s", task)
	// }
}
