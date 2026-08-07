package main

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"qsl-mail/backend/store"
)

type Printer struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
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

func (a *App) GetPrinters() []Printer {
	command, args := "lpstat", []string{"-a"}
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
		name := strings.Fields(line)[0]
		if runtime.GOOS != "windows" {
			name = strings.TrimSuffix(name, ":")
		}
		printers = append(printers, Printer{Name: name, DisplayName: name})
	}
	return printers
}
