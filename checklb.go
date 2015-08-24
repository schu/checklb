package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
)

var Version string

func main() {
	usage := `checklb - send a HTTP request to a number of targets

<target> can be an IP address or a domain name. The latter will be
resolved to send a request to each IP. If no target is given,
<host-header> will be resolved as target.

Usage:
  checklb [options] <host-header> [<target>...]

Options:
  --https        Send https request (default: http)
  --path=<path>  Path to request [default: /]
`

	args, _ := docopt.Parse(usage, nil, true, fmt.Sprintf("checklb %s", Version), false)

	formatError := color.New(color.FgRed).SprintfFunc()
	formatWarning := color.New(color.FgYellow).SprintfFunc()

	host := args["<host-header>"].(string)
	path := args["--path"].(string)

	proto := "http"
	if args["--https"].(bool) {
		proto = "https"
	}

	targets := args["<target>"].([]string)
	if len(targets) == 0 {
		targets = []string{host}
	}

	var resolvedTargets []net.IP
	var wg sync.WaitGroup

	wg.Add(len(targets))

	mutex := &sync.Mutex{}
	for _, target := range targets {
		go func(target string) {
			defer wg.Done()
			addr := net.ParseIP(target)
			if addr != nil {
				mutex.Lock()
				resolvedTargets = append(resolvedTargets, addr)
				mutex.Unlock()
				return
			}
			addrs, err := net.LookupIP(target)
			if err != nil {
				fmt.Fprint(os.Stderr, formatError("Failed to resolve %s\n%s\n", target, err))
				os.Exit(1)
			}
			mutex.Lock()
			resolvedTargets = append(resolvedTargets, addrs...)
			mutex.Unlock()
		}(target)
	}

	wg.Wait()
	wg.Add(len(resolvedTargets))

	type result struct {
		Response *http.Response
		Target   net.IP
	}

	results := make(chan *result)

	for _, target := range resolvedTargets {
		go func(target net.IP) {
			defer wg.Done()
			transport := &http.Transport{}
			if proto == "https" {
				transport.TLSClientConfig = &tls.Config{ServerName: host}
			}
			client := &http.Client{Transport: transport}
			req, err := http.NewRequest("GET", fmt.Sprintf("%s://[%s]%s", proto, target.String(), path), nil)
			if err != nil {
				fmt.Fprint(os.Stderr, formatError("Failed to prepare HTTP request\n%s\n", err))
				os.Exit(1)
			}
			req.Host = host
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprint(os.Stderr, formatError("HTTP request failed\n%s\n", err))
				os.Exit(1)
			}
			defer resp.Body.Close()
			results <- &result{resp, target}
		}(target)
	}

	go func() {
		for res := range results {
			output := fmt.Sprintf("%s\t%s\n", res.Target, res.Response.Status)
			if res.Response.StatusCode != 200 {
				fmt.Print(formatWarning(output))
			} else {
				fmt.Print(output)
			}
		}
	}()

	wg.Wait()
}
