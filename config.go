package log

import "flag"

type LogConfig struct {
	IsFileModel bool   `json:"isFile" ini:"isFile"`
	FileDir     string `json:"fileDir" ini:"fileDir"`
	MaxSize     int    `json:"maxSize" ini:"maxSize"`
	MaxHour     int    `json:"maxHour" ini:"maxHour"`
	Prefix      string `json:"prefix" ini:"prefix"`
	FileName    string `json:"fileName" ini:"fileName"`
}

func init() {
	logFile := LogConfig{}
	flag.BoolVar(&logFile.IsFileModel, "isFile", false, "Is or not file log")
	flag.StringVar(&logFile.FileDir, "fileDir", "", "Defaults to the current folder")
	flag.IntVar(&logFile.MaxSize, "maxSize", 0, "Maximum file size")
	flag.IntVar(&logFile.MaxHour, "maxHour", 0, "Maximum retention hours")

	flag.Parse()
}
