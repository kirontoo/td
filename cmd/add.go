package cmd

import (
	"fmt"

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

	var ids []int
	var failed []string
	for _, arg := range args {
		id, err := db.CreateTask(arg)
		if err != nil {
			failed = append(failed, arg)
		}
		ids = append(ids, id+1)
	}

	fmt.Println("Added new task(s):", util.SliceToString(ids))
	if len(failed) > 0 {
		fmt.Println("Failed to create task(s):")
		for _, t := range failed {
			fmt.Printf("* %s\n", t)
		}
	}
}
