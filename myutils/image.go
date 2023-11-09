package myutils

import (
	"os"
	"path/filepath"
)

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

func ImageExtValidation(filename string) bool {
	corrects := []string{".jpeg", ".tiff", ".ai", ".psd", ".gif", ".jpg", ".png"}

	ext := filepath.Ext(filename)

	for _, value := range corrects {
		if value == ext {
			return true
		}
	}

	return false
}