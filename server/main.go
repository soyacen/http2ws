package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"http2ws/conf"
	"http2ws/logger"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func StartServer() {
	var wait sync.WaitGroup
	msgChan := make(chan []byte)

	wait.Add(1)
	startHttpServer(&wait, msgChan)

	wait.Add(1)
	starWebsocketServer(&wait, msgChan)

	wait.Wait()
	logger.Info("The server is shutting down。。。")
}

func startHttpServer(wg *sync.WaitGroup, msgChan chan []byte) {
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", serveHttp(msgChan))
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HttpPort),
			Handler: mux,
		}
		logger.Info("Http Server listen on %d port", conf.HttpPort)
		logger.Fatal(server.ListenAndServe())
	}()
}

func serveHttp(msgChan chan []byte) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		msg, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		go func(msg []byte) {
			logger.Infof("server receive msg %s", msg)
			select {
			case msgChan <- msg:
			case <-msgChan:
				logger.Infof("server sweep msg %s", msg)
				msgChan <- msg
			}
		}(msg)
		writer.WriteHeader(http.StatusNoContent)
	}
}

func starWebsocketServer(wg *sync.WaitGroup, msgChan <-chan []byte) {
	go func() {
		defer wg.Done()
		var upgrader = websocket.Upgrader{
			HandshakeTimeout: 1 * time.Second,
			ReadBufferSize:   1048,
			WriteBufferSize:  1048,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				http.Error(w, reason.Error(), status)
			},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/", serveWs(msgChan, upgrader))
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.WebSocketPort),
			Handler: mux,
		}
		logger.Infof("Websocket Server listen on %d port", conf.WebSocketPort)
		logger.Fatal(server.ListenAndServe())
	}()

}

func serveWs(msgChan <-chan []byte, upgrader websocket.Upgrader) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Info(request.RemoteAddr, " is new remote addr")
		c, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Info(request.RemoteAddr, " ", err.Error())
			http.Error(writer, err.Error(), http.StatusUpgradeRequired)
			return
		}
		logger.Info("new remote addr is ",request.RemoteAddr)
		defer c.Close()
		for {
			msg := <-msgChan
			logger.Info(request.RemoteAddr, " send msg ", string(msg))
			err = c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Info(request.RemoteAddr, " send msg error", err.Error())
				break
			}
			logger.Info(request.RemoteAddr, " send msg success")
		}
	}
}
