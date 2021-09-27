package generator

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

//go:embed tmpl.tpl
var tplString string

var tpl *template.Template

func firstToLow(s string) string {
	s = strcase.LowerCamelCase(s)
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func firstToUpper(s string) string {
	s = strcase.LowerCamelCase(s)
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func init() {

	tpl1, err := template.New("").Funcs(template.FuncMap{
		"firstToLow":   firstToLow,
		"firstToUpper": firstToUpper,
	}).Parse(tplString)
	if err != nil {
		panic("should never happen: " + err.Error())
	}
	tpl = tpl1
}

type table struct {
	TableName string   `json:"table"`
	Columns   []string `json:"columns"`
}

func GenerateFromSchema(pkg string, i, o string) error {

	bs, err := ioutil.ReadFile(i)
	if err != nil {
		return fmt.Errorf("read file %s failed:%w", i, err)
	}

	var tables []table
	if err := json.Unmarshal(bs, &tables); err != nil {
		return err
	}
	code, err := generate(pkg, tables)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(o, []byte(code), 0664)
}

func generate(pkg string, tables []table) (string, error) {

	data := map[string]interface{}{
		"stmtPkgName": "",
		"package": pkg,
		"tables":  tables,
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, data); err != nil {
		return "", err
	}
	return buff.String(), nil
}
