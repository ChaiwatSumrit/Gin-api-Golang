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

	//mysql logger
	"github.com/rs/zerolog"
	// "github.com/rs/zerolog/hlog"
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

// fake alice to avoid dep
type middleware func(http.Handler) http.Handler
type alice struct {
	m []middleware
}

func (a alice) Append(m middleware) alice {
	a.m = append(a.m, m)
	return a
}
func (a alice) Then(h http.Handler) http.Handler {
	for i := range a.m {
		h = a.m[len(a.m)-1-i](h)
	}
	return h
}

func init() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
}


func main(){
		
	// host := "localhost"
	// log := zerolog.New(os.Stdout).With().
    // Timestamp().
    // Str("role", "my-service").
    // Str("host", host).
    // Logger()
		
	// // c := alice.New()
	// var c alice

	// // Install the logger handler with default output on the console
	// c = c.Append(hlog.NewHandler(log))

	// // Install some provided extra handler to set some request's context fields.
	// // Thanks to those handler, all our logs will come with some pre-populated fields.
	// c = c.Append(hlog.RemoteAddrHandler("ip"))
	// c = c.Append(hlog.UserAgentHandler("user_agent"))
	// c = c.Append(hlog.RefererHandler("referer"))
	// c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	// // Here is your final handler
	// h := c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// Get the logger from the request's context. You can safely assume it
	// 	// will be always there: if the handler is removed, hlog.FromRequest
	// 	// will return a no-op logger.
	// 	hlog.FromRequest(r).Info().
	// 		Str("user", "current user").
	// 		Str("status", "ok").
	// 		Msg("Something happened")

	// 	// Output: {"level":"info","time":"2001-02-03T04:05:06Z","role":"my-service","host":"local-hostname","req_id":"b4g0l5t6tfid6dtrapu0","user":"current user","status":"ok","message":"Something happened"}
	// }))
	// http.Handle("/", h)

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal().Err(err).Msg("Startup failed")
	// }


	// Example skip path request.
	app := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	app.Run(":8080")
	// log.Info().Str("language", "golang").Msg("Msg='Start API Server localhost:8080'")
	// log.Info().Str("language", "golang").Msg("Msg='Start API Server localhost:8080'")

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
		log.Error().Str("Database", "Initialize").Msg("Msg='Connect MYSQL Database Fail Error : "+err.Error()+"'")
	}

	log.Info().Str("Database", "Initialize").Msg("Msg='Connect Database Success root : best1459900574821@tcp(127.0.0.1:3306)/charset'")

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

		if c.Request.Method == http.MethodPost {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()

			log.Debug().Str("method", "POST").Msg("Msg='payload="+buf.String()+"'")


		}else if c.Request.Method == http.MethodPut {
			var buf bytes.Buffer
			newBody := &MyReadCloser{c.Request.Body, &buf}
			c.Request.Body = newBody
			c.Next()

			log.Debug().Str("method", "PUT").Msg("Msg='payload="+buf.String()+"'")

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

	//log fomat json
	log.Info().Str("language", "golang").Msg("Msg='Start API Server localhost:8080'")

	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
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
			Out: os.Stdout, // Stdout Stdin Stderr
			// TimeFormat specifies the format for timestamp in output.
			TimeFormat: time.RFC3339,
			// NoColor disables the colorized output.
			NoColor: true, // ปรับสี
		},
	)
	
	//เพื่อสร้าง Engine instance ของ Gin 
	//มี middleware Logger และ Recovery ติดตั้งมาให้
	// r := gin.Default() 
	//เหมือน gin.Default() ; Full
	app := gin.New()

	//middleware
	app.Use(DebugLogger())

	// log.Info().Str("language", "golang").Msg("Start API Server localhost:8080")
	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.


	//mysql
	system := CustomerHandler{}
	system.Initialize()

	//app router
	app.POST("/loging", system.Loging)
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
func (h *CustomerHandler) Loging(c *gin.Context) {
	log.Info().Str("functionName", "Loging").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	// if err := c.ShouldBindJSON(&customer); err != nil {
	// 	Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
	// 	log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
	// 	c.JSON(http.StatusBadRequest,Msg)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"msg":"ok"})
}

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	log.Info().Str("functionName", "GetAllCustomer").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	customers := []Customer{}

	h.DB.Find(&customers)

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	log.Debug().Str("method", "GET").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "GetCustomer").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "GetCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	log.Info().Str("functionName", "GetCustomer").Msg("Msg='CodeStarus:200 Request Success'")
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	log.Info().Str("functionName", "SaveCustomer").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "SaveCustomer").Msg(Msg)
		c.JSON(http.StatusBadRequest,Msg)
		return
	}
	// log.Debug().Str("method", "POST").Msg("payload="+customer)

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "SaveCustomer").Msg(Msg)
		c.Status(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest,Msg)
	}

	log.Info().Str("functionName", "SaveCustomer").Msg("Msg='CodeStarus:200 Request Success'")
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	log.Debug().Str("method", "PUT").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "UpdateCustomer").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		Msg := "Msg='CodeStarus:400 BadRequest "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
		c.JSON(http.StatusBadRequest,Msg)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		Msg := "Msg='Insert Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "UpdateCustomer").Msg(Msg)
		c.JSON(http.StatusBadRequest,Msg)
		return
	}
	
	log.Info().Str("functionName", "UpdateCustomer").Msg("Msg='CodeStarus:200 Request Success'")
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {

	log.Debug().Str("method", "DELETE").Msg("param="+c.Param("id"))

	log.Info().Str("functionName", "DeleteCustomer").Msg("Msg='Request Function'")
	log.Debug().Msg("path="+c.Request.RequestURI)

	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		Msg := "Msg='CodeStarus:404 Not Found id:"+id+" on Database Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "DeleteCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		Msg := "Msg='Delete Database Fail Error: "+err.Error()+"'"// err.Error() conv to string
		log.Error().Str("functionName", "DeleteCustomer").Msg(Msg)
		c.JSON(http.StatusNotFound,Msg)
		return
	}

	log.Info().Str("functionName", "DeleteCustomer").Msg("Msg='CodeStarus:200 Request Success'")
	c.Status(http.StatusNoContent)
}

























