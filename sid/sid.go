package main

import "etec3702"
import "math/rand"
import "sync"


var numSearchers int
var numInserters int
var condLock sync.Mutex
var condition = sync.NewCond(&condLock)

func main(){
    counter:=0	
    for {
        etec3702.Delay()
        switch(rand.Intn(3)){
            case 0:
                go searcher(counter)
            case 1:
                go inserter(counter)
            case 2:
                go deleter(counter)
        }
        counter++
        etec3702.Delay()
    }
}

//any number of searchers
//can be concurrent with ONE inserter
func searcher(id int){
	condLock.Lock()
	numSearchers++
    etec3702.Delay()
    etec3702.Output("Searching")
	condition.Broadcast()
	condLock.Unlock()
	
	condLock.Lock()
	numSearchers--
    etec3702.Delay()	
    etec3702.Output("Done Searching")
	condition.Broadcast()
	condLock.Unlock()
}

//only ONE inserter
//can be concurrent with searchers
func inserter(id int){
	condLock.Lock()
	for numInserters > 0{
		condition.Wait()
	}
	numInserters++
    etec3702.Delay()	
    etec3702.Output("Inserting")
	condition.Broadcast()
	condLock.Unlock()
	
	condLock.Lock()
	numInserters--
    etec3702.Delay()
    etec3702.Output("Done Inserting")
	condition.Broadcast()
	condLock.Unlock()
}

//only ONE at a time
//NO concurrency. forever alone
func deleter(id int){
	condLock.Lock()
	for numInserters > 0 || numSearchers > 0{
		condition.Wait()
	}
    etec3702.Delay()
    etec3702.Output("Deleting")
    etec3702.Delay()
	
    etec3702.Output("Done Deleting")
	condition.Broadcast()
	condLock.Unlock()
}
