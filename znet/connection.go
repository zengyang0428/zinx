package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct{
	//当前链接的socket TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool
	////当前链接所绑定的处理业务方法API
	//handleAPI ziface.HandleFunc

	//告知当前链接已经退出/停止 channnel
	ExitChan chan bool
	//该链接处理的方法Router
	Router ziface.IRouter
}
//初始化链接模块的方法
func NewConnection(conn *net.TCPConn,connID uint32,router ziface.IRouter) *Connection{
	c:= &Connection{
		Conn: conn,
		ConnID: connID,
		//handleAPI: callback_api,
		Router: router,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
	return c
}
//链接的读业务方法
func(c *Connection)StartReader(){
	fmt.Println("Reader Coroutine is running..")
	defer fmt.Println("connID=",c.ConnID,"Reader is exit,remote addr is",c.RemoteAddr().String())
	defer c.Stop()
	for{
		//读取客户端的数据到buf中，最大512字节
		// buf := make([]byte,512)
		// _,err := c.Conn.Read(buf)
		// if err != nil{
		// 	fmt.Println("connId=",c.ConnID," recv is err",err)
		// 	continue
		// }
		//创建一个拆包解包对象
		dp := NewDataPack()
		
		//读取客户端的Msg Head 二进制流8个字节
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(c.GetTCPConnection(),headData);err!=nil{
			fmt.Println("read msg head error",err)
		}
		//拆包，得到msgid和msgdatalen 放在msg消息中
		msg,err := dp.Unpack(headData)
		if err != nil{
			fmt.Println("unpack error",err)
			break
		}
		//根据datalen，再次读取Data，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _,err:= io.ReadFull(c.GetTCPConnection(),data) ;err!=nil{
				fmt.Println("read msg data error",err)
				break
			}

		}
		msg.SetData(data)
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg: msg,
		}
		//执行注册的路由方法
		go func (request ziface.IRequest)  {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//从路由中，找到注册绑定的Conn对应的router调用
	}
}

//启动链接让当前链接准备开始工作
func(c *Connection)Start(){
	fmt.Println("Conn Start()...connID=",c.ConnID)
	//启动当前链接的读数据业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务

}
//停止链接 结束当前链接的工作
func(c *Connection) Stop(){
	fmt.Println("Conn Stop()...ConnID=",c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true{
		return
	}
	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()
	//回收资源
	close((c.ExitChan))
}
//获取当前链接的绑定socket conn
func(c *Connection)GetTCPConnection() *net.TCPConn{
	return c.Conn
}

//获取当前链接模块的链接ID
func(c *Connection) GetConnID() uint32{
	return c.ConnID
}
//获取远程客户端的TCP状态 IP port
func(c *Connection) RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}
// //发送数据 将数据发送给远程客户端
// func(c *Connection) Send(data []byte) error{
// 	return nil
// }
//提供一个SenMsg方法，将我们要发送给客户端的数据，先进行封包，在发送
func(c *Connection) SenMsg(msgId uint32,data []byte)error{
	if c.isClosed == true{
		return errors.New("Connection closend when send msg")
	}
	//将data进行封包 MsgDataLen｜MsgID｜data
	dp := NewDataPack()
	binaryMsg,err := dp.Pack(MewMsgPackAge(msgId,data))
	if err != nil{
		fmt.Println("PAck error msg id=",msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	if _,err := c.Conn.Write(binaryMsg); err != nil{
		fmt.Println("Write msg id",msgId,"error",err)
		return errors.New("conn Write error")
	}
	return nil
}