package main

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	egui "github.com/alkresin/external"
)

func main() {
	// go run main.go d:/export/ //SUNSEY2/d$/Temp_Почта/
	t := time.Now()                                                         //currentTime
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 17, 15, 0, 0, time.Local) // time in future
	diff := t1.Sub(t)
	if diff < 0 {
		notify("Today is too late")
		return
	}

	time.Sleep(time.Duration(diff.Seconds()) * time.Second)
	err := movefilelist(os.Args[1], os.Args[2])
	if err != nil {
		notify(err.Error())
		return
	}
}

func notify(message string) {
	if egui.Init("") != 0 {
		return
	}
	pWindow := &egui.Widget{X: 100, Y: 100, W: 300, H: 120, Title: "Error"}
	egui.InitMainWindow(pWindow)
	pWindow.AddWidget(&egui.Widget{Type: "label", X: 20, Y: 20, W: 245, H: 44, Title: message})
	pWindow.Activate()
	egui.Exit()
}

func movefile(oldpath, newpath string) error {
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return syscall.MoveFile(from, to) //windows API
}

func filenamelist(filepath string) []string {
	var list []string
	rd, err := ioutil.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			list = append(list, fi.Name())
		}
	}
	return list
}

func movefilelist(oldpath, newpath string) error {
	for _, fi := range filenamelist(oldpath) {
		err := movefile(oldpath+fi, newpath+fi)
		if err != nil {
			return err
		}
	}
	return nil
}
