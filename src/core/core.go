package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/voodooEntity/go-tachicrypt/src/encryptor"
	"github.com/voodooEntity/go-tachicrypt/src/fileutils"
	"github.com/voodooEntity/go-tachicrypt/src/masterlock"
	"github.com/voodooEntity/go-tachicrypt/src/splitter"
	"github.com/voodooEntity/go-tachicrypt/src/utils"
	"github.com/voodooEntity/go-tachicrypt/src/zipper"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Core struct {
	PartSize  int
	KeySize   int
	SaltSize  int
	PartCount int
}

func New() *Core {
	return &Core{
		SaltSize: 100,
		KeySize:  32,
	}
}

func (c *Core) Hide(dataPath string, partCount int, outputDir string) error {
	c.PartCount = partCount

	zipr := zipper.New()
	base64Zip, err := zipr.ZipToString(dataPath)
	if err != nil {
		return fmt.Errorf("error zipping and encoding: %w", err)
	}

	parts, padding := splitter.SplitStringWithPadding(base64Zip, c.PartCount)

	var partInfos []masterlock.PartInfo
	for i, part := range parts {

		encryptedPart, key, err := encryptor.EncryptWithRandomKey(part)
		if err != nil {
			return fmt.Errorf("error encrypting part: %w", err)
		}

		filename, err := utils.GenerateRandomFilename()
		if err != nil {
			return fmt.Errorf("error generating filename: %w", err)
		}

		partPath := filepath.Join(outputDir, filename)
		err = fileutils.WriteToFile(partPath, []byte(encryptedPart))
		if err != nil {
			return fmt.Errorf("error writing encrypted part to file: %w", err)
		}

		partInfos = append(partInfos, masterlock.PartInfo{
			Index:    i,
			Filename: filename,
			Key:      key,
		})
	}

	masterLockData, err := masterlock.CreateMasterLock(partInfos, padding)
	if err != nil {
		return fmt.Errorf("error creating master lock file: %w", err)
	}

	password, err := utils.PromptForPassword()
	if err != nil {
		return fmt.Errorf("error prompting for password: %w", err)
	}

	encryptedMasterLock, err := encryptor.EncryptWithPassword(string(masterLockData), password)
	if err != nil {
		return fmt.Errorf("error encrypting master lock file: %w", err)
	}
	fmt.Println(string(encryptedMasterLock))
	masterLockPath := filepath.Join(outputDir, "masterlock")
	err = fileutils.WriteToFile(masterLockPath, []byte(encryptedMasterLock))
	if err != nil {
		return fmt.Errorf("error writing master lock file: %w", err)
	}

	err = fileutils.ObfuscateFileTimestamps(outputDir)
	if err != nil {
		return fmt.Errorf("error obfuscating file timestamps: %w", err)
	}

	return nil
}

// Unhide decrypts the data from the given parts and master lock file, storing the result in the output path.
func (c *Core) Unhide(partsDir, masterlockPath, outputPath string) error {
	// Step 1: Decrypt Master Lock File
	encryptedMasterLock, err := ioutil.ReadFile(filepath.Join(partsDir, "masterlock.enc"))
	if err != nil {
		return fmt.Errorf("error reading encrypted master lock file: %w", err)
	}

	password, err := utils.PromptForPassword()
	if err != nil {
		return fmt.Errorf("error prompting for password: %w", err)
	}

	decryptedMasterLock, err := encryptor.DecryptWithPassword(string(encryptedMasterLock), password)
	if err != nil {
		return fmt.Errorf("error decrypting master lock file: %w", err)
	}

	var masterLock masterlock.MasterLock
	err = json.Unmarshal([]byte(decryptedMasterLock), &masterLock)
	if err != nil {
		return fmt.Errorf("error unmarshaling master lock file: %w", err)
	}

	// Step 2: Decrypt Each Part
	var base64Parts []string
	for _, partInfo := range masterLock.Parts {
		partPath := filepath.Join(partsDir, partInfo.Filename)
		encryptedPart, err := ioutil.ReadFile(partPath)
		if err != nil {
			return fmt.Errorf("error reading encrypted part file: %w", err)
		}

		key, err := base64.StdEncoding.DecodeString(partInfo.Key)
		if err != nil {
			return fmt.Errorf("error decoding key: %w", err)
		}

		decryptedPart, err := encryptor.DecryptWithRandomKey(string(encryptedPart), string(key))
		if err != nil {
			return fmt.Errorf("error decrypting part: %w", err)
		}

		base64Parts = append(base64Parts, string(decryptedPart))
	}

	// Step 3: Reconstruct Base64 Data
	paddedData := strings.Join(base64Parts, "")
	paddedData = strings.TrimRight(paddedData, "=")

	// Step 4: Decode and Unzip
	decodedData, err := base64.StdEncoding.DecodeString(paddedData)
	if err != nil {
		return fmt.Errorf("error decoding base64 data: %w", err)
	}

	zipper := zipper.New()
	err = zipper.UnzipFromBytes(decodedData, outputPath)
	if err != nil {
		return fmt.Errorf("error unzipping data: %w", err)
	}

	return nil
}
