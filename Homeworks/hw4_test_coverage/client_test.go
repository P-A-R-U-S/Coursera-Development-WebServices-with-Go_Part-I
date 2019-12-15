package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"encoding/xml"
	"os"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"strings"
	"time"
)

type Root struct {
	Rows []Row `xml:"row"`
}

/*

 <row>
    <id>0</id>
    <guid>1a6fa827-62f1-45f6-b579-aaead2b47169</guid>
    <isActive>false</isActive>
    <balance>$2,144.93</balance>
    <picture>http://placehold.it/32x32</picture>
    <age>22</age>
    <eyeColor>green</eyeColor>
    <first_name>Boyd</first_name>
    <last_name>Wolf</last_name>
    <gender>male</gender>
    <company>HOPELI</company>
    <email>boydwolf@hopeli.com</email>
    <phone>+1 (956) 593-2402</phone>
    <address>586 Winthrop Street, Edneyville, Mississippi, 9555</address>
    <about>Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.
</about>
    <registered>2017-02-05T06:23:27 -03:00</registered>
    <favoriteFruit>apple</favoriteFruit>
  </row>

*/
type Row struct {
	About         string `xml:"about"`
	Address       string `xml:"address"`
	Age           int `xml:"age"`
	Balance       string `xml:"balance"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	EyeColor      string `xml:"eyeColor"`
	FavoriteFruit string `xml:"favoriteFruit"`
	FirstName     string `xml:"first_name"`
	Gender        string `xml:"gender"`
	GUID          string `xml:"guid"`
	ID            int `xml:"id"`
	IsActive      bool `xml:"isActive"`
	LastName      string `xml:"last_name"`
	Phone         string `xml:"phone"`
	Picture       string `xml:"picture"`
	Registered    string `xml:"registered"`
}

func SearchServer(w http.ResponseWriter, r *http.Request)  {

	//
	var accessToken = r.Header.Get("AccessToken")
	if len(accessToken) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized"))
		return
	}

	/*Params*/
	var limit = 0
	var lParam = r.URL.Query().Get("limit")
	if l, ok := strconv.Atoi(lParam); ok == nil {
		limit = l
	}

	//
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()



	v, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	q := new(Root)
	err = xml.Unmarshal(v, &q)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	var users []User
	for _, r := range q.Rows {


		u := User{r.ID,
		r.FirstName + " " + r.LastName,
		r.Age,
		r.About,
		r.Gender}

		users = append(users, u)
	}

	//Filter By limit
	if limit > 0 {
		users = users[:limit]
	}
	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshal json", err)
		return
	}

	w.Write(data)
}

func Test_FindUsers(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	var client = SearchClient{URL:ts.URL}
	client.AccessToken = "ACCESS1234567890"
	var req = SearchRequest{}



	_, err := client.FindUsers(req)

	if err != nil {
		t.Fail()
	}
}

func Test_FindUsers_Limit__non_negative_test(t *testing.T) {
	var client = SearchClient{}
	var req = SearchRequest{}
	req.Limit = -1

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}
	if err.Error() != "limit must be > 0" {
		t.Fail()
	}

}

func Test_FindUsers_Offset__non_negative_test(t *testing.T) {
	var client = SearchClient{}
	var req = SearchRequest{}
	req.Offset = -1

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}
	if err.Error() != "offset must be > 0" {
		t.Fail()
	}

}

func Test_FindUsers_Limit_more_than_25(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	var client = SearchClient{URL:ts.URL}
	client.AccessToken = "ACCESS1234567890"
	var req = SearchRequest{}
	req.Limit = 26

	res, err := client.FindUsers(req)
	if err != nil {
		t.Fail()
	}

	if len(res.Users) != 25 {
		t.Fail()
	}

}

func Test_FindUsers_Return_less_than_provided_in_limit(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal([]User{{Id: 0, Name: "FakeUser", Age: 34, About: "About", Gender: "Male"}})
		if err != nil {
			fmt.Println("Error marshal json", err)
			return
		}

		w.Write(data)
	}))

	var client = SearchClient{URL:ts.URL}
	client.AccessToken = "ACCESS1234567890"
	var req = SearchRequest{}
	req.Limit = 26

	res, err := client.FindUsers(req)
	if err != nil {
		t.Fail()
	}

	if len(res.Users) != 1 {
		t.Fail()
	}

	if res.NextPage {
		t.Fail()
	}
}

func Test_FindUsers_Bad_AccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if err.Error() != "Bad AccessToken" {
		t.Fail()
	}
}

func Test_FindUsers_Broken_JSON(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	}))

	var client= SearchClient{URL: ts.URL}
	client.AccessToken = "ACCESS1234567890"
	var req= SearchRequest{}
	req.Limit = 26
	req.Query = ""

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(),"cant unpack result json") {
		t.Fail()
	}
}

func Test_FindUsers_SearchServer_Fatal_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
	}))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if err.Error() != "SearchServer fatal error" {
		t.Fail()
	}
}

func Test_FindUser_Timeout(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	}))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(), "timeout for") {
		t.Fail()
	}
}

func Test_FindUser_StatusBadRequest(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(),"cant unpack error json: unexpected end of JSON input") {
		t.Fail()
	}
}

func Test_FindUser_StatusBadRequest_ErrorBadOrderField(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		res := SearchErrorResponse{Error:"ErrorBadOrderField"}

		data, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Error marshal json", err)
			return
		}

		w.Write(data)
	}))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}
	req.OrderField = "Name"

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("OrderFeld %s invalid", req.OrderField)) {
		t.Fail()
	}
}

func Test_FindUser_StatusBadRequest_Error(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		res := SearchErrorResponse{Error:"Field are not exist."}

		data, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Error marshal json", err)
			return
		}

		w.Write(data)
	}))

	var client = SearchClient{URL:ts.URL}
	var req = SearchRequest{}
	req.OrderField = "Name"

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("unknown bad request error: %s", "Field are not exist.")) {
		t.Fail()
	}
}

func Test_FindUser_Other(t *testing.T) {

	var client = SearchClient{URL:""}
	var req = SearchRequest{}

	_, err := client.FindUsers(req)
	if err == nil {
		t.Fail()
	}

	if !strings.HasPrefix(err.Error(), "unknown error ") {
		t.Fail()
	}
}