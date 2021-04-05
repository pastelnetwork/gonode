package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fogleman/gg"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/skratchdot/open-golang/open"
	"github.com/xlzd/gotp"
	"golang.org/x/image/font/inconsolata"
)

const (
	OTPSecretFile   = "otp_secret.txt"
	OTPSecretLength = 16
	OTPEnvVarName   = "PASTEL_OTP_SECRET"
)

func setupGoogleAuthenticatorForPrivateKeyEncryption() error {
	secret := gotp.RandomSecret(OTPSecretLength)
	fmt.Printf("\nThis is you Google Authenticor Secret:%v", secret)

	err := os.WriteFile(OTPSecretFile, []byte(secret), 0644)
	if err != nil {
		return fmt.Errorf("setupGoogleAuthenticatorForPrivateKeyEncryption: %w", err)
	}

	googleAuthUri := gotp.NewDefaultTOTP(secret).ProvisioningUri("user@user.com", "pastel")
	const qrCodeFileName = "qr.png"
	err = qrcode.WriteFile(googleAuthUri, qrcode.Medium, 256, qrCodeFileName)
	if err != nil {
		return fmt.Errorf("setupGoogleAuthenticatorForPrivateKeyEncryption: %w", err)
	}

	const W = 1200
	const H = 800
	im, err := gg.LoadImage(qrCodeFileName)
	if err != nil {
		return fmt.Errorf("setupGoogleAuthenticatorForPrivateKeyEncryption: %w", err)
	}

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

	dc.DrawStringWrapped(googleAuthUri, x, y+90, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.DrawImageAnchored(im, W/2, H/2, 0.5, 0.5)

	const resultFileName = "Google_Authenticator_QR_Code.png"
	dc.SavePNG(resultFileName)
	open.Start(resultFileName)

	return nil
}

func conformsCharacterSet(value string, characterSet string) bool {
	for _, rune := range characterSet {
		if !strings.ContainsRune(characterSet, rune) {
			return false
		}
	}
	return true
}

func generateCurrentOtpString() string {
	otp_secret := os.Getenv(OTPEnvVarName)
	if len(otp_secret) == 0 {
		otp_secret_file_data, err := ioutil.ReadFile(OTPSecretFile)
		if err != nil {
			return ""
		}
		otp_secret = string(otp_secret_file_data)
	}
	otp_secret_character_set := "ABCDEF1234567890"
	if !conformsCharacterSet(otp_secret, otp_secret_character_set) {
		return ""
	}
	return gotp.NewDefaultTOTP(otp_secret).Now()
}

func generateCurrentOtpStringFromUserInput() string {
	fmt.Println("\n\nEnter your Google Authenticator Secret in all upper case and numbers:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	otp_secret := scanner.Text()

	if len(otp_secret) != OTPSecretLength {
		return ""
	}
	otp_secret_character_set := "ABCDEF1234567890"
	if !conformsCharacterSet(otp_secret, otp_secret_character_set) {
		return ""
	}
	return gotp.NewDefaultTOTP(otp_secret).Now()
}
