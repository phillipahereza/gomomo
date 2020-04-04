# MTN MoMo API Go Client</h1>

<strong>Power your apps with our MTN MoMo API</strong>

# Usage

## Installation

```bash
 $ go get -u github.com/phillipahereza/mtnmomo
```

# Sandbox Environment

## Creating a sandbox environment API user 

Next, we need to get the `User ID` and `User Secret` and to do this we shall need to use the Primary Key for the Product to which we are subscribed, as well as specify a host. This package ships with a CLI tool that helps to create sandbox credentials. 
It assumes you have created an account on `https://momodeveloper.mtn.com` and have your `Ocp-Apim-Subscription-Key`. 

```bash
$ momoapi sandbox -callback http://ahereza.dev -key 0d31d966e5674a999c82772aa95f2cca
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

	client := momo.NewClient(collectionPK, "sandbox", "https://sandbox.momodeveloper.mtn.com/")
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

	balance, err := client.Collection.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", balance)

	active, err := client.Collection.IsPayeeActive(ctx, "256789997290")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(active)
}
```