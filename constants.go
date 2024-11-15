package angus

const (
	BufferSize = 262144 // 256K
	//BufferSize = 1048576 * 5 // 5M

	MSG           = 0x1 // send message to console log
	RUNJS         = 0x2 // run javascript
	APPLYCSS      = 0x3 // apply css
	APPLYHTML     = 0x4 // insert html at a specific element
	LOADJS        = 0x5 // load javascript file from url
	LOADCSS       = 0x6 // load css file from url
	LOADHTML      = 0x7 // load html file from url
	REGISTEREVENT = 0x8 // register event
)
