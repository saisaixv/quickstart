package utils

import(
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

func FileIsExist(filename string) bool {
	_,err:=os.Stat(filename)
	return err==nil || os.IsExist(err)
}

func WaitEndSignal()  {
	sigs:=make(chan os.Signal,1)
	done:=make(chan bool,1)
	
	signal.Notify(sigs,syscall.SIGINT,syscall.SIGTERM)

	go func ()  {
		sig:=<-sigs
		fmt.Println()
		fmt.Println(sig)
		done<-true
	}()
	<-done
}