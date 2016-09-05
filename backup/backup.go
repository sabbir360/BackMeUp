package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// CopyFile Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

//CompareFileAndProcess if new found
func CompareFileAndProcess(fromPath, toPath string) {
	var fcerror error
	copyOccured := false
	fileFrom, err := os.Stat(fromPath)

	fromFileName := fileFrom.Name()

	if err == nil {
		fileTo, err := os.Stat(toPath)

		if err != nil {
			fmt.Println("---->New file detected as", fromFileName)
			fcerror = CopyFile(fromPath, toPath)
			copyOccured = true
		} else if fileTo.Size() != fileFrom.Size() {
			fmt.Println("---->Modification detected for", fromFileName)
			fcerror = CopyFile(fromPath, toPath)
			copyOccured = true
		}

	} else {
		fmt.Println("||ERROR||: AH! Destination file missing.", err)
	}

	if fcerror != nil {
		fmt.Println("||ERROR||: Copy failed", fcerror)
	} else if copyOccured {
		fmt.Println("||INFO||: Successfully copied from", fromPath, " to ", toPath)
	}

}

// ReadDirectory reads a directory
func ReadDirectory(fileProgress chan bool, directory string, directoryTo string) {
	files, err := ioutil.ReadDir(directory)

	fmt.Println("\n-->Scanning:", directory)

	if err == nil {
		for _, file := range files {
			readObj := directory + "/" + file.Name()
			toReadObj := directoryTo + "/" + file.Name()

			// if prefix != "" {
			// 	readObj = prefix + "/" + readObj
			// 	toReadObj = prefix + "/" + readObj
			// } else {
			// 	readObj = directory + "/" + readObj
			// 	toReadObj = directoryTo + "/" + readObj
			// 	// fmt.Println(readObj)
			// 	// var input string
			// 	// fmt.Scanln(&input)
			// }

			if info, err1 := os.Stat(readObj); err1 == nil && info.IsDir() {
				// fmt.Println(readObj, "is", "directory.")
				_, desterr := os.Open(toReadObj)
				if os.IsNotExist(desterr) {
					os.Mkdir(toReadObj, file.Mode())
				}
				chiledFileProgres := make(chan bool, 1)
				go ReadDirectory(chiledFileProgres, readObj, toReadObj)
				<-chiledFileProgres
			} else {
				// modificationTime := file.ModTime()
				// fmt.Println(readObj, "is", "file.", "Modification time:", modificationTime)
				// fmt.Println("Processing file-->", readObj, "With", toReadObj)

				CompareFileAndProcess(readObj, toReadObj)
			}
		}
	} else {
		fmt.Println(err)
	}
	fileProgress <- true
}

// DirectoryDataConfiJSON is struct used by JSON Parse
type DirectoryDataConfiJSON struct {
	SourceDir      string `json:"source_dir"`
	DestinationDir string `json:"destination_dir"`
}

// ReadConfig is a helper which execute JSON file
func ReadConfig(path string) {

	jsonByte, err := ioutil.ReadFile(path)
	if err == nil {
		jo := []DirectoryDataConfiJSON{}

		marshalErr := json.Unmarshal(jsonByte, &jo)
		if marshalErr != nil {
			fmt.Println(marshalErr)
		} else {
			totalChannel := len(jo)
			fileProgress := make(chan bool, totalChannel)
			for _, r := range jo {
				go ReadDirectory(fileProgress, r.SourceDir, r.DestinationDir)
			}

			for i := 0; i < totalChannel; i++ {
				<-fileProgress
			}

		}
	} else {
		fmt.Println("||ERROR||: Config Read Error!", err)
	}
}