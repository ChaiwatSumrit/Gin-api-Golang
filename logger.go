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
	Level       		string   	`json:"level"`
	Timestamp			string 		`json:"timestamp"`
	Actor				string		`json:"actor"`
	Component			string      `json:"component"`
	FunctionName    	string 		`json:"function"`
	Message  			string 		`json:"message"`
	RequestMethod 		string 		`json:"method"`
	CodeStatus     		string    	`json:"status"`
	UUID				string     	`json:"uuid"`
}

//	Logger("level", "component",  "requestMethod", "functionName", "message", "codeStatus")

func Logger(level string, actor string,component string, requestMethod string, functionName string, message string, codeStatus string, channel chan<- string ) {

	// timestamp := time.Now()
	// timestamp.Format(time.RFC3339)
	timeZone , _ := time.LoadLocation("Asia/Bangkok")
	timestamp := time.Now().In(timeZone)

	loggerAsJson := LoggerModel{
		Timestamp		:timestamp.Format(time.RFC3339),
		Actor			:strings.ToLower(actor),
		Component		:component,
		RequestMethod 	:strings.ToUpper(requestMethod),
		FunctionName 	:functionName,	
		Message  		:message,	
		CodeStatus     	:codeStatus,
		UUID			:UUIR_LOGS,
	}
	
	level = strings.ToUpper(level)
	if level=="" || (level!="INFO" && level!="DEBUG" && level!="ERROR" && level!="WANNING" && level!="FATAL" ) {
		fmt.Printf("Logger require level='INFO','DEBUG','ERROR' or 'WANNING'")
		return
	}

	switch level {

	case "INFO":
		loggerAsJson.Level = "INFO"
		fmt.Printf(`%s %s |%s| "actor":"%s" "component":"%s" "function":"%s" %s %s`+"\n",
		 loggerAsJson.Level, loggerAsJson.Timestamp, message , loggerAsJson.Actor, component, functionName, requestMethod, codeStatus)	
	case "DEBUG":
		loggerAsJson.Level = "DEBUG"
		fmt.Printf(`%s %s |%s| "actor":"%s" "component":"%s" "function":"%s" %s  %s`+"\n",
		 loggerAsJson.Level, loggerAsJson.Timestamp, message , loggerAsJson.Actor, component, functionName, requestMethod, codeStatus)
	case "ERROR":
		loggerAsJson.Level = "ERROR"
		fmt.Printf(`%s %s |%s| "actor":"%s" "component":"%s" "function":"%s" %s  %s`+"\n",
		 loggerAsJson.Level, loggerAsJson.Timestamp, message , loggerAsJson.Actor, component, functionName, requestMethod, codeStatus)
	case "WANNING":
		loggerAsJson.Level = "WANNING"
		fmt.Printf(`%s %s |%s| "actor":"%s" "component":"%s" "function":"%s" %s  %s`+"\n",
		 loggerAsJson.Level, loggerAsJson.Timestamp, message , loggerAsJson.Actor, component, functionName, requestMethod, codeStatus)
	case "FATAL":
		loggerAsJson.Level = "FATAL"
		fmt.Printf(`%s %s |%s| "actor":"%s" "component":"%s" "function":"%s" %s  %s`+"\n",
		 loggerAsJson.Level, loggerAsJson.Timestamp, message , loggerAsJson.Actor, component, functionName, requestMethod, codeStatus)
	
	}
	// LoggerDriving(loggerAsJson)

	channel <- level+" OK"

}




// func LoggerDriving(payload LoggerModel) {
//     url := "0.0.0.0:9000"
// 	fmt.Println("URL:>", url)
// 	fmt.Println("payload:>", payload)

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
//     // resp, err := client.Do(req)
//     // if err != nil {
//     //     panic(err)
// 	// }
// 	res, _ := client.Do(req)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// defer res.Body.Close()
// 	fmt.Println("res:>", res)

// 	return 
// }

