//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto --micro_out=./proto node.proto
package main

import (
	"github.com/godcong/go-trait"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/service"
	_ "github.com/godcong/ipfs-media-service/statik"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use: "node",
}

func main() {

	configPath := rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "Config name for load config")
	elk := rootCmd.PersistentFlags().Bool("elk", true, "Log output with elk")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("elk", rootCmd.PersistentFlags().Lookup("elk"))
	Execute()

	if *elk {
		trait.InitElasticLog("ipfs-node-service", nil)
	}

	err := config.Initialize(*configPath)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//start
	service.Start()

	go func() {
		sig := <-sigs
		//bm.Stop()
		log.Info(sig, "exiting")
		service.Stop()
		done <- true
	}()
	<-done

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
