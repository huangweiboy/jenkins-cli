package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

//curl -k -u $JENKINS_USER:$JENKINS_TOKEN $JENKINS_URL/crumbIssuer/api/json -s

// Start contains the command line options
type CrumbIssuerOptions struct {
	Upload bool
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Print crumbIssuer of Jenkins",
	Long:  `Print crumbIssuer of Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		crumb, _ := getCrumb()
		fmt.Printf("%s=%s", crumb.CrumbRequestField, crumb.Crumb)
	},
}

type CrumbIssuer struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func getCrumb() (CrumbIssuer, *JenkinsServer) {
	config := getCurrentJenkins()

	jenkinsRoot := config.URL
	api := fmt.Sprintf("%s/crumbIssuer/api/json", jenkinsRoot)

	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.UserName, config.Token)

	var crumbIssuer CrumbIssuer
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if response, err := client.Do(req); err == nil {
		if data, err := ioutil.ReadAll(response.Body); err == nil {
			if response.StatusCode == 200 {
				json.Unmarshal(data, &crumbIssuer)
			} else {
				fmt.Println("get curmb error")
				log.Fatal(string(data))
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return crumbIssuer, config
}
