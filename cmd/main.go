package main

import (
	"fmt"
	"os"
)

func usage(err ...error) {
	fmt.Println("help...")

	errString := ""
	for _, e := range err {
		if err != nil {
			errString = fmt.Sprintf("%s\n%s", err)
		}
	}

	if errString != "" {
		fmt.Println(errString)
		os.Exit(1)
	}

	os.Exit(0)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		usage()
	}

	cmd := os.Args[1]

	switch cmd {
	case "new":
		if len(args) < 3 {
			usage()
		}

		name := os.Args[2]
		if err := cmdNew(name); err != nil {
			usage(err)
		}
	}
}

func cmdNew(name string) error {
	fmt.Println("Creating new FRESH project: ", name)
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		return err
	}

	directories := []string{"cmd"}

	for _, dir := range directories {
		if err := os.Mkdir(name+"/"+dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
