package main

import "os"
import "fmt"
import "math"
import "archive/zip"
import "image"
import _ "image/png"
import _ "image/jpeg"

var B = NewBlackboard()

func main(){
	if len(os.Args) < 2{
		fmt.Println("No file specified")
		return
	}
	fname := os.Args[1]
	
	go ioStuff(fname);
	
	const numThreads = 4
	for i:=0;i<numThreads;i++{
		go compareImages();
	}
	
	_,scanners := B.get("IO Done")
	numToScan := scanners.(int) - 1 //there are only this many comparisons
	
	for i:=0;i<numToScan;i++{
		//fmt.Println(i)
		B.get("scanned")
	}
	
	for i:=0;i<numThreads;i++{
		B.put("DIE", 0)
	}
	fmt.Println("Done")
	
}

func compareImages(){	
	for{
		tag,imgs := B.get("compare", "DIE")
		
		if tag == "DIE"{
			break
		}
		files := imgs.([]*zip.File)
		
		f1data,_ := (*files[0]).Open()
		f2data,_ := (*files[1]).Open()
		
		img1,_,_ := image.Decode(f1data)
		img2,_,_ := image.Decode(f2data)
		
		f1data.Close()
		f2data.Close()
		
		
		if !areSame(&img1, &img2){
			fmt.Println((*files[0]).Name, "and", (*files[1]).Name)
		}
		B.put("scanned", 0)
	}
}

func areSame(img1, img2 *image.Image) bool{
	var r1,g1,b1 uint32
	var r2,g2,b2 uint32	
	var rgb1,rgb2 float64
	var difference float64
	
	diffPixels := 0
	totalPixels := (*img1).Bounds().Max.X * (*img1).Bounds().Max.Y
	threshold := 5.0
	
	for x:=0;x<(*img1).Bounds().Max.X;x+=4{
		for y:=0;y<(*img1).Bounds().Max.Y;y+=4{
			r1,g1,b1,_ = (*img1).At(x, y).RGBA()
			r2,g2,b2,_ = (*img2).At(x, y).RGBA()
			
			rgb1 = float64((r1+g1+b1) >> 8)
			rgb2 = float64((r2+g2+b2) >> 8)
			difference = math.Abs(rgb2-rgb1)
			
			if difference > threshold{
				//fmt.Println(difference)
				diffPixels++
			}
		}
	}
	
	return (diffPixels < (totalPixels >> 7))
}

func ioStuff(input string){
	fp,e := zip.OpenReader(input)
	if e != nil{
		fmt.Println("Error opening file");
		os.Exit(-1)
	}

	var prevImg *zip.File
	for _,currentImg := range fp.File{
		//cycle through each of the zip file's files
		if prevImg != nil{ //only manage comparison with previous image
			B.put("compare", []*zip.File{prevImg, currentImg})
		}
		prevImg = currentImg
	}
	B.put("IO Done", len(fp.File))
}
