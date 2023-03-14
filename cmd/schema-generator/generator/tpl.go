package generator

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

//go:embed tmpl.tpl
var tplString string

var schemaTpl *template.Template

//go:embed go-proto-tmpl.tpl
var goProtoTplString string

var goProtoTpl *template.Template

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

	tpl1, err := template.New("t1").Funcs(template.FuncMap{
		"firstToLow":   firstToLow,
		"firstToUpper": firstToUpper,
	}).Parse(tplString)
	if err != nil {
		panic("should never happen: " + err.Error())
	}
	schemaTpl = tpl1

	tpl1, err = template.New("t2").Funcs(template.FuncMap{
		"firstToLow":   firstToLow,
		"firstToUpper": firstToUpper,
	}).Parse(goProtoTplString)
	if err != nil {
		panic("should never happen: " + err.Error())
	}
	goProtoTpl = tpl1
}

type Column struct {
	Name     string `json:"name"`
	Update   bool   `json:"update,omitempty"`
	GoType   string `json:"gotype,omitempty"`
	Sortable bool   `json:"sortable,omitempty"`
	Comment  string `json:"comment"`
	Tag      string `json:"tag"`
}

type table struct {
	TableName   string   `json:"table"`
	HasUpdateAt bool     `json:"hasUpdateAt"`
	Update      bool     `json:"-"`
	HasSortable bool     `json:"-"`
	Columns     []Column `json:"columns"`
}

func (t *table) preDo() {

	for i, _ := range t.Columns {
		c := &t.Columns[i]
		if c.Tag == "" {
			c.Tag = fmt.Sprintf(`json:"%s"`, c.Name)
		}
		if c.GoType == "" {
			c.GoType = "interface{}"
		}
		if c.Sortable {
			t.HasSortable = true
		}
		if c.Update {
			t.Update = true
			if c.GoType == "" {
				c.GoType = "interface{}"
			}
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

func GenerateSchemaFromSchema(pkg string, i, o string) error {
	bs, err := os.ReadFile(i)
	if err != nil {
		return fmt.Errorf("read file %s failed:%w", i, err)
	}
	var tables []table
	if err := json.Unmarshal(bs, &tables); err != nil {
		return err
	}

	code, err := generate(schemaTpl, pkg, tables)
	if err != nil {
		return err
	}
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return err
	}
	return os.WriteFile(o, formatted, 0664)
}

func GenerateGoProtoFromSchema(pkg string, i, o string) error {
	bs, err := os.ReadFile(i)
	if err != nil {
		return fmt.Errorf("read file %s failed:%w", i, err)
	}
	var tables []table
	if err := json.Unmarshal(bs, &tables); err != nil {
		return err
	}

	code, err := generate(goProtoTpl, pkg, tables)
	if err != nil {
		return err
	}
	formatted, err := format.Source([]byte(code))
	if err != nil {
		fmt.Println(code)
		return err
	}
	return os.WriteFile(o, formatted, 0664)
}

func generate(tpl *template.Template, pkg string, tables tables) (string, error) {

	sort.Sort(tables)
	for i := range tables {
		tables[i].preDo()
	}
	data := map[string]interface{}{
		"stmtPkgName": "",
		"package":     pkg,
		"tables":      tables,
	}
	//fmt.Println(data)

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, data); err != nil {
		return "", err
	}
	return buff.String(), nil
}
