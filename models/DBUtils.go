package models

import(
	"database/sql"
	"log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var mDB *sql.DB

func init()  {
	
	db,err:=sql.Open("mysql","root:saisai@@@/saisai")

	if err!=nil{
		log.Fatal(err)
	}

	mDB=db
	initTables()

}

func initTables()  {

	createUserTableSql:="create table if not exists User( `id` INT UNSIGNED AUTO_INCREMENT,`name` VARCHAR(100) NOT NULL,PRIMARY KEY ( `id` ) )"
	_,err:=mDB.Exec(createUserTableSql)
	if err!=nil{
		log.Fatal(err)
	}

	createTokenTableSql:="create table if not exists Token( `id` INT UNSIGNED AUTO_INCREMENT,`token` VARCHAR(100) NOT NULL,`expiration` int,`refreshtime` DATE,PRIMARY KEY ( `id` ) )"
	_,err=mDB.Exec(createTokenTableSql)
	if err!=nil{
		log.Fatal(err)
	}

	createPwdTableSql:="create table if not exists Password( `id` INT UNSIGNED AUTO_INCREMENT,`user_id` INT NOT NULL,`password` VARCHAR(100) NOT NULL,PRIMARY KEY ( `id` ) )"
	_,err=mDB.Exec(createPwdTableSql)
	if err!=nil{
		log.Fatal(err)
	}

}

func Excute(sql string,args ...interface{}) (sql.Result,error)  {
	fmt.Printf("Excute = %v\n",args)
	results,err:=mDB.Exec(sql,args...)
	return results,err
}

func Query(sql string,args ...interface{})  ([][]string,error){
	fmt.Printf("Query = %v\n",args)
	rows,err:=mDB.Query(sql,args...)
	if err!=nil{
		log.Fatal(err)
	}

	cols,_:=rows.Columns()
	values:=make([][]byte,len(cols))
	scans:=make([]interface{},len(cols))

	for i:=range values{
		scans[i]=&values[i]
	}

	results:=make ([][]string,0)
	i:=0
	for rows.Next(){
		err=rows.Scan(scans...)
		row:=make([]string,0)
		for _,v:=range values{
			row=append(row,string(v))
		}
		results=append(results,row)
		i++
	}
	rows.Close()

	return results,err

}