package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"math/rand"

	consulapi "github.com/hashicorp/consul/api"
)

type User struct {
	ID       uint64    `json:"id"`
	Username string    `json:"username"`
	Products []product `json:"products"`
	URL      string    `json:"url"`
}

type product struct {
	ID    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// dang ky service voi consul
func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "user-service"   // id
	registration.Name = "user-service" // name cua service
	address := hostname()
	registration.Address = address // ip
	p, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = p // port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", address, p) // healthcheck
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

// lookup service tu consul
func lookupServiceWithConsul(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
	return "", err
	}
	passingOnly := true
	serviceEntries, _, err := consul.Health().Service(serviceName, "", passingOnly, nil)
	if len(serviceEntries) == 0 && err == nil {
	return "", fmt.Errorf("service ( %s ) was not found", serviceName)
	}
		if err != nil {
		return "", err
	}
			
	instanceIdx := rand.Intn(len(serviceEntries))
	
	address := serviceEntries[instanceIdx].Service.Address
	port := serviceEntries[instanceIdx].Service.Port

	return fmt.Sprintf("http://%s:%v", address, port), nil
}

// check key-value moi tu consuk K-V
func Configuration(key string) (bool, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return false, err
	}
	kvpair, _, err := consul.KV().Get(key, nil)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return false, err
	}
	if kvpair.Value == nil {
		fmt.Fprintf(w, "Configuration empty")
		return false, nil
	}
	//val := string(kvpair.Value)
	//fmt.Fprintf(w, "%s", val)
	if kvpair.Value == "enable"{
		return true, nil
	}
}

func main() {
	registerServiceWithConsul()

	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/user-products", UserProduct)

	fmt.Printf("user service is up on port: %s", port())

	http.ListenAndServe(port(), nil)
}

// web handler "/healthcheck"
func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `user service is good`)
}

// web hander "/user-products"
func UserProduct(w http.ResponseWriter, r *http.Request) {
	p := []product{}

	url, err := lookupServiceWithConsul("product-service")

	new_product, _ := Configuration("new-product")

	fmt.Println("URL: ", url)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	client := &http.Client{}

	if new_product == false {
		resp, err := client.Get(url + "/products")
		u := User{
			ID:       1,
			Username: "user1@gmail.com",
		}
	} else {
		u := User{
			ID:       2,
			Username: "user2@gmail.com",
		}
		resp, err := client.Get(url + "/new-products")
	}

	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	u.Products = p
	u.URL = url
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&u)
}

// tra ve port tu bien moi truong "USER_SERVICE_PORT" | 8080
func port() string {
	p := os.Getenv("USER_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8080"
	}
	return fmt.Sprintf(":%s", p)
}

// tra ve hostname
func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}
