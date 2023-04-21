package service

import (
	"fmt"
	"log"
	"os"
)

type FileService struct {
}

func (f *FileService) WriteToFile(link string, hostIp string) error {
	file, err := os.OpenFile("/root/clientt.ovpn", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.WriteString(fmt.Sprintf("client\nproto v2r\nremote %s 443\n%s", hostIp, link))
	if err != nil {
		return err
	}

	return nil
}
