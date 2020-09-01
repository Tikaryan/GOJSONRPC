package main

import (
	"errors"
	"fmt"
	gojsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

var client struct {
	GetMessage func(key string) ([]string,error)
	AddKeyMessage func(key, msg string) (string,error)
}
var serverAdd, userName string
func main() {
	app := cli.App{
		Name: "Client for GO_JSONRPC_CLI_BADGER",
		Authors: []*cli.Author{
			{Name: "Rajesh Bhai", Email: "bhaiRajeshKhanprHo@gmail.com"},
		},
		Commands: []*cli.Command{
			{Name: "connect", Aliases: []string{"conn"}, Action: clientConn, Usage: "connect to the specified server address"},
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Destination: &serverAdd,Name: "serverAddress", Aliases: []string{"servadd"},Usage: "Providing which server address client want to connect"},
		&cli.StringFlag{Destination: &userName,Name: "username", Aliases: []string{"user"},Usage: "For login purpose username is required and for sending messages if method is called"},
	}
	err := app.Run(os.Args)
	if err != nil{
		log.Fatalln(err)
	}

}
func clientConn(ctx *cli.Context)error{
	var handler = &client
	userName = strings.ToLower(userName)

	if userName != "aryan" && userName != "ram" {
		return errors.New("incorrect username")
	}
	closer, err := gojsonrpc.NewClient("ws://"+serverAdd, "GORPC", handler, nil)
	if err != nil {
		return errors.New("error creating client maybe servAdd,nameSpace value is wrong")
	}
	defer closer()


		db, errs := handler.AddKeyMessage(userName,"It's Been fun ")
		if errs != nil {
			return errors.New("unable to add sorry for inconvenience")
		}
		fmt.Println(db)
		str,err := handler.GetMessage(userName)
		if err != nil {
			return errors.New("Invalid key or Key not found")
		}
		fmt.Println(str)

	return nil
}
