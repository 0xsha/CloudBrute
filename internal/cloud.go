package internal

import (
	"errors"
	"github.com/ipinfo/go-ipinfo/ipinfo"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"mvdan.cc/xurls/v2"
	"net"
	"net/http"
	"strings"
)

// next version maybe?
//func CloudDetectJS()  { }

func CloudDetectIP(domain string, key string) (string, error) {

	authTransport := ipinfo.AuthTransport{Token: key}
	httpClient := authTransport.Client()

	resp, err := net.LookupIP(domain)
	if err != nil {
		log.Error().Msg("Can't Connect to host")
		return "", err
	}
	if len(resp) < 1 {
		log.Error().Msg("Can't resolve target IP")
		return "", err
	}
	firstIp := resp[0]

	client := ipinfo.NewClient(httpClient)
	info, err := client.GetOrganization(firstIp)
	if err != nil {
		log.Error().Msg("Can't resolve target Organization")
	}

	log.Debug().Msg(strings.TrimSpace(info))

	info = strings.ToLower(info)

	if strings.Contains(info, "amazon") {
		return "amazon", nil
	}
	if strings.Contains(info, "linode") {
		return "linode", nil
	}
	if strings.Contains(info, "digitalocean") {
		return "digitalocean", nil
	}
	if strings.Contains(info, "google") {
		return "google", nil
	}

	if strings.Contains(info, "microsoft") {
		return "microsoft", nil
	}
	if strings.Contains(info, "alibaba") {
		return "alibaba", nil
	}
	if strings.Contains(info, "choopa") {
		return "vultr", nil
	}

	// CloudFlare detection if target is behind proxy it means we can't detect true provider
	if strings.Contains(info, "cloudflare") {
		return "", errors.New("CloudFlare detected target is behind proxy")
	}

	return info, nil

	// NO-API Idea
	//result, err := whois.Whois("0xsha.io")
	//if err == nil {
	//	res, err := whoisparser.Parse(result)
	//	fmt.Println(res.Administrative.Organization)
	//	if err!=nil{
	//		log.Fatal(err)
	//	}

}


func CloudDetectHTML(domain string, c *Config , providerPath string) (string, error) {

	if !strings.HasPrefix(domain, "https://") {

		domain = "https://" + domain
	}

	resp, err := http.Get(domain)
	if err != nil {
		return "", errors.New("can't connect to domain")
	}

	detected := ""

	for _, provider := range c.Providers {

		providerConfig, err := InitCloudConfig(provider, providerPath)

		var allUrls []string
		if err != nil {
			log.Fatal().Err(err).Msg("Exiting...")
		}

		// append all urls for each provider
		for _, itemURL := range providerConfig.APPUrls {
			allUrls = append(allUrls, itemURL)
		}
		for _, itemURL := range providerConfig.StorageUrls {
			allUrls = append(allUrls, itemURL)
		}
		for _, itemURL := range providerConfig.AppRegionUrls {
			allUrls = append(allUrls, itemURL)
		}
		for _, itemURL := range providerConfig.StorageRegionUrls {
			allUrls = append(allUrls, itemURL)
		}

		// make a dictionary
		providersUrl := make(map[string][]string)
		providersUrl[provider] = allUrls

		// let's make it relaxed and not miss any possibility
		rxRelaxed := xurls.Relaxed()
		rep, _ := ioutil.ReadAll(resp.Body)
		rx := rxRelaxed.FindAllString(string(rep), -1)

		for _, linkItem := range rx {

			for _, values := range providersUrl {

				for _, value := range values {

					if strings.Contains(linkItem, value) {
						detected = provider
						//log.Info().Msg(provider+" detected in source")
						break
					}
				}
			}
		}

		if detected != "" {
			return provider, nil
		}
	}
	return "", errors.New("no provider found in HTML")

}

func CheckSupportedCloud(org string, c *Config) (string, error) {

	for _, company := range c.Providers {

		if org == company {
			return org, nil
		}
	}
	return "", errors.New("unsupported cloud: ")

}
