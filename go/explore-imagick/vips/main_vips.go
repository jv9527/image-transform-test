package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"gopkg.in/h2non/bimg.v1"
	"github.com/gorilla/mux"
)

var (
	ROOT_FOLDER = "./images/100-images/"
	IMG_FOLDER = "images/jpg/"
)

const (
	NUM_OF_WORKER= 4
	PRE_SIZE= 700
	QUALITY= 85
)

type args struct {
	imgB    []byte
	input 	string
	output 	string
	size 	int
	chDone  chan interface{}
	resizer func([]byte, int) ([]byte, error)
}

var chAdd chan<- *args
var wg *sync.WaitGroup

func main(){
	chAdd, wg = runWorker(NUM_OF_WORKER)

	r := mux.NewRouter()

	r.HandleFunc("/resize", HandleResizeRequest).Methods("POST")

	if http.ListenAndServe(":8081", r) != nil {
		close(chAdd)
		wg.Wait()
	}
}

func HandleResizeRequest(w http.ResponseWriter, r *http.Request) {
	blob, fileHeader, err := r.FormFile("file_upload")
	if err != nil {
		_ = fmt.Errorf("error while read image blob: %v", err)

		_, _ = fmt.Fprint(w, err)
		return
	}

	// Get Bytes image
	buf, err := ioutil.ReadAll(blob)
	if err != nil {
		fmt.Printf("error while read blob: %v", err)
		_, _ = fmt.Fprint(w, err)
	}

	// Send to Image Resizer
	chDone := make(chan interface{})
	chAdd <- &args{
		imgB: buf,
		output: ROOT_FOLDER + "results/" + "vips/" + fmt.Sprintf("%d/%d_", PRE_SIZE, time.Now().UnixNano()) + fileHeader.Filename,
		size: PRE_SIZE,
		chDone: chDone,
		resizer: useResize,
	}

	// Wait to process
	<-chDone
}

func main2(){
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


func useResize(buf []byte, size int) ([]byte, error) {
	return bimg.Resize(buf, bimg.Options{

		Width:   size,
		Height:  size,
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
					TestResizeVips(a.imgB, &a.input, &a.output, a.size, a.resizer)

					close(a.chDone)
				}
			}
		}(i, &wg)
	}

	return chAdd, &wg
}

func TestResizeVips(imgB []byte, input, output *string, size int, resizer func([]byte, int) ([]byte, error)) {
	var err error
	if imgB == nil {
		// Open image and store on buffer
		imgB, err = bimg.Read(*input)
		if err != nil {
			fmt.Printf("Error while read images. err: %v\n", err)
			return
		}
	}

	// Resize using Thumbnail
	bNewImage, err := resizer(imgB, size)
	if err != nil {
		fmt.Printf("error while resize image using thumbnail. err: %v\n", err)
		return
	}

	//// Write image
	//if err := bimg.Write(*output, bNewImage); err != nil {
	//	fmt.Printf("error while write image. err: %v\n", err)
	//	return
	//}
}