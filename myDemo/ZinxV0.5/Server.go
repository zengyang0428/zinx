package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
基于Zinx框架来开发的 服务器端应用程序
*/
//ping test 自定义路由
type PingRouter struct{
	znet.BaseRouter
}

//Test Bandle
func(this *PingRouter)Handle(request ziface.IRequest){
	fmt.Println("Call Router Handle")
	//先读取客户端的数据，在回写ping..ping..ping
	fmt.Println("recv form client :msgID=",request.GetMsgId(),
	",data=",string(request.GetData()))
	err := request.GetConnection().SenMsg(1,[]byte("ping...ping..ping"))
	if err != nil{
		fmt.Println(err)
	}
}

func main() {
	//1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx v0.5]")
	//2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	
	//2.启动server
	s.Serve()
}
