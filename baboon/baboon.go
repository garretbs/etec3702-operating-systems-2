package main
import "etec3702"

func baboon(){
	
	monkLock.Lock()
	for numMacs > 0 || numBabs >= 3 || totalBabs >= 49{ //don't go while max baboons or any macaques
		condition.Wait()
	}
	numBabs+=1
	totalBabs+=1
	condition.Broadcast()
	etec3702.Delay();
	etec3702.Output("Baboon on rope");
	monkLock.Unlock()

    monkLock.Lock()
	numBabs-=1
	if totalBabs >= 49 && numBabs == 0{
		totalMacs = 0
	}
	condition.Broadcast()
	etec3702.Delay();
	etec3702.Output("Baboon off rope");
    monkLock.Unlock()
    
}
