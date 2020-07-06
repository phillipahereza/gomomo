# MTN MoMo API Go Client</h1>

<strong>Power your apps with our MTN MoMo API</strong>

<div>
  Join the active, engaged community: <br>
  <a href="https://momodeveloper.mtn.com/">Website</a>
  <span> | </span>
  <a href="https://spectrum.chat/momo-api-developers/">Spectrum</a>
  <br><br>
</div>

[![phillipahereza](https://circleci.com/gh/phillipahereza/gomomo.svg?style=svg)](https://github.com/phillipahereza/gomomo)
[![Go Report Card](https://goreportcard.com/badge/github.com/phillipahereza/momoapi-go)](https://goreportcard.com/badge/github.com/phillipahereza/momoapi-go)
[![Join the community on Spectrum](https://withspectrum.github.io/badge/badge.svg)](https://spectrum.chat/momo-api-developers/)

# Usage

## Installation

```bash
 $ go get -u github.com/phillipahereza/gomomo
```

# Sandbox Environment

## Creating a sandbox environment API user 

Next, we need to get the `User ID` and `User Secret` and to do this we shall need to use the Primary Key for the Product to which we are subscribed, as well as specify a host. 
This package ships with a CLI tool that helps to create sandbox credentials. 
It assumes you have created an account on `https://momodeveloper.mtn.com` and have your `Ocp-Apim-Subscription-Key`. 

To install the CLI
```bash
go get github.com/phillipahereza/gomomo/cmd/momocli 
```
Then run the following command to create a sandbox user and get their API key
```bash
$ momocli sandbox -callback http://ahereza.dev -key 0d31d966e5674a999c82772aa95f2cca
```

The `providerCallBackHost` is your callback host and `Ocp-Apim-Subscription-Key` is your API key for the specific product to which you are subscribed. 
The `API Key` is unique to the product and you will need an `API Key` for each product you use. You should get a response similar to the following:

```bash
API Key: cbd4aa5d0929439ab4760ec10762b9c5
User ID: ee67c3b9-0357-4351-ac88-e47340213bb1
```

## Configuration

Before we can fully utilize the library, we need to specify global configurations. The global configuration must contain the following:

* `BASE_URL`: An optional base url to the MTN Momo API. By default the staging base url will be used
* `ENVIRONMENT`: Either "sandbox" or "production". Default is 'sandbox'
* `CALLBACK_HOST`: The domain where you webhooks urls are hosted. This is mandatory.

Once you have specified the global variables, you can now provide the product-specific variables. 
Each MoMo API product requires its own authentication details i.e its own `Subscription Key`, `User ID` and 
`User Secret`, also sometimes refered to as the `API Secret`. As such, we have to configure subscription keys for 
each product as show below.

## Collection

* `collectionPK`: Primary Key for the `Collection` product on the developer portal.
* `userID`: For sandbox, use the one generated with the `mtnmomo` command.
* `apiKey`: For sandbox, use the one generated with the `mtnmomo` command.

*Note: `userID` and `apiKey` for production are provided on the MTN OVA dashboard*

An example of how to use the `Collection` product:
```go
func main() {
	collectionPK := "0d31d966e5674a999c82772aa95f2cca"
	userID := "356c0142-21a5-4ff3-9f00-87c268a60378"
	apiKey := "c0d857dba3944ce3b6d436c04963e1ea"

	ctx := context.Background()

	client := gomomo.NewClient(collectionPK, "sandbox", "https://sandbox.momodeveloper.mtn.com/")

	// Calling GetToken fetches an access token and sets it onto the client
	// All subsequent calls using the same client will be automatically set with the Authorization header using this token
	_, err := client.Collection.GetToken(ctx, apiKey, userID)
	if err != nil {
		log.Fatal(err)
	}

	transactionID, err := client.Collection.RequestToPay(ctx, "46733123453", 500, "2323", "payee Note", "Payer Message", "EUR")
	if err != nil {
		log.Fatal(err)
	}

	status, err := client.Collection.GetTransaction(ctx, transactionID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)
    
    // Get balance endpoint on sandbox is flaky, returns a different response every time is called
	balance, err := client.Collection.GetBalance(ctx)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v", balance)

	active, err := client.Collection.IsPayeeActive(ctx, "256789997290")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(active)
}
```

### Methods

1. `RequestToPay`: This operation is used to request a payment from a consumer (Payer). The payer will be asked to authorize the payment. The transaction is executed once the payer has authorized the payment. The transaction will be in status PENDING until it is authorized or declined by the payer or it is timed out by the system.

2. `GetTransaction`: Retrieve transaction information using the `transactionId` returned by `RequestToPay`. You can invoke it at intervals until the transaction fails or succeeds. If the transaction has failed, it will throw an appropriate error. 

3. `GetBalance`: Get the balance of the account.

4. `IsPayerActive`: check if an account holder is registered and active in the system.


## Disbursement

* `disbursementPK`: Primary Key for the `Disbursement` product on the developer portal.

An example of how to use the `Disbursement` product:

```go
func main() {
	disbursementPK := "50ff0a2926784599a8409f263a5cca6c"
	userID := "757e5ab2-88f6-413c-9d6b-3399560aa1df"
	apiKey := "b58e05c2ada542068e9dedbc64663c43"

	ctx := context.Background()

	disbursementClient := gomomo.NewClient(disbursementPK, "sandbox", "https://sandbox.momodeveloper.mtn.com/")
	_, err := disbursementClient.Disbursement.GetToken(ctx, apiKey, userID)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := disbursementClient.Disbursement.GetBalance(ctx)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("%+v\n", balance)
	}

	active, err := disbursementClient.Disbursement.IsPayeeActive(ctx, "256789997290")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payee is Active: %v\n", active)

	transactionID, err := disbursementClient.Disbursement.Transfer(ctx, "46733123453", 500, "2323", "", "", "EUR")
	if err != nil {
		log.Fatal(err)
	}

	status, err := disbursementClient.Disbursement.GetTransfer(ctx, transactionID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)

}
```

### Methods

1. `Transfer`: This operation is used to transfer an amount from the owner’s account to a payee account. The payer will be asked to authorize the payment. The transaction is executed once the payer has authorized the payment. The transaction will be in status PENDING until it is authorized or declined by the payer or it is timed out by the system.

2. `GetTransfer`: Retrieve transfer information using the `transactionId` returned by `Transfer`. You can invoke it at intervals until the transaction fails or succeeds. If the transaction has failed, it will throw an appropriate error. 

3. `GetBalance`: Get the balance of the account.

4. `IsPayerActive`: check if an account holder is registered and active in the system.

## Remittance

* `remittancePK`: Primary Key for the `Remittance` product on the developer portal.

An example of how to use the `Remittance` product:
```go
func main() {
	remittancePK := "24f5444da2c74fb59e6e4de50798db5d"
	userID := "712db4b5-f079-47bf-843b-c86f6784ebe0"
	apiKey := "3a8723537dc04729a83c1b7a09cc4e54"

	ctx := context.Background()

	remittanceClient := gomomo.NewClient(remittancePK, "sandbox", "https://sandbox.momodeveloper.mtn.com/")
	_, err := remittanceClient.Remittance.GetToken(ctx, apiKey, userID)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := remittanceClient.Remittance.GetBalance(ctx)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("%+v\n", balance)
	}

	active, err := remittanceClient.Remittance.IsPayeeActive(ctx, "256789997290")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payee is Active: %v\n", active)

	transactionID, err := remittanceClient.Remittance.Transfer(ctx, "46733123453", 500, "2323", "", "", "EUR")
	if err != nil {
		log.Fatal(err)
	}

	status, err := remittanceClient.Remittance.GetTransfer(ctx, transactionID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", status)

}
```

## License

GNU GPLv3

Copyright (C) 2020 Phillip Ahereza

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.