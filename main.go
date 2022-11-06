package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProductDetails struct {
	Sku       string `json:"sku"`
	Summaries []struct {
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
		} `json:"mainImage"`
	} `json:"summaries"`
	Attributes struct {
		ConditionType []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"condition_type"`
		MerchantShippingGroup []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"merchant_shipping_group"`
		MerchantSuggestedAsin []struct {
			Value         string `json:"value"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"merchant_suggested_asin"`
		PurchasableOffer []struct {
			Currency string `json:"currency"`
			StartAt  struct {
				Value time.Time `json:"value"`
			} `json:"start_at"`
			OurPrice []struct {
				Schedule []struct {
					ValueWithTax float64 `json:"value_with_tax"`
				} `json:"schedule"`
			} `json:"our_price"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"purchasable_offer"`
		FulfillmentAvailability []struct {
			FulfillmentChannelCode string `json:"fulfillment_channel_code"`
			Quantity               int    `json:"quantity"`
			MarketplaceID          string `json:"marketplace_id"`
		} `json:"fulfillment_availability"`
		MainProductImageLocator []struct {
			MediaLocation string `json:"media_location"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"main_product_image_locator"`
		OtherProductImageLocator1 []struct {
			MediaLocation string `json:"media_location"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"other_product_image_locator_1"`
		OtherProductImageLocator2 []struct {
			MediaLocation string `json:"media_location"`
			MarketplaceID string `json:"marketplace_id"`
		} `json:"other_product_image_locator_2"`
	} `json:"attributes"`
	Issues []struct {
		Message        string   `json:"message"`
		Severity       string   `json:"severity"`
		AttributeName  string   `json:"attributeName"`
		AttributeNames []string `json:"attributeNames"`
	} `json:"issues"`
	Offers []struct {
		MarketplaceID string `json:"marketplaceId"`
		OfferType     string `json:"offerType"`
		Price         struct {
			Currency string `json:"currency"`
			Amount   string `json:"amount"`
		} `json:"price"`
	} `json:"offers"`
	FulfillmentAvailability []struct {
		FulfillmentChannelCode string `json:"fulfillmentChannelCode"`
		Quantity               int    `json:"quantity"`
	} `json:"fulfillmentAvailability"`
}
type AwsRequestSigner struct {
	RegionName             string
	AwsSecurityCredentials map[string]string
}

const (
	AwsAlgorithm         = "AWS4-HMAC-SHA256"
	AwsRequestType       = "aws4_request"
	AwsAccessTokenHeader = "x-amz-access-token"
	AwsDateHeader        = "x-amz-date"
)

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

// ClientCredentialsGenerator returns Access Token and other shit.
func ClientCredentialsGenerator() ClientCredentials {
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
func Sha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
func HMACSha256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
func APIURIConstruct(endpoint string, requestPath string, parameters string, marketplace string, httpMethod string, credentials ClientCredentials, AppCred AppCredentials) {
	authS2ReqURL, err := url.Parse(endpoint + requestPath + "&" + "marketplaceIds=" + marketplace)
	if err != nil {
		panic(err)
	}
	fmt.Println(authS2ReqURL.String())
	authS2Req, err := http.NewRequest(httpMethod, authS2ReqURL.String(), nil)
	if err != nil {
		panic(err)
	}
	// CREATE CANONICAL URL
	var canonicalURL string

	// calculate the signature
	// WHAT THE FUCK IS THIS? A FUCKING RABBITHOLE? WERE YOU BORED WHILE THINKING THIS? SERIOUSLY?
	// derive a signing key first
	kSecret := AppCred.IAMSecretAccess
	kDate := HMACSha256([]byte("AWS4"+kSecret), []byte(time.Now().UTC().Format("20060102")))
	kRegion := HMACSha256(kDate, []byte("us-east-1"))
	kService := HMACSha256(kRegion, []byte("execute-api"))
	kSigning := HMACSha256([]byte("aws4_request"), kService)
	// create credential scope
	var awsRegion = "us-east-1"
	var exec = "execute-api"
	var terminationString = "aws4_request"
	var credentialScope = time.Now().UTC().Format("20060102") + "/" + awsRegion + "/" + exec + "/" + terminationString
	// create a string to sign and declare signedHeaders and canonicalheaders.
	var signedHeaders string
	signedHeaders += "host;"
	signedHeaders += "user-agent;"
	signedHeaders += "x-amz-access-token;"
	signedHeaders += "x-amz-date"
	var canonicalHeaders string
	canonicalHeaders += strings.ToLower("host") + ":" + strings.TrimSpace("sellingpartnerapi-na.amazon.com") + "\n"
	canonicalHeaders += strings.ToLower("User-Agent") + ":" + strings.TrimSpace("serpent-of-time/0.31 (Language=Go; Windows)") + "\n"
	canonicalHeaders += strings.ToLower("x-amz-access-token") + ":" + strings.TrimSpace(credentials.AccessToken) + "\n"
	canonicalHeaders += strings.ToLower("x-amz-date") + ":" + strings.TrimSpace(GetTime()) + "\n"
	// hash an empty string with SHA256 algorithm
	emptyHash := Sha256([]byte(""))
	// crete url encoded query string
	tempURL, _ := url.Parse(requestPath)
	// have the marketplaces in the query string, but escape the string.
	var query string
	query += url.QueryEscape("&marketplaceIds=" + marketplace)
	// create escaped string
	escapedURL := tempURL.String() + query
	canonicalURL = fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%x", httpMethod, escapedURL,
		strings.TrimSpace(canonicalHeaders), strings.TrimSpace(signedHeaders), string(emptyHash))
	canonicalURL256 := Sha256([]byte(canonicalURL))
	canonicalURLHexed := hex.EncodeToString(canonicalURL256)
	canonicalURLHexed = strings.ToLower(canonicalURLHexed)
	var stringToSign = "AWS4-HMAC-SHA256" + "\n" + GetTime() + "\n" + credentialScope + "\n" + canonicalURLHexed
	// signature = HexEncode(HMAC(derived signing key, string to sign))
	signature := hmac.New(sha256.New, kSigning)
	signature.Write([]byte(stringToSign))
	signatureSum := signature.Sum(nil)
	var signatureSumHexed = make([]byte, hex.EncodedLen(len(signatureSum)))
	hex.Encode(signatureSumHexed, signatureSum)
	fmt.Println("--------------------")
	fmt.Println(string(signatureSumHexed))
	fmt.Println("--------------------")
	// create canonicalheaders
	var authHeader = fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s, SignedHeaders=%s, Signature=%x, X-Amz-Date=%s", AppCred.IAMClientID+"/"+credentialScope, signedHeaders, signatureSumHexed, GetTime())
	// start writing to canonicalHeaders one by one
	// create authheader, god fuck me this is never-ending rabbit hole of... words arent enough to describe this mess.
	authS2Req.Header.Set("Authorization", authHeader)
	authS2Req.Header.Set("SignedHeaders", "host;user-agent;x-amz-access-token")
	authS2Req.Header.Set("Content-Type", "application/json")
	authS2Req.Header.Set("User-Agent", "serpent-of-time/0.31 (Language=Go; Windows)")
	authS2Req.Header.Set("host", "sellingpartnerapi-na.amazon.com")
	authS2Req.Header.Set("x-amz-date", GetTime())
	authS2Req.Header.Set("x-amz-access-token", credentials.AccessToken)
	resp, err := http.DefaultClient.Do(authS2Req)
	if err != nil {
		panic(err)
	}
	respData, err := io.ReadAll(resp.Body)
	fmt.Println(string(respData))
	var tempItemData ProductDetails
	err = json.Unmarshal(respData, &tempItemData)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respData))
}
func GetTime() string {
	return time.Now().UTC().Format("20060102T150405Z")
}
func main() {
	//
	// https://developer-docs.amazon.com/sp-api/docs/building-listings-management-workflows-guide#list-an-offer-for-an-item-that-already-exists-in-the-amazon-catalog
	//
	// WORKFLOW:
	// Call the searchCatalogItems operation to search for existing items in the Amazon catalog by product identifiers (UPC, EAN, etc.) or keywords.
	var tempASIN = "B00ABALPNA"
	var searchCatalogAbsolutePath = fmt.Sprintf("/catalog/2022-04-01/items")
	var parameters = fmt.Sprintf("?identifiers=%s?identifiersType=%s", tempASIN, "ASIN")
	APIURIConstruct("https://sellingpartnerapi-na.amazon.com", searchCatalogAbsolutePath, parameters, "ATVPDKIKX0DER", "GET", ClientCredentialsGenerator(), GetCredentials())
}
