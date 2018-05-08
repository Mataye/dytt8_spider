package util

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type(
	sites struct {
		Name string `json:"name"`
		Url string  `json:"url"`
	}

	
	confInfo struct {
		Sites []sites `json:"sites"`
		NeedSource []string `json:"need_source"`
	}

)

var (
	config = confInfo{}
)


func init() {
	initConf()
}

func initConf()  {
	configFile := filepath.Join(os.Getenv("GOPATH"),"src","dytt_spider","config.json")
	info ,err := ioutil.ReadFile(configFile)
	if nil != err {
		fmt.Println(err)
		os.Exit(3)
	}
	if err := json.Unmarshal(info,&config);nil != err {
		fmt.Println(err)
		os.Exit(3)
	}
}

func GetSites() []sites {
	return  config.Sites
}

func GetNeedSource() []string  {
	return  config.NeedSource
}


