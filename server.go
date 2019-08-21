package main

import (
	//gin framwork
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib"

	//http
	"net/http"

	//mysql lib
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	//mysql logger
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"flag"

	"os"
	// "regexp"
	"io"
	"time"
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
	log.Info().Str("language", "golang").Msg("Start API Server localhost:8080")

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
	db, err := gorm.Open("mysql", "root:best1459900574821@tcp(127.0.0.1:3306)/charset1")	
	if err != nil {
		log.Error().Str("functionName", "GetCustomer").Msg("Connect Database Fail Error : "+err.Error())
	}

	db.AutoMigrate(&Customer{})
	h.DB = db
}
	

type CustomerHandler struct {
	DB *gorm.DB
}

// var (
// 	rxURL = regexp.MustCompile(`^/regexp\d*`)
// )


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
	// fmt.Println("run here", n, err)
	// log.Info().Str("foo", "bar").Msg(fmt.Println("run here", n, err))

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

func DebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println(c.Request.RequestURI)
		log.Debug().Str("middleware", "DebugLogger").Msg("path="+c.Request.RequestURI)

		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()

			log.Debug().Str("method", "POST").Msg("payload="+buf.String())


		}else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()

			log.Debug().Str("method", "PUT").Msg("payload="+buf.String())

		// }else if c.Request.Method == http.MethodGet {
		// 	var buf bytes.Buffer
		// 	newBody := &MyReadCloser{c.Request.Body, &buf}
		// 	c.Request.Body = newBody
		// 	c.Next()
		// 	log.Debug().Str("method", "GET").Msg("payload="+buf.String())
	
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
	log.Info().Str("language", "golang").Msg("Start API Server localhost:8080")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    debug := flag.Bool("debug", true, "sets log level to debug")

	flag.Parse()
	// Default level for this example is info, unless debug flag is present
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// if gin.IsDebugging() {
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//Add file and line number to log
	// log.Logger = log.With().Caller().Logger()
	// log.Info().Msg("hello world")
	
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			// Out is the output destination.
			Out:     os.Stdout, // Stdout Stdin Stderr
			// TimeFormat specifies the format for timestamp in output.
			TimeFormat: time.RFC3339,
			// NoColor disables the colorized output.
			NoColor: false, // ปรับสี
		},
	)
	
	//เพื่อสร้าง Engine instance ของ Gin 
	//มี middleware Logger และ Recovery ติดตั้งมาให้
	// r := gin.Default() 
	//เหมือน gin.Default() ; Full
	app := gin.New()

	//middleware
	app.Use(DebugLogger())

	log.Info().Str("language", "golang").Msg("Start API Server localhost:8080")
	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	app.Use(logger.SetLogger() )

	//mysql
	system := CustomerHandler{}
	system.Initialize()

	//app router
	app.GET("/customers", system.GetAllCustomer)
	app.GET("/customers/:id", system.GetCustomer)
	app.POST("/customers", system.SaveCustomer)
	app.PUT("/customers/:id", system.UpdateCustomer)
	app.DELETE("/customers/:id", system.DeleteCustomer)

	return app
}

/* 	
	########################################################################################################
	######################################### ROUTER&CONTROLLER ############################################
	########################################################################################################
*/

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	log.Info().Str("functionName", "GetAllCustomer").Msg("request function")

	customers := []Customer{}

	h.DB.Find(&customers)

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	log.Debug().Str("method", "GET").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "GetCustomer").Msg("request function")

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Not Found id:"+id+" on Database Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "GetCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	log.Info().Str("functionName", "SaveCustomer").Msg("request function")

	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Bind JSON Fail Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "SaveCustomer").Msg("Bind JSON Fail Error: "+err.Error())
		c.JSON(http.StatusBadRequest,Msg)
		return
	}
	// log.Debug().Str("method", "POST").Msg("payload="+customer)

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Insert Database Fail Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "SaveCustomer").Msg(Msg)
		c.Status(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest,Msg)
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	log.Debug().Str("method", "PUT").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "UpdateCustomer").Msg("request function")

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Not Found id:"+id+" on Database Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Bind JSON Fail Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg("Bind JSON Fail Error: "+err.Error())
		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Insert Database Fail Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
		c.JSON(http.StatusBadRequest,Msg)
		return
	}
	
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {

	log.Debug().Str("method", "DELETE").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "DeleteCustomer").Msg("request function")

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Not Found id:"+id+" on Database Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "DeleteCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		Msg := "Delete Database Fail Error: "+err.Error()// err.Error() conv to string
		log.Error().Str("functionName", "DeleteCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	c.Status(http.StatusNoContent)
}

























