//go:build !darwin

package main

func nativePrintEnvelope(_, _ float64, _, _ bool, _ string) PrintResult {
	return PrintResult{Reason: "当前平台没有可用的 Wails 原生打印接口", Handled: false}
}
