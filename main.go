// http://stackoverflow.com/questions/18934327/golang-how-to-unmarshal-xml-attributes-with-colons
// http://play.golang.org/p/ZfQbJoSeQT
package main

import (
	"bytes"
	"code.google.com/p/gcfg"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Setting struct {
		Url  string
		User string
		Pw   string
	}
}

type Start1 struct {
	MultiStatus1 MultiStatus1 `xml:"multistatus"`
}
type MultiStatus1 struct {
	Hallo string `xml:"hallo"`
}

func getAddressBook(cfg Config) []byte {
	url := cfg.Setting.Url
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`<card:addressbook-query xmlns:d="DAV:" xmlns:card="urn:ietf:params:xml:ns:carddav">
    <d:prop>
        <d:getetag />
        <card:address-data />
    </d:prop>
	</card:addressbook-query>`)
	req, err := http.NewRequest("REPORT", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Depth", "1")
	req.Header.Add("Content-Type:", "application/xml; charset=utf-8")
	client := &http.Client{}
	req.SetBasicAuth(cfg.Setting.User, cfg.Setting.Pw)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//	fmt.Println("response Status:", resp.Status)
	//	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println("response Body:", string(body))
	// fmt.Println(string(body))
	return body
}

func main() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "myconfig.cfg")
	if err != nil {
		fmt.Println("hallo")
		fmt.Printf("error: %v", err)
		return
	}
	dada := getAddressBook(cfg)
	fmt.Println(string(dada))
	str1 := `<multistatus>
	<hallo>geht was</hallo>
</multistatus>`
	fmt.Println(str1)
	v := Start1{}
	err = xml.Unmarshal([]byte(str1), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	v.MultiStatus1.Hallo = "test"
	fmt.Println(v)
}
