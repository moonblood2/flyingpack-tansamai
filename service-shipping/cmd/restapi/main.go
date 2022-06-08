package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jna-distribution/service-shipping/cmd/restapi/handler"
	"github.com/jna-distribution/service-shipping/cmd/restapi/middleware"
	"github.com/jna-distribution/service-shipping/cmd/restapi/responses"
	"github.com/jna-distribution/service-shipping/internal"
	"github.com/jna-distribution/service-shipping/internal/contact"
	"github.com/jna-distribution/service-shipping/internal/courier"
	"github.com/jna-distribution/service-shipping/internal/destination"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/internal/order"
	"github.com/jna-distribution/service-shipping/internal/ordparcel"
	"github.com/jna-distribution/service-shipping/internal/ordparflash"
	"github.com/jna-distribution/service-shipping/internal/ordparshippop"
	"github.com/jna-distribution/service-shipping/internal/ordproduct"
	"github.com/jna-distribution/service-shipping/internal/origin"
	"github.com/jna-distribution/service-shipping/internal/product"
	"github.com/jna-distribution/service-shipping/internal/provider"
	"github.com/jna-distribution/service-shipping/internal/sender"
	"github.com/jna-distribution/service-shipping/internal/user"
	"github.com/jna-distribution/service-shipping/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	//Loading the environment variables from .env file at root project.
	env, err := pkg.LoadEnv()
	if err != nil {
		log.Fatalln(err)
	}

	//Initialize DB and try to connect.
	db, err := sql.Open("pgx", env["DB_DSN"])
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	//Read constant/courier.json
	courierJSONFilePath, _ := filepath.Abs("constant/courier.json")
	courierJSONFile, err := os.Open(courierJSONFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	//Courier repository
	courierRepository, err := courier.NewJSONRepository(courierJSONFile)
	if err != nil {
		log.Fatalln(err)
	}
	//Initializing repositories. PostgreSQL case.
	userRepository := user.NewPostgresRepository(db)
	contactRepository := contact.NewPostgresRepository(db)
	productRepository := product.NewPostgresRepository(db)
	orderParcelRepository := ordparcel.NewPostgresRepository(db)
	orderProductRepository := ordproduct.NewPostgresRepository(db)
	senderRepository := sender.NewPostgresRepository(db)
	originRepository := origin.NewPostgresRepository(db)
	destinationRepository := destination.NewPostgresRepository(db)
	//Provider repository
	orderShippopRepository := ordparshippop.NewPostgresRepository(db)
	orderFlashRepository := ordparflash.NewPostgresRepository(db)

	//Initializing shipping providers.
	providers := provider.NewProviderList(
		map[entity.ProviderCode]provider.Provider{
			entity.ProviderCodeShippop: provider.NewShippopProvider(
				env["PROVIDER_HTTP_API_PROXY_URI"],
				env["PROVIDER_SHIPPOP_API_URL"],
				env["PROVIDER_SHIPPOP_API_KEY"],
				orderShippopRepository,
			),
			entity.ProviderCodeFlash: provider.NewFlashProvider(
				env["PROVIDER_HTTP_API_PROXY_URI"],
				env["PROVIDER_FLASH_API_URL"],
				env["PROVIDER_FLASH_API_KEY"],
				env["PROVIDER_FLASH_MCH_ID"],
				orderFlashRepository,
			),
		},
	)

	//JNA address
	jnaAddress := entity.Origin{ContactInfo: entity.ContactInfo{
		Name:        env["JNA_NAME"],
		PhoneNumber: env["JNA_PHONE_NUMBER"],
		Address:     env["JNA_ADDRESS"],
		District:    env["JNA_DISTRICT"],
		State:       env["JNA_STATE"],
		Province:    env["JNA_PROVINCE"],
		Postcode:    env["JNA_POSTCODE"],
	}}
	jnaAddressAsOrigin := env["JNA_ADDRESS_AS_ORIGIN"] == "true"
	anUrl := env["AN_URL"]
	//Initializing services
	courierService := courier.NewService(courierRepository)
	userService := user.NewService(userRepository, contactRepository)
	productService := product.NewService(productRepository)
	orderService := order.NewService(
		productRepository,
		senderRepository,
		originRepository,
		destinationRepository,
		orderParcelRepository,
		orderProductRepository,
		providers,
		jnaAddress,
		jnaAddressAsOrigin,
		anUrl,
	)

	//Initializing rootMux
	rootMux := http.NewServeMux()

	//Building handlers for each module.
	handler.BuildUserHandler(rootMux, userService)
	handler.BuildAuthHandler(rootMux, userService)
	handler.BuildOrderHandler(rootMux, orderService)
	handler.BuildProductHandler(rootMux, productService)
	handler.BuildCourierProductHandler(rootMux, courierService, productService)

	//Handle Index
	rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rw := responses.NewJSONResponseWriter(w)
		if r.RequestURI != "/" {
			e := responses.BuildErrorResponse(internal.ErrNotFound{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		if r.Method != http.MethodGet {
			e := responses.BuildErrorResponse(responses.ErrMethodNotAllowed{})
			rw.WriteWithStatus(e.Status, e)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>service-shipping</h1>"))
	})

	//Setup middlewares for root rootMux
	corsMw := middleware.WithCORS(rootMux)
	loggingMw := middleware.WithLogging(corsMw)

	//Initializing server
	server := &http.Server{
		Handler: loggingMw,
		Addr:    fmt.Sprintf(":%v", env["PORT"]),
	}

	log.Printf("Server is running on: %v", server.Addr)

	//Starting to listen and serve.
	log.Fatalln(server.ListenAndServe())
}
