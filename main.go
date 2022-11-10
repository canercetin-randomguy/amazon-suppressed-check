package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//
// USED TO STORE TEAM NAMES, TO SORT OUT SEARCH RESULTS
//

type NBATeams []struct {
	TeamID       int    `json:"teamId"`
	Abbreviation string `json:"abbreviation"`
	TeamName     string `json:"teamName"`
	SimpleName   string `json:"simpleName"`
	Location     string `json:"location"`
}

type NFLTeams []struct {
	City string `json:"city"`
	Name string `json:"name"`
	Abr  string `json:"abr"`
	Conf string `json:"conf"`
	Div  string `json:"div"`
}

type NHLTeams []struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type MLBTeams struct {
	TeamAll struct {
		CopyRight    string `json:"copyRight"`
		QueryResults struct {
			Created   string `json:"created"`
			TotalSize string `json:"totalSize"`
			Row       []struct {
				PhoneNumber        string `json:"phone_number"`
				VenueName          string `json:"venue_name"`
				FranchiseCode      string `json:"franchise_code"`
				SportFull          string `json:"sport_full"`
				AllStarSw          string `json:"all_star_sw"`
				SportCode          string `json:"sport_code"`
				AddressCity        string `json:"address_city"`
				City               string `json:"city"`
				NameDisplayFull    string `json:"name_display_full"`
				SpringLeagueAbbrev string `json:"spring_league_abbrev"`
				TimeZoneAlt        string `json:"time_zone_alt"`
				SportID            string `json:"sport_id"`
				VenueID            string `json:"venue_id"`
				MlbOrgID           string `json:"mlb_org_id"`
				MlbOrg             string `json:"mlb_org"`
				LastYearOfPlay     string `json:"last_year_of_play"`
				LeagueFull         string `json:"league_full"`
				LeagueID           string `json:"league_id"`
				NameAbbrev         string `json:"name_abbrev"`
				AddressProvince    string `json:"address_province"`
				BisTeamCode        string `json:"bis_team_code"`
				League             string `json:"league"`
				SpringLeague       string `json:"spring_league"`
				BaseURL            string `json:"base_url"`
				AddressZip         string `json:"address_zip"`
				SportCodeDisplay   string `json:"sport_code_display"`
				MlbOrgShort        string `json:"mlb_org_short"`
				TimeZone           string `json:"time_zone"`
				AddressLine1       string `json:"address_line1"`
				MlbOrgBrief        string `json:"mlb_org_brief"`
				AddressLine2       string `json:"address_line2"`
				AddressLine3       string `json:"address_line3"`
				DivisionAbbrev     string `json:"division_abbrev"`
				SportAbbrev        string `json:"sport_abbrev"`
				NameDisplayShort   string `json:"name_display_short"`
				TeamID             string `json:"team_id"`
				ActiveSw           string `json:"active_sw"`
				AddressIntl        string `json:"address_intl"`
				State              string `json:"state"`
				AddressCountry     string `json:"address_country"`
				MlbOrgAbbrev       string `json:"mlb_org_abbrev"`
				Division           string `json:"division"`
				Name               string `json:"name"`
				TeamCode           string `json:"team_code"`
				SportCodeName      string `json:"sport_code_name"`
				WebsiteURL         string `json:"website_url"`
				FirstYearOfPlay    string `json:"first_year_of_play"`
				LeagueAbbrev       string `json:"league_abbrev"`
				NameDisplayLong    string `json:"name_display_long"`
				StoreURL           string `json:"store_url"`
				NameShort          string `json:"name_short"`
				AddressState       string `json:"address_state"`
				DivisionFull       string `json:"division_full"`
				SpringLeagueFull   string `json:"spring_league_full"`
				Address            string `json:"address"`
				NameDisplayBrief   string `json:"name_display_brief"`
				FileCode           string `json:"file_code"`
				DivisionID         string `json:"division_id"`
				SpringLeagueID     string `json:"spring_league_id"`
				VenueShort         string `json:"venue_short"`
			} `json:"row"`
		} `json:"queryResults"`
	} `json:"team_all"`
}

