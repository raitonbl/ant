package cmd

import (
	"fmt"
	"github.com/raitonbl/ant/internal"
	"os"
)

func exit(err *internal.Problem) {
	fmt.Println(err.Message)
	os.Exit(err.Code)
}