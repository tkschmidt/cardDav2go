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

type Prop struct {
	Status string `xml:"status"`
}
type Propstat struct {
	Prop Prop `xml:"prop"`
}

type Response struct {
	Href     string   `xml:"href"`
	Propstat Propstat `xml:"propstat"`
}
type Result struct {
	XMLName   xml.Name `xml:"multistatus"`
	Responses Response `xml:"response"`
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
	//dada := getAddressBook(cfg)

	data1 := `<?xml version="1.0" encoding="utf-8"?>
<d:multistatus xmlns:card="urn:ietf:params:xml:ns:carddav" xmlns:d="DAV:" xmlns:s="http://sabredav.org/ns">
	<d:response>
		<d:href>/remote.php/carddav/addressbooks/test/kontakte/cd28ea0d-83e1-45c0-8ce2-d29dfa3bf9b9%40cloud.tkschmidt.me.vcf</d:href>
		<d:propstat>
			<d:prop>
			<d:status>HTTP/1.1 200 OK</d:status>
			</d:prop>
		</d:propstat>
	</d:response>
</d:multistatus>`

	v := Result{}
	//err = xml.Unmarshal(data1, &v)
	err = xml.Unmarshal([]byte(data1), &v)

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println("##########################")
	//	fmt.Printf("XMLName: %#v\n", v.XMLName)
	fmt.Printf("%#v\n", v)
	fmt.Printf("%#q\n", v)
}
