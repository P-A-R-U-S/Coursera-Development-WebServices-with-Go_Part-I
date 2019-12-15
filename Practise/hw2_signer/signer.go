package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job ) {

	wg := &sync.WaitGroup{}

	var in, out chan interface{}

	for _,v := range jobs {
		in = out
		out = make(chan interface{}, 1)
		var j = v

		wg.Add(1)
		go func(wg *sync.WaitGroup, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)

			j(in, out)
		}(wg, in, out)
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {

	//0 SingleHash data 0
	//0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
	//0 SingleHash crc32(md5(data)) 502633748
	//0 SingleHash crc32(data) 4108050209
	//0 SingleHash result 4108050209~502633748

	wg := &sync.WaitGroup{}

	//var index = 0

	var mu = &sync.Mutex{}

	for data := range in {

		var out_crc32_data 		= make(chan string, 1)
		var out_md5_data 		= make(chan string, 1)
		var out_crc32_md5_data 	= make(chan string, 1)

		var s = strconv.Itoa(data.(int))
		//fmt.Printf("%d SingleHash data %s\n", index, s)

		wg.Add(1)
		go func(wg_singleHash *sync.WaitGroup, out_channel chan string) {

			defer wg.Done()

			var crc32_data = DataSignerCrc32(s)

			//fmt.Printf("%d SingleHash crc32(data) %s\n", index, crc32_data)

			out_channel <- crc32_data
		}(wg,out_crc32_data)

		wg.Add(1)
		go func(wg_singleHash *sync.WaitGroup, mu_singleHash *sync.Mutex, out_channel chan string) {
			defer wg.Done()
			mu.Lock()
			var md5_data = DataSignerMd5(s)
			mu.Unlock()
			//fmt.Printf("%d SingleHash md5(data) %s\n", index, md5_data)

			out_channel <- md5_data
		}(wg, mu, out_md5_data)

		wg.Add(1)
		go func(wg_singleHash *sync.WaitGroup, in_channel chan string, out_channel chan string) {
			defer wg.Done()

			var md5_data = <-in_channel
			var crc32_md5_data = DataSignerCrc32(md5_data)

			//fmt.Printf("%d SingleHash crc32(md5(data)) %s\n", index, crc32_md5_data)

			out_channel <- crc32_md5_data
		}(wg,out_md5_data, out_crc32_md5_data)

		wg.Add(1)
		go func(wg_singleHash *sync.WaitGroup, in_channel_crc32_data chan string, in_channel_crc32_md5_data chan string) {
			defer wg.Done()
			defer close(in_channel_crc32_data)
			defer close(in_channel_crc32_md5_data)


			var crc32_md5_data =<- in_channel_crc32_md5_data
			var crc32_data =<- in_channel_crc32_data


			var result = crc32_data + "~" + crc32_md5_data

			//fmt.Printf("%d SingleHash result %s\n", index, result)
			//index++
			out <- result

		}(wg, out_crc32_data,out_crc32_md5_data)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {

	type hash struct {
		i int
		value string
	}
	/*
	4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
	4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
	4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
	4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
	4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
	4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
	4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555
	 */

	wg := &sync.WaitGroup{}

	for data := range in {

		var s = data.(string)

		var multiHash_channel = make(chan hash)

		wg.Add(1)
		go func(wg_multiHash *sync.WaitGroup, in_multiHash chan hash) {

			defer wg_multiHash.Done()

			defer close(in_multiHash)

			var a [6]string
			for i := 0; i < 6; i++ {

				var h = <- in_multiHash
				a[h.i] = h.value
			}

			var result = strings.Join(a[:], "")
			//fmt.Printf("%s MultiHash result: %s\n", s,result)
			out <- result

		}(wg, multiHash_channel)

		for i := 0; i < 6; i++ {

			go func(i int, out_multiHash chan hash) {
				var value = DataSignerCrc32(strconv.Itoa(i) + s)

				//fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", s,i, value)

				out_multiHash <- hash{ i, value}

			}(i, multiHash_channel)
		}
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {

	/*
	CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
	*/
	var s []string
	for data := range in {
		s = append(s, data.(string))
	}

	sort.Strings(s)
	var result = strings.Join(s, "_")
	//fmt.Printf("CombineResults %s\n", result)
	out <- result
}