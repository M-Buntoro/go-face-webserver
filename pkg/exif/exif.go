package exif

import (
	"bytes"
	"encoding/json"
	"image"
	"log"

	"github.com/adityayuga/goexif/exif"
	"github.com/adityayuga/goexif/mknote"
	"gocv.io/x/gocv"
)

type ImageExif struct {
	ApertureValue           []string `json:"ApertureValue,omitempty"`
	BrightnessValue         []string `json:"BrightnessValue,omitempty"`         //: "209/100"
	ColorSpace              []int    `json:"ColorSpace,omitempty"`              //: 1
	DateTime                string   `json:"DateTime,omitempty"`                //: "2018:11:23 15:09:33"
	DateTimeOriginal        string   `json:"DateTimeOriginal,omitempty"`        //: "2018:11:23 15:09:33"
	DateTimeDigitizedstring string   `json:"DateTimeDigitizedstring,omitempty"` //: "2018:11:23 15:09:33"
	ExifVersion             string   `json:"ExifVersion,omitempty"`             //: "0220"
	ExposureProgram         []int    `json:"ExposureProgram,omitempty"`         //: 2
	ExposureMode            []int    `json:"ExposureMode,omitempty"`            //: 0
	ExposureBiasValue       []string `json:"ExposureBiasValue,omitempty"`       //: "0/10"
	Flash                   []int    `json:"Flash,omitempty"`                   // :0
	ISOSpeedRatings         []int    `json:"ISOSpeedRatings ,omitempty"`        //: 100
	ImageUniqueID           string   `json:"ImageUniqueID,omitempty"`           //: "W08LLKF00AM"
	Make                    string   `json:"Make,omitempty"`                    //: "samsung"
	MeteringMode            []int    `json:"MeteringMode,omitempty"`            //: 2
	Model                   string   `json:"Model,omitempty"`                   //: "SM-N960F"
	SceneCaptureType        []int    `json:"SceneCaptureType,omitempty"`        //: 2
	Orientation             []int    `json:"Orientation,omitempty"`             //: 8
	PixelXDimension         []int    `json:"PixelXDimension,omitempty"`         //: 3264
	PixelYDimension         []int    `json:"PixelYDimension,omitempty"`         //: 2448
	ResolutionUnit          []int    `json:"ResolutionUnit,omitempty"`          //: 2
	Software                string   `json:"Software,omitempty"`                //: "N960FXXS2ARJ4"
	WhiteBalance            []int    `json:"WhiteBalance,omitempty"`            //: 0 (Auto),  1(Manual)
}

type ExifModule struct {
}

func New() ExifModule {
	return ExifModule{}
}

// extract exif from bytes content of an image
func (e ExifModule) ExtractExifFromImageBytes(source []byte) (result ImageExif, err error) {
	exif.RegisterParsers(mknote.All...)

	imageReader := bytes.NewReader(source)
	if imageReader == nil {
		log.Println("[Error][EXIF-READER] Error on creating reader, ", err)
		return result, err
	}

	exifdata, err := exif.Decode(imageReader)
	if err != nil {
		log.Println("[Error][EXIF-READER] Error on reading exif, ", err)
		return result, err
	}

	marshaled, err := exifdata.MarshalJSON()
	if err != nil {
		log.Println("[Error][EXIF-READER] Error on marshaling exif, ", err)
		return result, err
	}

	err = json.Unmarshal(marshaled, &result)
	if err != nil {
		log.Println("[Error][EXIF-READER] Error on unmarshaling exif, ", err)
		return result, err
	}

	return result, nil
}

/*
 * AdjustOrientation
 * utilize opencv to correct image rotation from exif data
 * for more information on exif orientation try https://www.impulseadventure.com/photo/exif-orientation.html
 */
func (e ExifModule) AdjustOrientation(img image.Image, exifOrientation int) image.Image {
	// no need rotation
	if exifOrientation <= 1 {
		return img
	}

	imgmat, err := gocv.ImageToMatRGBA(img)
	if err != nil {
		log.Printf("[Error][REORIENTATION] Failed to reorient image, %v\n", err)
		return img
	}
	dest := gocv.NewMat()

	// closing matrix object
	defer func() {
		err = dest.Close()
		if err != nil {
			log.Println("[Error][REORIENTATION] Closing mat destination")
		}

		err = imgmat.Close()
		if err != nil {
			log.Println("[Error][REORIENTATION] Closing mat image")
		}
	}()

	switch exifOrientation {
	case 1:
		// orientation is normal, no need for rotation
		return img

	case 8:
		// image is rotated 90 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate90CounterClockwise)

	case 3:
		// image is rotated 180 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate180Clockwise)

	case 6:
		// image is rotated 270 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate90Clockwise)

	// from here image is horizontally flipped
	case 2:
		// image is flipped horizontally
		gocv.Flip(imgmat, &dest, 1)
	case 4:
		// image is flipped horizontally rotated 180 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate180Clockwise)
		gocv.Flip(dest, &dest, 1)

	case 5:
		// image is flipped horizontally rotated 270 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate90Clockwise)
		gocv.Flip(dest, &dest, 1)

	case 7:
		// image is flipped horizontally rotated 90 degrees
		gocv.Rotate(imgmat, &dest, gocv.Rotate90CounterClockwise)
		gocv.Flip(dest, &dest, 1)
	}

	imrslt, err := dest.ToImage()
	if err != nil {
		log.Printf("[Error][REORIENTATION] Failed to reorient image, %v\n", err)
		return img
	}

	return imrslt
}
