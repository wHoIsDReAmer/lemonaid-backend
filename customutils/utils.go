package customutils

import (
	"math/rand"
	"os"
	"time"
)

var (
	src = rand.NewSource(time.Now().Unix())
)

func RandI(a int, b int) int {
	r := rand.New(src)

	return r.Intn(b-a) + a
}

func RandF(a float64, b float64) float64 {
	r := rand.New(src)

	return r.Float64()*(b-a) + a
}

func ImageProcessing(buffer []byte, quality int, filename string) error {
	//filename := strings.Replace(uuid.New().String(), "-", "", -1) + ".webp"

	//io.Copy()

	//converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	//if err != nil {
	//	return filename, err
	//}
	//
	//processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	//if err != nil {
	//	return filename, err
	//}
	//
	//os.MkdirAll("./public/contents", 0777)
	//
	//writeError := bimg.Write(fmt.Sprintf("./public/contents/%s", filename), processed)
	//if writeError != nil {
	//	return filename, writeError
	//}

	os.MkdirAll("./public/contents", 0777)

	dst, _ := os.Create("./public/contents/" + filename)
	defer dst.Close()

	_, err := dst.Write(buffer)

	return err
}
