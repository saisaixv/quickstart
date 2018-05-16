package main

import (
	"fmt"
	"os"

	_ "quickstart/routers"
	"github.com/astaxie/beego"
)

func main() {
	fmt.Printf("main pid = %d\n",os.Getpid())
	// beego.SetStaticPath("/static","static")
	go beego.Run()
	
	i:=-1
	for i<0{
		// fmt.Printf("index = %d\n",i)
	}
}

