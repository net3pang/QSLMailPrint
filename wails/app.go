package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"qsl-mail/backend/store"
)

type Printer struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PrintResult struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
	Handled bool   `json:"handled"`
}

type App struct {
	ctx      context.Context
	database *store.Database
}

func NewApp() *App {
	database, err := store.Open("")
	if err != nil {
		panic(err)
	}
	return &App{database: database}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SaveRecord(collection string, record map[string]any) (map[string]any, error) {
	return a.database.SaveRecord(collection, record)
}

func (a *App) SaveTemplate() (string, error) {
	if a.ctx == nil {
		return "", errors.New("应用尚未完成初始化")
	}
	filename, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		Title:           "保存 CSV 导入模板",
		DefaultFilename: "qsl-mail-import-template.csv",
		Filters: []wailsRuntime.FileFilter{{
			DisplayName: "CSV 文件 (*.csv)",
			Pattern:     "*.csv",
		}},
		CanCreateDirectories: true,
	})
	if err != nil || filename == "" {
		return filename, err
	}
	if strings.ToLower(filepath.Ext(filename)) != ".csv" {
		filename += ".csv"
	}
	content := "\xEF\xBB\xBF姓名,呼号,地址,邮编,手机号\nTakahashi Ken,JA7QXG,1-2-3 Chiyoda Tokyo,100000,090-1234-5678\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		return "", err
	}
	return filename, nil
}

func (a *App) DeleteRecord(collection, id string) error {
	return a.database.DeleteRecord(collection, id)
}

func (a *App) PrintEnvelope(options map[string]any) PrintResult {
	width := optionNumber(options, "width", 220)
	height := optionNumber(options, "height", 110)
	landscape := optionBool(options, "landscape", false)
	showPrintPanel := optionBool(options, "showPrintDialog", !optionBool(options, "silent", false))
	printerName, _ := options["deviceName"].(string)
	return nativePrintEnvelope(width, height, landscape, showPrintPanel, printerName)
}

func optionNumber(options map[string]any, key string, fallback float64) float64 {
	value, ok := options[key]
	if !ok {
		return fallback
	}
	switch number := value.(type) {
	case float64:
		return number
	case float32:
		return float64(number)
	case int:
		return float64(number)
	default:
		return fallback
	}
}

func optionBool(options map[string]any, key string, fallback bool) bool {
	value, ok := options[key].(bool)
	if !ok {
		return fallback
	}
	return value
}

func (a *App) GetPrinters() []Printer {
	command, args := "lpstat", []string{"-p"}
	if runtime.GOOS == "windows" {
		command = "powershell"
		args = []string{"-NoProfile", "-Command", "Get-Printer | Select-Object -ExpandProperty Name"}
	}
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil
	}
	printers := make([]Printer, 0)
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name := line
		if runtime.GOOS == "windows" {
			name = strings.TrimSpace(line)
		} else if strings.HasPrefix(line, "printer ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "printer "))
			if index := strings.Index(name, " is "); index >= 0 {
				name = strings.TrimSpace(name[:index])
			}
		} else if index := strings.Index(name, " accepting requests"); index >= 0 {
			name = strings.TrimSpace(name[:index])
		}
		if name == "" {
			continue
		}
		printers = append(printers, Printer{Name: name, DisplayName: name})
	}
	return printers
}
