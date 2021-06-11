package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
)

type File struct {
	Filesize  int16 `json:"filesize"`
	Extension string `json:"extension"`
	FilePath  string `json:"filepath"`
}

type Files struct {
	Files []File `json:"files"`
}
var files Files

var wg sync.WaitGroup

// parsing the directories with go routines
// when in the flow if a directory is found that is parsed by go routine
// no of go-routines will be equal to the number of sub-directories in a directory
func parse(dir string) {
	defer wg.Done()

	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && path != dir {
			wg.Add(1)
			go parse(path)
			return filepath.SkipDir
		}
		if f.Mode().IsRegular() {
			files.Files = append(files.Files, File{Filesize: int16(f.Size()), Extension: filepath.Ext(path), FilePath: path})
		}
		return nil
	}

	filepath.Walk(dir, visit)
}

//for sending post requests to the server
// number of go routines will be decided by the user

func sendHttpRequest(num int, url string, data Files) {
	defer wg.Done()
	fmt.Println(num,len(data.Files))
	buff, err := json.Marshal(data)
	if err != nil {
		fmt.Println(num," : Couldn't send the Request",err)
	}
	filesPost,filesPostError:= http.Post(url,"application/json",bytes.NewBuffer(buff))
	if filesPostError != nil{
		fmt.Println(num,"Couldn't connect to Server :",filesPostError.Error())
		return
	}
	if filesPost.StatusCode == 200 {
		fmt.Println(num,"Post Successful")
	} else{
		var ErrorResp map[string]interface{}
		err = json.NewDecoder(filesPost.Body).Decode(ErrorResp)
		if err != nil{
			fmt.Println(num,"Couldn't decode the failed response from the server")
		}
		fmt.Println(num,"Error from Server :",ErrorResp)
	}
}

func main() {
	// for parsing user provided flags
	pathFlag := flag.String("path", `C:\Users\SiddhiManojManojVarm\Documents`, "path required")
	urlFlag := flag.String("url", "http://localhost:1999/", "url required")
	helpFlag := flag.Bool("help", false, "flag required")
	flag.Parse()
	path := *pathFlag
	url := *urlFlag
	help := *helpFlag
	fmt.Println(path)
	fmt.Println(url)
	fmt.Println(help)
	// if help is true the cpuprofile and memory profile will be saved to files
	if help {
		cpuprofile, err := os.Create("CPUProfile.prof")
		if err != nil {
			fmt.Println("could not create CPU profile: ", err)
		}
		defer cpuprofile.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(cpuprofile); err != nil {
			fmt.Println("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
		memoryprofile, err := os.Create("MemoryProfile.prof")
		if err != nil {
			fmt.Println("could not create memory profile: ", err)
		}
		defer memoryprofile.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(memoryprofile); err != nil {
			fmt.Println("could not write memory profile: ", err)
		}
	}
	// parsing the directories and saving all the files to a list
	wg.Add(1)
	parse(path)
	wg.Wait()
	fmt.Println(len(files.Files))
	// requests will be sent based on the number of routines
	routines := 5
	if len(files.Files) < 5 {
		func() {
			wg.Add(1)
			go sendHttpRequest(0, url, files)
		}()
		wg.Wait()
		return
	}

	for i := 0; i < routines; i = i + 1 {
		if i < 4 {
			func() {
				wg.Add(1)
				go sendHttpRequest(i, url, Files{Files: files.Files[int(i*len(files.Files)/5):int((i+1)*len(files.Files)/5)]})
			}()
		} else {
			func() {
				wg.Add(1)
				go sendHttpRequest(i, url, Files{Files: files.Files[int(i*len(files.Files)/5):]})
			}()
		}
	}
	wg.Wait()
}
