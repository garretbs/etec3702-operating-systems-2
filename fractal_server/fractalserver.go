package main

import "net/http"
import "fmt"
import "os"
import "io"
import "image"
import "image/png"
import "strconv"
import "time"
import "math"

func main(){
    http.HandleFunc("/", indexPage )
    http.HandleFunc("/fractal", fractalImage )
    http.ListenAndServe(":8989", nil)
}

func indexPage(w http.ResponseWriter, r *http.Request){
    fmt.Println(r.Method,r.URL)
    fp,err := os.Open("goindex.html")
    defer fp.Close()
    if err != nil {
        fmt.Println("File not found")
        w.WriteHeader(404)
    } else {
        io.Copy( w, fp )
    }
}

func fractalImage(wr http.ResponseWriter, r *http.Request ){
    
    fmt.Println(r.Method,r.URL)
    
    startTime := time.Now()
    
    w,_ := strconv.Atoi( r.URL.Query()["w"][0] )
    h,_ := strconv.Atoi( r.URL.Query()["h"][0] )
    xmin,_ := strconv.ParseFloat( r.URL.Query()["xmin"][0], 64 )
    xmax,_ := strconv.ParseFloat( r.URL.Query()["xmax"][0], 64 )
    ymin,_ := strconv.ParseFloat( r.URL.Query()["ymin"][0], 64 )
    ymax,_ := strconv.ParseFloat( r.URL.Query()["ymax"][0], 64 )
    maxiter,_ := strconv.ParseFloat( r.URL.Query()["maxiter"][0], 64 )
    
    img := image.NewRGBA( image.Rectangle{ image.Point{0,0}, image.Point{w,h} } );

	
	const num_threads = 4
	strip := w/num_threads
	deltaY := (ymax-ymin)/float64(h)
    deltaX := (xmax-xmin)/float64(w)
	
	var chans [num_threads]chan bool
	for i := range chans {
	   chans[i] = make(chan bool)
	}
	
	for i := range chans{
		go compute_parallel(img,h,strip,i,xmin,ymin,maxiter,deltaX,deltaY,chans[i])
	}
	
	for i := range chans{
		_ = <- chans[i]
	}
	
	// compute(img,w,h,xmin,xmax,ymin,ymax,maxiter)
    
    png.Encode( wr , img )
    
    totalTime := time.Since(startTime)
    fmt.Printf("Total time: %f seconds\n", totalTime.Seconds() )
    
}

func compute_parallel(img *image.RGBA, h,strip,idx int,xmin,ymin,maxiter,deltaX,deltaY float64,c chan bool ) image.Image {
	
	ix := strip*idx
	mx := strip*(idx+1)

    for y,py:=0,ymin; y<h ; y,py = y+1,py+deltaY {
        for x,px:=ix,xmin+(deltaX*float64(ix)); x<mx ; x,px = x+1,px+deltaX {
            idx := y*img.Stride + x*4
            iter := iterations_to_infinity( px,py,maxiter )
            r,g,b := map_color(iter, maxiter )
            img.Pix[idx] = r;      //red
            img.Pix[idx+1] = g;      //green
            img.Pix[idx+2] = b;    //blue
            img.Pix[idx+3] = 255;    //alpha
        }
    }
	c <- true
    return img
}

//Compute the fractal image for the given rectangle and image size.
func compute(img *image.RGBA, w,h int,xmin,xmax,ymin,ymax,maxiter float64 ) image.Image {
    
    deltaY := (ymax-ymin)/float64(h)
    deltaX := (xmax-xmin)/float64(w)
    
    for y,py:=0,ymin; y<h ; y,py = y+1,py+deltaY {
        for x,px:=0,xmin; x<w ; x,px = x+1,px+deltaX {
            idx := y*img.Stride + x*4
            iter := iterations_to_infinity( px,py,maxiter )
            r,g,b := map_color(iter, maxiter )
            img.Pix[idx] = r;      //red
            img.Pix[idx+1] = g;      //green
            img.Pix[idx+2] = b;    //blue
            img.Pix[idx+3] = 255;    //alpha
        }
    }

    return img
}

//For point x,y: See how many iterations we can do before we go off
//to infinity.
func iterations_to_infinity(x float64, y float64, maxi float64 ) float64{
    c := complex(x,y)
    z := complex(0,0)
    for k:=0.0;k<maxi;k++ {
        z = z*z
        z = z+c
        if real(z)*real(z) + imag(z)*imag(z) > 4 {
            return k
        }
    }
    return maxi
}


// Map a color to an RGB value
// When k=0, returns red
// As k approaches MAX_ITERS, the returned color
// will proceed through orange, yellow, green, blue, and purple
// If k >= MAX_ITERS, the returned color is black.
// Returned values: red, green, blue, in the range 0...255
func map_color(k float64, MAX_ITERS float64 ) (uint8,uint8,uint8) {
    //N. Schaller's algorithm to map
    //HSV to RGB values.
    //http://www.cs.rit.edu/~ncs/color/t_convert.html

    s := 0.8     //saturation
    v := 0.95    //value
    h := (float64(k)/MAX_ITERS) * 6.0       //hue
    
    if h >= 6 {
        v=0
    }
    
    ipart := math.Floor(h)
    fpart := h-ipart
    A := v*(1-s);
    B := v*(1-s*fpart);
    C := v*(1-s*(1-fpart));
    var r,g,b float64;
    
    if ipart == 0 {
        r=v; g=C; b=A;
    } else if ipart == 1 {
        r=B; g=v; b=A;
    } else if ipart == 2 {
        r=A; g=v; b=C;
    } else if ipart == 3 {
        r=A; g=B; b=v;
    } else if ipart == 4 {
        r=C; g=A; b=v;
    } else{
        r=v; g=A; b=B;
    }
    return uint8(math.Floor(r*255)), uint8(math.Floor(g*255)), uint8(math.Floor(b*255))
}

