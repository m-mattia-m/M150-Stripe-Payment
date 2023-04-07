package main

import (
	"fmt"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/stripe/stripe-go/v74/product"
)

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

//func (sc StripeConfig) AddDefaultPaymentMethodToCustomer(customerId, cardId string) (*stripe.Customer, error) {
//
//	pmAttachParams := &stripe.PaymentMethodAttachParams{
//		Customer: stripe.String(customerId),
//	}
//
//	res, err := paymentmethod.Attach(cardId, pmAttachParams)
//
//	if err != nil {
//		return nil, err
//	}
//
//	paymentMethodParam := &stripe.CustomerParams{
//		DefaultSource: stripe.String(cardId),
//		//PaymentMethod: nil,
//		//Source:        nil,
//	}
//
//	response, err := sc.StripeClient.Customers.Update(customerId, paymentMethodParam)
//	if err != nil {
//		return nil, err
//	}
//	_ = response
//	_ = res
//	return nil, nil
//
//}

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

func (sc StripeConfig) CreatePayment(amount int64, currency string, paymentMethodID string, customerId string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		PaymentMethod: stripe.String(paymentMethodID),
	}

	if customerId != "" {
		params.Customer = stripe.String(customerId)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (sc StripeConfig) ListProducts() ([]*stripe.Product, error) {
	params := &stripe.ProductListParams{
		ListParams: stripe.ListParams{
			Limit: stripe.Int64(10),
		},
	}

	productList := product.List(params)

	products := productList.ProductList().Data

	if err := productList.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (sc StripeConfig) CreateSubscription(customerID string, priceId string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Currency: stripe.String("CHF"),
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceId),
			},
		},
	}

	subscription, err := sc.StripeClient.Subscriptions.New(params)

	if err != nil {
		return nil, err
	}
	return subscription, nil
}
