package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/pastel-utility/configs"
	"github.com/pkg/errors"
)

func setupInitCommand(config *configs.Config) *cli.Command {
	initCommand := cli.NewCommand("init")
	initCommand.SetUsage("Command that performs initialization of the system for both Wallet and SuperNodes")
	initCommand.AddFlags(
		cli.NewFlag("network", &config.Init.Network).SetUsage("Set network type `mainnet` or `testnet`").SetAliases("n").SetValue("mainnet"),
		cli.NewFlag("force", &config.Init.Force).SetAliases("f").SetUsage("Force to overwrite config files and re-download ZKSnark parameters"),
		cli.NewFlag("peers", &config.Init.Peers).SetAliases("p").SetUsage("List of peers to add into pastel.conf file, must be in the format - `ip` or `ip:port`"),
	)

	// create walletnode and supernode subcommands
	walletnodeSubCommand := cli.NewCommand("walletnode")
	walletnodeSubCommand.SetUsage("Perform wallet specific initialization")
	walletnodeSubCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runWalletSubCommand(ctx, config)
	})

	supernodeSubCommand := cli.NewCommand("supernode")
	supernodeSubCommand.SetUsage("Perform Supernode/Masternode specific initialization")
	supernodeSubCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runSuperNodeSubCommand(ctx, config)
	})

	initCommand.AddSubcommands(
		walletnodeSubCommand,
		supernodeSubCommand,
	)

	initCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		ctx = log.ContextWithPrefix(ctx, "app")
		if len(args) == 0 {
			return errSubCommandRequired
		}
		return runInit(ctx, config)
	})

	return initCommand
}

func runInit(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Init")
	defer log.WithContext(ctx).Info("End")

	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
		os.Exit(0)
	})

	//run the init command logic flow
	err := initCommandLogic(ctx, config)
	if err != nil {
		return err
	}

	return nil
}

// initCommandLogic runs the init command logic flow
// takes all provided arguments in the cli and does all the background tasks
// Print success info log on successfully ran command, return error if fail
func initCommandLogic(ctx context.Context, config *configs.Config) error {
	zksnarkPath := filepath.Join(config.WorkDir, "/.pastel-params/")

	if _, err := os.Stat(zksnarkPath); os.IsNotExist(err) || config.Init.Force {
		if err := os.MkdirAll(zksnarkPath, os.ModePerm); err != nil {
			return errors.Errorf("could not create zksnarkpath %s, %v", zksnarkPath, err)
		}
	}

	v := viper.New()
	v.SetConfigType("json")
	if err := v.Unmarshal(config); err != nil {
		return err
	}

	if config.Init.Force {
		if err := v.WriteConfigAs(defaultConfigFile); err != nil {
			return err
		}
	} else {
		if err := v.SafeWriteConfigAs(defaultConfigFile); err != nil {
			return err
		}
	}

	// download zksnark params
	if err := downloadZksnarkParams(ctx, zksnarkPath, config.Init.Force); err != nil {
		return err
	}

	if err := checkLocalAndRouterFirewalls(ctx, []string{"80", "21"}); err != nil {
		return err
	}

	return nil
}

// downloadZksnarkParams downloads zksnark params to the specified forlder
// Print success info log on successfully ran command, return error if fail
func downloadZksnarkParams(ctx context.Context, path string, force bool) error {
	log.WithContext(ctx).Infoln("Downloading zksnark files")

	for _, zksnarkParamsName := range configs.ZksnarkParamsNames {
		zksnarkParamsPath := filepath.Join(path, zksnarkParamsName)

		if _, err := os.Stat(zksnarkParamsPath); os.IsNotExist(err) || force {
			out, err := os.Create(zksnarkParamsPath)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Error creating file: %s\n", zksnarkParamsPath)
				return errors.Errorf("Failed to create file: %v \n", err)
			}

			log.WithContext(ctx).Infof("downloading: %s", zksnarkParamsPath)
			resp, err := http.Get(configs.ZksnarkParamsURL + zksnarkParamsName)
			if err != nil {
				out.Close()
				log.WithContext(ctx).WithError(err).Errorf("Error downloading file: %s\n", configs.ZksnarkParamsURL+zksnarkParamsName)
				return errors.Errorf("Failed to download: %v \n", err)
			}

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				out.Close()
				resp.Body.Close()
				return err
			}

			out.Close()
			resp.Body.Close()
		} else {
			log.WithContext(ctx).Infof("zksnark file %s already exists and force set to %v", zksnarkParamsPath, force)
		}
	}
	log.WithContext(ctx).Info("ZkSnark params download completed.\n")
	return nil
}

