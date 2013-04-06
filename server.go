package main

import(
	"fmt"
	"net/http"
	"log"
	"time"
	"os"
)

type HandleFnc func(http.ResponseWriter, *http.Request)

func LogRequest(r *http.Request){
	var logText string = fmt.Sprintf(" - %v", r.URL)
	log.Println(logText)
}

func LogPanic(handlerFunction HandleFnc) HandleFnc{	
	return func(w http.ResponseWriter, r *http.Request){
				go LogRequest(r)

				defer func(){
					if err := recover(); err != nil{
						var errorMessage string = fmt.Sprintf("[Panic] %v - %v: %v", time.Now(), r.URL, err)
						log.Println("[Panic] %v", errorMessage)
					}
				}()
				handlerFunction(w, r)
			}
}


func HandleHome(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Connection", "keep-alive")
	fmt.Fprintf(w, "Hello World - %v", time.Now())
}

func initializeLogger(fileName string){
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
 	if err != nil{
		log.Println(err.Error())
 	}

	log.SetOutput(f)
}

func main(){
	fmt.Println("Starting Server...")
	initializeLogger("log/server.log");

	http.HandleFunc("/", LogPanic(HandleHome))

	if err := http.ListenAndServe(":8799", nil); err != nil{
		panic(err)
	}
}