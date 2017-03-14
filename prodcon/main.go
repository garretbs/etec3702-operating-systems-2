package main

//import "etec3702"
import "fmt"
import "os"
import "regexp"
import "net/url"
import "net/http"
import "io/ioutil"
import "sync"

var visited  = make(map[string]bool)
var output  = make(map[string]bool)

var activeThreads int
var sitesToVisit  = []string{originalDomain}

var globalsLock sync.Mutex
var condition = sync.NewCond(&globalsLock)

//READ ONLY
var rex = regexp.MustCompile(`(?i)<\s*a\s*href\s*=\s*('(\\'|[^'])*'|"(\\"|[^"])*"|[^ >]+)`)
var ignore = regexp.MustCompile(`(javascript:(.*))|(mailto:(.*))|(.*)((\.pdf)|(#)|(.pptx)|(.docx)|(.xlsx))`)

var originalDomain = os.Args[1:][0]
var originalHost,_ = url.Parse(originalDomain)

func main(){
	
	const numThreads = 4
	activeThreads = numThreads
	for i:=0;i<numThreads;i++{
		go spiderCrawl()
	}
	
	globalsLock.Lock()
	for activeThreads > 0{
		condition.Wait()
	}
	globalsLock.Unlock()
	
	fmt.Println("Done,", len(output), "sites")
}

func spiderCrawl(){

	var hyperlink string
	var currentSite string
	var ignored [][]string
	
	for(true){
		
		//Consume
		globalsLock.Lock()
		activeThreads--
		condition.Broadcast()
		for len(sitesToVisit) == 0{
			condition.Wait()
		}
		activeThreads++
		currentSite = sitesToVisit[0]
		sitesToVisit = sitesToVisit[1:]
		condition.Broadcast()
		globalsLock.Unlock()
		
		globalsLock.Lock()
		if !visited[currentSite]{
		
			visited[currentSite] = true
			globalsLock.Unlock()
			
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
						if(len(ignored) < 1){
							
							//Produce							
							u3,err := url.Parse(hyperlink)
							if err == nil{
								u3 = originalHost.ResolveReference(u3)
								hyperlink = u3.String()
								if u3.Host == originalHost.Host{
									globalsLock.Lock()
									sitesToVisit = append(sitesToVisit, hyperlink)
									condition.Broadcast()
									globalsLock.Unlock()
								}else{//if new link to external site
									globalsLock.Lock()
									if !output[hyperlink]{
										fmt.Println(hyperlink)
										output[hyperlink] = true
									}
									globalsLock.Unlock()
								}
							}
						}
					}
				}
			}
		}else{
			globalsLock.Unlock()
		}
	}
}
