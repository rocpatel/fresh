package main

import (
	"fmt"
	"os"
	"os/exec"
)

func usage(err ...error) {
	fmt.Println("help...")

	errString := ""
	for _, e := range err {
		if err != nil {
			errString = fmt.Sprintf("%s\n%s", errString, e)
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

	directories := []string{"cmd", "handler", "view", "view/hello", "view/layout", "model"}

	for _, dir := range directories {
		if err := os.Mkdir(name+"/"+dir, os.ModePerm); err != nil {
			return err
		}
	}

	files := map[string][]byte{
		name + "/go.mod":                 genGoModContent(name),
		name + "/cmd/main.go":            genCmdMainContent(name),
		name + "/handler/hello.go":       genIndexHandlerContent(name),
		name + "/view/layout/base.templ": genBaseLayoutContent(),
		name + "/view/hello/index.templ": genIndexContent(name),
	}

	for file, content := range files {
		if err := os.WriteFile(file, content, os.ModePerm); err != nil {
			return err
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(name); err != nil {
		return err
	}

	fmt.Println("installing project dependencies and generating templ files....")

	if err := exec.Command("go", "get", "github.com/rocpatel/fresh").Run(); err != nil {
		return err
	}

	fmt.Println("installing project dependencies templ....")
	if err := exec.Command("go", "get", "github.com/a-h/templ").Run(); err != nil {
		return err
	}

	if out, err := exec.Command("templ", "generate").Output(); err != nil {
		fmt.Println(string(out))
		return err
	}

	if err := os.Chdir(cwd); err != nil {
		return err
	}

	return nil
}

func genGoModContent(mod string) []byte {
	return []byte(fmt.Sprintf(`
module %s

go 1.22.0
	
`, mod))
}

func genCmdMainContent(name string) []byte {
	return []byte(fmt.Sprintf(`
package main

import (
	"github.com/rocpatel/fresh"
	"%s/handler"
)

func main() {
	app := fresh.New()
	app.Get("/",handler.HandleHelloIndex)
	app.Start(":3000")
}
	
`, name))
}

func genBaseLayoutContent() []byte {
	return []byte(`
package layout

templ Base() {
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>Slick Application</title>
	</head>
	<body>
		{ children... }
	</body>
</html>
}
`)
}

func genIndexContent(name string) []byte {
	return []byte(fmt.Sprintf(`
package hello

import (
	"%s/view/layout"
)

templ Index() {
	@layout.Base() {
		<h1>hello there</h1>
	}
}
`, name))
}

func genIndexHandlerContent(name string) []byte {
	return []byte(fmt.Sprintf(`
package handler

import (
	"github.com/rocpatel/fresh"
	"%s/view/hello"
)

func HandleHelloIndex(c *fresh.Context) error {
	return c.Render(hello.Index())
}
`, name))
}
