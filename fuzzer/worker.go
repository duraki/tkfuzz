package fuzzer

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"
	"tkfuzz/tk1"
	"tkfuzz/util"

	"github.com/spf13/cast"
)

// Internals
type Endpoint byte

const (
	destAFPGA Endpoint = 1
	DestFW    Endpoint = 2
	DestApp   Endpoint = 3
)

type CmdLen byte
type Cmd interface {
	Code() byte
	String() string

	CmdLen() CmdLen
	Endpoint() Endpoint
}

func (fwcmd FuzzFwCmd) CmdLen() CmdLen {
	return fwcmd.cmdLen
}

/** todo: use mode arg opts. */
func (fwcmd FuzzFwCmd) Endpoint() Endpoint {
	ary := []Endpoint{
		destAFPGA,
		DestFW,
		DestApp,
	}

	rng := 1
	/** correct way */
	result := ary[rng]
	//return result

	return Endpoint(result)
}

func (fwcmd FuzzFwCmd) Code() byte {
	return fwcmd.code
}

// type CMD_CTX<Cmd> struct {
// }
// End Of Internals

type Bitflip struct {
	m_id   int    // bitflip mode id
	mode   string // mode name (ie. reverser, encoder, ...)
	before string // ascii+hex
	after  string // ascii+hex
}

type Payload struct {
	ascii      string          /* ascii-basis */
	hex        string          /* hexadecimal data */
	binr       string          /* binary data */
	buff_size  int             /* total size of payload buff */
	buff_cmd   string          /* per fw and tkey registers */
	bitflip_mx map[int]Bitflip /* bit/byte operations */
}

func PreparePayload(nsize int, bitflips bool) []Payload {
	var ary []Payload // nill slices
	var fuzzdata Payload

	for n := 0; n <= nsize; n++ {
		fuzzdata = GenerateFuzzMXNormal(CMD_BYTEMAP[rand.Intn(len(CMD_BYTEMAP))], CMD_BYTEMAP[:], rand.Intn(len(CMD_MAXLEN)))
		Log.Printf("#prepare(%d/%d) payload id=%d data=%v size=%d\n", n, nsize, n, fuzzdata.hex, fuzzdata.buff_size)
		ary = append(ary, fuzzdata)
	}

	Log.Println("Prepared Payloads: ", "%v", ary)
	return ary
}

func (l CmdLen) ByteLen() int {

	if os.Getenv("MAX") == "true" {
		Log.Println("Environment has MAX=true, will return ByteLeb() w. rand(int)")
		return rand.Int()
	} else {
		Log.Println("Hint: Provide MAX=true as an environ variable for more fuzz")
	}

	switch l {
	case 0:
		return 1
	case 1:
		return 4
	case 2:
		return 32
	case 3:
		return 128
	}

	return 0
}

func CustomFrameBuf(cmd FuzzFwCmd, id int) ([]byte, error) {
	if id > 3 { /**/
		Err.Fatal("frame id must be 0..3")
	} else {
		Log.Println("frame id set to", id)
	} /* frame id must be 0..3 */
	if cmd.Endpoint() > 3 { /**/
	} /* endpoint must be 0..3 */
	if cmd.CmdLen() > 3 { /**/
	} /* CmdLen must be 0..3 */
	/**/ /* 		skipped.*/

	fbuf := (1 + CmdLen(cmd.cmdLen.ByteLen()))
	// Make a buffer with frame header + cmdLen payload
	frame := make([]byte, fbuf)
	frame[0] = (byte(id) << 5) | (byte(cmd.Endpoint()) << 3) | byte(cmd.CmdLen())

	// set command code
	frame[1] = cmd.Code()
	fmt.Printf("custom frame: %2x\n", frame)

	//_, _ = TKh.GetNameVersion()

	return frame, nil
}

func UUID() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

func RandomPick(obj []int) interface{} {
	//wr := []byte{}
	rand.Seed(time.Now().UnixNano())
	lucky := obj[rand.Intn(len(obj))]

	fmt.Println("RandomPick()= ", lucky, ">>", util.PrintType(lucky))

	return lucky
}

type FuzzFwCmd struct {
	code   byte
	name   string
	cmdLen CmdLen
}

//	type fwCmd struct {
//		code   byte
//		name   string
//		cmdLen CmdLen
//	}
func CreateFuzzFwCmd() FuzzFwCmd {
	// lensize := []int{0, 1, 4, 32, 128}
	//NewSlice(start, end, step int) []int
	/** fuzz cmdsize **/

	const INC = 256
	sz := []int{}
	def := cast.ToIntSlice(CMD_MAXLEN)
	sz = append(def)

	if os.Getenv("MAX") == "true" {
		max := cast.ToIntSlice(util.NewSlice(8, INC, 2*2)) // util.NewSlice(0, cast.ToIntSlice(INC_CMD_MAXLEN), 1)
		sz = append(max)

		Log.Println("✦ FUZZ = MAXOUT ✅ [via MAX=true][ENV] Using maxim: INC = CMD_MAXLEN() + ", INC)
	}

	fz_cmdsize := RandomPick(cast.ToIntSlice(sz)) //cast.ToIntSlice((byte[int]))))
	fz_name := fmt.Sprintf("fuzzfw.cmd-payload-{%s}_#max(%d)#", UUID(), len(sz))
	fz_payload := RandomPick(cast.ToIntSlice(CMD_BYTEMAP))

	Log.Printf("\n\n\tCreateFuzzFwCmd(\n\n\t\t%s\n\n\t) cmdsize=0x%d cmdid=0x%02x\n\n", fz_name, fz_cmdsize, fz_payload)

	ecmd, _ := cast.ToIntE(fz_payload)
	fzcmd := FuzzFwCmd{byte(ecmd), fz_name, CmdLen(cast.ToInt(fz_cmdsize))}
	tk1.Dump("tk1.Dump", bytes.NewBufferString(cast.ToString(fz_payload)).Bytes())
	Log.Printf("fuzzfw cmd randomised:\n\t %v", fzcmd)

	test_fwcmd()

	return fzcmd
}

func test_fwcmd() {
	id := DestFW                                     // id == endpoint
	fzcmd := FuzzFwCmd{0x01, "cmdGetNameVersion", 1} //CmdLen(CMD_MAXLEN[0])}
	tx, err := CustomFrameBuf(fzcmd, cast.ToInt(id))
	if err != nil {
		Err.Println("err while creating new frame buff: ", err)
	}

	tk1.Dump("test_fwcmd:dump", tx)
}

func GenerateFuzzMXNormal(cmd uint8, bytes_in []byte, buflen int /*cmd uint32, bytes_in string, len int*/) Payload {
	Log.Print("GenerateFuzzMXNormal", "_NO_CMD_", cmd, "_BYTES_IN=_", len(bytes_in)*rand.Intn(124), "_LEN=_", buflen)

	/** Using bytes.Buffer to write buf **/
	/** buf* bytes usage **/
	/**
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%x", "0x01")
		buf.WriteString("frame_str")
		buf.Write([]byte("str"))
		buf.WriteByte(32)
		buf.WriteRune('☉') // Unicode.UTF-8

		// Logging buf*
		// fmt.Println(buf.String())
	**/

	var fwcmd = CreateFuzzFwCmd()

	os.Exit(0)

	return Payload{
		hex: string(fwcmd.code),
	}
}

func (p Payload) Describe() string   { return "" }
func (p Payload) Permutate() Payload { return p }
