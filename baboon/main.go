package main

import "math/rand"
import "etec3702"

func main(){
    for true {
        if (rand.Int31() & 1) == 0 {
            go macaque()
        } else {
            go baboon()
        }
        etec3702.Delay()
    }
}

    
