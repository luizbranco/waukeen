package sqlite

import (
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/luizbranco/waukeen"
)

func TestCreateAccount(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	acc := &waukeen.Account{
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}
	err = db.CreateAccount(acc)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if acc.ID != "1" {
		t.Errorf("wants account id 1, got %s", acc.ID)
	}

	acc = &waukeen.Account{Number: ""}
	err = db.CreateAccount(acc)

	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestFindAccounts(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := []waukeen.Account{
		{
			ID:       "1",
			Number:   "12345",
			Name:     "Banking",
			Type:     waukeen.Checking,
			Currency: "CAD",
			Balance:  1000,
		},
		{
			ID:       "2",
			Number:   "67890",
			Name:     "Banking",
			Type:     waukeen.Savings,
			Currency: "USD",
			Balance:  -500,
		},
	}

	for _, a := range want {
		err = db.CreateAccount(&a)

		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	got, err := db.FindAccounts()
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestFindAccount(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err = db.CreateAccount(want)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	got, err := db.FindAccount("12345")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestUpdateAccount(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := &waukeen.Account{
		ID:       "1",
		Number:   "12345",
		Name:     "Banking",
		Type:     waukeen.Checking,
		Currency: "CAD",
		Balance:  1000,
	}

	err = db.CreateAccount(want)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	want.Number = "02468"
	err = db.UpdateAccount(want)
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	got, err := db.FindAccount("02468")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}
}

func TestCreateTransaction(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	tr := &waukeen.Transaction{
		AccountID:   "1",
		FITID:       "12345",
		Type:        waukeen.Debit,
		Title:       "First Transaction",
		Alias:       "Renamed Transaction",
		Description: "Surcharge",
		Amount:      9999,
		Date:        time.Now(),
	}
	err = db.CreateTransaction(tr)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if tr.ID != "1" {
		t.Errorf("wants transaction id 1, got %s", tr.ID)
	}

	tr = &waukeen.Transaction{FITID: ""}
	err = db.CreateTransaction(tr)

	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestCreateRule(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	r := &waukeen.Rule{
		AccountID: "1",
		Type:      waukeen.TagRule,
		Match:     "dominos",
		Result:    "pizza",
	}
	err = db.CreateRule(r)

	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if r.ID != "1" {
		t.Errorf("wants transaction id 1, got %s", r.ID)
	}

	r = &waukeen.Rule{Match: ""}
	err = db.CreateRule(r)

	if err == nil {
		t.Errorf("wants error, got none")
	}
}

func TestFindRules(t *testing.T) {
	db, err := New("waukeen_test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("waukeen_test.db")

	want := []waukeen.Rule{
		{
			AccountID: "1",
			ID:        "1",
			Type:      waukeen.TagRule,
			Match:     "dominos",
			Result:    "pizza",
		},
	}

	for _, r := range want {
		err = db.CreateRule(&r)

		if err != nil {
			t.Errorf("wants no error, got %s", err)
		}
	}

	got, err := db.FindRules("1")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wants %+v, got %+v", want, got)
	}

	got, err = db.FindRules("")
	if err != nil {
		t.Errorf("wants no error, got %s", err)
	}

	if len(got) != 0 {
		t.Errorf("wants empty rules, got %+v", got)
	}
}