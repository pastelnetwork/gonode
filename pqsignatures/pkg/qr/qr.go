// Package qr generates image sequences of QR-codes from the input messages of any length.
package qr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Imports jpeg decoders
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DataDog/zstd"
	"github.com/fogleman/gg"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	// MaxMsgLength defines the maximum character length of a message to be encoded with a single QR-code
	MaxMsgLength = 2953

	// DataQRImageSize defines a preferable size of a QR-code image
	DataQRImageSize = 100

	// PositionVectorQRImageSize defines the size for the QR-code of the position vector
	PositionVectorQRImageSize = 185

	textYPadding = 20
	rowYPadding  = 50
)

// Image contains a QR-code image data
type Image struct {
	raw   []byte
	image image.Image
	title string
	alias string
}

// DecodedMessage represents a single message decoded from the sequence of QR-codes
type DecodedMessage struct {
	Alias   string
	Content string
}

// PositionVectorEncodingError describes errors message if position vector can't encoded with a single QR-code image
var PositionVectorEncodingError = errors.Errorf("Position vector should be encoded as single qr code image")

// CroppingError describes errors message if Image doesn't support cropping
var CroppingError = errors.Errorf("Image interface doesn't support cropping")

// MalformedPositionVector describes errors message for a malformed position vector
var MalformedPositionVector = errors.Errorf("Malformed position vector")

// OutputSizeTooSmall describes errors message if output image size is too small to contain all input images
var OutputSizeTooSmall = errors.Errorf("Output size is too small to map input images")

func compress(src string) (string, error) {
	output, err := zstd.CompressLevel(nil, []byte(src), 22)
	if err != nil {
		return "", errors.New(err)
	}
	return base64.StdEncoding.EncodeToString(output), nil
}

func decompress(src string) (string, error) {
	input, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", errors.New(err)
	}
	uncompressed, err := zstd.Decompress(nil, input)
	if err != nil {
		return "", errors.New(err)
	}
	return string(uncompressed), nil
}

// Encode splits input msg into chunks to fit max supported length of QR code message and generates an array of QR codes images.
func Encode(msg string, alias string, outputDir string, outputFileTitle string, outputFileNamePattern string, outputFileNameSuffix string) ([]Image, error) {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return nil, errors.New(err)
		}
	}
	pngs, err := toPngs(msg, MaxMsgLength, DataQRImageSize)
	if err != nil {
		return nil, err
	}
	var images []Image
	for i, imageBytes := range pngs {
		filePathPartNumber := ""
		titlePartNumber := ""
		if len(pngs) > 1 {
			filePathPartNumber = fmt.Sprintf("__part_%v_of_%v", i+1, len(pngs))
			titlePartNumber = fmt.Sprintf(" part %v of %v", i+1, len(pngs))
		}
		partTitle := fmt.Sprintf("%v%v", outputFileTitle, titlePartNumber)
		filePath := filepath.Join(outputDir, fmt.Sprintf("%v%v%v.png", outputFileNamePattern, filePathPartNumber, outputFileNameSuffix))
		err := os.WriteFile(filePath, imageBytes, 0644)
		if err != nil {
			return nil, errors.New(err)
		}
		img, err := gg.LoadImage(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		images = append(images, Image{
			raw:   imageBytes,
			title: partTitle,
			image: img,
			alias: alias,
		})
	}
	return images, nil
}

// LoadImages loads all QR-Code images under input pattern file path
func LoadImages(pattern string) ([]Image, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errors.New(err)
	}
	var images []Image
	for _, filePath := range matches {
		img, err := gg.LoadImage(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			return nil, errors.New(err)
		}
		imageBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, errors.New(err)
		}

		images = append(images, Image{
			raw:   imageBytes,
			title: fileInfo.Name(),
			image: img,
		})
	}
	return images, nil
}

