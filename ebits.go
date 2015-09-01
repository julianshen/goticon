package main

import (
	"crypto/sha512"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Assets map[string]([]int)

const (
	IMAGE_WIDTH = 400
)

var (
	compose_order = [...]string{"background", "face", "eye", "mouth", "hair", "clothes"}
	cache         = make(map[string]image.Image)
	male          = make(Assets)
	female        = make(Assets)
)

func InitAssets() {
	if len(male) > 0 && len(female) > 0 {
		return
	}

	filepath.Walk("img", func(path string, info os.FileInfo, err error) error {
		re := regexp.MustCompile("/(male|female)/(\\D+)(\\d+)\\.png")
		found := re.FindStringSubmatch(path)

		if len(found) > 0 {
			gender := found[1]
			datatype := found[2]
			index, err := strconv.Atoi(found[3])

			if err == nil {
				if gender == "male" {
					male[datatype] = append(male[datatype], index)
				} else if gender == "female" {
					female[datatype] = append(female[datatype], index)
				}
			}
		}

		return nil
	})

	for _, value := range male {
		sort.Ints(value)
	}

	for _, value := range female {
		sort.Ints(value)
	}
}

func LoadImage(name string) image.Image {
	cached := cache[name]

	if cached != nil {
		return cached
	}

	imgReader, err := os.Open(name)
	if err != nil {
		log.Fatal("Asset " + name + " is not existed")
	}
	defer imgReader.Close()

	img, _, err := image.Decode(imgReader)

	if err != nil {
		log.Fatal("Asset " + name + " is corrupted?")
	}

	cache[name] = img
	return img
}

func GenerateIdenticon8bits(gender string, data []byte) image.Image {
	sum := sha512.Sum512(data)

	var assets Assets
	switch {
	case gender == "male":
		assets = male
	case gender == "female":
		assets = female
	default:
		log.Fatal("Unknow gender " + gender)
	}

	img := image.NewNRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_WIDTH))
	for i, layer_name := range compose_order {
		layer_sels := assets[layer_name]
		data := sum[i]
		num_sels := len(layer_sels)

		layer_index := layer_sels[int(data)%num_sels]

		log.Println("use img/" + gender + "/" + layer_name + strconv.Itoa(layer_index) + ".png")
		layer_img := LoadImage("img/" + gender + "/" + layer_name + strconv.Itoa(layer_index) + ".png")
		draw.Draw(img, image.Rect(0, 0, IMAGE_WIDTH, IMAGE_WIDTH), layer_img, image.Point{0, 0}, draw.Over)
	}

	return img
}
