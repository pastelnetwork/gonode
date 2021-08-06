package database

import (
	"context"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
)

const (
	logPrefix       = "database"
	queryLevelNone  = "none"
	schemaDelimiter = "---"

	userInfoWriteTemplate            = "user_info_write"
	userInfoQueryTemplate            = "user_info_query"
	artInfoWriteTemplate             = "art_info_write"
	artInstanceInfoWriteTemplate     = "art_instance_info_write"
	artLikeWriteTemplate             = "art_like_write"
	cumulatedSalePriceByUserTemplate = "cumulated_sale_price_by_user"
	getFolloweesTemplate             = "get_followees"
	getFollowersTemplate             = "get_followers"
	getFriendTemplate                = "get_friend"
	highestSalePriceByUserTemplate   = "highest_sale_price_by_user"
	nftCopiesExistTemplate           = "nft_copies_exist"
	nftCreatedByArtistTemplate       = "nft_created_by_artist"
	nftForSaleByArtistTemplate       = "nft_for_sale_by_artist"
	nftOwnedByUserTemplate           = "nft_owned_by_user"
	nftSoldByUserTemplate            = "nft_sold_by_user"
	transactionWriteTemplate         = "transaction_write"
	uniqueNftByUserTemplate          = "unique_nft_by_user"
	userFollowWriteTemplate          = "user_follow_write"
	usersLikeNftTemplate             = "users_like_nft"
)

