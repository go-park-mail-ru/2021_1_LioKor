package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"image"
	"image/jpeg"
	_ "image/png" // to allow png uploading and converting to jpg
)

type Config struct {
	Host              string
	Port              int
	AllowedOrigins    []string
	AvatarStoragePath string
}

/* Converts dataURL to file and saves it. Returns file path. Only jpg and png supported
Usage example:

path, err := dataURLToFile("wolchara", newData.AvatarURL, 500)
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println(path) // wolchara.jpg
}
*/
func DataURLToFile(path string, dataURL string, maxSizeKB int) (string, error) {
	if dataURL == "" {
		return "", nil
	}

	splittedURL := strings.Split(dataURL, ",")
	if len(splittedURL) != 2 {
		return "", errors.New("incorrect data url")
	}

	meta := splittedURL[0]
	var ext string
	if strings.Index(meta, "image/jpeg") != -1 {
		ext = "jpg"
	} else if strings.Index(meta, "image/png") != -1 {
		ext = "png"
	} else {
		return "", errors.New("forbidden data format")
	}

	base64Data := splittedURL[1]
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", errors.New("unable to base64 decode")
	}

	if len(decoded) > maxSizeKB*1024 {
		return "", errors.New("image is too big")
	}

	img, format, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return "", err
	}

	if (format == "jpeg") || (format == "png") {
		ext = "jpg" // because we convert both jpg and png to jpg
	} else {
		return "", errors.New("forbidden data format")
	}

	path += "." + ext
	f, err := os.Create(path)
	if err != nil {
		return "", errors.New("unable to save file")
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

	return path, nil
}

func GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	randNumStr := strconv.Itoa(rand.Intn(32000))

	h := sha256.New()
	h.Write([]byte(randNumStr))
	return hex.EncodeToString(h.Sum(nil))
}
