package controllers

import(
	// "github.com/astaxie/beego"
	"quickstart/models"
	"fmt"
	"io/ioutil"
	// "encoding/json"
	// "os"
	"strconv"
	"strings"
	"runtime"

	_ "github.com/go-sql-driver/mysql"

)

type UserController struct {
	MainController
	Token string
	Language string
}



func (c *UserController) Prepare()  {

	id:=Goid()
	fmt.Printf("prepare   ============== id = %d\n",id)
	c.Token=c.FetchHeader("Token")
	c.Language=c.FetchHeader("Language")
	
}

func (c *UserController) Get(){


	req:=c.Ctx.Request

	body,err:=ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err!=nil{
		fmt.Printf("error = %s\n",err)
	}
	
	fmt.Printf("request body = %s\n",body)

	u,err:=models.Query("select * from User")
	if err!=nil{
		fmt.Printf("error = %s\n",err)
	}

	fmt.Printf("values = %s\n",u)

	c.Ctx.WriteString("names")
}

func (c *UserController) Post(){

	result,_:=models.Query("select * from Token")

	fmt.Printf("token result = %s\n",result)

	login:=new(Login)

	defer func (){
		fmt.Println("excute this function")
		c.Data["json"]=login
		c.ServeJSON()
	}()

	// c.Token=c.FetchHeader("Token")
	// c.Language=c.FetchHeader("Language")

	// fmt.Printf("token = %s,language = %s\n",c.Token,c.Language)

	err:=c.FetchJsonBody(login)
	if err!=nil{
		fmt.Printf("error = %s\n",err)
	}

	fmt.Printf("body = %s\n",login)
}

func (this *UserController) FetchHeader(headerType string) string {
	header:=this.Ctx.Input.Header(headerType)
	if header!=""{
		return header
	}
	return ""
}



func Goid() int {
    defer func()  {
        if err := recover(); err != nil {
            fmt.Println("panic recover:panic info:%v", err)        }
    }()

    var buf [64]byte
    n := runtime.Stack(buf[:], false)
    idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
    id, err := strconv.Atoi(idField)
    if err != nil {
        panic(fmt.Sprintf("cannot get goroutine id: %v", err))
    }
    return id
}