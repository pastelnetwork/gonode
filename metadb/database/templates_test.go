package database

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	userdata1 = UserdataWriteCommand{
		Realname:           "cat",
		FacebookLink:       "fb.com",
		TwitterLink:        "tw.com",
		NativeCurrency:     "usd",
		Location:           "us",
		PrimaryLanguage:    "en",
		Categories:         "a",
		Biography:          "b",
		AvatarImage:        "",
		AvatarFilename:     "1234.jpg",
		CoverPhotoImage:    fmt.Sprintf("%x", []byte{4, 5, 6, 7}),
		CoverPhotoFilename: "4567.jpg",
		ArtistPastelID:     "abc",
		Timestamp:          123,
		Signature:          "xyz",
		PreviousBlockHash:  "hash",
	}

	userdata2 = UserdataWriteCommand{
		Realname:           "cat",
		FacebookLink:       "fb.com",
		TwitterLink:        "tw.com",
		NativeCurrency:     "usd",
		Location:           "us",
		PrimaryLanguage:    "en",
		Categories:         "a",
		Biography:          "b",
		AvatarImage:        fmt.Sprintf("%x", []byte{1, 2, 3, 4}),
		AvatarFilename:     "1234.jpg",
		CoverPhotoImage:    fmt.Sprintf("%x", []byte{4, 5, 6, 255}),
		CoverPhotoFilename: "4567.jpg",
		ArtistPastelID:     "xyz",
		Timestamp:          123,
		Signature:          "xyz",
		PreviousBlockHash:  "hash",
	}

	writeTemplateResult1 = `INSERT OR REPLACE INTO user_metadata (
		artist_pastel_id,
		real_name,
		facebook_link,
		twitter_link,
		native_currency,
		location,
		primary_language,
		categories,
		biography,
		timestamp,
		signature,
		previous_block_hash,
		avatar_image,
		avatar_filename,
		cover_photo_image,
		cover_photo_filename
	) VALUES (
		'abc',
		'cat',
		'fb.com',
		'tw.com',
		'usd',
		'us',
		'en',
		'a',
		'b',
		123,
		'xyz',
		'hash',
		NULL,
		'1234.jpg',
		x'04050607',
		'4567.jpg'
	);`

	writeTemplateResult2 = `INSERT OR REPLACE INTO user_metadata (
		artist_pastel_id,
		real_name,
		facebook_link,
		twitter_link,
		native_currency,
		location,
		primary_language,
		categories,
		biography,
		timestamp,
		signature,
		previous_block_hash,
		avatar_image,
		avatar_filename,
		cover_photo_image,
		cover_photo_filename
	) VALUES (
		'xyz',
		'cat',
		'fb.com',
		'tw.com',
		'usd',
		'us',
		'en',
		'a',
		'b',
		123,
		'xyz',
		'hash',
		x'01020304',
		'1234.jpg',
		x'040506ff',
		'4567.jpg'
	);`
)

type testSuiteTemplate struct {
	suite.Suite
	k *templateKeeper
}

func (ts *testSuiteTemplate) SetupSuite() {
	k, err := NewTemplateKeeper("./commands")
	assert.Nil(ts.T(), err)
	ts.k = k
}

func (ts *testSuiteTemplate) Test_templateKeeper_GetTemplate() {
	tests := []struct {
		key        string
		assertFunc func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			key:        "unique_nft_by_user",
			assertFunc: assert.NotNil,
		},
		{
			key:        "get_friend",
			assertFunc: assert.NotNil,
		},
		{
			key:        "abc",
			assertFunc: assert.Nil,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("Test_templateKeeper_GetTemplate-%d", i), func(t *testing.T) {
			tt.assertFunc(t, ts.k.GetTemplate(tt.key))
		})
	}
}

func (ts *testSuiteTemplate) TestNewTemplateKeeper() {
	tests := []struct {
		templateDir string
		wantErr     bool
		assertFunc  func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			templateDir: "./commands",
			wantErr:     false,
			assertFunc:  assert.NotNil,
		},
		{
			templateDir: "./abc",
			wantErr:     true,
			assertFunc:  assert.Nil,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestNewTemplateKeeper-%d", i), func(t *testing.T) {
			got, err := NewTemplateKeeper(tt.templateDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTemplateKeeper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.assertFunc(t, got)
		})
	}
}

func (ts *testSuiteTemplate) Test_substituteTemplate() {
	type args struct {
		tmpl *template.Template
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test_substituteTemplate1",
			args: args{
				tmpl: ts.k.GetTemplate("user_info_query"),
				data: "123123",
			},
			want:    `SELECT * FROM user_metadata WHERE artist_pastel_id = '123123';`,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate2",
			args: args{
				tmpl: ts.k.GetTemplate("user_info_query"),
				data: "abc123xyz",
			},
			want:    `SELECT * FROM user_metadata WHERE artist_pastel_id = 'abc123xyz';`,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate3",
			args: args{
				tmpl: ts.k.GetTemplate("user_info_write"),
				data: userdata1,
			},
			want:    writeTemplateResult1,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate4",
			args: args{
				tmpl: ts.k.GetTemplate("user_info_write"),
				data: userdata2,
			},
			want:    writeTemplateResult2,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplateError",
			args: args{
				tmpl: nil,
				data: map[string]string{
					"ArtistPastelID": "abc123xyz",
				},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			got, err := substituteTemplate(tt.args.tmpl, tt.args.data)
			if (err != nil) != tt.wantErr {
				ts.T().Errorf("substituteTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				ts.T().Errorf("substituteTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (ts *testSuiteTemplate) Test_templateKeeper_GetCommand(t *testing.T) {
	tests := []struct {
		key     string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			key:     "user_info_write",
			data:    userdata1,
			want:    writeTemplateResult1,
			wantErr: false,
		},
		{
			key:     "user_info_write",
			data:    userdata2,
			want:    writeTemplateResult2,
			wantErr: false,
		},
		{
			key:     "abc",
			data:    userdata2,
			want:    "",
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("Test_templateKeeper_GetCommand-%d", i), func(t *testing.T) {
			got, err := ts.k.GetCommand(tt.key, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("templateKeeper.GetCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("templateKeeper.GetCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
