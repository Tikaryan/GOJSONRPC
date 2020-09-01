package main

import (
	"fmt"
	gojsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	conn "gojsonrpc/DBConnection"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)
type RPC struct {
	n int
}
func (r *RPC) GetMessage(key string) ([]string,error){
	fmt.Println("***GetMessage***")
	message,err := conn.GetFromDB([]byte(key))
	if err != nil {
		return nil, err
	}
	str := strings.Split(string(message),"|")
	//removing the last element bcz empty above line will give empty value at last index
	str = str[:len(str)-1]
	return str, nil
}
func (r *RPC) AddKeyMessage(key, msg string) (string,error){
	fmt.Println("***AddMessage*** ")
	msg = msg + "|"
	err := conn.AddToDB([]byte(key), []byte(msg))
	if err != nil{
		return "Unable to receive message",err
	}
	fmt.Printf("\n%s\n","****called add message****")
	return "Successfully received the message",nil
}

func (R *RPC) DeleteALl()error{
	fmt.Println("***DRoping All*** ")
	err := conn.DropAll()
	if err != nil{
		return err
	}
	return nil
}
	var port string
var DbFlag = false
var ptr = &DbFlag
func main() {
	app := cli.NewApp()
	app.Name = "GO_JSONRPC_CLI_BADGER"
	app.Authors = []*cli.Author{
		{Name: "Tikaryan",Email: "aryantikarya9@gmail.com"},
	}

	// Receiving port number to run the server at that port
	app.Flags = []cli.Flag{
		&cli.StringFlag{Destination: &port, Name: "port",Usage: "To get the specific ports which needs to run"},
	}

	//Calling server based on command
	app.Commands = []*cli.Command{
		&cli.Command{Name: "start",Aliases: []string{"st"},Action: serverStart,Description: "running server"},
	}
	err := app.Run(os.Args)
	if err != nil{
		log.Fatalln(err)
	}

	if *ptr {
	defer conn.CloseDB()
	}
}
func serverStart(ctx *cli.Context)error{
	var gorpc = &RPC{}
	rpcServer := gojsonrpc.NewServer()
	rpcServer.Register("GORPC", gorpc)

	route := mux.NewRouter()

	fmt.Println("**starting-server***")
	route.Handle("/", rpcServer)
	httptest.NewServer(route)

	//COnnect DB
	ptr = conn.DbConnect(ptr)
	err := http.ListenAndServe(":"+port,route)

	if err != nil {
		return err
	}
	return nil
}