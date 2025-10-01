# An Investigation into Pastel's Steganography, QR Codes, and Image Fingerprinting

## Introduction

The Pastel Network is a decentralized platform for creating, trading, and collecting non-fungible tokens (NFTs). A key challenge in the NFT space is ensuring the authenticity and originality of the artwork. Pastel addresses this challenge through a sophisticated and multi-layered approach that combines steganography, QR codes, and robust image fingerprinting.

This document provides an extremely detailed and comprehensive explanation of these three core components, how they work together, and how they are implemented in the Pastel GoNode and `dd-service` codebase.

The three pillars of Pastel's NFT registration process are:

1.  **Robust Image Fingerprinting**: This is the process of generating a unique and compact "fingerprint" for an image. Unlike cryptographic hashes, which are sensitive to even the smallest changes in the input data, robust image fingerprints are designed to be similar for visually identical or near-identical images. This allows the Pastel Network to detect duplicate or near-duplicate artworks, even if they have been slightly modified.

2.  **QR Code Signatures**: Pastel uses QR codes to store a wealth of information directly within the artwork file. This includes the robust image fingerprint, the artist's digital signature, and other metadata. By embedding this data in the image itself, Pastel creates a self-contained and verifiable token that is not dependent on external servers or databases.

3.  **Steganography**: This is the art and science of hiding information within other information. In the context of the Pastel Network, steganography is used to embed the QR code signature into the original artwork image in a way that is imperceptible to the human eye. This ensures that the artwork's aesthetic qualities are not compromised, while still allowing the network to extract the hidden data for verification.

Together, these three technologies create a powerful and secure system for NFT registration that is designed to protect artists and collectors from art theft and duplication. This document will now delve into the technical details of each of these components, starting with steganography.

## Part 1: Steganography in Pastel

Steganography is the practice of concealing a file, message, image, or video within another file, message, image, or video. The word steganography comes from Greek steganographia, which combines the words steganos, meaning "covered or concealed", and -graphia meaning "writing".

The advantage of steganography over cryptography is that the intended secret message does not attract attention to itself as an object of scrutiny. Plainly visible encrypted messages, no matter how unbreakable, arouse interest and may in themselves be incriminating in countries where encryption is illegal.

### Least Significant Bit (LSB) Steganography

The Pastel Network uses a technique called Least Significant Bit (LSB) steganography to hide the QR code signature within the artwork image. LSB steganography is a simple yet effective method of embedding information in an image by modifying the least significant bits of the pixel color values.

In a digital image, each pixel is represented by a set of values that define its color. For example, in an RGB image, each pixel has a red, a green, and a blue component. Each of these components is typically represented by an 8-bit value, ranging from 0 to 255.

The least significant bit is the last bit in the binary representation of a number. For example, the binary representation of the number 200 is `11001000`. The LSB is the rightmost bit, which is 0. Changing the LSB of a number has a very small impact on its value. For example, if we change the LSB of 200 from 0 to 1, the new value is 201 (`11001001`). This small change in the color value of a pixel is imperceptible to the human eye.

LSB steganography works by replacing the LSB of each color component of each pixel with the bits of the message to be hidden. For example, if we want to hide the letter 'A', which has the ASCII value 65 and the binary representation `01000001`, we can embed it in the first 8 pixels of an image.

### The Go Implementation in `common/image/steganography/`

The core steganography logic in the Pastel GoNode is implemented in the `common/image/steganography/steganography.go` file. This package provides functions for encoding and decoding messages in images using LSB steganography.

#### Encoding a Message

The `Encode` function is the main entry point for encoding a message into an image. It takes an `image.Image` and a byte slice representing the message as input, and returns a new `image.NRGBA` image with the message embedded in it.

```go
// Encode encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 23 pixels
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		message []byte : byte slice of the message to be encoded
		pictureInputFile image.Image : image data used in encoding
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func Encode(pictureInputFile image.Image, message []byte) (*image.NRGBA, error) {

	rgbImage := imageToNRGBA(pictureInputFile)

	if err := EncodeNRGBA(rgbImage, message); err != nil {
		return nil, err
	}
	return rgbImage, nil
}
```

