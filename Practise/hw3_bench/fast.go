package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type User struct {
	Browsers []string `json:"Browsers"`
	Company  string   `json:"Company"`
	Email    string   `json:"Email"`
	Job      string   `json:"Job"`
	Name     string   `json:"Name"`
	Phone    string   `json:"Phone"`
}

func (u *User) Parse (line []byte) {
	var i, si int
	var s string
	i = 13//len(`{"browsers":["`)
	//i--
	for {
		i++
		si = i
		for line[i] != byte('"') {
			i++
		}

		s = string(line[si:i])
		u.Browsers = append(u.Browsers, s)

		i++

		if line[i] == byte(']') {
			break
		}

		i += 1
	}

	// Company
	i += 13 //len(`,"company":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	u.Company = string(line[si:i])
	//i = si
	//fmt.Printf("Company:%s\n", u.Company)

	// County
	i += 13 //len(`,"country":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	//u.Country = string(line[i:si])
	//fmt.Printf("Country:%s\n", string(line[i:si]))
	//i = si

	// Email
	i += 11 //len(`,"email":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	u.Email = string(line[si:i])
	//i = si
	//fmt.Printf("Email:%s\n", u.Email)

	// Job
	i += 9 //len(`,"job":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	u.Job = string(line[si:i])
	//i = si
	//fmt.Printf("Job:%s\n", u.Job)

	// Name
	i += 10 //len(`,"name":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	u.Name = string(line[si:i])
	//i = si
	//fmt.Printf("Name:%s\n", u.Name)

	// Phone
	i += 11 //len(`,"phone":"`)
	//i++
	si = i
	for line[i] != byte('"') {
		i++
	}
	u.Phone = string(line[si:i])
	//i = si
	//fmt.Printf("Phone:%s\n", u.Phone)
}

/*
BenchmarkSlow-8               20          92060028 ns/op        319476030 B/op    276175 allocs/op
BenchmarkFast-8              500           2654428 ns/op          1880622 B/op      8605 allocs/op
==================================================================================================

BenchmarkFast-8              500           2654428 ns/op          1880622 B/op      8605 allocs/op

BenchmarkFast-8              200           9949424 ns/op         2579351 B/op      17612 allocs/op
BenchmarkFast-8              200           8989090 ns/op         2316911 B/op      17600 allocs/op
BenchmarkFast-8              200           8744798 ns/op         2118523 B/op      15607 allocs/op
BenchmarkFast-8              200           8827188 ns/op         2118480 B/op      15607 allocs/op
BenchmarkFast-8              200           8788168 ns/op         2118526 B/op      15607 allocs/op
BenchmarkFast-8              200           8223974 ns/op         2118490 B/op      15607 allocs/op
BenchmarkFast-8              200           7580545 ns/op          928610 B/op      13603 allocs/op
BenchmarkFast-8              200           7766176 ns/op          928601 B/op      13603 allocs/op
BenchmarkFast-8              200           7567046 ns/op          928603 B/op      13603 allocs/op
BenchmarkFast-8             1000           2147625 ns/op          974253 B/op      13601 allocs/op
BenchmarkFast-8             1000           2000943 ns/op          974234 B/op      13601 allocs/op
BenchmarkFast-8             1000           2010017 ns/op          974250 B/op      13601 allocs/op
BenchmarkFast-8             1000           2012561 ns/op          974245 B/op      13601 allocs/op
BenchmarkFast-8             1000           2009935 ns/op          974258 B/op      13601 allocs/op

*/
// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	var user = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(User)
		},
	}

	//var line = sync.Pool{
	//	New: func() interface{} {
	//		// The Pool's New function should generally only return pointer
	//		// types, since a pointer can be put into the return interface
	//		// value without an allocation:
	//		return new([]byte)
	//	},
	//}

	scanner := bufio.NewScanner(file)
	i := -1
	var notSeenBefore, isAndroid, isMSIE, isAndroidContains, isMSIEContains bool

	for scanner.Scan() {
		line :=  scanner.Bytes()

		var u = user.Get().(*User)
		u.Browsers = make([]string, 1)
		// fmt.Printf("%v %v\n", err, line)
		//err := json.Unmarshal(line, &u)
		u.Parse(line)
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false

		//Browsers, ok := user["Browsers"].([]interface{})
		//if !ok {
		//	// log.Println("cant cast Browsers")
		//	continue
		//}
		l := len(u.Browsers)
		for j:=0; j<l; j++ {

			isAndroidContains = strings.Contains(u.Browsers[j], "Android")
			isMSIEContains =  strings.Contains(u.Browsers[j], "MSIE")

			if isAndroidContains || isMSIEContains {

				isAndroid = isAndroid || isAndroidContains
				isMSIE = isMSIE || isMSIEContains

				notSeenBefore = true
				for _, item := range seenBrowsers {
					if item == u.Browsers[j] {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["Name"])
					seenBrowsers = append(seenBrowsers, u.Browsers[j])
					uniqueBrowsers++
				}
			}

		}

		user.Put(u)
		i++
		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["Name"], user["Email"])
		email := strings.Replace( u.Email,"@", " [at] ", 1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, u.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

// {"browsers":["Mozilla/4.0 (compatible; MSIE 5.5; Windows NT 5.0 )","Opera/9.80 (S60; SymbOS; Opera Mobi/499; U; ru) Presto/2.4.18 Version/10.00","Mozilla/5.0 (X11; U; SunOS i86pc; en-US; rv:1.9.1b3) Gecko/20090429 Firefox/3.1b3","Mozilla/5.0 (X11; U; Linux x86_64; en-us) AppleWebKit/537.36 (KHTML, like Gecko)  Chrome/30.0.1599.114 Safari/537.36 Puffin/4.8.0.2965AT"],"company":"Brainverse","country":"Iraq","email":"hic_iusto@Yodoo.gov","job":"Office Assistant #{N}","name":"Roger Martinez","phone":"011-86-99"}
