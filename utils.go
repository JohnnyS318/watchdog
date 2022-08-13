package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
)

func CopyWeb(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: status code %d", resp.StatusCode)
		return make([]byte, 0), fmt.Errorf("error: status code: %v", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func CompareBuf(buf1, buf2 []byte) bool {
	return bytes.Equal(buf1, buf2)
}

func CheckEmailValidity(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
