/*
auth:   wuxun
date:   2019-12-20 20:02
mail:   lbwuxun@qq.com
desc:   handle user login message，send token to user client
*/

package handle

import "connQueue/conns"
import "github.com/gorilla/websocket"
import "log"

func GetToken(conn *conns.ClientConn){
	d := "LiuBaiSiMiDa123456"
	connClien := conn.GetConn()
	err :=connClien.WriteMessage(websocket.BinaryMessage, MsgAssemblerReader(d))
	if err != nil {
		log.Print("send token err:", err)
	}
}