func maxQRCodeSize(images []Image) image.Point {
	maxQRCodeImageSideSizeX := 0
	maxQRCodeImageSideSizeY := 0
	for _, qrImage := range images {
		x := qrImage.image.Bounds().Size().X
		y := qrImage.image.Bounds().Size().Y
		if maxQRCodeImageSideSizeX < x {
			maxQRCodeImageSideSizeX = x
		}
		if maxQRCodeImageSideSizeY < y {
			maxQRCodeImageSideSizeY = y
		}
	}
	return image.Point{maxQRCodeImageSideSizeX, maxQRCodeImageSideSizeY + textYPadding + rowYPadding}
}

// ImagesFitOutputSize verifies if containerSize has enough capacity to accommodate all QR-codes images
func ImagesFitOutputSize(images []Image, containerSize image.Point) error {
	maxQRCodeSize := maxQRCodeSize(images)

	capacityX := containerSize.X / maxQRCodeSize.X
	capacityY := containerSize.Y / maxQRCodeSize.Y

	if capacityX < 1 || capacityY < 1 {
		return OutputSizeTooSmall
	}

	// Subtract 1 to reserve it for position vector
	totalCapacity := capacityX*capacityY - 1

	if totalCapacity < len(images) {
		return OutputSizeTooSmall
	}

	return nil
}

// MapImages maps input images into output container image of specified size.
func MapImages(images []Image, containerSize image.Point, outputFilePath string) error {
	if err := ImagesFitOutputSize(images, containerSize); err != nil {
		return err
	}
	dc := gg.NewContext(containerSize.X, containerSize.Y)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	dc.SetColor(color.White)
	paddingPixels := 2
	currentX := 0
	currentY := 10
	positionVector := fmt.Sprintf("%v;", len(images))
	currentImageAlias := ""
	for _, image := range images {
		size := image.image.Bounds().Size()
		captionX := size.X / 2

		if currentX+size.X > containerSize.X {
			currentX = 0
			currentY += size.Y + rowYPadding
		}

		dc.DrawStringAnchored(image.title, float64(currentX+captionX), float64(currentY), 0.5, 0.5)
		dc.DrawImageAnchored(image.image, currentX+captionX, currentY+textYPadding+size.Y/2, 0.5, 0.5)

		if currentImageAlias != image.alias {
			currentImageAlias = image.alias
			positionVector += fmt.Sprintf("%v:", image.alias)
		}

		positionVector += fmt.Sprintf("%v,%v,%v;", currentX, currentY+textYPadding, size.X)
		currentX += size.X + paddingPixels
	}
	compressedPositionVector, err := compress(positionVector)
	if err != nil {
		return err
	}
	positionVectorPngsBytes, err := toPngs(compressedPositionVector, MaxMsgLength, PositionVectorQRImageSize)
	if err != nil {
		return err
	}

	if len(positionVectorPngsBytes) != 1 {
		return errors.New(PositionVectorEncodingError)
	}

	positionVectorImage, _, err := image.Decode(bytes.NewReader(positionVectorPngsBytes[0]))
	if err != nil {
		return errors.New(err)
	}

	dc.DrawImageAnchored(positionVectorImage, containerSize.X-positionVectorImage.Bounds().Size().X/2, containerSize.Y-positionVectorImage.Bounds().Size().Y/2, 0.5, 0.5)

	err = dc.SavePNG(outputFilePath)
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func breakStringIntoChunks(str string, size int) []string {
	if len(str) <= size {
		return []string{str}
	}
	var result []string
	chunk := make([]rune, size)
	pos := 0
	for _, rune := range str {
		chunk[pos] = rune
		pos++
		if pos == size {
			pos = 0
			result = append(result, string(chunk))
		}
	}
	if pos > 0 {
		result = append(result, string(chunk[:pos]))
	}
	return result
}

func toPngs(s string, dataSize int, qrImageSize int) ([][]byte, error) {
	chunks := breakStringIntoChunks(s, dataSize)
	var qrs [][]byte
	for _, chunk := range chunks {
		enc := qrcode.NewQRCodeWriter()
		image, err := enc.Encode(chunk, gozxing.BarcodeFormat_QR_CODE, qrImageSize, qrImageSize, nil)
		if err != nil {
			return nil, errors.New(err)
		}
		pngBuffer := new(bytes.Buffer)
		err = png.Encode(pngBuffer, image)
		if err != nil {
			return nil, errors.New(err)
		}
		qrs = append(qrs, pngBuffer.Bytes())
	}
	return qrs, nil
}

func cropImage(img image.Image, rect image.Rectangle) (image.Image, error) {
	type cropper interface {
		SubImage(r image.Rectangle) image.Image
	}
	cropperImg, ok := img.(cropper)
	if !ok {
		return nil, errors.New(CroppingError)
	}
	return cropperImg.SubImage(rect), nil
}

func decodeImage(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", errors.New(err)
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", errors.New(err)
	}

	return result.GetText(), nil
}

