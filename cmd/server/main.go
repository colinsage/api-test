package main

import (
	"fmt"
	"github.com/colinsage/api-test/model"
	"github.com/json-iterator/go"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main(){
	fmt.Println(os.Args)

	project := loadProject(os.Args[1])
	plan := loadPlan(os.Args[2])

	fmt.Println("project", project)
	fmt.Println("plan", plan)

	plan.Merge(project)

	executePlan(plan)
}


func loadProject(path string) *model.Project{
	project := model.Project{}

	f, _ := os.Open(path)
	data,_ := ioutil.ReadAll(f)


	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &project)
	if err != nil {
		fmt.Println("decode failed. ", err)
	}
	return &project
}

func loadPlan(path string) *model.Plan{
	plan := model.Plan{
		CurrentQps: make(map[string]int),
	}

	f, _ := os.Open(path)
	data,_ := ioutil.ReadAll(f)


	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &plan)
	if err != nil {
		fmt.Println("decode failed. ", err)
	}
	return &plan
}

func executePlan(plan *model.Plan){
	var wg sync.WaitGroup
	for _, link := range plan.LinkList {
		wg.Add(1)
		go func(){
			runLink(&link, plan.CurrentQps[link.Name])
			wg.Done()
		}()
	}
	wg.Wait()
}

func runLink(link *model.Link, qps int){
	//load query
	f, _ := os.Open(link.Query)
	data,_ := ioutil.ReadAll(f)

	qs := strings.Split(string(data), "\n")

	limiter := rate.NewLimiter(rate.Limit(qps), 2*qps)

	failed := 0
	success := 0

	omit := 0
	var last int64
	lastFailed := 0
	lastSuccess := 0
	lastOmit := 0
	last = time.Now().Unix()

	//
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	tcpAddress := link.Service.Address +
		":" + strconv.Itoa(link.Service.Port);
	conn, err := dial.Dial("tcp", tcpAddress )

	if err != nil {
		fmt.Printf("conn failed %s \n", tcpAddress )
		return
	}


	for {
		max := len(qs)
		for i:=0; i< max; i++{
			line := qs[i]
			if len(line) == 0 {
				continue
			}
			if limiter.Allow() {
				p := strings.Split(line, "?")
				if len(p) != 3 {
					fmt.Println("wrong query ", i,  line, len(line))
					time.Sleep(time.Microsecond*1)
					continue
				}
				url := "http://" + link.Service.Address +
					":" + strconv.Itoa(link.Service.Port) +
					p[0] + "?" + p[1]

				content := p[2]

				go func(u string, c string){
					resp, err := http.Post(u, "", strings.NewReader(c))
					if err != nil || resp.StatusCode != 200 {
						if resp != nil {
							fmt.Println("post failed", resp.StatusCode, err, c)
						}else{
							fmt.Println("post failed", err, c)

						}
						failed++
					}else{
						success++
					}
				}(url, content)
			}else{
				omit++
				time.Sleep(time.Microsecond*1)
				i--
			}

			now := time.Now().Unix()

			if (now - last) > 10{
				fmt.Printf("succ: %v, faild: %v, omit: %d \n", (success-lastSuccess)/10.0, (failed-lastFailed)/10.0 , (omit - lastOmit)/10)
				last = now
				lastSuccess = success
				lastFailed = failed
				lastOmit = omit
			}
		}
	}
	// send query

}