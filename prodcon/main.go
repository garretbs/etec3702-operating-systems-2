package main

import "etec3702"
import "fmt"
import "os"
import "regexp"
import "net/url"
import "net/http"
import "io/ioutil"
import "sync"

//var rex regexp
var sitesLock sync.Mutex
var visited  = make(map[string]bool)
var output  = make(map[string]bool)

var originalDomain = os.Args[1:][0]
var originalHost,_ = url.Parse(originalDomain)

var sitesToVisit  = []string{originalDomain}
var consume = etec3702.NewSemaphore(1)

var rex = regexp.MustCompile(`(?i)<\s*a\s*href\s*=\s*('(\\'|[^'])*'|"(\\"|[^"])*"|[^ >]+)`)
var ignore = regexp.MustCompile(`(mailto:(.*))|(.*)((\.pdf)|(#)|(.pptx)|(.docx)|(.xlsx))`)

//BE SURE to use locks when accessing global variables, ya dumb fuck
//for some reason, www.shawnee.edu is not the same hostname as shawnee.edu. wtf
//ditto with things like myssu.shawnee.edu

//test site:
// selenium.ssucet.org:8001
// ~50 links


func main(){
	
	const numThreads = 4
	var chans [numThreads]chan bool
	
	for i:=0;i<numThreads;i++{
		chans[i] = make(chan bool)
		go spiderCrawl(chans[i])
	}
	
	
	for i:=0;i<numThreads;i++{
		_ = <- chans[i]
	}
}

func spiderCrawl(c chan bool){

	var hyperlink string
	var currentSite string
	var ignored [][]string
	
	for(true){
		
		//Consume
		consume.Acquire()
		sitesLock.Lock()
		
		currentSite = sitesToVisit[0] //om nom nom
		sitesToVisit = sitesToVisit[1:] //clean my plate
		sitesLock.Unlock()
		
		//fmt.Println(currentSite)
		
		sitesLock.Lock()
		if !visited[currentSite]{
		
			visited[currentSite] = true
			sitesLock.Unlock()			
			
			currentUrl,_ := url.Parse(currentSite)
		
			resp,err := http.Get(currentUrl.String())
			if err == nil{
				data,err := ioutil.ReadAll(resp.Body)
				if err == nil{
					resp.Body.Close()
					
					s := string(data)
					sm := rex.FindAllStringSubmatch(s, -1)
					
					for i:=0;i<len(sm);i++{
						hyperlink = sm[i][1]
						hyperlink = hyperlink[1:len(hyperlink)-1] //remove first and last characters, which are quotation marks
						
						
						ignored = ignore.FindAllStringSubmatch(hyperlink, -1)
						if(len(ignored) < 1 && len(hyperlink) > 1){
							
							if hyperlink[0] == '/' {
								if currentSite[len(currentSite)-1] == '/'{
									currentSite = currentSite[:len(currentSite)-1]
								}//so we don't get double slashes
								hyperlink = originalDomain + hyperlink
							}							
							
							//Produce
							u3,err := url.Parse(hyperlink)
							if err == nil{
								if u3.Host == originalHost.Host{
									sitesLock.Lock()
									sitesToVisit = append(sitesToVisit, hyperlink)
									consume.Release()
									sitesLock.Unlock()
									//fmt.Println(hyperlink)
								}else{//if new link to external site
									sitesLock.Lock()
									if !output[hyperlink]{
										//fmt.Println(u3.Host, originalHost.Host)
										fmt.Println(hyperlink)
										//fmt.Println(len(output))
										output[hyperlink] = true
									}
									sitesLock.Unlock()
								}
							}
						}
					}
				}
			}
		}else{
			sitesLock.Unlock()
		}
	}
}
