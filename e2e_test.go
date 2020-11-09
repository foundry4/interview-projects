package main

import (
	"net/http/httptest"
	"testing"

	"io"

	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/boseabhishek/go-shopper/shopping"

	. "github.com/smartystreets/goconvey/convey"
)

func createProductSUT() (*shopping.ProductHandler, *shopping.ProductService) {

	pm := make(map[string]*shopping.Product)
	pm["1"] = &shopping.Product{Id: "1", Product: "abc", Price: "£1.00"}
	pm["2"] = &shopping.Product{Id: "2", Product: "def", Price: "£2.00"}

	mdb := shopping.NewDB(pm)
	ps := shopping.NewProductService(mdb)
	ph := shopping.NewProductHandler(ps)

	return ph, ps
}

func createCartSUT() (*shopping.CartHandler, *shopping.ProductService) {

	cm := make(map[string]*shopping.Product)

	mdb := shopping.NewDB(cm)
	ps := shopping.NewProductService(mdb)
	ch := shopping.NewCartHandler(ps)

	return ch, ps
}

func Test_View_List_Of_Products(t *testing.T) {

	Convey("Given there are products available in the shop", t, func() {
		req, rr := setup(http.MethodGet, "/products", nil, t)
		ph, _ := createProductSUT()

		handler := http.HandlerFunc(ph.Query)
		handler.ServeHTTP(rr, req)

		Convey("When I view the home page", func() {
			ps := make([]*shopping.Product, 2)
			err := json.Unmarshal([]byte(rr.Body.String()), &ps)
			if err != nil {
				t.Errorf("error converting to json: %v", err)
			}

			Convey("Then I can see the list of products", func() {
				So(len(ps), ShouldEqual, 2)

			})

		})

	})

}

func setup(verb, url string, body io.Reader, t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodGet, "/products", body)
	if err != nil {
		t.Errorf("error creating request %v", err)
	}
	rr := httptest.NewRecorder()

	return req, rr
}

func Test_Select_And_View_The_Details_Of_A_Product(t *testing.T) {

	Convey("Given there are products available in the shop", t, func() {

		req, rr := setup(http.MethodGet, "/products", nil, t)
		ph, _ := createProductSUT()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := http.HandlerFunc(ph.QueryByID)
		handler.ServeHTTP(rr, req)

		want := &shopping.Product{Id: "1", Product: "abc", Price: "£1.00"}

		p := new(shopping.Product)
		err := json.Unmarshal([]byte(rr.Body.String()), &p)
		if err != nil {
			t.Errorf("error converting to json: %v", err)
		}

		Convey("And I am viewing the home page", func() {
			Convey("When I select a product", func() {
				Convey("Then I am presented with the product details", func() {
					Convey("And I can see the option to buy the product", func() {

						So(p.Id, ShouldEqual, want.Id)
						So(p.Price, ShouldEqual, want.Price)
						So(p.Product, ShouldEqual, want.Product)

					})

				})
			})
		})

	})
}

func Test_Add_Item_To_Basket(t *testing.T) {

	Convey("Given I am viewing the details of a product", t, func() {

		Convey("When I choose to buy the product", func() {

			// get products from inventory
			req, rr := setup(http.MethodGet, "/products", nil, t)
			ph, _ := createProductSUT()

			req = mux.SetURLVars(req, map[string]string{"id": "1"})
			handler := http.HandlerFunc(ph.QueryByID)
			handler.ServeHTTP(rr, req)

			p := new(shopping.Product)
			err := json.Unmarshal([]byte(rr.Body.String()), &p)
			if err != nil {
				t.Errorf("error converting to json: %v", err)
			}

			// adding to cart
			buf := new(bytes.Buffer)
			if p != nil {
				err := json.NewEncoder(buf).Encode(p)
				if err != nil {
					t.Errorf("error encoding model %v", err)
				}
			}
			req, rr = setup(http.MethodPost, "/cart", buf, t)
			ch, ps := createCartSUT()

			handler = http.HandlerFunc(ch.Create)
			handler.ServeHTTP(rr, req)

			cp := new(shopping.Product)
			err = json.Unmarshal([]byte(rr.Body.String()), &cp)
			if err != nil {
				t.Errorf("error converting to json: %v", err)
			}

			Convey("Then I see confirmation that the product is added to my basket", func() {
				So(len(ps.Db.Type.(map[string]*shopping.Product)), ShouldEqual, 1)

				Convey("And I see the updated value of the items in my basket", func() {
					So(p.Id, ShouldEqual, cp.Id)
					So(p.Price, ShouldEqual, cp.Price)
					So(p.Product, ShouldEqual, cp.Product)
				})
			})
		})
	})

}