The `Encode` function first converts the input image to the `image.NRGBA` format using the `imageToNRGBA` helper function. It then calls the `EncodeNRGBA` function to perform the actual encoding.

The `EncodeNRGBA` function is where the LSB manipulation happens. Here's a breakdown of the process:

1.  **Message Length**: The function first gets the length of the message and splits it into four bytes. These four bytes are then prepended to the message itself. This is done so that the decoder knows how many bytes to read from the image.

    ```go
    var messageLength = uint32(len(message))
    // ...
    one, two, three, four := splitToBytes(messageLength)

    message = append([]byte{four}, message...)
    message = append([]byte{three}, message...)
    message = append([]byte{two}, message...)
    message = append([]byte{one}, message...)
    ```

2.  **Bit Channel**: A channel is used to feed the bits of the message to the encoding loop. The `getNextBitFromString` function is a goroutine that iterates over the message bytes and sends each bit to the channel.

    ```go
    ch := make(chan byte, 100)
    go getNextBitFromString(message, ch)
    ```

3.  **Pixel Iteration**: The function then iterates over each pixel of the image, and for each pixel, it iterates over the Red, Green, and Blue color components.

    ```go
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            c = rgbImage.NRGBAAt(x, y) // get the color at this pixel
            // ...
        }
    }
    ```

4.  **LSB Manipulation**: For each color component, the function reads a bit from the channel and sets the LSB of the color component to that bit using the `setLSB` function.

    ```go
    /*  RED  */
    bit, ok = <-ch
    if !ok { // if we don't have any more bits left in our message
        rgbImage.SetNRGBA(x, y, c)
        // return nil
    }
    setLSB(&c.R, bit)
    ```

The `setLSB` function is a simple bitwise operation:

```go
// setLSB given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}
```

#### Decoding a Message

The `Decode` function is used to extract a hidden message from an image. It takes the length of the message and the image as input.

```go
// Decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeNRGBA making the conversion from image.Image to image.NRGBA
/*
	Input:
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func Decode(msgLen uint32, pictureInputFile image.Image) (message []byte) {
	return decode(4, msgLen, pictureInputFile) // the offset of 4 skips the "header" where message length is defined
}
```

The `Decode` function calls the `decode` function with an offset of 4. This is to skip the first four bytes of the message, which contain the message length. The `decode` function, in turn, calls `decodeNRGBA` after converting the image to the `NRGBA` format.

The `decodeNRGBA` function works in a similar way to `EncodeNRGBA`. It iterates over the pixels and color components, but instead of setting the LSB, it reads the LSB using the `getLSB` function and reconstructs the message bytes.

```go
// getLSB given a byte, will return the least significant bit of that byte
func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}
```

#### Other Helper Functions

*   `MaxEncodeSize`: This function calculates the maximum number of bytes that can be hidden in an image. The formula is `((width * height * 3) / 8) - 4`. The `* 3` is because we are using the R, G, and B components. The `/ 8` is to convert bits to bytes. The `- 4` is to account for the four bytes used to store the message length.

*   `GetMessageSizeFromImage`: This function reads the first four bytes of the hidden message to determine its length.

This concludes the detailed explanation of the steganography functionality in the Pastel Network. The next part will focus on the QR code signature generation.

## Part 2: QR Code Signatures in Pastel

QR (Quick Response) codes are two-dimensional barcodes that can store a variety of data, including text, URLs, and binary data. The Pastel Network leverages QR codes to embed critical information directly into the artwork image, creating a self-contained and verifiable NFT.

This section provides a detailed explanation of how Pastel generates and uses QR code signatures, and how the Go implementation in `common/image/qrsignature/` works.

### The Purpose of QR Code Signatures

The primary purpose of the QR code signature is to store the following information within the artwork itself:

*   **Robust Image Fingerprint**: The unique fingerprint of the artwork, which is used for duplicate detection.
*   **Artist's Digital Signature**: A cryptographic signature from the artist, proving their ownership and the authenticity of the artwork.
*   **Timestamp**: A timestamp indicating when the artwork was registered.
*   **Other Metadata**: Additional metadata about the artwork, such as its title, description, and keywords.

