package main

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"encoding/json"
	"github.com/mailru/easyjson"
	"io"
	"os"
	"bufio"
	"strings"
	"fmt"
)

type User struct {
	Browsers []string `json:"browsers"`
	//Company  string   `json:"_"`
	//Country  string   `json:"_"`
	Email    string   `json:"email"`
	//Job      string   `json:"_"`
	Name     string   `json:"name"`
	//Phone    string   `json:"_"`
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
	scanner *bufio.Scanner
)

/*
{
   "browsers":[
      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
      "LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1",
      "Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1",
      "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"
   ],
   "company":"Flashpoint",
   "country":"Dominican Republic",
   "email":"JonathanMorris@Muxo.edu",
   "job":"Programmer Analyst #{N}",
   "name":"Sharon Crawford",
   "phone":"176-88-49"
}{
   "browsers":[
      "Mozilla/5.0 (X11; FreeBSD amd64) AppleWebKit/537.4 (KHTML like Gecko) Chrome/22.0.1229.79 Safari/537.4",
      "Mozilla/5.0 (Linux; U; Android 1.5; en-gb; T-Mobile_G2_Touch Build/CUPCAKE) AppleWebKit/528.5  (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1",
      "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Trident/5.0)",
      "Mozilla/5.0 (iPad; U; CPU OS 4_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8F190 Safari/6533.18.5"
   ],
   "company":"Jatri",
   "country":"Kenya",
   "email":"eum_rerum_explicabo@Topiczoom.info",
   "job":"Web Developer #{N}",
   "name":"Susan Ellis",
   "phone":"187-70-57"
}{
   "browsers":[
      "Mozilla/5.0 (X11; U; Linux x86_64; en-US) AppleWebKit/534.15 (KHTML, like Gecko) Chrome/10.0.613.0 Safari/534.15",
      "Mozilla/5.0 (X11; Linux i686; rv:49.0) Gecko/20100101 Firefox/49.0",
      "Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/601.1 (KHTML, like Gecko) CriOS/49.0.2623.109 Mobile/14A5335b Safari/601.1.46",
      "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.71 (KHTML like Gecko) WebVideo/1.0.1.10 Version/7.0 Safari/537.71"
   ],
   "company":"Dabtype",
   "country":"Ecuador",
   "email":"accusamus_et_magnam@Voonix.gov",
   "job":"Internal Auditor",
   "name":"Joshua Fisher",
   "phone":"655-76-81"
}{
   "browsers":[
      "Mozilla/5.0 (X11; Linux x86_64; en-US; rv:2.0b2pre) Gecko/20100712 Minefield/4.0b2pre",
      "Mozilla/5.0 (Symbian/3; Series60/5.2 NokiaE7-00/010.016; Profile/MIDP-2.1 Configuration/CLDC-1.1 ) AppleWebKit/525 (KHTML, like Gecko) Version/3.0 BrowserNG/7.2.7.3 3gpp-gba",
      "Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0)",
      "Mozilla/5.0 (iPhone; U; CPU iPhone OS 5_1_1 like Mac OS X; da-dk) AppleWebKit/534.46.0 (KHTML, like Gecko) CriOS/19.0.1084.60 Mobile/9B206 Safari/7534.48.3"
   ],
   "company":"Thoughtbeat",
   "country":"Thailand",
   "email":"sed_reiciendis_qui@Bluezoom.net",
   "job":"Office Assistant #{N}",
   "name":"Kathleen Williams",
   "phone":"8-912-514-24-86"
}{
   "browsers":[
      "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10240",
      "Mozilla/4.0 (compatible; GoogleToolbar 4.0.1019.5266-big; Windows XP 5.1; MSIE 6.0.2900.2180)",
      "SAMSUNG-S8000/S8000XXIF3 SHP/VPP/R5 Jasmine/1.0 Nextreaming SMM-MMS/1.2.0 profile/MIDP-2.1 configuration/CLDC-1.1 FirePHP/0.3",
      "Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.8.1) Gecko/20061024 Firefox/2.0 (Swiftfox)"
   ],
   "company":"Youfeed",
   "country":"Uruguay",
   "email":"gLong@Zoombox.org",
   "job":"Automation Specialist #{N}",
   "name":"Michael Davis",
   "phone":"946-49-01"
}{
   "browsers":[
      "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 7 Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.78 Safari/537.36 OPR/30.0.1856.93524",
      "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/6.0)",
      "Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US) AppleWebKit/532.9 (KHTML, like Gecko) Chrome/5.0.310.0 Safari/532.9",
      "Mozilla/5.0 (compatible; Yahoo! Slurp China; http://misc.yahoo.com.cn/help.html)"
   ],
   "company":"Feedbug",
   "country":"Germany",
   "email":"yKing@Leexo.net",
   "job":"Cost Accountant",
   "name":"Melissa Price",
   "phone":"3-981-961-68-80"
}{
   "browsers":[
      "SonyEricssonW580i/R6BC Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1",
      "Mozilla/5.0 (Windows NT 6.2; ARM; Trident/7.0; Touch; rv:11.0; WPDesktop; NOKIA; Lumia 920) like Geckoo",
      "Mozilla/5.0 (Linux; U; Android 2.0; en-us; Milestone Build/ SHOLS_U2_01.03.1) AppleWebKit/530.17 (KHTML, like Gecko) Version/4.0 Mobile Safari/530.17",
      "Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/51.0.2704.79 Chrome/51.0.2704.79 Safari/537.36"
   ],
   "company":"Wordware",
   "country":"Equatorial Guinea",
   "email":"est_ratione_maxime@Thoughtsphere.gov",
   "job":"Electrical Engineer",
   "name":"Maria Watkins",
   "phone":"5-525-472-31-88"
}{
   "browsers":[
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.78.1 (KHTML like Gecko) Version/7.0.6 Safari/537.78.1",
      "Mozilla/1.22 (compatible; MSIE 5.01; PalmOS 3.0) EudoraWeb 2.1",
      "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2049.0 Safari/537.36",
      "Mozilla/5.0 (SymbianOS/9.4; U; Series60/5.0 SonyEricssonP100/01; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) Version/3.0 Safari/525"
   ],
   "company":"Flashspan",
   "country":"Turkmenistan",
   "email":"est_adipisci@Brightbean.name",
   "job":"Research Assistant #{N}",
   "name":"Johnny Jones",
   "phone":"3-553-621-66-19"
}
*/
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	replacer := strings.NewReplacer("@" , " [at] ")
	var seenBrowsers []string
	foundUsers := ""

	scanner = bufio.NewScanner(file)
	i := 0
	var notSeenBefore bool
	var isAndroid, isMSIE, isAndroidContains, isMSIEContains bool
	var user User

	for scanner.Scan() {

		//user = userPool.Get().(User)
		user.UnmarshalJSON([]byte(scanner.Text()))

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {

			isAndroidContains = strings.Contains(browser, "Android")
			isMSIEContains = strings.Contains(browser, "MSIE")

			if isAndroidContains || isMSIEContains {

				isAndroid = isAndroidContains || isAndroid
				isMSIE = isMSIEContains || isMSIE

				notSeenBefore = true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
						continue
					}
				}

				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
				}
			}
		}

		if isAndroid && isMSIE {
			// log.Println("Android and MSIE user:", user["name"], user["email"])

			email := replacer.Replace(user.Email) //,"@" , " [at] ", 1)
			foundUsers += fmt.Sprintf("[%d] %v <%v>\n", i, user.Name, email)
		}

		i++
	}


	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

func easyjson3486653aDecodeCourseraHomeworksHw3Bench(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		//case "company":
		//	out.Company = string(in.String())
		//case "country":
		//	out.Country = string(in.String())
		case "email":
			out.Email = string(in.String())
		//case "job":
		//	out.Job = string(in.String())
		case "name":
			out.Name = string(in.String())
		//case "phone":
		//	out.Phone = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3486653aDecodeCourseraHomeworksHw3Bench(&r, v)
	return r.Error()
}
