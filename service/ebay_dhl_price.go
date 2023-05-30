package service

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"path"
	"strings"
	"sysafari.com/sysafari/fbaguard/utils"
)

// CheckEbayDHLPrice 从指定链接中获取文档的链接
func CheckEbayDHLPrice() {
	uri := viper.GetString("ebay.seller-url")
	docSelector := viper.GetString("ebay.doc-selector")

	href, err := utils.GetHrefFromUrl(uri, docSelector, "href")
	if err != nil {
		log.Println(err)
	}
	parse, err := url.Parse(href)
	if err != nil {
		log.Printf("Can't parse doc url '%s': %s\n", href, err)
	}

	docname := path.Base(parse.Path)
	log.Printf("Current document filename: %s\n", docname)
	docDate := docname[strings.LastIndex(docname, "-")+1 : strings.LastIndex(docname, ".")]

	lastDate := viper.GetString("ebay.last-date")
	if docDate != lastDate {
		_, err := utils.SendDingDingRobotMessage(href, utils.Amazon)
		if err != nil {
			log.Printf("Send dingding robot message failed: %s\n", err)
		} else {
			// 更新最新文档日期
			viper.Set("ebay.last-date", docDate)
			fmt.Println("Write config file.....")
			err = viper.WriteConfig()
			if err != nil {
				fmt.Printf("Save config file failed: %v\n", err)
			}
		}

	} else {
		log.Printf("Current document(%s) is valid, dont need to update.\n", lastDate)
	}
}
