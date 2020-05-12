package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type FileDb struct {
	File_Name string `json:"File_Name"`
}

type MongoDB struct {
	Mongo_Ip        string `json:"Mongo_Ip"`
	Mongo_Port      string `json:"Mongo_Port"`
	DB_Name         string `json:"DB_Name"`
	Collection_Name string `json:"Collection_Name"`
}

type Config struct {
	RecordType string  `json:"RecordType"`
	FileT      FileDb  `json:"FileDb"`
	MongoT     MongoDB `json:"MongoDb"`
}

func main() {
	// 打开json文件
	jsonFile, err := os.Open("conf.json")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()
	var new Config
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &new)

	fmt.Println(new.RecordType)

	new1 := new.FileT
	fmt.Println(reflect.TypeOf(new1))

}
