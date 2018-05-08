package dl

import (
	"net/http"
	"github.com/spacemonkeygo/errors"
	"strconv"
	"github.com/cihub/seelog/archive/gzip"
	"compress/flate"
	"io"
	"compress/zlib"
	"io/ioutil"
	"strings"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"fmt"
	"golang.org/x/text/transform"
	"bytes"
)

func GetBody(url string) (Bytes []byte ,err error)   {
	resp ,err := http.Get(url)

	if nil != err {
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 304 {
		err = errors.New("Request " + url + "Faild! StatusCode " + strconv.Itoa(resp.StatusCode))
		return
	}

	Bytes ,err = chaneEncoding(resp)
	return 
}

//页面解压.修改编码，防止乱码
func chaneEncoding(resp *http.Response) (Byte []byte,err error) {
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var gzipReader *gzip.Reader
		gzipReader ,err = gzip.NewReader(resp.Body,"gzip")
		if nil != err {
			return
		}
		defer gzipReader.Close()
		resp.Body = gzipReader
		break

	case "deflate":
		resp.Body = flate.NewReader(resp.Body)
		break
	case "zlib":
		var readCloser io.ReadCloser
		readCloser,err = zlib.NewReader(resp.Body)
		if nil != err {
			return
		}
		defer readCloser.Close()
		resp.Body = readCloser
		break
	}

	Byte,err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return
	}

	text := string(Byte)
	//fmt.Println(text)
	Byte = []byte(Change(text,GetCodeFormat(text)))
	return  
}


//去除&nbsp; 替换&ldquo;和&rdquo;
func replaceSpecialCharacter(text string) (newText string) {

	newText = strings.Replace(text,"&nbsp;","",-1)
	//替换&ldquo;和&rdquo;
	newText = strings.Replace(text, "“", "\"", -1)
	newText = strings.Replace(text, "”", "\"", -1)
	newText = strings.Replace(text, "…", "...", -1)
	return
}


//改变编码格式
func Change(text ,codeFormat string) (newText string) {
	text =  replaceSpecialCharacter(text)
	if "UTF-8" != codeFormat {
		var reader *transform.Reader

		reader = transform.NewReader(bytes.NewReader([]byte(text)), simplifiedchinese.GB18030.NewDecoder())
		Bytes, _ := ioutil.ReadAll(reader)
		text = string(Bytes)
		//text = mahonia.NewEncoder("gbk").ConvertString(text)
	}
	newText = text
	return
}

//获取页面编码格式
func GetCodeFormat(text string) (codeFormat string)  {
	detector := chardet.NewTextDetector()
	result ,err := detector.DetectBest([]byte(text))
	if nil == err {
		if strings.HasPrefix(result.Charset, "GB") &&
			(strings.HasSuffix(result.Charset, "18030") ||
				strings.HasSuffix(result.Charset, "2312") ||
				strings.HasSuffix(result.Charset, "13000")) {
			codeFormat = strings.Replace(result.Charset, "-", "", -1)
		} else {
			codeFormat = result.Charset
		}
	}
	fmt.Println(codeFormat)
	return
}

