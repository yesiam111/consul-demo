package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type product struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "product-service"
	registration.Name = "product-service"
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

// check key-value moi tu consuk K-V
func Configuration(key string) (bool, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Sprintf( "Error. %s", err)
		return false, err
	}
	kvpair, _, err := consul.KV().Get(key, nil)
	if err != nil {
		fmt.Sprintf("Error. %s", err) 
		return false, err
	}	
	if kvpair == nil {
		fmt.Sprintf( "Configuration empty")
		return false, nil
	}


	if string(kvpair.Value) == "enable" {
		return true, nil
	}

	return false,nil
}

// handle url "/product"
func Products(w http.ResponseWriter, r *http.Request) {
	products := []product{
		{
			ID:    1,
			Name:  "Acer Laptop",
			Price: 2000000.00,
		},
		{
			ID:    2,
			Name:  "Western Digital HDD",
			Price: 500.00,
		},
		{
			ID:    3,
			Name:  "Dell Laptop",
			Price: 1500000.00,
		},
		{
			ID:    4,
			Name:  "Casio Watch",
			Price: 50000.00,
		},
		{
			ID:    5,
			Name:  "Lenovo Laptop",
			Price: 20000000.00,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&products)
}

func newProducts(w http.ResponseWriter, r *http.Request) {

	new_product, _ := Configuration("new-product")

	if new_product == true {
		products := []product{
			{
				ID:    1,
				Name:  "Alienware Laptop",
				Price: 2000000.00,
			},
			{
				ID:    2,
				Name:  "Samsung HDD",
				Price: 500.00,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&products)
	} else {
		products := []product{
			{
				ID:    0,
				Name:  "Not released yet",
				Price: 0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&products)
	}
}

func main() {
	registerServiceWithConsul()

	http.HandleFunc("/healthcheck", healthcheck)
	
	http.HandleFunc("/products", Products)
	
	http.HandleFunc("/new-products", newProducts)

	fmt.Printf("product service is up on port: %s", port())
	http.ListenAndServe(port(), nil)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `product service is good`)
}

func port() string {
	p := os.Getenv("PRODUCT_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8100"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}
