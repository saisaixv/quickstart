package cache

import(
	
	"fmt"
	"time"
	"encoding/json"

	clog "github.com/cihub/seelog"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

func newPool(url string,max_idle,idle_timeout int,db int) *redis.Pool{

	timeout:=time.Duration(idle_timeout)*time.Second
	return &redis.Pool{
		MaxIdle:max_idle,
		IdleTimeout:timeout,
		Dial:func()(redis.Conn,error){
			c,err:=redis.DialURL(url)
			if err!=nil{
				fmt.Println("DialURL error")
				return nil,err
			}

			//使用 几号 数据库
			_,err=c.Do("SELECT",0)

			if err!=nil{
				fmt.Println("select db error")
				fmt.Println(err)
			}

			return c,err
		},
		TestOnBorrow:func(c redis.Conn,t time.Time)error{
			if time.Since(t)<time.Minute{
				fmt.Println("select db error")
				return nil
			}

			_,err:=c.Do("PING")
			return err
		},
	}
}

func Init(url string,max_idle,idle_timeout int,db int)  {
	pool=newPool(url,max_idle,idle_timeout,db)	
}

func Get() redis.Conn{
	if pool ==nil{
		clog.Critical("Please set cache pool first!")
		return nil
	}
	return pool.Get()
}

func Ping() error{
	conn:=Get()
	defer conn.Close()
	_,err:=conn.Do("PING")
	return err
}

func Close(){
	clog.Info("[redis pool close]")
	pool.Close()
}

func DoMapSet(key string,obj map[int]map[int]string,expire int)  {
	redisConn:=Get()
	defer redisConn.Close()

	value,_:=json.Marshal(obj)
	
	_,err:=redisConn.Do("SETEX",key,expire,value)
	if err!=nil{
		fmt.Println(err)
	}
	
}

func DoMapGet(key string)(obj map[int]map[int]string)  {
	redisConn:=Get()
	defer redisConn.Close()

	value,err:=redis.Bytes(redisConn.Do("GET",key))
	if err!=nil{
		fmt.Println(err)
	}

	errShal:=json.Unmarshal(value,&obj)
	if errShal!=nil{
		fmt.Println(errShal)
	}
	return obj
}

func DoSet(key string,obj interface{},expire int) bool  {
	redisConn:=Get()
	defer redisConn.Close()

	value,_:=json.Marshal(obj)

	_,err:=redisConn.Do("SETEX",key,expire,value)
	if err!=nil{
		return false
	}
	return true
}


func DoSetNx(key string,expire int) bool {
	redisConn:=Get()
	defer redisConn.Close()

	retSetNx,errSetNx:=redisConn.Do("SETNX",key,"1")
	if errSetNx!=nil{
		fmt.Println(errSetNx)
		return false
	}

	if retSetNx==0{
		return false
	}

	_,err2:=redisConn.Do("EXPIRE",key,expire)
	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	return true
}

func DoHSet(key string,field string,obj interface{},expire int)bool  {
	redisConn:=Get()
	defer redisConn.Close()

	value,_:=json.Marshal(obj)

	_,err:=redisConn.Do("HSET",key,field,value)
	if err!=nil{
		fmt.Println(err)
		return false
	}

	_,err2:=redisConn.Do("EXPIRE",key,expire)
	if err2!=nil{
		fmt.Println(err)
		return false
	}

	return true
}

func DoHDel(key string,field string) bool {
	redisConn:=Get()
	defer redisConn.Close()

	_,err:=redisConn.Do("HDEL",key,field)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	return true
}

func DoHGet(key string,field string,obj interface{}) bool {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err1:=redisConn.Do("HGET",key,field)
	if ret==nil{
		return false
	}

	value,err2:=redis.Bytes(ret,err1)
	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	if value==nil{
		return false
	}

	err3:=json.Unmarshal(value,obj)
	if err3!=nil{
		fmt.Println(err3)
		return false
	}
	return true
}

//set key expire date
func DoExpire(key string,expire int) bool {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err:=redisConn.Do("EXPIRE",key,expire)
	if err!=nil{
		return false
	}
	value,_:=redis.Int(ret,err)

	if value==1{
		return true
	}else{
		return false
	}
}

func DoDel(key string) bool {
	redisConn:=Get()
	defer redisConn.Close()

	_,err:=redisConn.Do("DEL",key)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	return true
}

func DoGet(key string,obj interface{}) bool {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err1:=redisConn.Do("GET",key)
	if ret==nil{
		return false
	}

	value,err2:=redis.Bytes(ret,err1)
	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	err3:=json.Unmarshal(value,obj)
	if err3!=nil{
		fmt.Println(err3)
		return false
	}

	return true
}

func DoFlushDb() bool  {
	redisConn:=Get()
	defer redisConn.Close()

	_,_=redisConn.Do("FLUSHDB")
	return true
}

func DoStrSet(key string,obj string,expire int) bool {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err:=redisConn.Do("SETEX",key,expire,obj)
	if ret==nil{
		fmt.Println(err)
		return false
	}else{
		return true
	}
}

func DoStrGet(key string) (ret bool,obj string) {
	redisConn:=Get()
	defer redisConn.Close()

	retGet,err1:=redisConn.Do("get",key)
	if retGet==nil{
		return false,""
	}

	if err1==nil{
		value,err2:=redis.String(retGet,err1)
		if err2==nil{
			fmt.Println(err2)
			return true,value
		}else{
			return false,""
		}
	}else{
		return false,""
	}
}

func DoKeys(key string) (ret bool,keys []string) {
	redisConn:=Get()
	defer redisConn.Close()

	retGet,err1:=redisConn.Do("keys",key)
	if retGet==nil{
		return false,nil
	}

	if err1==nil{
		value,err2:=redis.Strings(retGet,err1)
		if err2!=nil{
			fmt.Println(err2)
			return false,nil
		}else{
			return true,value
		}
	}else{
		return false,nil
	}
}


func DoStrHSetConn(key string,field string,value string,conn redis.Conn)  {
	_,err:=conn.Do("hset",key,field,value)
	if err!=nil{
		fmt.Println(err)
	}
}

func DoStrHGetConn(key string,field string,conn redis.Conn) (obj string) {
	ret,err1:=conn.Do("hget",key,field)
	if ret!=nil{
		value,err2:=redis.String(ret,err1)
		if err2!=nil{
			fmt.Println(err2)
		}
		return value
	}

	return ""
}


func DoDelConn(key string,conn redis.Conn) bool {
	_,err:=conn.Do("DEL",key)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	return true
}

func DoHDelConn(key string,field string,conn redis.Conn) bool {
	_,err:=conn.Do("HDEL",key)
	if err!=nil{
		return false
	}
	return true
}

func DoHkeys(key string) []string {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err:=redis.Strings(redisConn.Do("hkeys",key))
	if err!=nil{
		fmt.Println(err)
	}

	return ret
}

func  DoHVals(key string)(bool,[]interface{})  {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err1:=redisConn.Do("HVals",key)
	if ret==nil{
		return false,nil
	}

	value,err2:=redis.ByteSlices(ret,err1)
	if err2!=nil{
		return false,nil
	}

	obj:=make([]interface{},0)
	for _,v:=range value{
		var f interface{}

		err3:=json.Unmarshal(v,&f)
		if err3!=nil{
			return false,nil
		}
		obj =append(obj,f)
	}

	return true,obj
}

func DoHLen(key string)int64  {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err:=redis.Int64(redisConn.Do("hlen",key))
	
	if err!=nil{
		fmt.Println(err)
	}
	return ret
}

func DoSetConn(key string,obj interface{},expire int,conn redis.Conn) bool {
	value,_:=json.Marshal(obj)

	_,err:=conn.Do("SETEX",key,expire,value)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	return true
}

func DOGetConn(key string,obj interface{},conn redis.Conn) bool {
	ret,err1:=conn.Do("GET",key)
	if ret==nil{
		return false
	}

	value,err2:=redis.Bytes(ret,err1)

	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	err3:=json.Unmarshal(value,obj)
	if err3!=nil{
		fmt.Println(err3)
		return false
	}
	return true
}
func  DoRPush(key string,obj interface{},expire int) bool  {
	redisConn:=Get()
	defer redisConn.Close()

	value,_:=json.Marshal(obj)
	_,err:=redisConn.Do("RPUSH",key,value)

	if err!=nil{
		fmt.Println(err)
		return false
	}

	_,err2:=redisConn.Do("EXPIRE",key,expire)
	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	return true
}


func DoLPop(key string,obj interface{}) bool {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err1:=redisConn.Do("LPop",key)
	if ret==nil{
		return false
	}

	value,err2:=redis.Bytes(ret,err1)
	if err2!=nil{
		fmt.Println(err2)
		return false
	}

	if value==nil{
		return false
	}

	err3:=json.Unmarshal(value,obj)
	if err3!=nil{
		fmt.Println(err3)
		return false
	}

	return true
}



func DoZAdd(key string,score float64,obj interface{}) (int,bool) {
	redisConn:=Get()
	defer redisConn.Close()

	value,_:=json.Marshal(obj)

	ret,err1:=redisConn.Do("ZAdd",key,score,value)
	if err1!=nil{
		fmt.Println(err1)
		return -1,false
	}

	ret2,err2:=redis.Int(ret,err1)
	if err2!=nil{
		fmt.Println(err2)
		return -1,false
	}

	return ret2,true
}

func DoZRange(key string,start int ,stop int) (bool,[]interface{}) {
	redisConn:=Get()
	defer redisConn.Close()

	ret,err1:=redisConn.Do("ZRange",key,start,stop)
	if ret!=nil{
		return false,nil
	}

	value,err2:=redis.ByteSlices(ret,err1)
	if err2!=nil{
		fmt.Println(err2)
		return false,nil
	}

	obj:=make([]interface{},0)
	for _,v:=range value{
		var f interface{}
		err3:=json.Unmarshal(v,&f)
		if err3!=nil{
			return false,nil
		}

		obj=append(obj,f)
	}

	return true,obj
}

func DoExists(key string)bool  {
	redisConn:=Get()
	defer redisConn.Close()

	exists,err:=redis.Bool(redisConn.Do("EXISTS",key))
	if err!=nil{
		return false
	}
	return exists
}