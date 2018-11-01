package main

import (
	"github.com/tsingson/fastweb/fasthttputils"
	"os"
)

// get file info for VOD only
func getFileInfo(filename string) (md5checksum string, filesize int64, err error) {

	md5checksum, err = fasthttputils.Md5CheckSum(filename)
	if err != nil {
		return md5checksum, filesize, err
	}
	filesize, err = getFileSize(filename)
	if err != nil {
		return md5checksum, filesize, err
	}
	return md5checksum, filesize, nil
}

// get file size from os
func getFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err

	}
	fileSize := fileInfo.Size() //获取size
	return fileSize, nil
}
