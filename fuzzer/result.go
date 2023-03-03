package fuzzer

import "time"

// var Log = util.Log
// var Err = util.Err

const (
	desc string = "This code contans FuzzPipeline and the fuzz result triggers."
)

// #/	 /** fuzz modes */
// #/
// #/	 const (
// #/	 					F_FW  int = 0 // firmware fuzz mode
// #/	 					F_SYS     = 1 // post-firmware fuzzing
// #/	 )

// #/	/** 	This struct contains Fuzzing Results of each fuzz iteration,
// #/	/**		both for valid and invalid results. **/
type FuzzResult struct {
	created_at  time.Time
	proof_perct int    // percentage % of payload validity, ie. 25% = bad 	OR 	  85% = good
	byte_bin    string // fuzz payload (binary)
	byte_hex    string // fuzz payload (hex)
	rx_bus      string // rx line data
	tx_bus      string // tx line data
}

// #/	/**		FuzzResults Contextual container. Holds FuzzResult metadata, settings
// #/	/**		statistics, fuzzed payloads, so on.
type FuzzResultContext struct {
	// #/ 	   /** Todo: perhaps change the fr_id (int) type to some custom interface a la ObjID. */
	fr_id      int                // FuzzResults Identified
	fr_valid   []FuzzResult       // triggered fuzz payloads
	fr_invalid []FuzzResult       // invalid, but still iterated payloads (can be used later to skip over)
	setting    FuzzResultSettings //
	device_tty string             // same as FuzzPipeline- FuzzConfig- FTTY ... easier to indicate
	s_total    int                // size_ fr_valid + fr_invalid - total fuzz results executed
}

// #/	/**		FuzzResultsSettings as a part of the contextual interface.
type FuzzResultSettings struct {
	outfile  string // file path pointing to result output logs
	exported bool   // if exported, warn, and decline
}
