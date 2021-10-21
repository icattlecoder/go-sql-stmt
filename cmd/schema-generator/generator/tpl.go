package generator

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
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

type Column struct {
	Name     string `json:"name"`
	Sortable bool   `json:"sortable,omitempty"`
}

type table struct {
	TableName   string   `json:"table"`
	HasSortable bool     `json:"-"`
	Columns     []Column `json:"columns"`
}

func (t *table) preDo() {

	for _, c := range t.Columns {
		if c.Sortable {
			t.HasSortable = true
			break
		}
	}
}

type tables []table

func (t tables) Len() int {
	return len(t)
}

func (t tables) Less(i, j int) bool {
	return t[i].TableName < t[j].TableName
}

func (t tables) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
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

func generate(pkg string, tables tables) (string, error) {

	sort.Sort(tables)
	for i := range tables {
		tables[i].preDo()
	}
	data := map[string]interface{}{
		"stmtPkgName": "",
		"package":     pkg,
		"tables":      tables,
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, data); err != nil {
		return "", err
	}
	return buff.String(), nil
}
