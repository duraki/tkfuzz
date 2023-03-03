package main

/* tkfuzz utility */
/**
 * Developed by Halis D. from durakiconsulting, LLC in collaboration with
 * Skullkey AB and Tillitis AB.
 *
 * --
 * Optional ENVs:
 * 		DBG=1 tkfuzz <mode> <tty>
 *		FTTY=/dev/pts/0 tkfuzz <mode>
 *
 * Optional Arguments:
 * 		tkfuzz <mode> <tty>
 * 	  	Fuzz Modes: fuzzfw OR fuzzsys
 *			fuzzfw: This mode asks the TKey if device is in firmware mode currently, and if so,
 *							it engages fuzz pipeline.
 *			fuzzsys: This mode tries to fuzz on UART RX/TX during other system lifecycle and phases,
 *							 such is during the app. load, et al.
 * 		TTY:
 *			tty argument should point to serial interface in /pts/dev/x or /dev/tty
 **/

import (
	"fmt"
	"os"
	s "strings"
	"tkfuzz/fuzzer"
	"tkfuzz/tk1"
	"tkfuzz/util"
)

const (
	Author    = "Halis Duraki"
	Company_A = "durakiconsulting, LLC"
	Company_B = "Skullkey, AB"
	Year      = 2023
)

var stderr = util.Err
var Log = util.Log

var tty = os.Getenv("FTTY")
var fzmode = 0

func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("\t tkfuzz <fzmode> <tty>\n")
	fmt.Printf("\t tkfuzz fuzzfw|fuzzsys <tty>\n\n")
	fmt.Printf("Via FTTY ENV:\n")
	fmt.Printf("\tFTTY=/dev/pts/4\t tkfuzz fuzzfw\t # will fuzz fw mode\n")
	fmt.Printf("\tFTTY=/dev/pts/4\t tkfuzz fuzzsys\t # will fuzz sys/app mode\n")
	fmt.Printf("	 # FTTY= must be path to serial bus tty\n")
	fmt.Printf("\n")
	fmt.Printf("TK1 Protocol:\n")
	fmt.Printf("\t Bus Speed: %d\n", tk1.SerialSpeed)
	fmt.Printf("\t Serial Device: %s\n\n", "TILL/TKEY/TK")
}

func main() {

	util.Init()

	fmt.Printf("\n")
	fmt.Printf("Tillitis TK1 Firmware and OS fuzzer ❤️‍🔥\n")
	fmt.Printf("\t ... by %s,\n\tcollab. of %s & %s (%d)\n", s.ToLower(Author), Company_A, Company_B, Year)
	fmt.Printf("\n")

	args := os.Args[1:] // args, w/o binary

	/* preliminary arg check */
	if len(args) <= 0 {
		stderr.Fatal("incorrect CLI usage of tkfuzz binary, please use:\n\t", "tkfuzz fuzzfw <path/to/tty>")
	}

	mode := args[0]
	if mode == "" {
		stderr.Fatal("fuzz mode is required, pick 'fuzzfw' OR 'fuzzsys'")
	}
	if !(s.Contains(mode, "fuzzfw") || s.Contains(mode, "fuzzsys")) {
		stderr.Fatal("fuzz mode is not correct, pick 'fuzzfw' OR 'fuzzsys' as a 1st arg")
	} else {

	}

	if tty == "" { /* first conditional, tty is inited very early via FTTY= */
		/* otherwise, lets try to get tty from the CLI args */
		if len(args) > 1 {
			tty = args[1] // arg n+1 equals to serial tty path
		}
	}

	if tty == "" {
		printUsage()

		Log.Printf("tty has not been set, use FTTY env variable or pass it as argument:\n\t%s\n", "tkfuzz fuzzfw /dev/tty/0")
	} else {
		Log.Printf("serial tty is set to:\n\t%s\n", tty)
	}

	if s.Contains(mode, "fuzzfw") {
		fzmode = fuzzer.F_FW
	} else {
		fzmode = fuzzer.F_SYS
	}

	Log.Printf("fuzz mode selected:\n\t%s [IDx: %d]\n", mode, fzmode)

	iteration_no := 100 * 1024

	Log.Printf("fuzz iteration ctx:\n\t%d [100 * 1024]\n", iteration_no)

	/* contains fuzz cfg provided via cli args */
	fuzzer.CxFuzzcfg = fuzzer.InitFuzzCfg(iteration_no, fzmode, tty)
	fuzzer.CxFuzzer = fuzzer.InitFuzzer(fuzzer.CxFuzzcfg, fuzzer.FUZZ_THREADS)

	Log.Println("fuzzer initialized and configured properly. now dumping Fuzzer container ctx:")
	Log.Println(fuzzer.CxFuzzer)

	Log.Printf("%s", s.Repeat("*", 100))
	_ = fuzzer.CxFuzzer.Describe()
	Log.Printf("%s", s.Repeat("*", 100))
	_ = fuzzer.CxFuzzcfg.Describe()
	Log.Printf("%s", s.Repeat("*", 100))
	fuzzer.CxPipeline = fuzzer.CreateFuzzPipeline(fuzzer.CxFuzzer)
	Log.Printf("\t\tFuzzPipeline=%v", fuzzer.CxPipeline)
	_ = fuzzer.CxPipeline.Describe()
	Log.Printf("%s", s.Repeat("*", 100))

	if util.AskForConfirmation() {
		if fuzzer.CountDown() != false {
			_ = fuzzer.CxPipeline.Start()
		}
	} else {
		Log.Println("Fuzz Pipeline is destroyed during the confirmation. Will exit now.")
		return
	}
}
