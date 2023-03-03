# tkfuzz

Fuzzer for Tillitis TKey written in Go.

## Usage

Defined CLI Commands:

```
$ FTTY=<tty> DBG=1 tkfuzz <mode> <tty> 
```

Usage options:

```
$ tkfuzz fuzzfw /dev/tty* 			# starts fw fuzz on /dev/tty*
$ tkfuzz fuzzsys /dev/tty* 			# starts user mode fuzz on /dev/tty*


					# starts fw fuzz on (ENV)FTTY with maximum
					# permutation settings available 
$ FTTY=/dev/tty* MAX=true \
  tkfuzz fuzzfw
```

**Logs** will be stored in:

```
$ tree logs/
		├── err.log
		└── tkfuzz.log
```

### Command-Line Options

**Arguments:**

* `<mode>` is fuzz mode, either `fuzzfw` & `fuzzsys`
  - `fuzzfw` does preliminary smoke test to see if TKey is in Firmware mode. Executes commands and sends firmware related frames.
  - `fuzzsys` does fuzzing in the rest of the ring regions inside TKey, such is Application mode. 
* `tty` refers to the serial device, a FS path pointing to primary or replica (pseudo) terminal driver. See also `FTTY=` below.

**Environment Variables:**

* `FTTY=(STRING)` can be used instead of `<tty>`
* `DBG=(BOOL)` can be used to set state of `tkfuzz`
* `MAX=(BOOL)` can be used to set fuzzer to maximum capacity

## Run or Build

Run via `go run` without build:

```
$ go run main <mode> <tty>
```

Build with `go build` and execute the binary:

```
$ go build && ./tkfuzz
# ...
```

### Running Sample

```
$ DBG=true go run main.go fuzzfw /dev/tty.usbmodem1101
[log] 2023/03/01 11:01:45 log.go:48: tkfuzz is starting . logging initialised.
[err] 2023/03/01 11:01:45 log.go:49: tkfuzz stderr ~~~~                         valid.

Tillitis TK1 Firmware fuzzer
		 ... by halis duraki,
		collab. of durakiconsulting, LLC & Skullkey, AB (2023)

2023/03/01 11:01:45 serial tty is set to:
		/dev/tty.usbmodem1101
2023/03/01 11:01:45 fuzz mode selected:
		fuzzfw [IDx: 0]
2023/03/01 11:01:45 fuzz iteration ctx:
		102400 [100 * 1024]
2023/03/01 11:01:45 fuzzer initialized and configured properly. now dumping Fuzzer container ctx:
2023/03/01 11:01:45 {{102400 0 /dev/tty.usbmodem1101} 10}
2023/03/01 11:01:45 ****************************************************************************************************
2023/03/01 11:01:45 ## Fuzzer Container
2023/03/01 11:01:45 struct data->
2023/03/01 11:01:45  CFG ->  {102400 0 /dev/tty.usbmodem1101}
2023/03/01 11:01:45  Threads -> (10)
****************************************************************************************************
2023/03/01 11:01:45 ## Fuzzer Configuration
2023/03/01 11:01:45 struct data:
2023/03/01 11:01:45  Iter(n) => 102400
2023/03/01 11:01:45  Fuzz(Mode) => FW
****************************************************************************************************
2023/03/01 11:01:45             FuzzPipeline={{13904835822000419528 1694459 0x102fffda0} {13904835822000419528 1694501 0x102fffda0} 0  {{102400 0 /dev/tty.usbmodem1101} 10} false false {0 [] [] { false}  0}}
2023/03/01 11:01:45 ## Fuzz Pipeline
2023/03/01 11:01:45 struct data:
	// redacted ...
2023/03/01 11:01:45  fuzzer? => {{102400 0 /dev/tty.usbmodem1101} 10}
2023/03/01 11:01:45  Running? => false
****************************************************************************************************
Is this correct fuzz pipeline setup (Y/n):
✅ Fuzz pipeline is confirmed, starting soon ...
2023/03/01 11:01:46 delaying execution for 0 seconds, hold tight ...
2023/03/01 11:01:46 Fuzzing started
2023/03/01 11:01:46 starting fuzz pipeline .... #idx unknown
2023/03/01 11:01:46 CxProgress defined containing => {{{13904835823029014352 956549167 0x102fffda0} {13904835822000419528 1694501 0x102fffda0} 0  {{102400 0 /dev/tty.usbmodem1101} 10} true false {0 [] [] { false}  0}} 0x1400000c030 map[]}
2023/03/01 11:01:46 GenerateFuzzMXNormal_NO_CMD_0_BYTES_IN=_40_LEN=_0
2023/03/01 11:01:46 {96 fuzzfw.cmd-payload-_{a0255d9d-d9a8-6288-1dc2-230075627ae2}%!(EXTRA string= datalen= %d, int=120) 120}
	// redacted ...
```