func Test_Remove_Item_From_Basket(t *testing.T) {

	Convey("Given I have added several products to the basket", t, func() {

		// get products from inventory
		req, rr := setup(http.MethodGet, "/products", nil, t)
		ph, _ := createProductSUT()

		handler := http.HandlerFunc(ph.Query)
		handler.ServeHTTP(rr, req)

		ps := []*shopping.Product{}
		err := json.Unmarshal([]byte(rr.Body.String()), &ps)
		if err != nil {
			t.Errorf("error converting to json: %v", err)
		}

		// adding to cart
		ch, cs := createCartSUT()
		for _, p := range ps {
			buf := new(bytes.Buffer)
			if p != nil {
				err := json.NewEncoder(buf).Encode(p)
				if err != nil {
					t.Errorf("error encoding model %v", err)
				}
			}
			req, rr = setup(http.MethodPost, "/cart", buf, t)
			handler = http.HandlerFunc(ch.Create)
			handler.ServeHTTP(rr, req)
		}

		// check if basket contains the products added
		req, rr = setup(http.MethodGet, "/cart", nil, t)

		handler = http.HandlerFunc(ch.Query)
		handler.ServeHTTP(rr, req)

		cps := []*shopping.Product{}
		err = json.Unmarshal([]byte(rr.Body.String()), &cps)
		if err != nil {
			t.Errorf("error converting to json: %v", err)
		}

		Convey("And I am viewing the basket contents", func() {

			So(len(cps), ShouldEqual, 2)

			Convey("When I choose to remove an item from the basket", func() {

				req, rr = setup(http.MethodDelete, "/cart", nil, t)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				handler = http.HandlerFunc(ch.Delete)
				handler.ServeHTTP(rr, req)

				Convey("Then I see that the product is removed from the basket", func() {

					So(len(cs.Db.Type.(map[string]*shopping.Product)), ShouldEqual, 1)
				})
			})
		})
	})

}

func Test_Checkout(t *testing.T) {

	Convey("Given I have added at least one product to the basket", t, func() {

		// get products from inventory
		req, rr := setup(http.MethodGet, "/products", nil, t)
		ph, _ := createProductSUT()

		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		handler := http.HandlerFunc(ph.QueryByID)
		handler.ServeHTTP(rr, req)

		p := new(shopping.Product)
		err := json.Unmarshal([]byte(rr.Body.String()), &p)
		if err != nil {
			t.Errorf("error converting to json: %v", err)
		}

		// adding to cart
		buf := new(bytes.Buffer)
		if p != nil {
			err := json.NewEncoder(buf).Encode(p)
			if err != nil {
				t.Errorf("error encoding model %v", err)
			}
		}
		req, rr = setup(http.MethodPost, "/cart", buf, t)
		ch, cs := createCartSUT()

		handler = http.HandlerFunc(ch.Create)
		handler.ServeHTTP(rr, req)

		cp := new(shopping.Product)
		err = json.Unmarshal([]byte(rr.Body.String()), &cp)
		if err != nil {
			t.Errorf("error converting to json: %v", err)
		}

		Convey("When I choose to check out", func() {
			req, rr = setup(http.MethodDelete, "/checkout", nil, t)

			req = mux.SetURLVars(req, map[string]string{"id": "1"})
			handler = http.HandlerFunc(ch.Checkout)
			handler.ServeHTTP(rr, req)

			Convey("Then I see confirmation that I have checked out", func() {
				var msg string
				err = json.Unmarshal([]byte(rr.Body.String()), &msg)
				if err != nil {
					t.Errorf("error converting to json: %v", err)
				}

				So(msg, ShouldStartWith, "Item with id 1 has been successfully checked out")

				Convey("And I receive an email with the details of the purchase", func() {

					So(msg, ShouldEndWith, "Please check your email for confirmation.")

					Convey("And I see that my basket is empty", func() {

						So(len(cs.Db.Type.(map[string]*shopping.Product)), ShouldEqual, 0)

					})
				})
			})

		})
	})

}
