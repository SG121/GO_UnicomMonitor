package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// 发送消息
func sendMessage(wsHost string, paramMsg string) {
	//wss://vd-file-hnzz2-wcloud.wojiazongguan.cn:50443/h5player/live
	uri := url.URL{
		Scheme: "wss",
		Host:   wsHost,
		Path:   "/h5player/live",
	}
	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Fatalf("无法连接到 WebSocket 服务器：", err)
	}
	defer conn.Close()
	// 发送消息
	message := "_paramStr_=" + paramMsg
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		FmtPrint("发送消息失败：", err)
		return
	}
	// 接收消息
	_, response, err := conn.ReadMessage()
	if err != nil {
		FmtPrint("接收消息失败：", err)
		return
	}
	FmtPrint("收到的回复：", Decode(string(response)))
	// 发送消息
	cmdMessage := `{"time":1243,"cmd":3}`
	err = conn.WriteMessage(websocket.TextMessage, []byte(cmdMessage))
	if err != nil {
		FmtPrint("发送命令失败：", err)
		return
	}
	// 继续接收消息并处理数据
	for {
		_, response, err := conn.ReadMessage()
		if err != nil {
			FmtPrint("接收消息失败：", err)
			break
		}
		// 检查特定条件
		if len(response) > 1 && response[1] == 0x63 {
			// 将数据附加到全局变量
			globalVar = append(globalVar, response[0x4e:]...)
			// 打印数据的长度
			// FmtPrint("数据长度：", len(globalVar))
			// 保存数据
			SaveData()
		}
	}
}

// 主函数
func main() {
	FmtPrint("开始启动")
	// 读取配置文件
	config, err := ReadConfig()
	if err != nil {
		FmtPrint("读取配置文件出错：", err)
		return
	}
	// 发送消息
	wsHost := config.WsHost
	paramMsg := config.ParamMsg
	sendMessage(wsHost, paramMsg)
	//
	FmtPrint("已退出")
}
