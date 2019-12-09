/*
auth:   wuxun
date:   2019-12-09 20:39
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package handle
import (
	"connQueue/conns"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"connQueue/proto"
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
	connID := conns.GetLastestConnID()
	connClient := conns.NewClient(9999, conn, connID)
	conns.Push(connID,connClient)
	defer conn.Close()
	reader := make(chan string ,1)
	data := ""
	go func(){
		for{
			log.Printf("please input: 	")
			fmt.Scanf("%s",&data)
			reader <- data
			log.Printf("your input : %v",data)
		}
	}()
	go func(){
		d := ""
		for{
			select {
			case d =<- reader:
				log.Printf("----->send your input")
				err1 :=conn.WriteMessage(websocket.BinaryMessage, MsgAssemblerReader(d))
				if err1 != nil {
					log.Printf("write close:", err1)
				} else {
					log.Printf("send input over!")
				}
			}

		}
	}()
	for { 				//读消息，进行阻塞，
		_, buffer, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if err := proto.Unmarshal(buffer, &clientRes); err != nil {
			log.Printf("proto unmarshal: %s", err)
		}
		log.Printf("recv userId=%d MsgId=%d Data=%s", clientRes.UserId, clientRes.MsgId, clientRes.Data)
	}
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
