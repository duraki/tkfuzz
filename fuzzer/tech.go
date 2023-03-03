/**
// impl. Pinout Fuzzing
// impl. Various baud rates
// impl. Baud buffer size
// impl. Querying PINS for data retrival
// impl. serial read + serial cmd + incrementation techq.
// impl. fuzz functions:
//		* byte/data reverse
//		* test cases
//		* unicode/utf8 bytefuzz
// 		* write custom data encoder (reversed bits, bits encoding etc.)
**/

package fuzzer

import "math/rand"

// Rx/Tx Commands
var CMD_BYTEMAP [10]byte /* tk fw cmds */
var CMD_MAXLEN [5]int    /* max cmd buff size */
var CMD_EXTRA [3]byte
var CMD_EXTRA_RSP [3]byte

// Rx/Tx Frames
var FRM_BYTEMAP [30]byte /* fw frame parsing bytemaps */
var FRM_MAXLEN [4]int    /* fw frame maxlen */
var FWRSP_LEN [5]int     /* fw rsp len */

const ( /** enum endpoints */
	DST_HW_IFPGA int = 0
	DST_HW_AFPGA     = 1
	DST_FW           = 2
	DST_SW           = 3
)

func RAND(ary [10]byte) int {
	sz := len(ary) - rand.Int()
	return sz

	// rand.Int()
	//return sz[rand.Uint64()]
}

func InitByteMaps() {
	/** transmissions */
	InitCommands()
	InitFrames()
	InitExtraCommands()

	/** expectations */
	DefineFirmwareRSP()
}

func InitCommands() {
	cmd_bytes := [10]byte{
		0x01 /* CMD_GET_NAME_VERSION */, 0x02, /* RSP */
		0x03 /* CMD_LOAD_APP */, 0x04, /* RSP */
		0x05 /* CMD_LOAD_APP_DATA */, 0x06, /* RSP */
		/*						   */ 0x07,           /*RSP APP_DATA_READY */
		0x08 /* CMD_GET_UID */, 0x09, /*RSP */
	}
	CMD_BYTEMAP = cmd_bytes

	cmd_maxlen := [5]int{
		0, 1, 4, 32, 128,
	}
	CMD_MAXLEN = cmd_maxlen

	Log.Printf("cmd_bytes=%v cmd_maxlen=%d", cmd_bytes, cmd_maxlen)
}

func InitExtraCommands() {
	cmd_extra := [3]byte{
		0x00, /* CMD sent to the HW, with a single byte of data. Usable in TRNG*/
		0x13, /* CMD sent to the FW, with 128 bytes of data. Load app binary or other data in memory.*/
		0x1a, /* CMD sent to running APP, with 32 bytes of data. Used in key derivation, signature..*/
	}
	CMD_EXTRA = cmd_extra

	cmd_extra_rsp := [3]byte{
		0x01, /* RSP CMD sent to the HW which responds with 4 bytes of data. Used in VERSION retrieval.*/
		0x14, /* RSP CMD to the FW indicating unsuccessful command, which responds with 1 byte of data.*/
		0x1b, /* RSP CMD to the running APP, indicating successful command. The RSP contains 128 bytes of data.*
		/*****/ /* 	   	 used during EdDSA Ed25519 signature extraction */
	}
	CMD_EXTRA_RSP = cmd_extra_rsp
}

func InitFrames() {
	/** the fw's parseframe() function takes 2 args:
	 * 1.  uint8_t in = readbyte()
	 * OR  uint8_t in = readbyte_ledflash(...)
	 *
	 * 2.  reference to &HDR which contains
	 * 	struct frame_header hdr;
	 *
	 * which seems to act as a frame_header buffer malloc.
	 * used in both directions
	 */
	frm_parse := [30]byte{
		/* in & 0x80 */ 0x80, /* BAD VERSION */
		/* in & 0x4  */ 0x4, /* EXPECTS VAL 0 */
		/* (in & 0x60) >> 5 */ 0x60, 0x60 >> 5, /* HDR ID */
		/* (in & 0x18) >> 3 */ 0x18, 0x18 >> 3, /* HDR Endpoint - might be abused to jump in FW/APP/FPGA mode */
	}
	FRM_BYTEMAP = frm_parse

	/** extracted in parseframe() function, by AND operand */
	/** len = (in & 0x3) **/
	frm_len := [4]int{
		1, 4, 32, 128,
	}
	FRM_MAXLEN = frm_len
}

func DefineFirmwareRSP() {
	fwrsp_bytes := [5]int{
		/* FW_RSP_NAME_VERSION */ 32,
		/* FW_RSP_LOAD_APP */ 4,
		/* FW_RSP_LOAD_APP_DATA */ 4,
		/* FW_RSP_LOAD_APP_DATA_READY */ 128,
		/* FW_RSP_GET_UDI */ 32,
	}
	FWRSP_LEN = fwrsp_bytes
}

type FrameHeader struct {
	id       uint8
	endpoint int
	cmdlen   int
}

func fwfuzz(hdr FrameHeader, rspcode int /* fwcmd */, buf uint8) {
	/** the fwreply() function works like this **/
	/**
		switch (rspcode): 			// depending on rspcode provided

		case FW_RSP_CMD_X:
			len, nbytes = 128
			break;
		case FW_RSP....
			...

			// Frame Protocol Header
		writebyte(
			genhdr(
				hdr.id,
				hdr.endpoint,
				0x0,
				len
			)
		)

			// Firmware Protocol Header
		writebyte(rspcode)

			// Do same for remaining bytes
		nbytes = nbytes - 1
			// The remaining bytes and buffer is sent in one chunk
		write(
			buf,
			nbytes
		)
	**/

}
