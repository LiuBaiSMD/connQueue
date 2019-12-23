/*
auth:   wuxun
date:   2019-12-09 17:25
mail:   lbwuxun@qq.com
desc:   1.store conns (push and pop)
		2.return the current connID which is pop just now
		3.return the current lenth of channel of conns
*/

package conns

import "sync"

type ConnMap struct {
	connChan chan int
	connMap  sync.Map
	curConnID int
	connCMap chan *ClientConn
}

var connIDCreator chan int
var cMap ConnMap

func init() {
	cMap.connCMap = make(chan *ClientConn, 10000)
	cMap.connChan = make(chan int, 10000)
	connIDCreator = make(chan int, 1)
	cMap.curConnID = -1
	connIDCreator <- 1

}

func Push(connID int, connValue *ClientConn){
	cMap.connChan <- connID
	cMap.connCMap <- connValue
}

func PushChan(connID int, connValue interface{}){
	cMap.connMap.Store(connID, connValue)
	cMap.connChan <- connID
}

func Pop()(int, *ClientConn){
	connID := <-cMap.connChan
	cMap.curConnID = connID
	select {
		case connValue := <- cMap.connCMap:
			return connID, connValue
		default:
			return -1, nil
	}
}

func PopChan()(int, *ClientConn){
	connID := <-cMap.connChan
	cMap.curConnID = connID
	connValueITF, isOK := cMap.connMap.Load(connID)
	if !isOK{
		return -1, nil
	}
	Delete(connID)
	connValue := connValueITF.(*ClientConn)
	return connID, connValue
}

func Delete(connID int){
	cMap.connMap.Delete(connID)
}

func LenthConn()int{
	return len(cMap.connChan)
}

func GetCurConnID()int{
	if len(cMap.connChan)==0{
		cMap.curConnID = -1
	}
	return cMap.curConnID
}

func GetLastestConnID()int{
	last := <- connIDCreator
	next := last+1
	connIDCreator <- next
	return last
}