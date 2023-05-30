package service

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"log"
	"net/url"
	"path"
	"strings"
	"sysafari.com/sysafari/fbaguard/model"
	"sysafari.com/sysafari/fbaguard/utils"
)

/**
判断Amazon FBA 是否更新计算方式
通过访问@link{https://sell.amazon.de/versand-durch-amazon?ref_=asde_soa_rd&} 网站来获取信息
*/

// sendMail 发送通知邮件
func sendMail(body string) (err error) {
	from := viper.GetStringMap("email.from")
	sender := model.MailSender{}
	err = mapstructure.Decode(from, &sender)
	if err != nil {
		return err
	}
	fmt.Println("sender:", sender)

	to := viper.GetStringSlice("email.to")
	fmt.Println("to:", to)
	if len(to) == 0 {
		return errors.New("至少设置一位收件人")
	}

	cc := viper.GetStringSlice("email.cc")
	fmt.Println("cc:", cc)

	m := gomail.NewMessage()
	m.SetHeader("From", sender.User)
	m.SetHeader("To", to...)

	for _, s := range cc {
		c := strings.Split(s, ",")
		m.SetAddressHeader("Cc", c[0], c[1])
	}

	m.SetHeader("Subject", viper.GetString("email.subject"))
	m.SetBody("text/html", body)

	d := gomail.NewDialer(sender.Host, sender.Port, sender.User, sender.Password)

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// CheckAmazonFBA Check Amazon FBA document whether to update
func CheckAmazonFBA() {
	uri := viper.GetString("amazon.seller-url")
	docSelector := viper.GetString("amazon.doc-selector")

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
	docDate := strings.Split(docname, "-")[0]

	lastDate := viper.GetString("amazon.last-date")

	if docDate != lastDate {
		//body := viper.GetString("email.body")
		//body = strings.ReplaceAll(body, "DOC_DATE", time.Now().Format("2006-01-02"))
		//
		//body = strings.ReplaceAll(body, "DOC_HREF", href)

		//err = sendMail(body)
		_, err := utils.SendDingDingRobotMessage(href, utils.Amazon)
		if err != nil {
			log.Printf("Send dingding robot message failed: %s\n", err)
		} else {
			// 更新最新文档日期
			viper.Set("amazon.last-date", docDate)
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
