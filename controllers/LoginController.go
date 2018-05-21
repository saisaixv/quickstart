package controllers

import(
	// "github.com/astaxie/beego"
	"fmt"
	"quickstart/models"
	"strconv"
	
)

type LoginController struct{
	MainController
}

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginRsp struct {
	Error int `json:"errnum"`
	UserID int64 `json:"user_id"`
	Token string `json:"token"`
}

type User struct {
	Id int
	UserName string	
}

func (this *LoginController) Post()  {

	rsp:=new(LoginRsp)
	rsp.Error=0

	defer func (){
		this.Data["json"]=rsp
		this.ServeJSON()
	}()

	login:=new(Login)

	err:=this.FetchJsonBody(login)
	if err!=nil{
		fmt.Printf("error = %s\n",err)
		rsp.Error=101
		return
	}

	result,err:=models.Query("select * from User where name=?",login.UserName)
	if err!=nil{
		fmt.Printf("error = %s\n",err)
		rsp.Error=102
		return
	}

	if len(result)!=1{
		fmt.Printf("error = %s\n",err)
		rsp.Error=102
		return
	}

	u:=new(User)

	u.Id,_=strconv.Atoi(result[0][0])
	u.UserName=result[0][1]

	result,err=models.Query("select password from Password where user_id = ?",u.Id)
	if err!=nil{
		fmt.Printf("error = %s\n",err)
		rsp.Error=102
		return
	}

	if len(result)!=1{
		fmt.Printf("error = %s\n",err)
		rsp.Error=102
		return
	}

	pwd:=result[0][0]

	if login.Password!=pwd{
		fmt.Printf("error = %s\n",err)
		rsp.Error=103
		return
	}

	token:=login.UserName+strconv.FormatInt(rsp.UserID,10)
	
	rsp.Token=token

	
}
