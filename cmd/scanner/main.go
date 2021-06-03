package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"OpenZeppelin/fotify-node/clients"
	"OpenZeppelin/fotify-node/config"
	"OpenZeppelin/fotify-node/feeds"
	"OpenZeppelin/fotify-node/services"
	"OpenZeppelin/fotify-node/services/scanner"
)

func loadKey() (*keystore.Key, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return nil, err
	}

	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	passphrase := string(pw)

	files, err := ioutil.ReadDir("/.keys")
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, errors.New("there must be only one key in key directory")
	}

	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", "/.keys", files[0].Name()))
	if err != nil {
		return nil, err
	}

	return keystore.DecryptKey(keyBytes, passphrase)
}

func initTxStream(ctx context.Context, cfg config.Config) (*scanner.TxStreamService, error) {
	url := cfg.Scanner.Ethereum.JsonRpcUrl
	startBlock := config.ParseBigInt(cfg.Scanner.StartBlock)
	endBlock := config.ParseBigInt(cfg.Scanner.EndBlock)
	chainID := config.ParseBigInt(cfg.Scanner.ChainID)

	if url == "" {
		return nil, fmt.Errorf("ethereum.jsonRpcUrl is required")
	}

	tracing := true
	if cfg.Scanner.DisableTracing {
		tracing = false
	}

	return scanner.NewTxStreamService(ctx, scanner.TxStreamServiceConfig{
		Url: url,
		BlockFeedConfig: feeds.BlockFeedConfig{
			Start:   startBlock,
			End:     endBlock,
			ChainID: chainID,
			Tracing: tracing,
		},
	})
}

func initTxAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService) (*scanner.TxAnalyzerService, error) {
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return scanner.NewTxAnalyzerService(ctx, scanner.TxAnalyzerServiceConfig{
		TxChannel:    stream.ReadOnlyTxStream(),
		AgentConfigs: cfg.Agents,
		AlertSender:  as,
	})
}

func initBlockAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService) (*scanner.BlockAnalyzerService, error) {
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return scanner.NewBlockAnalyzerService(ctx, scanner.BlockAnalyzerServiceConfig{
		BlockChannel: stream.ReadOnlyBlockStream(),
		AgentConfigs: cfg.Agents,
		AlertSender:  as,
	})
}

func initAlertSender(ctx context.Context) (clients.AlertSender, error) {
	key, err := loadKey()
	if err != nil {
		return nil, err
	}
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return clients.NewAlertSender(ctx, clients.AlertSenderConfig{
		Key:           key,
		QueryNodeAddr: qn,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	as, err := initAlertSender(ctx)
	if err != nil {
		return nil, err
	}
	txStream, err := initTxStream(ctx, cfg)
	if err != nil {
		return nil, err
	}
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream)
	if err != nil {
		return nil, err
	}

	return []services.Service{
		txStream,
		txAnalyzer,
		blockAnalyzer,
		scanner.NewTxLogger(ctx),
	}, nil
}

func main() {
	services.ContainerMain("scanner", initServices)
}
