//go:build darwin

package main

/*
#cgo CFLAGS: -fblocks
#cgo LDFLAGS: -framework AppKit -framework WebKit -framework PDFKit -framework Foundation
#include <stdlib.h>

int qsl_print_webview(double width_mm, double height_mm, int landscape, int show_print_panel, const char *printer_name);
*/
import "C"

import "unsafe"

func nativePrintEnvelope(width, height float64, landscape, showPrintPanel bool, printerName string) PrintResult {
	name := C.CString(printerName)
	defer C.free(unsafe.Pointer(name))

	status := C.qsl_print_webview(
		C.double(width),
		C.double(height),
		boolToCInt(landscape),
		boolToCInt(showPrintPanel),
		name,
	)
	switch int(status) {
	case 1:
		return PrintResult{Success: true, Handled: true}
	case 2:
		return PrintResult{Reason: "用户取消打印", Handled: true}
	default:
		return PrintResult{Reason: "找不到 Wails 打印窗口或打印视图", Handled: false}
	}
}

func boolToCInt(value bool) C.int {
	if value {
		return 1
	}
	return 0
}
