package cmd

import (
	"fmt"
	"strings"

	"github.com/kirontoo/td/db"
	"github.com/kirontoo/td/util"
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

	id, err := db.CreateTask(task)
	if err != nil {
		util.ExitIfErr(err)
	}

	fmt.Println("Added new task:", id)
}
