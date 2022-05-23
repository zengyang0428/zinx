package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

//封包，拆包的具体模块
type DataPack struct{}
//拆包封包实例的一个初始化方法
func NewDataPack() *DataPack{
	return &DataPack{}
}

//获取包的头的长度方法
func (dp *DataPack)GetHeadLen() uint32{
	//Datalen uint32(4字节) + ID uint32(4字节)
	return 8
}
//封包方法
//|data|len|msgid|data|
func (dp *DataPack)Pack(msg ziface.IMessage)([]byte,error){
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将datalen写进databuff中
	if err  := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgLen());err != nil{
		return nil,err
	}
	
	//将msgid 写进databuff中
	if err  := binary.Write(dataBuff,binary.LittleEndian,msg.GetMsgId());err != nil{
		return nil,err
	}
	//将data数据 写进databuff中
	if err  := binary.Write(dataBuff,binary.LittleEndian,msg.GetData());err != nil{
		return nil,err
	}
	return dataBuff.Bytes(),nil
}

//拆包方法(将包的Head信息都出来)，之后在根据head信息里的data的长度，在进行一次读
//Unpack 拆包方法(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if utils.Globalobject.MaxPackageSize > 0 && msg.DataLen > utils.Globalobject.MaxPackageSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}




