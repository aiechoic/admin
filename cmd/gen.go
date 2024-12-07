package main

import (
	"bytes"
	"flag"
	"fmt"
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
	flag.Parse()
	if *pkg == "" || *model == "" || *dir == "" {
		flag.PrintDefaults()
		return
	}
	err := os.MkdirAll(*dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	singularName := getSingularName(*model)
	pluralName := getPluralName(*model)
	var data = map[string]string{
		"package": *pkg,
		"model":   strings.Title(singularName),
		"models":  strings.Title(pluralName),
		"value":   singularName,
		"values":  pluralName,
	}

	executeTemplate(*tplDir, *dir, "handler", singularName, data)
	executeTemplate(*tplDir, *dir, "model", singularName, data)
	executeTemplate(*tplDir, *dir, "routes", singularName, data)
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

func getSingularName(name string) string {
	s := schema.NamingStrategy{
		SingularTable: true,
	}
	return s.TableName(name)
}

func getPluralName(name string) string {
	s := schema.NamingStrategy{
		SingularTable: false,
	}
	return s.TableName(name)
}
