package util

import(
	"log"
	"github.com/garyburd/redigo/redis"
)

var (
	Conn *redis.Conn
)

func init()  {
	conn,err:=redis.Dial("tcp","127.0.0.1:6379")
	if err!=nil{
		log.Fatal(err)
	}

	Conn=&conn
}

func DoGet(key string,obj interface{})  bool{
	
	DoGet()
}