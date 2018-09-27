package mysqlsrv

import (
	"database/sql"
	"log"
	// Mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

type qryDetails struct {
	dbNm   string
	qryCnt string
	qry    string
}

//GetDBs get all databases
func (*Mysqlconn) GetDBs(w http.ResponseWriter, r *http.Request) error {
	qDt := new(qryDetails)
	qDt.dbNm = "information_schema"
	qDt.qryCnt = "SELECT COUNT(*) as count FROM schemata"
	qDt.qry = "SELECT * FROM schemata" + GetPagntStr(r)

	return makeList(w, r, qDt)
}

//ListDB gets table names
func (*Mysqlconn) ListDB(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	dbname := vars["dbname"]
	qDt := new(qryDetails)
	qDt.dbNm = dbname
	qDt.qryCnt = "SELECT COUNT(*) as count FROM information_schema.tables WHERE tables.table_schema = '" + dbname + "'"
	qDt.qry = "SELECT * FROM information_schema.tables WHERE tables.table_schema = '" + dbname + "' " + GetPagntStr(r)

	return makeList(w, r, qDt)
}

//TblCont gets DB table content
func (*Mysqlconn) TblCont(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	dbname := vars["dbname"]
	dbtable := vars["dbtable"]
	qDt := new(qryDetails)
	qDt.dbNm = dbname
	qDt.qryCnt = "SELECT COUNT(*) as count FROM `" + dbname + "`.`" + dbtable + "`"
	qDt.qry = "SELECT * FROM `" + dbname + "`.`" + dbtable + "` " + GetPagntStr(r)
	return makeList(w, r, qDt)
}

func makeList(w http.ResponseWriter, r *http.Request, qDt *qryDetails) error {
	pConn := GetMyConn(w, r)
	pConn.MakeConnection(w, r, qDt.dbNm)
	FchData := NewMyData(w, r)
	//Get DATA COUNT for pagination
	rowsC := pConn.DB.QueryRow(qDt.qryCnt)
	err := rowsC.Scan(&FchData.Pagnt.RwsTot)
	if err != nil {
		log.Print(err)
		return err
	}
	SetPagnt(r, FchData)

	//Get DATA
	rows, err := pConn.DB.Query(qDt.qry)
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Print(err)
		return err
	}
	FchData.ColNames = make([]string, 0)
	for _, nms := range columns {
		FchData.ColNames = append(FchData.ColNames, nms)
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Get the values
	FchData.ShData = make([][]string, 0)
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		//err = errors.New("Error Reading scanArgs")
		if err != nil {
			log.Print(err)
			return err
		}
		// a []byte can hold null values.
		var colValue []byte
		var dataRow []string
		for _, col := range values {
			colValue = col
			dataRow = append(dataRow, string(colValue))
		}
		FchData.ShData = append(FchData.ShData, dataRow)
	}
	return err
}
