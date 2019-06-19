package main

import (
	"flag"
	"sops/utils"
	"net/http"
	"fmt"
	"log"
	"errors"
)

var (
	welcome           string
	help              bool
	port              string
	role              string
	token             string
	serverHttpAddress string
	clientHttpAddress string
	err               error
	server			 *utils.ServerController
)
func init (){
	flag.BoolVar(&help,"h",false,"this help")
	flag.StringVar(&port,"p","1821","sops http port")
	flag.StringVar(&role,"r","server","sops role server or client")
	flag.StringVar(&token,"auth","","sops client auth token value,use for client")
	flag.StringVar(&serverHttpAddress,"serverHttpAddress","127.0.0.1","http address,use for client")
	flag.Usage = utils.FlagUsage

	server,_ = utils.NewServerController()
}

func main(){
	base := new(utils.BaseController)
	flag.Parse()

	if help {
		flag.Usage()
		return
	}
	switch role {
	case "server":
		welcome += fmt.Sprintf("port：%v,role:%v\n",port,role)
		token, _ := utils.GenerateToken()
		//log.Println(token)

		welcome += fmt.Sprintf("client join to server neet auth token,for example:\n")
		welcome += fmt.Sprintf("	./sops -serverHttpAddress=\"%s:%s\" -r=\"client\" -auth=\"sops:%s\"\n", serverHttpAddress,port,token)
	case "client":
		if serverHttpAddress == "" || token == "" {
			flag.Usage()
			return
		}
		_, err := utils.ParseToken(token)
		if err != nil {
			log.Fatal(err)
			return
		}
		welcome += fmt.Sprintf("port：%v,role:%v\n",port,role)
		welcome += fmt.Sprintf("join to server, address:%s\n", serverHttpAddress)

		// join the client to server
		client := new(utils.ClientInfo)
		client.HttpAddress = fmt.Sprintf("http://%s:%s",clientHttpAddress,port)
		server.AddClient(client)
		fmt.Println(server)
	default:
		welcome += fmt.Sprintf("角色定义参数有错，需要server|client")
		err = errors.New("role error")
	}
	fmt.Println(welcome)
	if err != nil {
		return
	}

	http.HandleFunc("/init/test",base.InitTest)

	err = http.ListenAndServe(fmt.Sprintf(":%v",port),nil)
	if err != nil {
		log.Fatal("ListenAndServe:",err)
	}
}