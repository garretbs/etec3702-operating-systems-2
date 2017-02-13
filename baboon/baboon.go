package main
import "etec3702"

func baboon(){
    etec3702.Delay();
	babEntry.Acquire()
	etec3702.Output("Baboon on rope");
    etec3702.Delay();
	
	babLock.Acquire()
	numBabs++
	if numBabs == 3{
		babExit.Release() //let the monkeys run through
	}else{
		babEntry.Release() //if not the third monkey, let another in
	}
	babLock.Release()
	
	babExit.Acquire()
	babExit.Release()
	etec3702.Output("Baboon off rope");
	
	babLock.Acquire()
	numBabs--
	if numBabs == 0{
		babExit.Acquire()
		macEntry.Release()
	}
	babLock.Release()
}
