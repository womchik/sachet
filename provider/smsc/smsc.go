package smsc

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"github.com/messagebird/sachet"
)

type SmscConfig struct {
	User     string `yaml:"username"`
	Password string `yaml:"password"`
}

const SmscRequestTimeout = time.Second * 60

type Smsc struct {
	User     string
	Password string
	ApiId    string
}

func NewSmsc(config SmscConfig) *Smsc {
	Smsc := &Smsc{User: config.User, Password: config.Password, ApiId: config.ApiId}
	return Smsc
}

func (c *Smsc) Send(message sachet.Message) (err error) {
	for _, number := range message.To {
		err = c.SendOne(message, number)
		if err != nil {
			return fmt.Errorf("Failed to make API call to smsc:%s", err)
		}
	}
	return
}

func (c *Smsc) SendOne(message sachet.Message, PhoneNumber string) (err error) {
	fmt.Printf("ALERT : %s\n", message.Text)
	encoded_message := url.QueryEscape(message.Text)
	smsURL := fmt.Sprintf("https://smsc.ru/sys/send.php?login=%s&psw=%s&phones=%s&mes=%s", c.User, c.Password, PhoneNumber, encoded_message)
	var request *http.Request
	var resp *http.Response
	request, err = http.NewRequest("GET", smsURL, nil)
	if err != nil {
		return
	}
	httpClient := &http.Client{}
	httpClient.Timeout = SmscRequestTimeout
	resp, err = httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var body []byte
	resp.Body.Read(body)
	if resp.StatusCode == http.StatusOK && err == nil {
		return
	}
	return fmt.Errorf("Failed sending sms:Reason: %s, StatusCode : %d", string(body), resp.StatusCode)
}