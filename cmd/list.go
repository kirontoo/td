/*
Copyright © 2021 kirontoo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kirontoo/td/db"
	"github.com/kirontoo/td/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all current uncompleted tasks",
	Run:   runListCmd,
}

var completed bool
var all bool

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&completed, "completed", "c", false, "List completed tasks")
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "List all tasks")
}

func runListCmd(cmd *cobra.Command, args []string) {
	tasks, err := db.GetUncompletedTasks()
	if err != nil {
		util.ExitIfErr(err)
	}

	if len(tasks) == 0 {
		pterm.Info.Println("No tasks!")
		return
	}

	if completed {
		tasks, err = db.GetCompletedTasks()
		if len(tasks) == 0 {
			pterm.Info.Println("Could not find completed tasks")
		}
	} else if all {
		tasks, err = db.GetAllTasks()
		handleError(err)
	}

	for _, task := range tasks {
		fmt.Printf("%d. %s\n", task.Key+1, task.Value)
	}
}

func handleError(err error) {
	if err != nil {
		pterm.Error.Println("Could not find completed tasks")
		os.Exit(1)
	}
}
