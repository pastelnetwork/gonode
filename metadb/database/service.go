package database

import (
	"bytes"
	"context"
	"io/ioutil"
	"regexp"
	"strings"

	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
)

var (
	logPrefix       = "database"
	queryLevelNone  = "none"
	schemaDelimiter = "---"
)

type Config struct {
	SchemaPath        string `mapstructure:"schema-path" json:"schema-path,omitempty"`
	WriteTemplatePath string `mapstructure:"write-template-path" json:"write-template-path,omitempty"`
	QueryTemplatePath string `mapstructure:"query-template-path" json:"query-template-path,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

type DatabaseOps struct {
	metaDB        metadb.MetaDB
	writeTemplate *template.Template
	queryTemplate *template.Template
	config        *Config
}

func substituteTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var templateBuffer bytes.Buffer
	if tmpl == nil || data == nil {
		return "", errors.Errorf("input nil template or data")
	}
	if err := tmpl.Execute(&templateBuffer, data); err != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

func (db *DatabaseOps) IsLeader() bool {
	return db.metaDB.IsLeader()
}

func (db *DatabaseOps) LeaderAddress() string {
	re := regexp.MustCompile(":[0-9]+$")
	address := re.Split(db.metaDB.LeaderAddress(), -1)[0]
	return address
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) WriteUserData(ctx context.Context, data *pb.UserdataRequest) error {
	if data == nil {
		return errors.Errorf("input nil data")
	}
	command, err := substituteTemplate(db.writeTemplate, pbToWriteCommand(data))
	if err != nil {
		return errors.Errorf("error while subtitute template: %w", err)
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

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) ReadUserData(ctx context.Context, artistPastelID string) (userdata.ProcessRequest, error) {
	command, err := substituteTemplate(db.queryTemplate, artistPastelID)
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

	db.writeTemplate, err = template.ParseFiles(db.config.WriteTemplatePath)
	if err != nil {
		return errors.Errorf("error while parsing write template: %w", err)
	}

	db.queryTemplate, err = template.ParseFiles(db.config.QueryTemplatePath)
	if err != nil {
		return errors.Errorf("error while parsing query template: %w", err)
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
