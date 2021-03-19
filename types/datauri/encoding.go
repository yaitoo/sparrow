package datauri

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
)

var encodings map[string]Encoding

func init() {
	encodings = make(map[string]Encoding)
	encodings["png"] = &PngEncoding{}
	encodings["jpeg"] = &JpegEncoding{}
	encodings["jpg"] = &JpegEncoding{}
	encodings["gif"] = &GifEncoding{}
	encodings["bmp"] = &BitmapEncoding{}
	encodings["amr"] = &AmrEncoding{}
}

func createEncoding(ext string) Encoding {

	e, ok := encodings[ext]
	if ok {
		return e
	}

	return nil

}

//Encoding 图片编码器
type Encoding interface {
	Decode(dataURI []byte) (mediaType string, buf []byte)
}

type AmrEncoding struct {
}

//Decode 转换dataURI格式到bytes
func (enc *AmrEncoding) Decode(dataURI []byte) (string, []byte) {

	return "audio/amr", dataURI
}

//JpegEncoding jpeg图片编码器
type JpegEncoding struct {
}

//Decode 转换dataURI格式到bytes
func (enc *JpegEncoding) Decode(dataURI []byte) (string, []byte) {

	img, err := jpeg.Decode(bytes.NewReader(dataURI))

	if err != nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	return "image/jpeg", buf.Bytes()
}

//PngEncoding png图片编码器
type PngEncoding struct {
}

//Decode 转换dataURI格式到bytes
func (enc *PngEncoding) Decode(dataURI []byte) (string, []byte) {
	img, err := png.Decode(bytes.NewReader(dataURI))

	if err != nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	return "image/png", buf.Bytes()
}

//BitmapEncoding bmp图片编码器
type BitmapEncoding struct {
}

//Decode 转换dataURI格式到bytes
func (enc *BitmapEncoding) Decode(dataURI []byte) (string, []byte) {
	img, err := bmp.Decode(bytes.NewReader(dataURI))

	if err != nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	err = bmp.Encode(buf, img)
	return "image/bmp", buf.Bytes()
}

//GifEncoding gif图片编码器
type GifEncoding struct {
}

//Decode 转换dataURI格式到bytes
func (enc *GifEncoding) Decode(dataURI []byte) (string, []byte) {
	img, err := gif.DecodeAll(bytes.NewReader(dataURI))

	if err != nil {
		return "", nil
	}

	buf := new(bytes.Buffer)
	err = gif.EncodeAll(buf, img)
	return "image/gif", buf.Bytes()
}
