package mysqlsrv

import (
	"database/sql"
	"log"
	// Mysql driver
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/igorbalden/crudtool/services/pagination"
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
	qDt.qry = "SELECT * FROM schemata" + pagination.GetPagntStr(r)

	return makeSessList(w, r, qDt)
}

//ListDB gets table names
func (*Mysqlconn) ListDB(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	dbname := vars["dbname"]
	qDt := new(qryDetails)
	qDt.dbNm = dbname
	qDt.qryCnt = "SELECT COUNT(*) as count FROM information_schema.tables WHERE tables.table_schema = '" + dbname + "'"
	qDt.qry = "SELECT * FROM information_schema.tables WHERE tables.table_schema = '" + dbname + "' " + pagination.GetPagntStr(r)

	return makeSessList(w, r, qDt)
}

//TblCont gets DB table content
func (*Mysqlconn) TblCont(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	dbname := vars["dbname"]
	dbtable := vars["dbtable"]
	// make sql search string
	whereStr := ""
	schVars := r.URL.Query()
	for ink := range schVars {
		// ink is the operator input name, actually "op_[COLUMN_NAME]"
		// skey in the search input name, or column name
		if len(ink) > 2 {
			skey := ink[3:]
			if ink[0:3] == "op_" && schVars[skey][0] != "" {
				if schVars[ink][0] == "" {
					err := errors.New("Operator not selected. Click the back button")
					return err
				}
				whereStr += " AND " + skey + " " + schVars[skey][0]
			}
		}
	}
	if whereStr != "" {
		whereStr = " WHERE " + whereStr[4:]
	}

	qDt := new(qryDetails)
	qDt.dbNm = dbname
	qDt.qryCnt = "SELECT COUNT(*) as count FROM `" + dbname + "`.`" + dbtable + "` " +
		whereStr
	qDt.qry = "SELECT * FROM `" + dbname + "`.`" + dbtable + "` " + whereStr +
		pagination.GetPagntStr(r)
	return makeSessList(w, r, qDt)
}

//GetColmnsMeta gets data type, and collation, about table's columns
func (*Mysqlconn) GetColmnsMeta(w http.ResponseWriter, r *http.Request) *DbData {
	vars := mux.Vars(r)
	dbname := vars["dbname"]
	dbtable := vars["dbtable"]
	metaqDt := new(qryDetails)
	metaqDt.dbNm = dbname
	metaqDt.qryCnt = ""
	metaqDt.qry = "SELECT `COLUMN_NAME`, `COLUMN_TYPE`, `COLLATION_NAME` FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` LIKE '" + dbname + "' AND `TABLE_NAME` LIKE '" + dbtable + "'"

	return makeList(w, r, metaqDt)
}

/*
makeList gets query info by a qryDetails type input and
returns the result data as a pointer to a DbData type var
Errors are fatal
*/
func makeList(w http.ResponseWriter, r *http.Request, qDt *qryDetails) *DbData {
	pConn := GetMyConn(w, r)
	pConn.MakeConnection(w, r, qDt.dbNm)
	retData := new(DbData)
	//Get DATA COUNT for pagination
	if qDt.qryCnt != "" {
		rowsC := pConn.DB.QueryRow(qDt.qryCnt)
		err := rowsC.Scan(&retData.Totalrows)
		if err != nil {
			log.Fatal(err)
		}
		pagination.SetPagnt(r, *retData.Totalrows)
	}

	//Get DATA
	rows, err := pConn.DB.Query(qDt.qry)
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	retData.ColNames = make([]string, 0)
	for _, nms := range columns {
		retData.ColNames = append(retData.ColNames, nms)
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Get the values
	retData.ShData = make([][]string, 0)
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		//err = errors.New("Error Reading scanArgs")
		if err != nil {
			log.Fatal(err)
		}
		// a []byte can hold null values.
		var colValue []byte
		var dataRow []string
		for _, col := range values {
			colValue = col
			dataRow = append(dataRow, string(colValue))
		}
		retData.ShData = append(retData.ShData, dataRow)
	}
	return retData
}

/*
makeSessList gets query info by a qryDetails type input and
writes the result data in the session field MyDbData
Returns the error
*/
func makeSessList(w http.ResponseWriter, r *http.Request, qDt *qryDetails) error {
	pConn := GetMyConn(w, r)
	pConn.MakeConnection(w, r, qDt.dbNm)
	FchData := NewMyData(w, r)
	//Get DATA COUNT for pagination
	if qDt.qryCnt != "" {
		rowsC := pConn.DB.QueryRow(qDt.qryCnt)
		err := rowsC.Scan(&FchData.Totalrows)
		if err != nil {
			log.Print(err)
			return err
		}
		pagination.SetPagnt(r, *FchData.Totalrows)
	}
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
