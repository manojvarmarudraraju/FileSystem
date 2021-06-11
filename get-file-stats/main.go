package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// stats struct for marshaling and unmarshalling
type Stats struct {
	FilesCount int64 `json:"FilesCount"`
	MaxFileSize int64 `json:"MaxFileSize"`
	AverageFileSize float64 `json:"AverageFileSize"`
	ListExtensions []string `json:"ListExtensions"`
	FileExtensionsCount map[string]int `json:"FileExtensionsCount"`
	LatestPath []string `json:"LatestPath"`
}

func main(){



	// parsing flags provided by the user while running main.go
	logger := log.New(os.Stdout,"file-stats : ",log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	intervalFlag := flag.Int("interval",10,"Interval at which get requests happen")
	hostFlag := flag.String("host","http://localhost:1999/Stats","host flag")
	helpFlag := flag.Bool("help",false,"help flag")
	flag.Parse()

	help := *helpFlag
	// for saving memory profile and cpu profile into files
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
	interval := *intervalFlag
	host := *hostFlag
	logger.Println(interval)
	logger.Println(host)
	// http client to connect to the server and get the stats data at every interval
	httpClient := http.Client{}
	for ;true;{
		// creating request
		req,err := http.NewRequest(http.MethodGet,host,nil)
		if err != nil{
			logger.Println("failed to create request")
			break
		}
		// sending get request to server
		statsInfo,statsFetchError := httpClient.Do(req)
		if statsFetchError != nil{
			logger.Println("Couldn't connect to Server :",statsFetchError.Error())
		}
		// if get request is successful
		if statsInfo.StatusCode == 200 {
			var StatsInfo Stats
			err = json.NewDecoder(statsInfo.Body).Decode(&StatsInfo)
			if err != nil {
				logger.Println("Couldn't decode the successful response from the server")
			}
			logger.Println("Data from Server :",StatsInfo)
		} else{
			var ErrorResp map[string]interface{}
			err = json.NewDecoder(statsInfo.Body).Decode(ErrorResp)
			if err != nil{
				logger.Println("Couldn't decode the failed response from the server")
			}
			logger.Println("Error from Server :",ErrorResp)
		}
		// stop the execution for interval provided by the user
		time.Sleep(time.Duration(interval) * time.Second)
	}

}
