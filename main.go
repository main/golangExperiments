package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hmage/goexif/exif"
	"github.com/hmage/goexif/mknote"
)

func main() {
	fname := "MVI_8511.MOV"

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	freader, err := reentrantReader(f, func(reader io.Reader) error {
		buf := make([]byte, 512)
		n, err := reader.Read(buf)
		if err != nil {
			return err
		}

		log.Println("!!!!!!!!!!!!!!!!", http.DetectContentType(buf[:n])) //application/octet-stream

		return nil
	})

	if err != nil {
		log.Println("reentrantReader error ", err)
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(freader)
	if err != nil {
		log.Fatal(err)
	}

	ExifIFDPointer, _ := x.Get(exif.ExifIFDPointer) // normally, don't ignore errors!
	fmt.Println(ExifIFDPointer.StringVal())

	LensModel, _ := x.Get(exif.LensModel) // normally, don't ignore errors!
	fmt.Println(LensModel.StringVal())

	camModel, _ := x.Get(exif.Model) // normally, don't ignore errors!
	fmt.Println(camModel.StringVal())

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, _ := focal.Rat2(0) // retrieve first (only) rat. value
	fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, _ := x.DateTime()
	fmt.Println("Taken: ", tm)

	lat, long, _ := x.LatLong()
	fmt.Println("lat, long: ", lat, ", ", long)

}

func reentrantReader(r io.Reader, f func(io.Reader) error) (io.Reader, error) {
	buf := bytes.NewBuffer([]byte{})
	err := f(io.TeeReader(r, buf))
	return io.MultiReader(buf, r), err
}
