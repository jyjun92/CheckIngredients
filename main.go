package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/*
TO DO:
- add radio button in UI: hair, face, body, lips
- get lots of info from director pi and add facts page maybe?
- for ulta, var x = document.getElementsByClassName("ProductDetail__ingredients"), x[0].textContent gives ingredients
*/

var numOfFlags int = 0

func main() {
	var baseURL string
	fmt.Println("Enter url here: ")
	fmt.Scanf("%s", &baseURL)

	getIngredients(baseURL)
}

func getIngredients(baseURL string) {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchIngredients := doc.Find("#tabpanel2") //tabpanel2
	//css-pz80c5
	var contents string = ""
	children := searchIngredients.Children()
	contents = children.Text()
	res1 := strings.Split(contents, "  ") //get rid of explanations
	if len(res1) > 1 {
		contents = res1[1]
	}
	res2 := strings.Split(contents, ", ")

	res2[len(res2)-1] = strings.Split(res2[len(res2)-1], ".")[0] // removing . from last item
	for i := 0; i < len(res2); i++ {
		//fmt.Println("Ingredient Name: ", res2[i], "  INCI Name: ", getINCI(res2[i]))
		go getINCI(res2[i])
	}
	time.Sleep(time.Second * 20)
	defer fmt.Println(numOfFlags, " harmful ingredient(s) found in this product.")

}

type dataDTO struct {
	Nhits      int
	Parameters struct {
		Dataset  string
		Timezone string
		Q        string
		Rows     int
		Format   string
		Facet    []string
	}
	Records []struct {
		Datasetid string
		Recordid  string
		Fields    struct {
			Inci_name                    string
			Function                     string
			Update_date                  string
			Cosing_ref_no                string
			Chem_iupac_name_descpription string
			Restriction                  string
			Cas_no                       string
			Ec_no                        string
		}
		Record_timestamp string
	}
	Facet_groups []struct {
		Faucets []struct {
			Count int
			Path  string
			State string
			Name  string
		}
		Name string
	}
}

func getINCI(ingredient string) {
	var baseURL string = "https://public.opendatasoft.com/api/records/1.0/search/?dataset=cosmetic-ingredient-database-ingredients-and-fragrance-inventory&q="
	var endURL string = "&facet=update_date&facet=restriction&facet=function"
	//url should be baseURL + ingredient + endURL
	var completeURL = baseURL + url.QueryEscape(ingredient) + endURL

	res, err := http.Get(completeURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, _ := ioutil.ReadAll(res.Body)

	var rec dataDTO
	err2 := json.Unmarshal(doc, &rec)

	checkErr(err2)

	if len(rec.Records) > 0 {
		flag := checkIfHarmful(rec.Records[0].Fields.Inci_name)
		if flag {
			fmt.Println(rec.Records[0].Fields.Inci_name, " ----------- CANCEROUS!!!!!!!!")
			numOfFlags++
		} else {
			fmt.Println(rec.Records[0].Fields.Inci_name)
		}
	} else {
		var splitIngredientName = strings.Split(ingredient, "(")
		if len(splitIngredientName) > 0 {
			getINCI(splitIngredientName[0])
		} else {
			fmt.Println(ingredient, " NOT found!")
		}
	}

}

func checkIfHarmful(ingredient string) bool {
	result := false
	for i := 0; i < len(avoidChemicals); i++ {
		if strings.Contains(ingredient, avoidChemicals[i]) {
			return true
		}
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with StatusCode:", res.StatusCode)
	}
}

var avoidChemicals = [...]string{"AVOBENZONE", "ISOPROPYL ALCOHOL", "SLS", "SODIUM LAURYL SULFATE",
	"SLES", "SODIUM LAURETH SULFATE", "TEA", "TRIETHANOLAMINE", "PEG", "POLYETHYLENE GLYCOL",
	"COLOR", "ISOPROPYL METHYPHENOL", "SORBIC ACID", "DHT", "PARABEN",
	"TRICLOSAN", "BHA", "BUTYL HYDROXY ANISOLE", "BUTYLATED HYDROXYANISOLE", "OXYBENZONE", "UREA",
	"MINERAL OIL", "THYMOL", "TRIISOPROPANOLAMINE", "FRAGRANCE", "PHENOXY ETHANOL", "DEA", "DIETHANOLAMINE",
	"DMP", "FORMALDEHYDE", "PERFUME", "PARFUM", "CETEARETH", "SILOXANE", "SILICA", "ARSENIC", "TALC",
	"DIMETHICONE", " TOLUENE", "PHTHALATE"}
