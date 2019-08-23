package main

import (

	//gin framwork
	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib"

	//http
	"net/http"

	//mysql lib
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	
	// "flag"

	"io"
	"bytes"
)

/* 	
	########################################################################################################
	################################################## MAIN ################################################
	########################################################################################################
*/	




func main(){
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
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

func (h *CustomerHandler) Initialize() {
	// "user:password@/dbname?charset=utf8&parseTime=True&loc=Local
	// db, err := gorm.Open("mysql", "root:best1459900574821@dbname?charset=utf8&parseTime=True&loc=Local")
	db, err := gorm.Open("mysql", "root:best1459900574821@tcp(127.0.0.1:3306)/charset")	
	if err != nil {
		Logger("ERROR", "sample_server", "", "Initialize", "Connect MYSQL Database Fail Error :"+err.Error(), "")
	}
	Logger("INFO", "sample_server", "", "Initialize", "Connect Database Success root : best1459900574821@tcp(127.0.0.1:3306)/charset", "")

	db.AutoMigrate(&Customer{})
	h.DB = db
}
	

type CustomerHandler struct {
	DB *gorm.DB
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

func LoggerPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			Logger("DEBUG", "sample_server", "POST", "LoggerPayload", "payload="+buf.String(), "")

		}else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()
			Logger("DEBUG", "sample_server", "PUT", "LoggerPayload", "payload="+buf.String(), "")

		}else {
			c.Next()
		}
	}
}

/* 	
	########################################################################################################
	############################################## GIN FRANWORK ############################################
	########################################################################################################
*/

func setupRouter() *gin.Engine {

	//log fomat json
	Logger("INFO", "sample_server", "", "setupRouter", "Start API Server localhost:8080", "")

    // debug := flag.Bool("debug", true, "sets log level to debug")

	// flag.Parse()
	
	//เพื่อสร้าง Engine instance ของ Gin 
	//มี middleware Logger และ Recovery ติดตั้งมาให้
	app := gin.Default() 
	//เหมือน gin.Default() ; Full
	// app := gin.New()

	//middleware
	app.Use(LoggerPayload())

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.


	//mysql
	system := CustomerHandler{}
	system.Initialize()

	//app router
	app.GET("/customers", system.GetAllCustomer)
	app.GET("/customers/:id", system.GetCustomer)
	app.POST("/customers", system.SaveCustomer)
	app.PUT("/customers/:id", system.UpdateCustomer)
	app.DELETE("/customers/:id", system.DeleteCustomer)

	// app.Use(logger.SetLogger() )

	return app
}

/* 	
	########################################################################################################
	######################################### ROUTER&CONTROLLER ############################################
	########################################################################################################
*/

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {

	Logger("INFO", "sample_server", "GET", "GetAllCustomer", "Request Function", "")
	Logger("DEBUG", "sample_server", "GET", "GetAllCustomer", "path="+c.Request.RequestURI, "")

	customers := []Customer{}

	h.DB.Find(&customers)

	Logger("INFO", "sample_server", "GET", "GetAllCustomer", "Request Success", "200")
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	Logger("DEBUG", "sample_server", "GET", "GetCustomer", "param="+c.Param("id"), "")

	Logger("INFO", "sample_server", "GET", "GetCustomer", "Request Function", "")
	Logger("DEBUG", "sample_server", "GET", "GetCustomer", "path="+c.Request.RequestURI, "")


	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "GET", "GetCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")

		c.JSON(http.StatusNotFound,Msg)
		return
	}

	Logger("INFO", "sample_server", "GET", "GetCustomer", "Request Success", "200")
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	Logger("INFO", "sample_server", "POST", "SaveCustomer", "Request Function", "")
	Logger("DEBUG", "sample_server", "POST", "SaveCustomer", "path="+c.Request.RequestURI, "")

	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "POST", "SaveCustomer", "BadRequest "+err.Error(), "400")

		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "POST", "SaveCustomer", "Insert Database Fail Error: "+err.Error(), "400")

		c.JSON(http.StatusBadRequest,Msg)
	}

	Logger("INFO", "sample_server", "POST", "SaveCustomer", "Request Success", "200")

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	Logger("DEBUG", "sample_server", "PUT", "UpdateCustomer", "param="+c.Param("id"), "")

	Logger("INFO", "sample_server", "PUT", "UpdateCustomer", "Request Function", "")
	Logger("DEBUG", "sample_server", "PUT", "UpdateCustomer", "path="+c.Request.RequestURI, "")

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "PUT", "UpdateCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "PUT", "UpdateCustomer", "Insert Database Fail Error: "+err.Error(), "400")

		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "PUT", "UpdateCustomer", "Insert Database Fail Error: "+err.Error(), "400")

		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	Logger("INFO", "sample_server", "PUT", "UpdateCustomer", "Request Success", "200")
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	Logger("DEBUG", "sample_server", "DELETE", "DeleteCustomer", "param="+c.Param("id"), "")

	Logger("INFO", "sample_server", "DELETE", "DeleteCustomer", "Request Function", "")
	Logger("DEBUG", "sample_server", "DELETE", "DeleteCustomer", "path="+c.Request.RequestURI, "")

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "DELETE", "DeleteCustomer", "DB Not Found id:"+id+" on Database Error: "+err.Error(), "404")

		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		Msg := "Msg='Delete Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		Logger("ERROR", "sample_server", "DELETE", "DeleteCustomer", "Delete Database Fail Error: "+err.Error(), "400")

		c.JSON(http.StatusNotFound,Msg)
		return
	}

	Logger("INFO", "sample_server", "DELETE", "DeleteCustomer", "Request Success", "200")
	c.Status(http.StatusNoContent)
}




















