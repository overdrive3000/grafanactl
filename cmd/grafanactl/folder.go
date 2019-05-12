// Copyright © 2019 Juan Mesa <linuxven@gmail.com>
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

package grafanactl

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Grafana Folders",
	Long: `Perform operations against Grafana Folders:

* Create Folders
* List Folders`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("folder invoked!")
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// folderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// folderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
