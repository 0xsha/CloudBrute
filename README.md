# CloudBrute


A tool to find a company (target) infrastructure, files, and apps on the top cloud providers (Amazon, Google, Microsoft, DigitalOcean, Alibaba, Vultr, Linode). 
The outcome is useful for bug bounty hunters, red teamers, and penetration testers alike.  

##### The complete writeup is available [here](https://0xsha.io/posts/introducing-cloudbrute-wild-hunt-on-the-clouds)



## Cloud brute on the clouds?

<p align="center">
  <img alt="huntersuite" src="https://user-images.githubusercontent.com/23289085/101143253-35ea6b80-3649-11eb-9130-d1fc306c9a76.png" height="200" />
  <p align="center">
Enjoying this tool? Support it's development and take your game to the next level by using <a href="https://huntersuite.io">HunterSuite.io</a>
  </p>
</p>



## At a glance 

![CloudBrute](./assets/cloudbrute_digram.png)



## Motivation 

While working on [HunterSuite](https://huntersuite.io), and as part of the job, we are always thinking of something we can automate to make black-box security testing easier. We discussed this idea of creating a multiple platform cloud brute-force hunter.mainly to find open buckets, apps, and databases hosted on the clouds and possibly app behind proxy servers.   
Here is the list issues we tried to fix:

- separated wordlists 
- lack of proper concurrency 
- lack of supporting all major cloud providers 
- require authentication or keys or cloud CLI access
- outdated endpoints and regions 
- Incorrect file storage detection 
- lack support for proxies (useful for bypassing region restrictions) 
- lack support for user agent randomization (useful for bypassing rare restrictions) 
- hard to use, poorly configured

## Features
- Cloud detection (IPINFO API and Source Code)
- Supports all major providers
- Black-Box (unauthenticated)
- Fast (concurrent)
- Modular and easily customizable 
- Cross Platform (windows, linux, mac)
- User-Agent Randomization 
- Proxy Randomization (HTTP, Socks5) 

## Supported Cloud Providers

Microsoft:
- Storage
- Apps

Amazon: 
- Storage
- Apps

Google: 
- Storage
- Apps 

DigitalOcean: 
- storage

Vultr:
- Storage 

Linode:
- Storage

Alibaba:
- Storage 

## Version
1.0.0


## Usage
Just download the latest [release](https://github.com/0xsha/CloudBrute/releases) for your operation system and follow the usage.

To make the best use of this tool, you have to understand how to configure it correctly. When you open your downloaded version, there is a config folder, and there is a config.YAML file in there.

It looks like this 
```yaml
providers: ["amazon","alibaba","amazon","microsoft","digitalocean","linode","vultr","google"] # supported providers
environments: [ "test", "dev", "prod", "stage" , "staging" , "bak" ] # used for mutations
proxytype: "http"  # socks5 / http
ipinfo: ""      # IPINFO.io API KEY
```

For IPINFO API, you can register and get a free key at [IPINFO](https://ipinfo.io), the environments used to generate URLs, such as test-keyword.target.region and test.keyword.target.region, etc.

We provided some wordlist out of the box, but it's better to customize and minimize your wordlists (based on your recon) before executing the tool.

After setting up your API key, you are ready to use CloudBrute. 


```
 ██████╗██╗      ██████╗ ██╗   ██╗██████╗ ██████╗ ██████╗ ██╗   ██╗████████╗███████╗
██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗██╔══██╗██╔══██╗██║   ██║╚══██╔══╝██╔════╝
██║     ██║     ██║   ██║██║   ██║██║  ██║██████╔╝██████╔╝██║   ██║   ██║   █████╗  
██║     ██║     ██║   ██║██║   ██║██║  ██║██╔══██╗██╔══██╗██║   ██║   ██║   ██╔══╝  
╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝██████╔╝██║  ██║╚██████╔╝   ██║   ███████╗
 ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚══════╝
                                                V 1.0.7
usage: CloudBrute [-h|--help] -d|--domain "<value>" -k|--keyword "<value>"
                  -w|--wordlist "<value>" [-c|--cloud "<value>"] [-t|--threads
                  <integer>] [-T|--timeout <integer>] [-p|--proxy "<value>"]
                  [-a|--randomagent "<value>"] [-D|--debug] [-q|--quite]
                  [-m|--mode "<value>"] [-o|--output "<value>"]
                  [-C|--configFolder "<value>"]

                  Awesome Cloud Enumerator

Arguments:

  -h  --help          Print help information
  -d  --domain        domain
  -k  --keyword       keyword used to generator urls
  -w  --wordlist      path to wordlist
  -c  --cloud         force a search, check config.yaml providers list
  -t  --threads       number of threads. Default: 80
  -T  --timeout       timeout per request in seconds. Default: 10
  -p  --proxy         use proxy list
  -a  --randomagent   user agent randomization
  -D  --debug         show debug logs. Default: false
  -q  --quite         suppress all output. Default: false
  -m  --mode          storage or app. Default: storage
  -o  --output        Output file. Default: out.txt
  -C  --configFolder  Config path. Default: config


```

for example 
```
CloudBrute -d target.com -k target -m storage -t 80 -T 10 -w "./data/storage_small.txt"
```
please note -k keyword used to generate URLs, so if you want the full domain to be part of mutation, you have used it for both domain (-d) and keyword (-k) arguments 

If a cloud provider not detected or want force searching on a specific provider, you can use -c option.
```
CloudBrute -d target.com -k keyword -m storage -t 80 -T 10 -w -c amazon -o target_output.txt
```


## Dev 
- Clone the repo 
- go build -o CloudBrute main.go
- go test internal 



## in action

[![asciicast](https://asciinema.org/a/QIYRNgJMKhGX3woUTB3kh0HmC.svg)](https://asciinema.org/a/QIYRNgJMKhGX3woUTB3kh0HmC)

##  How to contribute
- Add a module or fix something and then pull request.
- Share it with whomever you believe can use it.
- Do the extra work and share your findings with community &hearts;


## FAQ

##### How to make the best out of this tool? 
Read the usage.

##### I get errors; what should I do? 
Make sure you read the usage correctly, and if you think you found a bug open an issue. 

##### When I use proxies, I get too many errors, or it's too slow?
It's because you use public proxies, use private and higher quality proxies. You can use [ProxyFor](https://github.com/0xsha/proxyfor) to verify the good proxies with your chosen provider. 

##### too fast or too slow ?
change -T (timeout) option to get best results for your run.

## Credits 

Inspired by every single repo listed  [here](https://github.com/mxm0z/awesome-sec-s3)