type Config struct {
	SchemaPath   string `mapstructure:"schema-path" json:"schema-path,omitempty"`
	TemplatePath string `mapstructure:"template-path" json:"write-template-path,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

type DatabaseOps struct {
	metaDB    metadb.MetaDB
	templates *templateKeeper
	config    *Config
}

func (db *DatabaseOps) IsLeader() bool {
	return db.metaDB.IsLeader()
}

func (db *DatabaseOps) LeaderAddress() string {
	re := regexp.MustCompile(":[0-9]+$")
	address := re.Split(db.metaDB.LeaderAddress(), -1)[0]
	return address
}

func (db *DatabaseOps) writeData(ctx context.Context, command string) error {
	if len(command) == 0 {
		return errors.Errorf("invalid command")
	}
	result, err := db.metaDB.Write(ctx, command)
	if err != nil {
		return errors.Errorf("error while writting to db: %w", err)
	}
	if result.Error != "" {
		return errors.Errorf("error while writting to db: %s", result.Error)
	}

	return nil
}

func (db *DatabaseOps) WriteArtInfo(ctx context.Context, data ArtInfo) error {
	command, err := db.templates.GetCommand(artInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return errors.Errorf("error while writing art info %w", db.writeData(ctx, command))
}

func (db *DatabaseOps) WriteArtInstanceInfo(ctx context.Context, data ArtInstanceInfo) error {
	command, err := db.templates.GetCommand(artInstanceInfoWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return errors.Errorf("error while writing art info %w", db.writeData(ctx, command))
}

func (db *DatabaseOps) WriteArtLike(ctx context.Context, data ArtLike) error {
	command, err := db.templates.GetCommand(artLikeWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return errors.Errorf("error while writing art info %w", db.writeData(ctx, command))
}

func (db *DatabaseOps) WriteTransaction(ctx context.Context, data ArtTransaction) error {
	command, err := db.templates.GetCommand(artLikeWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return errors.Errorf("error while writing art info %w", db.writeData(ctx, command))
}

func (db *DatabaseOps) WriteUserFollow(ctx context.Context, data UserFollow) error {
	command, err := db.templates.GetCommand(userFollowWriteTemplate, data)
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}
	return errors.Errorf("error while writing art info %w", db.writeData(ctx, command))
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) WriteUserData(ctx context.Context, data *pb.UserdataRequest) error {
	if data == nil {
		return errors.Errorf("input nil data")
	}
	command, err := db.templates.GetCommand(userInfoWriteTemplate, pbToWriteCommand(data))
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
	}

	return db.writeData(ctx, command)
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) ReadUserData(ctx context.Context, artistPastelID string) (userdata.ProcessRequest, error) {
	command, err := db.templates.GetCommand(userInfoQueryTemplate, artistPastelID)
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while subtitute template: %w", err)
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while querying db: %w", err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return userdata.ProcessRequest{}, errors.Errorf("no artist with pastel id = %s", artistPastelID)
	} else if nrows > 1 {
		return userdata.ProcessRequest{}, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()
	resultMap, err := queryResult.Map()
	if err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while extracting result: %w", err)
	}

	var dbResult UserdataReadResult
	if err := mapstructure.Decode(resultMap, &dbResult); err != nil {
		return userdata.ProcessRequest{}, errors.Errorf("error while decoding result: %w", err)
	}

	return dbResult.ToUserData(), nil
}

func (db *DatabaseOps) salePriceByUserQuery(ctx context.Context, command string) (float64, error) {
	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return 0.0, errors.Errorf("error while querying db: %w", err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return 0.0, errors.Errorf("pastel id does not exist or no transaction have been made, query: %s", command)
	} else if nrows > 1 {
		return 0.0, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()

	var price float64
	if err := queryResult.Scan(&price); err != nil {
		return 0.0, errors.Errorf("error occurs while scanning query result: %w", err)
	}
	return price, nil
}

func (db *DatabaseOps) GetCumulatedSalePriceByUser(ctx context.Context, pastelID string) (float64, error) {
	if pastelID == "" {
		return 0.0, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(cumulatedSalePriceByUserTemplate, pastelID)
	if err != nil {
		return 0.0, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.salePriceByUserQuery(ctx, command)
}

func (db *DatabaseOps) queryPastelID(ctx context.Context, command string) ([]string, error) {
	if len(command) == 0 {
		return nil, errors.Errorf("invalid command")
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}

	result := make([]string, 0)
	for queryResult.Next() {
		var id string
		if err := queryResult.Scan(&id); err != nil {
			return nil, errors.Errorf("error while scaning db result: %w", err)
		}
		result = append(result, id)
	}

	return result, nil
}

func (db *DatabaseOps) GetFollowees(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFolloweesTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *DatabaseOps) GetFollowers(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFollowersTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *DatabaseOps) GetFriends(ctx context.Context, pastelID string) ([]string, error) {
	if pastelID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(getFriendTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *DatabaseOps) GetHighestSalePriceByUser(ctx context.Context, pastelID string) (float64, error) {
	if pastelID == "" {
		return 0.0, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(highestSalePriceByUserTemplate, pastelID)
	if err != nil {
		return 0.0, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.salePriceByUserQuery(ctx, command)
}

func (db *DatabaseOps) GetExistingNftCopies(ctx context.Context, artID string) (int, error) {
	if artID == "" {
		return 0, errors.Errorf("invalid art ID")
	}

	command, err := db.templates.GetCommand(nftCopiesExistTemplate, artID)
	if err != nil {
		return 0, errors.Errorf("error while subtitute template: %w", err)
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return 0, errors.Errorf("error while querying db: %w", err)
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return 0, errors.Errorf("art id does not exist or no transaction have been made: %s", artID)
	} else if nrows > 1 {
		return 0, errors.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()

	var copies int
	if err := queryResult.Scan(&copies); err != nil {
		return 0, errors.Errorf("error occurs while scanning query result: %w", err)
	}
	return copies, nil
}

func (db *DatabaseOps) queryToInterface(ctx context.Context, command string) ([]map[string]interface{}, error) {
	if len(command) == 0 {
		return nil, errors.Errorf("invalid command")
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelNone)
	if err != nil {
		return nil, errors.Errorf("error while querying db: %w", err)
	}

	result := make([]map[string]interface{}, 0)
	for queryResult.Next() {
		resultMap, err := queryResult.Map()
		if err != nil {
			return nil, errors.Errorf("error while extracting result: %w", err)
		}
		result = append(result, resultMap)
	}

	return result, nil
}

func (db *DatabaseOps) GetNftCreatedByArtist(ctx context.Context, artistPastelID string) ([]NftCreatedByArtistQueryResult, error) {
	command, err := db.templates.GetCommand(nftCopiesExistTemplate, artistPastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]NftCreatedByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record NftCreatedByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *DatabaseOps) GetNftForSaleByArtist(ctx context.Context, artistPastelID string) ([]NftForSaleByArtistQueryResult, error) {
	command, err := db.templates.GetCommand(nftForSaleByArtistTemplate, artistPastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]NftForSaleByArtistQueryResult, 0)
	for _, mp := range allResults {
		var record NftForSaleByArtistQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *DatabaseOps) GetNftOwnedByUser(ctx context.Context, pastelID string) ([]NftOwnedByUserQueryResult, error) {
	command, err := db.templates.GetCommand(nftOwnedByUserTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]NftOwnedByUserQueryResult, 0)
	for _, mp := range allResults {
		var record NftOwnedByUserQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *DatabaseOps) GetNftSoldByUser(ctx context.Context, pastelID string) ([]NftSoldByUserQueryResult, error) {
	command, err := db.templates.GetCommand(nftSoldByUserTemplate, pastelID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]NftSoldByUserQueryResult, 0)
	for _, mp := range allResults {
		var record NftSoldByUserQueryResult
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *DatabaseOps) GetUniqueNftByUser(ctx context.Context, query UniqueNftByUserQuery) ([]ArtInfo, error) {
	command, err := db.templates.GetCommand(uniqueNftByUserTemplate, query)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	allResults, err := db.queryToInterface(ctx, command)
	if err != nil {
		return nil, errors.Errorf("error while query db: %w", err)
	}

	result := make([]ArtInfo, 0)
	for _, mp := range allResults {
		var record ArtInfo
		if err := mapstructure.Decode(mp, &record); err != nil {
			return nil, errors.Errorf("error while decoding result: %w", err)
		}
		result = append(result, record)
	}

	return result, nil
}

func (db *DatabaseOps) GetUsersLikeNft(ctx context.Context, artID string) ([]string, error) {
	if artID == "" {
		return nil, errors.Errorf("invalid pastel ID")
	}

	command, err := db.templates.GetCommand(usersLikeNftTemplate, artID)
	if err != nil {
		return nil, errors.Errorf("error while subtitute template: %w", err)
	}

	return db.queryPastelID(ctx, command)
}

func (db *DatabaseOps) Run(ctx context.Context) error {

	ctx = log.ContextWithPrefix(ctx, logPrefix)
	log.WithContext(ctx).Info("start initialization")

	content, err := ioutil.ReadFile(db.config.SchemaPath)
	if err != nil {
		return errors.Errorf("error while reading schema file: %w", err)
	}

	db.metaDB.WaitForStarting()
	if db.metaDB.IsLeader() {
		listOfCommands := strings.Split(string(content), schemaDelimiter)
		for _, cmd := range listOfCommands {
			if _, err := db.metaDB.Write(ctx, cmd); err != nil {
				return errors.Errorf("error while creating db schema: %w", err)
			}
		}
	}

	db.templates, err = NewTemplateKeeper(db.config.TemplatePath)
	if err != nil {
		return errors.Errorf("error while creating new template keeper: %w", err)
	}

	log.WithContext(ctx).Info("done initialization")
	// block until context is done
	<-ctx.Done()
	log.WithContext(ctx).Info("userdata service is stopped")
	return nil
}

func NewDatabaseOps(metaDB metadb.MetaDB, config *Config) *DatabaseOps {
	return &DatabaseOps{
		metaDB: metaDB,
		config: config,
	}
}
