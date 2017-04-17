// Copyright © 2017 Michael Morello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	"github.com/barkbay/zac/http"
	"github.com/barkbay/zac/k8s"
	"github.com/barkbay/zac/rate"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api/v1"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Démarre un serveur web pour exposer des métriques comme les taux d'alerte",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wRate := rate.NewRateWarningCounter()
		serve(wRate)
		listenEvents(wRate)
	},
}

func serve(wRate *rate.WarningRates) {

	// Create HTTP server
	httpServer := http.NewHttpServer(wRate)
	go httpServer.Listen()

}

func listenEvents(wRate *rate.WarningRates) {

	// Get REST client
	clientset := k8s.NewClientSet()

	errorCount := 0
	// Main event loop
	for {
		pods, err := clientset.Core().Events("").Watch(v1.ListOptions{})
		if err != nil {
			errorCount++
			if errorCount > 6*10 {
				// We have waited 10 minutes, exit now
				panic(err.Error())
			} else {
				fmt.Printf("{\"message\" : \"%s\" , \"retry\":%d}\n", err.Error(), errorCount)
				// Wait 10 seconds and retry
				time.Sleep(10 * time.Second)
			}
		} else {
			// Reset error count because cnx seems OK
			errorCount = 0
			// Ok, listen for events
			fmt.Printf("OK, listen for events\n")
			for evt := range pods.ResultChan() {
				wRate.Register(evt.Object.(*v1.Event))
			}
		}

	}
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
