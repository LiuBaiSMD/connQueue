package main

import (
	"connQueue/proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/util/log"
	"time"
	"net/http"
	"os"
	"os/signal"
	"fmt"
)

const (
	CLIENTID = 10
	USERID   = 10001
)

var userIDCreator chan int

var (
	clientRes heartbeat.Request
	wsHost          = "127.0.0.1:8080"
	wsPath          = "/heartbeat"
	msgSeqId uint64 = 0
)

type Client struct {
	Host string
	Path string
}

func main() {
	userIDCreator = make(chan int, 1)
	userIDCreator <- 10001
	var count int
	go func(){
		for{
			fmt.Println("count: ", count)
			//if count>100{
			//	break
			//}
			count ++
			time.Sleep(time.Microsecond * 10)
			go msgHandler()
	}
	}()
	time.Sleep(time.Second*10)
	log.Log("----------->over")
}


func NewWebsocketClient(host, path string) *Client {
	return &Client{
		Host: host,
		Path: path,
	}
}

func (this *Client) SendMessage() error {

	// 增加一个信号监控,检测各种退出的情况,方便通知服务器断开连接
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dialer := &websocket.Dialer{
		HandshakeTimeout:time.Second * 10,
	}
	conn, _, err := dialer.Dial("ws://"+this.Host+this.Path, http.Header{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close() //关闭连接

	done := make(chan struct{})
	// 另外其一个goroutine处理接收消息
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Log("read:", err)
				return
			}
			if err := proto.Unmarshal(message, &clientRes); err != nil {
				log.Logf("proto unmarshal: %s", err)
			}
			log.Logf("recv: %v", clientRes)
		}
	}()
	//进行发送输入功能
	reader:= make(chan string, 1)
	reader <- "10001"
	d := ""
	for {
		select {
		case <-done:
			return nil
		case d=<-reader:
			err1 :=conn.WriteMessage(websocket.BinaryMessage, MsgAssemblerReader(d))
			if err1 != nil {
				log.Logf("write close:", err1)
			} else {
				continue
			}
		case <-interrupt:
			// 发送 CloseMessage 类型的消息来通知服务器关闭连接，不然会报错CloseAbnormalClosure 1006错误
			// 等待服务器关闭连接，如果超时自动关闭.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "9999"))
			if err != nil {
				log.Fatalf("write close:", err)
				return nil
			}
			log.Fatalf("write close!")
			return nil
		}
	}
}

func getLatestUserID()uint64{
	lUID := <-userIDCreator
	nextUID := lUID+1
	userIDCreator <- nextUID
	return uint64(lUID)
}

func msgHandler() {
	clientWrapper := NewWebsocketClient(wsHost, wsPath)
	if err := clientWrapper.SendMessage(); err != nil {
		log.Logf("SendMessage: errr%v", err)
	}
}

func MsgAssemblerReader(data string) []byte {
	msgSeqId += 1
	retPb := &heartbeat.Request{
		ClientId: CLIENTID,
		UserId:   getLatestUserID(),
		MsgId:    msgSeqId,
		Data:     data,
	}
	byteData, err := proto.Marshal(retPb)
	if err != nil {
		log.Fatal("pb marshaling error: ", err)
	}
	return byteData
}