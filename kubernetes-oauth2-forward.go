package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var clientId string
var clientSecret string

var grantType = "authorization_code"
var responseType = "code"
var sessions = make(map[int]string)
var scope = "default"

var redirectUri string
var authTokenUri = "https://auth.viarezo.fr/oauth/token"
var authAPIUri = "https://auth.viarezo.fr/api/user/show/me"
var authAuthorizeUri = "https://auth.viarezo.fr/oauth/authorize"

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type User struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

func main() {
	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectUri = os.Getenv("REDIRECT_URI")
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateState() string {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	hash := sha256.Sum256(data)
	base64String := base64.StdEncoding.EncodeToString(hash[:])
	return base64String[:len(base64String)-2]
}

func getId() int {
	for i := 0; i < 10000; i++ {
		_, ok := sessions[i]
		if !ok {
			return i
		}
	}
	return -1
}
func redirectAuth(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	id := getId()
	sessions[id] = state

	cookie := &http.Cookie{
		Name:  "id",
		Value: fmt.Sprint(id),
	}
	http.SetCookie(w, cookie)

	parameters := url.Values{}
	parameters.Add("redirect_uri", redirectUri)
	parameters.Add("client_id", clientId)
	parameters.Add("response_type", responseType)
	parameters.Add("state", state)
	parameters.Add("scope", scope)

	http.Redirect(w, r, authAuthorizeUri+"?"+parameters.Encode(), 302)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	codeArray, hasCode := query["code"]
	stateArray, hasState := query["state"]

	if len(query) == 0 {
		redirectAuth(w, r)
		return
	}
	if hasCode && hasState {
		id, cookieErr := r.Cookie("id")
		if cookieErr != nil {
			log.Fatal(cookieErr)
		}
		name, ok := checkIdentity(codeArray[0], stateArray[0], id.Value)
		if ok {
			fmt.Fprintf(w, "Hello %s!\n", name)
			return
		}
		http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
		return
	}

	http.Error(w, "400 bad request.", http.StatusBadRequest)
}

func checkIdentity(code, state, id string) (string, bool) {
	intId, _ := strconv.Atoi(id)
	if state != sessions[intId] {
		return "", false
	}
	token, ok := getTokens(code)
	if !ok {
		return "", false
	}
	delete(sessions, intId)
	name := getUsername(token.AccessToken)
	return name, true
}

func getTokens(code string) (*Token, bool) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	form := url.Values{}
	form.Set("grant_type", grantType)
	form.Set("code", code)
	form.Set("redirect_uri", redirectUri)
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)
	encodedForm := form.Encode()
	req, err := http.NewRequest("POST", authTokenUri, strings.NewReader(encodedForm))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		log.Fatal(reqErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return &Token{}, false
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	result := &Token{}
	decodeErr := decoder.Decode(result)
	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	return result, true
}

func getUsername(accessToken string) string {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", authAPIUri, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, respErr := client.Do(req)
	if respErr != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	result := &User{}
	decodeErr := decoder.Decode(result)
	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	return result.FirstName + " " + result.LastName
}
