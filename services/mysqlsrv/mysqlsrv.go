package mysqlsrv

import (
	"database/sql"
	"log"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/igorbalden/crudtool/services/ibsession"
	"net/http"
)

var err error

//Mysqlconn database connection string
type Mysqlconn struct {
	ConStr string
	DB     *sql.DB
}

//DbData passes data to views
type DbData struct {
	Totalrows *int
	ColNames  []string
	ShData    [][]string
}

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
