# File Splitter
A file splitting utility to break large files into smaller chunks

## Overview
A simple microservice that handles the process of chunking a file and rebuilding it back up into a single file.

Allows for encryption of each piece, with automatic passkey generation or extendible to allow for manual key insertion.

## How to run it
Navigate to the run directory and execute the appropriate commands

Application needs 4 inputs to run
1. file path to which file want to split
2. chunk format (K, KB, MB, GB)
3. chunk size in that format
4. encryption flag
5. Decryption file location

When rebuilding the file, automatic checksum verification is run to validate all pieces came together with notification of success/failure.

## Operations

#### Spliting a file
```
Basic split - 10MB chunks

go run main.go --op split --filePath ~/Downloads/testfile.txt
```

```
Specified chunk size - 100KB

go run main.go --op split --filePath ~/Downloads/testfile.txt --size 100 --format KB
```

```
Encryption enabled

go run main.go --op split --filePath ~/Downloads/testfile.txt -e
```

#### Merge File
```
Basic merge

go run main.go --op merge --filePath file-data_testfile.txt
```

```
Decrypting merge

go run main.go --op merge --filePath file-data_testfile.txt -e
```

```
Decrypting merge with separate key file

go run main.go --op merge --filePath file-data_testfile.txt -e -ef ~/keyFile.txt
```

## License
This package is licensed under the MIT License - see the LICENSE file for details.
