//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto node.proto
package main

import (
	"fmt"
	"github.com/godcong/go-trait"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/service"
	_ "github.com/godcong/ipfs-media-service/statik"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "node",
	}

	configPath := rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "Config name for load config")
	elk := rootCmd.PersistentFlags().Bool("elk", true, "Log output with elk")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("elk", rootCmd.PersistentFlags().Lookup("elk"))
	_ = rootCmd.Execute()
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
		fmt.Println(sig, "exiting")
		service.Stop()
		done <- true
	}()
	<-done

}
