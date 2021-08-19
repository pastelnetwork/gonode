package database

import (
	"reflect"
	"testing"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

func Test_pbToWriteCommand(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.UserdataRequest
		want UserdataWriteCommand
	}{
		{
			name: "TestUserdata_ToWriteCommand1",
			req: &pb.UserdataRequest{
				RealName:        "cat",
				FacebookLink:    "fb.com",
				TwitterLink:     "tw.com",
				NativeCurrency:  "usd",
				Location:        "us",
				PrimaryLanguage: "en",
				Categories:      "a",
				Biography:       "b",
				AvatarImage: &pb.UserdataRequest_UserImageUpload{
					Content:  nil,
					Filename: "123",
				},
				CoverPhoto: &pb.UserdataRequest_UserImageUpload{
					Content:  []byte{4, 5, 6},
					Filename: "456",
				},
				ArtistPastelID:    "id",
				Timestamp:         123,
				Signature:         "sign",
				PreviousBlockHash: "hash",
			},
			want: UserdataWriteCommand{
				RealName:           "cat",
				FacebookLink:       "fb.com",
				TwitterLink:        "tw.com",
				NativeCurrency:     "usd",
				Location:           "us",
				PrimaryLanguage:    "en",
				Categories:         "a",
				Biography:          "b",
				AvatarImage:        "",
				AvatarFilename:     "123",
				CoverPhotoImage:    "040506",
				CoverPhotoFilename: "456",
				ArtistPastelID:     "id",
				Timestamp:          123,
				Signature:          "sign",
				PreviousBlockHash:  "hash",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pbToWriteCommand(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pbToWriteCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
