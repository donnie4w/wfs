package image

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_imagetype(t *testing.T) {
	bs, _ := ioutil.ReadFile(`C:\Users\dong\Desktop\wfs\1.jpg`)
	fmt.Println(imageType(bs))
}

func _Test_resize(t *testing.T) {
	bs, _ := ioutil.ReadFile(`C:\Users\dong\Desktop\wfs\1.png`)
	dest, err := Resize(bs, 100, 1000, Mode5)
	if err == nil {
		ioutil.WriteFile(`C:\Users\dong\Desktop\wfs\temp.png`, dest, 0644)
	} else {
		fmt.Println(err.Error())
	}
}
