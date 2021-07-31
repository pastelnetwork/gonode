package database

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"text/template"

	"github.com/mitchellh/mapstructure"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
)

var (
	logPrefix        = "database"
	queryLevelStrong = "strong"
	schemaDelimiter  = "---"
)

type Config struct {
	SchemaPath         string `mapstructure:"schema-path" json:"schema-path,omitempty"`
	WriteTemplatePath  string `mapstructure:"write-template-path" json:"write-template-path,omitempty"`
	QueryTemplatePath  string `mapstructure:"query-template-path" json:"query-template-path,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

type DatabaseOps struct {
	metaDB         metadb.MetaDB
	writeTemplate  *template.Template
	queryTemplate  *template.Template
	updateTemplate *template.Template
	config         *Config
}

func substituteTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var templateBuffer bytes.Buffer
	if tmpl == nil || data == nil {
		return "", fmt.Errorf("input nil template or data")
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
	return db.metaDB.LeaderAddress()
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) WriteUserData(ctx context.Context, data pb.UserdataRequest) error {
	command, err := substituteTemplate(db.writeTemplate, pbToWriteCommand(data))
	if err != nil {
		return err
	}

	if _, err := db.metaDB.Write(ctx, command); err != nil {
		return err
	}

	return nil
}

// WriteUserData writes metadata in the struct UserdataProcessRequest to metadb
func (db *DatabaseOps) ReadUserData(ctx context.Context, artistPastelID string) (userdata.UserdataProcessRequest, error) {
	command, err := substituteTemplate(db.queryTemplate, artistPastelID)
	if err != nil {
		return userdata.UserdataProcessRequest{}, err
	}

	queryResult, err := db.metaDB.Query(ctx, command, queryLevelStrong)
	if err != nil {
		return userdata.UserdataProcessRequest{}, err
	}

	nrows := queryResult.NumRows()
	if nrows == 0 {
		return userdata.UserdataProcessRequest{}, fmt.Errorf("no artist with pastel id = %s", artistPastelID)
	} else if nrows > 1 {
		return userdata.UserdataProcessRequest{}, fmt.Errorf("upto %d records are returned", nrows)
	}

	//right here we make sure that there is just 1 row in the result
	queryResult.Next()
	resultMap, err := queryResult.Map()
	if err != nil {
		return userdata.UserdataProcessRequest{}, err
	}

	var dbResult UserdataReadResult
	if err := mapstructure.Decode(resultMap, &dbResult); err != nil {
		return userdata.UserdataProcessRequest{}, err
	}

	return dbResult.ToUserData(), nil
}

func (db *DatabaseOps) Run(ctx context.Context) error {

	ctx = log.ContextWithPrefix(ctx, logPrefix)
	log.WithContext(ctx).Info("start initialization")

	content, err := ioutil.ReadFile(db.config.SchemaPath)
	if err != nil {
		return err
	}

	db.metaDB.WaitForStarting()
	if db.metaDB.IsLeader() {
		listOfCommands := strings.Split(string(content), schemaDelimiter)
		for _, cmd := range listOfCommands {
			if _, err := db.metaDB.Write(ctx, cmd); err != nil {
				return err
			}
		}
	}

	db.writeTemplate, err = template.ParseFiles(db.config.WriteTemplatePath)
	if err != nil {
		return err
	}

	db.queryTemplate, err = template.ParseFiles(db.config.QueryTemplatePath)
	if err != nil {
		return err
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
