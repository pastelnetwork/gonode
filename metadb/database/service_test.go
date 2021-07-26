package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"text/template"

	"github.com/pastelnetwork/gonode/metadb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	schema = `CREATE TABLE IF NOT EXISTS user_metadata (
		artist_pastel_id TEXT PRIMARY KEY UNIQUE,
		real_name TEXT,
		facebook_link TEXT,
		twitter_link TEXT,
		native_currency TEXT,
		location TEXT,
		primary_language TEXT,
		categories TEXT,
		biography TEXT,
		timestamp INTEGER NOT NULL,
		signature TEXT NOT NULL,
		previous_block_hash TEXT NOT NULL,
		avatar_image BLOB,
		avatar_filename TEXT,
		cover_photo_image BLOB,
		cover_photo_filename TEXT
	)`

	queryTemplate = `SELECT * FROM user_metadata WHERE artist_pastel_id = '{{.}}'`

	writeTemplate = `INSERT OR REPLACE INTO user_metadata (
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
		'{{.ArtistPastelID}}',
		'{{.Realname}}',
		'{{.FacebookLink}}',
		'{{.TwitterLink}}',
		'{{.NativeCurrency}}',
		'{{.Location}}',
		'{{.PrimaryLanguage}}',
		'{{.Categories}}',
		'{{.Biography}}',
		{{.Timestamp}},
		'{{.Signature}}',
		'{{.PreviousBlockHash}}',
		{{ if (eq .AvatarImage "")}}NULL,{{ else }}x'{{.AvatarImage}}',{{ end }}
		'{{.AvatarFilename}}',
		{{ if (eq .CoverPhotoImage "")}}NULL,{{ else }}x'{{.CoverPhotoImage}}',{{ end }}
		'{{.CoverPhotoFilename}}'
	)`

	updateTemplate = `UPDATE
    user_metadata
	SET
		real_name = '{{.Realname}}',
		facebook_link = '{{.FacebookLink}}',
		twitter_link = '{{.TwitterLink}}',
		native_currency = '{{.NativeCurrency}}',
		location = '{{.Location}}',
		primary_language = '{{.PrimaryLanguage}}',
		categories = '{{.Categories}}',
		biography = '{{.Biography}}',
		timestamp = {{.Timestamp}},
		signature = '{{.Signature}}',
		previous_block_hash = '{{.PreviousBlockHash}}',
		avatar_image = {{ if (eq .AvatarImage "")}}NULL,{{ else }}x'{{.AvatarImage}}',{{ end }}
		avatar_filename = '{{.AvatarFilename}}',
		cover_photo_image = {{ if (eq .CoverPhotoImage "")}}NULL,{{ else }}x'{{.CoverPhotoImage}}',{{ end }}
		cover_photo_filename = '{{.CoverPhotoFilename}}'
	WHERE
		artist_pastel_id = '{{.ArtistPastelID}}'
	;`

	data1 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "abc",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data2 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "xyz",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data3 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "qwe",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data4 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "rty",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data5 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "data5_1.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "data5",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data6 = Userdata{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "data6_1.jpg",
		},
		CoverPhoto: UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "data6_2.jpg",
		},
		ArtistPastelID:    "data6",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

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
	)`

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
	)`

	updateTemplateResult1 = `UPDATE
    user_metadata
	SET
		real_name = 'cat',
		facebook_link = 'fb.com',
		twitter_link = 'tw.com',
		native_currency = 'usd',
		location = 'us',
		primary_language = 'en',
		categories = 'a',
		biography = 'b',
		timestamp = 123,
		signature = 'xyz',
		previous_block_hash = 'hash',
		avatar_image = NULL,
		avatar_filename = '1234.jpg',
		cover_photo_image = x'04050607',
		cover_photo_filename = '4567.jpg'
	WHERE
		artist_pastel_id = 'abc'
	;`

	updateTemplateResult2 = `UPDATE
    user_metadata
	SET
		real_name = 'cat',
		facebook_link = 'fb.com',
		twitter_link = 'tw.com',
		native_currency = 'usd',
		location = 'us',
		primary_language = 'en',
		categories = 'a',
		biography = 'b',
		timestamp = 123,
		signature = 'xyz',
		previous_block_hash = 'hash',
		avatar_image = x'01020304',
		avatar_filename = '1234.jpg',
		cover_photo_image = x'040506ff',
		cover_photo_filename = '4567.jpg'
	WHERE
		artist_pastel_id = 'xyz'
	;`
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite

	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	workDir string
	ops     *DatabaseOps
}

func (ts *testSuite) SetupSuite() {
	queryTmpl, err := template.New("query").Parse(queryTemplate)
	assert.Nil(ts.T(), err)
	writeTmpl, err := template.New("write").Parse(writeTemplate)
	assert.Nil(ts.T(), err)
	updateTmpl, err := template.New("update").Parse(updateTemplate)
	assert.Nil(ts.T(), err)

	workDir, err := ioutil.TempDir("", "metadb-*")
	assert.NoError(ts.T(), err)

	config := metadb.NewConfig()
	config.SetWorkDir(workDir)
	ts.workDir = workDir

	db := metadb.New(config, "uuid", []string{})
	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	ts.wg.Add(1)
	go func() {
		defer ts.wg.Done()
		// start the rqlite server
		ts.Nil(db.Run(ts.ctx), "run service")
	}()
	db.WaitForStarting()
	_, err = db.Write(ts.ctx, schema)

	ts.Nil(err)
	ts.ops = &DatabaseOps{
		metaDB:         db,
		writeTemplate:  writeTmpl,
		queryTemplate:  queryTmpl,
		updateTemplate: updateTmpl,
	}
	ts.Nil(ts.ops.WriteUserData(ts.ctx, data3))
	ts.Nil(ts.ops.WriteUserData(ts.ctx, data4))
	ts.Nil(ts.ops.WriteUserData(ts.ctx, data5))
	ts.Nil(ts.ops.WriteUserData(ts.ctx, data6))
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.wg.Wait()

	// remove the data directly
	os.RemoveAll(filepath.Join(ts.workDir))
}

func (ts *testSuite) Test_substituteTemplate() {
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
				tmpl: ts.ops.queryTemplate,
				data: "123123",
			},
			want:    `SELECT * FROM user_metadata WHERE artist_pastel_id = '123123'`,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate2",
			args: args{
				tmpl: ts.ops.queryTemplate,
				data: "abc123xyz",
			},
			want:    `SELECT * FROM user_metadata WHERE artist_pastel_id = 'abc123xyz'`,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate3",
			args: args{
				tmpl: ts.ops.writeTemplate,
				data: userdata1,
			},
			want:    writeTemplateResult1,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate4",
			args: args{
				tmpl: ts.ops.writeTemplate,
				data: userdata2,
			},
			want:    writeTemplateResult2,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate5",
			args: args{
				tmpl: ts.ops.updateTemplate,
				data: userdata1,
			},
			want:    updateTemplateResult1,
			wantErr: false,
		},
		{
			name: "Test_substituteTemplate6",
			args: args{
				tmpl: ts.ops.updateTemplate,
				data: userdata2,
			},
			want:    updateTemplateResult2,
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

func (ts *testSuite) TestDatabaseOps_WriteUserData() {
	type args struct {
		ctx  context.Context
		data Userdata
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestDatabaseOps_WriteUserData1",
			args: args{
				ctx:  ts.ctx,
				data: data1,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData2",
			args: args{
				ctx:  ts.ctx,
				data: data2,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData3",
			args: args{
				ctx:  ts.ctx,
				data: Userdata{},
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData4",
			args: args{
				ctx:  ts.ctx,
				data: data1,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData5",
			args: args{
				ctx:  ts.ctx,
				data: data2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			if err := ts.ops.WriteUserData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				ts.T().Errorf("DatabaseOps.WriteUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_ReadUserData() {
	type args struct {
		ctx            context.Context
		artistPastelID string
	}
	tests := []struct {
		name    string
		args    args
		want    Userdata
		wantErr bool
	}{
		{
			name: "TestDatabaseOps_ReadUserData1",
			args: args{
				ctx:            ts.ctx,
				artistPastelID: "qwe",
			},
			want:    data3,
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_ReadUserData2",
			args: args{
				ctx:            ts.ctx,
				artistPastelID: "rty",
			},
			want:    data4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			got, err := ts.ops.ReadUserData(tt.args.ctx, tt.args.artistPastelID)
			if (err != nil) != tt.wantErr {
				ts.T().Errorf("DatabaseOps.ReadUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				ts.T().Errorf("DatabaseOps.ReadUserData() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_UpdateUserData() {
	type args struct {
		ctx  context.Context
		data Userdata
	}

	data5.Realname = "dog"
	data5.Location = "us"
	data6.Categories = "c"
	data6.Biography = "d"

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestDatabaseOps_UpdateUserData1",
			args: args{
				ctx:  ts.ctx,
				data: data5,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_UpdateUserData2",
			args: args{
				ctx:  ts.ctx,
				data: data6,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			if err := ts.ops.UpdateUserData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				ts.T().Errorf("DatabaseOps.WriteUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
			data, err := ts.ops.ReadUserData(tt.args.ctx, tt.args.data.ArtistPastelID)
			if (err != nil) != tt.wantErr {
				ts.T().Errorf("DatabaseOps.WriteUserData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(data, tt.args.data) {
				ts.T().Errorf("DatabaseOps.ReadUserData() = %+v, want %+v", data, tt.args.data)
			}
		})
	}
}
