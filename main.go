package main

import (
	"bufio"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"

	"fmt"

	"github.com/PastelNetwork/pqSignatures/legroast"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"

	"github.com/skratchdot/open-golang/open"

	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	"golang.org/x/image/font/inconsolata"
)

func set_up_google_authenticator_for_private_key_encryption_func() {
	secretLength := 16
	secret := gotp.RandomSecret(secretLength)
	fmt.Printf("\nThis is you Google Authenticor Secret:%v", secret)

	err := os.WriteFile("otp_secret.txt", []byte(secret), 0644)
	if err != nil {
		log.Fatal(err)
	}

	google_auth_uri := gotp.NewDefaultTOTP(secret).ProvisioningUri("user@user.com", "pastel")

	err = qrcode.WriteFile(google_auth_uri, qrcode.Medium, 256, "qr.png")

	const W = 1200
	const H = 800
	im, err := gg.LoadImage("qr.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(W, H)
	dc.SetRGB(255, 255, 255)
	dc.Clear()
	textColor := color.White

	dc.SetFontFace(inconsolata.Regular8x16)

	warning_message := "You should take a picture of this screen on your phone, but make sure your camera roll is secure first!\nYou can also write down your Google Auth URI string (shown below) as a backup, which will allow you to regenerate the QR code image.\n\n"
	textRightMargin := 60.0
	textTopMargin := 90.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin

	dc.SetColor(textColor)
	dc.DrawStringWrapped(warning_message, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.SetFontFace(inconsolata.Bold8x16)

	dc.DrawStringWrapped(google_auth_uri, x, y+90, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.DrawImageAnchored(im, W/2, H/2, 0.5, 0.5)
	dc.SavePNG("Google_Authenticator_QR_Code.png")

	open.Start("Google_Authenticator_QR_Code.png")
}

func generate_and_store_key_for_nacl_box_func() {
	box_key := nacl.NewKey()
	box_key_base64 := base64.StdEncoding.EncodeToString(box_key[:])
	err := os.WriteFile("box_key.bin", []byte(box_key_base64), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nThis is the key for encrypting the pastel ID private key (using NACL box) in Base64: %v", box_key)
	fmt.Printf("\nThe key has been saved as a file in the working directory. You should also write this key down as a backup.")
}

func get_image_hash_from_image_file_path_func(sample_image_file_path string) (string, error) {
	f, err := os.Open(sample_image_file_path)
	if err != nil {
		return "", err
	}

	defer f.Close()
	hash := sha3.New256()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func pastel_id_keypair_generation_func() (string, string) {
	fmt.Println("Generating LegRoast keypair now...")
	pk, sk := legroast.Keygen()
	pastel_id_private_key_b16_encoded := base64.StdEncoding.EncodeToString(sk)
	pastel_id_public_key_b16_encoded := base64.StdEncoding.EncodeToString(pk)
	return pastel_id_private_key_b16_encoded, pastel_id_public_key_b16_encoded
}

func get_nacl_box_key_from_file_func(box_key_file_path string) []byte {
	box_key_base64, err := ioutil.ReadFile(box_key_file_path)
	if err != nil {
		panic(err)
	}
	//assert(len(box_key_base64)==44)
	box_key, _ := base64.StdEncoding.DecodeString(string(box_key_base64))
	return box_key
}

func write_pastel_public_and_private_key_to_file_func(pastel_id_public_key_b16_encoded string, pastel_id_private_key_b16_encoded string, box_key_file_path string) {
	pastel_id_public_key_export_format := "-----BEGIN LEGROAST PUBLIC KEY-----\n" + pastel_id_public_key_b16_encoded + "\n-----END LEGROAST PUBLIC KEY-----"
	pastel_id_private_key_export_format := "-----BEGIN LEGROAST PRIVATE KEY-----\n" + pastel_id_private_key_b16_encoded + "\n-----END LEGROAST PRIVATE KEY-----"
	box_key := get_nacl_box_key_from_file_func(box_key_file_path)
	var key [32]byte
	copy(key[:], box_key)
	encrypted := secretbox.EasySeal(([]byte)(pastel_id_private_key_export_format[:]), &key)

	key_file_path := "pastel_id_key_files"
	if _, err := os.Stat(key_file_path); os.IsNotExist(err) {
		os.MkdirAll(key_file_path, 0770)
	}
	err := os.WriteFile(key_file_path+"/pastel_id_legroast_public_key.pem", []byte(pastel_id_public_key_export_format), 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(key_file_path+"/pastel_id_legroast_private_key.pem", []byte(encrypted), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func generate_qr_codes_from_pastel_keypair_func(pastel_id_public_key_b16_encoded string, pastel_id_private_key_b16_encoded string) {

	key_file_path := "pastel_id_key_files"
	if _, err := os.Stat(key_file_path); os.IsNotExist(err) {
		os.MkdirAll(key_file_path, 0770)
	}
	err := qrcode.WriteFile(pastel_id_public_key_b16_encoded, qrcode.Medium, 3000, key_file_path+"/pastel_id_legroast_public_key_qr_code.png")
	if err != nil {
		panic(err)
	}
	err = qrcode.WriteFile(pastel_id_private_key_b16_encoded, qrcode.Medium, 1000, key_file_path+"/pastel_id_legroast_private_key_qr_code.png")
	if err != nil {
		panic(err)
	}
}

func pastel_id_write_signature_on_data_func(input_data_or_string string, pastel_id_private_key_b16_encoded string, pastel_id_public_key_b16_encoded string) string {
	fmt.Printf("\nGenerating LegRoast signature now...")
	//with MyTimer():

	pastel_id_private_key, _ := base64.StdEncoding.DecodeString(pastel_id_private_key_b16_encoded)
	pastel_id_public_key, _ := base64.StdEncoding.DecodeString(pastel_id_public_key_b16_encoded)
	//sleep(0.1*random.random()) #To combat side-channel attacks
	pastel_id_signature := legroast.Sign(pastel_id_public_key, pastel_id_private_key, ([]byte)(input_data_or_string[:]))
	pastel_id_signature_b16_encoded := base64.StdEncoding.EncodeToString(pastel_id_signature)
	//sleep(0.1*random.random())
	return pastel_id_signature_b16_encoded
}

func pastel_id_verify_signature_with_public_key_func(input_data_or_string string, pastel_id_signature_b16_encoded string, pastel_id_public_key_b16_encoded string) int {
	fmt.Printf("\nVerifying LegRoast signature now...")
	//with MyTimer():

	pastel_id_signature, _ := base64.StdEncoding.DecodeString(pastel_id_signature_b16_encoded)
	pastel_id_public_key, _ := base64.StdEncoding.DecodeString(pastel_id_public_key_b16_encoded)
	//sleep(0.1*random.random())
	verified := legroast.Verify(pastel_id_public_key, ([]byte)(input_data_or_string[:]), pastel_id_signature)
	//sleep(0.1*random.random())
	if verified > 0 {
		fmt.Printf("\nSignature is valid!")
	} else {
		fmt.Printf("\nWarning! Signature was NOT valid!")
	}
	return verified
}

func conformsCharacterSet(value string, characterSet string) bool {
	for _, rune := range characterSet {
		if strings.ContainsRune(characterSet, rune) == false {
			return false
		}
	}
	return true
}

func generate_current_otp_string_func() string {
	otp_secret := os.Getenv("PASTEL_OTP_SECRET")
	if len(otp_secret) == 0 {
		otp_secret_file_data, err := ioutil.ReadFile("otp_secret.txt")
		if err != nil {
			return ""
		}
		otp_secret = string(otp_secret_file_data)
	}
	otp_secret_character_set := "ABCDEF1234567890"
	if conformsCharacterSet(otp_secret, otp_secret_character_set) == false {
		return ""
	}
	return gotp.NewDefaultTOTP(otp_secret).Now()
}

func generate_current_otp_string_from_user_input_func() string {
	fmt.Println("\n\nEnter your Google Authenticator Secret in all upper case and numbers:\n\n")
	var otp_secret string
	fmt.Scanln(&otp_secret)
	if len(otp_secret) != 16 {
		return ""
	}
	otp_secret_character_set := "ABCDEF1234567890"
	if conformsCharacterSet(otp_secret, otp_secret_character_set) == false {
		return ""
	}
	return gotp.NewDefaultTOTP(otp_secret).Now()
}

func import_pastel_public_and_private_keys_from_pem_files_func(box_key_file_path string) (string, string) {
	key_file_path := "pastel_id_key_files"
	if info, err := os.Stat(key_file_path); os.IsNotExist(err) || !info.IsDir() {
		fmt.Printf("\nCan't find key storage directory, trying to use current working directory instead!")
		key_file_path = "."
	}
	pastel_id_public_key_pem_filepath := key_file_path + "/pastel_id_legroast_public_key.pem"
	pastel_id_private_key_pem_filepath := key_file_path + "/pastel_id_legroast_private_key.pem"

	pastel_id_public_key_b16_encoded := ""
	pastel_id_private_key_b16_encoded := ""
	pastel_id_public_key_export_format := ""
	pastel_id_private_key_export_format__encrypted := ""
	otp_correct := false

	infoPK, errPK := os.Stat(pastel_id_public_key_pem_filepath)
	infoSK, errSK := os.Stat(pastel_id_private_key_pem_filepath)
	if !os.IsNotExist(errPK) && !infoPK.IsDir() && !os.IsNotExist(errSK) && !infoSK.IsDir() {
		pastel_id_public_key_export_data, errPK := ioutil.ReadFile(pastel_id_public_key_pem_filepath)
		pastel_id_private_key_export_data_encrypted, errSK := ioutil.ReadFile(pastel_id_private_key_pem_filepath)
		if errPK == nil && errSK == nil {
			pastel_id_public_key_export_format = string(pastel_id_public_key_export_data)
			pastel_id_private_key_export_format__encrypted = string(pastel_id_private_key_export_data_encrypted)

			otp_string := generate_current_otp_string_func()
			if otp_string == "" {
				otp_string = generate_current_otp_string_from_user_input_func()
			}
			fmt.Println("\n\nPlease Enter your pastel Google Authenticator Code:")
			fmt.Println(otp_string)

			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			otp_from_user_input := scanner.Text()

			if len(otp_from_user_input) != 6 {
				return "", ""
			}
			otp_correct = (otp_from_user_input == otp_string)
		}
	}

	if otp_correct {
		box_key := get_nacl_box_key_from_file_func(box_key_file_path)
		var key [32]byte
		copy(key[:], box_key)
		pastel_id_private_key_export_format, _ := secretbox.EasyOpen(([]byte)(pastel_id_private_key_export_format__encrypted[:]), &key)
		pastel_id_private_key := strings.ReplaceAll(string(pastel_id_private_key_export_format), "-----BEGIN LEGROAST PRIVATE KEY-----\n", "")
		pastel_id_private_key_b16_encoded = strings.ReplaceAll(string(pastel_id_private_key), "\n-----END LEGROAST PRIVATE KEY-----", "")

		pastel_id_public_key_export_format = strings.ReplaceAll(pastel_id_public_key_export_format, "-----BEGIN LEGROAST PUBLIC KEY-----\n", "")
		pastel_id_public_key_b16_encoded = strings.ReplaceAll(pastel_id_public_key_export_format, "\n-----END LEGROAST PUBLIC KEY-----", "")
	}
	return pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded
}

func main() {
	if _, err := os.Stat("otp_secret.txt"); os.IsNotExist(err) {
		set_up_google_authenticator_for_private_key_encryption_func()
	}

	box_key_file_path := "box_key.bin"
	if _, err := os.Stat(box_key_file_path); os.IsNotExist(err) {
		generate_and_store_key_for_nacl_box_func()
	}

	sample_image_file_path := "sample_image2.png"

	fmt.Printf("\nApplying signature to file %v", sample_image_file_path)
	sha256_hash_of_image_to_sign, err := get_image_hash_from_image_file_path_func(sample_image_file_path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSHA256 Hash of Image File: %v", sha256_hash_of_image_to_sign)
	//sample_input_data_to_be_signed := sha256_hash_of_image_to_sign.encode('utf-8')

	pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded := import_pastel_public_and_private_keys_from_pem_files_func(box_key_file_path)
	if pastel_id_public_key_b16_encoded == "" {
		pastel_id_private_key_b16_encoded, pastel_id_public_key_b16_encoded = pastel_id_keypair_generation_func()
		write_pastel_public_and_private_key_to_file_func(pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded, box_key_file_path)
		generate_qr_codes_from_pastel_keypair_func(pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded)
	}

	pastel_id_signature_b16_encoded := pastel_id_write_signature_on_data_func(sha256_hash_of_image_to_sign, pastel_id_private_key_b16_encoded, pastel_id_public_key_b16_encoded)
	_ = pastel_id_verify_signature_with_public_key_func(sha256_hash_of_image_to_sign, pastel_id_signature_b16_encoded, pastel_id_public_key_b16_encoded)
}
