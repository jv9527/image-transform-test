package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"
	"gopkg.in/h2non/bimg.v1"
)

var (
	ROOT_FOLDER = "./images/100-images/"
	IMG_FOLDER = "images/jpg/"
)

const (
	NUM_OF_WORKER= 4
	PRE_SIZE= 700
	QUALITY= 70
)

func main(){
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
				output: ROOT_FOLDER + "results/" + "vips/" + fmt.Sprintf("%d/%d_", PRE_SIZE, i) + file.Name(),
				size: PRE_SIZE,
				resizer: useResize,
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
	resizer func([]byte, int) ([]byte, error)
}


func useResize(buf []byte, size int) ([]byte, error) {
	return bimg.Resize(buf, bimg.Options{

		Width:   size,
		//Height:  size,
		StripMetadata: true,
		Quality: QUALITY,
		Type: bimg.JPEG,
	})
}

func runWorker(numOfWorker int) (chan<- *args, *sync.WaitGroup) {
	var wg sync.WaitGroup

	chAdd := make(chan *args)
	wg.Add(numOfWorker)

	for i := 0 ; i < numOfWorker ; i++ {
		go func(i int, wg *sync.WaitGroup){
			fmt.Println("Running worker: ", i)
			defer fmt.Println("Stopping worker: ", i)
			for {
				select {
				case a, ok := <-chAdd:
					if !ok {
						wg.Done()
						return
					}
					TestResizeVips(&a.input, &a.output, a.size, a.resizer)
				}
			}
		}(i, &wg)
	}

	return chAdd, &wg
}

func TestResizeVips(input, output *string, size int, resizer func([]byte, int) ([]byte, error)) {
	// Open image and store on buffer
	buff, err := bimg.Read(*input)
	if err != nil {
		fmt.Printf("Error while read images. err: %v\n", err)
		return
	}

	// Resize using Thumbnail
	bNewImage, err := resizer(buff, size)
	if err != nil {
		fmt.Printf("error while resize image using thumbnail. err: %v\n", err)
		return
	}

	// Write image
	if err := bimg.Write(*output, bNewImage); err != nil {
		fmt.Printf("error while write image. err: %v\n", err)
		return
	}
}