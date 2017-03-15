package main

import "bufio"
import "os"
import "log"
import "fmt"
import "image"
import "image/gif"
import "image/color"




func main(){
    grid1,grid2 := readInput(os.Args[1]);

    var src, dest, tmp [][]byte
    src=grid1
    dest=grid2
    
    numrows := len(src)
    numcols := len(src[0])

    pic,pal := makeGif(numrows,numcols)
        
    //let the game begin!
    
    for iters:=0;iters<200;iters++ {
        updateState(src,dest,numrows,numcols);

        outputAnimationFrame( &pic, &pal, dest, numrows, numcols );
        
        //swap roles of source and destination
        tmp = src
        src = dest
        dest = tmp
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

func readInput(filename string) ([][]byte, [][]byte) {
    
    fp,err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal("Unable to open")
    }
    
    //read the input file, dimensioning both buffers simultaneously
    scanner := bufio.NewScanner(fp)
    var grid1 [][]byte;
    var grid2 [][]byte;
    
    for scanner.Scan() {
        line := scanner.Text()
        grid1 = append(grid1,[]byte(line))
        grid2 = append(grid2,[]byte(line))
    }
    
    for r := 0; r < len(grid1); r++ {
        for c:=0;c<len(grid1[r]);c++{
            if grid1[r][c] == 42 {
                grid1[r][c] = 1
            } else {
                grid1[r][c] = 0
            }
        }
    }
    return grid1,grid2
}


func makeGif(numrows int , numcols int ) (gif.GIF, color.Palette) {
    
    //create a color palette for the GIF.
    //color 0 = black, 1=white
    var pal color.Palette
    var pix color.RGBA
    pix.R = 0
    pix.G = 0
    pix.B = 0
    pix.A = 255
    pal = append(pal,pix)
    pix.R = 255
    pix.G = 255
    pix.B = 255
    pix.A = 255
    pal = append(pal,pix)
    
    var g gif.GIF;
    
    //scale the output image size so it's easier to see
    multiplier := 8
    width := numcols*multiplier
    height := numrows*multiplier
    
    g.Config.Width = width
    g.Config.Height = height
    g.LoopCount = -1
    return g,pal
}

func updateState(src [][]byte, dest [][]byte, numrows int, numcols int ){
    for r:=0;r<numrows;r++ {
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
        }
    }
}

func outputAnimationFrame(g *gif.GIF, pal *color.Palette, 
                data [][]byte, numrows int, numcols int ){
                    
    img := image.NewPaletted( 
        image.Rectangle{ 
            image.Point{0,0}, 
            image.Point{g.Config.Width,g.Config.Height} },
        *pal )
    
    multiplier := 8
    
    pr:=0
    pc:=0
    for r:=0;r<numrows;r++ {
        for rr:=0;rr<multiplier; rr,pr = rr+1,pr+1 {
            pc=0
            for c:=0;c<numcols;c++ {
                for cc:=0;cc<multiplier; cc,pc = cc+1,pc+1 {
                    if data[r][c] == 0 {
                        img.SetColorIndex(pc,pr,0)
                    } else {
                        img.SetColorIndex(pc,pr,1)
                    }   
                }
            }
        }
    }

    g.Image = append(g.Image,img)
    g.Delay = append(g.Delay, 10)
}
    
