package main

import (
	"fmt"
	"strconv"
	"runtime"

	"quickstart/cache"
	"quickstart/utils"

	_ "quickstart/routers"
	"github.com/astaxie/beego"
	"github.com/go-ini/ini"
	clog "github.com/cihub/seelog"
)

const(
	LINUX_PROFILE = "/home/saisai/gosource/src/quickstart/conf/profile.ini"
	WIN_PROFILE   = "C:/Caton/goSrc/src/caton/cydex_saas/cydex_manager/deploy/profile_test.ini"
	MAC_PROFILE   = "/Users/yingyue/go/src/caton/cydex_saas/cydex_manager/deploy/profile_mac.ini"
)

var(
	http_addr string
	beego_loglevel int
	//https config
	beego_enable_https bool
	beego_graceful bool
	beego_https_addr string
	beego_https_port string
	beego_https_certfile string
	beego_https_keyfile string
)

func main() {

	var config *utils.Config

	if runtime.GOOS=="windows"{
		fmt.Println("is windows........")
		config=utils.NewConfig(WIN_PROFILE)
	}else if  runtime.GOOS=="darwin"{
		fmt.Println("is macOS........")
		config=utils.NewConfig(MAC_PROFILE)
	}else{
		fmt.Println("is linux........")
		config=utils.NewConfig(LINUX_PROFILE)
	}

	if config==nil{
		clog.Critical("new config failed")
		return
	}

	utils.MakeDefaultConfig(config)
	cfg,err:=config.Load()
	if err!=nil{
		clog.Critical(err)
		return
	}

	if err=setupHttp(cfg);err!=nil{
		fmt.Println("finish")
		return
	}

	if err=setupRedis(cfg);err!=nil{
		fmt.Println("finish")
		return
	}

	fmt.Println("run............................")

	
	// run()

	// beego.SetStaticPath("/static","static")
	beego.Run()
	
	// utils.WaitEndSignal()

}

func run()  {
	
	beego.SetLevel(beego_loglevel)

	if beego_enable_https==true{
		beego.BConfig.Listen.EnableHTTP=false
		go beego.Run(beego_https_addr+":"+beego_https_port)
	}else{
		go beego.Run(http_addr)
	}
	clog.Info("Run...")
}

func setupHttp(cfg *ini.File) (err error) {
	sec,err:=cfg.GetSection("http")
	if err!=nil{
		return err
	}
	addr:=sec.Key("addr").String()
	if addr ==""{
		err=fmt.Errorf("HTTP addr can't be empty")	
		return err
	}

	// show_req:=sec.Key("show_req").MustBool(false)
	// show_rsp:=sec.Key("show_rsp").MustBool(false)

	// clog.Infof("[setup http] addr:'%s',show:[%t,%t]",addr,show_req,show_rsp)
	http_addr=addr

	beego_loglevel=sec.Key("beego_loglevel").MustInt(beego.LevelInformational)

	fmt.Println("setup http 1")

	//https config
	beego_enable_https=sec.Key("beego_enable_https").MustBool(false)
	beego_graceful=sec.Key("beego_graceful").MustBool(false)
	beego_https_addr=sec.Key("beego_https_addr").String()
	beego_https_port:=sec.Key("beego_https_port").String()
	beego_https_certfile=sec.Key("beego_https_certfile").String()
	beego_https_keyfile=sec.Key("beego_https_keyfile").String()

	beego.BConfig.Listen.EnableHTTPS=beego_enable_https
	beego.BConfig.Listen.Graceful=beego_graceful
	beego.BConfig.Listen.HTTPAddr=beego_https_addr
	beego.BConfig.Listen.HTTPSPort,_=strconv.Atoi(beego_https_port)
	beego.BConfig.Listen.HTTPSCertFile=beego_https_certfile
	beego.BConfig.Listen.HTTPSKeyFile=beego_https_keyfile

	fmt.Println("setup http finish")

	return
}

func setupRedis(cfg *ini.File) (err error)  {
	
	sec,err:=cfg.GetSection("redis")
	if err!=nil{
		return err
	}
	url:=sec.Key("url").String()
	max_idles:=sec.Key("max_idles").MustInt(3)
	idle_timeout:=sec.Key("idle_timeout").MustInt(240)

	fmt.Printf("url = %s,max_idles = %d,idle_timeout = %d\n",url,max_idles,idle_timeout)

	//9表示 创建 redis要访问的是9号数据库
	cache.Init(url,max_idles,idle_timeout,9)


	err=cache.Ping()
	if err!=nil{
		fmt.Println("ping error")
		return err
	}
	return nil
}
