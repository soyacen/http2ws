/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"http2ws/conf"
	"http2ws/logger"
	"http2ws/server"
)

var rootCmd = &cobra.Command{
	Use:   "http2ws",
	Short: "把 http 转发到 websocket",
	Long: `当用户A请求 http 服务时，把所带的数据在通过 websocket 服务发送给用户B`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init()
		server.StartServer()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVar(&conf.HttpPort, "httpPort", 4567, "http服务端口")
	rootCmd.PersistentFlags().IntVar(&conf.WebSocketPort, "wsPort", 5678, "websocket服务端口")
	rootCmd.PersistentFlags().StringVar(&conf.LogFile, "logFile", "./http2ws.log", "日志文件路径")

}