By embedding this information in the image, Pastel ensures that the NFT is self-contained and can be verified without relying on external servers or databases. This is a crucial aspect of Pastel's decentralized approach to NFT registration.

### The Structure of the QR Code Signature

A single QR code has a limited data capacity. To store all the necessary information, Pastel uses a clever system of multiple QR codes arranged on a single canvas. This canvas is then embedded in the artwork image using steganography.

The QR code signature consists of the following components:

*   **Payload QR Codes**: These are the QR codes that contain the actual data, such as the image fingerprint and the artist's signature. Each payload is split into multiple QR codes if the data exceeds the capacity of a single QR code.
*   **Metadata QR Code**: This is a special QR code that contains metadata about the other QR codes on the canvas. This includes the position and size of each payload QR code. The metadata QR code is always placed at the top-left corner of the canvas, making it easy for the decoder to find and read it first.

### The Go Implementation in `common/image/qrsignature/`

The QR code signature generation logic is implemented in the `common/image/qrsignature/` directory. This package provides the necessary tools for creating, arranging, and decoding the QR codes.

#### `qrcode.go`: The Core QR Code Logic

This file defines the `QRCode` struct, which represents a single QR code. It uses the `github.com/makiuchi-d/gozxing/qrcode` library for the underlying QR code operations.

```go
// QRCode represents QR code image and its coordinates on the canvas.
type QRCode struct {
	image.Image
	imageSize int
	X         int
	Y         int
	writer    *qrcode.QRCodeWriter
	reader    gozxing.Reader
}
```

The `QRCode` struct has methods for encoding and decoding data:

*   `Encode(title, data string) error`: This method takes a title and data as input and generates a new QR code image. The title is displayed above the QR code.
*   `Decode() (string, error)`: This method reads the data from the QR code image and returns it as a string.

#### `payload.go`: Handling the Data

This file defines the `Payload` struct, which represents the data to be encoded in the QR codes. A payload can be split into multiple QR codes if the data is too large for a single QR code.

```go
// Payload represents data of the signature.
type Payload struct {
	name    PayloadName
	data    []byte
	qrCodes []*QRCode
}
```

The `Encode` method of the `Payload` struct is responsible for splitting the data into chunks and creating a QR code for each chunk.

#### `canvas.go`: Arranging the QR Codes

This file defines the `Canvas` struct, which is responsible for arranging the QR codes on a single image. The `FillPos` method takes a slice of `QRCode` structs and sets their `X` and `Y` coordinates, ensuring that they are placed on the canvas without overlapping.

#### `metadata.go`: The Metadata QR Code

This file defines the `Metadata` struct, which represents the metadata QR code. The `Encode` method of the `Metadata` struct takes a slice of `Payload` structs and generates a QR code that contains the positions and sizes of all the payload QR codes.

#### `signature.go`: Orchestrating the Process

This file brings everything together. The `Signature` struct represents the entire QR code signature, including all the payloads and the metadata.

The `Encode` method of the `Signature` struct performs the following steps:

1.  **Encode Payloads**: It first calls the `Encode` method of each payload to generate the payload QR codes.
2.  **Arrange on Canvas**: It then uses the `Canvas` to arrange the QR codes on a single image.
3.  **Encode Metadata**: It generates the metadata QR code, which contains the positions of the payload QR codes.
4.  **Draw on Image**: It draws all the QR codes (payloads and metadata) onto a new image.
5.  **Steganographic Embedding**: Finally, it calls the `steganography.Encode` function to embed the QR code signature image into the original artwork image.

The `Decode` method performs the reverse process. It first extracts the QR code signature image from the artwork using `steganography.Decode`. It then reads the metadata QR code to get the positions of the payload QR codes, and then decodes each payload QR code to retrieve the original data.

This detailed process ensures that a large amount of data can be securely and reliably stored within the artwork itself, making the Pastel NFT a truly self-contained and verifiable digital asset. The next part will explore the robust image fingerprinting technology that is at the heart of Pastel's duplicate detection system.

## Part 3: Robust Image Fingerprinting in Pastel

