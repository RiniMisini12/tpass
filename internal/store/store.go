package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rinimisini112/tpass/internal/config"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

type PasswordEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Store struct {
	Entries []PasswordEntry `json:"entries"`
}

const (
	saltLength  = 16
	nonceLength = 12
	timeCost    = 1
	memoryCost  = 64 * 1024
	parallelism = 4
	keyLength   = 32
)

func ConstructMasterPassword(cfg config.SecurityConfig) (string, error) {
	osSec, err := config.GetOSSecurity()
	if err != nil {
		return "", err
	}

	cfg.UID = osSec.UID
	cfg.GID = osSec.GID
	cfg.Username = osSec.Username

	saltSource := cfg.UID + cfg.GID + cfg.Username
	saltHash := sha256.Sum256([]byte(saltSource))
	salt := saltHash[:16]

	key := argon2.IDKey([]byte(cfg.MainPassword), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%x", key), nil
}

func GetMasterPassword() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	masterPassword, err := ConstructMasterPassword(cfg.Security)
	if err != nil {
		log.Fatalf("Error constructing master password: %v", err)
	}

	return masterPassword
}

func deriveKey(masterPassword string, salt []byte) []byte {
	return argon2.IDKey([]byte(masterPassword), salt, timeCost, memoryCost, uint8(parallelism), keyLength)
}

func encryptData(plaintext []byte, masterPassword string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := deriveKey(masterPassword, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, nonceLength)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	full := append(salt, nonce...)
	full = append(full, ciphertext...)

	return base64.StdEncoding.EncodeToString(full), nil
}

func decryptData(encoded string, masterPassword string) ([]byte, error) {
	full, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(full) < saltLength+nonceLength {
		return nil, errors.New("encrypted data too short")
	}

	salt := full[:saltLength]
	nonce := full[saltLength : saltLength+nonceLength]
	ciphertext := full[saltLength+nonceLength:]

	key := deriveKey(masterPassword, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func storeFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	storeDir := filepath.Join(home, ".tpass")
	if err := os.MkdirAll(storeDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(storeDir, "store.json"), nil
}

func SaveStore(store Store, masterPassword string) error {
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}

	encrypted, err := encryptData(data, masterPassword)
	if err != nil {
		return err
	}

	filePath, err := storeFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(encrypted), 0600)
}

func LoadStore(masterPassword string) (Store, error) {
	filePath, err := storeFilePath()
	if err != nil {
		return Store{}, err
	}

	encrypted, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Store{Entries: []PasswordEntry{}}, nil
		}
		return Store{}, err
	}

	data, err := decryptData(string(encrypted), masterPassword)
	if err != nil {
		return Store{}, err
	}

	var store Store
	if err := json.Unmarshal(data, &store); err != nil {
		return Store{}, err
	}
	return store, nil
}
