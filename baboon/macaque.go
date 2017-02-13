package main
import "etec3702"

func macaque(){
    etec3702.Delay();
	macEntry.Acquire()
    etec3702.Output("Macaque on rope");
    etec3702.Delay();
    
	macLock.Acquire()
	numMacs++
	if numMacs == 3{
		macExit.Release() //let the monkeys run through		
	}else{
		macEntry.Release() //if not the third monkey, let another in
	}
	macLock.Release()
	
	macExit.Acquire()
	macExit.Release()
	etec3702.Output("Macaque off rope");
	
	macLock.Acquire()
	numMacs--
	if numMacs == 0{
		macExit.Acquire()
		babEntry.Release()
	}
	macLock.Release()
}
