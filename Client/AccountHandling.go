package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	httpurl  = `http://%s:%s`
	httpsurl = `https://%s:%s`
)

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func auth(opt *FlagOptions) bool {

	fmt.Println("[+] Authenticating to:", opt.host, "as", opt.username)

	data := url.Values{

		"username": {opt.username},
		"password": {opt.password},
	}

	var url string
	sess.jar = NewJar()

	if opt.ssl {

		sess.url = fmt.Sprintf(httpsurl, opt.host, strconv.Itoa(opt.port))
		url = sess.url + "/login"

	} else {

		sess.url = fmt.Sprintf(httpurl, opt.host, strconv.Itoa(opt.port))
		url = sess.url + "/login"

	}
	sess.client = http.Client{nil, nil, sess.jar, (60 * time.Second)}

	resp, err := sess.client.PostForm(url, data)
	if err != nil {

		fmt.Println("Auth Error")
		log.Fatal(err)
		return false

	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Resp Processing Error")
		log.Fatal(err)
		return false

	}

	if strings.Contains(string(body), "Welcome") {

		fmt.Println("[+] Successful Login")
		return true

	} else {

		return false

	}

}
