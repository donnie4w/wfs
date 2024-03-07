// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package tc

import (
	"strconv"
	"strings"

	"github.com/donnie4w/gofer/image"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

type imageMode struct {
	mode   int
	width  int
	height int
	gray   bool
	invert bool
	format string
	rotate int
	fliph  bool
	flipv  bool
	cropA  []int
	cropS  []int
	blur   float64
	scaleU []int
	scaleL []int
}

func (i *imageMode) getOptions() *image.Options {
	return &image.Options{Gray: i.gray, Invert: i.invert, Format: i.format, Rotate: i.rotate, FlipH: i.fliph, FlipV: i.flipv, Blur: i.blur, CropAnchor: i.cropA, CropSide: i.cropS, ScaleUpper: i.scaleU, ScaleLower: i.scaleL}
}

func parseUriToImagemode(uri string) (iv2 *imageMode, err error) {
	pairs := strings.Split(uri, "/")
	length := len(pairs)
	iv2 = &imageMode{}
	for i := 0; i < length; i += 2 {
		if i+1 >= length {
			break
		}
		key := pairs[i]
		valueStr := pairs[i+1]
		switch key {
		case sys.IMAGEMODE, sys.IMAGEVIEW, sys.IMAGEVIEW2:
			iv2.mode, _ = strconv.Atoi(valueStr)
		case "w", "width":
			iv2.width, _ = parseSide(valueStr)
		case "h", "height":
			iv2.height, _ = parseSide(valueStr)
		case "gray", "grey", "g":
			if value, _ := strconv.Atoi(valueStr); value == 1 {
				iv2.gray = true
			}
		case "i", "invert":
			if value, _ := strconv.Atoi(valueStr); value == 1 {
				iv2.invert = true
			}
		case "f", "format":
			iv2.format = valueStr
		case "r", "rotate":
			iv2.rotate, _ = parseSide(valueStr)
		case "fliph":
			if value, _ := strconv.Atoi(valueStr); value == 1 {
				iv2.fliph = true
			}
		case "flipv":
			if value, _ := strconv.Atoi(valueStr); value == 1 {
				iv2.flipv = true
			}
		case "c", "crop":
			if t, crop := praseCrop(valueStr); crop != nil {
				if t == 0 {
					iv2.cropA = crop
				} else if t == 1 {
					iv2.cropS = crop
				}
			}
		case "b", "blur":
			iv2.blur, _ = parseSigma(valueStr)
		case "s", "scale":
			if t, scale := praseScale(valueStr); scale != nil {
				if t == 0 {
					iv2.scaleU = scale
				} else if t == 1 {
					iv2.scaleL = scale
				}
			}
		}
	}
	return
}

func parseSide(valueStr string) (i int, err error) {
	if i, err = strconv.Atoi(valueStr); i > sys.MaxSide {
		i = sys.MaxSide
	}
	return
}

func parseSigma(valueStr string) (i float64, err error) {
	if i, err = strconv.ParseFloat(valueStr, 64); i > sys.MaxSigma {
		i = sys.MaxSigma
	}
	return
}

func praseCrop(valueStr string) (t int, crop []int) {
	defer util.Recover()
	if valueStr == "" {
		return
	}
	if valueStr[0] == 's' {
		t = 1
	}
	if valueStr[0] == 's' || valueStr[0] == 'a' {
		valueStr = valueStr[1:]
	}
	parts := strings.Split(valueStr, "-")
	dims := strings.Split(parts[0], "x")
	crop = make([]int, 4)
	if len(parts) == 1 {
		crop[0], _ = parseSide(dims[0])
		if len(dims) > 1 {
			crop[1], _ = parseSide(dims[1])
		}
	} else if len(parts) == 2 {
		crop[0], _ = parseSide(dims[0])
		if len(dims) > 1 {
			crop[1], _ = parseSide(dims[1])
		}
		crop[2], _ = parseSide(parts[1])
	} else if len(parts) == 3 {
		crop[0], _ = parseSide(dims[0])
		if len(dims) > 1 {
			crop[1], _ = parseSide(dims[1])
		}
		crop[2], _ = parseSide(parts[1])
		crop[3], _ = parseSide(parts[2])
	}
	if crop[0] == 0 && crop[1] == 0 && crop[2] == 0 && crop[3] == 0 {
		return t, nil
	}
	return
}

func praseScale(valueStr string) (t int, scale []int) {
	defer util.Recover()
	if valueStr == "" {
		return
	}
	if valueStr[0] == 's' {
		t = 1
		valueStr = valueStr[1:]
	}
	dims := strings.Split(valueStr, "x")
	scale = make([]int, 3)
	if len(dims) == 1 {
		scale[0], _ = parseSide(dims[0])
	} else if len(dims) > 1 {
		scale[0], _ = parseSide(dims[0])
		scale[1], _ = parseSide(dims[1])
	}
	if scale[0] == 0 && scale[1] == 0 {
		return t, nil
	}
	scale[2] = sys.MaxPixel
	return
}
