package main

import (
	"bytes"
	"log"
	"os"
	"strings"
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

func Mapping(f func(string) DefaultValue, array []string) []DefaultValue {
	buff := make([]DefaultValue, len(array))
	for i, v := range array {
		buff[i] = f(v)
	}
	return buff
}

func main() {
	var args struct {
		TargetFile    string   `arg:"-t,--target-file,required" help:"file path of target 'tasks.json' "`
		DefaultValues []string `arg:"-v,--default-value,separate,required" help:"id and default values to update, the format is '${ID}:${DEFAULT_VALUE}'"`
	}
	arg.MustParse(&args)

	filePath := args.TargetFile
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defaultValues := Mapping(func(s string) DefaultValue {
		strs := strings.Split(s, "=")
		return DefaultValue{Id: strs[0], Value: strs[1]}
	}, args.DefaultValues)

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

	// 2. バイト文字列に変換
	d := []byte(m)

	// 3. 書き込み
	err = os.WriteFile(filePath, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
