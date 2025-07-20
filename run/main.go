package main

import (
	"fmt"
	"runtime"
	"splitter"
)

func main() {
	// spltFunc()
	mergeFunc()
}

func spltFunc() {

	// inputFileName := flag.String("input file name", "", "Input file to split")
	// chunkSize := flag.Int("size", 10, "Size of each file in format (default: 10MB)")
	// chunkFormat := flag.String("format", "MB", "File format [B, KB, MB, GB] (default: MB)")
	// flag.Parse()

	split := splitter.Split{
		FileName: "zoom_amd64.deb",
		FilePath: "/home/alex/Downloads/",
		EncryptionConfig: splitter.EncryptionConfig{
			IsEncrypted: false,
		},
	}
	configs := splitter.SplitConfigs{
		ChunkSize: 1,
		Format:    "MB",
	}
	err := split.Split(configs)
	if err != nil {
		fmt.Println("Error while splitting file:", err)
	}
	if runtime.NumGoroutine() > 1 {
		fmt.Println("Zombie goroutines exist!!", runtime.NumGoroutine())
	}

	fmt.Printf("Successfully split file %s", split.FileName)
}

func mergeFunc() {
	fullPath := "/home/alex/Development/learning-go/splitter/run/file-data_zoom_amd64.deb"

	split := splitter.Split{
		EncryptionConfig: splitter.EncryptionConfig{
			IsEncrypted: false,
		},
		FilePath: fullPath,
	}

	err := split.Merge()
	if err != nil {
		fmt.Println("Errors encountered during merge:", err)
	}
	fmt.Println("Successful file merge")
}
