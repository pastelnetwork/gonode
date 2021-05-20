package pqsignatures

import (
	"bufio"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/xlzd/gotp"
	"golang.org/x/image/font/inconsolata"
)

const (
	// OTPSecretLength hints the length of the OTP secret
	OTPSecretLength       = 16
	otpSecretCharacterSet = "ABCDEF1234567890"
	otpURIQRCodeSize      = 256
)

// SetupOTPAuthenticator configures a one-time password authentication.
func SetupOTPAuthenticator(userEmail, otpSecretFilePath, qrCodeFilePath string) error {
	secret := gotp.RandomSecret(OTPSecretLength)
	fmt.Printf("\nThis is you Google Authenticator Secret:%v", secret)

	err := os.WriteFile(otpSecretFilePath, []byte(secret), 0644)
	if err != nil {
		return errors.New(err)
	}

	googleAuthURI := gotp.NewDefaultTOTP(secret).ProvisioningUri(userEmail, "pastel")

	enc := qrcode.NewQRCodeWriter()
	im, err := enc.Encode(googleAuthURI, gozxing.BarcodeFormat_QR_CODE, otpURIQRCodeSize, otpURIQRCodeSize, nil)
	if err != nil {
		return errors.New(err)
	}

	const W = 1200
	const H = 800
	dc := gg.NewContext(W, H)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	textColor := color.White

	dc.SetFontFace(inconsolata.Regular8x16)

	warningMessage := "You should take a picture of this screen on your phone, but make sure your camera roll is secure first!\nYou can also write down your Google Auth URI string (shown below) as a backup, which will allow you to regenerate the QR code image.\n\n"
	textRightMargin := 60.0
	textTopMargin := 90.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin

	dc.SetColor(textColor)
	dc.DrawStringWrapped(warningMessage, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.SetFontFace(inconsolata.Bold8x16)

	dc.DrawStringWrapped(googleAuthURI, x, y+90, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.DrawImageAnchored(im, W/2, H/2, 0.5, 0.5)

	dc.SavePNG(qrCodeFilePath)
	open.Start(qrCodeFilePath)

	return nil
}

func conformsCharacterSet(value string, characterSet string) bool {
	for _, rune := range value {
		if !strings.ContainsRune(characterSet, rune) {
			return false
		}
	}
	return true
}

func generateCurrentOtpString(otpSecretFilePath string) string {
	otpSecretFileData, err := ioutil.ReadFile(otpSecretFilePath)
	if err != nil {
		return ""
	}
	otpSecret := string(otpSecretFileData)
	if !conformsCharacterSet(otpSecret, otpSecretCharacterSet) {
		return ""
	}
	return gotp.NewDefaultTOTP(otpSecret).Now()
}

func generateCurrentOtpStringFromUserInput() string {
	fmt.Println("\n\nEnter your Google Authenticator Secret in all upper case and numbers:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	otpSecret := scanner.Text()

	if len(otpSecret) != OTPSecretLength {
		return ""
	}
	if !conformsCharacterSet(otpSecret, otpSecretCharacterSet) {
		return ""
	}
	return gotp.NewDefaultTOTP(otpSecret).Now()
}
