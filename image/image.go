/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package image

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type ResizeType int
type Mode int

const (
	SCALE ResizeType = iota
	THUMBNAIL
)

const (
	Mode0 Mode = iota
	Mode1
	Mode2
	Mode3
	Mode4
	Mode5
)

func Resize(srcData []byte, width, height int, mode Mode) (destData []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
		}
	}()
	if width == 0 && height == 0 {
		return srcData, nil
	}
	img, _, er := image.Decode(bytes.NewReader(srcData))
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	nw, nh, resizeType := praseMode(mode, w, h, width, height)
	if er == nil {
		//		fmt.Println("===>", nw, " ", nh, " ", resizeType)
		var nrgba *image.NRGBA
		switch resizeType {
		case SCALE:
			nrgba = imaging.Resize(img, nw, nh, imaging.Lanczos)
		case THUMBNAIL:
			nrgba = imaging.Fill(img, nw, nh, imaging.Center, imaging.Lanczos)
		}
		var buf bytes.Buffer
		switch imageType(srcData) {
		case "jpeg":
			err = jpeg.Encode(&buf, nrgba, nil)
		case "png":
			err = png.Encode(&buf, nrgba)
		case "gif":
			err = gif.Encode(&buf, nrgba, nil)
		case "bmp":
			err = bmp.Encode(&buf, nrgba)
		case "tif":
			err = tiff.Encode(&buf, nrgba, nil)
		default:
			return srcData, nil
		}
		if err == nil {
			return buf.Bytes(), nil
		}
	}
	return srcData, nil
}

func imageType(srcData []byte) (s string) {
	readlen := 8
	length := len(srcData)
	if length < readlen {
		readlen = length
	}
	prx := strings.ToUpper(hex.EncodeToString(srcData[0:readlen]))
	switch {
	case strings.HasPrefix(prx, "FF"):
		s = "jpeg"
	case strings.HasPrefix(prx, "89504E470D0A1A0A"):
		s = "png"
	case strings.HasPrefix(prx, "474946"):
		s = "gif"
	case strings.HasPrefix(prx, "424D"):
		s = "bmp"
	case strings.HasPrefix(prx, "4949") || strings.HasPrefix(prx, "4D4D"):
		s = "tif"
	}
	return
}

func praseMode(mode Mode, w, h, width, height int) (nw, nh int, resizeType ResizeType) {
	if width > w && height > h {
		return w, h, SCALE
	}
	switch mode {
	case Mode0:
		nw, nh = getMin(w, h, width, height, false)
		resizeType = SCALE
	case Mode1:
		nw, nh = getMax(w, h, width, height, true)
		resizeType = THUMBNAIL
	case Mode2:
		nw, nh = getMin(w, h, width, height, false)
		resizeType = SCALE
	case Mode3:
		nw, nh = getMax(w, h, width, height, false)
		resizeType = SCALE
	case Mode4:
		nw, nh = getMax(w, h, width, height, false)
		resizeType = SCALE
	case Mode5:
		nw, nh = getMax(w, h, width, height, true)
		resizeType = THUMBNAIL
	default:
		return w, h, SCALE
	}
	return
}

func getMin(w, h, width, height int, isThubnail bool) (nw, nh int) {
	if width > w && height > h || (width == 0 && height == 0) {
		return w, h
	}
	if isThubnail {
		nw, nh = w, h
		if width < w {
			nw = width
		}
		if height < h {
			nh = height
		}
		if nw == 0 {
			nw = nh
		}
		if nh == 0 {
			nh = nw
		}
		return
	}
	if float32(width)/float32(w) > float32(height)/float32(h) {
		if height > 0 {
			return 0, height
		} else if width < w {
			return width, 0
		}
	} else {
		if width > 0 {
			return width, 0
		} else if height < h {
			return 0, height
		}
	}
	return w, h
}

func getMax(w, h, width, height int, isThubnail bool) (nw, nh int) {
	if width > w && height > h {
		return w, h
	}
	if isThubnail {
		nw, nh = w, h
		if width < w {
			nw = width
		}
		if height < h {
			nh = height
		}
		if nw == 0 {
			nw = nh
		}
		if nh == 0 {
			nh = nw
		}
		return
	}
	if float32(width)/float32(w) > float32(height)/float32(h) {
		if width < w {
			return width, 0
		} else if height > 0 {
			return 0, height
		}
	} else {
		if height < h {
			return 0, height
		} else if width > 0 {
			return width, 0
		}
	}
	return w, h
}
