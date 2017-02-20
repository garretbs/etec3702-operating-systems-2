package main
import "etec3702"

func macaque(){

	monkLock.Lock()
	for numBabs > 0 || numMacs >= 3 || totalMacs >= 49{ //don't go while max macaques or any baboons
		condition.Wait()
	}
	numMacs+=1
	totalMacs+=1
	condition.Broadcast()
	etec3702.Delay();
	etec3702.Output("Macaque on rope");
	monkLock.Unlock()

    monkLock.Lock()
	numMacs-=1
	if totalMacs >= 49 && numMacs == 0{
		totalBabs = 0
	}
	condition.Broadcast()
	etec3702.Delay();
	etec3702.Output("Macaque off rope");
	monkLock.Unlock()    
}
