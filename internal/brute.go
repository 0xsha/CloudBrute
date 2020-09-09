package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func HandleHTTPRequests(reqs, results chan string, quit chan int, bar *pb.ProgressBar , details *RequestDetails) {


	for link := range reqs {


		log.Debug().Msg(link)
		if len(details.ProxyList) > 0 {

		  	chosenProxy := SelectRandomItem(details.ProxyList)

		  	if details.ProxyType == "socks5" {


				log.Debug().Msg("requesting through socks5 proxy : " + chosenProxy)


				dialSocksProxy, err := proxy.SOCKS5("tcp", chosenProxy, nil, proxy.Direct)
				socksTransport := &http.Transport{Dial: dialSocksProxy.Dial, DisableKeepAlives: true , TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

				if err != nil {
					continue
				}

				socksClient := &http.Client{
					Transport: socksTransport,
				}

				req , _ := http.NewRequest("HEAD","https://" +link ,nil)

				if len(details.RandomAgent) >0 {

					chosenAgent := SelectRandomItem(details.RandomAgent)
					req.Header.Set("User-Agent",chosenAgent)

				}

				resp, err := socksClient.Do(req)

				if err != nil {

					log.Err(err).Msg("err")

					results <- "err"
					bar.Increment()
					continue
				}
				bar.Increment()
				results <- link + ":" + strconv.Itoa(resp.StatusCode)

			}

			if details.ProxyType == "http"{


				proxyURL, _ := url.Parse("http://"+chosenProxy)

				log.Debug().Msg("requesting through http proxy : " + chosenProxy)

				httpProxyClient := &http.Client{
					Transport: &http.Transport{
						DisableKeepAlives: true,
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
						Proxy:             http.ProxyURL(proxyURL),
					},
				}



				req , _ := http.NewRequest("HEAD","http://" +link ,nil)

				if len(details.RandomAgent) >1 {

					chosenAgent := SelectRandomItem(details.RandomAgent)
					req.Header.Set("User-Agent",chosenAgent)
					log.Debug().Msg("user-agent : " + chosenAgent)

				}

				resp, err := httpProxyClient.Do(req)

				if err != nil {

					log.Err(err).Msg("proxy error")

					results <- "err"
					bar.Increment()
					continue
				}
				bar.Increment()
				results <- link + ":" + strconv.Itoa(resp.StatusCode)


			}

		} else {

			client := http.Client{
				Transport: &http.Transport{
					DisableKeepAlives: true },
			 }

			req , _ := http.NewRequest("HEAD","https://" +link ,nil)

			if len(details.RandomAgent) >0 {

				chosenAgent := SelectRandomItem(details.RandomAgent)
				req.Header.Set("User-Agent",chosenAgent)
			}

			resp, err := client.Do(req)


			if err != nil {

				results <- "err"
				bar.Increment()
				continue
			}

			//log.Debug().Msg(strconv.Itoa(resp.StatusCode))

			bar.Increment()
			results <- link + ":" + strconv.Itoa(resp.StatusCode)
		}


		if len(reqs) == len(results) {
			quit <- 0
		}


	}


}

func AsyncHTTPHead(urls []string, threads int, timeout int , details RequestDetails) {

	result := make(chan string)
	reqs := make(chan string, len(urls)) // buffered
	quit := make(chan int)

	bar := pb.StartNew(len(urls))

	for i := 0; i < threads; i++ {
		go HandleHTTPRequests(reqs, result, quit , bar , &details)
	}

	go func() {
		for _, link := range urls {
			reqs <- link
		}
	}()




	// parsing http codes
	// 500 , 502 server error
	// 404 not found
	// 200 found
	// 400, 401 , 403  protected
	// 302 , 301 redirect

	for {
		select {
		case res := <-result:
			if res != "err" {
				if strings.Contains(res, "200") {
					log.Info().Msg("Open : " + "[response code :" + res +"]")
				}
				if strings.Contains(res, "301") || strings.Contains(res, "402") {
					log.Warn().Msg("Redirect : " + "[response code :" + res +"]")
				}
				if strings.Contains(res, "400") || strings.Contains(res, "401") ||
					strings.Contains(res, "403") {
					log.Warn().Msg("Protected : " + "[response code :" + res +"]")
				}
				if strings.Contains(res, "500") || strings.Contains(res, "502") {
					log.Warn().Msg("Server error :" + "[response code :" + res +"]")
				}
			}

		case <-time.After(time.Duration(timeout) * time.Second):
			fmt.Println("timeout")
			bar.Increment()
		case <-quit:
			bar.Set(len(urls))
			bar.Finish()

			return
		}
	}



}

func GenerateMutatedUrls(wordListPath string, provider string, providerPath string, target string , environments []string) ([]string, error) {


	//envs := []string{"test", "dev", "prod", "stage"}
	words ,err := ReadTextFile(wordListPath)

	if err!=nil{
		log.Fatal().Err(err).Msg("Exiting ...")
	}
	permutations := []string{"%s-%s-%s", "%s-%s.%s", "%s-%s%s", "%s.%s-%s", "%s.%s.%s"}

	var compiled []string

	for _, env := range environments {

		for _, word := range words {

			for _, permutation := range permutations {
				formatted := fmt.Sprintf(permutation, target, word, env)
				compiled = append(compiled, formatted)
			}

		}
	}

	urlPermutations := []string{"%s.%s", "%s-%s", "%s%s"}
	for _, word := range words {

		for _, permutation := range urlPermutations {
			formatted := fmt.Sprintf(permutation, target, word)
			compiled = append(compiled, formatted)
		}

	}

	providerConfig, err := InitCloudConfig(provider ,providerPath )

	if err != nil {

		log.Fatal().Err(err).Msg("Exiting...")
	}

	log.Info().Msg("Initialized " + provider + " config")

	var finalUrls []string

	if len(providerConfig.Regions) > 0 {

		for _, region := range providerConfig.Regions {

			for _, regionUrl := range providerConfig.RegionUrls {

				for _, word := range compiled {

					finalUrls =  append ( finalUrls, word + "." + region + "." + regionUrl)
				}
			}
		}
	}

	if len(providerConfig.StorageUrls) > 0 {
		for _, storage := range providerConfig.StorageUrls {

			for _, word := range compiled {

				finalUrls = append(finalUrls, word+"."+storage)

			}

		}
	}

	if len(providerConfig.APPUrls) > 0 {
		for _, app := range providerConfig.APPUrls {

			for _, word := range compiled {

				finalUrls =  append ( finalUrls, word + "." + app)
			}

		}
	}

	return finalUrls, nil

}
