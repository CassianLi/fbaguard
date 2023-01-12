/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sysafari.com/sysafari/fbaguard/service"

	"github.com/spf13/cobra"
)

// amazonCmd represents the amazon command
var amazonCmd = &cobra.Command{
	Use:   "amazon",
	Short: "获取Amazon的FBA计算方式的更新通知",
	Long: `通过访问网站：https://sell.amazon.de/versand-durch-amazon?ref_=asde_soa_rd& ，
判断FBA费用计算的文档来判断当前是否更新了FBA计算方法，如已更新，则发送邮件给指定人员. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("amazon called")

		service.CheckAmazonFBA()
	},
}

func init() {
	rootCmd.AddCommand(amazonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// amazonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// amazonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
