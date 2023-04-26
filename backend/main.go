package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"stripe/docs"
)

// @title Stripe Payment Sample
// @version 1.0
// @description this is an stripe payment example

// @contact.name API Support
// @contact.email dev@mattiamueggler.ch

// @host      localhost:3000
// @BasePath  /api/v1

func main() {
	godotenv.Load()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.POST("/api/v1/pay-with-customer", CreatePaymentWithCustomer)
	r.POST("/api/v1/pay", CreatePaymentWithoutCustomer)
	r.POST("/api/v1/pay-subscription", PaySubscription)
	r.POST("/api/v1/create-subscription", CreateSubscription)
	r.GET("/api/v1/products", ListProducts)

	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":3000")
}

// CreatePaymentWithCustomer 	godoc
// @title           			CreatePaymentWithCustomer
// @description     			Create a payment with customer
// @Tags 						Payment
// @Router  					/pay-with-customer [post]
// @Accept 						json
// @Produce						json
// @Param        				request    body     PaymentRequestWithCustomer  true  "PaymentRequest"
// @Success      				200  {object} SuccessPayment
// @Failure      				400  {object} Error
// @Failure      				404  {object} Error
// @Failure      				500  {object} Error
func CreatePaymentWithCustomer(c *gin.Context) {
	var paymentRequest PaymentRequestWithCustomer
	err := c.Bind(&paymentRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	client, err := NewStripeClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	creditCard, err := client.CreateCreditCard(
		paymentRequest.CreditCard.CardNumber,
		paymentRequest.CreditCard.ExpireMonth,
		paymentRequest.CreditCard.ExpireYear,
		paymentRequest.CreditCard.Cvc,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	customer, err := client.CreateCustomer(
		paymentRequest.Customer,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	//response, err := client.AddDefaultPaymentMethodToCustomer(
	//	customer.ID,
	//	creditCard.ID,
	//)
	//_ = response
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
	//		Message: fmt.Sprintf("%s", err),
	//	})
	//	return
	//}

	payment, err := client.CreatePayment(
		paymentRequest.Amount,
		paymentRequest.Currency,
		creditCard.ID,
		customer.ID,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessPayment{
		Message: fmt.Sprintf("payment with id '%s' was successfully", payment.ID),
	})
}

// CreatePaymentWithoutCustomer godoc
// @title           			CreatePaymentWithoutCustomer
// @description     			Create a payment without customer
// @Tags 						Payment
// @Router  					/pay [post]
// @Accept 						json
// @Produce						json
// @Param        				request    body     PaymentRequest  true  "PaymentRequest"
// @Success      				200  {object} SuccessPayment
// @Failure      				400  {object} Error
// @Failure      				404  {object} Error
// @Failure      				500  {object} Error
func CreatePaymentWithoutCustomer(c *gin.Context) {
	var paymentRequest PaymentRequest
	err := c.Bind(&paymentRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	client, err := NewStripeClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	creditCard, err := client.CreateCreditCard(
		paymentRequest.CreditCard.CardNumber,
		paymentRequest.CreditCard.ExpireMonth,
		paymentRequest.CreditCard.ExpireYear,
		paymentRequest.CreditCard.Cvc,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	payment, err := client.CreatePayment(
		paymentRequest.Amount,
		paymentRequest.Currency,
		creditCard.ID,
		"",
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	c.JSON(http.StatusOK, SuccessPayment{
		Message: fmt.Sprintf("payment with id '%s' was successfully", payment.ID),
	})
}

// PaySubscription				godoc
// @title           			PaySubscription
// @description     			Pay a Subscription
// @Tags 						Payment
// @Router  					/pay-subscription [post]
// @Accept 						json
// @Produce						json
// @Param        				request    body     PaymentRequest  true  "PaymentRequest"
// @Success      				200  {object} SuccessPayment
// @Failure      				400  {object} Error
// @Failure      				404  {object} Error
// @Failure      				500  {object} Error
func PaySubscription(c *gin.Context) {

}

// CreateSubscription			godoc
// @title           			CreateSubscription
// @description     			Create a Subscription
// @Tags 						Product
// @Router  					/create-subscription [post]
// @Accept 						json
// @Produce						json
// @Param        				request    body     PaymentRequest  true  "PaymentRequest"
// @Success      				200  {object} SuccessPayment
// @Failure      				400  {object} Error
// @Failure      				404  {object} Error
// @Failure      				500  {object} Error
func CreateSubscription(c *gin.Context) {
	var paymentRequest PaymentRequestWithCustomer
	err := c.Bind(&paymentRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	client, err := NewStripeClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	customer, err := client.CreateCustomer(
		paymentRequest.Customer,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
		return
	}

	payment, err := client.CreateSubscription(
		customer.ID,
		paymentRequest.ProductId,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	c.JSON(http.StatusOK, SuccessPayment{
		Message: fmt.Sprintf("payment with id '%s' was successfully", payment.ID),
	})
}

// ListProducts					godoc
// @title           			ListProducts
// @description     			List all products
// @Tags 						Products
// @Router  					/products [get]
// @Accept 						json
// @Produce						json
// @Success      				200  {object} []stripe.Product
// @Failure      				400  {object} Error
// @Failure      				404  {object} Error
// @Failure      				500  {object} Error
func ListProducts(c *gin.Context) {
	client, err := NewStripeClient()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	products, err := client.ListProducts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("%s", err),
		})
	}

	c.JSON(http.StatusOK, products)
}

func getEnv(key string) (string, error) {
	env, err := os.LookupEnv(key) // GetEnv() does not work
	if !err || env == "" {
		return env, errors.New(fmt.Sprintf("getenv: environment variable '%s' empty", key))
	}
	return env, nil
}

type PaymentRequestWithCustomer struct {
	Currency   string     `json:"currency"`
	Amount     int64      `json:"amount"`
	ProductId  string     `json:"product_id"`
	Customer   Customer   `json:"customer"`
	CreditCard CreditCard `json:"credit_card"`
}

type PaymentRequest struct {
	Currency   string     `json:"currency"`
	Amount     int64      `json:"amount"`
	ProductId  string     `json:"product_id"`
	CreditCard CreditCard `json:"credit_card"`
}

type Customer struct {
	Description string `json:"description"`
	Email       string `json:"email"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Street      string `json:"street"`
	StreetNr    string `json:"street_nr"`
	PostCod     string `json:"post_cod"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Province    string `json:"province"`
}

type CreditCard struct {
	CardNumber  string `json:"cardnumber"`
	ExpireMonth int64  `json:"expire_month"`
	ExpireYear  int64  `json:"expire_year"`
	Cvc         string `json:"cvc"`
}

type SuccessPayment struct {
	Message string `json:"message"`
}

type Error struct {
	Message string `json:"message"`
}
