package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ethereum struct {
	EthereumAddress string `json:"ethereum_address"`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("RootHandler started and redirecting to the index page")

	t, _ := template.ParseFiles(INDEX)

	t.Execute(w, nil)

}

func QueryHandler(w http.ResponseWriter, r *http.Request) {

	//Retrieve the HTML form parameter of POST method
	//	e := r.FormValue("ethereum-data")

	//	g, err := GetGist(user)

	//	if err != nil {
	//		log.Printf("GetGistError: %s, %v",
	//			err.Error(), http.StatusInternalServerError)
	//
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}

	t, err := template.ParseFiles(INDEX)

	if err != nil {
		fmt.Println(err.Error())
	}

	g := GetEthereumData()
	fmt.Println(g)

	t.Execute(w, g)

}

func ApiHandler(w http.ResponseWriter, r *http.Request) {

	//	decoder := json.NewDecoder(r.Body)

	//	var request ethereum

	//	err := decoder.Decode(&request)

	//	g, err := GetGist(request.User)

	//	if err != nil {
	//		log.Printf("GetGistError: %s, %v",
	//			err.Error(), http.StatusInternalServerError)
	//
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}

	//	if _, exists := GlobalStore[request.User]; !exists {
	//		GlobalStore[request.User] = g.Counter
	//	}
	//
	//	client, err := rpc.Dial("https://mainnet.infura.io/v3/70ca5e48e7cb47079e0377062b461ec2")
	//	if err != nil {
	//		fmt.Println("rpc.Dial err", err)
	//		return
	//	}

	//	var result map[string]string
	//	err = client.Call(&result, "eth_getTransactionByBlockNumberAndIndex", "latest", "0x0")

	//	if err != nil {
	//		fmt.Println("client.Call err", err)
	//		return
	//	}
	//  fmt.Printf("%+v", result)
	//fmt.Printf("accounts: %s\n", account[0])
	//	jsonData, err := json.Marshal(result)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	g := GetEthereumData()
	jsonData, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err.Error())
	}
	//	fmt.Fprint(w, bytes.NewBuffer(g["eth_getTransactionByBlockNumberAndIndex"]))
	fmt.Fprint(w, string(jsonData))
	return
}
