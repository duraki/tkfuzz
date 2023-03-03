package fuzzer

import (
	"fmt"
	"os"
	s "strings"
	"time"
	"tkfuzz/tk1"
	"tkfuzz/util"
)

var max_pgrs int = CxFuzzcfg.iteration // max = iteration * sizeof(total_payload_bytes) + encoders
// #/								    /* max permutations during fuzz */
var crr_pgrs int = 0 /* global prog event handler, increases on each loop */
// todo: add remaining_time = maxim.

// var Log = util.Log
var Err = util.Err

var fh *os.File     // contains fh of the fuzzing exec. results
var TKh = tk1.New() // provides interface to communicate with TKey

type ResultMap struct {
	valid_sum       int           // total valid as of now
	invalid_sum     int           // total invalid as of now
	itr_no          int           // iteration number (from fuzzloop size)
	approx_timeleft time.Duration // approx. elapsed time
}

type FuzzingProgress struct {
	pipeline FuzzPipeline
	tk1h     *tk1.TillitisKey     // pointer to GLOBAL tk1h
	Results  map[string]ResultMap //{total int, valid int}
	// #/	 /* 			r_rm = make(map([string]ResultMap)
	// #/	 /* 			r_rm[util.GEN_RM_KEY_ID] = ResultMap {
	// #/						VALID_SUM, INVALID_SUM, FUZZ_ITERATION_NUMBER, TIME_LEFT_CALC()
	// #/	 /* 			}
	// #/
	/** Important to note that the ResultMap does not contain all fuzz
	 * information, but only a small part of it, used as a progress
	 * calculator.
	 *
	 * The Results struct variable is a kvmap where key acts as a
	 *	identifier for multiple ResultMap appends.
	 * The form of the ResultMap is as following:
	 *
	 * 		TKFUZZ_(device)_(fuzzmode)_BAR_(bitrate)_ // (maybe?) _DateTime 	OR 	_iteration_index ...
	 *
	 * Each string based key contains FuzzRM as per schematics above.
	 *
	 *				* do not append manually, use a Builder abstraction *
	 **/
}

type FuzzPipeline struct {
	/** process timing */
	t_start    time.Time     // from time.Now()
	t_end      time.Time     // from time.Now()
	t_elapsed  time.Duration // Using time.Sub() to return duration diff.
	time_range string        // Using time.String() to return ascii time

	/** container holder **/
	fuzzer Fuzzer

	/** pipeline status */
	isrunning bool // is pipeline running?
	stopped   bool

	/** fuzz results ctx **/
	results FuzzResultContext // #/ /** updated once each X iteration 		 **/
	// #/							/** unlike log message buffering and 	 **/
	// #/							/** chunked out file writes.		 	 **/
}

func CreateFuzzPipeline(withFuzzer Fuzzer) FuzzPipeline {
	return FuzzPipeline{
		t_start:    time.Now(),
		t_end:      time.Now(),
		t_elapsed:  time.Duration(0),
		time_range: "",
		fuzzer:     withFuzzer,
		isrunning:  false,
		stopped:    false,
	}
}

func (fp FuzzPipeline) Start() bool {
	InitByteMaps()

	if fp.stopped == true {
		Err.Println("Unable to Start fuzzer pipeline. The pipeline was already stopped completely.")
		// todo: DumpDebug(...) and file write out (defer).
		os.Exit(0)
	}

	Log.Println("starting fuzz pipeline .... #idx unknown")
	fp.isrunning = true
	fp.t_start = time.Now()

	/* create a fuzz progress container */
	CxProgress = FuzzingProgress{
		pipeline: fp,
		tk1h:     TKh,
	}
	Log.Println("CxProgress defined containing =>\n\t", fmt.Sprintf("%v", CxProgress))
	if device_valid {
	}

	fp.SmokeTest()
	//_, _ = TKh.GetNameVersion()

	payload := []byte{0x50, 0x01}
	tk1.Dump("payload tx", payload)
	TKh.Write(payload)
	TKh.SetReadTimeout(2)
	rx, _, _ := TKh.ReadFrameCustom(0)
	TKh.SetReadTimeout(0)
	tk1.Dump("GetNameVersion rx", rx)

	os.Exit(0)

	PreparePayload(10, true)

	return true
}

