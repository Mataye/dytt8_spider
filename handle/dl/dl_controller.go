package dl

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"bytes"
	"sync"
	"strings"
	"dytt8_spider/util"
	"encoding/json"
	"github.com/spacemonkeygo/errors"
)





type (
	Movie struct {
		m_name string
		m_ftps []string
	}


	Video struct {
		M_name   string
		Ftps_url string
	}
)

//func DLrun(title,url ,host string,cto *sync.WaitGroup)   {
//
//	next_url := getNextPage(url)
//	fmt.Println(next_url)
//	go DLdata(title,url ,host)
//	if "" != next_url {
//	  DLrun(title,next_url,host,cto)
//	}
//
//	cto.Done()
//	return
//}
func DLrun(title,url ,host string) {

	next_url := getNextPage(url)
	fmt.Println(next_url)
	go DLdata(title, url, host)
	if "" != next_url {
		DLrun(title, next_url, host)
	}
	return
}


func DLdata(title,url ,host string) {
	doc,err := getDoc(url)
	if nil != err {
		fmt.Print("create list doc faild")
		fmt.Println(err)
		return
	}

	var t_urls []string
	doc.Find(".co_content8 b").Each(func(i int, s *goquery.Selection) {
		d := s.Find("a")
		//t_name := d.Text()
		t_url := d.AttrOr("href","")
		t_url = "http://"+host + t_url
		t_urls = append(t_urls,t_url)
	})

	//获取当前页的所有数据
	//设定管道和锁
	control := new(sync.WaitGroup)
	data_channel := make(chan Movie,len(t_urls))
	control.Add(len(t_urls))
	for i:=0; i<len(t_urls) ; i++  {
		go getDetailPageInfo(control,t_urls[i],data_channel)
	}
	control.Wait()
	close(data_channel)
	//var current_movies  []Movie
	if 0 == len(data_channel) {
		fmt.Println("empty Data")
		return
	}


	mydb := util.GetDb()
	var insertData Video
	for val:= range data_channel {
		insertData = Video{}
		insertData = combineInsertData(val)
		//if  != insertData {
			if err := mydb.Create(insertData).Error; nil != err {
				fmt.Println(err)
			}
		//}

	//	insertData = combineInsertData(insertData,val)
		//current_movies = append(current_movies , i)
	}

	//if nil != insertData {
	//	mydb := util.GetDb()
	//	if err := mydb.Create(insertData).Error; nil != err {
	//		fmt.Println(err)
	//	}
	//}
	dd ,_ := json.Marshal(insertData)
	fmt.Println(string(dd))
	close(data_channel)
}


func combineInsertData(detail  Movie)  Video {

	jsonStr ,_ := json.Marshal(detail.m_ftps)
	//if nil ==err {
		mid := Video{detail.m_name,string(jsonStr)}
		//insertData = append(insertData,mid)
	//}
	return mid
	//return insertData
}

//func combineInsertData(insertData []InsertData,detail  Movie)  []InsertData {
//
//	jsonStr ,err := json.Marshal(detail.m_ftps)
//	if nil == err {
//		mid := InsertData{detail.m_name,string(jsonStr)}
//		insertData = append(insertData,mid)
//	}
//	return insertData
//}

func getNextPage(url string) (next_url string) {
	doc,err := getDoc(url)
	if nil != err {
		fmt.Print("create list doc faild ")
		fmt.Println(err)
		return
	}

	doc.Find(".co_content8 .x a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		next_path := s.AttrOr("href","")
		btn_name := s.Text()
		if "" != next_path && "下一页" == btn_name {
			next_url = combineNewUrl(url,next_path)
			return false
		}
		return  true
	})
	return
}


func combineNewUrl(old_url , new_path string) (new_url string) {
	mid := strings.Split(old_url,"/")
	mid =mid[1:len(mid)-1]

	new_url = ""
	new_url =  strings.Join(mid,"/")
	new_url = "http:/"+new_url +"/"+ new_path
	return

}

//func getDetailPageInfo(control *sync.WaitGroup,t_url string) {
func getDetailPageInfo(control *sync.WaitGroup,t_url string,data_channel chan Movie) {
	doc,err := getDoc(t_url)

	if nil != err || nil == doc {
		fmt.Print("create detail page "+t_url +" doc faild ")
		fmt.Println(err)
		control.Done()
		return
	}

	letfDoc :=  doc.Find(".bd3r")

	if nil == letfDoc {
		fmt.Println( t_url + "can not find .db3r")
		control.Done()
		return
	}

	title :=  letfDoc.Find(".title_all").Text()
	//ftps := letfDoc.Find("#Zoom").Find("table").Find("tbody")

	var  ftps = Movie{}
	ftps.m_name = title
	letfDoc.Find("#Zoom").Find("table").Each(func(i int, s *goquery.Selection) {
		d_url := s.Find("a").AttrOr("href","")
		//d_name := s.Find("a").Text()
		if d_url != "" {
			ftps.m_ftps = append(ftps.m_ftps,d_url)
		} else  {
			fmt.Println("****" +title)
		}
	})
	data_channel<-ftps
	//fmt.Println(ftps)
	control.Done()
	return
}



func getDoc(url string) (doc *goquery.Document,err error) {
	//var (
	//	ctx context.Context
	//	cancel context.CancelFunc
	//)
	//ctx,cancel = context.WithTimeout(context.Background(),30 * time.Second)
	//defer  cancel()
	var body []byte
	body, err = GetBody(url)
	if nil != err{
		return
	}
	if  nil == body {
		err = errors.New("get Page " + url + " Data Empty")
		return
	}


	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if nil != err {
		return
	}


	return
}
