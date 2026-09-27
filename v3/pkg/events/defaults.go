package events

import "runtime"

var defaultWindowEventMapping = map[string]map[WindowEventType]WindowEventType{
	"windows": {
		Windows.WindowClosing:      Common.WindowClosing,
		Windows.WindowInactive:     Common.WindowLostFocus,
		Windows.WindowClickActive:  Common.WindowFocus,
		Windows.WindowActive:       Common.WindowFocus,
		Windows.WindowMaximise:     Common.WindowMaximise,
		Windows.WindowMinimise:     Common.WindowMinimise,
		Windows.WindowRestore:      Common.WindowRestore,
		Windows.WindowUnMaximise:   Common.WindowUnMaximise,
		Windows.WindowUnMinimise:   Common.WindowUnMinimise,
		Windows.WindowFullscreen:   Common.WindowFullscreen,
		Windows.WindowUnFullscreen: Common.WindowUnFullscreen,
		Windows.WindowShow:         Common.WindowShow,
		Windows.WindowHide:         Common.WindowHide,
		Windows.WindowDidMove:      Common.WindowDidMove,
		Windows.WindowDidResize:    Common.WindowDidResize,
		Windows.WindowSetFocus:     Common.WindowFocus,
		Windows.WindowKillFocus:    Common.WindowLostFocus,
		Windows.WindowDPIChanged:   Common.WindowDPIChanged,
	},
	"darwin": {
		Mac.WindowDidResignKey:       Common.WindowLostFocus,
		Mac.WindowDidBecomeKey:       Common.WindowFocus,
		Mac.WindowDidMiniaturize:     Common.WindowMinimise,
		Mac.WindowDidDeminiaturize:   Common.WindowUnMinimise,
		Mac.WindowDidEnterFullScreen: Common.WindowFullscreen,
		Mac.WindowDidExitFullScreen:  Common.WindowUnFullscreen,
		Mac.WindowMaximise:           Common.WindowMaximise,
		Mac.WindowUnMaximise:         Common.WindowUnMaximise,
		Mac.WindowDidMove:            Common.WindowDidMove,
		Mac.WindowDidResize:          Common.WindowDidResize,
		Mac.WindowDidZoom:            Common.WindowMaximise,
		Mac.WindowShow:               Common.WindowShow,
		Mac.WindowHide:               Common.WindowHide,
		Mac.WindowZoomIn:             Common.WindowZoomIn,
		Mac.WindowZoomOut:            Common.WindowZoomOut,
		Mac.WindowZoomReset:          Common.WindowZoomReset,
		Mac.WindowShouldClose:        Common.WindowClosing,
	},
	"linux":   linuxWindowEventMapping,
	"freebsd": linuxWindowEventMapping,
}

// linuxWindowEventMapping is shared with FreeBSD, which reports the same GTK
// window events. It has to be named rather than inlined twice, because
// DefaultWindowEventMapping indexes by runtime.GOOS and a missing key returns
// nil: setupEventMapping would then iterate over an empty map and register no
// window event handlers at all, so close, focus, move, resize and load-finished
// would never reach the frontend.
var linuxWindowEventMapping = map[WindowEventType]WindowEventType{
	Linux.WindowDeleteEvent:  Common.WindowClosing,
	Linux.WindowFocusIn:      Common.WindowFocus,
	Linux.WindowFocusOut:     Common.WindowLostFocus,
	Linux.WindowDidMove:      Common.WindowDidMove,
	Linux.WindowDidResize:    Common.WindowDidResize,
	Linux.WindowLoadFinished: Common.WindowShow,
}

func DefaultWindowEventMapping() map[WindowEventType]WindowEventType {
	platform := runtime.GOOS
	return defaultWindowEventMapping[platform]
}
