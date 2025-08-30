package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/AlexanderEl/splitter"
)

func main() {
	operation := flag.String("op", "split", "The operation to perform (split/merge)")
	filePath := flag.String("filePath", "", "Path to file to split or directory to merge")
	chunkSize := flag.Int("size", 10, "Size of each file in format")
	chunkFormat := flag.String("format", "MB", "File format [B, KB, MB, GB]")
	encryption := flag.Bool("e", false, "Encryption flag")
	decryptionFilePath := flag.String("ef", "", "Decryption key file (default: current directory)")
	flag.Parse()

	switch strings.ToLower(*operation) {
	case "split":
		splitFunc(filePath, chunkFormat, chunkSize, encryption)
	case "merge":
		mergeFunc(filePath, decryptionFilePath, encryption)
	default:
		fmt.Println("Invalid operation given")
	}
}

func splitFunc(filePath, chunkFormat *string, chunkSize *int, encryption *bool) {
	paths := strings.Split(*filePath, "/")
	fileName := paths[len(paths)-1]
	path := strings.Join(paths[:len(paths)-1], "/")

	split := splitter.Split{
		FileName:    fileName,
		FilePath:    path,
		IsEncrypted: *encryption,
	}
	configs := splitter.SplitConfigs{
		ChunkSize: *chunkSize,
		Format:    *chunkFormat,
	}

	enabledStr := "enabled"
	if !*encryption {
		enabledStr = "disabled"
	}
	fmt.Printf("Splitting file '%s' with encryption %s\n", fileName, enabledStr)

	err := split.Split(configs)
	if err != nil {
		fmt.Println("Error while splitting file:", err)
	} else {
		fmt.Printf("Successfully split file '%s'\n", split.FileName)
	}
	if runtime.NumGoroutine() > 1 {
		fmt.Println("Zombie goroutines exist!!", runtime.NumGoroutine())
	}
}

func mergeFunc(filePath, decryptionFilePath *string, encryption *bool) {
	split := splitter.Split{
		IsEncrypted: *encryption,
		FilePath:    *filePath,
	}

	enabledStr := "enabled"
	if !*encryption {
		enabledStr = "disabled"
	}
	fmt.Printf("Merging files in '%s' with encryption %s\n", *filePath, enabledStr)

	err := split.Merge(*decryptionFilePath)
	if err != nil {
		fmt.Println("Errors encountered during merge:", err)
	} else {
		fmt.Printf("Successfully merged file from '%s'\n", *filePath)
	}
}
