package main
import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"golang.org/x/net/dict"
	"github.com/MakeNowJust/heredoc/v2"
)

func main() {
	args := os.Args[1:]
	l := len(args)
	if l < 1 {
		os.Exit(showHelp())
	}
	if args[0] == "help" {
		os.Exit(showHelp())
	}
	if args[0] == "weather" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showWeather(args[1:]...))
	} else if args[0] == "subnet" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showSubnet(args[1]))
	} else if args[0] == "asn" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showASN(args[1]))
	} else if args[0] == "upordown" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showUpOrDown(args[1]))
	} else if args[0] == "scan" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showScan(args[1]))
	} else if args[0] == "qrcode" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showQRCode(args[1]))
	} else if args[0] == "cheat" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showCheat(args[1]))
	} else if args[0] == "dict" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showDict(args[1]))
	} else if args[0] == "locate" {
		if l < 2 {
			os.Exit(showHelp())
		}
		os.Exit(showLocate(args[1]))
	} else if args[0] == "myip" {
		os.Exit(showMyIP())
	}
}

func showHelp() int {
	fmt.Println(heredoc.Doc(`
Quick help (arguments):
help
    Display this help text.
weather <location> [<format>]
    Display current weather in location (e.g. 'Los_Angeles') and format='1' to '3'
subnet <x.x.x.x/y>
    Subnet calculator.
asn <x.x.x.x>
    Find AS Number for a given IP address.
upordown <http[s]://example.com>
    Is a host up?
scan <x.x.x.x>
    Scan a host's open ports
locate <x.x.x.x>
    Geo locate an IP address
myip
    My IP: What is it?
qrcode <string>
    Generate a QR code
cheat <string>
    Documentation on a command/concept
dict <string>
    Word/phrase definition
`))
	return 0
}

func doQuery(queryVerb string, queryString string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(queryVerb, queryString, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "curl/7.79.1")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func output(out string) {
	fmt.Println(out)
}

func showWeather(args ...string) int {
	var queryString string
	if len(args) > 1 {
		queryString = fmt.Sprintf("https://wttr.in/%s?u&format=%s", args[0], args[1])
	} else {
		queryString = fmt.Sprintf("https://wttr.in/%s?u", args[0])
	}
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

func showSubnet(arg string) int {
	queryString := fmt.Sprintf("https://api.hackertarget.com/subnetcalc/?q=%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

func showASN(arg string) int {
	queryString := fmt.Sprintf("https://api.hackertarget.com/aslookup/?q=%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

// TODO configure email address
func showUpOrDown(arg string) int {
	prefix := ""
	if ! strings.HasPrefix(arg, "http://") && ! strings.HasPrefix(arg, "https://") {
		prefix = "https://"
	}
        email := "you@example.com"
	queryString := fmt.Sprintf("https://ping.gl/%s/%s%s", email, prefix, arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

func showScan(arg string) int {
	queryString := fmt.Sprintf("https://api.hackertarget.com/nmap/?q=%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

// TODO accept file/stdin
func showQRCode(arg string) int {
	queryString := fmt.Sprintf("qrenco.de/%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

func showCheat(arg string) int {
	queryString := fmt.Sprintf("https://cheat.sh/%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

// TODO Support multiple dictionaries, e.g. gcide, etc.
func showDict(arg string) int {
	client, err := dict.Dial("tcp", "dict.org:2628")
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer client.Close()
	defn, err := client.Define("english", arg)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	for i := range defn {
		fmt.Println(string(defn[i].Text))
	}
/*
	dicts, err := client.Dicts()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	for i := range dicts {
		fmt.Println(dicts[i])
	}
*/
	return 0
}

func showLocate(arg string) int {
	queryString := fmt.Sprintf("http://ip-api.com/%s", arg)
	out, err := doQuery("GET", queryString)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}

func showMyIP() int {
	out, err := doQuery("GET", "https://ident.me")
	if err != nil {
		fmt.Println(err)
		return -1
	}
	output(out)
	return 0
}
