package grpc

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	json "github.com/json-iterator/go"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/rqstore"
	"github.com/pastelnetwork/gonode/common/utils"
	pb "github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/raptorq/node"
)

const (
	inputEncodeFileName = "input.data"
	symbolIDFileSubDir  = "meta"
	symbolFileSubdir    = "symbols"
	concurrency         = 4
)

type raptorQ struct {
	conn      *clientConn
	client    pb.RaptorQClient
	config    *node.Config
	semaphore chan struct{} // Semaphore to control concurrency
}

func randID() string {
	id := uuid.NewString()
	return id[0:8]
}

func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0750)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func createTaskFolder(base string, subDirs ...string) (string, error) {
	taskID := randID()
	taskPath := filepath.Join(base, taskID)
	taskPath = filepath.Join(taskPath, filepath.Join(subDirs...))

	err := os.MkdirAll(taskPath, 0750)

	if err != nil {
		return "", err
	}

	return taskPath, nil
}

func createInputEncodeFile(base string, data []byte) (taskPath string, inputFile string, err error) {
	taskPath, err = createTaskFolder(base)

	if err != nil {
		return "", "", errors.Errorf("create task folder: %w", err)
	}

	inputFile = filepath.Join(taskPath, inputEncodeFileName)
	err = writeFile(inputFile, data)

	if err != nil {
		return "", "", errors.Errorf("write file: %w", err)
	}

	return taskPath, inputFile, nil
}

func createInputDecodeSymbols(base string, symbols map[string][]byte) (path string, err error) {
	path, err = createTaskFolder(base, symbolFileSubdir)

	if err != nil {
		return "", errors.Errorf("create task folder: %w", err)
	}

	for id, data := range symbols {
		symbolFile := filepath.Join(path, id)
		err = writeFile(symbolFile, data)

		if err != nil {
			return "", errors.Errorf("write symbol file: %w", err)
		}
	}

	return path, nil
}

// scan symbol id files in "meta" folder, return map of file Ids & contents of file (as list of line)
func scanSymbolIDFiles(dirPath string) (map[string]node.RawSymbolIDFile, error) {
	filesMap := make(map[string]node.RawSymbolIDFile)

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Errorf("scan a path %s: %w", path, err)
		}

		if info.IsDir() {
			// TODO - compare it to root
			return nil
		}

		fileID := filepath.Base(path)

		configFile, err := os.Open(path)
		if err != nil {
			return errors.Errorf("opening file: %s - err: %w", path, err)
		}
		defer configFile.Close()

		file := node.RawSymbolIDFile{}
		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&file); err != nil {
			return errors.Errorf("parsing file: %s - err: %w", path, err)
		}

		filesMap[fileID] = file

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesMap, nil
}

// Encode data, and return a list of identifier of symbols
func (service *raptorQ) RQEncode(ctx context.Context, data []byte, id string, store rqstore.Store) (*node.Encode, error) {
	service.semaphore <- struct{}{} // Acquire slot
	defer func() {
		<-service.semaphore // Release the semaphore slot
	}()

	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	_, inputPath, err := createInputEncodeFile(service.config.RqFilesDir, data)
	if err != nil {
		return nil, errors.Errorf("create input file: %w", err)
	}

	req := pb.EncodeRequest{
		Path: inputPath,
	}

	res, err := service.client.Encode(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send encode request: %w", err)
	}

	fileMap, err := utils.ReadDirFilenames(res.Path)
	if err != nil {
		return nil, errors.Errorf("scan symbol folder %s: %w", res.Path, err)
	}

	if err := store.StoreSymbolDirectory(id, res.Path); err != nil {
		return nil, errors.Errorf("store symbol directory: %w", err)
	}

	if len(fileMap) != int(res.SymbolsCount) {
		return nil, errors.Errorf("symbol count not match : expect %d, result %d", res.SymbolsCount, len(fileMap))
	}
	output := &node.Encode{
		Symbols: fileMap,
		EncoderParam: node.EncoderParameters{
			Oti: res.EncoderParameters,
		},
	}

	if err := os.Remove(inputPath); err != nil {
		log.WithContext(ctx).WithError(err).Error("rq-encode: error removing input file")
	}

	return output, nil
}

func (service *raptorQ) EncodeInfo(ctx context.Context, data []byte, copies uint32, blockHash string, pastelID string) (*node.EncodeInfo, error) {
	service.semaphore <- struct{}{} // Acquire a semaphore slot
	defer func() {
		<-service.semaphore // Release the semaphore slot
	}()

	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	_, inputPath, err := createInputEncodeFile(service.config.RqFilesDir, data)
	if err != nil {
		return nil, errors.Errorf("create input file: %w", err)
	}

	req := pb.EncodeMetaDataRequest{
		Path:        inputPath,
		FilesNumber: copies,
		BlockHash:   blockHash,
		PastelID:    pastelID,
	}

	res, err := service.client.EncodeMetaData(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send encode info request: %w", err)
	}

	// scan return symbol Id files
	filesMap, err := scanSymbolIDFiles(res.Path)
	if err != nil {
		return nil, errors.Errorf("scan symbol id files folder %s: %w", res.Path, err)
	}

	if len(filesMap) != int(copies) {
		return nil, errors.Errorf("symbol id files count not match: expect %d, output %d", copies, len(filesMap))
	}

	output := &node.EncodeInfo{
		SymbolIDFiles: filesMap,
		EncoderParam: node.EncoderParameters{
			Oti: res.EncoderParameters,
		},
	}

	if err := os.Remove(inputPath); err != nil {
		log.WithContext(ctx).WithError(err).Error("encode info: error removing input file")
	}

	return output, nil
}

func (service *raptorQ) Decode(ctx context.Context, encodeInfo *node.Encode) (*node.Decode, error) {
	if encodeInfo == nil {
		return nil, errors.Errorf("invalid encode info")
	}

	ctx = service.contextWithLogPrefix(ctx)

	symbolsDir, err := createInputDecodeSymbols(service.config.RqFilesDir, encodeInfo.Symbols)
	if err != nil {
		return nil, errors.Errorf("create symbol files: %w", err)
	}

	req := pb.DecodeRequest{
		EncoderParameters: encodeInfo.EncoderParam.Oti,
		Path:              symbolsDir,
	}

	res, err := service.client.Decode(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("send decode request: %w", err)
	}

	restoredFile, err := readFile(res.Path)
	if err != nil {
		return nil, errors.Errorf("read file: %w, Path: %s", err, res.Path)
	}

	output := &node.Decode{
		File: restoredFile,
	}

	return output, nil
}

func (service *raptorQ) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRaptorQ(conn *clientConn, config *node.Config) node.RaptorQ {
	return &raptorQ{
		conn:      conn,
		client:    pb.NewRaptorQClient(conn),
		config:    config,
		semaphore: make(chan struct{}, concurrency),
	}
}
