package commands

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
	"pcfdev/provisioner"
)

type DisableUAAHSTS struct {
	WebXMLPath string
}

func (d *DisableUAAHSTS) Run() error {
	var webXMLData WebApp

	webXMLContents, err := ioutil.ReadFile(d.WebXMLPath)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(bytes.NewReader(webXMLContents))
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(&webXMLData); err != nil {
		return err
	}

	hstsFilter := Filter{
		FilterName:  "httpHeaderSecurity",
		FilterClass: "org.apache.catalina.filters.HttpHeaderSecurityFilter",
		InitParam: InitParam{
			ParamName:  "hstsEnabled",
			ParamValue: "false",
		},
		AsyncSupported: true,
	}
	hstsFilterExists := false
	for _, filter := range webXMLData.Filters {
		if strings.TrimSpace(filter.FilterName) == strings.TrimSpace(hstsFilter.FilterName) &&
			strings.TrimSpace(filter.FilterClass) == strings.TrimSpace(hstsFilter.FilterClass) &&
			strings.TrimSpace(filter.InitParam.ParamName) == strings.TrimSpace(hstsFilter.InitParam.ParamName) {
			hstsFilterExists = true
		}
	}

	if hstsFilterExists {
		webXMLData.Filters = nil
	} else {
		webXMLData.Filters = []Filter{hstsFilter}
	}

	webXMLFile, err := os.OpenFile(d.WebXMLPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer webXMLFile.Close()

	encoder := xml.NewEncoder(webXMLFile)
	encoder.Indent("", "    ")
	if err := encoder.Encode(&webXMLData); err != nil {
		panic(err)
	}

	return nil
}

func (*DisableUAAHSTS) Distro() string {
	return provisioner.DistributionPCF
}

type WebApp struct {
	XMLName xml.Name `xml:"web-app"`
	Filters []Filter `xml:"filter"`
	AllXML  string   `xml:",innerxml"`
}

type Filter struct {
	FilterName     string    `xml:"filter-name"`
	FilterClass    string    `xml:"filter-class"`
	InitParam      InitParam `xml:"init-param"`
	AsyncSupported bool      `xml:"async-supported"`
}

type InitParam struct {
	ParamName  string `xml:"param-name"`
	ParamValue string `xml:"param-value"`
}