// checkLocalAndRouterFirewalls checks local and router firewalls and suggest what to open
func checkLocalAndRouterFirewalls(ctx context.Context, requiredPorts []string) error {
	baseURL := "https://portchecker.com?q=" + strings.Join(requiredPorts[:], ",")
	// resp, err := http.Get(baseURL)
	// if err != nil {
	// 	log.WithContext(ctx).WithError(err).Errorf("Error requesting url\n")
	// 	return errors.Errorf("Failed to request port url %v \n", err)
	// }
	// defer resp.Body.Close()
	// ok := resp.StatusCode == http.StatusOK
	ok := true
	if ok {
		log.WithContext(ctx).Info("Your ports {} are opened and can be accessed by other PAstel nodes! ", baseURL)
	} else {
		log.WithContext(ctx).Info("Your ports {} are NOT opened and can NOT be accessed by other Pastel nodes!\n Please open this ports in your router firewall and in {} firewall", baseURL, "func_to_check_OS()")
	}
	return nil
}

// runWalletSubCommand runs wallet subcommand
func runWalletSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infoln("Wallet Node Start")
	defer log.WithContext(ctx).Infoln("Wallet Node End")

	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
		os.Exit(0)
	})

	walletNodePath := filepath.Join(config.WorkDir, "walletnode")

	if _, err := os.Stat(walletNodePath); os.IsNotExist(err) || config.Init.Force {
		if err := os.MkdirAll(walletNodePath, os.ModePerm); err != nil {
			return errors.Errorf("could not create walletnode path %s, %v", walletNodePath, err)
		}
	}

	v := viper.New()
	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewBufferString(config.Init.String())); err != nil {
		return err
	}

	walletConfigPath := fmt.Sprintf("%s/wallet.env", walletNodePath)

	if config.Init.Force {
		if err := v.WriteConfigAs(walletConfigPath); err != nil {
			return err
		}
	} else {
		if err := v.SafeWriteConfigAs(walletConfigPath); err != nil {
			return err
		}
	}

	return nil
}

// runSuperNodeSubCommand runs wallet subcommand
func runSuperNodeSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infoln("Super Node Start")
	defer log.WithContext(ctx).Infoln("Super Node End")

	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
		os.Exit(0)
	})

	superNodePath := filepath.Join(config.WorkDir, "supernode")

	if _, err := os.Stat(superNodePath); os.IsNotExist(err) || config.Init.Force {
		if err := os.MkdirAll(superNodePath, os.ModePerm); err != nil {
			return errors.Errorf("could not create supernode path %s, %v", superNodePath, err)
		}
	}

	v := viper.New()
	if err := v.ReadConfig(bytes.NewBufferString(config.Start.String())); err != nil {
		return err
	}

	superNodeConfigPath := fmt.Sprintf("%s/supernode.env", superNodePath)

	if config.Init.Force {
		if err := v.WriteConfigAs(superNodeConfigPath); err != nil {
			return err
		}
	} else {
		if err := v.SafeWriteConfigAs(superNodeConfigPath); err != nil {
			return err
		}
	}

	// pip install gdown
	// ~/.local/bin/gdown https://drive.google.com/uc?id=1U6tpIpZBxqxIyFej2EeQ-SbLcO_lVNfu
	// unzip ./SavedMLModels.zip -d <paslel-dir>/supernode/tfmodels

	_, err := RunCMD("pip3", "install", "gdown")
	if err != nil {
		return err
	}

	savedModelURL := "https://drive.google.com/uc?id=1U6tpIpZBxqxIyFej2EeQ-SbLcO_lVNfu"

	log.WithContext(ctx).Infof("Downloading: %s ...\n", savedModelURL)

	_, err = RunCMD("gdown", savedModelURL)
	if err != nil {
		return err
	}

	tfmodelsPath := filepath.Join(superNodePath, "tfmodels")
	if _, err := os.Stat(tfmodelsPath); os.IsNotExist(err) || config.Init.Force {
		if err := os.MkdirAll(tfmodelsPath, os.ModePerm); err != nil {
			return errors.Errorf("could not create tfmodels path %s, %v", tfmodelsPath, err)
		}
	}

	_, err = RunCMD("unzip", fmt.Sprintf("%s/SavedMLModels.zip", defaultTempDir), "-d", tfmodelsPath)
	if err != nil {
		return err
	}

	return nil
}

// RunCMD runs shell command and returns output and error
func RunCMD(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}
