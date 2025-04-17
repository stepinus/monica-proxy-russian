package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

var (
	RestySSEClient = resty.New().
			SetTimeout(3 * time.Minute).
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
			SetDoNotParseResponse(true). // 告诉 Resty，不要自动读取/解析 Body，让我们自己来处理流
			SetHeaders(map[string]string{
			"Content-Type":    "application/json",
			"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
			"x-client-locale": "zh_CN", // 可以不传，但默认会返回英文回答
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// 如果不是 200，尝试把 body 打印出来
			if resp.StatusCode() != 200 {
				return fmt.Errorf("monica API error: status %d, body: %s",
					resp.StatusCode(), resp.String())
			}
			return nil
		})

	RestyDefaultClient = resty.New().
				SetTimeout(time.Second * 30).
				SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
				SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// 如果不是 200，尝试把 body 打印出来
			if resp.StatusCode() != 200 {
				return fmt.Errorf("monica API error: status %d, body: %s",
					resp.StatusCode(), resp.String())
			}
			return nil
		})
)