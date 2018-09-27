package mysqlsrv

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"strings"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/igorbalden/crudtool/services/ibsession"
	"net/http"
	"strconv"
)

var err error

//Mysqlconn database connection string
type Mysqlconn struct {
	ConStr string
	DB     *sql.DB
}

//DbData passes data to views
type DbData struct {
	Pagnt    Pgn
	ColNames []string
	ShData   [][]string
}

//Pgn passes pagination data to views
type Pgn struct {
	PgURL   string
	RwsTot  int
	PgCurr  int
	PgLast  int
	PgPrev  int
	PgNext  int
	PgOne   int
	PgTwo   int
	PgThree int
	PgFour  int
	PgFive  int
}

const limitQry = 20

func init() {
}

//GetMyConn get the DB Coneection for this session
func GetMyConn(w http.ResponseWriter, r *http.Request) *Mysqlconn {
	sess, _ := ibsession.GetSession(w, r)
	if value, ok := sess.Get("MyConn").(*Mysqlconn); ok {
		return value
	}
	return nil
}

//NewMyData create new DbData for this session
func NewMyData(w http.ResponseWriter, r *http.Request) *DbData {
	sess, _ := ibsession.GetSession(w, r)
	FchData := new(DbData)
	sess.Set("MyDbData", FchData)
	return FchData
}

//GetMyData get the DbData for this session
func GetMyData(w http.ResponseWriter, r *http.Request) *DbData {
	sess, _ := ibsession.GetSession(w, r)
	if value, ok := sess.Get("MyDbData").(*DbData); ok {
		return value
	}
	return nil
}

//Connect to the server
func (*Mysqlconn) Connect(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	servername := r.Form.Get("servername")
	protocol := r.Form.Get("protocol")
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	tmpStr := string(username + ":" + password + "@" + protocol + "(" + servername + ")/")
	pConn := GetMyConn(w, r)
	pConn.ConStr = tmpStr
	dbname := r.Form.Get("dbname")
	if dbname == "" {
		dbname = "information_schema"
	}
	return pConn.MakeConnection(w, r, dbname)
}

//MakeConnection makes the actual connection to DB
func (*Mysqlconn) MakeConnection(w http.ResponseWriter, r *http.Request, dbname string) error {
	pConn := GetMyConn(w, r)
	pConn.DB, err = sql.Open("mysql", pConn.ConStr+dbname)
	if err != nil {
		log.Print("Authentication failed: ", err)
		return err
	}
	err = pConn.DB.Ping() // test connection
	if err != nil {
		log.Print("Database connection failed, ", err)
		return err
	}

	return err
}

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
func SetPagnt(r *http.Request, FchData *DbData) {
	Limit := limitQry
	page := 1
	vars := mux.Vars(r)
	curPg, _ := strconv.Atoi(vars["p"])
	if curPg != 0 {
		page = curPg
	}
	ptot := FchData.Pagnt.RwsTot / Limit
	url := r.RequestURI
	if strings.Index(url, "/p/") > -1 {
		FchData.Pagnt.PgURL = url[:strings.Index(url, "/p/")]
	} else {
		FchData.Pagnt.PgURL = url
	}

	FchData.Pagnt.PgCurr = page
	FchData.Pagnt.PgPrev = page - 1
	FchData.Pagnt.PgNext = page + 1
	FchData.Pagnt.PgOne = page - 2
	FchData.Pagnt.PgTwo = page - 1
	FchData.Pagnt.PgThree = page
	FchData.Pagnt.PgFour = page + 1
	FchData.Pagnt.PgFive = page + 2
	FchData.Pagnt.PgLast = ptot + 1
	return
}
