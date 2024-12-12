package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jinzhu/inflection"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	pkg := flag.String("p", "", "package name")
	model := flag.String("m", "", "model name")
	dir := flag.String("d", "", "destination directory")
	tplDir := flag.String("td", "templates/gen", "template directory, default is 'templates/gen'")
	// Customize the usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "gen")
		fmt.Fprintln(os.Stderr, "This program generates code based on the provided model and package name.")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "gen -p mypackage -m MyModel -d ./output")
	}
	flag.Parse()
	if *pkg == "" || *model == "" || *dir == "" {
		flag.PrintDefaults()
		return
	}
	err := os.MkdirAll(*dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	singularSnakeCaseName := getSnakeCaseName(*model, true)
	pluralSnakeCaseName := getSnakeCaseName(*model, false)

	singularCamelCaseName := getCamelCaseName(*model, true)
	pluralCamelCaseName := getCamelCaseName(*model, false)

	singularCamelCaseTitle := strings.Title(singularCamelCaseName)
	pluralCamelCaseTitle := strings.Title(pluralCamelCaseName)

	var data = map[string]string{
		"package":         *pkg,
		"snakeCaseModel":  singularSnakeCaseName,
		"snakeCaseModels": pluralSnakeCaseName,
		"camelCaseModel":  singularCamelCaseName,
		"camelCaseModels": pluralCamelCaseName,
		"CamelCaseModel":  singularCamelCaseTitle,
		"CamelCaseModels": pluralCamelCaseTitle,
	}

	executeTemplate(*tplDir, *dir, "handler", singularSnakeCaseName, data)
	executeTemplate(*tplDir, *dir, "model", singularSnakeCaseName, data)
	executeTemplate(*tplDir, *dir, "routes", singularSnakeCaseName, data)
}

func executeTemplate(tplDir, dir, file, model string, data map[string]string) {
	tplFile := filepath.Join(tplDir, fmt.Sprintf("%s.tmpl", file))
	tpl := template.Must(template.ParseFiles(tplFile))
	filename := filepath.Join(dir, fmt.Sprintf("%s_%s.go", model, file))
	_, err := os.Stat(filename)
	if err == nil {
		log.Printf("file '%s' already exists\n", filename)
		return
	}
	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, data)
	if err != nil {
		panic(fmt.Errorf("failed to execute template '%s': %w", tplFile, err))
	}
	err = os.WriteFile(filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
	logrus.Printf("created file '%s'\n", filename)
}

func getSnakeCaseName(model string, singular bool) string {
	s := schema.NamingStrategy{
		SingularTable: singular,
	}
	return s.TableName(model)
}

func getCamelCaseName(model string, singular bool) string {
	if singular {
		return inflection.Singular(model)
	} else {
		return inflection.Plural(model)
	}
}