func (fp FuzzPipeline) Pause() bool {
	Log.Println("pausing fuzz pipeline .... idling")
	fp.isrunning = false
	fp.t_end = time.Now()
	fp.t_elapsed = fp.calculate_elapsed()

	return true
}

func (fp FuzzPipeline) Stop() bool {
	Log.Println("Stopping fuzz pipeline .... If fuzzer is stopped, it can't be resumed; use SIG:'pause' instead.")
	// XXX: AskForConfirmation() ?
	fp.isrunning = false
	fp.t_end = time.Now()
	fp.t_elapsed = fp.calculate_elapsed()
	fp.stopped = true /* after stop(), the pipeline cant be recovered **/

	return true
}

func (fp FuzzPipeline) calculate_elapsed() time.Duration {
	difference := fp.t_end.Sub(fp.t_start)
	Log.Println("Time Start - Time End (Diff):")
	Log.Printf("\t diff = %v\n", difference)

	fp.t_elapsed = difference

	return difference
}

func (fp FuzzPipeline) Describe() string {
	Log.Println("")
	Log.Println("## Fuzz Pipeline")
	Log.Println("struct data:")
	Log.Printf(" time_start => %v", fp.t_start.String())
	Log.Printf(" time_end => %v", fp.t_end.String())
	Log.Printf(" elapsed => %v", fp.calculate_elapsed())
	Log.Println(" 	time_range NULL")
	Log.Println(" fuzzer? =>", fp.fuzzer)
	Log.Println(" Running? =>", fp.isrunning)

	return fmt.Sprintf("start=%v end=%v running=%v\n", fp.t_start, fp.t_end, fp.isrunning)
}

var once_asked = false   /* used for CountDown Preliminary Checks, ask retry only once */
var device_valid = false /* flag device validity */
func (fp FuzzPipeline) SmokeTest() bool {
	Log.Println("... pipeline smoke test ...")
	tk := tk1.New()
	devpath := fp.fuzzer.fzcfg.ttypath
	Log.Println("Connecting to device on serial port ... ", devpath)

	if err := tk.Connect(devpath, tk1.WithSpeed(tk1.SerialSpeed)); err != nil {
		Err.Println("Could not open Serial TTY device: ", devpath)
		Err.Println("resp err: ", err)
		device_valid = false
		return false
	}

	Log.Println("All seems good, assigning ptr tk => TKh ... 🔱")
	device_valid = true
	TKh = tk

	InitByteMaps() // calls tech.go# and initialises the registered bytes in tk fw

	return device_valid
}

func CountDown() bool {
	var T_SEC_DELAY = 6
	if util.DEBUG != false {
		T_SEC_DELAY = 0
	}

	Log.Printf("delaying execution for %d seconds, hold tight ...\n", T_SEC_DELAY)
	for i := 1; i <= T_SEC_DELAY; i++ {
		Log.Println("Starting Tkey fuzzer in ... ", i)
		time.Sleep(1 * time.Second)
		if i >= 2 && device_valid != true {
			if CxPipeline.SmokeTest() != true {
				Err.Println("Smoke test failed. 💔 Make sure to properly connect device or emulation")

				if once_asked { /* allow to retry again, once **/
					Err.Println("❌ Retry was not helpful. Will exit the pipeline. Make sure to use correct information.")
					return false
				} else {
					once_asked = true
					if util.AskForConfirmation() {
						if CxPipeline.SmokeTest() != true {
							Err.Println("❌ Retry was not helpful. Will exit the pipeline. Make sure to use correct information.")
						}
					}
				}

				return false
			} else {
				Log.Println("✅ Superb! Seems like the Pipeline Smoke Test passed: SUCCESS")
				//return true
			}
		}
	}

	Log.Println("Fuzzing started")
	Log.Printf("%s", s.Repeat(" ❤️‍🔥 ", 10))
	return true
}
