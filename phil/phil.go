package main

import "etec3702"
import "sync"
import "runtime"

var footman = etec3702.NewSemaphore(4)
var S [5]sync.Mutex

func left(num int) int{
    if num == 0 {
        return 4
    } else {
        return num-1
    }
}
func right(num int) int {
    return (num+1)%5
}

func phil(num int){
    for {
        etec3702.Output("Philosopher",num,"is thinking")
        etec3702.Delay()
        footman.Acquire()
        S[left(num)].Lock()
        S[right(num)].Lock()
        etec3702.Output("Philosopher",num,"is eating")
        etec3702.Delay()
        S[left(num)].Unlock()
        S[right(num)].Unlock()
        footman.Release()
    }
}

func main(){
    go phil(0)
    go phil(1)
    go phil(2)
    go phil(3)
    go phil(4)
    runtime.Goexit()    //main thread can exit now; children stay alive
}
