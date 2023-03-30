package hermes

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/domain"
)

func (s *service) restartPastelID(ctx context.Context) error {
	if s.restartPasteldCheckTimer != pastelIDRestartInterval {
		s.restartPasteldCheckTimer += 2
	}

	s.restartPasteldCheckTimer = 0
	if isAvailable := s.checkNextBlockAvailable(ctx); !isAvailable {
		log.WithContext(ctx).Info("next block not available after 10 minutes, should restart pastelid")

		pastelD, err := getPastelDPath(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error retrieving pasteld path")
			return nil
		}

		pastelCli, err := getPastelCliPath(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error retrieving pastel-cli path")
			return nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error getting home dir")
			return nil
		}

		cmd := exec.Command(pastelCli, "stop")
		cmd.Dir = homeDir

		log.WithContext(ctx).Info("pastel-cli has been stopped")
		time.Sleep(15 * time.Second)

		var extIP string
		if extIP, err = utils.GetExternalIPAddress(); err != nil {
			log.WithContext(ctx).WithError(err).Error("Could not get external IP address")
			return nil
		}

		dataDir, err := getDataDir(ctx, filepath.Join(homeDir, ".pastel/supernode.yml"))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error retrieving data-dir from supernode.yml")
			return nil
		}

		mnPrivKey, err := getMasternodePrivKey(ctx, filepath.Join(dataDir, "testnet3/masternode.conf"), extIP)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error retrieving masternode private key")
			return nil
		}

		var pasteldArgs []string
		pasteldArgs = append(pasteldArgs,
			fmt.Sprintf("--datadir=%s", dataDir),
			fmt.Sprintf("--externalip=%s", extIP),
			fmt.Sprintf("--masternodeprivkey=%s", mnPrivKey),
			"--txindex=1",
			"--reindex",
			"--masternode",
			"--daemon",
		)

		log.WithContext(ctx).Infof("Starting -> %s %s", pastelD, strings.Join(pasteldArgs, " "))
		cmd = exec.Command(pastelD, pasteldArgs...)
		cmd.Dir = homeDir

		if err := cmd.Run(); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error starting pastelid")
			return nil
		}

		log.WithContext(ctx).Info("pasteld has been restarted")
	} else {
		log.WithContext(ctx).Info("block count has been updated")
	}

	return nil
}

func getDataDir(ctx context.Context, confFilePath string) (dataDir string, err error) {
	log.WithContext(ctx).Info("opening the supernode.yml")
	file, err := os.OpenFile(confFilePath, os.O_RDWR, 0644)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Could not open supernode.yml - %s", confFilePath)
		return "", err
	}
	log.WithContext(ctx).Info("File opened")
	defer file.Close()

	log.WithContext(ctx).Info("Parsing data-dir from supernode.yml")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "work-dir:") {
			dataDir = line
			if err != nil {
				log.WithContext(ctx).Error(fmt.Sprintf("Could not parse work-dir: from file: - %s", confFilePath))
				return "", err
			}

			dataDirPath := strings.Split(dataDir, ": ")
			return dataDirPath[1], nil
		}
	}

	return dataDir, nil
}

func loadMasternodeConfFile(ctx context.Context, masternodeConfFilePath string) (map[string]domain.MasterNodeConf, error) {
	// Read ConfData from masternode.conf
	confFile, err := os.ReadFile(masternodeConfFilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Failed to read existing masternode.conf file - %s", masternodeConfFilePath))
		return nil, err
	}

	var conf map[string]domain.MasterNodeConf
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Invalid existing masternode.conf file - %s", masternodeConfFilePath))
		return nil, err
	}

	return conf, nil
}

func getMasternodePrivKey(ctx context.Context, masterNodeConffilePath, extIP string) (string, error) {
	conf, err := loadMasternodeConfFile(ctx, masterNodeConffilePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to load existing masternode.conf file")
		return "", err
	}

	log.WithContext(ctx).Info("Attempting to load existing masternode.conf file using external IP address...")
	for mnName, mnConf := range conf {
		extAddrPort := strings.Split(mnConf.MnAddress, ":")
		extAddr := extAddrPort[0] // get Ext IP and Port
		if extAddr == extIP {
			log.WithContext(ctx).Info(fmt.Sprintf("Loading masternode.conf file using %s conf", mnName))
			return mnConf.MnPrivKey, nil
		}
	}

	return "", errors.New("not able to found private key")
}

// checkNextBlockAvailable calls pasteld and checks if a new block is available
func (s *service) checkNextBlockAvailable(ctx context.Context) bool {
	blockCount, err := s.pastelClient.GetBlockCount(ctx)
	if err != nil {
		return false
	}
	if blockCount > s.currentBlockCount {
		atomic.StoreInt32(&s.currentBlockCount, blockCount)
		return true
	}

	return false
}

func getPastelDPath(ctx context.Context) (path string, err error) {
	//create command
	findCmd := exec.Command("find", ".", "-print")
	grepCmd := exec.Command("grep", "-x", "./pastel/pasteld")

	findCmd.Dir, err = os.UserHomeDir()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting home dir")
		return "", err
	}

	//make a pipe and set the input and output to reader and writer
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	findCmd.Stdout = writer
	grepCmd.Stdin = reader

	//cache the output of "grep" to memory
	grepCmd.Stdout = &buf

	//starting the commands
	findCmd.Start()
	grepCmd.Start()

	//waiting for commands to complete and close the reader & writer
	findCmd.Wait()
	writer.Close()

	grepCmd.Wait()
	reader.Close()

	pathWithEscapeCharacter := buf.String()
	return strings.Replace(pathWithEscapeCharacter, "\n", "", 1), nil
}

func getPastelCliPath(ctx context.Context) (path string, err error) {
	//create command
	findCmd := exec.Command("find", ".", "-print")
	grepCmd := exec.Command("grep", "-x", "./pastel/pastel-cli")

	findCmd.Dir, err = os.UserHomeDir()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error getting home dir")
		return "", err
	}

	//make a pipe and set the input and output to reader and writer
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	findCmd.Stdout = writer
	grepCmd.Stdin = reader

	//cache the output of "grep" to memory
	grepCmd.Stdout = &buf

	//starting the commands
	findCmd.Start()
	grepCmd.Start()

	//waiting for commands to complete and close the reader & writer
	findCmd.Wait()
	writer.Close()

	grepCmd.Wait()
	reader.Close()

	pathWithEscapeCharacter := buf.String()
	return strings.Replace(pathWithEscapeCharacter, "\n", "", 1), nil
}
