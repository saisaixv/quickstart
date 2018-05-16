package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type MainController struct {
	beego.Controller

}

type JSONStruct struct {
    Code int
    Msg  string
}

func (c *MainController) Get() {

	// mystruct := &JSONStruct{0, "hello"}

	var dat []map[string]interface{}
	json.Unmarshal([]byte(`[{"name":"saisai","age":25},{"name":"aaa","age":28}]`),&dat)

	for _,value:= range dat{
		fmt.Printf("name = %s,age = %f",value["name"],value["age"])
	}
	

    c.Data["json"] = &dat
    c.ServeJSON()

	// var dat []map[string]interface{}
	// err:=json.Unmarshal([]byte(`{"name":"saisai","age":25}`),&dat)
	// if err!=nil{
	// 	fmt.Println("Json  format error")
	// 	return
	// }
	// c.Data["json"]=&dat
	// c.ServeJSON()
	// c.Ctx.Output.JSON(dat,false,false)
	// c.Ctx.WriteString("bbb")
}

// ServeJSON sends a json response with encoding charset.
func (c *MainController) ServeJSON(encoding ...bool) {
    var (
        hasIndent   = true
        hasEncoding = false
    )
    
        hasIndent = false
    
    if len(encoding) > 0 && encoding[0] {
        hasEncoding = true
    }
    c.Ctx.Output.JSON(c.Data["json"], hasIndent, hasEncoding)
}

func (c *MainController) Post(){
	c.Ctx.Output.Body([]byte(`{"name":"saisai","age":27}`))
}

func (this *MainController) FetchJsonBody(v interface{}) error {
	req:=this.Ctx.Request
	defer req.Body.Close()

	body,err:=ioutil.ReadAll(req.Body)
	if err!=nil{
		return err
	}

	if len(body)==0{
		return nil
	}

	if err:=json.Unmarshal(body,v);err!=nil{
		return err
	}
	return nil
}