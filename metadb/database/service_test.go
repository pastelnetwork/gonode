package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
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
		user_data_hash TEXT NOT NULL,
		avatar_image BLOB,
		avatar_filename TEXT,
		cover_photo_image BLOB,
		cover_photo_filename TEXT
	);`

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

	data1 = pb.UserdataRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "abc",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data2 = pb.UserdataRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "xyz",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data3 = pb.UserdataRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "qwe",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data4 = pb.UserdataRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "rty",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}

	data3ReadResult = userdata.ProcessRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: userdata.UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: userdata.UserImageUpload{
			Content:  nil,
			Filename: "",
		},
		ArtistPastelID:    "qwe",
		Timestamp:         123,
		PreviousBlockHash: "hash",
	}

	data4ReadResult = userdata.ProcessRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: userdata.UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: userdata.UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "rty",
		Timestamp:         123,
		PreviousBlockHash: "hash",
	}

	userDataFrame = pb.UserdataRequest{
		Realname:        "cat",
		FacebookLink:    "fb.com",
		TwitterLink:     "tw.com",
		NativeCurrency:  "usd",
		Location:        "us",
		PrimaryLanguage: "en",
		Categories:      "a",
		Biography:       "b",
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{1, 2, 3, 4},
			Filename: "1234.jpg",
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  []byte{4, 5, 6, 7},
			Filename: "4567.jpg",
		},
		ArtistPastelID:    "rty",
		Timestamp:         123,
		Signature:         "xyz",
		PreviousBlockHash: "hash",
	}
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
	ops     *Ops
}

func (ts *testSuite) setupTableTest() {
	ts.Nil(ts.ops.WriteUserData(ts.ctx, &data3))
	ts.Nil(ts.ops.WriteUserData(ts.ctx, &data4))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art10_qwe", ArtistPastelID: "qwe", Copies: 2, CreatedTimestamp: 1}))

	for i := 1; i <= 10; i++ {
		userDataFrame.ArtistPastelID = fmt.Sprintf("id%d", i)
		ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))
	}

	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art1_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 1}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art2_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 2}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art1_id2", ArtistPastelID: "id2", Copies: 1, CreatedTimestamp: 2}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art2_id2", ArtistPastelID: "id2", Copies: 2, CreatedTimestamp: 3}))

	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art1_id1", ArtID: "art1_id1", Price: 10.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins2_art1_id1", ArtID: "art1_id1", Price: 20.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art2_id1", ArtID: "art2_id1", Price: 30.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins2_art2_id1", ArtID: "art2_id1", Price: 40.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art1_id2", ArtID: "art1_id2", Price: 50.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art2_id2", ArtID: "art2_id2", Price: 60.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins2_art2_id2", ArtID: "art2_id2", Price: 70.0}))

	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t1", InstanceID: "ins1_art1_id1", Timestamp: 20, SellerPastelID: "id1", BuyerPastelID: "id5", Price: 20.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t2", InstanceID: "ins1_art1_id1", Timestamp: 21, SellerPastelID: "id5", BuyerPastelID: "id6", Price: 40.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t3", InstanceID: "ins2_art2_id2", Timestamp: 23, SellerPastelID: "id2", BuyerPastelID: "id6", Price: 40.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t4", InstanceID: "ins2_art2_id2", Timestamp: 25, SellerPastelID: "id6", BuyerPastelID: "id7", Price: 45.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t5", InstanceID: "ins1_art1_id1", Timestamp: 26, SellerPastelID: "id6", BuyerPastelID: "id7", Price: 45.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t6", InstanceID: "ins1_art1_id2", Timestamp: 30, SellerPastelID: "id2", BuyerPastelID: "id4", Price: 20.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t7", InstanceID: "ins1_art1_id2", Timestamp: 31, SellerPastelID: "id4", BuyerPastelID: "id1", Price: 20.0}))
	ts.Nil(ts.ops.WriteTransaction(ts.ctx, ArtTransaction{TransactionID: "t8", InstanceID: "ins1_art1_id1", Timestamp: 32, SellerPastelID: "id7", BuyerPastelID: "id1", Price: 20.0}))

	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id1", FolloweePastelID: "id2"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id1", FolloweePastelID: "id3"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id1", FolloweePastelID: "id4"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id1", FolloweePastelID: "id5"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id1", FolloweePastelID: "id6"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id2", FolloweePastelID: "id1"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id2", FolloweePastelID: "id3"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id2", FolloweePastelID: "id4"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id2", FolloweePastelID: "id5"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id2", FolloweePastelID: "id6"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id3", FolloweePastelID: "id1"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id3", FolloweePastelID: "id2"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id3", FolloweePastelID: "id4"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id3", FolloweePastelID: "id5"}))
	ts.Nil(ts.ops.WriteUserFollow(ts.ctx, UserFollow{FollowerPastelID: "id3", FolloweePastelID: "id6"}))
}

func (ts *testSuite) SetupSuite() {
	workDir, err := ioutil.TempDir("", "servicetest-metadb-*")
	assert.NoError(ts.T(), err)

	config := metadb.NewConfig()
	config.SetWorkDir(workDir)
	config.HTTPPort = 4003
	config.RaftPort = 4004
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

	ts.Nil(db.EnableFKConstraints(true))

	content, err := ioutil.ReadFile("./commands/schema.sql")
	ts.Nil(err)
	listOfCommands := strings.Split(string(content), schemaDelimiter)
	for _, cmd := range listOfCommands {
		result, err := db.Write(ts.ctx, cmd)
		ts.Nil(err)
		assert.Equal(ts.T(), result.Error, "")
	}

	tmpls, err := NewTemplateKeeper("./commands")
	ts.Nil(err)
	ts.ops = &Ops{
		metaDB:    db,
		templates: tmpls,
	}

	ts.setupTableTest()
}

func (ts *testSuite) TearDownSuite() {
	if ts.cancel != nil {
		ts.cancel()
	}
	ts.wg.Wait()

	// remove the data directly
	os.RemoveAll(filepath.Join(ts.workDir))
}

func (ts *testSuite) TestDatabaseOps_WriteUserData() {
	type args struct {
		ctx  context.Context
		data *pb.UserdataRequest
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
				data: &data1,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData2",
			args: args{
				ctx:  ts.ctx,
				data: &data2,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData3",
			args: args{
				ctx:  ts.ctx,
				data: &pb.UserdataRequest{},
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData4",
			args: args{
				ctx:  ts.ctx,
				data: &data1,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData5",
			args: args{
				ctx:  ts.ctx,
				data: &data2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			if err := ts.ops.WriteUserData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				ts.T().Errorf("Ops.WriteUserData() error = %v, wantErr %v", err, tt.wantErr)
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
		want    userdata.ProcessRequest
		wantErr bool
	}{
		{
			name: "TestDatabaseOps_ReadUserData1",
			args: args{
				ctx:            ts.ctx,
				artistPastelID: "qwe",
			},
			want:    data3ReadResult,
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_ReadUserData2",
			args: args{
				ctx:            ts.ctx,
				artistPastelID: "rty",
			},
			want:    data4ReadResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			got, err := ts.ops.ReadUserData(tt.args.ctx, tt.args.artistPastelID)
			if (err != nil) != tt.wantErr {
				ts.T().Errorf("Ops.ReadUserData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				ts.T().Errorf("Ops.ReadUserData() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_WriteArtInfo() {
	tests := []struct {
		data    ArtInfo
		wantErr bool
	}{
		{
			data:    ArtInfo{ArtID: "art1_qwe", ArtistPastelID: "qwe", Copies: 2, CreatedTimestamp: 5},
			wantErr: false,
		},
		{
			data:    ArtInfo{ArtID: "art1_rty", ArtistPastelID: "rty", Copies: 2, CreatedTimestamp: 7},
			wantErr: false,
		},
		{
			data:    ArtInfo{ArtID: "art1_qwe", ArtistPastelID: "qwe", Copies: 2, CreatedTimestamp: 10},
			wantErr: true,
		},
		{
			data:    ArtInfo{ArtID: "art2_rty", ArtistPastelID: "dfkd-skfjsdk", Copies: 2, CreatedTimestamp: 11},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_WriteArtInfo-%d", i), func(t *testing.T) {
			if err := ts.ops.WriteArtInfo(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				fmt.Printf("%s/%s\n", fmt.Sprintf("TestDatabaseOps_WriteArtInfo-%d", i), err)
				ts.T().Errorf("Ops.WriteArtInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOps_WriteArtInstanceInfo(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx  context.Context
		data ArtInstanceInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			if err := db.WriteArtInstanceInfo(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteArtInstanceInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOps_WriteArtLike(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx  context.Context
		data ArtLike
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			if err := db.WriteArtLike(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteArtLike() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOps_WriteTransaction(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx  context.Context
		data ArtTransaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			if err := db.WriteTransaction(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOps_WriteUserFollow(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx  context.Context
		data UserFollow
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			if err := db.WriteUserFollow(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteUserFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabaseOps_salePriceByUserQuery(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx     context.Context
		command string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.salePriceByUserQuery(tt.args.ctx, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.salePriceByUserQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ops.salePriceByUserQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetCumulatedSalePriceByUser(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetCumulatedSalePriceByUser(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetCumulatedSalePriceByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ops.GetCumulatedSalePriceByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_queryPastelID(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx     context.Context
		command string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.queryPastelID(tt.args.ctx, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.queryPastelID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.queryPastelID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetFollowees(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetFollowees(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetFollowees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetFollowees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetFollowers(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetFollowers(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetFollowers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetFollowers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetFriends(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetFriends(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetFriends() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetFriends() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetHighestSalePriceByUser(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetHighestSalePriceByUser(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetHighestSalePriceByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ops.GetHighestSalePriceByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetExistingNftCopies(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx   context.Context
		artID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetExistingNftCopies(tt.args.ctx, tt.args.artID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetExistingNftCopies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ops.GetExistingNftCopies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_queryToInterface(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx     context.Context
		command string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.queryToInterface(tt.args.ctx, tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.queryToInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.queryToInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetNftCreatedByArtist(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx            context.Context
		artistPastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []NftCreatedByArtistQueryResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetNftCreatedByArtist(tt.args.ctx, tt.args.artistPastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetNftCreatedByArtist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetNftCreatedByArtist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetNftForSaleByArtist(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx            context.Context
		artistPastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []NftForSaleByArtistQueryResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetNftForSaleByArtist(tt.args.ctx, tt.args.artistPastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetNftForSaleByArtist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetNftForSaleByArtist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetNftOwnedByUser(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []NftOwnedByUserQueryResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetNftOwnedByUser(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetNftOwnedByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetNftOwnedByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetNftSoldByUser(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx      context.Context
		pastelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []NftSoldByUserQueryResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetNftSoldByUser(tt.args.ctx, tt.args.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetNftSoldByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetNftSoldByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetUniqueNftByUser(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx   context.Context
		query UniqueNftByUserQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []ArtInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetUniqueNftByUser(tt.args.ctx, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetUniqueNftByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetUniqueNftByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_GetUsersLikeNft(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx   context.Context
		artID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			got, err := db.GetUsersLikeNft(tt.args.ctx, tt.args.artID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetUsersLikeNft() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetUsersLikeNft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseOps_Run(t *testing.T) {
	type fields struct {
		metaDB    metadb.MetaDB
		templates *templateKeeper
		config    *Config
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Ops{
				metaDB:    tt.fields.metaDB,
				templates: tt.fields.templates,
				config:    tt.fields.config,
			}
			if err := db.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Ops.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDatabaseOps(t *testing.T) {
	type args struct {
		metaDB metadb.MetaDB
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *Ops
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDatabaseOps(tt.args.metaDB, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDatabaseOps() = %v, want %v", got, tt.want)
			}
		})
	}
}
