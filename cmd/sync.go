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

	"github.com/barkbay/zac/k8s"
	"github.com/barkbay/zac/zabbix"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronise des éléments de la configuration Zabbix en fonction de celle d'un cluster",
	Long:  `Synchronise des éléments de la configuration Zabbix en fonction de celle d'un cluster`,
	/*Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("sync called")
	},*/
}

var syncRatesCmd = &cobra.Command{
	Use:   "rates",
	Short: "Création de la configuration Zabbix pour vérifier les taux d'alerte exposé",
	Long: `Cette commande permet la création des scénarios web qui configure Zabbix pour
requêter régulièrement l'API REST exposé par la commande 'rate'.
Les triggers associés sont aussi créés pour déclencher si nécessaire des actions (e.g. mail)
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("==> Add or Update Zabbix web scenario and triggers")
		// Get REST client
		clientset := k8s.NewClientSet()

		synchronizer, err := zabbix.NewZabbixSynchronizer(clientset)
		if err != nil {
			panic(err.Error())
		}
		synchronizer.Sync()

	},
}

var syncRoutesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Création de la configuration Zabbix pour vérifier la disponibilité des routes Openshift",
	Long: `Cette commande permet la création des scénarios web qui configure Zabbix pour
requêter régulièrement les routes déclarées dans un cluster Openshift.
Les triggers associés sont aussi créés pour déclencher si nécessaire des actions (e.g. mail)
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sync routes called")
	},
}

func init() {
	syncCmd.AddCommand(syncRatesCmd)
	syncCmd.AddCommand(syncRoutesCmd)
	RootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
