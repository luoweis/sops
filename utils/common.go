package utils

import (
	"net/http"
	"encoding/json"
	"log"
)

type BaseController struct{}
type ClientInfo struct{
	HttpAddress		string	`json:"httpAddress"`
	SSL				bool	`json:"ssl"`
}
type ServerController struct {
	Clients 	[]*ClientInfo
}


func NewServerController()(server *ServerController, err error){
	return
}

// 测试方法，可以删除
func (handler *BaseController) InitTest(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm() // 解析参数，默认不进行参数解析
	dataMap := map[string]interface{}{
		"state":200,
		"code":200,
		"info":"init test url",
	}
	dataStr,_ := json.Marshal(dataMap)
	resp.Write([]byte(dataStr))
	log.Println("request path:",req.URL.Path)
}

func (server *ServerController) AddClient(client *ClientInfo) (s *ServerController){
	server.Clients = append(server.Clients, client)
	s = server
	return
}

func (server *ServerController) ShowClient(){

}
