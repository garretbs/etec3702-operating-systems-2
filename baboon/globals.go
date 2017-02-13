package main
import "etec3702"

var babEntry = etec3702.NewSemaphore(1)
var babExit = etec3702.NewSemaphore(0)

var macEntry = etec3702.NewSemaphore(0)
var macExit = etec3702.NewSemaphore(0)

var babLock = etec3702.NewSemaphore(1)
var macLock = etec3702.NewSemaphore(1)

var numBabs,numMacs int

