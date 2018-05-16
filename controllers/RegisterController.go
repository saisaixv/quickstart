package controllers

import(
	"fmt"
	"quickstart/models"
)

type RegisterController struct {
	MainController
}

type Register struct {

	UserName string `json:"username"`
	Password string `json:"password"`
}

type RegisterRsp struct {

	Error int `json:"errnum"`
	UserID int64 `json:"user_id"`
	UserName string `json:"username"`
	
}

func (this *RegisterController)Post()  {
	
	rsp:=new(RegisterRsp)

	defer func(){
		this.Data["json"]=rsp
		this.ServeJSON()
	}()

	register:=new(Register)

	err:=this.FetchJsonBody(register)
	if err!=nil{
		rsp.Error=101
		return
	}

	

	result,err:=models.Query("select * from User where name=?",register.UserName)
	if err!=nil{
		rsp.Error=102
		return
	}

	if len(result)!=0{
		rsp.Error=103
		return 
	}
	fmt.Printf("register = %s\n",register.UserName)
	results,err:=models.Excute("insert into User (name) values(?)",register.UserName)
	if err!=nil{
		rsp.Error=104
		fmt.Printf("models.Excute error = %s\n",err)
		return
	}
	id,err:=results.LastInsertId()
	if err!=nil{
		rsp.Error=104
		fmt.Printf("results.LastInsertId error = %s\n",err)
		return
	}
	fmt.Printf("insert id = %d\n",id)
	rsp.UserID=id

	results,err=models.Excute("insert into Password (user_id,password) values(?,?)",id,register.Password)
	if err!=nil{
		rsp.Error=104
		fmt.Printf("models.Excute error = %s\n",err)
		return
	}
	rsp.Error=0
	rsp.UserName=register.UserName
	
}