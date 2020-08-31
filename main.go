package main

import (
	"fmt"
	gojsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/gorilla/mux"
	conn "gojsonrpc/DBConnection"
	"log"
	"net/http"
	"net/http/httptest"
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
func (R *RPC) DeleteByKey(){

}
func main() {
	var gorpc = &RPC{}

	rpcServer := gojsonrpc.NewServer()
	rpcServer.Register("GORPC", gorpc)

	route := mux.NewRouter()

	fmt.Println("**starting-server***")
	route.Handle("/", rpcServer)
	httptest.NewServer(route)
	err := http.ListenAndServe(":8080",route)
	if err != nil{
		log.Fatal("error listening",err)
	}
	defer conn.CloseDB()
}