Robust image fingerprinting, also known as perceptual hashing, is a cornerstone of the Pastel Network's duplicate detection system. Unlike cryptographic hashing algorithms like SHA-256, which are designed to be sensitive to even the slightest change in the input data, robust image fingerprinting algorithms are designed to produce similar or identical fingerprints for visually similar images. This is crucial for detecting duplicate or near-duplicate artworks, even if they have been slightly altered (e.g., resized, compressed, or color-adjusted).

This section provides a deep dive into the theory and practice of robust image fingerprinting as implemented in the Pastel `dd-service`.

### The Need for Robust Image Fingerprinting

Traditional cryptographic hashes are unsuitable for image similarity detection because they are designed to be avalanche-effect compliant. This means that a single-bit change in the input image will result in a completely different hash. While this is desirable for data integrity verification, it is the opposite of what is needed for image similarity detection.

Robust image fingerprinting algorithms, on the other hand, are designed to be resistant to content-preserving modifications. This means that if an image is resized, compressed, or has its colors slightly adjusted, the resulting fingerprint will be very similar to the fingerprint of the original image. This allows the Pastel Network to identify duplicate artworks even if they are not pixel-for-pixel identical.

### Deep Learning and Image Fingerprinting

The Pastel `dd-service` uses a deep learning-based approach to image fingerprinting. This approach leverages the power of convolutional neural networks (CNNs) to learn a rich and discriminative representation of the image content.

The core idea is to use a pre-trained CNN as a feature extractor. The image is fed into the network, and the activations of one of the intermediate layers are used as the image's feature vector. This feature vector is a high-dimensional representation of the image that captures its semantic content.

### The Custom Model in `dd-service`

As of the latest version, the `dd-service` uses a custom deep learning model for feature extraction, which is a significant improvement over the previous ResNet-based approach. The model architecture is defined in the `lib/models/` directory of the `dd-service` repository.

The model is a convolutional neural network (CNN) that has been specifically designed for image retrieval and similarity tasks. It incorporates several advanced techniques, including:

*   **Generalized Mean (GeM) Pooling**: As discovered in our research, GeM pooling is a trainable pooling layer that generalizes max and average pooling. It has been shown to be highly effective in image retrieval tasks. The `dd-service` uses a GeM pooling layer to aggregate the feature maps from the convolutional layers into a compact and discriminative global descriptor. The implementation can be found in `lib/models/gem.py`.

*   **IBN-Net**: The model architecture is based on the IBN-Net (Instance-Batch Normalization Network), which is a deep residual network that incorporates both Instance Normalization (IN) and Batch Normalization (BN) layers. This has been shown to improve the model's generalization ability by learning features that are robust to both appearance and content variations. The IBN-Net implementation can be found in `lib/models/resnet_ibn.py` and `lib/models/resnet_ibn_a.py`.

### The Fingerprinting Process in `dd-service`

The image fingerprinting process in the `dd-service` is handled by the `DupeDetectionTask` class in `lib/dupe_detection.py`. Here's a step-by-step breakdown of the process:

1.  **Image Preprocessing**: The input image is first preprocessed to a standard size and format. This includes resizing the image, converting it to the RGB color space, and normalizing the pixel values. The `prepare_image_for_model_func` in `lib/dupe_detection.py` handles this process.

2.  **Feature Extraction**: The preprocessed image is then fed into the custom deep learning model to extract a feature vector. The `compute_image_deep_learning_features_func` in `lib/dupe_detection.py` is responsible for this step.

3.  **Principal Component Analysis (PCA)**: The extracted feature vector is high-dimensional. To reduce the dimensionality and create a more compact fingerprint, Principal Component Analysis (PCA) is applied. The `dd-service` uses a pre-trained PCA matrix, which is loaded from the `faiss_trained_pca.dat` file. The `PcaImageTransformer` class in `lib/image_utils.py` handles the PCA transformation.

4.  **L2 Normalization**: After PCA, the fingerprint is L2-normalized. This ensures that the fingerprint has a unit length, which is important for the cosine similarity calculation.

The resulting L2-normalized vector is the final robust image fingerprint.

