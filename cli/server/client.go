/*
 * Copyright © 2019 Hedzr Yeh.
 */

package server

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://localhost:5151"

func RunClient() {
	if err := client1(); err != nil {
		logrus.Fatal("RunClient failed.", err)
	}
}

func client1() (err error) {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     idleTimeout,
	}

	if err = http2.ConfigureTransport(transport); err != nil {
		logrus.Fatal("Cannot configure http2 transport.", err)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   activeTimeout,
	}

	var resp *http.Response
	resp, err = client.Get(url)
	// resp, err = http.Get(url)

	if err != nil {
		logrus.Fatalf("Failed get: %s", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Fatal("Cannot close resp.body.", err)
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body: %s", err)
	}
	fmt.Printf(
		"Got response %d: %s, length=%d, Body is:\n%s\n",
		resp.StatusCode, resp.Proto,
		resp.ContentLength,
		string(body))
	return
}
