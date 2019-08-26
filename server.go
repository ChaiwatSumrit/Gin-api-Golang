package main

import (

	//gin framwork
	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib"

	//http
	"net/http"

	//mysql lib
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql"

	//mongo
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	
	// "flag"
	"os"
	"io"
	"bytes"
	"context"
	"fmt"
	"log"

	//uuid
	"os/exec"
)

/* 	
	########################################################################################################
	################################################## MAIN ################################################
	########################################################################################################
*/	



var UUIR_LOGS string
var ACTOR = "robot_test"
func main(){

	UUIR_LOGS, err := exec.Command("uuidgen").Output()
    if err != nil {
        log.Fatal(err)
	}
	fmt.Printf("ACTOR : %s", ACTOR)
	fmt.Printf("UUIR_LOGS : %s", UUIR_LOGS)
	
	// Example skip path request.
	app := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	app.Run(":8080")
}

/* 	
	########################################################################################################
	################################################# MYSQL ################################################
	########################################################################################################
*/	

//MODEL
type Customer struct {
	Id        uint   `json:"_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

func (h *CustomerHandler) InitializeMYSQL() {
// 	// "user:password@/dbname?charset=utf8&parseTime=True&loc=Local
// 	// db, err := gorm.Open("mysql", "root:best1459900574821@dbname?charset=utf8&parseTime=True&loc=Local")
// 	db, err := gorm.Open("mysql", "root:best1459900574821@tcp(127.0.0.1:3306)/charset")	
// 	if err != nil {
// 		go Logger("ERROR", "sample_server", "", "InitializeMYSQL", "Connect MYSQL Database Fail Error :"+err.Error(), "")
// 	}
// 	go Logger("INFO", "sample_server", "", "InitializeMYSQL", "Connect Database Success root : best1459900574821@tcp(127.0.0.1:3306)/charset", "")

// 	db.AutoMigrate(&Customer{})
// 	h.DB = db
}


type CustomerHandler struct {
	Collection *mongo.Collection
	Channel chan string
}


/* 	
	########################################################################################################
	############################################## MONGO DB ################################################
	########################################################################################################
*/	

//mongodb+srv://development:<password>@clustermaster-zvis2.mongodb.net/test?retryWrites=true&w=majority
//ZGFMzUvDJ745GFDq
var username = "development"
var host = "clustermaster-zvis2.mongodb.net/test?retryWrites=true&w=majority"  // of the form foo.mongodb.net

func (h *CustomerHandler) InitializeMongoDB() {

	ctx := context.TODO()
	password	:= "ZGFMzUvDJ745GFDq"

    mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s", username, password, host)
	go Logger("DEBUG","Admin" ,"sample_server", "", "InitializeMongoDB", "connection string is:"+mongoURI, "", h.Channel)

    // Set client options and connect
    clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
		go Logger("FATAL","Admin" ,"sample_server", "", "InitializeMongoDB", err.Error(), "", h.Channel)
        os.Exit(1)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
		go Logger("FATAL","Admin" ,"sample_server", "", "InitializeMongoDB", err.Error(), "", h.Channel)
        os.Exit(1)
	}
	
	collection := client.Database("logistics").Collection("customer")

	h.Collection = collection
	go Logger("INFO","Admin" ,"sample_server", "", "InitializeMongoDB", "Connected to MongoDB!", "", h.Channel)

}

/* 	
	########################################################################################################
	############################################## MIDELWARE ###############################################
	########################################################################################################
*/
type MyReadCloser struct {
	rc io.ReadCloser
	w  io.Writer
}

func (rc *MyReadCloser) Read(p []byte) (n int, err error) {
	n, err = rc.rc.Read(p)
	if n > 0 {
		if n, err := rc.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}

func (rc *MyReadCloser) Close() error {
	return rc.rc.Close()
}

func (h *CustomerHandler)LoggerPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG",ACTOR ,"sample_server", "POST", "LoggerPayload", "payload="+buf.String(), "", h.Channel)

		}else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			go Logger("DEBUG",ACTOR , "sample_server", "PUT", "LoggerPayload", "payload="+buf.String(), "", h.Channel)

		}else {
			c.Next()
		}
	}
}

func (h *CustomerHandler) InitializeChannel() {
	h.Channel = make(chan string)
	return
}
/* 	
	########################################################################################################
	############################################## GIN FRANWORK ############################################
	########################################################################################################
*/

func setupRouter() *gin.Engine {

	//log fomat json

    // debug := flag.Bool("debug", true, "sets log level to debug")

	// flag.Parse()
	
	//เพื่อสร้าง Engine instance ของ Gin 
	//มี middleware Logger และ Recovery ติดตั้งมาให้
	app := gin.Default() 
	//เหมือน gin.Default() ; Full
	// app := gin.New()

	//middleware

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.


	//mysql
	system := CustomerHandler{}
	// system.InitializeMYSQL()

	system.InitializeMongoDB()
	system.InitializeChannel()

	app.Use(system.LoggerPayload())

	// result := make(chan string)
	go Logger("INFO",ACTOR , "sample_server", "", "setupRouter", "Start API Server localhost:8080", "", system.Channel)

	//app router
	app.GET("/customers", system.GetAllCustomer)
	app.GET("/customers/:id", system.GetCustomer)
	app.POST("/customers", system.SaveCustomer)
	app.PUT("/customers/:id", system.UpdateCustomer)
	app.DELETE("/customers/:id", system.DeleteCustomer)

	// app.Use(logger.Setgo Logger() )

	return app
}

/* 	
	########################################################################################################
	######################################### ROUTER&CONTROLLER ############################################
	########################################################################################################
*/

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {

	go Logger("INFO",ACTOR , "sample_server", "GET", "GetAllCustomer", "Request Function", "", h.Channel)
	go Logger("DEBUG",ACTOR , "sample_server", "GET", "GetAllCustomer", "path="+c.Request.RequestURI, "", h.Channel)

	customers := []*Customer{}

    cur, err := h.Collection.Find(context.TODO(), bson.M{})
    if err != nil {
        log.Fatal("Error on Finding all the documents", err)
    }
    for cur.Next(context.TODO()) {
        var customer Customer
        err = cur.Decode(&customer)
        if err != nil {
            log.Fatal("Error on Decoding the document", err)
        }
        customers = append(customers, &customer)
    }

	go Logger("INFO",ACTOR , "sample_server", "GET", "GetAllCustomer", "Request Success", "200", h.Channel)
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	go Logger("DEBUG",ACTOR , "sample_server", "GET", "GetCustomer", "param="+c.Param("id"), "", h.Channel)

	go Logger("INFO",ACTOR , "sample_server", "GET", "GetCustomer", "Request Function", "", h.Channel)
	go Logger("DEBUG",ACTOR , "sample_server", "GET", "GetCustomer", "path="+c.Request.RequestURI, "", h.Channel)


	// id := c.Param("id")
	customer := Customer{}

	// if err := h.DB.Find(&customer, id).Error; err != nil {
	// 	Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
	// 	go Logger("ERROR", "sample_server", "GET", "GetCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")

	// 	c.JSON(http.StatusNotFound,Msg)
	// 	return
	// }

	go Logger("INFO",ACTOR , "sample_server", "GET", "GetCustomer", "Request Success", "200", h.Channel)
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	go Logger("INFO",ACTOR , "sample_server", "POST", "SaveCustomer", "Request Function", "", h.Channel)

	go Logger("DEBUG",ACTOR , "sample_server", "POST", "SaveCustomer", "path="+c.Request.RequestURI, "", h.Channel)
	// fmt.Println(<-h.Channel)

	customer := Customer{}


	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		go Logger("ERROR",ACTOR , "sample_server", "POST", "SaveCustomer", "BadRequest "+err.Error(), "400", h.Channel)

		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	_, err := h.Collection.InsertOne(context.TODO(), customer)
	if err != nil {
		Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		go Logger("ERROR",ACTOR , "sample_server", "POST", "SaveCustomer", "Insert Database Fail Error: "+err.Error(), "400", h.Channel)

		c.JSON(http.StatusBadRequest,Msg)
	}

	go Logger("INFO",ACTOR , "sample_server", "POST", "SaveCustomer", "Request Success:", "200", h.Channel)

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	go Logger("DEBUG",ACTOR , "sample_server", "PUT", "UpdateCustomer", "param="+c.Param("id"), "", h.Channel)

	go Logger("INFO",ACTOR , "sample_server", "PUT", "UpdateCustomer", "Request Function", "", h.Channel)
	
	go Logger("DEBUG",ACTOR , "sample_server", "PUT", "UpdateCustomer", "path="+c.Request.RequestURI, "", h.Channel)

	// id := c.Param("id")
	customer := Customer{}

	// if err := h.DB.Find(&customer, id).Error; err != nil {
	// 	Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
	// 	go Logger("ERROR", "sample_server", "PUT", "UpdateCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")
	// 	c.JSON(http.StatusNotFound,Msg)
	// 	return
	// }

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		go Logger("ERROR",ACTOR , "sample_server", "PUT", "UpdateCustomer", "Insert Database Fail Error: "+err.Error(), "400", h.Channel)

		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	// if err := h.DB.Save(&customer).Error; err != nil {
	// 	Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
	// 	go Logger("ERROR", "sample_server", "PUT", "UpdateCustomer", "Insert Database Fail Error: "+err.Error(), "400")

	// 	c.JSON(http.StatusBadRequest,Msg)
	// 	return
	// }

	go Logger("INFO",ACTOR , "sample_server", "PUT", "UpdateCustomer", "Request Success", "200", h.Channel)
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	go Logger("DEBUG",ACTOR , "sample_server", "DELETE", "DeleteCustomer", "param="+c.Param("id"), "", h.Channel)

	go Logger("INFO",ACTOR , "sample_server", "DELETE", "DeleteCustomer", "Request Function", "", h.Channel)
	go Logger("DEBUG",ACTOR , "sample_server", "DELETE", "DeleteCustomer", "path="+c.Request.RequestURI, "", h.Channel)

	// id := c.Param("id")
	// customer := Customer{}

	// if err := h.DB.Find(&customer, id).Error; err != nil {
	// 	Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
	// 	go Logger("ERROR", "sample_server", "DELETE", "DeleteCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")

	// 	c.JSON(http.StatusNotFound,Msg)
	// 	return
	// }

	// if err := h.DB.Delete(&customer).Error; err != nil {
	// 	Msg := "Msg='Delete Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
	// 	go Logger("ERROR", "sample_server", "DELETE", "DeleteCustomer", "Delete Database Fail Error: "+err.Error(), "400")

	// 	c.JSON(http.StatusNotFound,Msg)
	// 	return
	// }

	go Logger("INFO",ACTOR , "sample_server", "DELETE", "DeleteCustomer", "Request Success", "200", h.Channel)
	c.Status(http.StatusNoContent)
}




















