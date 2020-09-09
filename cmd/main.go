package main
// an early version of cloud brute plugin for HunterSuite.io
import (
	"fmt"
	engine "github.com/0xsha/cloudbrute/internal"
	"github.com/akamensky/argparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)




func main() {

	//banner
	banner := ` ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ██████╗ ██████╗ ██╗   ██╗████████╗███████╗
██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗██╔══██╗██╔══██╗██║   ██║╚══██╔══╝██╔════╝
██║     ██║     ██║   ██║██║   ██║██║  ██║██████╔╝██████╔╝██║   ██║   ██║   █████╗  
██║     ██║     ██║   ██║██║   ██║██║  ██║██╔══██╗██╔══██╗██║   ██║   ██║   ██╔══╝  
╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝██████╔╝██║  ██║╚██████╔╝   ██║   ███████╗
 ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝
						V 1.0.0`
	fmt.Println(banner)

	// beautify the results
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	parser := argparse.NewParser("CloudBrute", "Awesome Cloud Enumerator")

	domain := parser.String("d", "domain",
		&argparse.Options{
		Required: true,
		Help: "domain"})


	keyword := parser.String("k", "keyword",
		&argparse.Options{
			Required: true,
			Help: "keyword used to generator urls"})

	wordList := parser.String("w", "wordlist",
		&argparse.Options{
			Required: true,
			Help: "path to wordlist"})

	useProvider := parser.String("c", "cloud",
		&argparse.Options{
			Required: false,
			Help: "force a search, check config.yaml providers list"})

	threads := parser.Int("t", "threads",
		&argparse.Options{
			Required: false,
			Help: "number of threads",
			Default: 80})

	timeout := parser.Int("T", "timeout",
		&argparse.Options{
			Required: false,
			Help: "timeout per request in seconds",
			Default: 10})

	useProxy := parser.String("p", "proxy",
		&argparse.Options{
			Required: false,
			Help: "use proxy list"})


	useAgents := parser.String("a", "randomagent",
		&argparse.Options{
			Required: false,
			Help: "user agent randomization"})



	debug := parser.Flag("D", "debug",
		&argparse.Options{
			Required: false,
			Help: "show debug logs",
			Default: false})



	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	//initialize config.yaml
	config := engine.InitConfig("./config/config.yaml")
	apiKey := config.IPInfo
	log.Info().Msg("Initialized scan config")


	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}



	// check out args
	var details engine.RequestDetails

	if *useAgents != "" {
		userAgents , err :=  engine.ReadTextFile(*useAgents)

		if err!=nil{

			log.Fatal().Err(err).Msg("Can't read agents file, check config.yaml")
		}

		details.RandomAgent = userAgents
	}

	if *useProxy != "" {

		proxyList , err :=  engine.ReadTextFile(*useProxy)


		if err!=nil{

			log.Fatal().Err(err).Msg("Can't read proxy file , check config.yaml")
		}

		details.ProxyList = proxyList
		details.ProxyType = config.ProxyType
	}



	var cloud string
	if *useProvider != "" {

		cloud = *useProvider

		}else {

		// Detect the cloud
		cloud, err = engine.CloudDetect(*domain, apiKey)
		if err != nil {
			log.Fatal().Err(err).Msg("Exiting...")
		}

	}

	// Do we support the provider?
	provider, err := engine.CheckSupportedCloud(cloud, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Exiting...")
	}

	log.Info().Msg(provider + " detected")






	urls, err := engine.GenerateMutatedUrls(*wordList, provider, "./config/modules/", *keyword , config.Environments)
	//var p cb.Progress
	//p.CurrentProgress = 0
	//p.TotalProgress =  float64(len(urls))

	if err != nil {
		log.Fatal().Err(err).Msg("Exiting...")

	}


	output := engine.GenerateOutputName(*keyword)

	engine.AsyncHTTPHead(urls, *threads, *timeout , details , output)


}
