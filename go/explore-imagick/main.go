package main_2

import (
	"fmt"
	"gopkg.in/h2non/bimg.v1"
	"gopkg.in/gographics/imagick.v2/imagick"
	"github.com/gographics/gmagick"
	"sync"

	//"os"
	"time"
)

var (
	ROOT_FOLDER = "./images/100-images/"
	FILE_NAME = []string {
		"1_test.jpg", "2_test.jpg", "3_test.jpg", "5_test.jpg",
		"6_test.jpg", "8_test.jpg", "9_test.jpg", "10_test.jpg",
		"11_test.jpg", "12_test.jpg", "13_test.jpg", "14_test.jpg", "15_test.jpg",
		"16_test.jpg", "17_test.jpg", "18_test.jpg", "19_test.jpg", "20_test.jpg",
		"21_test.jpg", "22_test.jpg", "23_test.jpg", "24_test.jpg", "25_test.jpg",
		"26_test.jpg", "27_test.jpg", "28_test.jpg", "29_test.jpg", "30_test.jpg",
	}
)

func main20(){
	imagick.Initialize()
	defer imagick.Terminate()

	gmagick.Initialize()
	defer gmagick.Terminate()

	// Vips
	fmt.Println("===== Vips - Thumbnail =====")
	//TestResizeVips(200, "vips/thumbnail/200/", useThumbnail)
	//TestResizeVips(700, "vips/thumbnail/700/", useThumbnail)
	fmt.Println("===== Vips - Resize =====")
	//TestResizeVips(200, "vips/resize/200/", useResize)
	//TestResizeVips(700, "vips/resize/700/", useResize)

	fmt.Println("===== Gmagick =====")
	TestResizeGmagick(200, "gmagick/200/")
	TestResizeGmagick(700, "gmagick/700/")

	fmt.Println("===== Imagick =====")
	TestResizeImagick(200,"imagick/200/")
	TestResizeImagick(700,"imagick/700/")
}

func main(){
	chAdd, wg := runWorker(1)
	//time.Sleep(5 * time.Second)
	start := time.Now()

	for i := 0 ; i <  1 ; i++ {
		for j := 0 ; j < 1 ; j++ {
			chAdd <- &args{
				input: ROOT_FOLDER + FILE_NAME[j],
				output: ROOT_FOLDER + "results/" + "vips/resize/200/" + fmt.Sprintf("%d_", i) + FILE_NAME[j],
				size: 700,
				resizer: useResize,
			}
		}
	}

	close(chAdd)
	wg.Wait()
	fmt.Println(time.Since(start))


	//TestResizeVips(200, "thumbnail/200/", useThumbnail)
	//TestResizeVips(700, "thumbnail/700/", useThumbnail)
	//TestResizeVips(700, "vips/resize/700/", useResize)
	//TestResizeVips(700, "resize/700/", useResize)
}

func main4(){
	imagick.Initialize()
	defer imagick.Terminate()

	TestResizeImagick(200,"imagick/200/")
	TestResizeImagick(700,"imagick/700/")
}

func main3(){
	gmagick.Initialize()
	defer gmagick.Terminate()

	TestResizeGmagick(200, "gmagick/200/")
	TestResizeGmagick(700, "gmagick/700/")
}

type args struct {
	input 	string
	output 	string
	size 	int
	resizer func([]byte, int) ([]byte, error)
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
	//buff, err := bimg.Read(ROOT_FOLDER + FILE_NAME[i])
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
	//if err := bimg.Write(ROOT_FOLDER + "results/" + prefix + fmt.Sprintf("%d", j) + FILE_NAME[i], bNewImage); err != nil {
	if err := bimg.Write(*output + ".webp", bNewImage); err != nil {
		fmt.Printf("error while write image. err: %v\n", err)
		return
	}
}

func TestResizeVips2(size int, prefix string, resizer func([]byte, int) ([]byte, error)) {
	start := time.Now()

	for j:= 0 ; j < 100; j++{
		for i:=0 ; i < 30; i++ {
			// Open image and store on buffer
			buff, err := bimg.Read(ROOT_FOLDER + FILE_NAME[i])
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
			if err := bimg.Write(ROOT_FOLDER + "results/" + prefix + fmt.Sprintf("%d", j) + FILE_NAME[i], bNewImage); err != nil {
				fmt.Printf("error while write image. err: %v\n", err)
				return
			}
		}
	}
	fmt.Println(time.Since(start))
}

func TestResizeGmagick(size uint, prefix string) {
	mw := gmagick.NewMagickWand()
	defer mw.Destroy()

	start := time.Now()

	// Resize image
	for j:= 0 ; j < 10; j++{
		for i := 0; i < 30 ; i++ {
			if err := mw.ReadImage(ROOT_FOLDER + FILE_NAME[i]); err != nil {
				fmt.Printf("error while read image. err: %v\n", err)
				return
			}

			// Resize
			if err := mw.ResizeImage(size, size, gmagick.FILTER_UNDEFINED, 1); err != nil {
				fmt.Printf("error while resize image: %v\n", err)
				return
			}

			// Save image
			if err := mw.WriteImage(ROOT_FOLDER + "results/" + prefix + FILE_NAME[i]); err != nil {
				fmt.Printf("error while save image: %v\n", err)
				return
			}
		}
	}

	fmt.Println(time.Since(start))
}

func TestResizeImagick(size uint, prefix string) {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	start := time.Now()

	// Resize Image
	for j:= 0 ; j < 10; j++{
		for i := 0; i < 30; i++ {
			if err := mw.ReadImage(ROOT_FOLDER+FILE_NAME[i]); err != nil {
				fmt.Printf("error while read image. err: %v\n", err)
				return
			}

			_ = mw.NormalizeImage()

			// Resize using thumbnail
			if err := mw.ThumbnailImage(size, size); err != nil {
				fmt.Printf("error while resize image. err: %v\n", err)
				return
			}



			if err := mw.WriteImage(ROOT_FOLDER + "results/" + prefix + FILE_NAME[i]); err != nil {
				fmt.Printf("error while image: %v\n", err)
				return
			}
		}
	}

	fmt.Println(time.Since(start))
}



func useThumbnail(buf []byte, size int) ([]byte, error) {
	// Create new Image
	newImage := bimg.NewImage(buf)

	return newImage.Thumbnail(size)
}

func useResize(buf []byte, size int) ([]byte, error) {
	return bimg.Resize(buf, bimg.Options{
		Width:   size,
		Height:  size,
		Quality: 90,
		Type: bimg.WEBP,
	})
}