// USED FOR LOOKING UP TO PRODUCT IF EXISTS.
type ItemforLookup struct {
	Asin        string `json:"asin"`
	Identifiers []struct {
		MarketplaceID string `json:"marketplaceId"`
		Identifiers   []struct {
			IdentifierType string `json:"identifierType"`
			Identifier     string `json:"identifier"`
		} `json:"identifiers"`
	} `json:"identifiers"`
	ProductTypes []struct {
		MarketplaceID string `json:"marketplaceId"`
		ProductType   string `json:"productType"`
	} `json:"productTypes"`
	Summaries []struct {
		MarketplaceID        string `json:"marketplaceId"`
		AdultProduct         bool   `json:"adultProduct"`
		Autographed          bool   `json:"autographed"`
		Brand                string `json:"brand"`
		BrowseClassification struct {
			DisplayName      string `json:"displayName"`
			ClassificationID string `json:"classificationId"`
		} `json:"browseClassification"`
		Color                   string `json:"color"`
		ItemClassification      string `json:"itemClassification"`
		ItemName                string `json:"itemName"`
		Manufacturer            string `json:"manufacturer"`
		Memorabilia             bool   `json:"memorabilia"`
		ModelNumber             string `json:"modelNumber"`
		PackageQuantity         int    `json:"packageQuantity"`
		PartNumber              string `json:"partNumber"`
		Size                    string `json:"size"`
		Style                   string `json:"style"`
		TradeInEligible         bool   `json:"tradeInEligible"`
		WebsiteDisplayGroup     string `json:"websiteDisplayGroup"`
		WebsiteDisplayGroupName string `json:"websiteDisplayGroupName"`
	} `json:"summaries"`
}
type LookupPayload struct {
	NumberOfResults int             `json:"numberOfResults"`
	Items           []ItemforLookup `json:"items"`
}

