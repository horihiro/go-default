package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/alexflint/go-arg"
	"github.com/dlclark/regexp2"
)

type DefaultValue struct {
	Id    string
	Value string
}

func Reduce(f func(string, DefaultValue) (string, error), a string, array []DefaultValue) (string, error) {
	var err error
	for _, x := range array {
		a, err = f(a, x)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
	}
	return a, nil
}

func main() {
	var args struct {
		TargetFolder    string            `arg:"-t,--target-folder,required" help:"target path of folder containing '.vscode/tasks.json'" placeholder:"/PATH/TO/PROJECT/FOLDER"`
		DefaultValues   map[string]string `arg:"-s,--set,required" help:"pairs of id and default values to update" placeholder:"id1=value1 ... idN=valueN"`
		BackupTasksJson string            `arg:"-b,--backup-file" help:"path of backup file of '.vscode/tasks.json'" placeholder:"/PATH/TO/tasks.json.bak"`
	}
	arg.MustParse(&args)

	defaultValues := make([]DefaultValue, 0, len(args.DefaultValues))
	for k, v := range args.DefaultValues {
		defaultValues = append(defaultValues, DefaultValue{Id: k, Value: v})
	}
	filePath := filepath.Join(args.TargetFolder, ".vscode", "tasks.json")
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	backupPath := args.BackupTasksJson
	if backupPath != "" {
		os.WriteFile(backupPath, file, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	var buf bytes.Buffer
	buf.WriteString("(?<=\"id\"\\s*:\\s*\"{{.}}\",?\\s*(?:\"(?:[a-z]+)\"\\s*:\\s*(?:\"[^\"]+\"|\\{[^}]+\\}|\\[[^\\]]+\\]),?\\s*)*\"default\"\\s*:\\s*\")[^\"]*(?=\".*)")
	buf.WriteString("|")
	buf.WriteString("(?<=\"default\"\\s*:\\s*\")[^\"]*(?=\",?\\s*(?:\"(?:[a-z]+)\"\\s*:\\s*(?:\"[^\"]+\"|\\{[^}]+\\}|\\[[^\\]]+\\]),?\\s*)*\"id\"\\s*:\\s*\"{{.}}\",?)")

	tmpl := buf.String()
	t, err := template.New("regexpPattern").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	text := string(file)
	m, err := Reduce(func(tasksJson string, defaultValue DefaultValue) (string, error) {
		Id := defaultValue.Id
		Value := defaultValue.Value

		var buf bytes.Buffer
		err := t.Execute(&buf, Id)
		if err != nil {
			log.Fatal(err)
		}
		re := regexp2.MustCompile(buf.String(), regexp2.RE2)
		m, err := re.Replace(tasksJson, Value, -1, -1)
		if err != nil {
			log.Fatal(err)
		}

		return m, nil
	}, text, defaultValues)
	if err != nil {
		log.Fatal(err)
	}

	d := []byte(m)
	err = os.WriteFile(filePath, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
