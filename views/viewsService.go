package views

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//StrFuncs type to define functions for use in templates
type StrFuncs struct {
	Descr string
	ToLow func(string) string
}

//PStrFuncs var for use in templates
var PStrFuncs *StrFuncs

func viewDataInit() map[string]interface{} {
	var viewData = make(map[string]interface{})
	// define all the functions to use in the templates, here.
	PStrFuncs = &StrFuncs{
		Descr: "Tmpl Functions",
		ToLow: func(origStr string) string {
			return strings.ToLower(origStr)
		},
	}
	viewData["funcs"] = PStrFuncs
	return viewData
}

var templates *template.Template
var loginTemplate *template.Template
var dbSelectTemplate *template.Template
var sqlStmtTemplate *template.Template
var tablesTemplate *template.Template

//PopulateTemplates is used to parse all templates present in
//the templates folder
func PopulateTemplates() {
	var allFiles []string
	templatesDir := "./templates/"
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		log.Println(err)
		os.Exit(1) // No point in running app if templates aren't read
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templatesDir+filename)
		}
	}

	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	loginTemplate = templates.Lookup("login.html")
	dbSelectTemplate = templates.Lookup("dbSelect.html")
	sqlStmtTemplate = templates.Lookup("sqlStmt.html")
	tablesTemplate = templates.Lookup("tables.html")

}
