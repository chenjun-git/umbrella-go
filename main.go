package main

import (
	"flag"
	"fmt"
	"net/http"

	"umbrella-go/umbrella-common/monitor"

	"umbrella-go/common"
	"umbrella-go/handler"
	"umbrella-go/utils"
)

var (
	BuildTime    = "No Build Time"
	BuildGitHash = "No Build Git Hash"
	BuildGitTag  = "No Build Git Tag"
)

func main() {
	initConfig()
	initMonitor()

	startHTTPServer()
}

func initConfig() {
	configPath := flag.String("config", "", "config file's path")
	flag.Parse()

	common.InitConfig(*configPath)
}

func initMonitor() {
	monitor.Init(common.Config.Monitor.Namespace, common.Config.Monitor.Subsystem)
	monitor.Monitor.SetVersion(monitor.Version{GitHash: BuildGitHash, GitTag: BuildGitTag, BuildTime: BuildTime})
	monitor.InitAndListen(common.Config.HTTP.Listen)
}

func startHTTPServer() {
	router := handler.RegisterBackendRouter()
	utils.RegisterPProf(router)
	monitor.RegisterHandlers(router)

	httpServer := &http.Server{
		Addr:    common.Config.HTTP.Listen,
		Handler: router,
	}
	fmt.Printf("start http server, listen: %s!\n", common.Config.HTTP.Listen)
	if err := httpServer.ListenAndServe(); err != nil {
		panic("start http server failed")
	}
}
