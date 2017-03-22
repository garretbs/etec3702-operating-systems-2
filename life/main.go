package main

import "etec3702"
import "bufio"
import "os"
import "log"
import "fmt"
import "image"
import "image/gif"
import "image/color"
import "sync"


const numThreads = 7

var tmp [][]byte
var src,dest = readInput(os.Args[1])
var numrows = len(src)
var numcols = len(src[0])
var pic,pal =  makeGif()

var globalsLock sync.Mutex


func main(){

	
	//create threads
	var chans [numThreads]chan bool
	
	for i:=0;i<numThreads;i++{
		chans[i] = make(chan bool)
		go lifeThread(i, chans[i]);
	}
	
	for i:=0;i<numThreads;i++{//wait for all threads to finish all iterations
		_ = <- chans[i]
	}
	
	//write output image
    fp,err := os.Create("out.gif")
    if err != nil {
        log.Fatal("Could not open output")
    }
    gif.EncodeAll(fp,&pic)
    fp.Close()
    fmt.Println("Done")
}

func lifeThread(threadID int, c chan bool){
	//let the game begin!
    
    for iters:=0;iters<200;iters++ {
		
        updateState(threadID);
		
		barrier();
		
		
		//only one thread needs to do the things below
		if(threadID == 0){
			globalsLock.Lock()
			outputAnimationFrame();
			
			//swap roles of source and destination
			tmp = src
			src = dest
			dest = tmp
			globalsLock.Unlock()
			
			
		}
		resetBarrier();
		
		
    }
	c <- true
}

func readInput(filename string) ([][]byte, [][]byte) {
    
    fp,err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal("Unable to open")
    }
    
    //read the input file, dimensioning both buffers simultaneously
    scanner := bufio.NewScanner(fp)
    var g1 [][]byte;
    var g2 [][]byte;
    
    for scanner.Scan() {
        line := scanner.Text()
        g1 = append(g1,[]byte(line))
        g2 = append(g2,[]byte(line))
    }
    
    for r := 0; r < len(g1); r++ {
        for c:=0;c<len(g1[r]);c++{
            if g1[r][c] == 42 {
                g1[r][c] = 1
            } else {
                g1[r][c] = 0
            }
        }
    }
    return g1,g2
}


func makeGif() (gif.GIF, color.Palette) {
    
    //create a color palette for the GIF.
    //color 0 = black, 1=white
    var palette color.Palette
    var pix color.RGBA
    pix.R = 0
    pix.G = 0
    pix.B = 0
    pix.A = 255
    palette = append(palette,pix)
    pix.R = 255
    pix.G = 255
    pix.B = 255
    pix.A = 255
    palette = append(palette,pix)
    
    var g gif.GIF;
    
    //scale the output image size so it's easier to see
    multiplier := 8
    width := numcols*multiplier
    height := numrows*multiplier
    
    g.Config.Width = width
    g.Config.Height = height
    g.LoopCount = -1
    return g,palette
}

func updateState(threadID int ){
	//needs to only do certain rows based on threadID
    for r:=threadID;r<numrows;r+=numThreads {
        for c:=0;c<numcols;c++ {
            var up,down,left,right int
            up = r-1
            if up < 0 {
                up = numrows-1
            }
            down = r+1
            if down >= numrows {
                down = 0
            }
            left = c-1
            if left < 0 {
                left = numcols-1
            }
            right = c+1
            if right >= numcols {
                right = 0
            }
            
			globalsLock.Lock()
            n := src[up][c]
            s := src[down][c]
            w := src[r][left]
            e := src[r][right]
            nw := src[up][left]
            ne := src[up][right]
            sw := src[down][left]
            se := src[down][right]
            total := nw+n+ne+w+e+sw+s+se
			
			
            if total < 2 {
                dest[r][c]=0
            } else if total == 2 {
                dest[r][c] = src[r][c]
            } else if total == 3 {
                dest[r][c] = 1
            } else {
                dest[r][c] = 0
            }
			globalsLock.Unlock()
        }
    }
}

func outputAnimationFrame(){
	
	img := image.NewPaletted( 
		image.Rectangle{ 
			image.Point{0,0}, 
			image.Point{pic.Config.Width,pic.Config.Height} },
			pal )
    
    multiplier := 8
    
    pr:=0
    pc:=0
    for r:=0;r<numrows;r++ {
        for rr:=0;rr<multiplier; rr,pr = rr+1,pr+1 {
            pc=0
            for c:=0;c<numcols;c++ {
                for cc:=0;cc<multiplier; cc,pc = cc+1,pc+1 {
                    if dest[r][c] == 0 {
                        img.SetColorIndex(pc,pr,0)
                    } else {
                        img.SetColorIndex(pc,pr,1)
                    }   
                }
            }
        }
    }

    pic.Image = append(pic.Image,img)
    pic.Delay = append(pic.Delay, 10)
}

//Barrier stuff
var sem1 = etec3702.NewSemaphore(0)
var sem2 = etec3702.NewSemaphore(0)
var barrierLock sync.Mutex
var threadsWaiting = numThreads

func barrier(){
	barrierLock.Lock()
	threadsWaiting--
	if threadsWaiting == 0{
		for i:=0;i<numThreads;i++{
			sem1.Release()
		}
	}
	barrierLock.Unlock()
	sem1.Acquire()
}

func resetBarrier(){
	barrierLock.Lock()
	threadsWaiting++
	if threadsWaiting == numThreads{
		for i:=0;i<numThreads;i++{
			sem2.Release()
		}
	}
	barrierLock.Unlock()
	sem2.Acquire()
}

