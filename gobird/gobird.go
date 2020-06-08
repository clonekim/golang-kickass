package main

import (
	"gobird/cmd"
	"gobird/server"
	"log"
	"os"
)

func main() {

	var bird cmd.GoBird

	//cli 파라미터 검증
	err := cmd.RunCmd(&bird).Run(os.Args); if err != nil {
		log.Fatal(err)
	}

	//서버 시작
	server.StartHTTPD(&bird)

}
