package main

// an early version of cloud brute plugin for HunterSuite.io
import (
	"fmt"
	engine "github.com/0xsha/cloudbrute/internal"
	"github.com/akamensky/argparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path"
)

func main() {


	// parse arguments
	parser := argparse.NewParser("CloudBrute", "Awesome Cloud Enumerator")
	domain := parser.String("d", "domain",
		&argparse.Options{
			Required: true,
			Help:     "domain"})

	keyword := parser.String("k", "keyword",
		&argparse.Options{
			Required: true,
			Help:     "keyword used to generator urls"})

	wordList := parser.String("w", "wordlist",
		&argparse.Options{
			Required: true,
			Help:     "path to wordlist"})

	useProvider := parser.String("c", "cloud",
		&argparse.Options{
			Required: false,
			Help:     "force a search, check config.yaml providers list"})

	threads := parser.Int("t", "threads",
		&argparse.Options{
			Required: false,
			Help:     "number of threads",
			Default:  80})

	timeout := parser.Int("T", "timeout",
		&argparse.Options{
			Required: false,
			Help:     "timeout per request in seconds",
			Default:  10})

	useProxy := parser.String("p", "proxy",
		&argparse.Options{
			Required: false,
			Help:     "use proxy list"})

	useAgents := parser.String("a", "randomagent",
		&argparse.Options{
			Required: false,
			Help:     "user agent randomization"})

	debug := parser.Flag("D", "debug",
		&argparse.Options{
			Required: false,
			Help:     "show debug logs",
			Default:  false})
	quite := parser.Flag("q", "quite",
		&argparse.Options{
			Required: false,
			Help:     "suppress all output",
			Default:  false})

	mode := parser.String("m", "mode",
		&argparse.Options{
			Required: false,
			Default: "storage",
			Help:     "storage or app"})

	output := parser.String("o", "output",
		&argparse.Options{
			Default:  "out.txt",
			Required: false,
			Help:     "Output file"})

	configFolder := parser.String("C", "configFolder",
		&argparse.Options{
			Default:  "config",
			Required: false,
			Help:     "Config path"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	//banner
	banner := ` ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ██████╗ ██████╗ ██╗   ██╗████████╗███████╗
██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗██╔══██╗██╔══██╗██║   ██║╚══██╔══╝██╔════╝
██║     ██║     ██║   ██║██║   ██║██║  ██║██████╔╝██████╔╝██║   ██║   ██║   █████╗  
██║     ██║     ██║   ██║██║   ██║██║  ██║██╔══██╗██╔══██╗██║   ██║   ██║   ██╔══╝  
╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝██████╔╝██║  ██║╚██████╔╝   ██║   ███████╗
 ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝
						V 1.0.7
`
	if !*quite {
		_, _ = fmt.Fprintf(os.Stderr, banner)
	}

	// beautify the results
	if !*quite {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: ioutil.Discard})
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

	// check mode

	if *mode != "storage" && *mode != "app" {

		log.Fatal().Msg("Invalid mode use app or storage")
	}

	// initialize config.yaml
	configPath := path.Join(*configFolder, "config.yaml")
	providerPath := path.Join(*configFolder, "modules")
	log.Info().Msg(fmt.Sprintf("Detect config path: %s", configPath))
	log.Info().Msg(fmt.Sprintf("Detect provider path: %s", providerPath))
	config := engine.InitConfig(configPath)
	apiKey := config.IPInfo
	log.Info().Msg("Initialized scan config")

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// check out args
	var details engine.RequestDetails

	if *useAgents != "" {
		userAgents, err := engine.ReadTextFile(*useAgents)
		if err != nil {
			log.Fatal().Err(err).Msg("Can't read agents file, check config.yaml")
		}
		details.RandomAgent = userAgents
	}

	if *useProxy != "" {

		proxyList, err := engine.ReadTextFile(*useProxy)

		if err != nil {

			log.Fatal().Err(err).Msg("Can't read proxy file , check config.yaml")
		}

		details.ProxyList = proxyList
		details.ProxyType = config.ProxyType
	}

	var cloud string
	if *useProvider != "" {

		cloud = *useProvider

	} else {

		// Detect the cloud using ip info
		cloud, err = engine.CloudDetectIP(*domain, apiKey)
		if err != nil {
			log.Error().Err(err).Msg("IP detection failed")
		}


	}

	// Do we support the provider?
	provider, err := engine.CheckSupportedCloud(cloud, config)
	if err != nil {
		log.Warn().Msg("IP detection failed")

		// Detect the cloud from HTML and JavaScript source codes
		provider, err = engine.CloudDetectHTML(*domain, config , providerPath)

		if err != nil {

			log.Fatal().Err(err).Msg("Source detection failed as well use -c .")
		}

	}


	log.Info().Msg(provider + " detected")

	urls, err := engine.GenerateMutatedUrls(*wordList, *mode, provider, providerPath, *keyword, config.Environments)

	if err != nil {
		log.Fatal().Err(err).Msg("Exiting...")

	}

	//output := engine.GenerateOutputName(*keyword)
	engine.AsyncHTTPHead(urls, *threads, *timeout, details, *output)
}
