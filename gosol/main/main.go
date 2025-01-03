package main

import (
	"fmt"
	"github.com/slawomir-pryczek/HSServer/handler_socket2"
	"github.com/slawomir-pryczek/HSServer/handler_socket2/config"
	"github.com/slawomir-pryczek/HSServer/handler_socket2/handle_echo"
	"github.com/slawomir-pryczek/HSServer/handler_socket2/handle_profiler"
	"gosol/handle_kvstore"
	"gosol/passthrough"
	pluginmanager "gosol/plugins"
	"gosol/solana/handle_solana_01"
	"gosol/solana/handle_solana_admin"
	"gosol/solana/handle_solana_info"
	"gosol/solana/handle_solana_raw"
	"os"
	"runtime"
	"strings"
)

func readNodeConfig() {

	fmt.Println("\nReading node config...")
	nodes := (config.Config().GetRawData("SOL_NODES", "")).([]interface{})
	if len(nodes) <= 0 {
		fmt.Println("ERROR: No nodes defined, please define at least one solana node to connect to")
		os.Exit(10)
		return
	}

	for _, v := range nodes {
		handle_solana_admin.NodeRegisterFromConfig(v.(map[string]interface{}))
	}
	fmt.Println("")
}

func main() {

	pluginmanager.RegisterAll()
	readNodeConfig()

	numCpu := runtime.NumCPU() * 2
	runtime.GOMAXPROCS(numCpu) // register handlers
	var handlers []handler_socket2.ActionHandler
	handlers = append(handlers, &handle_echo.HandleEcho{})
	handlers = append(handlers, &handle_profiler.HandleProfiler{})
	handlers = append(handlers, &handle_solana_raw.Handle_solana_raw{})
	handlers = append(handlers, &handle_solana_01.Handle_solana_01{})
	handlers = append(handlers, &handle_solana_info.Handle_solana_info{})
	handlers = append(handlers, &handle_passthrough.Handle_passthrough{})
	handlers = append(handlers, &handle_solana_admin.Handle_solana_admin{})
	handlers = append(handlers, &handle_kvstore.Handle_kvstore{})

	if len(config.Config().Get("RUN_SERVICES", "")) > 0 && config.Config().Get("RUN_SERVICES", "") != "*" {
		var hModified []handler_socket2.ActionHandler
		_tmp := strings.Split(config.Config().Get("RUN_SERVICES", ""), ",")
		supported := make(map[string]bool)
		for _, v := range _tmp {
			supported[strings.Trim(v, "\r\n \t")] = true
		}

		for _, v := range handlers {
			shouldEnable := false
			for _, action := range v.GetActions() {
				if supported[action] {
					shouldEnable = true
					break
				}
			}

			if shouldEnable {
				hModified = append(hModified, v)
			}
		}

		handlers = hModified
	}

	// start the server
	handler_socket2.RegisterHandler(handlers...)
	handler_socket2.StartServer(strings.Split("0.0.0.0:7777,h0.0.0.0:7778", ","))
}
