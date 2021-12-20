package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program /path/to/env/dir command [arg1 arg2 ...]")
		log.Fatal("not enough arguments")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}
