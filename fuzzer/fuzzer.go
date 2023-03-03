package fuzzer

import (
	"fmt"
	"log"
	"tkfuzz/util"
)

var Log = util.Log

var CxFuzzcfg Fuzzcfg          /* global context container */
var CxFuzzer Fuzzer            /* ... */
var CxPipeline FuzzPipeline    /* pipeline */
var CxProgress FuzzingProgress /* fuzz progress */

const (
	FUZZ_THREADS int = 10
)

/** fuzz modes */
const (
	F_FW  int = 0 // firmware fuzz mode
	F_SYS     = 1 // scope fuzzing
)

type Fuzzcfg struct {
	iteration int
	fzmode    int // Fuzzfw | Fuzzsys
	ttypath   string
}

type Fuzzer struct {
	fzcfg   Fuzzcfg
	threads int
}

func (cfg Fuzzcfg) Describe() string {
	Log.Println("")
	Log.Println("## Fuzzer Configuration")
	Log.Println("struct data:")
	Log.Printf(" Iter(n) => %v", cfg.iteration)
	fzmd := util.Ternary(cfg.fzmode != 1, "FW", "SYS")
	Log.Printf(" Fuzz(Mode) => %s", fzmd)

	return fmt.Sprintf("%d %d %s\n", cfg.iteration, cfg.fzmode, cfg.ttypath)
}

func (f Fuzzer) Describe() string {
	Log.Println("## Fuzzer Container")
	log.Println("struct data->")
	log.Println(" CFG -> ", f.fzcfg)
	log.Printf(" Threads -> (%v)", f.threads)

	return fmt.Sprintf("%v %d\n", f.fzcfg, f.threads)
}

func InitFuzzCfg(itx int, m int, tty string) Fuzzcfg {
	return Fuzzcfg{
		iteration: itx,
		fzmode:    m,
		ttypath:   tty,
	}
}

func InitFuzzer(cfg Fuzzcfg, threads int) Fuzzer {

	return Fuzzer{
		fzcfg:   cfg,
		threads: FUZZ_THREADS,
	}
}
