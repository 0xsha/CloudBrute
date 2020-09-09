package internal

import (
	"errors"
	"github.com/ipinfo/go-ipinfo/ipinfo"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
)

func CloudDetect(domain string , key string) (string , error) {

	authTransport := ipinfo.AuthTransport{Token: key}
	httpClient := authTransport.Client()


	resp, err := net.LookupIP(domain)
	if err != nil {
		log.Error().Msg("Can't Connect to host")
		return "",err
	}
	if len(resp) < 1{
		log.Error().Msg("Can't resolve target IP")
		return "",err
	}
	firstIp := resp[0]

	client := ipinfo.NewClient(httpClient)
	info, err := client.GetOrganization(firstIp)
	if err != nil {
		log.Error().Msg("Can't resolve target Organization")
	}

	log.Debug().Msg( strings.TrimSpace( info))

	info = strings.ToLower(info)



	if strings.Contains(info,"amazon") {
		return "amazon",nil
	}
	if strings.Contains(info,"linode") {
		return "linode",nil
	}
	if strings.Contains(info,"digitalocean") {
		return "digitalocean",nil
	}
	if strings.Contains(info,"google") {
		return "google",nil
	}

	if strings.Contains(info,"microsoft") {
		return "microsoft",nil
	}
	if strings.Contains(info,"alibaba") {
		return "alibaba",nil
	}
	if strings.Contains(info,"choopa") {
		return "vultr",nil
	}

	// CloudFlare detection if target is behind proxy it means we can't detect true provider
	if strings.Contains(info,"cloudflare") {
		return "",errors.New("CloudFlare detected target is behind proxy")
	}

	return info,nil


	// NO-API Idea
	//result, err := whois.Whois("0xsha.io")
	//if err == nil {
	//	res, err := whoisparser.Parse(result)
	//	fmt.Println(res.Administrative.Organization)
	//	if err!=nil{
	//		log.Fatal(err)
	//	}

	}

func CheckSupportedCloud(org string, c *Config) (string,error) {

	for _,company := range c.Providers{

		if org==company {
			return org,nil
		}
	}
	return "",errors.New("unsupported cloud: ")

}