package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数是可以通过Zinx，json由用户进行配置
*/
type Globalobj struct{
	/*
	Server
	*/
	TcpServer ziface.IServer//当前Zinx全局的Server对象
	Host string 		//当前服务器主机监听的IP
	TcpPort int 		//当前服务器主机监听的端口号
	Name string 		//当前服务器的名称
	/*
	Zinx
	*/
	Version string 		//当前Zinx的版本号
	MaxConn int		//当前服务器主机允许的最大链接数
	MaxPackageSize uint32 	//当前Zinx框架数据包的最大值

}
/*
定义一个全局的对外Globalobj
*/
var Globalobject *Globalobj
/*
从zinx.json去加载用于自定义的参数
*/
func(g *Globalobj) Reload(){
	data,err := ioutil.ReadFile("conf/conf.json")
	if err != nil{
		panic(err)
	}
	//将json文件数据解析道srtuct中
	err = json.Unmarshal(data,&Globalobject)
	if err != nil{
		panic(err)
	}
}
/*
提供一个init方法，初始化当前的Globalobject
*/
func init(){
	//如果配置文件没有加载，默认的值
	Globalobject = &Globalobj{
		Name:"ZinxServerApp",
		Version: "V0.5",
		TcpPort: 8999,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从conf/Zinx，json去加载一些用户自定义的参数
	Globalobject.Reload()
}

