package splitter

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Splitter interface {
	// Split the file into chunks based on provided configs
	Split(configs SplitConfigs)

	// Merge the provided directory into a single file
	Merge()
}

type Split struct {
	EncryptionConfig EncryptionConfig
	FileName         string // Name of the file at FilePath location
	FilePath         string // Path to the location of the file
}

type EncryptionConfig struct {
	IsEncrypted bool
	Password    string
}

type SplitConfigs struct {
	ChunkSize int    // The size of each chunk in format (ChunkSize:1, Format:MB => 1MB)
	Format    string // B, KB, MB, GB
}

type splitterDetails struct {
	numChunks     int64
	totalFileSize int64
	fileSize      int
	dirPath       string
}

var outputDirPrefix string = "file-data_"
var checksumFileName string = "checksum"

func (s *Split) Split(configs SplitConfigs) error {
	filePath := cleanFilePath(s.FilePath, s.FileName)

	// Prepare for data copying
	details, err := prepareForSplitting(s, configs)
	if err != nil {
		return fmt.Errorf("failures encountered during preparation phase with error: %s", err)
	}
	checksumFilePath := details.dirPath + "/" + checksumFileName

	var wg sync.WaitGroup
	errorCh, workerCount := make(chan error, 2), 2

	wg.Add(workerCount)
	go splitFile(details, filePath, &wg, errorCh)
	go createChecksumFile(filePath, checksumFilePath, &wg, errorCh)
	go func() {
		wg.Wait()
		close(errorCh)
	}() // Receiver process cannot close, must be done by sender (child-goroutine)

	for err = range errorCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareForSplitting(s *Split, configs SplitConfigs) (*splitterDetails, error) {
	filepath := cleanFilePath(s.FilePath, s.FileName)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("provided filepath does not exist with error: %s", err)
	}
	defer file.Close()

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failure to read file details with error: %s", err)
	}

	fileFormatSize, err := getFileFormatSize(configs.Format)
	if err != nil {
		return nil, err
	}
	fileSize := configs.ChunkSize * fileFormatSize

	dirPath, err := createOutputDirectory(s.FileName)
	if err != nil {
		return nil, err
	}

	return &splitterDetails{
		numChunks:     fileInfo.Size()/int64(fileSize) + 1,
		totalFileSize: fileInfo.Size(),
		fileSize:      fileSize,
		dirPath:       dirPath,
	}, nil
}

func cleanFilePath(path, fileName string) string {
	lastChar := path[len(path)-1]
	if lastChar == '/' { // Trailing slash exists
		return path + fileName
	}
	return path + "/" + fileName // Add trailing slach to path
}

func createOutputDirectory(fileName string) (string, error) {
	newDirPath := filepath.Join(".", outputDirPrefix+fileName)

	if err := os.Mkdir(newDirPath, os.ModePerm); os.IsExist(err) {
		return "",
			fmt.Errorf("directory for file '%s' already exists at '%s'", fileName, newDirPath)
	}
	return newDirPath, nil
}

func numDigits(i int64) int {
	if i < 10 {
		return 1
	}
	count := 0
	for i != 0 {
		i /= 10
		count++
	}
	return count
}

func getFileFormatSize(format string) (int, error) {
	kilo := 1024
	switch format {
	case "B":
		return 1, nil
	case "KB":
		return kilo, nil
	case "MB":
		return kilo * kilo, nil
	case "GB":
		return kilo * kilo * kilo, nil
	default:
		return 0, fmt.Errorf("Invalid provided format:" + format)
	}
}

func splitFile(details *splitterDetails, filePath string, wg *sync.WaitGroup, ch chan error) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		ch <- fmt.Errorf("failure to open file, error: %s", err)
		return
	}
	defer file.Close()

	fileBuffer := make([]byte, details.fileSize)
	formatDigitsStr := "%0" + strconv.Itoa(numDigits(details.numChunks)) + "d"
	outputFileNamePrefix := details.dirPath + "/" + "data_"

	for chunkCount := range details.numChunks {
		numReadBytes, err := file.Read(fileBuffer)
		if err != nil && err != io.EOF {
			ch <- fmt.Errorf("error while reading file - error: %s", err)
			return
		}
		if numReadBytes == 0 {
			break
		}

		name := outputFileNamePrefix + fmt.Sprintf(formatDigitsStr, chunkCount)
		dstFile, err := os.Create(name)
		if err != nil {
			ch <- fmt.Errorf("failure to create a destination file with error: %s", err)
			return
		}
		defer dstFile.Close()

		_, err = dstFile.Write(fileBuffer[:numReadBytes])
		if err != nil {
			ch <- fmt.Errorf("failure to write destination file with error: %s", err)
			return
		}
	}
}

func createChecksumFile(filePath, outputPath string, wg *sync.WaitGroup, ch chan error) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		ch <- fmt.Errorf("provided filename does not exist with error: %s", err)
		return
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		ch <- fmt.Errorf("error while hashing the input file with error: %s", err)
		return
	}

	checkSumStr := fmt.Sprintf("%x", hash.Sum(nil))
	if err = os.WriteFile(outputPath, []byte(checkSumStr), 0644); err != nil {
		ch <- fmt.Errorf("failure to write checksum file: %s", err)
	}
}

func (s *Split) Merge() error {
	if _, err := os.Stat(s.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("provided directory path does not exist with error: %s", err)
	}

	fileName := s.FileName
	if len(fileName) == 0 {
		paths := strings.Split(s.FilePath, "/")
		if len(paths) > 1 {
			fileName = paths[len(paths)-1]
		} else {
			fileName = paths[0]
		}
		fileName = strings.Split(fileName, outputDirPrefix)[1] // Remove the prefix
	}

	entries, err := os.ReadDir(s.FilePath)
	if err != nil {
		return fmt.Errorf("failure to read contents of directory %s with error: %s", s.FilePath, err)
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create the new output file with error: %s", err)
	}
	defer outputFile.Close()

	for _, e := range entries {
		if e.Name() == checksumFileName {
			continue // Skip checksum file until all pieces are merged together
		}

		file, err := os.Open(cleanFilePath(s.FilePath, e.Name()))
		if err != nil {
			return fmt.Errorf("failed to open file with error: %s", err)
		}

		bytesCopied, err := io.Copy(outputFile, file)
		if err != nil {
			return fmt.Errorf("failed to write from '%s' to '%s' with error: %s",
				e.Name(), outputFile.Name(), err)
		}
		fmt.Printf("bytesCopied: %v\n", bytesCopied)
	}

	// Compare checksums
	doesChecksumMatch, err := compareChecksums(fileName, cleanFilePath(s.FilePath, checksumFileName))
	if err != nil {
		return fmt.Errorf("error comparing checksums: %s", err)
	}
	if !doesChecksumMatch {
		return fmt.Errorf("checksums do not match")
	}

	fmt.Printf("doesChecksumMatch: %v\n", doesChecksumMatch)

	return nil
}

func compareChecksums(outputFilePath, checksumFilePath string) (bool, error) {
	file, err := os.Open(outputFilePath)
	if err != nil {
		return false, fmt.Errorf("failure to open file for checksum comparison with error: %s", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("failure to copy file contents for hashing: %s", err)
	}

	newFileHash := fmt.Sprintf("%x", hash.Sum(nil))

	fileBytes, err := os.ReadFile(checksumFilePath)
	if err != nil {
		return false, fmt.Errorf("error reading checksum file: %s", err)
	}

	return newFileHash == string(fileBytes), nil
}
