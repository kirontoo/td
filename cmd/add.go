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
	task := strings.Join(args, " ")
	fmt.Printf("Added \"%s\" to your task list. \n", task)
}
