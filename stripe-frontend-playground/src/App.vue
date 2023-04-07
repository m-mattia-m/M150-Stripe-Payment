<template>
  <div class="pa-16 w-50">
    <h1>Demo: Stripe-Playground</h1>
    <h2>Customer</h2>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.customer.firstname"
          label="Firstname"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.customer.lastname"
          label="Lastname"
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.customer.street"
          label="Street"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.customer.street_nr"
          label="Number"
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.customer.post_cod"
          label="Post code"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.customer.city"
          label="City"
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.customer.province"
          label="State"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.customer.country"
          label="Country"
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.customer.description"
          label="Description"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.customer.email"
          label="E-Mail"
          required
      ></v-text-field>
    </div>
    <div class="pl-8"></div>
    <h2>Credit card</h2>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.credit_card.cardnumber"
          label="Card number"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.credit_card.cvc"
          label="CVC"
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.credit_card.expire_month"
          label="Month"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.credit_card.expire_year"
          label="Year"
          required
      ></v-text-field>
    </div>
    <div class="pl-8"></div>
    <h2>Product</h2>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.amount"
          label="Amount"
          required
      ></v-text-field>
      <div class="pl-4"></div>
      <v-text-field
          v-model="requestData.product_id"
          label="Product"
          disabled
          required
      ></v-text-field>
    </div>
    <div class="v-row pt-8">
      <v-text-field
          v-model="requestData.currency"
          label="currency"
          disabled
          required
      ></v-text-field>
    </div>
    <v-btn @click="submit">
      Submit
    </v-btn>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from 'vue-class-component';
import {MainPaymentRequestWithCustomer, PaymentApi} from "./api";
import {generateAxiosConfig} from "./helpers/axios-config";


@Options({
  components: {},
})
export default class App extends Vue {
  api = new PaymentApi(generateAxiosConfig())

  requestData: MainPaymentRequestWithCustomer = {
    amount: 100,
    credit_card: {
      cardnumber: "4242424242424242",
      cvc: "123",
      expire_month: 1,
      expire_year: 2030
    },
    currency: "CHF",
    customer: {
      city: "Samplecity",
      country: "Switzerland",
      description: "John Doe Develop-Sample",
      email: "john.doe@sapmle.com",
      firstname: "John",
      lastname: "Doe",
      post_cod: "1234",
      province: "SampleState",
      street: "Samplestreet",
      street_nr: "1"
    },
    product_id: "prod_NP1Y2WB3NQSJYA"
  }



  async submit(){
    const response = await this.api.payWithCustomerPost(this.requestData)
    console.log(response)
  }

}
</script>

<style scoped lang="scss">


</style>
