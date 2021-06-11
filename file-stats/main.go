package main

import (
	"file-stats/handlers"
	"file-stats/internal/handlerFunctions"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var StatsVar *handlerFunctions.Stats
//StatsVar = new(handlerFunctions.Stats)
var count int

func AlterStats(logger *log.Logger, reqs <- chan handlerFunctions.FilesStruct){
	for rows := range reqs{
		presentAvg := StatsVar.AverageFileSize
		for i:=0; i < len(rows.Files);i++{
			presentAvg = float64(StatsVar.FilesCount)*presentAvg + float64(rows.Files[i].Filesize)
			StatsVar.FilesCount += 1
			StatsVar.AverageFileSize = presentAvg / float64(StatsVar.FilesCount)
			if StatsVar.MaxFileSize < int64(rows.Files[i].Filesize){
				StatsVar.MaxFileSize = int64(rows.Files[i].Filesize)
			}
			if StatsVar.FileExtensionsCount == nil{
				StatsVar.FileExtensionsCount = make(map[string]int64)
			}
			if StatsVar.ListExtensions == nil{
				StatsVar.ListExtensions= make([]string,0)
			}
			if StatsVar.LatestPath == nil{
				StatsVar.LatestPath= make([]string,0)
			}
			if StatsVar.FileExtensionsCount == nil{
				StatsVar.FileExtensionsCount = make(map[string]int64)
			}
			if _,ok := StatsVar.FileExtensionsCount[rows.Files[i].Extension];ok{
				StatsVar.FileExtensionsCount[rows.Files[i].Extension] += 1
			} else{
				StatsVar.FileExtensionsCount[rows.Files[i].Extension] = 1
				StatsVar.ListExtensions = append(StatsVar.ListExtensions, rows.Files[i].Extension)
			}
			if len(StatsVar.LatestPath) < 10{
				StatsVar.LatestPath = append(StatsVar.LatestPath , rows.Files[i].FilePath)
			} else{
				_,StatsVar.LatestPath = StatsVar.LatestPath[0],StatsVar.LatestPath[1:]
				StatsVar.LatestPath = append(StatsVar.LatestPath,rows.Files[i].FilePath)
			}
		}
		logger.Println(StatsVar)
	}
}

func main() {
	portFlag := flag.String("port", "1999", "Port number of the server")
	helpFlag := flag.Bool("help", false, "Statistics to be displayed")
	httpsFlag := flag.Bool("https", false, "HTTPS or not")

	flag.Parse()

	port := *portFlag
	help := *helpFlag
	https := *httpsFlag

	fmt.Println(port)
	fmt.Println(help)
	fmt.Println(https)
	// for cpu profiling and memory profiling
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

	logger := log.New(os.Stdout,"file-stats : ",log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// channel to avoid memory access at the same time by the handler go routines
	var requestInput = make(chan handlerFunctions.FilesStruct)
	StatsVar = new(handlerFunctions.Stats)
	go AlterStats(logger,requestInput)
	api := http.Server{
		Addr:    port,
		Handler: handlers.API(logger,requestInput,StatsVar),
	}

	serverError := make(chan error, 1)

	go func() {
		if https {
			// https server
			log.Println("HTTPS Server is listening at localhost: ",port)
			serverError <- api.ListenAndServeTLS("server.crt","server.key")
		} else {
			// http server
			log.Println("HTTP Server is listening at localhost: ",port)
			serverError <- api.ListenAndServe()
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <- serverError:
		log.Fatal(err)
	case <- shutdown:
		log.Fatal("Main: shutting down the api server")
		api.Close()
	}

}
