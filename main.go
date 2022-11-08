package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

// USED FOR LOOKING UP TO PRODUCT IF EXISTS.
type TeamNameforLookup struct {
	LanguageTag   string `json:"language_tag"`
	Value         string `json:"value"`
	MarketplaceID string `json:"marketplace_id"`
}
type ItemforLookup struct {
	Asin       string `json:"asin"`
	Attributes struct {
		SupplierDeclaredMaterialRegulation []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"supplier_declared_material_regulation"`
		IncludedComponents []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"included_components"`
		ModelName []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"model_name"`
		BulletPoint []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"bullet_point"`
		Brand []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"brand"`
		Closure []struct {
			Type []struct {
				LanguageTag string `json:"language_tag"`
				Value       string `json:"value"`
			} `json:"type"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"closure"`
		GenericKeyword []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"generic_keyword"`
		ImportDesignation []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"import_designation"`
		CpsiaCautionaryStatement []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"cpsia_cautionary_statement"`
		ExternallyAssignedProductIdentifier []struct {
			Value         string `json:"value"`
			Type          string `json:"type"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"externally_assigned_product_identifier"`
		TeamName        []TeamNameforLookup `json:"team_name"`
		ItemTypeKeyword []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"item_type_keyword"`
		NumberOfItems []struct {
			Value         int    `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"number_of_items"`
		ItemPackageDimensions []struct {
			Length struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"length"`
			Width struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"width"`
			Height struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"height"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"item_package_dimensions"`
		Size []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"size"`
		PartNumber []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"part_number"`
		Style []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"style"`
		Color []struct {
			LanguageTag        string   `json:"language_tag"`
			StandardizedValues []string `json:"standardized_values"`
			Value              string   `json:"value"`
			MarketplaceID      string   `json:"marketplace_id"`
		} `json:"color"`
		ItemTypeName []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"item_type_name"`
		ItemPackageWeight []struct {
			Unit          string  `json:"unit"`
			Value         float64 `json:"value"`
			MarketplaceID string  `json:"marketplace_id"`
		} `json:"item_package_weight"`
		ItemShape []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"item_shape"`
		ModelNumber []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"model_number"`
		Department []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"department"`
		SupplierDeclaredDgHzRegulation []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"supplier_declared_dg_hz_regulation"`
		LeagueName []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"league_name"`
		ItemName []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"item_name"`
		ListPrice []struct {
			Currency      string  `json:"currency"`
			Value         float64 `json:"value"`
			MarketplaceID string  `json:"marketplace_id"`
		} `json:"list_price"`
		BatteriesRequired []struct {
			Value         bool   `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"batteries_required"`
		Material []struct {
			LanguageTag   string `json:"language_tag"`
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"material"`
	} `json:"attributes"`
	ProductTypes []struct {
		MarketplaceID string `json:"marketplaceId"`
		ProductType   string `json:"productType"`
	} `json:"productTypes"`
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
	MarketplaceID   string        `json:"marketplaceId"`
	Asin            string        `json:"asin"`
	ProductType     string        `json:"productType"`
	ConditionType   string        `json:"conditionType"`
	Status          []interface{} `json:"status"`
	ItemName        string        `json:"itemName"`
	CreatedDate     time.Time     `json:"createdDate"`
	LastUpdatedDate time.Time     `json:"lastUpdatedDate"`
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

//
//
// USED FOR CREATING LISTINGS
//
//

type ExternalProductIdentifier struct {
	Value         string `json:"value"`
	Type          string `json:"type"`
	MarketplaceID string `json:"marketplace_id"`
}
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
		Conditions        []ConditionType             `json:"condition_type"`
		ProductIdentifier []ExternalProductIdentifier `json:"externally_assigned_product_identifier"`
		ASIN              []MerchantSuggestedASIN     `json:"merchant_suggested_asin"`
		Offer             []Price                     `json:"purchasable_offer"`
	} `json:"attributes"`
}
type ProductTypeData struct {
	MetaSchema struct {
		Link struct {
			Resource string `json:"resource"`
			Verb     string `json:"verb"`
		} `json:"link"`
		Checksum string `json:"checksum"`
	} `json:"metaSchema"`
	Schema struct {
		Link struct {
			Resource string `json:"resource"`
			Verb     string `json:"verb"`
		} `json:"link"`
		Checksum string `json:"checksum"`
	} `json:"schema"`
	Requirements         string `json:"requirements"`
	RequirementsEnforced string `json:"requirementsEnforced"`
	PropertyGroups       struct {
		Offer struct {
			Title         string   `json:"title"`
			Description   string   `json:"description"`
			PropertyNames []string `json:"propertyNames"`
		} `json:"offer"`
		ProductIdentity struct {
			Title         string   `json:"title"`
			Description   string   `json:"description"`
			PropertyNames []string `json:"propertyNames"`
		} `json:"product_identity"`
	} `json:"propertyGroups"`
	Locale             string   `json:"locale"`
	MarketplaceIds     []string `json:"marketplaceIds"`
	ProductType        string   `json:"productType"`
	ProductTypeVersion struct {
		Version          string `json:"version"`
		Latest           bool   `json:"latest"`
		ReleaseCandidate bool   `json:"releaseCandidate"`
	} `json:"productTypeVersion"`
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
		signer.SignHTTP(context.Background(), aws.Credentials{
			AccessKeyID:     AppCred.IAMClientID,
			SecretAccessKey: AppCred.IAMSecretAccess,
		}, authS2ReqURLREQ, fmt.Sprintf("%x", bodySHA256), "execute-api",
			"us-east-1", time.Now().UTC())
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
			fmt.Println(string(respData))
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
	parameters := fmt.Sprintf("marketplaceIds=%s&includedData=%s", "ATVPDKIKX0DER", "issues,summaries")
	respData, success := APIURIConstruct("getListing", "https://sellingpartnerapi-na.amazon.com", getListingsAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCreds, tempCred, "")
	var tempListingData ListingsDetailPayload
	err := json.Unmarshal(respData, &tempListingData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respData))
	return tempListingData, success
}
func main() {
	var ClientCreds = ClientCredentialGenerator()
	// create a timer that will fire when the credentials expire
	ExpirationTimer := time.NewTimer(time.Duration(ClientCreds.ExpiresIn) * time.Second)
	var logFile, err = os.Create("run.log")
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetOutput(logFile)
	// OperationTimer := time.NewTimer(time.Duration(1) * time.Minute)
	// If product exists in Amazon:
	// Call the getListingsRestrictions operation with the ASIN identifier to retrieve any eligibility requirements that must be met before listing an item in the applicable condition.
	select {
	case <-ExpirationTimer.C:
		ClientCreds = ClientCredentialGenerator()
		ExpirationTimer.Reset(time.Duration(ClientCreds.ExpiresIn) * time.Second)
	default:
		//
		// https://developer-docs.amazon.com/sp-api/docs/building-listings-management-workflows-guide#list-an-offer-for-an-item-that-already-exists-in-the-amazon-catalog
		//
		// LISTING CREATION PROCESS
		//
		//
		// Call the searchCatalogItems operation to search for existing items in the Amazon catalog by product identifiers (UPC, EAN, etc.) or keywords.
		var tempASIN = "767345396723"
		var productName = "SOUTH DAKOTA STATE UNIVERSITY ROYAL BOTTLE OPENER KEYCHAIN"
		var searchCatalogAbsolutePath = fmt.Sprintf("/catalog/2022-04-01/items")
		identifier := "UPC"
		var parameters = fmt.Sprintf("identifiers=%s&identifiersType=%s&includedData=%s", tempASIN, identifier, "productTypes,attributes")
		productData, exists := APIURIConstruct("lookup", "https://sellingpartnerapi-na.amazon.com", searchCatalogAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCreds, GetCredentials(), "")
		var productLookupData LookupPayload
		err := json.Unmarshal(productData, &productLookupData)
		if err != nil {
			panic(err)
		}
		fmt.Println(productLookupData.Items[0].Attributes.TeamName[0].Value)
		if productLookupData.NumberOfResults > 1 {
			for i := range productLookupData.Items {
				if len(productLookupData.Items[i].Attributes.TeamName) > 0 {
					if strings.Contains(productLookupData.Items[i].Attributes.TeamName[0].Value, productName) {
						fmt.Println("valla var")
						break
					}
				}
			}
		} else {
			if exists {
				var getListingsRestrictionsAbsolutePath = fmt.Sprintf("/listings/2021-08-01/restrictions")
				tempCred := GetCredentials()
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
					fmt.Println("Printing `productLookupData.Items[0].ProductTypes[0].ProductType`")
					fmt.Println(productLookupData.Items[0].ProductTypes[0].ProductType)
					tempProductPUT.ProductType = productLookupData.Items[0].ProductTypes[0].ProductType
					// this is the only thing tempProductType is required for
					tempProductPUT.Requirements = "LISTING_OFFER_ONLY"
					tempProductPUT.Attributes.ASIN = append(tempProductPUT.Attributes.ASIN, MerchantSuggestedASIN{
						Value:         "          ",
						MarketplaceID: "ATVPDKIKX0DER",
					})
					tempProductPUT.Attributes.ProductIdentifier = append(tempProductPUT.Attributes.ProductIdentifier, ExternalProductIdentifier{
						Value:         tempASIN,
						Type:          identifier,
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
					var createListingAbsolutePath = fmt.Sprintf("/listings/2021-08-01/items/%s/%s", tempCred.SellerID, "asd-676")
					parameters := fmt.Sprintf("&marketplaceIds=%s", "ATVPDKIKX0DER")
					_, success := APIURIConstruct("createListing", "https://sellingpartnerapi-na.amazon.com", createListingAbsolutePath, parameters, "ATVPDKIKX0DER", "PUT", ClientCreds, tempCred, string(tempProductJSON))
					//
					//
					// LISTING CREATION PROCESS IS DONE, SLEEP A BIT TO WAIT FOR THE LISTING TO BE CREATED
					//
					//
					time.Sleep(5 * time.Second)
					if success {
						ListingDetails, success := ReturnListingDetails(tempCred, ClientCreds, "asd-676")
						if success {
							isSummariesEmpty := len(ListingDetails.Summaries) > 0
							isIssuesEmpty := len(ListingDetails.Issues) > 0
							// Summaries may be empty but issues may have data, case 1.
							if isSummariesEmpty && isIssuesEmpty == false {
								logger.Errorf("Listing %s is suppressed and not discorable nor buyable due to: ", tempASIN)
								for _, issue := range ListingDetails.Issues {
									logger.Errorln("Code: ", issue.Code, "Message: ", issue.Message)
									logger.Errorln("Severity: ", issue.Severity)
								}
								logger.Errorln("If message is blank, then only vendor can fix the problem.")
								break
								// Summaries and issues may not have data, case 2. Listing is created but not updated, yet. Sleep a bit and try again.
							} else if isSummariesEmpty == true && isIssuesEmpty == true {
								time.Sleep(5 * time.Second)
								ListingDetails, _ := ReturnListingDetails(tempCred, ClientCreds, "asd-676")
								if reflect.DeepEqual(ListingDetails.Summaries, []ListingsDetailPayloadSummaries{}) == true &&
									reflect.DeepEqual(ListingDetails.Issues, []ListingsDetailPayloadIssuses{}) == true {
									logger.Errorln("There was an error creating the listing for the ASIN: ", tempASIN)
									break
								}
								// Both summaries and issues have data, case 3. Listing is created and updated.
							} else if isSummariesEmpty == false && isIssuesEmpty == true { // check if summaries and issues are empty
								logger.Errorf("Listing %s is suppressed and not discorable nor buyable due to: ", tempASIN)
								for _, issue := range ListingDetails.Issues {
									logger.Errorln("Code: ", issue.Code, "Message: ", issue.Message)
									logger.Errorln("Severity: ", issue.Severity)
								}
								logger.Errorln("Summary status: ", ListingDetails.Summaries)
								logger.Errorln("If message is blank, then only vendor can fix the problem.")
								break
							}

						}
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
	}

}
