package ufm

import (
	"fmt"
	"net/http"
	"net/url"
	"crypto/tls"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
	//"time"
	"os"
)

func DoNothing() {
	fmt.Println("nothing.")
	os.Exit(0)
}
// need a client object that executes http requests returns json

type UfmClient struct {
	Insecure bool
	Username string
	Password string
	Endpoint string
	CurrentCookie *http.Cookie
}

func (u *UfmClient) writeCookieFile(cookieFile string){
	// write out current cookie to cookie file
	bytes, err := json.Marshal(u.CurrentCookie)
	if err != nil {
		fmt.Println("couldn't retrieve current cookie: ", err)
	}
	fmt.Println("writing out cookies.")
	err = ioutil.WriteFile(cookieFile, bytes, 0644)
	if err != nil {
		fmt.Println("couldn't write out cookiefile: ", err)
	}
}

func GetClient(username string, password string, endpoint string, insecure bool, cookieFile string) (*UfmClient, error){
	u:= &UfmClient{
		Username: username,
		Password: password,
		Endpoint: endpoint,
		Insecure: insecure,
		//AccessCookies: []*http.Cookie{},
		CurrentCookie: nil,
	}
	// Load cookies file check expiry
	//fmt.Println("reading cookie file")
	bytes, err := ioutil.ReadFile(cookieFile)
	if err == nil && len(bytes) > 0{
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
		//fmt.Println("succeeded using existing cookie.")
	} else {
		// Else authenticate and write out cookies file and set current cookie
		err := u.Auth()
		if err != nil {
			return nil, err
		}
		u.writeCookieFile(cookieFile)
	}
	return u, nil
}

func (u *UfmClient) Auth()(error) {
	// Check if cookies file exists and has content
	// UFM uses cookie after user/pass auth
	path := "/dologin"	
	form := url.Values{
		"httpd_username": {u.Username},
		"httpd_password": {u.Password},
	}
	//form := url.Values{}
	//form.Add("httpd_username", u.Username)
	//form.Add("httpd_password", u.Password)
	fmt.Println("form:", form)
	req, err := http.NewRequest("POST", u.Endpoint+path, strings.NewReader(form.Encode()))
	fmt.Println(req.FormValue("httpd_username"))
	fmt.Println(req.FormValue("httpd_password"))
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
	fmt.Println("sending request to", u.Endpoint+path)
	//fmt.Println("req", req)
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

// raw get with queries
func (u *UfmClient) Get(path string, queries []string) (*http.Response, error) {
	tr := &http.Transport{
        	TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
    	}
	req, err := http.NewRequest("GET", u.Endpoint+path, nil)
	req.AddCookie(u.CurrentCookie)

	q := req.URL.Query()
	for _,query := range queries {
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
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	//bodyBytes, _ := io.ReadAll(resp.Body)
	//jsonData, _ := json.Marshal(bodyBytes)
	//fmt.Println("resp:", string(bodyBytes))
	return resp, err	
	
}

// raw put
func (u *UfmClient) Post(path string, data io.Reader) (*http.Response, error) {
	tr := &http.Transport{
        	TLSClientConfig: &tls.Config{InsecureSkipVerify: u.Insecure},
    	}
	req, err := http.NewRequest("POST", u.Endpoint+path, data)
	req.AddCookie(u.CurrentCookie)

	//fmt.Println("req:", req)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//fmt.Fprintln(os.Stderr, "req:", req)
	return resp, err	
	
}


//func (u *ufmClient) Post(insecure bool, path string, args ...string) (*http.Response, error) {
//	tr := &http.Transport{
//        	TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
//    	}
//	client := &http.Client{Transport: tr}
//	contentType := "application/x-www-form-urlencoded"
//	return client.Post(u.Endpoint+path, contentType, )
//}

