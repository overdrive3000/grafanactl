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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gosuri/uitable"

	gapi "github.com/overdrive3000/go-grafana-api"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
func dashboardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard <get, list, create, delete>",
		Short: "Grafana Dashboards",
		Long: `Perform operations against Grafana Dashboards:

* Create Dashboards
* Delete Dashboards
* Search Dashboards`,
	}
	cmd.AddCommand(getDashboardCmd())
	cmd.AddCommand(createDashboardCmd())
	cmd.AddCommand(deleteDashboardCmd())

	return cmd
}

func deleteDashboardCmd() *cobra.Command {
	var uid string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Dashboard",
		Long: `Delete a grafana dashboard
grafanactl dashboard delete --uid <value>`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Deleting Dashboard")
			client, _ := SetUpClient()
			title, err := client.DeleteDashboard(uid)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Printf("Dashboard %s deleted\n", title)
		},
	}

	cmd.Flags().StringVar(&uid, "uid", "", "Dashboard UID")
	cmd.MarkFlagRequired("uid")

	return cmd
}

func createDashboardCmd() *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Dashboard",
		Long: `Creates a new Grafana Dashboard
grafanactl dashboard create --file|-f <value>`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Creating a new Dashboard")
			client, _ := SetUpClient()
			var dashboard gapi.Dashboard
			var err error
			file, err := ioutil.ReadFile(cmd.Flag("file").Value.String())
			if err != nil {
				log.Debug(err)
				os.Exit(1)
			}
			_ = json.Unmarshal(file, &dashboard)
			resp, err := client.NewDashboard(dashboard)
			if err != nil {
				log.Debug(err)
				os.Exit(1)
			}
			fmt.Printf("dashboard id %d created at %s\n", resp.ID, resp.URL)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Create dashboard from json file <value>")
	cmd.MarkFlagRequired("file")

	return cmd
}

// in progress
func getDashboardCmd() *cobra.Command {
	var uid string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Search a dashboard",
		Long: `Search a dashboard by UID
grafanactl dashboard get --uid <value>`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Getting Grafana Dashboard")
			client, _ := SetUpClient()
			var out []byte

			dashboard, err := client.GetDashboard(uid)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			switch cmd.Flag("output").Value.String() {
			case "json":
				out, err = json.Marshal(dashboard)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			case "table":
				dashboards := []gapi.Dashboard{*dashboard}
				out = dashboardAsTable(dashboards, 90)
			default:
				log.Error(errors.New(fmt.Sprintf("unknown output format %q", cmd.Flag("output").Value.String())))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, string(out))
		},
	}
	cmd.Flags().StringVar(&uid, "uid", "", "Dashboard UID to search")
	cmd.MarkFlagRequired("uid")

	return cmd
}

func dashboardAsTable(dashboards []gapi.Dashboard, colWidth uint) []byte {
	tbl := uitable.New()
	log.Debugf("%#v", dashboards)
	tbl.MaxColWidth = colWidth
	tbl.AddRow("ID", "UID", "TITLE", "FOLDER")
	for i := 0; i <= len(dashboards)-1; i++ {
		d := dashboards[i]
		tbl.AddRow(d.Model["id"], d.Model["uid"], d.Model["title"], d.Meta.FolderTitle)
	}
	return tbl.Bytes()
}
