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

	dbInfo struct {
		Host string `json:"host"`
		Port int `json:"port"`
		Db string `json:"db"`
		Usr string `json:"usr"`
		Pwd string `json:"pwd"`
		DriverName string `json:"driver_name"`
		MaxOpenConns int `json:"max_open_conns"`
		MaxIdleConns int `json:"max_idle_conns"`
	}

	confInfo struct {
		Sites []sites `json:"sites"`
		NeedSource []string `json:"need_source"`
		Db dbInfo `json:"db"`
	}

)

var (
	config = confInfo{}
)


func init() {
	initConf()
	init_Db()
}

func initConf()  {
	configFile := filepath.Join(os.Getenv("GOPATH"),"src","dytt8_spider","config.json")
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



