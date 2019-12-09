/*
auth:   wuxun
date:   2019-12-09 19:54
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package conns

import (
	"github.com/gorilla/websocket"
)

type ClientConn struct{
	userId int
	connID int
	conn *websocket.Conn
}
func NewClient(uId int, con *websocket.Conn, cId int)  *ClientConn{
	return &ClientConn{
		userId:uId,
		conn: con,
		connID:cId,
	}
}

func (c ClientConn)GetUserID()int{
	return c.userId
}

func (c ClientConn)GetConnID()int{
	return c.connID
}

func (c ClientConn)GetConn()*websocket.Conn{
	return c.conn
}