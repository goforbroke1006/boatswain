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

func GetKeysPair() (crypto.PrivKey, crypto.PubKey, error) {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		return nil, nil, homeDirErr
	}

	var (
		rootDir        = homeDir + "/.boatswain/"
		keysDir        = homeDir + "/.boatswain/keys"
		privateKeyPath = homeDir + "/.boatswain/keys/private_key"
		publicKeyPath  = homeDir + "/.boatswain/keys/public_key"
	)

	if mkdirErr := os.MkdirAll(rootDir, os.ModePerm); mkdirErr != nil {
		return nil, nil, mkdirErr
	}
	if mkdirErr := os.MkdirAll(keysDir, os.ModePerm); mkdirErr != nil {
		return nil, nil, mkdirErr
	}

	var (
		privateKey crypto.PrivKey = nil
		publicKey  crypto.PubKey  = nil
	)

	if _, err := os.Stat(privateKeyPath); err == nil {
		file, err := os.Open(privateKeyPath)
		if err != nil {
			return nil, nil, err
		}

		keyContent, err := io.ReadAll(file)
		if err != nil {
			return nil, nil, err
		}

		privKey, privKeyErr := crypto.UnmarshalRsaPrivateKey(keyContent)
		if privKeyErr != nil {
			return nil, nil, privKeyErr
		}
		privateKey = privKey
	}

	if _, err := os.Stat(publicKeyPath); err == nil {
		file, err := os.Open(publicKeyPath)
		if err != nil {
			return nil, nil, err
		}

		keyContent, err := io.ReadAll(file)
		if err != nil {
			return nil, nil, err
		}

		pubKey, pubKeyErr := crypto.UnmarshalRsaPublicKey(keyContent)
		if pubKeyErr != nil {
			return nil, nil, pubKeyErr
		}
		publicKey = pubKey
	}

	if privateKey == nil || publicKey == nil {
		var generateErr error
		privateKey, publicKey, generateErr = crypto.GenerateKeyPair(crypto.RSA, 2048)
		if generateErr != nil {
			return nil, nil, generateErr
		}
		privKeyBytes, err := privateKey.Raw()
		if err != nil {
			return nil, nil, err
		}
		pubKeyBytes, err := publicKey.Raw()
		if err != nil {
			return nil, nil, err
		}

		_, _ = os.OpenFile(privateKeyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		_, _ = os.OpenFile(publicKeyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

		if err := os.WriteFile(privateKeyPath, privKeyBytes, os.ModePerm); err != nil {
			return nil, nil, err
		}
		if err := os.WriteFile(publicKeyPath, pubKeyBytes, os.ModePerm); err != nil {
			return nil, nil, err
		}
	}

	return privateKey, publicKey, nil
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
