package main
import "sync"

var monkLock sync.Mutex
var condition = sync.NewCond(&monkLock)

var numBabs,numMacs int
var totalBabs,totalMacs int
