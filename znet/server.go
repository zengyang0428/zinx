package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前的Server添加一个router，server注册的链接对应的处理业务
	Router ziface.IRouter
}
// //定义当前客户端链接的所绑定handle api(目前这个hamdle是写死的，以后优化应该由用户自定义handle方法)
// func CallBackToClient(conn *net.TCPConn,data []byte,cnt int) error{
// 	//回显的业务
// 	fmt.Println("[Conn Handle]CallbackToClient..")
// 	if _,err:=conn.Write(data[:cnt]);err != nil{
// 		fmt.Println("write back buf err",err)
// 		return errors.New("CallBackToCloent error")
// 	}
// 	return nil
// }
//启动服务器
func (s *Server) Start(){
	fmt.Printf("Zinx Server Name:%s,listenner at IP:%s,Port:%d is starting\n",s.Name,s.IP,s.Port)
	fmt.Printf("Zinx version %s,MaxConn:%d,MaxPackeetSize:%d\n",s.IPVersion,
	utils.Globalobject.MaxConn,utils.Globalobject.MaxPackageSize)
	fmt.Printf("[Start] Server Listenner at IP:%s,Port %d,is starting\n",s.IP,s.Port)
	go func ()  {
		//1.获取一个TCP的Addr
		addr,err :=  net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil{
			fmt.Println("resolve tcp addt error",err)
			return
		}
		//2.监听服务器的地址
		Listenner,err :=net.ListenTCP(s.IPVersion,addr)
		if err != nil{
			fmt.Println("listen",s.IPVersion,"err",err)
			return
		}
		fmt.Println("start Zinx server succ",s.Name,"succ,Listenning")
		var cid uint32
		cid = 0
		//3.阻塞的等待客户端链接，处理客户端链接业务(读写)
		for{
			//如果用客户端链接过来，阻塞会返回
			conn,err := Listenner.AcceptTCP()
			if err != nil{
				fmt.Println("Accept err",err)
				continue
			}
			//将处理新链接的业务方法和conn进行绑定得到我们的链接模块
			dealConn := NewConnection(conn,cid,s.Router)
			cid++
			//启动当前的链接业务处理
			go dealConn.Start()
		}

	}()
	
}
//停止服务器
func (s *Server) Stop(){
	//TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行亭子或者回收
}	
//运行服务器
func (s *Server) Serve(){
	//启动server的服务器功能
	s.Start()
	//TODO做一些启动服务器之后的额外业务
	//阻塞状态
	select{}
}
//路由功能：给当前的服务器注册一个路由方法，供客户端的链接处理使用
func(s *Server)AddRouter(router ziface.IRouter){
	s.Router = router
	fmt.Println("Add Router Succ!")
}
/*
初始化Server 的方法
*/
func NewServer(name string) ziface.IServer{
	S := &Server{
		Name : utils.Globalobject.Name,
		IPVersion: "tcp4",
		IP: utils.Globalobject.Host,
		Port: utils.Globalobject.TcpPort,
		Router: nil,
	}
	return S
}