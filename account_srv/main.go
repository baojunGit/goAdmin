package main

import (
	"flag"
	"fmt"
	"github.com/baojunGit/goAdmin/account_srv/biz"
	"github.com/baojunGit/goAdmin/account_srv/initialize"
	"github.com/baojunGit/goAdmin/account_srv/proto/pb"
	"google.golang.org/grpc"
	"net"
)

func init() {
	// 将相对路径转为绝对路径
	initialize.InitDB("./config")
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入ip")
	port := flag.Int("port", 9095, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
