package util

import (
	"fmt"
	"os"

	"github.com/kirontoo/td/db"
)

func ExitIfErr(err error) {
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		db.Close()
		os.Exit(1)
	}
}
