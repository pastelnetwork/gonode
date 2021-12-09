package pastel

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeNFTTicket(t *testing.T) {
	inputAppData := AppTicket{
		PreviewHash:           []byte{1},
		Thumbnail1Hash:        []byte{2},
		Thumbnail2Hash:        []byte{3},
		DataHash:              []byte{4},
		FingerprintsHash:      []byte{6},
		FingerprintsSignature: []byte{7},

		RQIDs: []string{"9", "10"},
	}

	inputTicket := NFTTicket{
		Version:       1,
		Author:        string([]byte{2, 3, 4}),
		BlockNum:      5,
		BlockHash:     string([]byte{6, 7, 8}),
		Copies:        9,
		Royalty:       10,
		Green:         false,
		AppTicketData: inputAppData,
	}

	encoded, err := EncodeNFTTicket(&inputTicket)
	assert.Nil(t, err)
	outputTicket, err := DecodeNFTTicket(encoded)
	outputAppData := outputTicket.AppTicketData
	assert.Nil(t, err)
	fmt.Println(string(encoded))
	assert.Equal(t, inputTicket.Version, outputTicket.Version)
	assert.Equal(t, inputTicket.Author, outputTicket.Author)
	assert.Equal(t, inputTicket.BlockNum, outputTicket.BlockNum)
	assert.Equal(t, inputTicket.BlockHash, outputTicket.BlockHash)
	assert.Equal(t, inputTicket.Copies, outputTicket.Copies)
	assert.Equal(t, inputTicket.Royalty, outputTicket.Royalty)
	assert.Equal(t, inputTicket.Green, outputTicket.Green)

	assert.Equal(t, inputAppData.PreviewHash, outputAppData.PreviewHash)
	assert.Equal(t, inputAppData.Thumbnail1Hash, outputAppData.Thumbnail1Hash)
	assert.Equal(t, inputAppData.DataHash, outputAppData.DataHash)
}

