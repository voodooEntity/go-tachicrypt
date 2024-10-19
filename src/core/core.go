package core

import (
	"encoding/json"
	"fmt"
	"github.com/voodooEntity/go-tachicrypt/src/encryptor"
	"github.com/voodooEntity/go-tachicrypt/src/fileutils"
	"github.com/voodooEntity/go-tachicrypt/src/masterlock"
	"github.com/voodooEntity/go-tachicrypt/src/prettywriter"
	"github.com/voodooEntity/go-tachicrypt/src/splitter"
	"github.com/voodooEntity/go-tachicrypt/src/utils"
	"github.com/voodooEntity/go-tachicrypt/src/zipper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

func (c *Core) Hide(dataPath string, partCount int, outputDir string, prefilledPassword string) error {
	prettywriter.WriteInBox(40, "Configuration", prettywriter.Green, prettywriter.BlackBG, prettywriter.DoubleLine)
	prettywriter.Writeln("[==] Chosen mode: hide (encrypting)", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[==] Input directory: "+dataPath, prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[==] Output directory: "+outputDir, prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[==] Amount of parts: "+strconv.Itoa(partCount), prettywriter.Green, prettywriter.BlackBG)
	fmt.Println("")

	prettywriter.WriteInBox(40, "Starting Encryption Process", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	prettywriter.Writeln("[>>] Zipping input data to base64 string", prettywriter.Green, prettywriter.BlackBG)
	// Step 1: Create the zip base64 string data path
	c.PartCount = partCount

	zipr := zipper.New()
	base64Zip, err := zipr.ZipToString(dataPath)
	if err != nil {
		return fmt.Errorf("error zipping and encoding: %w", err)
	}
	prettywriter.Writeln("[>>] Splitting base64 zip into padded parts", prettywriter.Green, prettywriter.BlackBG)

	// Step 2: Split the base64 zip string into parts with padding
	parts, padding := splitter.SplitStringWithPadding(base64Zip, c.PartCount)

	// Step 3: Run encryption on all the parts, store them encrypted and add the info to masterlock
	var partInfos []masterlock.PartInfo
	for i, part := range parts {
		fmt.Print("\r")
		prettywriter.Write("[>>] Encrypt and store parts : "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(parts)), prettywriter.Green, prettywriter.BlackBG)
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
	fmt.Println("")
	prettywriter.Writeln("[**] All parts successfully encrypted and stored.", prettywriter.BlackBG, prettywriter.Green)
	fmt.Println("")
	// Step 4: Create Masterlock, prompt user for pwd and encrypt and store the masterlock
	masterLockData, err := masterlock.CreateMasterLock(partInfos, padding)
	if err != nil {
		return fmt.Errorf("error creating master lock file: %w", err)
	}

	password := ""
	if "" == prefilledPassword {
		password = utils.PromptForPassword("Please enter a password to encrypt the masterlock: ")
	} else {
		password = prefilledPassword
	}
	fmt.Println("")
	prettywriter.WriteInBox(40, "Handle masterlock file", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	prettywriter.Writeln("[>>] Encrypting masterlock", prettywriter.BlackBG, prettywriter.Green)
	encryptedMasterLock, err := encryptor.EncryptWithPassword(string(masterLockData), password)
	if err != nil {
		return fmt.Errorf("error encrypting master lock file: %w", err)
	}
	masterLockPath := filepath.Join(outputDir, "masterlock")
	prettywriter.Writeln("[>>] Writing masterlock", prettywriter.BlackBG, prettywriter.Green)
	err = fileutils.WriteToFile(masterLockPath, []byte(encryptedMasterLock))
	if err != nil {
		return fmt.Errorf("error writing master lock file: %w", err)
	}
	prettywriter.Writeln("[**] Masterlock successful written", prettywriter.BlackBG, prettywriter.Green)
	fmt.Println("")
	prettywriter.WriteInBox(40, "Final Shenanigans", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	prettywriter.Writeln("[>>] Obfuscating timestamps ", prettywriter.BlackBG, prettywriter.Green)
	// Step 5: Obfuscate timestamps to hide theoriginal encrypted parts order
	err = fileutils.ObfuscateFileTimestamps(outputDir)
	if err != nil {
		return fmt.Errorf("error obfuscating file timestamps: %w", err)
	}
	prettywriter.Writeln("[**] Timestamps successful altered", prettywriter.BlackBG, prettywriter.Green)
	fmt.Println("")
	prettywriter.WriteInBox(40, "Encryption finished", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	return nil
}

func (c *Core) Unhide(partsDir, outputPath string, prefilledPassword string) error {
	prettywriter.WriteInBox(40, "Configuration", prettywriter.Green, prettywriter.BlackBG, prettywriter.DoubleLine)
	prettywriter.Writeln("[==] Chosen mode: unhide (decrypting)", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[==] Input directory: "+partsDir, prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[==] Output directory: "+outputPath, prettywriter.Green, prettywriter.BlackBG)
	fmt.Println("")

	prettywriter.WriteInBox(40, "Starting Decryption Process", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	fmt.Println("")

	// Step 1: Decrypt Master Lock File
	password := ""
	if "" == prefilledPassword {
		password = utils.PromptForPassword("Enter the password to decrypt the masterlock: ")
	} else {
		password = prefilledPassword
	}

	fmt.Println()
	prettywriter.WriteInBox(40, "Handling masterlock", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	prettywriter.Writeln("[>>] Reading masterlock", prettywriter.BlackBG, prettywriter.Green)
	encryptedMasterLock, err := ioutil.ReadFile(filepath.Join(partsDir, "masterlock"))
	if err != nil {
		return fmt.Errorf("error reading encrypted master lock file: %w", err)
	}

	prettywriter.Writeln("[>>] Decrypting masterlock", prettywriter.BlackBG, prettywriter.Green)
	decryptedMasterLock, err := encryptor.DecryptWithPassword(string(encryptedMasterLock), password)
	if err != nil {
		return fmt.Errorf("error decrypting master lock file: %w", err)
	}

	prettywriter.Writeln("[>>] Unpacking masterlock", prettywriter.BlackBG, prettywriter.Green)
	var mlock masterlock.MasterLock
	err = json.Unmarshal([]byte(decryptedMasterLock), &mlock)
	if err != nil {
		return fmt.Errorf("error unmarshaling master lock file: %w", err)
	}
	prettywriter.Writeln("[**] Masterlock handled successful", prettywriter.BlackBG, prettywriter.Green)

	// Step 2: Decrypt Each Part
	var base64Parts []string
	fmt.Println("")
	prettywriter.WriteInBox(40, "Handling encrypted parts", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	for i, partInfo := range mlock.Parts {
		fmt.Print("\r")
		prettywriter.Write("[>>] Decrypting parts : "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(mlock.Parts)), prettywriter.Green, prettywriter.BlackBG)
		partPath := filepath.Join(partsDir, partInfo.Filename)
		encryptedPart, err := os.ReadFile(partPath)
		if err != nil {
			return fmt.Errorf("error reading encrypted part file: %w", err)
		}

		decryptedPart, err := encryptor.DecryptWithRandomKey(string(encryptedPart), string(partInfo.Key))
		if err != nil {
			return fmt.Errorf("error decrypting part: %w", err)
		}

		base64Parts = append(base64Parts, string(decryptedPart))
	}
	fmt.Println("")
	prettywriter.Writeln("[**] Parts decrypted successful ", prettywriter.Green, prettywriter.BlackBG)
	fmt.Println("")

	prettywriter.WriteInBox(40, "Handling zip", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)
	// Step 3: Reconstruct Base64 Data
	paddedData := strings.Join(base64Parts, "")
	unpaddedData := paddedData[:len(paddedData)-mlock.Padding]
	prettywriter.Writeln("[**] Reconstructed base64 string without padding ", prettywriter.Green, prettywriter.BlackBG)
	prettywriter.Writeln("[>>] Unpacking zip data ", prettywriter.Green, prettywriter.BlackBG)
	zipper := zipper.New()
	err = zipper.ExtractBase64ZipToDir(unpaddedData, outputPath)
	if err != nil {
		return fmt.Errorf("error unzipping data: %w", err)
	}
	prettywriter.Writeln("[**] Successfully unpacked zip data ", prettywriter.Green, prettywriter.BlackBG)
	fmt.Println()
	prettywriter.WriteInBox(40, "Decryption finished", prettywriter.Black, prettywriter.Green, prettywriter.DoubleLine)

	return nil
}