### Similarity Comparison

Once the fingerprint of a candidate image has been generated, it is compared to the fingerprints of all the other artworks in the Pastel Network's database. The `dd-service` uses a combination of three similarity metrics to produce a final similarity score:

1.  **Cosine Similarity**: This is the primary metric used for comparing the fingerprints. It measures the cosine of the angle between two vectors, which is a measure of their similarity. A cosine similarity of 1 means the vectors are identical, while a cosine similarity of 0 means they are orthogonal.

2.  **Hoeffding's D**: This is a non-parametric measure of the distance between two probability distributions. It is used to capture more subtle differences between the fingerprints that may not be captured by cosine similarity alone.

3.  **Hilbert-Schmidt Independence Criterion (HSIC)**: This is another non-parametric measure of the statistical independence of two random variables. It is used to further enhance the accuracy of the similarity comparison.

The `compute_similarity_between_image_and_df_of_comparison_images_func` in `lib/compute_similarity.py` implements the similarity comparison logic. It first calculates the cosine similarity between the candidate fingerprint and all the fingerprints in the database. It then selects the top N most similar fingerprints and calculates Hoeffding's D and HSIC for these top candidates. Finally, it uses a Bayesian Ridge Regression model to combine the three similarity scores into a final, highly accurate similarity score.

This multi-metric approach ensures that the duplicate detection system is both robust and accurate, and can reliably identify duplicate and near-duplicate artworks.

This concludes the detailed explanation of the robust image fingerprinting functionality in the Pastel Network. The final part of this document will tie everything together and describe the end-to-end NFT registration workflow.

## Part 4: Tying It All Together - The Full Workflow

In the previous sections, we have dissected the three core components of Pastel's image processing and registration system: steganography, QR codes, and robust image fingerprinting. Now, let's bring it all together and look at the complete workflow, from the moment an artist decides to register their artwork as an NFT on the Pastel network.

This entire process is designed to be as secure and decentralized as possible, involving multiple actors and checks and balances to ensure the integrity of the system.

### The Actors

Before we dive into the workflow, let's identify the key actors involved:

*   **The Artist:** The creator of the artwork.
*   **The User's Wallet (`walletnode`):** The software that the artist uses to interact with the Pastel network.
*   **Supernodes (`supernode`):** A network of powerful nodes that are responsible for processing and verifying NFT registration requests.
*   **The `dd-service`:** The dupe-detection service that we explored in the previous section.
*   **The Pastel Blockchain:** The decentralized ledger that stores the final NFT registration data.

### The Step-by-Step Workflow

Here is a step-by-step breakdown of the entire process:

**Step 1: The Artist Initiates the Registration**

The process begins when the artist decides to register their artwork as an NFT. They use their Pastel wallet (`walletnode`) to initiate the registration process, providing the artwork image and other relevant metadata (e.g., title, description, etc.).

**Step 2: The Wallet Prepares the Registration Ticket**

The user's wallet then prepares a "registration ticket". This is a data structure that contains all the information needed to register the NFT. The wallet performs the following actions:

1.  **Generate Image Fingerprint:** The wallet communicates with the `dd-service` to generate a robust image fingerprint for the artwork. As we saw in the previous section, this involves running the image through a custom deep learning model, applying PCA, and creating a feature vector.
2.  **Create QR Code Signature:** The wallet creates a QR code signature that contains the following information:
    *   The image fingerprint.
    *   A digital signature from the artist, proving their ownership of the artwork.
    *   A timestamp.
    *   Other relevant metadata.
    As we saw in Part 2, this may involve creating multiple QR codes that are arranged on a canvas, with a special metadata QR code to map their positions.
3.  **Embed QR Code Signature:** The wallet then uses the steganography functionality (as described in Part 1) to embed the QR code signature image into the original artwork image. This creates a new image that looks identical to the original but contains the hidden QR code signature.

**Step 3: The Wallet Submits the Registration Ticket**

The wallet then submits the registration ticket to the Pastel network. The ticket contains the modified image (with the embedded QR code signature) and other relevant information.

**Step 4: The Supernodes Process the Registration Ticket**

