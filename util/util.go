package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kirontoo/td/db"
)

func ExitIfErr(err error) {
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		db.Close()
		os.Exit(1)
	}
}

func SliceToString(s []int) string {
	text := []string{}
	for i := range s {
		number := s[i]
		nStr := strconv.Itoa(number)
		text = append(text, nStr)
	}

	return strings.Join(text, " ")
}
