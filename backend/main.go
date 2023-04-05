package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.POST("/api/v1/pay", CreatePayment)

	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":3000")
}

// CreatePayment 			godoc
// @title           		GetCustomer
// @description     		Return specific customer
// @Tags 					Customer
// @Router  				/pay [post]
// @Accept 					json
// @Produce					json
// @Security				Bearer
// @Param        			request    body     PaymentRequest  true  "PaymentRequest"
// @Success      			200  {object} SuccessPayment
// @Failure      			400  {object} Error
// @Failure      			404  {object} Error
// @Failure      			500  {object} Error
func CreatePayment(c *gin.Context) {
	var paymentRequest PaymentRequest
	err := c.Bind(&paymentRequest)
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewStripeClient()
	if err != nil {
		log.Fatal(err)
	}

	customer, err := client.CreateCustomer(
		paymentRequest.Customer,
	)
	if err != nil {
		log.Fatal(err)
	}

	creditCard, err := client.CreateCreditCard(
		paymentRequest.CreditCard.CardNumber,
		paymentRequest.CreditCard.ExpireMonth,
		paymentRequest.CreditCard.ExpireYear,
		paymentRequest.CreditCard.Cvc,
	)
	if err != nil {
		log.Fatal(err)
	}

	payment, err := client.CreatePayment(
		paymentRequest.Amount,
		paymentRequest.Currency,
		creditCard.ID,
	)

	c.JSON(http.StatusOK, SuccessPayment{
		Message: "payment was successfully",
	})
	_ = customer
	_ = creditCard
	_ = payment

}

type StripeConfig struct {
	StripeClient client.API
}

func NewStripeClient() (StripeConfig, error) {
	stripeSecretKey, err := getEnv("STRIPE_SECRET_KEY")
	if err != nil {
		return StripeConfig{}, err
	}
	params := &stripe.ChargeParams{}
	sc := &client.API{}
	sc.Init(fmt.Sprintf("%s", stripeSecretKey), nil)
	sc.Charges.Get("ch_3MWQSxEqHN3Vq4Xe1BKWhWAq", params)
	stripe.Key = stripeSecretKey

	return StripeConfig{StripeClient: *sc}, nil
}

func (sc StripeConfig) CreateCustomer(customer Customer) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Address: &stripe.AddressParams{
			City:       stripe.String(customer.City),
			Country:    stripe.String(customer.Country),
			Line1:      stripe.String(fmt.Sprintf("%s %s", customer.Street, customer.StreetNr)),
			PostalCode: stripe.String(customer.PostCod),
			State:      stripe.String(customer.Province),
		},
		Description: stripe.String(customer.Description),
		Email:       stripe.String(customer.Email),
		Name:        stripe.String(fmt.Sprintf("%s %s", customer.Firstname, customer.Lastname)),
	}

	stripeCustomer, err := sc.StripeClient.Customers.New(params)
	if err != nil {
		return nil, err
	}
	return stripeCustomer, nil
}

func (sc StripeConfig) CreateCreditCard(cardNumber string, expMonth int64, expYear int64, cvc string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cardNumber),
			ExpMonth: stripe.Int64(expMonth),
			ExpYear:  stripe.Int64(expYear),
			CVC:      stripe.String(cvc),
		},
	}
	pm, err := paymentmethod.New(params)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

func (sc StripeConfig) CreatePayment(amount int64, currency string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		PaymentMethod: stripe.String(paymentMethodID),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func getEnv(key string) (string, error) {
	env, err := os.LookupEnv(key) // GetEnv() does not work
	if !err || env == "" {
		return env, errors.New(fmt.Sprintf("getenv: environment variable '%s' empty", key))
	}
	return env, nil
}

type PaymentRequest struct {
	Currency   string     `json:"currency"`
	Amount     int64      `json:"amount"`
	Customer   Customer   `json:"customer"`
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
