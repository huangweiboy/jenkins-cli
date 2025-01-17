package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/linuxsuren/jenkins-cli/client"
	"github.com/linuxsuren/jenkins-cli/util"
	"github.com/spf13/cobra"
)

type JobTypeOption struct {
	OutputOption
}

var jobTypeOption JobTypeOption

func init() {
	jobCmd.AddCommand(jobTypeCmd)
	jobTypeCmd.Flags().StringVarP(&jobTypeOption.Format, "output", "o", "table", "Format the output")
}

var jobTypeCmd = &cobra.Command{
	Use:   "type",
	Short: "Print the types of job which in your Jenkins",
	Long:  `Print the types of job which in your Jenkins`,
	Run: func(cmd *cobra.Command, args []string) {
		jenkins := getCurrentJenkins()
		jclient := &client.JobClient{}
		jclient.URL = jenkins.URL
		jclient.UserName = jenkins.UserName
		jclient.Token = jenkins.Token
		jclient.Proxy = jenkins.Proxy
		jclient.ProxyAuth = jenkins.ProxyAuth

		if status, err := jclient.GetJobTypeCategories(); err == nil {
			var data []byte
			if data, err = jobTypeOption.Output(status); err == nil {
				if len(data) > 0 {
					fmt.Printf("%s\n", string(data))
				}
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

func (o *JobTypeOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		jobCategories := obj.([]client.JobCategory)
		table := util.CreateTable(os.Stdout)
		table.AddRow("number", "name", "type")
		for _, jobCategory := range jobCategories {
			for i, item := range jobCategory.Items {
				table.AddRow(fmt.Sprintf("%d", i), item.DisplayName,
					jobCategory.Name)
			}
		}
		table.Render()
		err = nil
		data = []byte{}
	}
	return
}
