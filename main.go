package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/evilsocket/brutemachine"
	"github.com/fatih/color"
	"regexp"
	"net/http"
	"io/ioutil"
	"unicode/utf8"
	"github.com/axgle/mahonia"
	"net"
	"time"
)

const Version = "1.0.3"

type Result struct {
	domain   string
	title    string
	Charset  string
	Language string
}

var (
	m         *brutemachine.Machine
	g         = color.New(color.FgGreen)
	r         = color.New(color.FgRed)
	domains   = flag.String("f", "domains.txt", "Input from list of domain list.")
	consumers = flag.Int("c", 30, "Number of concurrent consumers.")
	output    = flag.String("o", "", "Output results to file")
	timeout   = flag.Int64("t", 10, "http request timeout")
	pTitle    = regexp.MustCompile(`(?i:)<title>(.*?)</title>`)
)

func DoRequest(domain string) interface{} {
	result := Result{}
	if strings.Contains(domain, "http://") == false && strings.Contains(domain, "https://") == false {
		domain = "http://" + domain
	}
	client := &http.Client{Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			deadline := time.Now().Add(time.Duration(*timeout) * time.Second)
			c, err := net.DialTimeout(network, addr, time.Second*20)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
	},}
	result.domain = domain
	req, err := http.NewRequest("GET", domain, nil)
	req.Header.Set("Accept-Encoding", "")
	resp, err := client.Do(req)
	if err != nil {
		r.Println(result.domain, err.Error())
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Println(result.domain, err.Error())
		return nil
	}

	titleArr := pTitle.FindStringSubmatch(string(body))
	if titleArr != nil {
		if len(titleArr) == 2 {
			sTitle := titleArr[1]
			if !utf8.ValidString(sTitle) {
				sTitle = mahonia.NewDecoder("gb18030").ConvertString(sTitle)
			}
			result.title = sTitle
		} else {
			result.title = "无标题"
		}
	} else {
		result.title = "无标题"
	}

	return result
}

// OnResult prints out the results of a lookup
func OnResult(res interface{}) {
	result, ok := res.(Result)
	if !ok {
		r.Printf("Error while converting result.\n")
		return
	}
	g.Printf("%25s", result.domain)
	fmt.Println(" : ", result.title)
	fd, _ := os.OpenFile(*output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	s := strings.Join([]string{result.domain, "\t", result.title, "\n"}, "")
	buf := []byte(s)
	fd.Write(buf)
	fd.Close()
}

func main() {

	flag.Parse()

	if *domains == "" || *output == "" {
		flag.Usage()
		os.Exit(1)
	}

	m = brutemachine.New(*consumers, *domains, DoRequest, OnResult)
	if err := m.Start(); err != nil {
		panic(err)
	}

	m.Wait()

	g.Println("\nDONE")

	printStats()
}

// Print some stats
func printStats() {
	m.UpdateStats()
	fmt.Println("")
	fmt.Println("Requests :", m.Stats.Execs)
	fmt.Println("Results  :", m.Stats.Results)
	fmt.Println("Time     :", m.Stats.Total.Seconds(), "s")
	fmt.Println("Req/s    :", m.Stats.Eps)
}
