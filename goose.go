package main

import (
	"bytes"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func writeGaggle(filenames []string) {
	joinedFiles := strings.Join(filenames, "\n")
	newLineSlice := []string{joinedFiles, "\n"}
	content := strings.Join(newLineSlice, "")
	file, err := os.OpenFile(".gaggle", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	checkError(err)
	file.WriteString(content)
	file.Close()
}

func conditionallyCreateMigrationsDirectory() {
	path := "migrations/"
	_, err := os.Stat(path)

	if err != nil {
		os.Mkdir(path, 0755)
	}
}

func create() {
	now := time.Now()
	code := now.Format("20060102150405")
	filename := fmt.Sprintf("migrations/%s.sql", code)
	content := []byte("SQL goes here...")
	conditionallyCreateMigrationsDirectory()
	err := ioutil.WriteFile(filename, content, 0644)
	checkError(err)
}

func allMigrations() []string {
	files, err := ioutil.ReadDir("migrations")
	checkError(err)
	filenames := make([]string, 0)
	for i := 0; i < (len(files)); i++ {
		filenames = append(filenames, files[i].Name())
	}

	return filenames
}

func lastMigration(gaggle string) string {
	if gaggle != "" {
		runMigrations := strings.Split(gaggle, "\n")
		return runMigrations[len(runMigrations)-1]
	} else {
		return ""
	}
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
		executedFiles := make([]string, 0)

		for i := 0; i < (len(filenames)); i++ {
			fmt.Println(fmt.Sprintf("Running Migration: %v", filenames[i]))
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
				fmt.Println(out.String())
				executedFiles = append(executedFiles, filenames[i])
			}
		}

		writeGaggle(executedFiles)
	}
}

func migrate() {
	data, err := ioutil.ReadFile(".gaggle")
	lastMigration := lastMigration(strings.Trim(string(data), "\n"))

	if (err != nil) || (lastMigration == "") {
		unrunMigrations := allMigrations()
		runMigrations(unrunMigrations)
		fmt.Println("Migrations complete.")
	} else {
		filenames := allMigrations()
		i := sort.SearchStrings(filenames, lastMigration)

		if (i + 1) == len(filenames) {
			fmt.Println("Migrations up to date.")
		} else {
			unrunMigrations := filenames[(i + 1):]
			fmt.Println(unrunMigrations)
			runMigrations(unrunMigrations)
			fmt.Println("Migrations complete.")
		}
	}
}

func checkForRc() {
	_, err := os.Stat(".gooserc")

	if err != nil {
		fmt.Println("Can't find .gooserc. Please create it with with your database connection settings and try again.")
		os.Exit(0)
	}
}

func main() {
	checkForRc()
	command := os.Args[1]

	if command == "migrate" {
		migrate()
	} else if command == "hatch" {
		create()
	}
}
