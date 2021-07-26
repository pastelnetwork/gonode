package database

import (
	"reflect"
	"testing"
)

func TestUserdata_ToWriteCommand(t *testing.T) {
	type fields struct {
		Realname          string
		FacebookLink      string
		TwitterLink       string
		NativeCurrency    string
		Location          string
		PrimaryLanguage   string
		Categories        string
		Biography         string
		AvatarImage       UserImageUpload
		CoverPhoto        UserImageUpload
		ArtistPastelID    string
		Timestamp         int
		Signature         string
		PreviousBlockHash string
	}
	tests := []struct {
		name   string
		fields fields
		want   UserdataWriteCommand
	}{
		{
			name: "TestUserdata_ToWriteCommand1",
			fields: fields{
				Realname:        "cat",
				FacebookLink:    "fb.com",
				TwitterLink:     "tw.com",
				NativeCurrency:  "usd",
				Location:        "us",
				PrimaryLanguage: "en",
				Categories:      "a",
				Biography:       "b",
				AvatarImage: UserImageUpload{
					Content:  nil,
					Filename: "123",
				},
				CoverPhoto: UserImageUpload{
					Content:  []byte{4, 5, 6},
					Filename: "456",
				},
				ArtistPastelID:    "id",
				Timestamp:         123,
				Signature:         "sign",
				PreviousBlockHash: "hash",
			},
			want: UserdataWriteCommand{
				Realname:           "cat",
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
			d := &Userdata{
				Realname:          tt.fields.Realname,
				FacebookLink:      tt.fields.FacebookLink,
				TwitterLink:       tt.fields.TwitterLink,
				NativeCurrency:    tt.fields.NativeCurrency,
				Location:          tt.fields.Location,
				PrimaryLanguage:   tt.fields.PrimaryLanguage,
				Categories:        tt.fields.Categories,
				Biography:         tt.fields.Biography,
				AvatarImage:       tt.fields.AvatarImage,
				CoverPhoto:        tt.fields.CoverPhoto,
				ArtistPastelID:    tt.fields.ArtistPastelID,
				Timestamp:         tt.fields.Timestamp,
				Signature:         tt.fields.Signature,
				PreviousBlockHash: tt.fields.PreviousBlockHash,
			}
			if got := d.ToWriteCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Userdata.ToWriteCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
