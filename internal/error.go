package internal

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ",err)
		os.Exit(1)
	}
}

func GiveError(msg string){
	fmt.Print(msg)
	os.Exit(1)
}