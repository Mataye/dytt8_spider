package main

import (
	"runtime"
	_ "fmt"
	"regexp"
	"strings"
	"dytt8_spider/handle/dl"
	"dytt8_spider/util"
)





func main()  {
	//最大开两个原生线程,以达到真正的并行
	runtime.GOMAXPROCS(2)
	run()
}

var (
	//匹配出需要的html
	title_pattern = regexp.MustCompile(`<div class="title_all">(.*?)</div>`)
	//标题名字匹配
	name_pattern =  regexp.MustCompile(`<strong>(.*?)</strong>`)
	//更多按钮的url 匹配
	url_pattern =  regexp.MustCompile(`href="(.*?)"`)
)


func run()  {
	sites := util.GetSites()
	//needSource := util.GetNeedSource()

	//url := "http://" + sites[0].Url
	//body ,err := dl.GetBody(url)
	//if nil != err {
	//	fmt.Println(err)
	//	return
	//}
	//
	//title_groups := title_pattern.FindAllSubmatch(body,len(body))
	//var dataMap map[string]string
	//dataMap = make(map[string]string)
	//
	////获取有效的数据
	//for i:=0 ; i< len(title_groups) ; i++ {
	//	if len(title_groups[i]) >= 1 {
	//		mt_name, mt_url := getTitleNameAndUrl(title_groups[i][1])
	//		if "" == mt_url || "" == mt_name {
	//			continue
	//		}
	//		if isIn := In_array(mt_name,needSource);true == isIn {
	//			dataMap[mt_name] = resetUrl(mt_url,url)
	//		}
	//	}
	//}

	dl.DLrun("2018新片精品","http://www.dytt8.net/html/gndy/dyzz/index.html",sites[0].Url)
	//dl.DLrun("华语电视剧","http://www.dytt8.net/html/tv/hytv/",sites[0].Url)
	//dl.DLrun("日韩电视剧","http://www.dytt8.net/html/tv/rihantv/index.html",sites[0].Url)
	//dl.DLrun("欧美电视剧","http://www.dytt8.net/html/tv/oumeitv/index.html",sites[0].Url)

	//cto := new(sync.WaitGroup)
	//cto.Add(len(dataMap))
	//for n, v:= range dataMap  {
		// dl.DLrun(n,v,sites[0].Url)
	//}
	//cto.Wait()

}

func resetUrl(old_url,host string) (new_url string)  {
	new_url = old_url
	isOk1 := strings.Contains(old_url,host)
	if false == isOk1 {
		new_url = host + old_url
	}

	isOk2 := strings.Contains(new_url,"http://")
	if false == isOk2 {
		new_url += "http://"
	}
	return

}

func getTitleNameAndUrl(Bytes []byte) (name ,url string){
	name_groups := name_pattern.FindSubmatch(Bytes)
	url_groups := url_pattern.FindSubmatch(Bytes)

	if len(name_groups) >=1 && len(url_groups) >=1 {
		name = string(name_groups[1])
		url = string(url_groups[1])
	}
	return
}


func In_array(hack string,needle []string) bool  {
	for _, val := range needle {
		if val == hack {
			return  true
		}
	}
	return false
}

