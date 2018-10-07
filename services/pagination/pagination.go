package pagination

import (
	"fmt"
	"github.com/gorilla/mux"
	//"log"
	"strings"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

//Pgn passes pagination data to views
type Pgn struct {
	PgURLPath string
	PgURLQry  string
	RwsTot    int
	PgCurr    int
	PgLast    int
	PgPrev    int
	PgNext    int
	PgOne     int
	PgTwo     int
	PgThree   int
	PgFour    int
	PgFive    int
}

const limitQry = 20

//GetPagnt sets StartRow, and Limit, variables for DB query, from request parameters
//Not in use anymore. GetPagntStr is more useful
func GetPagnt(r *http.Request) (StartRow int, Limit int) {
	Limit = limitQry
	page := 1
	vars := mux.Vars(r)
	curPg, _ := strconv.Atoi(vars["p"])
	if curPg != 0 {
		page = curPg
	}
	StartRow = Limit * (page - 1)
	return StartRow, Limit
}

//GetPagntStr sets the Limit portion for DB query.
func GetPagntStr(r *http.Request) string {
	limit := limitQry
	page := 1
	vars := mux.Vars(r)
	curPg, _ := strconv.Atoi(vars["p"])
	if curPg != 0 {
		page = curPg
	}
	startRow := limit * (page - 1)
	return " LIMIT " + fmt.Sprintf("%d", startRow) + ", " + fmt.Sprintf("%d", limit)
}

//SetPagnt sets current page, and total pages, numbers for the pagination template
func SetPagnt(r *http.Request, totalRows int) *Pgn {
	Pagnt := new(Pgn)
	Pagnt.RwsTot = totalRows
	Limit := limitQry
	page := 1
	vars := mux.Vars(r)
	curPg, _ := strconv.Atoi(vars["p"])
	if curPg != 0 {
		page = curPg
	}
	ptot := Pagnt.RwsTot / Limit
	if Pagnt.RwsTot%Limit > 0 {
		ptot++
	}
	//log.Print("r.URL Query:", r.URL.Query())
	if strings.Index(r.URL.Path, "/p/") > -1 {
		Pagnt.PgURLPath = r.URL.Path[:strings.Index(r.URL.Path, "/p/")]
	} else {
		Pagnt.PgURLPath = r.URL.Path
	}
	Pagnt.PgURLQry = r.URL.RawQuery
	if Pagnt.PgURLQry != "" {
		Pagnt.PgURLQry = "?" + Pagnt.PgURLQry
	}
	Pagnt.PgCurr = page
	Pagnt.PgPrev = page - 1
	Pagnt.PgNext = page + 1
	Pagnt.PgOne = page - 2
	Pagnt.PgTwo = page - 1
	Pagnt.PgThree = page
	Pagnt.PgFour = page + 1
	Pagnt.PgFive = page + 2
	Pagnt.PgLast = ptot

	return Pagnt
}
