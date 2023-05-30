package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SalesChannel uint

const (
	Amazon SalesChannel = iota
	Ebay
)

const (
	// DingDingRobotURL 钉钉机器人的URL
	DingDingRobotURL = "https://oapi.dingtalk.com/robot/send?access_token=027cd5e031ea02df11334baa7bcd41bf378b411ad92076f66244bbcd2d182a5c"

	// DingDingRobotSecret 钉钉机器人的Secret
	DingDingRobotSecret = "SEC2bb7fedd4438bb2202b9a0f60a9d114845ee8a5fb16f9de9d7902bdb40fa7848"
)

// AmazonFBANotifyBody 亚马逊FBA更新通知的消息体
var AmazonFBANotifyBody = `{
  "at": {
    "isAtAll": true
  },
  "markdown": {
    "title": "亚马逊FBA更新",
    "text": "# 亚马逊FBA文档更新 \n 更新日期：*UPDATE_DATE*更新，请及时根据更新的文档内容，调整**Category FBA 计算器**，详细内容：\n - 更新后文档链接：DOCUMENT_LINK \n - 详情请查看亚马逊卖家中心：https://sell.amazon.de/versand-durch-amazon?ref_=asde_soa_rd& \n \n 消息来自：*fbaguard（计算文档监控模块）* \n"
  },
  "msgtype": "markdown"
}`

// EbayFBANotifyBody Ebay Fulfillment更新通知的消息体
var EbayFBANotifyBody = `{
  "at": {
    "isAtAll": true
  },
  "markdown": {
    "title": "Ebay Fulfillment 计算更新",
    "text": "# Ebay DHL Price 更新 \n 更新日期：*UPDATE_DATE*，请及时根据更新的文档内容，调整**Ebay 的Fulfillment 计算**，详细内容：\n - 更新后文档链接：DOCUMENT_LINK \n - 详情请查看DHL网站详细内容：https://www.dhl.de/en/privatkunden/pakete-versenden/deutschlandweit-versenden/preise-national.html \n \n 消息来自：*fbaguard（计算文档监控模块）* \n"
  },
  "msgtype": "markdown"
}
`

// SignURL 传入URL，并为其增加参数，获取timestamp，使用timestamp+"\n"+secret，使用HmacSHA256算法计算签名，将参数进行urlEncode，最后将签名作为参数加入到URL中
func SignURL(uri string, secret string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	text := timestamp + "\n" + secret
	signature := HmacSha256(text, secret)
	return uri + "&timestamp=" + timestamp + "&sign=" + url.QueryEscape(signature)
}

// HmacSha256 传入text和key，使用HmacSha256算法计算签名
func HmacSha256(text string, key string) string {
	// 1.计算HMAC-SHA256
	hmacEncoder := hmac.New(sha256.New, []byte(key))
	hmacEncoder.Write([]byte(text))
	// 2.base64 encode
	return base64.StdEncoding.EncodeToString(hmacEncoder.Sum(nil))
}

// SendRequest 创建方法通过url发送请求，提交body，并返回响应结果
func SendRequest(method string, url string, body string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	fmt.Println("response Status:", resp.Status)
	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, errors.New("Request failed with status code: " + resp.Status)
}

// SendDingDingRobotMessage 发送钉钉机器人消息,传入更新后的文档链接
func SendDingDingRobotMessage(documentLink string, channel SalesChannel) (bool, error) {
	// 1.获取签名链接
	signedURL := SignURL(DingDingRobotURL, DingDingRobotSecret)

	var body string
	switch channel {
	case Amazon:
		body = AmazonFBANotifyBody
	case Ebay:
		body = EbayFBANotifyBody
	default:
		return false, errors.New("unknown channel")

	}

	// 2.替换消息体中的文档链接和更新日期
	notifyBody := strings.Replace(body, "DOCUMENT_LINK", documentLink, 1)
	notifyBody = strings.Replace(notifyBody, "UPDATE_DATE", time.Now().Format("2006-01-02"), 1)

	// 3.发送请求
	return SendRequest("POST", signedURL, notifyBody)
}
