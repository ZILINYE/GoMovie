package Process

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type StoreType interface {
	CheckRecord([]Movie_info) []Movie_info
}
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

var new []Movie_info

func ReadConf() StoreType {
	// 打开json文件
	jsonFile, err := os.Open("conf.json")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()
	var new Config
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err!=nil{
		fmt.Print(err)
	}
	json.Unmarshal([]byte(byteValue), &new)
	new1 := new.RecordType

	if new1 == "FileDb" {
		return new.FileT
	} else {
		return new.MongoT
	}
}
func (t FileDb) CheckRecord(m []Movie_info) []Movie_info {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	f, _ := os.OpenFile(t.File_Name, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)


	defer f.Close()

	scanner := bufio.NewScanner(f)

	for _, value := range m {
		var i = false
		for scanner.Scan() {
			if strings.TrimSpace(value.Title) == scanner.Text() {
				i = true
				break
			}
		}
		if i == false {
			new = append(new, value)
			f.WriteString(strings.TrimSpace(value.Title))
			f.WriteString("\n")

		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return new

}

func (t MongoDB) CheckRecord(m []Movie_info) []Movie_info {
	// Initialize Mongo DB Connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("Movie").Collection("Download_info")
	// For Loop Check if Movie had been spider before
	for _, value := range m {
		var result_db Movie_info
		filter := bson.D{{"d_url", value.D_url}}
		// Check if data exist
		err = collection.FindOne(context.TODO(), filter).Decode(&result_db)
		// If Data not been spider before Do Next
		if err != nil {
			// Insert Movie into Mongo DB
			_, err := collection.InsertOne(context.TODO(), value)
			if err != nil {
				fmt.Println(err)
			}
			new = append(new, value)

		}

	}
	return new
}
