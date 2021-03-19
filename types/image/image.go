package image

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

func ExtractImgExt(fileId string) string {
	ext := path.Ext(fileId)
	if ext == "" {
		return ext
	}

	ext = strings.ToLower(ext[1:])
	switch ext {
	case "png", "jpg", "gif", "bmp", "jpeg":
		return ext
	default:
		return ""
	}
}

//制作缩略图
func Thumbnail(i image.Image, ext string, height, width int) ([]byte, error) {
	var maxHeight, maxWidth int
	maxHeight = height
	maxWidth = width
	if height == 0 {
		maxHeight = i.Bounds().Max.Y
	}
	if width == 0 {
		maxWidth = i.Bounds().Max.X
	}
	var r image.Image
	r = resize.Thumbnail(uint(maxWidth), uint(maxHeight), i, resize.Lanczos3)
	if height != 0 && width != 0 {
		r = imaging.Fill(r, maxWidth, maxHeight, imaging.Center, imaging.Lanczos)
	}
	return encode(r, strings.ToLower(ext))
}

// ThumbnailGif 制作 Gif 的缩略图
func ThumbnailGif(buf []byte, height, width int) ([]byte, error) {
	gifObj, err := gif.DecodeAll(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	var maxHeight, maxWidth int
	maxHeight = height
	maxWidth = width
	if height == 0 {
		maxHeight = gifObj.Config.Height
	}
	if width == 0 {
		maxWidth = gifObj.Config.Width
	}

	var palettedList []*image.Paletted
	for _, paletted := range gifObj.Image {
		// r := resize.Thumbnail(uint(maxWidth), uint(maxHeight), paletted, resize.Lanczos3)
		// p := image.NewPaletted(r.Bounds(), palette.Plan9)
		// r := imaging.Fit(paletted, maxWidth, maxHeight, imaging.Lanczos)
		r := imaging.Resize(paletted, maxWidth, maxHeight, imaging.Lanczos)
		// r := imaging.Fill(paletted, maxWidth, maxHeight, imaging.Center, imaging.Lanczos)
		p := image.NewPaletted(r.Rect, paletted.Palette)
		draw.Draw(p, p.Rect, r, r.Bounds().Min, draw.Over)
		palettedList = append(palettedList, p)
	}
	gifObj.Image = palettedList
	gifObj.Config.Width = maxWidth
	gifObj.Config.Height = maxHeight

	dstBuf := new(bytes.Buffer)
	if err = gif.EncodeAll(dstBuf, gifObj); err != nil {
		return nil, err
	}
	return dstBuf.Bytes(), nil
}

func Decode(buf []byte, ext string) (image.Image, error) {
	var origin image.Image
	switch strings.ToLower(ext) {
	case "gif":
		img, err := gif.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		origin = img
	case "jpg", "jpeg", "png", "bmp":
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
		origin = img
	default:
		return nil, errors.New("localstorage: image decode error")
	}
	return origin, nil
}

func encode(img image.Image, ext string) ([]byte, error) {
	r := new(bytes.Buffer)
	switch ext {
	case "jpg", "jpeg":
		jpeg.Encode(r, img, nil)
	case "png":
		png.Encode(r, img)
	case "gif":
		gif.Encode(r, img, nil)
	case "bmp":
		bmp.Encode(r, img)
	default:
		return nil, errors.New("localstorage: image encode error")
	}
	return r.Bytes(), nil
}
