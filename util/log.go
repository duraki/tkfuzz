package util

import (
	"io"
	"log"
	"os"
)

var Err = log.Default()
var Log = log.Default()

/**
var Err = log.Default()
var Log = log.Default()
**/

// var Log log.Default()
var LoggingEnabled = true

var DEBUG = Ternary(os.Getenv("DBG") != "", true, false) // os.Getenv("DBG")

func Init() {
	Err = log.New(os.Stderr, "", 0)
	Log = log.New(os.Stdout, "", 0)
	initLogging()
	LoggingEnabled = true
}

func initLogging() {
	/* stderr log init */
	efile, eerr := os.OpenFile("logs/err.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if eerr != nil {
		Err.Fatal(eerr)
	}
	wre := io.MultiWriter(os.Stdout, efile)
	Err.SetOutput(wre)
	Err.SetFlags(log.LstdFlags | log.Lshortfile)
	Err.SetPrefix("[err] ")

	/* SPAC=0 */

	/* stdout log init */
	file, err := os.OpenFile("logs/tkfuzz.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Err.Fatal(err)
	}
	defer file.Close()
	defer efile.Close()
	wrt := io.MultiWriter(os.Stdout, file)
	Log.SetOutput(wrt)
	Log.SetFlags(log.LstdFlags | log.Lshortfile)
	Log.SetPrefix("[log] ")
	Log.Println("tkfuzz is starting . logging initialised.")
	Err.Print("tkfuzz stderr ~~~~ 			valid. ")
}
