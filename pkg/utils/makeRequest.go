package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Todo Возможно, потребуется переопределить сигнатуру метода MakeRequest на что-то похожее как GetPage https://uproger.com/veb-skrejping-s-konkurencziej-v-golang/
func MakeRequest(siteURL string, headers map[string]string, timeout int) []byte {
	body := io.Reader(nil)
	req, err := http.NewRequest(http.MethodGet, siteURL, body)
	if err != nil {
		log.Fatalln(err)
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	// Use the default timeout if the timeout parameter isn't configured.
	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}

	// Use default http Client.
	httpClient := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       reqTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return respBody
}
