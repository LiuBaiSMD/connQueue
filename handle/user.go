/*
auth:   wuxun
date:   2019-12-09 20:39
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle

import (
	"connQueue/conns"
	"connQueue/proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func Login(w http.ResponseWriter, r *http.Request) {
	//
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %s", err)
		return
	}
	connID := conns.GetLastestConnID()   //获取最新的connID，进行连接排队
	connClient := conns.NewClient(connID, conn, connID)
	conns.Push(connID,connClient)
}

func MsgAssemblerReader(data string) []byte {
	msgSeqId += 1
	retPb := &heartbeat.Response{
		ClientId: CLIENTID,
		UserId:   USERID,
		MsgId:    msgSeqId,
		SessionId: 1000,
		Data:     data,
	}
	byteData, err := proto.Marshal(retPb)
	if err != nil {
		log.Fatal("pb marshaling error: ", err)
	}
	return byteData
}

func Chat2User(userId int, message string){

}