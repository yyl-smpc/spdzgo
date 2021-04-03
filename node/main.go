package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/yyl-smpc/spdzgo/node/http"
	"github.com/yyl-smpc/spdzgo/node/ws"
	"net"
	"os"
	"strings"
)

type Options struct {
	Host string //主机
	Port string //端口
}

var options = Options{
	Host: "0.0.0.0", //默认
	Port: "1008",
}

func Init()  {
	if strings.Compare(os.Getenv("HOST"),"")  != 0 {
		options.Host = os.Getenv("HOST")
	}
	if strings.Compare(os.Getenv("PORT"),"")  != 0  {
		options.Port = os.Getenv("PORT")
	}
	flag.StringVar(&(options.Host),"host",options.Host,"本地主机")
	flag.StringVar(&(options.Port),"port",options.Port,"本地主机端口")
	flag.Parse()
}

type wr struct {

}

func (w wr)Write(p []byte) (n int, err error)  {
	return len(p), nil
}


func main(){
	Init()
	router := gin.Default()
	http.Handler{}.Init(router)
	ws.Handler{}.Init(router)
	router.Run(net.JoinHostPort(options.Host, options.Port))
}