A group of Supernodes is randomly selected to process the registration ticket. These Supernodes perform a series of checks to verify the authenticity and originality of the artwork.

1.  **Extract QR Code Signature:** The Supernodes extract the QR code signature from the image using the steganography decoding functionality.
2.  **Decode QR Code Signature:** They then decode the QR codes to retrieve the image fingerprint, the artist's digital signature, and the other metadata.
3.  **Verify Digital Signature:** The Supernodes verify the artist's digital signature to ensure that the registration request is legitimate.
4.  **Duplicate Detection:** This is the most critical step. The Supernodes use the extracted image fingerprint to query the `dd-service` and check for duplicates. The `dd-service` compares the fingerprint of the new artwork to the fingerprints of all the other artworks that have already been registered on the network. As we saw in the previous section, this involves calculating a similarity score based on cosine similarity, Hoeffding's D, and HSIC.

**Step 5: The Supernodes Reach a Consensus**

Based on the results of the duplicate detection process, the Supernodes vote on whether to approve or reject the registration request. If a sufficient number of Supernodes agree that the artwork is original, the registration is approved.

**Step 6: The NFT is Registered on the Blockchain**

If the registration is approved, a new NFT is created on the Pastel blockchain. The registration data, including the image fingerprint and a link to the artwork file, is permanently stored on the blockchain.

### The Role of this Process in the Pastel Ecosystem

This entire process, from image fingerprinting to steganographic embedding, plays a crucial role in the Pastel ecosystem. It provides a robust and secure solution to some of the most pressing challenges in the NFT space:

*   **Authenticity:** The use of digital signatures and steganography ensures that the NFT is genuinely created by the artist and has not been tampered with.
*   **Originality:** The robust image fingerprinting and duplicate detection system provides a high degree of confidence that the NFT is a unique and original artwork.
*   **Verifiability:** The entire registration process is transparent and verifiable. Anyone can independently verify the authenticity and originality of an NFT by examining the data on the blockchain and in the artwork file itself.
*   **Decentralization:** The use of a network of Supernodes to process and verify registration requests ensures that the system is not controlled by any single entity.

By combining these different technologies in a novel and innovative way, the Pastel network has created a system that is greater than the sum of its parts. It is a system that provides a solid foundation for a more secure, trustworthy, and vibrant NFT ecosystem.

## Conclusion

In this document, we have taken a deep dive into the steganography, QR code generation, and robust image fingerprinting functionalities within the Pastel network. We have seen how these three components work together to create a sophisticated and secure system for NFT registration and duplicate detection.

We started by exploring the theoretical concepts of LSB steganography and then delved into the practical implementation in the GoNode codebase. We saw how Pastel uses this technique to hide data within images in a way that is imperceptible to the human eye.

Next, we looked at how Pastel uses QR codes to structure and store complex data, including digital signatures and metadata. We examined the `common/image/qrsignature` package and saw how it uses a clever metadata QR code to manage multiple QR codes on a single canvas.

Then, we moved on to the most critical part of the system: robust image fingerprinting. We explored the theoretical concepts of perceptual hashing and then dissected the `dd-service` codebase. We saw how Pastel uses a state-of-the-art custom deep learning model, PCA, and an ensemble of advanced similarity metrics to create a highly accurate duplicate detection system.

Finally, we brought everything together and looked at the complete workflow, from the moment an artist initiates the registration process to the point where the NFT is registered on the Pastel blockchain. We saw how the different components of the system work in concert to ensure the authenticity, originality, and verifiability of the NFTs on the platform.

The Pastel network's innovative use of these technologies is a testament to its commitment to building a more secure and trustworthy NFT ecosystem. By providing a robust solution to the problem of duplicate and near-duplicate NFTs, Pastel is helping to create a more level playing field for artists and collectors alike.

As the NFT space continues to evolve, it is clear that solutions like the one developed by Pastel will be essential for ensuring the long-term health and viability of the market. The combination of steganography, QR codes, and robust image fingerprinting provides a powerful and flexible framework for building a new generation of decentralized applications that are more secure, transparent, and trustworthy than ever before.
