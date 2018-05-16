package models

import(
	"fmt"
	"reflect"
)

type User struct{
	Id int64
	Name string
	Age int64
	Local string
	Token string
}

func Delete(id string) bool{

	_,err:=mDB.Exec("delete from user where id=?",id)
	if err!=nil{
		fmt.Printf("delete mDB.Exec error =%s\n",err)
		return false
	}
	return true
}

func Update(u *User)  bool{
	_,err:=mDB.Exec("update User set name=?,age=?,local=? where id=?",u.Name,u.Age,u.Local,u.Id)
	if err!=nil{
		fmt.Printf("Update mDB.Exec error =%s\n",err)
		return false
	}
	return true
}

func Insert(u *User)  bool{
	_,err:=mDB.Exec("insert into User (name,age,local) values (?,?,?)",u.Name,u.Age,u.Local)
	if err!=nil{
		fmt.Printf("Insert mDB.Exec error =%s\n",err)
		return false
	}
	return true
}

func QueryAll()  []*User{
	rows,err:=mDB.Query("select id Id,name Name,age Age,local Local from user")
	defer rows.Close()
	if err!=nil{
		fmt.Printf("QueryAll mDB.Query error = %s\n",err)
		return nil
	}

	columns,err:=rows.Columns()
	if err!=nil{
		fmt.Printf("QueryAll rows.Columns error = %s\n",err)
		return nil
	}

	values:=make([]interface{},len(columns))

	var uArr []*User

	for rows.Next() {
		u:=User{}
		reflectStruct:=reflect.ValueOf(&u).Elem()

		for i,v:=range columns{
				values[i] = reflectStruct.FieldByName(v).Addr().Interface()
		}

		rows.Scan(values...)
		uArr=append(uArr,&u)
	}

	return uArr

}

func QueryUserById(id string)  *User{
	rows,err:=mDB.Query("select id Id,name Name,age Age,local Local from User where id = ?",id)
	defer rows.Close()
	if err!=nil{
		fmt.Printf("QueryUserById mDB.Query error = %s\n",err)
		return nil
	}

	columns,err:=rows.Columns()
	if err!=nil{
		fmt.Printf("QueryUserById rows.Columns error = %s\n",err)
		return nil
	}

	values:=make([]interface{},len(columns))

	u:=User{}

	reflectStruct:=reflect.ValueOf(&u).Elem()

	for rows.Next(){

		for i,v:=range columns {
			values[i]=reflectStruct.FieldByName(v).Addr().Interface()
		}
		rows.Scan(values)
	}

	return &u
	

}