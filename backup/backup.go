package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/kardianos/osext"
)

var totalFile, totalCopiedFile int

// var modifyLog, modLogErr os.Open
// var totalLog, totalLogErr os.Open

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
			log.Println("---->New file detected as", fromFileName)
			fcerror = CopyFile(fromPath, toPath)
			copyOccured = true
		} else if fileTo.Size() != fileFrom.Size() {
			fmt.Println("---->Modification detected for", fromFileName)
			log.Println("---->Modification detected for", fromFileName)
			fcerror = CopyFile(fromPath, toPath)
			copyOccured = true
		}

	} else {
		fmt.Println("||ERROR||: AH! Destination file missing.", err)
	}

	if fcerror != nil {
		fmt.Println("||ERROR||: Copy failed", fcerror)
		log.Println("Copied Failed", fcerror)
	} else if copyOccured {
		totalCopiedFile = totalCopiedFile + 1
		fmt.Println("||INFO||: Successfully copied from", fromPath, " to ", toPath)
		log.Println("Copied from", fromPath, " to ", toPath)
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
				totalFile = totalFile + 1
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
	modifyFileName := GetLogPath()
	os.Remove(modifyFileName)
	mlogf, err := os.OpenFile(modifyFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err == nil {
		defer mlogf.Close()

		log.SetOutput(mlogf)
		// log.Println("This is a test log entry")

		jsonByte, err := ioutil.ReadFile(path)
		if err == nil {
			jo := []DirectoryDataConfiJSON{}

			marshalErr := json.Unmarshal(jsonByte, &jo)
			if marshalErr != nil {
				fmt.Println("JSON Unmarshed error. ", marshalErr)
			} else {
				totalChannel := len(jo)
				fileProgress := make(chan bool, totalChannel)

				// var modifyLog, modLogErr = os.Open("./modified.log")
				// var totalLog, totalLogErr os.OpenFile

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

		fmt.Println("\nTotal files:", totalFile, ", New/Modified:", totalCopiedFile)

	} else {
		fmt.Println("Log open failed. For path ", path, ". Error::", err)

	}

}

// GetLogPath will return current directory
func GetLogPath() string {
	var filename string
	// if runtime.GOOS == "windows" {
	// 	filename = "BackMeUpModify.log"
	// } else {
	// 	filename = "/BackMeUpModify.log"
	// }
	filename = "BackMeUpModify.log"
	//GetMyPath Returns the absolute directory of this(pathfind.go) file

	return GetConfigPath(true) + "/" + filename
}

func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}

// GetConfigPath will return current directory
func GetConfigPath(pathOnly bool) string {
	var filename string
	// if runtime.GOOS == "windows" {
	// 	filename = "/backmeup.config.json"
	// } else {
	// 	filename = "/backmeup.config.json"
	// }
	filename = "backmeup.config.json"
	//GetMyPath Returns the absolute directory of this(pathfind.go) file

	dir, err := osext.ExecutableFolder()

	if err != nil {
		if pathOnly {
			return ""
		}

		return filename
	}
	// fmt.Println(folderPath)

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]) + filename)
	// if err != nil {
	// 	if pathOnly {
	// 		return ""
	// 	}

	// 	return "" + filename
	// }

	if _, err := os.Stat(dir + "/" + filename); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if !pathOnly {
			return filename
		}
		// return filename
	}
	if pathOnly {
		//dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		return dir
	}
	return dir + "/" + filename
}
