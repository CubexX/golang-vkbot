package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type Api struct {
	AccessToken string
	UserId      string
	ExpiresIn   string
}

func (vk Api) Request(methodName string, params map[string]string) string {
	nmn, err := url.Parse("https://api.vk.com/method/" + methodName)
	if err != nil {
		panic(err)
	}

	query := nmn.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	query.Set("access_token", vk.AccessToken)
	nmn.RawQuery = query.Encode()

	resp, err := http.Get(nmn.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(text)
}
