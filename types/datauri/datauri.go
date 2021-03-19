package datauri

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
)

//ErrInvalidDataURI invalid DataURI string
var ErrInvalidDataURI = errors.New("datauri: invalid DataURI")

//ErrUnsupportedEncoding encoding is unsupported
var ErrUnsupportedEncoding = errors.New("datauri: unsupported encoding")

//DecodeDataURI decode DataURI into image bytes
func DecodeDataURI(ctx context.Context, dataURI string) ([]byte, string, error) {

	ext, data, err := extractDataURI(dataURI)

	if err != nil {
		return nil, "", err
	}

	e := createEncoding(ext)
	if e == nil {
		return nil, "", ErrUnsupportedEncoding
	}

	dataBuf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, "", err
	}

	mediaType, buf := e.Decode(dataBuf)
	if buf == nil || len(buf) == 0 {
		return nil, "", ErrInvalidDataURI
	}

	return buf, mediaType, nil

}

func extractDataURI(dataURI string) (string, string, error) {
	items := strings.Split(dataURI, ";base64,")

	if len(items) == 2 {
		mediaType := strings.TrimPrefix(items[0], "data:")
		if strings.HasPrefix(mediaType, "image/") {
			return strings.Replace(mediaType, "image/", "", -1), items[1], nil
		}
		if strings.HasPrefix(mediaType, "audio/") {
			return strings.Replace(mediaType, "audio/", "", -1), items[1], nil
		}
	}
	return "", "", ErrInvalidDataURI
}

func extractDataFromDataURI(dataURI string) string {
	s := strings.Split(dataURI, ";base64,")
	if len(s) > 1 {
		return s[1]
	}
	return ""
}