func TestDecodeNFTTicket(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in             string
		wantAuthor     string
		wantTitle      string
		wantCreatorURL string
	}{
		"a": {
			in:             `eyJuZnRfdGlja2V0X3ZlcnNpb24iOjEsImF1dGhvciI6ImpYWTc0amE3ZWRwdzFteGZIempKZ1V3Y0RmMmhKY2dRY2pVZHVaRVl0VXVuUWs3dkZQQ2tCOHN0N0Q5ODZLc3B4c1N1TEFKWHRjSkEzUFR5Z1NxcktQIiwiYmxvY2tudW0iOjExOTAwMiwiYmxvY2tfaGFzaCI6IjAwMDAwM2MxMThkM2I5ODI3ODdhMmE2MjVkNzEyYThmNzBhNzA0YThhOTdiZjMyODY3NWM3Mzg0NDk0Y2UzYzYiLCJjb3BpZXMiOjEsInJveWFsdHkiOjAsImdyZWVuIjpmYWxzZSwiYXBwX3RpY2tldCI6IkhRbG1ARkQsVDU/WjlGcEZDZjhxQmsvREssXHUwMDI2MTRfMmAsOFgyZT9OPUc7M3JaQWs3NV84blYuSkBtam9hQk0pY0k7LkY/QEE5RCdFPWBkWWRER2IxLkdcIiw6cUNJMVwiJUZcdTAwM2VsVEozXHUwMDI2WEVORS1RR1x1MDAyNkZcXEcobj1FSiVGNXIyXHUwMDNlOkhcIlVcdTAwMjZxRV9CNFEvMFxcUFx1MDAzZURlIVFpRkVoLyxBME9cdTAwMjZZK3RPcElDaVx1MDAzY2chP1ooXCInLCElRzUzQVx1MDAzYy1CLzBcXFNFQVJUXFwhRWFOXW1ELk9uUCxcInVRZERJSVFyRF1pUyErQkUyb0BxOTpLR1x1MDAyNkNyOEhcIkNrbUFUTXJMK3RPcEpFYi9bJERmVDMvQVJdXnBGQ2NYVywldSg/RS1cIlx1MDAyNm4wNl9WYS9vNSo7RElJUXJEZSo6XCJCbDdFcC9vRzZCK3RPcEpFYi9bJERmVDMvRWJUVztBU3RqckZDQjlcdTAwMjZELlJVLCwhJG9JQk9yO1xcRGYvdSs5UEpRVitER21cdTAwM2VAO1szK0JsXCJvMkA7MFYkQmw3UStAWDAoZkZDQiFcIkA7WzJyRWNjIzpGPVxcUEtES11UL0ZEbDFOK0NULnUrRXFqP0ZDZko4OVBcdTAwMjYtW0BcdTAwM2MsbXMrQ25uJ1x1MDAzY2JaR11CZUNOJUJrcTlyR3AkVThEL2FcdTAwM2NcdTAwMjZGQ2V1KkJsNSVKMikkcFx1MDAzZStFVjouK0UyJStAcT9jcEFSVFxcJ0ViMDs3QDNCTi1FY2NcdTAwM2UxQmxrSjlBZnRNKUYoSmQlQVJscC1EZlx1MDAyNnAjK0VNKzkrQ1Nua0JsOCE2K0NULjFCbC45a0JsN0BcIkdwJHAzRElkZnVAO29kMSt0T3BVQW9xQyVCbG4wXHUwMDI2LCEkb0JEZi91KzlQSlFWK3RPcFVBb3FDJUgjUj09M1svOVRAOnMtcS8wXFx0REZDMCojRWJUKis/Wic6aEFLai9aN1VeLlBGYFNbNkVjYyM6Rlx1MDAzY1c3W0RKIXV0QHJjLWhGRDVaMj9acC1vQVQocSREZnAvRUBWXHUwMDI2bmhFYmxwXFwsJXUoP0UtXCJcdTAwMjZuMDZfVmEvcFZBWEZFTVx1MDAyNigvbjhnOjA2XmlIQHFcXClNNFlWK3MySD1ALTFnNCE4K3RPcFVBb3FCcUFVL0tcdTAwM2NFYlx1MDAyNk51QVRUK1csI1Z1aUA7MGUoQDtdUmQvMEtcIkdGX3RRL0RfKiNNRGZUcj9AOzBnX0FUTXJMK3RPcFtEZmZFKD9YblwibkJrOz8wM1xcYDAwRSxvTjVCaztLcUJPUHNxLCEkb0Q4aUxeMVx1MDAzYy1MdTc9W0BcIjEzK09hN0E5Vz9cIkIwZXNxQjM5RlQ2VWwyYjJLKjR0OFBhMVI3ODVvaC8wXTFMRl90VCpAOzBPND9ZRWtoQkhmSl06SipRYkZcdTAwMjZbP1ZcdTAwM2MqO00rPXVeJUk4N1pwalx1MDAzY0crVGRAVjhlZzpJQFRvNjhpYXA3Uy1CbS5uKi9xK3RPcFtCUVxcMCRESUk2cTExK14nRihjYVksXCJcImg5ODdRK0E5bTktXCJHXVx1MDAzZWBEQVEyY1NIXCJnXHUwMDNlbDlrblE4O2VTb286SnUlJEUsbmBNRCx0Qy40WClGJEE3OVJnP1lFa2hCSGZKXTYkRzFoQ0xwcGBAXHUwMDNjLENlM0NibyM5MXIwITtHcTBxOiw/YHQ6Zl9kITBKYlhpSCFzQjs6ZTQvQCt0T3BNQmw3UXBFYz9cdTAwMjY1REtLcixCT1BzcSwhJG9HMUoxQT9DR2VSYjEvXW9cdTAwMjZDTGdGSlx1MDAzYylRYi0wUEcyczdcdTAwM2NnUmNDTiskIT0oSDhdQ2RMbHA9RFZxYlx1MDAzY0UxaFE3bixzS1x1MDAzY0RsTGJDSXBaNS8wXFxcXD9ESipPJEUsb1oxRkU6ZjFCa01cdTAwM2NsRkVNVjgsISRvRjBQI0U6MTIoSyg9RFRpaDZzPVxcMFx1MDAzYytvQi4xTjY4akNmM28rPUBtTzU4NTtRUTEuXHUwMDNjcnE9I2knNzc6XjQ4XHUwMDNjYVx1MDAyNlg6XHUwMDNjKS1McTBMZm5pN1JfVDNDLy45LUQtQ0BPQ01TKlx1MDAyNjFLSjZkXHUwMDNjKj1FQ1x1MDAzY0U9RUE3UmY3RENlQDktN3FhVlM9I2s9dDdzXHUwMDI2UVc3XHUwMDNjaVdLQ0ptJyRDZnJUJzEsMVhzMU1nbDIxL18oSTdSZmpWNzYsRWxcdTAwM2NgWDA3MUoyWHIxMFJPY0QvPS9rXHUwMDNjLD1tJVx1MDAzY2AyXmtDZHFpNkMxS1xcMz1CXFw2cVx1MDAzY2IrZytDaHVPT0MuX29EPSpSLEE9J0A6ZjdvaSNdMExcXCxLLzBcXFZJRStOQmVBVFZLbkZENVoyP1pVTDZGQ2Y7ckclR104QmxAbDUzWy06MjBIciVsRC4uTnJCT3U2bEFvRGcwQTcnN20/WTRcIm1GKihjLkA6cy1xM1xcVyovREtUZipBVEQtckFtXUxjQjRaLWtEZVx1MDAzYz9zQVROITFGRThXZTBlPU1rRWJvKiRBbV0uYUVjWlx1MDAzZTBELi5OckJPXHUwMDNlSWs/WjlGZEFLai9aOi1nJ1QsXHUwMDI2Z3QzRkNmOHFFYWEhXCJESW1tMT9aVF5xRWItRlUwSjVAQTFHZ2dDMStYVmBES0tIMURJbW9zRWFhIVwiREltbTE/WlRecUViLUZVMEpcIkRkRStOb29ES0JFNj9aVF5xRWItRlUwSjUlOjBmXydJMEtEJEIsJTU7MEFURFpzRkNlZnNGKFRXJ0YoOS0vQVRLJVZIUW0hQEBcdTAwM2NaRidCLUtBai9oZlwiXHUwMDNjMEsxc04yRFptLkJPdSgnQDstb0gwSjUlNjFHTGpHMEprVTssXHUwMDI2VXQ3RkUxZissISVEMTNBV1dSMWJnXkIvMF0lT0VjKideMEo1JTUzQWlLSjJFKlRILCcuPT9HcDU6Jy9oZiU5MWNbNkUzLjNcdTAwM2UlQmwuOWtBUkIrWkYoZi0rLCEoXHUwMDI2cEUrRXJxQk9Qc3EsISRvKzNGdDMsMGtFMHQxR1dLITJgXHUwMDNlMlcyYEVRTkFOO1NWQWhjK3MxLDFbcjNcdTAwMjYqMEUyKVxcbycxMWA3T0BVXyxKMi5nO1czK2NTVTNcdTAwMjZrRCtAbC5bXHUwMDNlLzBdJUVFYXJbXCJGRU0jLj9ZRWtoQkhmSl1BaVwiJFRBaT0nVTNcdTAwMjZFPSExY1I0Iyt0T3BIRyVHXVx1MDAyNkI0WUZgQFx1MDAzYzYqKzNbLzBKQVMqVlUwSkcxNzBKRzE3MEpGXFwlLCVQRFwiQW4/IW9ESVs2YkJPUHNxLCEkbyVAbC5hV0A6TTJRQDpOXyQxY1tFOy8wXFx0Q0ZgTG8sQk9Qc3E/WUVraEJIZkpdMEpZST8yRVx1MDAzZTsuMUxhaVMwa05GXCJAVi5XM0FpWD9VLCgyIWRCbGRXdEJrcTlySCFiKilFK0w0U0Ftb0xzQUxvJFx1MDAzZUYnaXJyRWIvVHJESTcqcUZDZksxQVRUK1dGRTJNOC8wXStTP1lPJWwsIVx1MDAyNnBQN29YRyVGI0AsaUBUYyxZR1hiLDY9REUjJEg9OUhIQTQoS105UEBsaUJRbk1QRC4sO0VFYWkzaCt0T29xOiw1YSk2cGpDYjtjWm5DR3VSXmA6ZGZVaDkxczRhXHUwMDNlJCtqTUVDT0N0Nzk9UFRHXHUwMDI2OllTRSV1R1kvMFpsaTMpWDBiOjJzSFY4OWRVLDk0VnNFMWcrcmk7YEpyXHUwMDNlODZcdTAwMjZaJEYjJGpcdTAwM2M7X2A1V0RGZTJORDAtc1EsXCI0dDc7ZHNCNz0qSnQrPSlySW43OUR1ajEsNCNxXHUwMDNjKF1HM0NOcktcdTAwM2NHIVNXO0JMXHUwMDNlWzdIOXRiTUYhXHUwMDNjLlo2dEpBOUIvTlVVR0JrckQ7RmoyZjg3dHQ3RztoajM3c0FLXHUwMDNjQjYuVSRHQTlcdTAwMjZtPSlgSUUyYEYjcSt0T3AxM0QrTnBcdTAwM2VcdTAwMjYkXHUwMDI2VEJtMzAqRWA2SSw4altcdTAwM2MsN3F1WzpCZlRUWkZER1x1MDAyNnBHPSxNMDFObDhZN1Q9SVIvMFtBakUpVVhhQywnSUIwaWd0WDFPO101RClQNmI6MDA4dDk0VUlHQmw/c2dBbDInNTMscm5ZRWFnNEQsXCJRWlZHJ1opRkZcXFhrU0ddWkQ3NzpVOllAVkFTTTc7WipkOGpibWM9XlwiTC9BNW04YkFNXFw2XUNFYjtSM0FhYDA7K2MrSEhcdTAwM2VZSD09KSg2TUBcdTAwM2NrT1Y6aHJbOTZwYmxOQW1vWElCL19cXEk5bTkjXjJmMVx1MDAzY04rdE9wKzNHcWVTQVUlOlwiN1VlbCo3bEc/V0ZcdTAwM2VAISQ7TCg5QT1cdTAwM2VrNURcdTAwM2NIYC44NzU4W1FHWEhbT0MwPVVeXHUwMDNlcUAxP0VGM1grQmQsU141c1tlVDMpTUBxRzsyQk01c1txMiwoMEZrIn0=`,
			wantTitle:      "Mona Lisa",
			wantCreatorURL: "https://www.leonardodavinci.net",
			wantAuthor:     "jXY74ja7edpw1mxfHzjJgUwcDf2hJcgQcjUduZEYtUunQk7vFPCkB8st7D986KspxsSuLAJXtcJA3PTygSqrKP",
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nftTicketIn, err := base64.StdEncoding.DecodeString(tc.in)
			assert.Nil(t, err)

			nft, err := DecodeNFTTicket(nftTicketIn)

			assert.Nil(t, err)
			assert.Equal(t, nft.AppTicketData.CreatorWebsite, "https://www.leonardodavinci.net")
			assert.Equal(t, "jXY74ja7edpw1mxfHzjJgUwcDf2hJcgQcjUduZEYtUunQk7vFPCkB8st7D986KspxsSuLAJXtcJA3PTygSqrKP", nft.Author)
			assert.Equal(t, nft.AppTicketData.NFTTitle, "Mona Lisa")
		})
	}
}
