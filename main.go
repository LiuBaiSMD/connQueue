package main

import (
	"connQueue/conns"
	"connQueue/handle"
	"connQueue/proto"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/web"
	"fmt"
	"log"
)

var upGrader = websocket.Upgrader{
	//对请求头进行检查
	//CheckOrigin: func(r *http.Request) bool { return true },
}
var (
	clientRes heartbeat.Request
	serverRsp heartbeat.Response
	msgSeqId uint64 = 0
	USERID uint64 = 666
	CLIENTID uint64 = 678

)

func main() {
	// New web service

	service := web.NewService(
		web.Name("go.micro.web.heartbeat"),
		web.Address(":8080"),
	)
	go func(){
		for{
			connID, conn := conns.Pop()
			fmt.Println("connClient: ", connID, conn)
		}
	}()
	if err := service.Init(); err != nil {
		log.Fatal("Init", err)
	}
	// websocket 连接接口 web.name注册根据.分割路由路径，所以注册的路径要和name对应上
	service.HandleFunc("/heartbeat", handle.Login)
	if err := service.Run(); err != nil {
		log.Fatal("Run: ", err)
	}
}