type ListingsDetailPayloadIssuses struct {
	Code           string   `json:"code"`
	Message        string   `json:"message"`
	Severity       string   `json:"severity"`
	AttributeNames []string `json:"attributeNames"`
}
type ListingsDetailPayloadSummaries struct {
	MarketplaceID   string    `json:"marketplaceId"`
	Asin            string    `json:"asin"`
	ProductType     string    `json:"productType"`
	ConditionType   string    `json:"conditionType"`
	Status          []string  `json:"status"`
	ItemName        string    `json:"itemName"`
	CreatedDate     time.Time `json:"createdDate"`
	LastUpdatedDate time.Time `json:"lastUpdatedDate"`
	MainImage       struct {
		Link   string `json:"link"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	}
}
type ListingsDetailPayload struct {
	Sku       string                           `json:"sku"`
	Summaries []ListingsDetailPayloadSummaries `json:"summaries"`
	Issues    []ListingsDetailPayloadIssuses   `json:"issues"`
}

// USED FOR CREATING LISTINGS
type ConditionType struct {
	Value         string `json:"value"`
	MarketplaceID string `json:"marketplace_id"`
}
type MerchantSuggestedASIN struct {
	Value         string `json:"value"`
	MarketplaceID string `json:"marketplace_id"`
}
type Price struct {
	Value         string `json:"value"`
	MarketplaceID string `json:"marketplace_id"`
}
type PUTRequestData struct {
	ProductType  string `json:"productType"`
	Requirements string `json:"requirements"`
	Attributes   struct {
		Conditions []ConditionType         `json:"condition_type"`
		ASIN       []MerchantSuggestedASIN `json:"merchant_suggested_asin"`
		Offer      []Price                 `json:"purchasable_offer"`
	} `json:"attributes"`
}
type ListingRestriction struct {
	Restrictions []struct {
		MarketplaceID string `json:"marketplaceId"`
		ConditionType string `json:"conditionType"`
		Reasons       []struct {
			ReasonCode string        `json:"reasonCode"`
			Message    string        `json:"message"`
			Links      []interface{} `json:"links"`
		} `json:"restrictions"`
	}
}
type ProductDetails struct {
	NumberOfResults int `json:"numberOfResults"`
	Items           []struct {
		Asin      string `json:"asin"`
		Summaries []struct {
			MarketplaceID        string `json:"marketplaceId"`
			AdultProduct         bool   `json:"adultProduct"`
			Autographed          bool   `json:"autographed"`
			Brand                string `json:"brand"`
			BrowseClassification struct {
				DisplayName      string `json:"displayName"`
				ClassificationID string `json:"classificationId"`
			} `json:"browseClassification"`
			Color                   string `json:"color"`
			ItemClassification      string `json:"itemClassification"`
			ItemName                string `json:"itemName"`
			Manufacturer            string `json:"manufacturer"`
			Memorabilia             bool   `json:"memorabilia"`
			ModelNumber             string `json:"modelNumber"`
			PackageQuantity         int    `json:"packageQuantity"`
			PartNumber              string `json:"partNumber"`
			Size                    string `json:"size"`
			Style                   string `json:"style"`
			TradeInEligible         bool   `json:"tradeInEligible"`
			WebsiteDisplayGroup     string `json:"websiteDisplayGroup"`
			WebsiteDisplayGroupName string `json:"websiteDisplayGroupName"`
		} `json:"summaries"`
	} `json:"items"`
}

type ClientCredentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
type AppCredentials struct {
	RefreshToken    string
	ClientID        string
	IAMClientID     string
	IAMSecretAccess string
	SellerID        string
	ClientSecret    string
}

func GetCredentials() AppCredentials {
	viper.SetConfigName("credentials")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	var tempCred AppCredentials
	tempCred.RefreshToken = viper.Get("REFRESH_TOKEN").(string)
	tempCred.ClientID = viper.Get("CLIENT_ID").(string)
	tempCred.ClientSecret = viper.Get("CLIENT_SECRET").(string)
	tempCred.SellerID = viper.Get("SELLER_ID").(string)
	tempCred.IAMClientID = viper.Get("IAM_CLIENT_ID").(string)
	tempCred.IAMSecretAccess = viper.Get("IAM_SECRET_ACCESS").(string)
	return tempCred
}

// ClientCredentialGenerator returns Access Token and other shit.
func ClientCredentialGenerator() ClientCredentials {
	var AppCred = GetCredentials()
	LWAAuth, err := url.Parse("https://api.amazon.com/auth/o2/token" + "?" + "grant_type=refresh_token" +
		"&" + "refresh_token=" + AppCred.RefreshToken +
		"&" + "client_id=" + AppCred.ClientID +
		"&" + "client_secret=" + AppCred.ClientSecret)
	if err != nil {
		panic(err)
	}
	//
	// https://developer-docs.amazon.com/sp-api/docs/connecting-to-the-selling-partner-api#step-1-request-a-login-with-amazon-access-token
	//
	authS1Req, err := http.NewRequest("POST", LWAAuth.String(), nil)
	if err != nil {
		panic(err)
	}
	authS1Req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	resp, err := http.DefaultClient.Do(authS1Req)
	if err != nil {
		panic(err)
	}
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var ClientCred ClientCredentials
	err = json.Unmarshal(respData, &ClientCred)
	if err != nil {
		panic(err)
	}
	return ClientCred
}

// APIURIConstruct https://developer-docs.amazon.com/sp-api/docs/connecting-to-the-selling-partner-api#step-2-construct-a-selling-partner-api-uri,
//
// endpoint =  https://sellingpartnerapi-na.amazon.com
//
// requestPath = /listings/2021-08-01/items/AXXXXXXXXXXXX/50-TS3D-QEPT, or any other usage.
//
// marketplace = marketplaceIds = 	ATVPDKIKX0DER for US, https://developer-docs.amazon.com/sp-api/docs/marketplace-ids for rest.
//
// httpMethod = GET, POST, PUT, DELETE, etc.
//
// Feed both of the Client and App credentials.
//
// Feed return value of GetTime() function.
//
// and fire it up.
func APIURIConstruct(operation string, endpoint string, requestPath string, parameters string, marketplace string, httpMethod string, ClientCred ClientCredentials, AppCred AppCredentials, body string) ([]byte, bool) {
	authS2ReqURL, err := url.Parse(endpoint + requestPath + "?" + "marketplaceIds=" + marketplace + "&" + parameters)
	if err != nil {
		panic(err)
	}
	authS2ReqURLREQ, err := http.NewRequest(httpMethod, authS2ReqURL.String(), nil)
	if err != nil {
		panic(err)
	}
	authS2ReqURLREQ.Header.Set("x-amz-access-token", ClientCred.AccessToken)
	authS2ReqURLREQ.Header.Set("Content-Type", "application/json")
	authS2ReqURLREQ.Header.Set("x-amz-date", GetTime())
	authS2ReqURLREQ.Header.Set("host", "sellingpartnerapi-na.amazon.com")
	authS2ReqURLREQ.Header.Set("user-agent", "XXXXXXXXXX/1.0 (Language=Go; Platform=Windows)")
	// put body to request
	authS2ReqURLREQ.Body = io.NopCloser(strings.NewReader(body))
	// declare a new signer signer := v4.NewSigner(&credentials.Credentials{}) but only fill the third parameter.
	signer := v4.NewSigner()
	if body == "" {
		err := signer.SignHTTP(context.Background(), aws.Credentials{
			AccessKeyID:     AppCred.IAMClientID,
			SecretAccessKey: AppCred.IAMSecretAccess,
		}, authS2ReqURLREQ, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "execute-api",
			"us-east-1", time.Now().UTC())
		if err != nil {
			panic(err)
		}
	} else {
		bodySHA256 := sha256.Sum256([]byte(body))
		err := signer.SignHTTP(context.Background(), aws.Credentials{
			AccessKeyID:     AppCred.IAMClientID,
			SecretAccessKey: AppCred.IAMSecretAccess,
		}, authS2ReqURLREQ, fmt.Sprintf("%x", bodySHA256), "execute-api",
			"us-east-1", time.Now().UTC())
		if err != nil {
			panic(err)
		}
	}
	resp, err := http.DefaultClient.Do(authS2ReqURLREQ)
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// TODO: change this to select case.
	switch operation {
	case "lookup":
		var tempItemData ProductDetails
		err = json.Unmarshal(respData, &tempItemData)
		if err != nil {
			panic(err)
		}
		if tempItemData.NumberOfResults == 0 {
			return respData, false
		} else {
			return respData, true
		}
	case "restriction":
		var tempRestrictionData ListingRestriction
		err = json.Unmarshal(respData, &tempRestrictionData)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(tempRestrictionData.Restrictions); i++ {
			for j := 0; j < len(tempRestrictionData.Restrictions[i].Reasons); j++ {
				if tempRestrictionData.Restrictions[i].Reasons[j].ReasonCode == "NOT_ELIGIBLE" {
					return []byte(""), false
				} else if tempRestrictionData.Restrictions[i].Reasons[j].ReasonCode != "NOT_ELIGIBLE" && tempRestrictionData.Restrictions[i].ConditionType == "new_new" {
					return []byte(""), true
				}
			}
		}
	case "productType":
		if string(respData) == "" {
			return respData, false
		} else {
			return respData, true
		}
	case "createListing":
		fmt.Println(string(respData))
		if strings.Contains(string(respData), `"status":"INVALID"`) {
			return respData, false
		} else {
			return respData, true
		}
		return respData, true
	case "getListing":
		return respData, true
	}
	return []byte(""), true
}
func GetTime() string {
	return time.Now().UTC().Format("20060102T150405Z")
}

// ReturnListingDetails performs getListingDetails operation and returns a ListingsDetailPayload struct.
func ReturnListingDetails(tempCred AppCredentials, ClientCreds ClientCredentials, sku string) (ListingsDetailPayload, bool) {
	var getListingsAbsolutePath = fmt.Sprintf("/listings/2021-08-01/items/%s/%s", tempCred.SellerID, sku)
	parameters := fmt.Sprintf("marketplaceIds=%s&includedData=%s", "ATVPDKIKX0DER", "summaries,issues")
	respData, success := APIURIConstruct("getListing", "https://sellingpartnerapi-na.amazon.com", getListingsAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCreds, tempCred, "")
	var tempListingData ListingsDetailPayload
	err := json.Unmarshal(respData, &tempListingData)
	if err != nil {
		panic(err)
	}
	return tempListingData, success
}
func GetSKU(cnt *int) string {
	randomSKU := "MCS-" + strconv.Itoa(*cnt)
	*cnt++
	return randomSKU
}
func GetDataRows() []*xlsx.Row {
	dataFile, err := xlsx.OpenFile("temp.xlsx")
	if err != nil {
		panic(err)
	}
	sh, ok := dataFile.Sheet["Sayfa1"]
	if !ok {
		panic("sheet not found")
	}
	return sh.Rows
}
func main() {
	var skuList []string
	var tryCount = 0
	var detailsExt bool
	var cnt = 32
	// on multiple result queries, these teams are needed to sort out what we need.
	var NBATeamList NBATeams
	// open nba.json and unmarshal it into NBATeamList
	jsonFile, err := os.Open("nba.json")
	if err != nil {
		panic(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &NBATeamList)
	if err != nil {
		panic(err)
	}
	var NCAATeamList []string
	// open ncaa.txt and append each line to NCAATeamList
	ncaaFile, err := os.Open("ncaa.txt")
	if err != nil {
		panic(err)
	}
	defer func(ncaaFile *os.File) {
		err := ncaaFile.Close()
		if err != nil {
			panic(err)
		}
	}(ncaaFile)
	scanner := bufio.NewScanner(ncaaFile)
	for scanner.Scan() {
		NCAATeamList = append(NCAATeamList, scanner.Text())
	}
	var NHLTeamList NHLTeams
	// open nhl.json and unmarshal it into NHLTeamList
	nhlFile, err := os.Open("nhl.json")
	if err != nil {
		panic(err)
	}
	defer func(nhlFile *os.File) {
		err := nhlFile.Close()
		if err != nil {
			panic(err)
		}
	}(nhlFile)
	nhlByteValue, err := io.ReadAll(nhlFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(nhlByteValue, &NHLTeamList)
	if err != nil {
		panic(err)
	}
	var NFLTeamList NFLTeams
	// open nfl.json and unmarshal it into NFLTeamList
	nflFile, err := os.Open("nfl.json")
	if err != nil {
		panic(err)
	}
	defer func(nflFile *os.File) {
		err := nflFile.Close()
		if err != nil {
			panic(err)
		}
	}(nflFile)
	nflByteValue, err := io.ReadAll(nflFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(nflByteValue, &NFLTeamList)
	if err != nil {
		panic(err)
	}
	var MLBTeamList MLBTeams
	// open mlb.json and unmarshal it into MLBTeamList
	mlbFile, err := os.Open("mlb.json")
	if err != nil {
		panic(err)
	}
	defer func(mlbFile *os.File) {
		err := mlbFile.Close()
		if err != nil {
			panic(err)
		}
	}(mlbFile)
	mlbByteValue, err := io.ReadAll(mlbFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(mlbByteValue, &MLBTeamList)
	if err != nil {
		panic(err)
	}

	var ClientCreds = ClientCredentialGenerator()
	// create a timer that will fire when the credentials expire
	ExpirationTimer := time.NewTimer(time.Duration(ClientCreds.ExpiresIn) * time.Second)
	var logFile, _ = os.Create("run.log")
	logger := logrus.New()
	logger.SetOutput(logFile)
	var teamPrefix string
	var found bool
	var ASINValue string
	var status string
	var teamSuffix string
	var tempCred = GetCredentials()
	// create a output.csv
	outputFile, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			panic(err)
		}
	}(outputFile)
	tempColumnNames := []string{
		"Product Title", "ASIN", "Sellable",
	}
	csvReader := csv.NewWriter(outputFile)
	csvReader.Write(tempColumnNames)
	csvReader.Flush()
	rows := GetDataRows()
	for i := range rows {
		tempSKU := GetSKU(&cnt)
		if i == 0 {
			continue
		} else {
			select {
			case <-ExpirationTimer.C:
				ClientCreds = ClientCredentialGenerator()
				ExpirationTimer.Reset(time.Duration(ClientCreds.ExpiresIn) * time.Second)
				logger.Infoln("Client credentials refreshed.")
			default:
				//
				// https://developer-docs.amazon.com/sp-api/docs/building-listings-management-workflows-guide#list-an-offer-for-an-item-that-already-exists-in-the-amazon-catalog
				//
				// LISTING CREATION PROCESS
				//
				//
				// Call the searchCatalogItems operation to search for existing items in the Amazon catalog by product identifiers (UPC, EAN, etc.) or keywords.
				var tempASIN = rows[i].Cells[2].Value
				var productName = rows[i].Cells[1].Value
				logger.Infoln("Searching for product with UPC: ", tempASIN)
				logger.Infoln("With product name: ", productName)
				// search for bucks in the structs declared earlier
				if rows[i].Cells[3].Value == "NBA" {
					for _, team := range NBATeamList {
						if strings.Contains(strings.ToLower(productName), strings.ToLower(team.SimpleName)) ||
							strings.Contains(strings.ToLower(productName), strings.ToLower(team.Location)) {
							teamPrefix = team.Location
							teamSuffix = team.SimpleName
						}
					}
				}
				if teamPrefix == "" && rows[i].Cells[3].Value == "NHL" {
					for _, team := range NHLTeamList {
						if strings.Contains(strings.ToLower(productName), strings.ToLower(team.Name)) ||
							strings.Contains(strings.ToLower(productName), strings.ToLower(team.City)) {
							teamPrefix = team.City
							teamSuffix = team.Name
						}
					}
				}
				if teamPrefix == "" && rows[i].Cells[3].Value == "NFL" {
					for _, team := range NFLTeamList {
						if strings.Contains(strings.ToLower(productName), strings.ToLower(team.Name)) ||
							strings.Contains(strings.ToLower(productName), strings.ToLower(team.City)) {
							teamPrefix = team.City
							teamSuffix = team.Name
						}
					}
				}
				if teamPrefix == "" && rows[i].Cells[3].Value == "MLB" {
					for _, team := range MLBTeamList.TeamAll.QueryResults.Row {
						if strings.Contains(strings.ToLower(productName), strings.ToLower(team.Name)) ||
							strings.Contains(strings.ToLower(productName), strings.ToLower(team.City)) {
							teamPrefix = team.City
							teamSuffix = team.Name
						}
					}
				}
				if teamPrefix == "" && rows[i].Cells[3].Value == "NCAA" {
					// open ncaa.txt and append each line to NCAATeamList
					ncaaFile, err := os.Open("ncaa.txt")
					if err != nil {
						panic(err)
					}
					defer func(ncaaFile *os.File) {
						err := ncaaFile.Close()
						if err != nil {
							panic(err)
						}
					}(ncaaFile)
					scanner := bufio.NewScanner(ncaaFile)
					for scanner.Scan() {
						tempText := strings.Split(scanner.Text(), " ")
						if strings.Contains(strings.ToLower(productName), strings.ToLower(tempText[0])) {
							teamPrefix = tempText[0]
							for i := 1; i < len(tempText); i++ {
								teamSuffix = teamSuffix + tempText[i]
							}
						}
					}
				}
				var searchCatalogAbsolutePath = fmt.Sprintf("/catalog/2022-04-01/items")
				identifier := "UPC"
				var parameters = fmt.Sprintf("identifiers=%s&identifiersType=%s&includedData=%s", tempASIN, identifier, "productTypes,summaries,identifiers")
				productData, exists := APIURIConstruct("lookup", "https://sellingpartnerapi-na.amazon.com", searchCatalogAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCreds, GetCredentials(), "")
				var productLookupData LookupPayload
				err := json.Unmarshal(productData, &productLookupData)
				if err != nil {
					panic(err)
				}
				fmt.Println(productLookupData.NumberOfResults)
				if productLookupData.NumberOfResults > 1 {
					for i := range productLookupData.Items {
						fmt.Println("Item> ", i)
						for j := range productLookupData.Items[i].Summaries {
							fmt.Println(productLookupData.Items[i].Summaries[j].ItemName, strings.ToLower(teamPrefix), strings.ToLower(teamSuffix))
							if strings.Contains(strings.ToLower(productLookupData.Items[i].Summaries[j].ItemName), strings.ToLower(teamPrefix)) ||
								strings.Contains(strings.ToLower(productLookupData.Items[i].Summaries[j].ItemName), strings.ToLower(teamSuffix)) {
								fmt.Println("Found a match!")
								fmt.Println("City> ", teamPrefix)
								fmt.Println("Team> ", teamSuffix)
								ASINValue = productLookupData.Items[i].Asin
								found = true
								break
							}
						}
					}
					if found == false {
						status = "NOT_FOUND"
					}
				} else if productLookupData.NumberOfResults == 0 {
					logger.Errorln("No results found for UPC: ", tempASIN)
					status = "NOT_FOUND"
					tempRowData := []string{
						productName, tempASIN, status,
					}
					csvReader.Write(tempRowData)
					csvReader.Flush()
				} else {
					if len(productLookupData.Items) > 0 {
						if len(productLookupData.Items[0].Asin) > 0 {
							// save productLookupData to a file
							f, err := os.Create("productLookupData.json")
							if err != nil {
								panic(err)
							}
							defer func(f *os.File) {
								err := f.Close()
								if err != nil {
									panic(err)
								}
							}(f)
							_, err = f.WriteString(string(productData))
							if err != nil {
								return
							}
							ASINValue = productLookupData.Items[0].Asin
							found = true
						} else {
							found = false
							break
						}
					} else {
						found = false
					}
				}
				fmt.Println("ASINValue: ", ASINValue)
				if productLookupData.NumberOfResults == 1 || found {
					if exists {
						// If product exists in Amazon:
						// Call the getListingsRestrictions operation with the ASIN identifier to retrieve any eligibility requirements that must be met before listing an item in the applicable condition.
						var getListingsRestrictionsAbsolutePath = fmt.Sprintf("/listings/2021-08-01/restrictions")
						parameters := fmt.Sprintf("asin=%s&sellerId=%s&marketplaceIDs=%s", tempASIN, tempCred.SellerID, "ATVPDKIKX0DER")
						// If we are able to sell it in Amazon:
						// Call the getProductType operation to retrieve the product type for the item.
						_, ableToSell := APIURIConstruct("restriction", "https://sellingpartnerapi-na.amazon.com", getListingsRestrictionsAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCreds, tempCred, "")
						if ableToSell {
							// If request is successful:
							// Unmarshal the ProductTypeData struct
							if err != nil {
								panic(err)
							}
							// Create a JSON file for submitting product
							var tempProductPUT PUTRequestData
							tempProductPUT.ProductType = "PRODUCT"
							tempProductPUT.Requirements = "LISTING_OFFER_ONLY"
							tempProductPUT.Attributes.ASIN = append(tempProductPUT.Attributes.ASIN, MerchantSuggestedASIN{
								Value:         ASINValue,
								MarketplaceID: "ATVPDKIKX0DER",
							})
							tempProductPUT.Attributes.Conditions = append(tempProductPUT.Attributes.Conditions, ConditionType{
								Value:         "new_new",
								MarketplaceID: "ATVPDKIKX0DER",
							})
							tempProductPUT.Attributes.Offer = append(tempProductPUT.Attributes.Offer, Price{
								Value:         "0.0",
								MarketplaceID: "ATVPDKIKX0DER",
							})
							// convert the struct to json
							tempProductJSON, err := json.Marshal(tempProductPUT)
							if err != nil {
								panic(err)
							}
							// create the listing
							var createListingAbsolutePath = fmt.Sprintf("/listings/2021-08-01/items/%s/%s", tempCred.SellerID, tempSKU)
							parameters := fmt.Sprintf("&marketplaceIds=%s", "ATVPDKIKX0DER")
							_, _ = APIURIConstruct("createListing", "https://sellingpartnerapi-na.amazon.com", createListingAbsolutePath, parameters, "ATVPDKIKX0DER", "PUT", ClientCreds, tempCred, string(tempProductJSON))
							//
							//
							// LISTING CREATION PROCESS IS DONE, SLEEP A BIT TO WAIT FOR THE LISTING TO BE CREATED
							//
							//
						} else {
							logger.Errorln("Listing creation failed.")
							break
						}
					} else {
						logger.Errorln("Unable to sell: ", tempASIN, " in Amazon.")
						break
					}
				} else {
					logger.Errorln("ASIN: ", tempASIN, " does not exist in Amazon.")
					break
				}
			}
			skuList = append(skuList, tempSKU)
			time.Sleep(500 * time.Millisecond)
		}
	}
	for SKU := range skuList {
		listingsDetail, _ := ReturnListingDetails(tempCred, ClientCreds, skuList[SKU])
		// marshal listingsDetail to json
		listingsDetailJSON, _ := json.Marshal(listingsDetail)
		if len(listingsDetail.Summaries) > 0 {
			for j := range listingsDetail.Summaries[0].Status {
				if listingsDetail.Summaries[0].Status[j] == "DISCOVERABLE" {
					logger.Infoln("Listing is created.")
					status = "DISCOVERABLE"
					tempRowData := []string{
						listingsDetail.Summaries[0].ItemName, listingsDetail.Summaries[0].Asin, status, string(listingsDetailJSON),
					}
					csvReader.Write(tempRowData)
					csvReader.Flush()
				} else {
					logger.Infoln("Listing is search suppressed.")
					status = "SEARCH_SUPPRESSED"
					tempRowData := []string{
						listingsDetail.Summaries[0].ItemName, listingsDetail.Summaries[0].Asin, status, string(listingsDetailJSON),
					}
					csvReader.Write(tempRowData)
					csvReader.Flush()
				}
			}
		} else {
			for {
				listingsDetail, _ = ReturnListingDetails(tempCred, ClientCreds, skuList[SKU])
				listingsDetailJSON, _ := json.Marshal(listingsDetail)
				if len(listingsDetail.Summaries) > 0 && (len(listingsDetail.Issues) == 0) {
					fmt.Println("SUMMARY EXISTS")
					for j := range listingsDetail.Summaries[0].Status {
						fmt.Println(listingsDetail.Summaries[0].Status[j])
						if listingsDetail.Summaries[0].Status[j] == "DISCOVERABLE" {
							logger.Infoln("Listing is created.")
							detailsExt = true
							status = "DISCOVERABLE"
						} else if listingsDetail.Summaries[0].Status[j] == "" {
							logger.Infoln("Listing is search suppressed.")
							detailsExt = true
							status = "SEARCH_SUPPRESSED"
						}
					}
				}
				if len(listingsDetail.Issues) > 0 && len(listingsDetail.Summaries) == 0 {
					logger.Errorln("Errors: ")
					for j := range listingsDetail.Issues {
						detailsExt = true
						logger.Error(listingsDetail.Issues[j].Code)
						logger.Errorln(listingsDetail.Issues[j].Message)
						logger.Errorln("Listing is search suppressed.")
						status = "DISCOVERABLE_WITH_ERRORS"
					}
				}
				if len(listingsDetail.Issues) > 0 && len(listingsDetail.Summaries) > 0 {
					for j := range listingsDetail.Summaries[0].Status {
						fmt.Println(listingsDetail.Summaries[0].Status[j])
						if listingsDetail.Summaries[0].Status[j] == "DISCOVERABLE" {
							logger.Infoln("Listing is created.")
							detailsExt = true
							status = "DISCOVERABLE"
						} else {
							logger.Infoln("Listing is search suppressed.")
							detailsExt = true
							status = "SEARCH_SUPPRESSED"
						}
					}
				}
				if len(listingsDetail.Issues) == 0 && len(listingsDetail.Summaries) > 1 {
					for j := range listingsDetail.Summaries[0].Status {
						fmt.Println(listingsDetail.Summaries[0].Status[j])
						if listingsDetail.Summaries[0].Status[j] == "DISCOVERABLE" {
							logger.Infoln("Listing is created.")
							detailsExt = true
							status = "DISCOVERABLE"
						} else {
							logger.Infoln("Listing is search suppressed.")
							detailsExt = true
							status = "SEARCH_SUPPRESSED"
						}
					}
				}
				if detailsExt {
					tempRowData := []string{
						listingsDetail.Summaries[0].ItemName, listingsDetail.Summaries[0].Asin, status, string(listingsDetailJSON),
					}
					csvReader.Write(tempRowData)
					csvReader.Flush()
					break
				}
				tryCount++
				if tryCount > 240 {
					status = "TIMEOUT"
					break
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
		detailsExt = false
		tryCount = 0
	}
}
