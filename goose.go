package main

import (
	"bytes"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func writeGaggle(filenames []string) {
	file, err := os.OpenFile(".gaggle", os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)

	for i := 0; i < (len(filenames)); i++ {
		fmt.Println(filenames[i])
		file.WriteString(fmt.Sprintf("%s\n", filenames[i]))
	}

	file.Close()
}

func create() {
	now := time.Now()
	code := now.Format("20060102150405")
	filename := fmt.Sprintf("migrations/%s.sql", code)
	content := []byte("SQL goes here...")
	err := ioutil.WriteFile(filename, content, 0644)
	checkError(err)
}

func allMigrations() []string {
	files, err := ioutil.ReadDir("migrations")
	checkError(err)
	filenames := make([]string, len(files)-1)
	for i := 0; i < (len(files)); i++ {
		filenames = append(filenames, files[i].Name())
	}

	return filenames
}

func runMigrations(filenames []string) {
	config, err := toml.LoadFile(".gooserc")
	checkError(err)

	executable := config.Get("goose.path").(string)
	database := config.Get("goose.db").(string)
	hostname := config.Get("goose.hostname").(string)
	port := config.Get("goose.port").(string)
	name := config.Get("goose.name").(string)
	username := config.Get("goose.username").(string)
	password := config.Get("goose.password").(string)

	if database == "postgres" {
		executedFiles := make([]string, len(filenames)-1)

		for i := 0; i < (len(filenames)); i++ {
			cmd := exec.Command(executable, "-h", hostname, "-p", port, "-d", name, "-U", username, "-f", fmt.Sprintf("migrations/%s", strings.TrimLeft(filenames[i], " ")))
			cmd.Env = []string{fmt.Sprintf("PGPASSWORD=%s", password)}
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				checkError(err)
			} else {
				fmt.Printf("%q\n", out.String())
				executedFiles = append(executedFiles, filenames[i])
			}
		}

		writeGaggle(executedFiles)
	}
}

func migrate() {
	data, err := ioutil.ReadFile(".gaggle")
	lastMigration := string(data)

	if (err != nil) || (lastMigration == "") {
		fmt.Println("can't find gaggle")
		filenames := allMigrations()
		runMigrations(filenames)
	} else {
		fmt.Println("found gaggle")
		filenames := allMigrations()
		runMigrations(filenames)
	}
}

func main() {
	command := os.Args[1]

	if command == "migrate" {
		migrate()
	} else if command == "create" {
		create()
	}
}
