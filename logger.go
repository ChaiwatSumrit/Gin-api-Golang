package main

import (
	// "encoding/json"
    // "net/http"
	// "bytes"

	"strings"
	"time"
	"fmt"
)

type LoggerModel struct {
	Timestamp			string 		`json:"timestamp"`
	Level       		string   	`json:"level"`
	Component			string      `json:"component"`
	RequestMethod 		string 		`json:"method"`
	FunctionName    	string 		`json:"function"`
	Message  			string 		`json:"message"`
	CodeStatus     		string    	`json:"status"`
}

//	Logger("level", "component",  "requestMethod", "functionName", "message", "codeStatus")

func Logger(level string,component string, requestMethod string, functionName string, message string, codeStatus string) {

	// timestamp := time.Now()
	// timestamp.Format(time.RFC3339)
	timeZone , _ := time.LoadLocation("Asia/Bangkok")
	timestamp := time.Now().In(timeZone)

	loggerAsJson := LoggerModel{
		Timestamp		:timestamp.Format(time.RFC3339),
		Component		:component,
		RequestMethod 	:strings.ToUpper(requestMethod),
		FunctionName 	:functionName,	
		Message  		:message,	
		CodeStatus     	:codeStatus,	
	}

	level = strings.ToUpper(level)
	if level=="" || (level!="INFO" && level!="DEBUG" && level!="ERROR" && level!="WANNING" ) {
		fmt.Printf("Logger require level='INFO','DEBUG','ERROR' or 'WANNING'")
		return
	}

	switch level {

	case "INFO":
		loggerAsJson.Level = "INFO"
		fmt.Printf(`"Timestamp":"%s" "level":"%s" "component":"%s" "method":"%s" "function:"%s" "message:"%s" "status":"%s"`+"\n",
		loggerAsJson.Timestamp, loggerAsJson.Level, component, requestMethod, functionName, message, codeStatus)	
	case "DEBUG":
		loggerAsJson.Level = "DEBUG"
		fmt.Printf(`"Timestamp":"%s" "level":"%s" "component":"%s" "method":"%s" "function:"%s" "message:"%s" "status":"%s"`+"\n",
		loggerAsJson.Timestamp, loggerAsJson.Level, component, requestMethod, functionName, message, codeStatus)
	case "ERROR":
		loggerAsJson.Level = "ERROR"
		fmt.Printf(`"Timestamp":"%s" "level":"%s" "component":"%s" "method":"%s" "function:"%s" "message:"%s" "status":"%s"`+"\n",
		loggerAsJson.Timestamp, loggerAsJson.Level, component, requestMethod, functionName, message, codeStatus)
	case "WANNING":
		loggerAsJson.Level = "WANNING"
		fmt.Printf(`"Timestamp":"%s" "level":"%s" "component":"%s" "method":"%s" "function:"%s" "message:"%s" "status":"%s"`+"\n",
		loggerAsJson.Timestamp, loggerAsJson.Level, component, requestMethod, functionName, message, codeStatus)
	}
	// LoggerDriving(loggerAsJson)
	return
}

// func LoggerDriving(payload LoggerModel) {
//     url := "localhost:9000"
//     fmt.Println("URL:>", url)


// 	//Json to byteArray
// 	payloadAsBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		fmt.Println("Marshal is error" + err.Error())
// 		return 
// 	}

//     // var jsonAsStr = []byte(payload)
//     req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadAsBytes))
//     req.Header.Set("Content-Type", "application/json")

//     client := &http.Client{}
//     resp, err := client.Do(req)
//     if err != nil {
//         panic(err)
//     }
// 	defer resp.Body.Close()
// 	return
// }