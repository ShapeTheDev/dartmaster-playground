package main

import (
	"time"

	cameraadmin "github.com/One-Hundred-Eighty/Circle/pkg/camera-admin"
)

func main() {
	cameraAdmin := cameraadmin.NewCameraAdmin()

	err := cameraAdmin.Start(1920, 1080)
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	cam1Ch := cameraAdmin.Subscribe(1, "maint")
	time.Sleep(10 * time.Second)
	cameraAdmin.Unsubscribe(1, cam1Ch, "maint")
	time.Sleep(2 * time.Second)
	cameraAdmin.ShutDown()
}
