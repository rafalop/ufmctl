package ufm

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	//"time"
	"errors"
	"bytes"
	"os"
)

func DoNothing() {
	fmt.Println("nothing.")
	os.Exit(0)
}

// need a client object that executes http requests returns json

type UfmClient struct {
	Insecure      bool
	Username      string
	Password      string
	Endpoint      string
	CurrentCookie *http.Cookie
	AuthToken	string
}

func (u *UfmClient) writeCookieFile(cookieFile string) {
	// write out current cookie to cookie file
	bytes, err := json.Marshal(u.CurrentCookie)
	if err != nil {
		fmt.Println("couldn't retrieve current cookie: ", err)
	}
	err = ioutil.WriteFile(cookieFile, bytes, 0644)
	if err != nil {
		fmt.Println("couldn't write out cookiefile: ", err)
	}
	fmt.Println("Wrote cookiefile to", cookieFile)
}

func GetClient(username string, password string, endpoint string, insecure bool, cookieFile string, authToken string) (*UfmClient, error) {
	u := &UfmClient{
		Username: username,
		Password: password,
		Endpoint: endpoint,
		Insecure: insecure,
		//AccessCookies: []*http.Cookie{},
		CurrentCookie: nil,
		AuthToken:  authToken,
	}
	if authToken == "" {
		// if we didn't supply authtoken, have to try and use cookie file
		bytes, err := ioutil.ReadFile(cookieFile)
		if err == nil && len(bytes) > 0 {
			err = json.Unmarshal(bytes, &u.CurrentCookie)
			if err != nil {
				return nil, err
				//}
				//if ! time.Now().Before(u.CurrentCookie.Expires) {
				//	fmt.Println(time.Now(), u.CurrentCookie.Expires)
				//	err := u.Auth(insecure)
				//	if err != nil {
				//		return nil, err
				//	}
				//	u.writeCookieFile(cookieFile)
			}
			if len(u.CurrentCookie.Value) > 0 {
				return u, nil
			} else {
				fmt.Fprint(os.Stderr, "Cookies file found, but cookie value is empty. Attempting re-auth with user/pass authentication\n")
			}
		} else {
			fmt.Fprint(os.Stderr, "No valid cookie file found, attempting user/pass authentication.\n")
		}

		if username == "" || password == "" {
			return nil, errors.New("No username or password not supplied, aborting.")
		}
		err = u.Auth()
		if err != nil {
			return nil, err
		}
		// write out cookie file if we are using cookies
		if u.CurrentCookie != nil {
			u.writeCookieFile(cookieFile)
		}
	} else {
		//fmt.Fprint(os.Stderr, "Authorization token supplied, using authtoken for authentication\n")
	}

	return u, nil
}

func (u *UfmClient) Auth() error {
	// if we are using token, no need to auth
	if u.AuthToken != "" {
		return nil
	}

	// try and auth for session (cookie)
	path := "/dologin"
	form := url.Values{
		"httpd_username": {u.Username},
		"httpd_password": {u.Password},
	}
	req, err := http.NewRequest("POST", u.Endpoint+path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
	}
	checkRedirect := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	client := &http.Client{Transport: tr, CheckRedirect: checkRedirect}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//bodyBytes, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(bodyBytes))
	//fmt.Println("req: ", req)
	//fmt.Println("resp: ", resp)
	u.CurrentCookie = resp.Cookies()[0]
	//fmt.Println("resp.cookies:", resp.Cookies())
	//fmt.Println("u cookies:", u.AccessCookies)
	//fmt.Println("current cookie:", u.CurrentCookie)
	return nil
}

func (u *UfmClient) AddAuth(req *http.Request) {
	// if we have token, just use it.
	if u.AuthToken != "" {
		req.Header.Set("Authorization", "Basic "+u.AuthToken)
		return
	}
	// else assume we have obtained a cookie
	req.AddCookie(u.CurrentCookie)
	return
}

// raw get with queries
func (u *UfmClient) Get(path string, queries []string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
	}
	req, err := http.NewRequest("GET", u.Endpoint+path, nil)
	u.AddAuth(req)

	q := req.URL.Query()
	for _, query := range queries {
		//fmt.Fprintln(os.Stderr, "processing query:", query)
		a := strings.Split(query, "=")[0]
		b := strings.Split(query, "=")[1]
		q.Add(a, b)
	}
	req.URL.RawQuery = q.Encode()
	//fmt.Println("req:", req)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Raw req: ", req)
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	//bodyBytes, _ := io.ReadAll(resp.Body)
	//jsonData, _ := json.Marshal(bodyBytes)
	//fmt.Println("resp:", string(bodyBytes))
	return resp, err

}

// raw put
func (u *UfmClient) Post(path string, contentType string, data string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
	}
	var ctype string
	if contentType == "" {
		ctype = "application/json"
	} else{
		ctype = contentType
	}
	req, err := http.NewRequest("POST", u.Endpoint+path, bytes.NewReader([]byte(data)))
	req.Header.Set("Content-type", ctype)
	u.AddAuth(req)

	//fmt.Println("req:", req)
	//fmt.Println("req data:", data)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	return resp, err

}

// raw delete
func (u *UfmClient) Delete(path string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
	}
	req, err := http.NewRequest("DELETE", u.Endpoint+path, nil)
	u.AddAuth(req)

	//fmt.Println("req:", req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	return resp, err

}

func (u *UfmClient) Put(path string, data io.Reader) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
	}
	req, err := http.NewRequest("PUT", u.Endpoint+path, data)
	u.AddAuth(req)

	//fmt.Println("req:", req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	return resp, err

}
