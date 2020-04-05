package gomomo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRemittanceServiceOp_GetToken(t *testing.T) {
	t.Run("Get token returns a 200_OK", func(t *testing.T) {
		setup()
		defer teardown()
		mux.HandleFunc(remittancesTokenURL, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"access_token": "token", "token_type": "tokenType", "expires_in": 0}`)
			w.WriteHeader(http.StatusOK)
			testMethod(t, r, http.MethodPost)
		})
		token, err := client.Remittance.GetToken(ctx, "234343434", "43434343434")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}

		if token != "token" {
			t.Errorf("Expected 'token' but got %s", token)
		}

		if client.Token != token {
			t.Errorf("Expected 'token' to be set on client but got %s", token)
		}
	})

	t.Run("GetToken returns a non 200_OK status", func(t *testing.T) {
		setup()
		defer teardown()
		mux.HandleFunc(remittancesTokenURL, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			testMethod(t, r, http.MethodPost)
		})
		_, err := client.Remittance.GetToken(ctx, "234343434", "43434343434")

		if err == nil {
			t.Errorf("Expected a non nil error but got %s", err)
		}
	})
}

func TestRemittanceServiceOp_Transfer(t *testing.T) {
	t.Run("Remittance.Transfer returns 202_ACCEPTED", func(t *testing.T) {
		setup()
		defer teardown()

		mux.HandleFunc(remittancesTransferURL, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			testMethod(t, r, http.MethodPost)
		})
		client.Token = "34534523243"
		transactionID, err := client.Remittance.Transfer(ctx, "25678999720", 500, "34232", "payee", "payer", "UGX")
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		if transactionID == "" {
			t.Errorf("Expected transactionID to be a non empty string")
		}
	})
}

func TestRemittanceServiceOp_GetTransfer(t *testing.T) {
	setup()
	defer teardown()

	expectedStatus := PaymentStatusResponse{
		Amount:                 "500",
		Currency:               "UGX",
		FinancialTransactionID: 2312,
		ExternalID:             "3232",
		Payer: paymentDetails{
			PartyIDType: "MSISDN",
			PartyID:     "4656473839",
		},
		Status: "SUCCESSFUL",
		Reason: "",
	}

	transactionID := "6c6eb16c-8b34-4d5d-bd41-2a9303f65075"

	urlStr := fmt.Sprintf("%s/%s", remittancesTransferURL, transactionID)

	mux.HandleFunc(urlStr, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expectedStatus)
		w.WriteHeader(http.StatusOK)
		testMethod(t, r, http.MethodGet)
	})

	actualStatus, err := client.Remittance.GetTransfer(ctx, transactionID)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if !reflect.DeepEqual(actualStatus, &expectedStatus) {
		t.Errorf("GetTransaction\n got=%#v\nwant=%#v", actualStatus, expectedStatus)
	}
}

func TestRemittanceServiceOp_GetBalance(t *testing.T) {
	setup()
	defer teardown()

	expectedBalance := BalanceResponse{
		AvailableBalance: "500",
		Currency:         "UGX",
	}

	mux.HandleFunc(remittancesBalanceURL, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expectedBalance)
		w.WriteHeader(http.StatusOK)
		testMethod(t, r, http.MethodGet)
	})

	actualBalance, err := client.Remittance.GetBalance(ctx)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if !reflect.DeepEqual(actualBalance, &expectedBalance) {
		t.Errorf("GetBalance\n got=%#v\nwant=%#v", actualBalance, expectedBalance)
	}
}

func TestRemittanceServiceOp_IsPayeeActive(t *testing.T) {
	setup()
	defer teardown()
	mobileNumber := "256789997290"
	urlStr := fmt.Sprintf("%s%s/active", remittancesIsAccountActiveURL, mobileNumber)

	mux.HandleFunc(urlStr, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		testMethod(t, r, http.MethodGet)
	})

	_, err := client.Remittance.IsPayeeActive(ctx, mobileNumber)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
