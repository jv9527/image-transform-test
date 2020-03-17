package main

import (
	"fmt"
	"github.com/discordapp/lilliput"
	"io/ioutil"
	"sync"
	"time"
)

var (
	ROOT_FOLDER = "./images/100-images/"
	IMG_FOLDER = "images/jpg/"
)

const (
	NUM_OF_WORKER= 4
	PRE_SIZE= 700
	QUALITY= 60
)

func main() {
	// Get All Files
	files, _ := ioutil.ReadDir(ROOT_FOLDER + IMG_FOLDER)

	chAdd, wg := runWorker(NUM_OF_WORKER)
	//time.Sleep(5 * time.Second)
	start := time.Now()

	for i := 0 ; i <  1 ; i++ {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			chAdd <- &args{
				input: ROOT_FOLDER + IMG_FOLDER + file.Name(),
				output: ROOT_FOLDER + "results/" + "lilliput/" + fmt.Sprintf("%d/%d_", PRE_SIZE, i) + file.Name(),
				size: PRE_SIZE,
			}
		}
	}

	close(chAdd)
	wg.Wait()
	fmt.Println(time.Since(start))
}

type args struct {
	input 	string
	output 	string
	size 	int
}

func runWorker(numOfWorker int) (chan<- *args, *sync.WaitGroup) {
	var wg sync.WaitGroup

	chAdd := make(chan *args)
	wg.Add(numOfWorker)

	for i := 0 ; i < numOfWorker ; i++ {
		go func(i int, wg *sync.WaitGroup){
			fmt.Println("Running worker: ", i)
			defer fmt.Println("Stopping worker: ", i)

			imgOps := lilliput.NewImageOps(8000)
			defer imgOps.Close()

			outputImg := make([]byte, 1*1024*1024)

			for {
				select {
				case a, ok := <-chAdd:
					if !ok {
						wg.Done()
						return
					}
					TestResizeLiliput(&a.input, &a.output, a.size, imgOps, outputImg)
				}
			}
		}(i, &wg)
	}

	return chAdd, &wg
}

func TestResizeLiliput(input, output *string, size int, imgOps *lilliput.ImageOps, outputImg []byte) {
	inputBuf, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("Error while read images. err: %v\n", err)
		return
	}

	// Create Decoder
	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		fmt.Printf("error while decode image: %v\n", err)
		return
	}
	defer decoder.Close()

	opts := &lilliput.ImageOptions{
		NormalizeOrientation:true,
		FileType:             ".jpg",
		Width:                size,
		Height:               size,
		ResizeMethod:         lilliput.ImageOpsResize,
		//NormalizeOrientation: true,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: QUALITY},
	}

	// Resize
	outputImg, err = imgOps.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error while resize image: %v\n", err)
		return
	}

	// Write image
	if err = ioutil.WriteFile(*output, outputImg, 0664); err != nil {
		fmt.Printf("error while write image. err: %v\n", err)
		return
	}
}

//func TestResizeLiliput(size int, prefix string, _ *lilliput.ImageOps) {
//	start := time.Now()
//
//	imgOps := lilliput.NewImageOps(8000)
//	defer imgOps.Close()
//
//	outputImg := make([]byte, 1*1024*1024)
//
//	// Resize
//	for j := 0 ; j < 10; j++ {
//		for i := 0; i < 30 ; i++ {
//			inputBuf, err := ioutil.ReadFile(ROOT_FOLDER + FILE_NAME[i])
//			if err != nil {
//				fmt.Printf("Error while read images. err: %v\n", err)
//				return
//			}
//
//			// Create Decoder
//			decoder, err := lilliput.NewDecoder(inputBuf)
//			if err != nil {
//				fmt.Printf("error while decode image: %v\n", err)
//				return
//			}
//			defer decoder.Close()
//
//			opts := &lilliput.ImageOptions{
//				FileType:             ".jpg",
//				Width:                size,
//				Height:               size,
//				ResizeMethod:         lilliput.ImageOpsFit,
//				NormalizeOrientation: true,
//				EncodeOptions:        map[int]int{lilliput.JpegQuality: 80},
//			}
//
//			// Resize
//			outputImg, err = imgOps.Transform(decoder, opts, outputImg)
//			if err != nil {
//				fmt.Printf("error while resize image: %v\n", err)
//				return
//			}
//
//			// Write image
//			if err = ioutil.WriteFile(ROOT_FOLDER + "results/" + prefix + fmt.Sprintf("%d_", j) + FILE_NAME[i], outputImg, 0400); err != nil {
//				fmt.Printf("error while write image. err: %v\n", err)
//				return
//			}
//		}
//	}
//
//	fmt.Println(time.Since(start))
//}
