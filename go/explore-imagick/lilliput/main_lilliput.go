package main

import (
	"fmt"
	"github.com/discordapp/lilliput"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	ROOT_FOLDER = "./images/100-images/"
	IMG_FOLDER  = "images/jpg/"
)

const (
	NUM_OF_WORKER = 4
	PRE_SIZE      = 200
	QUALITY       = 85
)

type args struct {
	imgB   []byte
	input  string
	output string
	size   int
	chDone chan interface{}
}

var chAdd chan<- *args
var wg *sync.WaitGroup

func main() {
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
	defer blob.Close()

	// Get Bytes image
	buf, err := ioutil.ReadAll(blob)
	if err != nil {
		fmt.Printf("error while read blob: %v", err)
		_, _ = fmt.Fprint(w, err)
	}

	// Send to Image Resizer
	chDone := make(chan interface{})
	chAdd <- &args{
		imgB:   buf,
		output: ROOT_FOLDER + "results/" + "lilliput/" + fmt.Sprintf("%d/%d_", PRE_SIZE, time.Now().UnixNano()) + fileHeader.Filename,
		size:   PRE_SIZE,
		chDone: chDone,
	}

	// Wait to process
	<-chDone
}

func main2() {
	// Get All Files
	files, _ := ioutil.ReadDir(ROOT_FOLDER + IMG_FOLDER)

	chAdd, wg := runWorker(NUM_OF_WORKER)
	//time.Sleep(5 * time.Second)
	start := time.Now()

	for i := 0; i < 100; i++ {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			chAdd <- &args{
				input:  ROOT_FOLDER + IMG_FOLDER + file.Name(),
				output: ROOT_FOLDER + "results/" + "lilliput/" + fmt.Sprintf("%d/%d_", PRE_SIZE, i) + file.Name(),
				size:   PRE_SIZE,
			}
		}
	}

	close(chAdd)
	wg.Wait()
	fmt.Println(time.Since(start))
}

func runWorker(numOfWorker int) (chan<- *args, *sync.WaitGroup) {
	var wg sync.WaitGroup

	chAdd := make(chan *args)
	wg.Add(numOfWorker)

	for i := 0; i < numOfWorker; i++ {
		go func(i int, wg *sync.WaitGroup) {
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
					TestResizeLiliput(a.imgB, &a.input, &a.output, a.size, imgOps, outputImg)

					close(a.chDone)
				}
			}
		}(i, &wg)
	}

	return chAdd, &wg
}

func TestResizeLiliput(imgB []byte, input, output *string, size int, imgOps *lilliput.ImageOps, outputImg []byte) {
	defer imgOps.Clear()

	var err error
	if imgB == nil {
		imgB, err = ioutil.ReadFile(*input)
		if err != nil {
			fmt.Printf("Error while read images. err: %v\n", err)
			return
		}
	}

	// Create Decoder
	decoder, err := lilliput.NewDecoder(imgB)
	if err != nil {
		fmt.Printf("error while decode image: %v\n", err)
		return
	}
	defer decoder.Close()

	opts := &lilliput.ImageOptions{
		NormalizeOrientation: true,
		FileType:             ".jpg",
		Width:                size,
		Height:               size,
		ResizeMethod:         lilliput.ImageOpsResize,
		EncodeOptions:        map[int]int{lilliput.JpegQuality: QUALITY},
	}

	// Resize
	outputImg, err = imgOps.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error while resize image: %v\n", err)
		return
	}

	//// Write image
	//if err = ioutil.WriteFile(*output, outputImg, 0664); err != nil {
	//	fmt.Printf("error while write image. err: %v\n", err)
	//	return
	//}
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
