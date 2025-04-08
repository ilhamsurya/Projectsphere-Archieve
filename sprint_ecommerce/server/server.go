package server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ramadhan1445sprint/sprint_ecommerce/config"
	"github.com/ramadhan1445sprint/sprint_ecommerce/controller"
	"github.com/ramadhan1445sprint/sprint_ecommerce/customErr"
	"github.com/ramadhan1445sprint/sprint_ecommerce/middleware"
	"github.com/ramadhan1445sprint/sprint_ecommerce/repo"
	"github.com/ramadhan1445sprint/sprint_ecommerce/svc"
)

type Server struct {
	db  *sqlx.DB
	app *fiber.App
}

// var (
// 	reqDurationHist = promauto.NewHistogramVec(prometheus.HistogramOpts{
// 		Name:    "shopifyx_http_request_duration_seconds",
// 		Help:    "Duration of HTTP requests.",
// 		Buckets: prometheus.LinearBuckets(1, 1, 10),
// 	}, []string{"path", "method", "status"})
// )

func NewServer(db *sqlx.DB) *Server {
	reqDurationHist := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "shopifyx_http_request_duration_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Use(func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		method := ctx.Method()
		path := ctx.Route().Path
		rawCode := ctx.Response().StatusCode()
		statusCode := strconv.Itoa(rawCode)
		fmt.Println("Default Status Code", ctx.Response().StatusCode())

		if err != nil {
			if customError, ok := err.(customErr.CustomError); ok {
				statusCode = strconv.Itoa(customError.Status())
			} else if rawCode == fiber.StatusOK || rawCode == fiber.StatusCreated {
				statusCode = "500"
			}
		}

		elapsedDuration := time.Since(start).Seconds()

		fmt.Println(method, path, statusCode, elapsedDuration)
		reqDurationHist.WithLabelValues(path, method, statusCode).Observe(elapsedDuration)

		return err
	})

	return &Server{
		app: app,
		db:  db,
	}
}

func (s *Server) Run() {
	s.app.Listen(":8000")
}

func (s *Server) RegisterRoute() {
	mainRoute := s.app.Group("/v1")
	registerUserRouter(mainRoute, s.db)
	registerHealthCheckRoute(mainRoute, s.db)
	registerProductRouter(mainRoute, s.db)
	mainRoute.Use(middleware.Authorization)
	// put another route below
	registerImageRouter(mainRoute)
	registerBankAccountRoute(mainRoute, s.db)
}

func registerProductRouter(r fiber.Router, db *sqlx.DB) {
	c := controller.NewController(svc.NewSvc(repo.NewRepo(db)))
	payment := controller.NewPaymentController(svc.NewPaymentSvc(repo.NewPaymentRepo(db)))
	stock := controller.NewStockController(svc.NewStockSvc(repo.NewStockRepo(db)))

	prodRouter := r.Group("/product")
	prodRouter.Get("/", middleware.ProductPageAuth, c.GetListProduct)
	prodRouter.Get("/:productId", c.GetDetailProduct)

	prodRouter.Use(middleware.Authorization)
	prodRouter.Post("/", c.CreateProduct)
	prodRouter.Patch("/:productId", c.UpdateProduct)
	prodRouter.Delete("/:productId", c.DeleteProduct)
	prodRouter.Post("/:productId/buy", payment.CreatePayment)
	prodRouter.Post("/:productId/stock", stock.UpdateStock)
}
func registerUserRouter(r fiber.Router, db *sqlx.DB) {
	ctr := controller.NewUserController(svc.NewUserSvc(repo.NewUserRepo(db)))
	userGroup := r.Group("/user")

	userGroup.Post("/register", ctr.Register)
	userGroup.Post("/login", ctr.Login)
}

func registerImageRouter(r fiber.Router) {
	cfg, err := awsCfg.LoadDefaultConfig(
		context.Background(),
		awsCfg.WithRegion("ap-southeast-1"),
		awsCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.GetString("S3_ID"),
				config.GetString("S3_SECRET_KEY"),
				"",
			),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	ctr := controller.NewImageController(svc.NewImageSvc(cfg))

	r.Post("/image", ctr.UploadImage)
}

func registerHealthCheckRoute(r fiber.Router, db *sqlx.DB) {
	c := controller.NewHealthController(svc.NewHealthSvc(repo.NewHealthRepo(db)))
	healthRouter := r.Group("/health")
	healthRouter.Get("/", c.Get)
}

func registerBankAccountRoute(r fiber.Router, db *sqlx.DB) {
	bankAccount := controller.NewBankAccountController(svc.NewBankAccounthSvc(repo.NewBankAccountRepo(db)))
	productRouter := r.Group("/bank")
	productRouter.Post("/account", bankAccount.CreateBankAccount)
	productRouter.Get("/account", bankAccount.GetBankAccount)
	productRouter.Patch("/account/:bankAccountId", bankAccount.UpdateBankAccount)
	productRouter.Delete("/account/:bankAccountId", bankAccount.DeleteBankAccount)
}
