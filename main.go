package main

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"

	"fmt"
	"time"

	"github.com/PastelNetwork/pqSignatures/legroast"
	"github.com/PastelNetwork/pqSignatures/qr"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"

	"github.com/skratchdot/open-golang/open"

	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/sha3"

	"golang.org/x/image/font/inconsolata"

	"github.com/auyer/steganography"
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
	if err != nil {
		log.Fatal(err)
	}

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
	err := qr.Encode(pastel_id_public_key_b16_encoded, key_file_path, "pastel_id_legroast_public_key_qr_code")
	if err != nil {
		panic(err)
	}
	err = qr.Encode(pastel_id_private_key_b16_encoded, key_file_path, "pastel_id_legroast_private_key_qr_code")
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
		if !strings.ContainsRune(characterSet, rune) {
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
	if !conformsCharacterSet(otp_secret, otp_secret_character_set) {
		return ""
	}
	return gotp.NewDefaultTOTP(otp_secret).Now()
}

func generate_current_otp_string_from_user_input_func() string {
	fmt.Println("\n\nEnter your Google Authenticator Secret in all upper case and numbers:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	otp_secret := scanner.Text()

	if len(otp_secret) != 16 {
		return ""
	}
	otp_secret_character_set := "ABCDEF1234567890"
	if !conformsCharacterSet(otp_secret, otp_secret_character_set) {
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

func loadImageSizes(list_of_qr_code_file_paths []string) []image.Point {
	var sizes []image.Point
	for _, imagePath := range list_of_qr_code_file_paths {
		image, err := gg.LoadImage(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		sizes = append(sizes, image.Bounds().Size())
	}
	return sizes
}

func totalWidth(sizes []image.Point) int {
	width := 0
	for _, size := range sizes {
		width += size.X
	}
	width += 3 * (len(sizes) - 1)
	return width
}

func maxHeight(sizes []image.Point) int {
	maxHeight := 0
	for _, size := range sizes {
		if maxHeight < size.Y {
			maxHeight = size.Y
		}
	}
	return maxHeight
}

func generate_signature_image_layer_func(pastel_id_public_key_b16_encoded string, pastel_id_signature_b16_encoded string, sample_image_file_path string) string {
	max_chunk_length_for_qr_code := 2200

	list_of_pastel_id_signature_b16_encoded_qr_code, err := qr.ToPngs(pastel_id_signature_b16_encoded, max_chunk_length_for_qr_code)
	number_of_signature_qr_codes := len(list_of_pastel_id_signature_b16_encoded_qr_code)
	current_datetime_string := time.Now().Format("Jan_02_2006_15_04_05")

	pastel_id_signatures_storage_file_path := "pastel_id_signature_files"
	if _, err := os.Stat(pastel_id_signatures_storage_file_path); os.IsNotExist(err) {
		os.MkdirAll(pastel_id_signatures_storage_file_path, 0770)
	}
	for i, imageData := range list_of_pastel_id_signature_b16_encoded_qr_code {
		err := os.WriteFile(pastel_id_signatures_storage_file_path+"/pastel_id_legroast_signature_qr_code__part_"+strconv.Itoa(i)+"_of_"+strconv.Itoa(number_of_signature_qr_codes)+current_datetime_string+".png", []byte(imageData), 0644)
		if err != nil {
			panic(err)
		}
	}

	err = qrcode.WriteFile(pastel_id_public_key_b16_encoded, qrcode.Medium, 250, pastel_id_signatures_storage_file_path+"/pastel_id_legroast_public_key_qr_code"+current_datetime_string+".png")
	if err != nil {
		panic(err)
	}

	sample_image, err := gg.LoadImage(sample_image_file_path)
	if err != nil {
		log.Fatal(err)
	}
	imageSize := sample_image.Bounds().Size()
	list_of_qr_code_file_paths, err := filepath.Glob(pastel_id_signatures_storage_file_path + "/*" + current_datetime_string + "*")
	if err != nil {
		panic(err)
	}

	list_of_qr_code_image_sizes := loadImageSizes(list_of_qr_code_file_paths)

	sum_of_widths_plus_borders := totalWidth(list_of_qr_code_image_sizes)
	max_height_of_qr_codes := maxHeight(list_of_qr_code_image_sizes)

	if max_height_of_qr_codes > imageSize.Y {
		panic(errors.New("image height is too small"))
	}

	if imageSize.X < sum_of_widths_plus_borders {

	}

	signature_layer_image_dc := gg.NewContext(imageSize.X, imageSize.Y)
	signature_layer_image_dc.SetRGB(255, 255, 255)
	signature_layer_image_dc.Clear()
	signature_layer_image_dc.SetColor(color.White)
	padding_pixels := 2
	current_qr_image_x_coordinate := 0
	for i, current_qr_code_image := range list_of_qr_code_file_paths {
		caption_text := ""
		caption_x_position := 0

		if i == 0 {
			caption_text = "Pastel Public Key"
			caption_x_position = list_of_qr_code_image_sizes[i].X / 2
		} else {
			current_qr_image_x_coordinate += list_of_qr_code_image_sizes[i-1].X + padding_pixels
			caption_text = "Pastel Signature Part " + strconv.Itoa(i) + " of " + strconv.Itoa(len(list_of_qr_code_image_sizes)-1)
			caption_x_position = list_of_qr_code_image_sizes[i-1].X / 2
		}

		image, err := gg.LoadImage(current_qr_code_image)
		if err != nil {
			log.Fatal(err)
		}
		signature_layer_image_dc.DrawStringAnchored(caption_text, float64(current_qr_image_x_coordinate+caption_x_position), 10, 0.5, 0.5)
		signature_layer_image_dc.DrawImageAnchored(image, current_qr_image_x_coordinate+caption_x_position, 30+list_of_qr_code_image_sizes[i].Y/2, 0.5, 0.5)
	}

	signature_layer_image_output_filepath := pastel_id_signatures_storage_file_path + "/Complete_Signature_Image_Layer__" + current_datetime_string + ".png"
	signature_layer_image_dc.SavePNG(signature_layer_image_output_filepath)

	return signature_layer_image_output_filepath
}

func hide_signature_image_in_sample_image_func(sample_image_file_path string, signature_layer_image_output_filepath string, signed_image_output_path string) {
	img, err := gg.LoadImage(sample_image_file_path)
	if err != nil {
		log.Fatal(err)
	}

	signature_layer_image_data, err := ioutil.ReadFile(signature_layer_image_output_filepath)
	if err != nil {
		panic(err)
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, signature_layer_image_data)
	if err != nil {
		panic(err)
	}
	outFile, _ := os.Create(signed_image_output_path)
	w.WriteTo(outFile)
	outFile.Close()
}

func extract_signature_image_in_sample_image_func(signed_image_output_path string, extracted_signature_layer_image_output_filepath string) {
	img, err := gg.LoadImage(signed_image_output_path)
	if err != nil {
		panic(err)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)

	decodedData := steganography.Decode(sizeOfMessage, img)
	err = os.WriteFile(extracted_signature_layer_image_output_filepath, decodedData, 0644)
	if err != nil {
		panic(err)
	}
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

	pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded := import_pastel_public_and_private_keys_from_pem_files_func(box_key_file_path)
	if pastel_id_public_key_b16_encoded == "" {
		pastel_id_private_key_b16_encoded, pastel_id_public_key_b16_encoded := pastel_id_keypair_generation_func()
		write_pastel_public_and_private_key_to_file_func(pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded, box_key_file_path)
		generate_qr_codes_from_pastel_keypair_func(pastel_id_public_key_b16_encoded, pastel_id_private_key_b16_encoded)
	}

	pastel_id_signature_b16_encoded := pastel_id_write_signature_on_data_func(sha256_hash_of_image_to_sign, pastel_id_private_key_b16_encoded, pastel_id_public_key_b16_encoded)
	_ = pastel_id_verify_signature_with_public_key_func(sha256_hash_of_image_to_sign, pastel_id_signature_b16_encoded, pastel_id_public_key_b16_encoded)

	signature_layer_image_output_filepath := generate_signature_image_layer_func(pastel_id_public_key_b16_encoded, pastel_id_signature_b16_encoded, sample_image_file_path)
	signed_image_output_path := "final_watermarked_image.png"
	hide_signature_image_in_sample_image_func(sample_image_file_path, signature_layer_image_output_filepath, signed_image_output_path)

	extracted_signature_layer_image_output_filepath := "extracted_signature_image.png"
	extract_signature_image_in_sample_image_func(signed_image_output_path, extracted_signature_layer_image_output_filepath)
}
