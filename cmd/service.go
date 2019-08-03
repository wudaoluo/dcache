/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/dcache/service"
	"sync"
	"fmt"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg = sync.WaitGroup{}

		switch  {
		case serviceFlag.TCP:
			runService(&wg,service.NewTcpServer(internal.TCP_PORT.GetAddr(serviceFlag.Listen)))
			fallthrough

		case serviceFlag.GRPC:
			fallthrough

		default:

		}

		wg.Wait()
	},
}

func runService(wg *sync.WaitGroup,s service.Service) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Run()
	}()
}


func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serviceCmd.Flags().BoolVar(&serviceFlag.TCP,"tcp",true,"run tcp server")
	serviceCmd.Flags().BoolVar(&serviceFlag.GRPC,"grpc",false,"run grpc server")
	serviceCmd.Flags().BoolVar(&serviceFlag.ALL,"all",false,"run all server")

	serviceCmd.Flags().StringVar(&serviceFlag.Listen,"listen","0.0.0.0","listen addr")

}


var serviceFlag = &internal.Services{}