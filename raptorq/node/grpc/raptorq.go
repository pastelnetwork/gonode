package grpc

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/raptorq/node"
)

const (
	inputEncodeFileName = "input.data"
	symbolIdFileSubDir  = "meta"
	symbolFileSubdir    = "symbols"
)

type raptorQ struct {
	conn   *clientConn
	client pb.RaptorQClient
	config *node.Config
}

func randId() string {
	id := uuid.NewString()
	return id[0:8]
}

func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0777)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func readFileLines(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, errors.Errorf("failed to read file: %w", err)
	}
	return strings.Split(string(b), "/n"), nil
}

func createTaskFolder(base string) (string, error) {
	taskId := randId()
	taskPath := filepath.Join(base, taskId)
	err := os.MkdirAll(taskPath, 0777)

	if err != nil {
		return "", err
	}

	return taskPath, nil
}

func createInputEncodeFile(base string, data []byte) (taskPath string, inputFile string, err error) {
	taskPath, err = createTaskFolder(base)

	if err != nil {
		return "", "", errors.Errorf("failed to create task folder: %w", err)
	}

	inputFile = filepath.Join(taskPath, inputEncodeFileName)
	err = writeFile(inputFile, data)

	if err != nil {
		return "", "", errors.Errorf("failed to create task folder: %w", err)
	}

	return taskPath, inputFile, nil
}

// scan symbol id files in "meta" folder, return map of file Ids & contents of file (as list of line)
func scanSymbolIdFiles(dirPath string) (map[string][]string, error) {
	filesMap := make(map[string][]string)

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Errorf("failed to scan a path %s: %w", path, err)
		}

		if info.IsDir() {
			// TODO - compare it to root
			return nil
		}

		fileId := filepath.Base(path)

		lines, err := readFileLines(path)
		if err != nil {
			return errors.Errorf("failed to read file %s: %w", path, err)
		}

		filesMap[fileId] = lines

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesMap, nil
}

// scan symbol  files in "symbols" folder, return map of file Ids & contents of file (as list of line)
func scanSymbolFiles(dirPath string) (map[string][]byte, error) {
	filesMap := make(map[string][]byte)

	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Errorf("failed to scan a path %s: %w", path, err)
		}

		if info.IsDir() {
			// TODO - compare it to root
			return nil
		}

		fileId := filepath.Base(path)

		data, err := readFile(path)
		if err != nil {
			return errors.Errorf("failed to read file %s: %w", path, err)
		}

		filesMap[fileId] = data

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesMap, nil
}

// Encode data, and return a list of identifier of symbols
func (service *raptorQ) Encode(ctx context.Context, data []byte) (*node.Encode, error) {
	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	_, inputPath, err := createInputEncodeFile(service.config.RqFilesDir, data)
	if err != nil {
		return nil, errors.Errorf("failed to create input file: %w", err)
	}

	req := pb.EncodeRequest{
		Path: inputPath,
	}

	res, err := service.client.Encode(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("failed to send request: %w", err)
	}

	fileMap, err := scanSymbolFiles(res.Path)
	if err != nil {
		return nil, errors.Errorf("failed to scan symbol folder %s: %w", res.Path, err)
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

	return output, nil
}

func (service *raptorQ) EncodeInfo(ctx context.Context, data []byte, copies uint32, blockHash string, pastelID string) (*node.EncodeInfo, error) {
	if data == nil {
		return nil, errors.Errorf("invalid data")
	}

	ctx = service.contextWithLogPrefix(ctx)

	_, inputPath, err := createInputEncodeFile(service.config.RqFilesDir, data)
	if err != nil {
		return nil, errors.Errorf("failed to create input file: %w", err)
	}

	req := pb.EncodeMetaDataRequest{
		Path:        inputPath,
		FilesNumber: copies,
		BlockHash:   blockHash,
		PastelId:    pastelID,
	}

	res, err := service.client.EncodeMetaData(ctx, &req)
	if err != nil {
		return nil, errors.Errorf("failed to send request: %w", err)
	}

	// scan return symbol Id files
	symbolCnt := res.SymbolsCount
	filesMap, err := scanSymbolIdFiles(res.Path)
	if err != nil {
		return nil, errors.Errorf("failed to scan symbol id files folder %s: %w", res.Path, err)
	}

	if len(filesMap) != int(copies) {
		return nil, errors.Errorf("symbol id files count not match: expect %d, output %d", copies, len(filesMap))
	}

	symbolIdFiles := make(map[string]node.RawSymbolIdFile)
	for fileID, lines := range filesMap {
		if len(lines) != int(symbolCnt+3) {
			return nil, errors.Errorf("file length not match: file %s, expect %d, output %d", fileID, symbolCnt+3, len(lines))
		}

		symbolIdFiles[fileID] = node.RawSymbolIdFile{
			Id:                lines[0],
			BlockHash:         lines[1],
			PastelId:          lines[2],
			SymbolIdentifiers: lines[3:],
		}

		// TODO : validate Id, blockHash, pastelId
	}

	output := &node.EncodeInfo{
		SymbolIdFiles: symbolIdFiles,
		EncoderParam: node.EncoderParameters{
			Oti: res.EncoderParameters,
		},
	}

	return output, nil
}

func (service *raptorQ) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newRaptorQ(conn *clientConn, config *node.Config) node.RaptorQ {
	return &raptorQ{
		conn:   conn,
		client: pb.NewRaptorQClient(conn),
		config: config,
	}
}
