/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"strconv"
	"strings"

	"wfs/image"
)

type Spec struct {
	RT     image.ResizeType
	Mode   image.Mode
	Width  int
	Height int
	src    []byte
}

func NewSpec(bs []byte, arg string) (spec *Spec) {
	spec = new(Spec)
	spec.src = bs
	//	fmt.Println("arg===>", arg)
	if strings.HasPrefix(arg, "?imageView2") {
		ss := strings.Split(arg, "/")
		if ss != nil && len(ss) > 3 {
			spec.Mode = GetMode(atoi(ss[1]))
			switch ss[2] {
			case "w":
				spec.Width = atoi(ss[3])
			case "h":
				spec.Height = atoi(ss[3])
			}
			if len(ss) > 5 {
				switch ss[4] {
				case "w":
					spec.Width = atoi(ss[5])
				case "h":
					spec.Height = atoi(ss[5])
				}
			}
		}
	}
	return
}

func (this *Spec) GetData() (bs []byte) {
	bs, _ = image.Resize(this.src, this.Width, this.Height, this.Mode)
	return
}

func GetMode(mode int) image.Mode {
	switch mode {
	case 0:
		return image.Mode0
	case 1:
		return image.Mode1
	case 2:
		return image.Mode2
	case 3:
		return image.Mode3
	case 4:
		return image.Mode4
	case 5:
		return image.Mode5
	default:
		return image.Mode0
	}
}

func atoi(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}
