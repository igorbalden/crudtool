package views

import (
	"github.com/gorilla/mux"
	"github.com/igorbalden/crudtool/services/mysqlsrv"
	"github.com/igorbalden/crudtool/services/pagination"
	"github.com/igorbalden/crudtool/utils"

	"log"
	"net/http"

	//no new variables
	_ "github.com/go-sql-driver/mysql"
)

//DbSelect shows the databases list
func DbSelect(w http.ResponseWriter, r *http.Request) {
	viewData := viewDataInit()
	viewData["Navigation"] = "dbselect"
	viewData["title"] = "Select Database"

	if r.Method == "GET" {
		myConn := mysqlsrv.GetMyConn(w, r)
		err := myConn.GetDBs(w, r)
		if err != nil {
			log.Print(err)
			viewData["Message"] = "Database error" + " - " + err.Error()
			dbSelectTemplate.Execute(w, viewData)
			return
		}
		myDt := mysqlsrv.GetMyData(w, r)
		viewData["Pagnt"] = pagination.SetPagnt(r, *myDt.Totalrows)
		viewData["ColNames"] = myDt.ColNames
		viewData["ShData"] = myDt.ShData
		dbSelectTemplate.Execute(w, viewData)
	}
}

//ListTables is used to handle the "/dbname/{dbname}" URL
func ListTables(w http.ResponseWriter, r *http.Request) {
	viewData := viewDataInit()
	vars := mux.Vars(r)
	viewData["Navigation"] = "listTables"
	viewData["title"] = "List " + vars["dbname"] + " Tables"
	if r.Method == "GET" {
		myConn := mysqlsrv.GetMyConn(w, r)
		err := myConn.ListDB(w, r)
		if err != nil {
			log.Print(err)
			viewData["Message"] = "Error opening database " + vars["dbname"] + " - " + err.Error()
			tablesTemplate.Execute(w, viewData)
			return
		}
		myDt := mysqlsrv.GetMyData(w, r)
		viewData["Pagnt"] = pagination.SetPagnt(r, *myDt.Totalrows)
		viewData["ColNames"] = myDt.ColNames
		viewData["ShData"] = myDt.ShData
		viewData["dbname"] = vars["dbname"]
		tablesTemplate.Execute(w, viewData)
	}
	return
}

//TblContent lists the contents of a DB Table
func TblContent(w http.ResponseWriter, r *http.Request) {
	viewData := viewDataInit()
	vars := mux.Vars(r)
	viewData["Navigation"] = "tblContent"
	viewData["title"] = "Table " + vars["dbtable"]
	if r.Method == "GET" {
		myConn := mysqlsrv.GetMyConn(w, r)
		err := myConn.TblCont(w, r)
		if err != nil {
			log.Print(err)
			viewData["Message"] = "Error opening database " + vars["dbname"] + " - " + err.Error()
			tablesTemplate.Execute(w, viewData)
			return
		}
		myDt := mysqlsrv.GetMyData(w, r)
		myMetaDt := myConn.GetColmnsMeta(w, r)
		viewData["MetaDt"] = myMetaDt
		viewData["SrchVars"] = r.URL.Query()
		viewData["Pagnt"] = pagination.SetPagnt(r, *myDt.Totalrows)
		viewData["ColNames"] = myDt.ColNames
		viewData["ShData"] = myDt.ShData
		viewData["dbname"] = vars["dbname"]
		viewData["dbtable"] = vars["dbtable"]
		tablesTemplate.Execute(w, viewData)
	}
	return
}

//SQLStmt executes an Sql statement
func SQLStmt(w http.ResponseWriter, r *http.Request) {
	viewData := viewDataInit()
	viewData["Navigation"] = "sqlstmt"
	viewData["title"] = "SQL Statement"

	myConn := mysqlsrv.GetMyConn(w, r)
	viewData["dbname"] = myConn.ConStr

	if r.Method == "GET" {
		//time.Sleep(3 * time.Second)
		viewData["csrfToken"] = utils.MakeCsrf(w, r)
		sqlStmtTemplate.Execute(w, viewData)
	}

	if r.Method == "POST" {
		r.ParseForm()
		if utils.ChkCsrf(w, r) {
			stmtIn := r.Form.Get("sqlstmt")
			viewData["message"] = "SQL Statement is :" + stmtIn
		} else {
			viewData["message"] = "Form data error."
		}
		viewData["csrfToken"] = utils.MakeCsrf(w, r)
		sqlStmtTemplate.Execute(w, viewData)
	}
}
