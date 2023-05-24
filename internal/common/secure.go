package common

import (
	"bufio"
	"io"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
)

func ReadPrivateKey() (crypto.PrivKey, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var (
		privateKeyPath = homeDir + "/.boatswain/keys/private_key"
		publicKeyPath  = homeDir + "/.boatswain/keys/public_key"
	)

	if _, err := os.Stat(privateKeyPath); err == nil {
		file, err := os.Open(privateKeyPath)
		if err != nil {
			return nil, err
		}

		keyContent, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		privateKey, err := crypto.UnmarshalRsaPrivateKey(keyContent)

		return privateKey, err
	}

	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return nil, err
	}
	privKeyBytes, err := privKey.Raw()
	if err != nil {
		return nil, err
	}
	pubKeyBytes, err := pubKey.Raw()
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(privateKeyPath, privKeyBytes, 0777); err != nil {
		return nil, err
	}
	if err := os.WriteFile(publicKeyPath, pubKeyBytes, 0777); err != nil {
		return nil, err
	}

	return privKey, nil
}

func LoadPeersList() ([]multiaddr.Multiaddr, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var listPath = homeDir + "/.boatswain/peers.txt"

	if _, err := os.Stat(listPath); err == nil {
		//
	} else {
		return nil, nil
	}

	file, err := os.Open(listPath)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	var bootstrapPeers []multiaddr.Multiaddr

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		line = []byte(strings.TrimSpace(string(line)))
		if len(line) == 0 {
			continue
		}

		addr, parseAddrErr := multiaddr.NewMultiaddr(string(line))
		if parseAddrErr != nil {
			zap.L().Warn("invalid multi-address",
				zap.String("line", string(line)),
				zap.Error(parseAddrErr))
			continue
		}

		bootstrapPeers = append(bootstrapPeers, addr)
	}

	return bootstrapPeers, nil
}
