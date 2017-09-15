package sdl

// #cgo windows LDFLAGS: -lSDL2
// #cgo linux freebsd darwin pkg-config: sdl2
// #include "sdl_wrapper.h"
import "C"

import (
	"runtime"
)

const (
	INIT_TIMER          = 0x00000001
	INIT_AUDIO          = 0x00000010
	INIT_VIDEO          = 0x00000020
	INIT_JOYSTICK       = 0x00000200
	INIT_HAPTIC         = 0x00001000
	INIT_GAMECONTROLLER = 0x00002000
	INIT_NOPARACHUTE    = 0x00100000
	INIT_EVERYTHING     = INIT_TIMER | INIT_AUDIO | INIT_VIDEO | INIT_JOYSTICK |
		INIT_HAPTIC | INIT_GAMECONTROLLER
)

const (
	RELEASED = 0
	PRESSED  = 1
)

var (
	// Queue of functions that are thread-sensitive
	callQueue chan func() // initialized only if Main(..) is called
	// Channel to signal that a thread-sensitive function has finished execution
	callDone  chan bool   // initialized only if Main(..) is called
)

func init() {
	// Make sure the main goroutine is bound to the main thread.
	runtime.LockOSThread()
}

// Main entry point. Run this function at the beginning of main(), and pass your
// own main body to it as a function. E.g.:
//
// 	func main() {
// 		sdl.Main(func() {
// 			// Your code here....
// 			// [....]
//
// 			// Calls to SDL can be made by any goroutine, but always guarded by sdl.Do()
// 			sdl.Do(func() {
// 				sdl.Init(0)
// 			})
// 		})
// 	}
//
// Avoid calling functions like os.Exit(..) within your passed-in function since
// they don't respect deferred calls. Instead, do this:
//
// 	func main() {
// 		var exitcode int
// 		sdl.Main(func() {
// 			exitcode = run()) // assuming run has signature func() int
// 		})
// 		os.Exit(exitcode)
// 	}
func Main(main func()) {
	callQueue, callDone = make(chan func()), make(chan bool, 1)

	go func() {
		defer close(callQueue)
		main()
		// fmt.Println("END") // to check if os.Exit(..) is called by main() above
	}()

	for f := range callQueue { f() }
}

// Do the specified function in the main thread.
// For this function to work, you must have correctly used sdl.Main(..) in your
// main() function. Calling this function before/without sdl.Main(..) will cause
// your code to block indefinitely.
func Do(f func()) {
	callQueue <- func() {
		f()
		callDone <- true
	}
	<-callDone
}

// Init (https://wiki.libsdl.org/SDL_Init)
func Init(flags uint32) (err error) {
	if C.SDL_Init(C.Uint32(flags)) != 0 {
		err = GetError()
	}
	return
}

// Quit (https://wiki.libsdl.org/SDL_Quit)
func Quit() {
	C.SDL_Quit()

	eventFilterCache = nil
	for k := range eventWatches {
		delete(eventWatches, k)
	}
}

// InitSubSystem (https://wiki.libsdl.org/SDL_InitSubSystem)
func InitSubSystem(flags uint32) (err error) {
	if C.SDL_InitSubSystem(C.Uint32(flags)) != 0 {
		err = GetError()
	}
	return
}

// QuitSubSystem (https://wiki.libsdl.org/SDL_QuitSubSystem)
func QuitSubSystem(flags uint32) {
	C.SDL_QuitSubSystem(C.Uint32(flags))
}

// WasInit (https://wiki.libsdl.org/SDL_WasInit)
func WasInit(flags uint32) uint32 {
	return uint32(C.SDL_WasInit(C.Uint32(flags)))
}

// GetPlatform (https://wiki.libsdl.org/SDL_GetPlatform)
func GetPlatform() string {
	return string(C.GoString(C.SDL_GetPlatform()))
}
