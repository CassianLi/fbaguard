/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sysafari.com/sysafari/fbaguard/service"

	"github.com/spf13/cobra"
)

// ebayCmd represents the ebay command
var ebayCmd = &cobra.Command{
	Use:   "ebay",
	Short: "检查Ebay DHL费用是否更新",
	Long: `通过每天访问DHL网站:https://www.dhl.de/en/privatkunden/pakete-versenden/deutschlandweit-versenden/preise-national.html, 
指定位置检查文档链接后缀中的日期是否发生变化，如果发生变化则发送通知. For example:
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ebay called")
		service.CheckEbayDHLPrice()
	},
}

func init() {
	rootCmd.AddCommand(ebayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ebayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ebayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
