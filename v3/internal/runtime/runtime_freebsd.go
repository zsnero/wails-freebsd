//go:build freebsd

package runtime

// FreeBSD runs the same WebKit2GTK build as Linux, so the message-handler
// shim is identical. On webkit2gtk, `messageHandlers.external.postMessage` only
// works when `this` is bound to the handler object. Assigning the bare function
// reference (as we did historically) silently swallows messages when called as
// `window._wails.invoke(msg)` — the page's invoke loses the receiver. Wrapping
// it like darwin does so callers can invoke without thinking about receiver
// binding.
var invoke = "window._wails.invoke=function(msg){window.webkit.messageHandlers.external.postMessage(msg);};"