// Decode decodes QR-codes sequences into messages.
func Decode(signatureLayerFilePath string) ([]DecodedMessage, error) {
	signatureLayerImage, err := gg.LoadImage(signatureLayerFilePath)
	if err != nil {
		return nil, errors.New(err)
	}

	signatureRect := signatureLayerImage.Bounds()
	positionVectorImageRectMin := image.Point{signatureRect.Max.X - PositionVectorQRImageSize, signatureRect.Max.Y - PositionVectorQRImageSize}
	positionVectorImage, err := cropImage(signatureLayerImage, image.Rectangle{positionVectorImageRectMin, signatureRect.Max})
	if err != nil {
		return nil, err
	}
	positionVectorCompressed, err := decodeImage(positionVectorImage)
	if err != nil {
		return nil, err
	}
	positionVector, err := decompress(positionVectorCompressed)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n" + positionVector)
	tokens := strings.Split(positionVector, ";")
	if len(tokens) <= 1 {
		return nil, errors.New(MalformedPositionVector)
	}
	tokenCount, err := strconv.Atoi(tokens[0])
	if err != nil {
		return nil, errors.New(err)
	}
	if tokenCount != len(tokens)-2 {
		return nil, errors.New(MalformedPositionVector)
	}

	var decoded []DecodedMessage
	currentAlias := ""
	for i := 1; i < len(tokens)-1; i++ {
		token := tokens[i]
		idx := strings.IndexByte(token, ':')
		if idx >= 0 {
			currentAlias = token[:idx]
			decoded = append(decoded, DecodedMessage{
				Alias: currentAlias,
			})
		}
		if currentAlias == "" {
			return nil, errors.New(MalformedPositionVector)
		}
		imageRectValues := strings.Split(token[idx+1:], ",")
		if len(imageRectValues) != 3 {
			return nil, errors.New(MalformedPositionVector)
		}
		rectX, err := strconv.Atoi(imageRectValues[0])
		if err != nil {
			return nil, errors.New(err)
		}
		rectY, err := strconv.Atoi(imageRectValues[1])
		if err != nil {
			return nil, errors.New(err)
		}
		rectSize, err := strconv.Atoi(imageRectValues[2])
		if err != nil {
			return nil, errors.New(err)
		}
		qrImageRect := image.Rectangle{image.Point{rectX, rectY}, image.Point{rectX + rectSize, rectY + rectSize}}
		qrImage, err := cropImage(signatureLayerImage, qrImageRect)
		if err != nil {
			return nil, err
		}

		fd, err := os.Create("last-decoded.png")
		if err != nil {
			return nil, errors.New(err)
		}
		defer fd.Close()
		png.Encode(fd, qrImage)

		decodedString, err := decodeImage(qrImage)
		if err != nil {
			return nil, err
		}
		decoded[len(decoded)-1].Content += decodedString
		fmt.Printf("\n" + decodedString)
	}

	return decoded, nil
}
