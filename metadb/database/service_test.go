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

	for i := 1; i <= 10; i++ {
		userDataFrame.ArtistPastelID = fmt.Sprintf("id%d", i)
		ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))
	}

	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art1_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 5, GreenNft: true, RarenessScore: 0.9, RoyaltyRatePercentage: 10.0}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art2_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 10, GreenNft: false, RarenessScore: 0.8, RoyaltyRatePercentage: 9.0}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art1_id2", ArtistPastelID: "id2", Copies: 1, CreatedTimestamp: 15, GreenNft: true, RarenessScore: 0.7, RoyaltyRatePercentage: 8.0}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art2_id2", ArtistPastelID: "id2", Copies: 2, CreatedTimestamp: 20, GreenNft: false, RarenessScore: 0.6, RoyaltyRatePercentage: 7.0}))

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

	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art1_id1", PastelID: "id2"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art1_id1", PastelID: "id3"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art1_id1", PastelID: "id4"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art1_id1", PastelID: "id5"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art2_id1", PastelID: "id6"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art2_id1", PastelID: "id7"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art1_id2", PastelID: "id7"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art2_id2", PastelID: "id8"}))
	ts.Nil(ts.ops.WriteArtLike(ts.ctx, ArtLike{ArtID: "art2_id2", PastelID: "id9"}))
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
				data: &data1,
			},
			wantErr: false,
		},
		{
			name: "TestDatabaseOps_WriteUserData4",
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
				ts.T().Errorf("Ops.WriteArtInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_WriteArtInstanceInfo() {
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art10_qwe", ArtistPastelID: "qwe", Copies: 2, CreatedTimestamp: 1}))

	askingPrice := 20.0
	tests := []struct {
		data    ArtInstanceInfo
		wantErr bool
	}{
		{
			data:    ArtInstanceInfo{InstanceID: "ins1_art10_qwe", ArtID: "art10_qwe", Price: 10.0, AskingPrice: &askingPrice},
			wantErr: false,
		},
		{
			data:    ArtInstanceInfo{InstanceID: "ins1_art10_qwe", ArtID: "art10_qwe", Price: 10.0, AskingPrice: &askingPrice},
			wantErr: true,
		},
		{
			data:    ArtInstanceInfo{InstanceID: "ins2_art10_qwe", ArtID: "art10_qwe", Price: 10.0},
			wantErr: false,
		},
		{
			data:    ArtInstanceInfo{InstanceID: "ins3_art10_qwe", ArtID: "art10_qwe", Price: 10.0},
			wantErr: true,
		},
		{
			data:    ArtInstanceInfo{InstanceID: "ins4_art10_qwe", ArtID: "art10_qwe_fff", Price: 10.0},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_WriteArtInstanceInfo-%d", i), func(t *testing.T) {
			err := ts.ops.WriteArtInstanceInfo(ts.ctx, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteArtInstanceInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				command := fmt.Sprintf(`select owner_pastel_id, asking_price from art_instance_metadata where instance_id = '%s'`, tt.data.InstanceID)
				queryResult, err := ts.ops.metaDB.Query(ts.ctx, command, queryLevelNone)
				ts.Nil(err)
				ts.Equal(int64(1), queryResult.NumRows())
				var owner string
				var askingPrice *float64
				queryResult.Next()
				queryResult.Scan(&owner, askingPrice)
				ts.Equal("qwe", owner)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_WriteArtLike() {
	userDataFrame.ArtistPastelID = "like"
	ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art1_like", ArtistPastelID: "like", Copies: 2, CreatedTimestamp: 12}))

	tests := []struct {
		data    ArtLike
		wantErr bool
	}{
		{
			data:    ArtLike{ArtID: "art1_like", PastelID: "id1"},
			wantErr: false,
		},
		{
			data:    ArtLike{ArtID: "art1_like", PastelID: "id2"},
			wantErr: false,
		},
		{
			data:    ArtLike{ArtID: "art1_like", PastelID: "id3"},
			wantErr: false,
		},
		{
			data:    ArtLike{ArtID: "art1_like", PastelID: "id3"},
			wantErr: true,
		},
		{
			data:    ArtLike{ArtID: "art1_like_213", PastelID: "id3"},
			wantErr: true,
		},
		{
			data:    ArtLike{ArtID: "art1_like", PastelID: "id3231"},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_WriteArtLike-%d", i), func(t *testing.T) {
			if err := ts.ops.WriteArtLike(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteArtLike() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_WriteTransaction() {
	for i := 0; i < 3; i++ {
		userDataFrame.ArtistPastelID = fmt.Sprintf("id_transaction_%d", i)
		ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))
	}

	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art0_id_transaction_0", ArtistPastelID: "id_transaction_0", Copies: 1, CreatedTimestamp: 15}))
	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art0_id_transaction_1", ArtistPastelID: "id_transaction_1", Copies: 2, CreatedTimestamp: 15}))

	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins0_art0_id_transaction_0", ArtID: "art0_id_transaction_0", Price: 10.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins0_art0_id_transaction_1", ArtID: "art0_id_transaction_1", Price: 10.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art0_id_transaction_1", ArtID: "art0_id_transaction_1", Price: 10.0}))

	tests := []struct {
		data    ArtTransaction
		wantErr bool
	}{
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_0",
				InstanceID:     "ins0_art0_id_transaction_0",
				Timestamp:      100,
				SellerPastelID: "id_transaction_0",
				BuyerPastelID:  "id_transaction_2",
				Price:          10.0,
			},
			wantErr: false,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_1",
				InstanceID:     "ins0_art0_id_transaction_1",
				Timestamp:      101,
				SellerPastelID: "id_transaction_1",
				BuyerPastelID:  "id_transaction_2",
				Price:          10.0,
			},
			wantErr: false,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_2",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      102,
				SellerPastelID: "id_transaction_1",
				BuyerPastelID:  "id_transaction_2",
				Price:          10.0,
			},
			wantErr: false,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_0",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      103,
				SellerPastelID: "id_transaction_2",
				BuyerPastelID:  "id_transaction_1",
				Price:          10.0,
			},
			wantErr: true,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_3",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      104,
				SellerPastelID: "id_transaction_1",
				BuyerPastelID:  "id_transaction_2",
				Price:          10.0,
			},
			wantErr: true,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_4",
				InstanceID:     "sdasd",
				Timestamp:      105,
				SellerPastelID: "id_transaction_2",
				BuyerPastelID:  "id_transaction_1",
				Price:          10.0,
			},
			wantErr: true,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_4",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      106,
				SellerPastelID: "adsd",
				BuyerPastelID:  "id_transaction_1",
				Price:          10.0,
			},
			wantErr: true,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_4",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      105,
				SellerPastelID: "id_transaction_2",
				BuyerPastelID:  "1231",
				Price:          10.0,
			},
			wantErr: true,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_4",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      105,
				SellerPastelID: "id_transaction_2",
				BuyerPastelID:  "id_transaction_1",
				Price:          10.0,
			},
			wantErr: false,
		},
		{
			data: ArtTransaction{
				TransactionID:  "test_trans_5",
				InstanceID:     "ins1_art0_id_transaction_1",
				Timestamp:      106,
				SellerPastelID: "id_transaction_2",
				BuyerPastelID:  "id_transaction_1",
				Price:          20.0,
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_WriteTransaction-%d", i), func(t *testing.T) {
			if err := ts.ops.WriteTransaction(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_WriteUserFollow() {
	for i := 0; i < 3; i++ {
		userDataFrame.ArtistPastelID = fmt.Sprintf("id_write_follow_%d", i)
		ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))
	}

	tests := []struct {
		data    UserFollow
		wantErr bool
	}{
		{data: UserFollow{FollowerPastelID: "id_write_follow_0", FolloweePastelID: "id_write_follow_1"}, wantErr: false},
		{data: UserFollow{FollowerPastelID: "id_write_follow_0", FolloweePastelID: "id_write_follow_1"}, wantErr: true},
		{data: UserFollow{FollowerPastelID: "id_write_follow_1", FolloweePastelID: "id_write_follow_1"}, wantErr: true},
		{data: UserFollow{FollowerPastelID: "id_write_follow_1", FolloweePastelID: "id_write_follow_0"}, wantErr: false},
		{data: UserFollow{FollowerPastelID: "id_write_follow_1", FolloweePastelID: "id_write_follow_2"}, wantErr: false},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_WriteUserFollow-%d", i), func(t *testing.T) {
			if err := ts.ops.WriteUserFollow(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.WriteUserFollow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_GetCumulatedSalePriceByUser() {
	tests := []struct {
		pastelID string
		want     float64
		wantErr  bool
	}{
		{
			pastelID: "",
			want:     0,
			wantErr:  true,
		},
		{
			pastelID: "id8",
			want:     0,
			wantErr:  false,
		},
		{
			pastelID: "id1",
			want:     20.0,
			wantErr:  false,
		},
		{
			pastelID: "id6",
			want:     90.0,
			wantErr:  false,
		},
		{
			pastelID: "id2",
			want:     60.0,
			wantErr:  false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetCumulatedSalePriceByUser-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetCumulatedSalePriceByUser(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetFollowees() {
	tests := []struct {
		pastelID string
		want     []string
		wantErr  bool
	}{
		{
			pastelID: "id1",
			want:     []string{"id2", "id3", "id4", "id5", "id6"},
			wantErr:  false,
		},
		{
			pastelID: "id2",
			want:     []string{"id1", "id3", "id4", "id5", "id6"},
			wantErr:  false,
		},
		{
			pastelID: "id3",
			want:     []string{"id1", "id2", "id4", "id5", "id6"},
			wantErr:  false,
		},
		{
			pastelID: "id4",
			want:     []string{},
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     nil,
			wantErr:  true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetFollowees-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetFollowees(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetFollowers() {
	tests := []struct {
		pastelID string
		want     []string
		wantErr  bool
	}{
		{
			pastelID: "id1",
			want:     []string{"id2", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id2",
			want:     []string{"id1", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id3",
			want:     []string{"id1", "id2"},
			wantErr:  false,
		},
		{
			pastelID: "id4",
			want:     []string{"id1", "id2", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id5",
			want:     []string{"id1", "id2", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id6",
			want:     []string{"id1", "id2", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id7",
			want:     []string{},
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     nil,
			wantErr:  true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetFollowers-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetFollowers(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetFriends() {
	tests := []struct {
		pastelID string
		want     []string
		wantErr  bool
	}{
		{
			pastelID: "id1",
			want:     []string{"id2", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id2",
			want:     []string{"id1", "id3"},
			wantErr:  false,
		},
		{
			pastelID: "id3",
			want:     []string{"id1", "id2"},
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     nil,
			wantErr:  true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetFriends-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetFriends(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetHighestSalePriceByUser() {
	tests := []struct {
		pastelID string
		want     float64
		wantErr  bool
	}{
		{
			pastelID: "id2",
			want:     40.0,
			wantErr:  false,
		},
		{
			pastelID: "id6",
			want:     45.0,
			wantErr:  false,
		},
		{
			pastelID: "id10",
			want:     0.0,
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     0.0,
			wantErr:  true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetHighestSalePriceByUser-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetHighestSalePriceByUser(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetExistingNftCopies() {
	tests := []struct {
		artID   string
		want    int
		wantErr bool
	}{
		{
			artID:   "art1_id1",
			want:    2,
			wantErr: false,
		},
		{
			artID:   "art2_id1",
			want:    2,
			wantErr: false,
		},
		{
			artID:   "art1_id2",
			want:    1,
			wantErr: false,
		},
		{
			artID:   "art2_id2",
			want:    2,
			wantErr: false,
		},
		{
			artID:   "s2123",
			want:    0,
			wantErr: true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetExistingNftCopies-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetExistingNftCopies(ts.ctx, tt.artID)
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

func (ts *testSuite) TestDatabaseOps_GetNftCreatedByArtist() {
	tests := []struct {
		artistPastelID string
		want           []NftCreatedByArtistQueryResult
		wantErr        bool
	}{
		{
			artistPastelID: "id1",
			want: []NftCreatedByArtistQueryResult{
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins1_art1_id1",
					ArtID:                 "art1_id1",
					Copies:                2,
					CreatedTimestamp:      5,
					GreenNft:              true,
					RarenessScore:         0.9,
					RoyaltyRatePercentage: 10.0,
				},
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins1_art2_id1",
					ArtID:                 "art2_id1",
					Copies:                2,
					CreatedTimestamp:      10,
					GreenNft:              false,
					RarenessScore:         0.8,
					RoyaltyRatePercentage: 9.0,
				},
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins2_art1_id1",
					ArtID:                 "art1_id1",
					Copies:                2,
					CreatedTimestamp:      5,
					GreenNft:              true,
					RarenessScore:         0.9,
					RoyaltyRatePercentage: 10.0,
				},
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins2_art2_id1",
					ArtID:                 "art2_id1",
					Copies:                2,
					CreatedTimestamp:      10,
					GreenNft:              false,
					RarenessScore:         0.8,
					RoyaltyRatePercentage: 9.0,
				},
			},
			wantErr: false,
		},
		{
			artistPastelID: "id2",
			want: []NftCreatedByArtistQueryResult{
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins1_art1_id2",
					ArtID:                 "art1_id2",
					Copies:                1,
					CreatedTimestamp:      15,
					GreenNft:              true,
					RarenessScore:         0.7,
					RoyaltyRatePercentage: 8.0,
				},
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins1_art2_id2",
					ArtID:                 "art2_id2",
					Copies:                2,
					CreatedTimestamp:      20,
					GreenNft:              false,
					RarenessScore:         0.6,
					RoyaltyRatePercentage: 7.0,
				},
				NftCreatedByArtistQueryResult{
					InstanceID:            "ins2_art2_id2",
					ArtID:                 "art2_id2",
					Copies:                2,
					CreatedTimestamp:      20,
					GreenNft:              false,
					RarenessScore:         0.6,
					RoyaltyRatePercentage: 7.0,
				},
			},
			wantErr: false,
		},
		{
			artistPastelID: "id3",
			want:           []NftCreatedByArtistQueryResult{},
			wantErr:        false,
		},
		{
			artistPastelID: "",
			want:           nil,
			wantErr:        true,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetNftCreatedByArtist-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetNftCreatedByArtist(ts.ctx, tt.artistPastelID)
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

func (ts *testSuite) TestDatabaseOps_GetNftForSaleByArtist() {
	tests := []struct {
		artistPastelID string
		want           []NftForSaleByArtistQueryResult
		wantErr        bool
	}{
		{
			artistPastelID: "id1",
			want: []NftForSaleByArtistQueryResult{
				NftForSaleByArtistQueryResult{
					InstanceID: "ins1_art2_id1",
					ArtID:      "art2_id1",
					Price:      30.0,
				},
				NftForSaleByArtistQueryResult{
					InstanceID: "ins2_art1_id1",
					ArtID:      "art1_id1",
					Price:      20.0,
				},
				NftForSaleByArtistQueryResult{
					InstanceID: "ins2_art2_id1",
					ArtID:      "art2_id1",
					Price:      40.0,
				},
			},
			wantErr: false,
		},
		{
			artistPastelID: "id2",
			want: []NftForSaleByArtistQueryResult{
				NftForSaleByArtistQueryResult{
					InstanceID: "ins1_art2_id2",
					ArtID:      "art2_id2",
					Price:      60.0,
				},
			},
			wantErr: false,
		},
		{
			artistPastelID: "",
			want:           nil,
			wantErr:        true,
		},
		{
			artistPastelID: "id10",
			want:           []NftForSaleByArtistQueryResult{},
			wantErr:        false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetNftForSaleByArtist-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetNftForSaleByArtist(ts.ctx, tt.artistPastelID)
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

func (ts *testSuite) TestDatabaseOps_GetNftOwnedByUser() {
	tests := []struct {
		pastelID string
		want     []NftOwnedByUserQueryResult
		wantErr  bool
	}{
		{
			pastelID: "id10",
			want:     []NftOwnedByUserQueryResult{},
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     nil,
			wantErr:  true,
		},
		{
			pastelID: "id7",
			want: []NftOwnedByUserQueryResult{
				NftOwnedByUserQueryResult{
					ArtID: "art1_id1",
					Count: 1,
				},
				NftOwnedByUserQueryResult{
					ArtID: "art2_id2",
					Count: 1,
				},
			},
			wantErr: false,
		},
		{
			pastelID: "id1",
			want: []NftOwnedByUserQueryResult{
				NftOwnedByUserQueryResult{
					ArtID: "art1_id1",
					Count: 1,
				},
				NftOwnedByUserQueryResult{
					ArtID: "art1_id2",
					Count: 1,
				},
				NftOwnedByUserQueryResult{
					ArtID: "art2_id1",
					Count: 2,
				},
			},
			wantErr: false,
		},
		{
			pastelID: "id2",
			want: []NftOwnedByUserQueryResult{
				NftOwnedByUserQueryResult{
					ArtID: "art2_id2",
					Count: 1,
				},
			},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetNftOwnedByUser-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetNftOwnedByUser(ts.ctx, tt.pastelID)
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

func (ts *testSuite) TestDatabaseOps_GetNftSoldByArtID() {
	tests := []struct {
		pastelID string
		want     NftSoldByArtIDQueryResult
		wantErr  bool
	}{
		{
			pastelID: "dadsa",
			want:     NftSoldByArtIDQueryResult{},
			wantErr:  false,
		},
		{
			pastelID: "",
			want:     NftSoldByArtIDQueryResult{},
			wantErr:  true,
		},
		{
			pastelID: "art1_id1",
			want: NftSoldByArtIDQueryResult{
				TotalCopies: 2,
				SoldCopies:  1,
			},
			wantErr: false,
		},
		{
			pastelID: "art2_id1",
			want: NftSoldByArtIDQueryResult{
				TotalCopies: 2,
				SoldCopies:  0,
			},
			wantErr: false,
		},
		{
			pastelID: "art1_id2",
			want: NftSoldByArtIDQueryResult{
				TotalCopies: 1,
				SoldCopies:  1,
			},
			wantErr: false,
		},
		{
			pastelID: "art2_id2",
			want: NftSoldByArtIDQueryResult{
				TotalCopies: 2,
				SoldCopies:  1,
			},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetNftSoldByArtID-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetNftSoldByArtID(ts.ctx, tt.pastelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetNftSoldByArtID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetNftSoldByArtID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (ts *testSuite) TestDatabaseOps_GetUniqueNftByUser() {
	tests := []struct {
		query   UniqueNftByUserQuery
		want    []ArtInfo
		wantErr bool
	}{
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id1",
				LimitTimestamp: 7,
			},
			want: []ArtInfo{
				ArtInfo{ArtID: "art2_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 10, GreenNft: false, RarenessScore: 0.8, RoyaltyRatePercentage: 9.0},
			},
			wantErr: false,
		},
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id1",
				LimitTimestamp: 1,
			},
			want: []ArtInfo{
				ArtInfo{ArtID: "art2_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 10, GreenNft: false, RarenessScore: 0.8, RoyaltyRatePercentage: 9.0},
				ArtInfo{ArtID: "art1_id1", ArtistPastelID: "id1", Copies: 2, CreatedTimestamp: 5, GreenNft: true, RarenessScore: 0.9, RoyaltyRatePercentage: 10.0},
			},
			wantErr: false,
		},
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id1",
				LimitTimestamp: 11,
			},
			want:    []ArtInfo{},
			wantErr: false,
		},
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id2",
				LimitTimestamp: 11,
			},
			want: []ArtInfo{
				ArtInfo{ArtID: "art2_id2", ArtistPastelID: "id2", Copies: 2, CreatedTimestamp: 20, GreenNft: false, RarenessScore: 0.6, RoyaltyRatePercentage: 7.0},
				ArtInfo{ArtID: "art1_id2", ArtistPastelID: "id2", Copies: 1, CreatedTimestamp: 15, GreenNft: true, RarenessScore: 0.7, RoyaltyRatePercentage: 8.0},
			},
			wantErr: false,
		},
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id2",
				LimitTimestamp: 17,
			},
			want: []ArtInfo{
				ArtInfo{ArtID: "art2_id2", ArtistPastelID: "id2", Copies: 2, CreatedTimestamp: 20, GreenNft: false, RarenessScore: 0.6, RoyaltyRatePercentage: 7.0},
			},
			wantErr: false,
		},
		{
			query: UniqueNftByUserQuery{
				ArtistPastelID: "id2",
				LimitTimestamp: 30,
			},
			want:    []ArtInfo{},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetUniqueNftByUser-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetUniqueNftByUser(ts.ctx, tt.query)
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

func (ts *testSuite) TestDatabaseOps_GetUsersLikeNft() {
	tests := []struct {
		artID   string
		want    []string
		wantErr bool
	}{
		{
			artID:   "id10",
			want:    []string{},
			wantErr: false,
		},
		{
			artID:   "",
			want:    nil,
			wantErr: true,
		},
		{
			artID:   "art1_id1",
			want:    []string{"id2", "id3", "id4", "id5"},
			wantErr: false,
		},
		{
			artID:   "art2_id1",
			want:    []string{"id6", "id7"},
			wantErr: false,
		},
		{
			artID:   "art1_id2",
			want:    []string{"id7"},
			wantErr: false,
		},
		{
			artID:   "art2_id2",
			want:    []string{"id8", "id9"},
			wantErr: false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestDatabaseOps_GetUsersLikeNft-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetUsersLikeNft(ts.ctx, tt.artID)
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

func (ts *testSuite) TestOps_UpdateAskingPrice() {
	userDataFrame.ArtistPastelID = "id_ap_0"
	ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))

	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art0_id_ap_0", ArtistPastelID: "id_ap_0", Copies: 2, CreatedTimestamp: 15}))

	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins0_art0_id_ap_0", ArtID: "art0_id_ap_0", OwnerPastelID: "ble", Price: 10.0}))
	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art0_id_ap_0", ArtID: "art0_id_ap_0", OwnerPastelID: "ble", Price: 10.0}))

	tests := []struct {
		data          AskingPriceUpdateRequest
		ownerPastelID string
		wantErr       bool
	}{
		{
			data: AskingPriceUpdateRequest{
				InstanceID:  "ins0_art0_id_ap_0",
				AskingPrice: 10.0,
			},
			ownerPastelID: "id_ap_0",
			wantErr:       false,
		},
		{
			data: AskingPriceUpdateRequest{
				InstanceID:  "ins1_art0_id_ap_0",
				AskingPrice: 20.0,
			},
			ownerPastelID: "id_ap_0",
			wantErr:       false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_UpdateAskingPrice-%d", i), func(t *testing.T) {
			if err := ts.ops.UpdateAskingPrice(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.UpdateAskingPrice() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				info, err := ts.ops.GetArtInstanceInfo(ts.ctx, tt.data.InstanceID)
				ts.Nil(err)
				ts.Equal(*info.AskingPrice, tt.data.AskingPrice)
				ts.Equal(info.OwnerPastelID, "id_ap_0")
			}
		})
	}
}

func (ts *testSuite) TestOps_GetArtInstanceInfo() {
	userDataFrame.ArtistPastelID = "id_ai_0"
	ts.Nil(ts.ops.WriteUserData(ts.ctx, &userDataFrame))

	ts.Nil(ts.ops.WriteArtInfo(ts.ctx, ArtInfo{ArtID: "art0_id_ai_0", ArtistPastelID: "id_ai_0", Copies: 2, CreatedTimestamp: 150}))

	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins0_art0_id_ai_0", ArtID: "art0_id_ai_0", Price: 10.0}))
	ts.Nil(ts.ops.UpdateAskingPrice(ts.ctx, AskingPriceUpdateRequest{InstanceID: "ins0_art0_id_ai_0", AskingPrice: 20.0}))

	ts.Nil(ts.ops.WriteArtInstanceInfo(ts.ctx, ArtInstanceInfo{InstanceID: "ins1_art0_id_ai_0", ArtID: "art0_id_ai_0", Price: 10.0}))

	price20 := 20.0
	tests := []struct {
		instanceID string
		want       ArtInstanceInfo
		wantErr    bool
	}{
		{
			instanceID: "231dsq",
			want:       ArtInstanceInfo{},
			wantErr:    true,
		},
		{
			instanceID: "",
			want:       ArtInstanceInfo{},
			wantErr:    true,
		},
		{
			instanceID: "ins0_art0_id_ai_0",
			want:       ArtInstanceInfo{InstanceID: "ins0_art0_id_ai_0", ArtID: "art0_id_ai_0", OwnerPastelID: "id_ai_0", Price: 10.0, AskingPrice: &price20},
			wantErr:    false,
		},
		{
			instanceID: "ins1_art0_id_ai_0",
			want:       ArtInstanceInfo{InstanceID: "ins1_art0_id_ai_0", ArtID: "art0_id_ai_0", OwnerPastelID: "id_ai_0", Price: 10.0, AskingPrice: nil},
			wantErr:    false,
		},
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_GetArtInstanceInfo-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetArtInstanceInfo(ts.ctx, tt.instanceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetArtInstanceInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetArtInstanceInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (ts *testSuite) TestOps_EndArtAuction() {
	tests := []struct {
		auctionID int64
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_EndArtAuction-%d", i), func(t *testing.T) {
			if err := ts.ops.EndArtAuction(ts.ctx, tt.auctionID); (err != nil) != tt.wantErr {
				t.Errorf("Ops.EndArtAuction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestOps_ArtPlaceBid() {
	tests := []struct {
		data    ArtPlaceBidRequest
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_ArtPlaceBid-%d", i), func(t *testing.T) {
			if err := ts.ops.ArtPlaceBid(ts.ctx, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Ops.ArtPlaceBid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (ts *testSuite) TestOps_GetAuctionInfo() {
	tests := []struct {
		auctionID int64
		want      ArtAuctionInfo
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_GetAuctionInfo-%d", i), func(t *testing.T) {
			got, err := ts.ops.GetAuctionInfo(ts.ctx, tt.auctionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.GetAuctionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ops.GetAuctionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (ts *testSuite) TestOps_NewArtAuction() {
	tests := []struct {
		data    NewArtAuctionRequest
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for i, tt := range tests {
		ts.T().Run(fmt.Sprintf("TestOps_NewArtAuction-%d", i), func(t *testing.T) {
			got, err := ts.ops.NewArtAuction(ts.ctx, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ops.NewArtAuction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Ops.NewArtAuction() = %v, want %v", got, tt.want)
			}
		})
	}
}
