package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// ProduceItem Struct
type ProduceItem struct {
	ProduceCode string  `json:"produce_code"`
	Name        string  `json:"name"`
	UnitPrice   float64 `json:"unit_price"`
}

// Starting Array of Produce Data
var produce = []ProduceItem{
	{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: 3.46},
	{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: 2.99},
	{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: 0.79},
	{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: 3.59},
}

// getAllProduce responds with the list of all Produce as JSON.
func getAllProduce(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, produce)
}

//Produce Code Collision Detection
func produceCheck(itemCode string) (present bool) {
	present = false
	for _, item := range produce {
		if item.ProduceCode == itemCode {
			present = true
			return
		}
	}
	return
}

//RegEx Match to make sure Produce Codes are formatted correctly
func checkCode(code string) (match bool) {
	re := regexp.MustCompile(`([a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9]-[a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9][a-zA-Z0-9])`)
	match = re.MatchString(code)
	if match == false {
		fmt.Printf("The following code is formatted wrong: %v\n", code)
	}
	return
}

func checkName(name string) (match bool) {
	re := regexp.MustCompile(`(^[A-Za-z0-9]*)$`)
	match = re.MatchString(name)
	if match == false {
		fmt.Printf("The following name is formatted wrong: %v\n", name)
	}
	return
}

//RegEx Match to make sure Prices are formatted correctly
func checkPrice(price float64) (match bool) {
	decimal := regexp.MustCompile(`(^[0-9]*\.[0-9][0-9])$`)
	decMatch := decimal.MatchString(fmt.Sprintf("%f", price))
	wholeNum := regexp.MustCompile(`(^[0-9]*)$`)
	wholeMatch := wholeNum.MatchString(fmt.Sprintf("%f", price))
	if decMatch == false && wholeMatch == false {
		fmt.Printf("The following price is formatted wrong: %v\n", price)
	}
	return
}

//Checks to make sure a Produce Item doesn't have empty values
func produceIntegrity(item ProduceItem) (correct bool) {
	correct = true
	if item.ProduceCode == "" || item.Name == " " || item.Name == "" || item.UnitPrice == 0 {
		correct = false
		fmt.Printf("The produce you are trying to upload is incomplete\n")
		return
	}
	return
}

//Iterates through produceItem array to check all the produce codes
func checkItems(items []ProduceItem) (format bool) {
	format = true
	for _, item := range items {
		if checkCode(item.ProduceCode) == false || checkName(item.Name) == false || produceIntegrity(item) == false || checkPrice(item.UnitPrice) == false {
			format = false
			return
		}
	}
	return
}

//Handle the POST Calls for adding single produce or more
//Must be an array of Produce Objects either 1 or more
func postProduce(c *gin.Context) {

	var newProduce []ProduceItem

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = json.Unmarshal(body, &newProduce)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	if checkItems(newProduce) == true {
		//Counter for number of produce that exist in the Mock DB
		existingProduce := 0

		for _, item := range newProduce {
			if produceCheck(item.ProduceCode) == false {
				produce = append(produce, item)
			} else {
				existingProduce = existingProduce + 1
				fmt.Printf("The following item already exists: %v\n", item)
			}
		}

		//If the payload only contains items that already exist
		//Return a conflict error
		//Else update DB with items from payload that don't conflict
		if existingProduce == len(newProduce) {
			c.String(http.StatusConflict, "All uploaded items already exist. Please change payload and try again")
		} else {
			c.IndentedJSON(http.StatusCreated, produce)
		}
	} else {
		c.String(http.StatusBadRequest, "Either the produce codes or unit prices are incorrectly formatted")
	}
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getProduceByID(c *gin.Context) {
	produce_code := c.Param("produce_code")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, p := range produce {
		if p.ProduceCode == produce_code {
			c.IndentedJSON(http.StatusOK, p)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "produce not found"})
}

func main() {
	router := gin.Default()
	router.GET("/produce", getAllProduce)
	router.GET("/produce/:produce_code", getProduceByID)
	router.POST("/produce", postProduce)

	router.Run("localhost:8080")
}
