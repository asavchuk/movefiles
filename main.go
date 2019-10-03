package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	egui "github.com/alkresin/external"
)

var theFileExist = errors.New("The file exists.")

var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func errorHandler() {
	if r := recover(); r != nil {
		notify(fmt.Sprintf("%v", r))
		errorLog.Fatal(r)
	}
}

func main() {
	// go run main.go c:/export/ //SUNSEY2/d$/Temp_Почта/
	// for testing: go run main.go c:/temp/ e:/temp1/
	defer errorHandler()                                                    // in case of runtime errors
	t := time.Now()                                                         //currentTime
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 17, 15, 0, 0, time.Local) // time in future
	diff := t1.Sub(t)
	if diff >= 0 {
		time.Sleep(time.Duration(diff.Seconds()) * time.Second)
	}
	err := movefilelist(os.Args[1], os.Args[2])
	if err != nil {
		errorLog.Println(err)
		notify(err.Error())
		return
	}
}

func moveFile(oldpath, newpath string) error {
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	err = syscall.MoveFile(from, to) //windows API
	if err != nil {
		errorLog.Println(err) // i.e. the file exists
		return err
	}
	return nil
}

func fileNameList(filepath string) []string {
	var list []string
	rd, err := ioutil.ReadDir(filepath)
	if err != nil {
		errorLog.Fatal(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			list = append(list, fi.Name())
		}
	}
	return list
}

func movefilelist(oldpath, newpath string) error {
	for _, fi := range fileNameList(oldpath) {
		// log.Println("from", oldpath+fi)
		// log.Println("to", newpath+fi)

		err := moveFile(oldpath+fi, newpath+fi)
		if err == nil {
			continue
		}
		if err.Error() == theFileExist.Error() {
			errorLog.Println(theFileExist)
			renameAndMoveUntilSuccess(fi, oldpath, newpath)
		} else {
			errorLog.Println(err)
			return err
		}
	}
	return nil
}

func renameAndMoveUntilSuccess(fi, oldpath, newpath string) {
	renamed := fi
	for {
		// renaming the file until it will be successfully moved
		renamed = "_" + renamed
		err := moveFile(oldpath+fi, newpath+renamed)
		if err == nil {
			errorLog.Println("moved success", newpath+renamed)
			break // success
		}
		if err.Error() == theFileExist.Error() {
			continue
		}
		if err != nil && err.Error() != theFileExist.Error() { // any other error
			notify(err.Error())
			errorLog.Fatal(err)
		}
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
