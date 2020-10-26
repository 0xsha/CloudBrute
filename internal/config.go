package internal

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// Providers []string{"amazon","alibaba","amazon","microsoft","digitalocean","linode","vultr"}

type Config struct {
	Author       string   `yaml:"author"` //
	IPInfo       string   `yaml:"ipinfo"` // API KEY
	ProxyType    string   `yaml:"proxytype"`
	Providers    []string `yaml:"providers"`
	Environments []string `yaml:"Environments"`
}

type RequestDetails struct {
	ProxyList   []string
	ProxyType   string
	RandomAgent []string
}

type CloudConfig struct {
	Regions     []string `yaml:"regions"`
	APPUrls     []string `yaml:"app_urls"`
	StorageUrls []string `yaml:"storage_urls"`
	//RegionUrls  []string `yaml:"region_urls"`
	AppRegionUrls  []string `yaml:"app_region_urls"`
	StorageRegionUrls  []string `yaml:"storage_region_urls"`
}

func InitConfig(path string) *Config {

	//log.Print("hello world")
	var config Config
	filename, _ := filepath.Abs(path)
	configFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func InitCloudConfig(cloud string, modulePath string) (*CloudConfig, error) {
	if !strings.Contains(cloud,".yaml") {
		cloud = cloud + ".yaml"
	}
	var cloudConfig CloudConfig
	filename, _ := filepath.Abs(path.Join(modulePath, cloud))
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configFile, &cloudConfig)
	if err != nil {
		return nil, err
	}
	return &cloudConfig, nil

}
