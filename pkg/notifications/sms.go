package notifications

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"shop/pkg/constants"
)

func SendSms(to string, msg string) bool {
	values := url.Values{}
	values.Add("to", to)
	values.Add("msg", msg)
	values.Add("api_id", constants.SmsRuId)
	values.Add("json", "1")

	req := "https://sms.ru/sms/send?" + values.Encode()
	log.Println(req)

	resp, err := http.DefaultClient.Post(req, "application/json", nil)

	if err != nil {
		log.Printf("Failed to send sms: %s\n", err.Error())
		body := resp.Body
		defer body.Close()
		b, _ := ioutil.ReadAll(resp.Body)

		var res interface{}
		_ = json.Unmarshal(b, res)
		log.Println(b)
		log.Println(res)
		return false
	}
	return true
}

