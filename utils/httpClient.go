package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	POST_DATA_TYPE_JSON = 1
	POST_DATA_TYPE_FORM = 2
)

type HttpClient struct {
	postData      map[string]interface{} //
	postContents  []byte
	headers       map[string]string
	timeOut       time.Duration
	postDataType  int
	useGZip       bool
	receiveCookie string
	queryMap      map[string]interface{}
}

func NewHttpClient() *HttpClient {
	hClient := HttpClient{}
	hClient.timeOut = time.Second * 30
	hClient.postDataType = POST_DATA_TYPE_FORM
	return &hClient
}

func (hClient *HttpClient) AddQuery(key string, value interface{}) {
	if hClient.queryMap == nil {
		hClient.queryMap = make(map[string]interface{})
	}
	hClient.queryMap[key] = value
}

func (hClient *HttpClient) getQuery() string {
	if hClient.queryMap == nil {
		return ""
	}
	str := ""
	for k, v := range hClient.queryMap {
		if str != "" {
			str = fmt.Sprintf("%v&%v=%v", str, k, v)
		} else {
			str = fmt.Sprintf("%v=%v", k, v)
		}
	}
	return str
}

// Set contents type
func (hClient *HttpClient) SetPostDataType(dataType int) {
	hClient.postDataType = dataType
}

func (hClient *HttpClient) SetPostData(postData interface{}) {
	switch vv := postData.(type) {
	case string:
		hClient.postContents = []byte(vv)
	default:
		hClient.postContents, _ = json.Marshal(postData)
	}
}

// add
func (hClient *HttpClient) AddFormData(key string, value interface{}) {
	if hClient.postData == nil {
		hClient.postData = make(map[string]interface{})
	}
	hClient.postData[key] = value
}

func (hClient *HttpClient) SetCookie(cookie string) {
	hClient.AddHeader("Cookie", cookie)
}

func (hClient *HttpClient) GetCookie() string {
	return hClient.receiveCookie
}

func (hClient *HttpClient) AddHeader(key, value string) {
	if hClient.headers == nil {
		hClient.headers = make(map[string]string)
	}
	hClient.headers[key] = value
}

//
func (hClient *HttpClient) EncodingGZip(bUse bool) {
	hClient.useGZip = bUse
}

// Post
func (hClient *HttpClient) Post(link string) ([]byte, error) {
	if hClient.postContents == nil || len(hClient.postContents) == 0 {
		hClient.postContents = hClient.GetPostData()
	}
	return hClient.do("POST", link, hClient.postContents)
}

func (hClient *HttpClient) Get(link string) ([]byte, error) {
	strForm := string(hClient.GetPostData())
	if strForm != "" {
		if !strings.Contains(link, "?") {
			link = link + "?"
		} else {
			if link[len(link)-1] != '&' {
				link += "&"
			}
		}
		link = link + strForm
	}
	return hClient.do("GET", link, nil)
}

func (hClient *HttpClient) GetPostData() []byte {
	if hClient.postData == nil || len(hClient.postData) == 0 {
		return []byte("")
	}

	if hClient.postDataType == POST_DATA_TYPE_JSON {
		data, _ := json.Marshal(hClient.postData)

		// clean postdata
		hClient.postData = nil
		return data
	} else {
		var data string
		for key, value := range hClient.postData {
			separate := "&"
			if len(data) == 0 {
				separate = ""
			}
			data += fmt.Sprintf("%s%s=%v", separate, key, value)
		}
		// clean postdata
		hClient.postData = nil
		return []byte(data)
	}
}

func (hClient *HttpClient) SetReferer(refUrl string) {
	hClient.AddHeader("Referer", refUrl)
}

func (hClient *HttpClient) setHeaders(request *http.Request) {
	for k, v := range hClient.headers {
		request.Header.Set(k, v)
	}
}

func (hClient *HttpClient) do(method string, link string, data []byte) ([]byte, error) {
	queryParams := hClient.getQuery()
	if queryParams != "" {
		if !strings.Contains(link, "?") {
			link += "?" + queryParams
		} else {
			link += "&" + queryParams
		}
	}

	var request *http.Request
	var err error
	if data != nil && len(data) != 0 {
		// gzip
		if hClient.useGZip {
			var zBuf bytes.Buffer
			zipWrite := gzip.NewWriter(&zBuf)

			if _, err = zipWrite.Write(data); err != nil {
				fmt.Println("-----gzip is faild,err:", err)
			}
			zipWrite.Close()
			request, err = http.NewRequest(method, link, &zBuf)
			request.Header.Add("Content-Encoding", "gzip")
			//request.Header.Add("Accept-Encoding", "gzip")
		} else {
			request, err = http.NewRequest(method, link, bytes.NewReader(data))
		}
	} else {
		request, err = http.NewRequest(method, link, nil)
	}

	// clean postdata
	// hClient.postContents = nil

	if err != nil {
		return nil, err
	} else {
		/*netClient := &http.Client{
			Timeout: hClient.timeOut,
		}
		var transport *http.Transport = nil
		if true {
			URL := url.URL{}
			urlProxy, _ := URL.Parse("http://127.0.0.1:8888")
			transport = &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			}
		} else {
			transport = &http.Transport{}
		}*/

		netClient := &http.Client{
			Timeout: hClient.timeOut,
			//Transport: transport,
		}

		// set header
		hClient.setHeaders(request)

		if response, err := netClient.Do(request); err != nil {
			return nil, err
		} else {
			// save recevie cookie
			for _, v := range response.Cookies() {
				separate := "; "
				if hClient.receiveCookie == "" {
					separate = ""
				}
				hClient.receiveCookie += fmt.Sprintf("%s%s=%s", separate, v.Name, v.Value)
			}

			data, err := ioutil.ReadAll(response.Body)
			response.Body.Close()

			if err == nil {
				// gzip decompress
				if strings.Contains(response.Header.Get("Accept-Encoding"), "gzip") {
					gzipReader, err := gzip.NewReader(bytes.NewReader(data))
					if err != nil {
						return data, nil
					}
					unBody, err := ioutil.ReadAll(gzipReader)
					gzipReader.Close()

					if err != nil {
						return data, nil
					} else {
						return unBody, nil
					}
				}
				return data, err
			} else {
				return nil, err
			}
		}
	}
}
