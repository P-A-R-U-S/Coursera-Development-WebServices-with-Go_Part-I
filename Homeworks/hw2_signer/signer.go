package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
func ExecutePipeline(jobs ...job)  {
	var in chan interface{}
	wg := &sync.WaitGroup{}
	for _, f := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func(plf func(input, output chan interface{}), i, o chan interface{}, wg *sync.WaitGroup) {
			plf(i, o)
			defer close(o)
			defer wg.Done()
		}(f, in, out, wg)
		in = out
	}
	wg.Wait()
}
*/

func ExecutePipeline(jobs ...job)  {

}

/*
func SingleHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for d := range in {
		v, ok := d.(int)
		if !ok {
			fmt.Println("cant convert result data to string")
		}
		data := strconv.Itoa(v)
		cv1 := make(chan string)
		wg.Add(1)
		go func(vdata string, out_cv1 chan<- string, wg1 *sync.WaitGroup, m *sync.Mutex) {

			m.Lock()
			md5 := DataSignerMd5(vdata)
			m.Unlock()

			out_cv1 <- DataSignerCrc32(md5)
			defer wg1.Done()
			defer close(out_cv1)
		}(data, cv1, wg, mutex)

		cv2 := make(chan string)
		wg.Add(1)
		go func(vdata string, out_cv2 chan<- string, wg2 *sync.WaitGroup) {

			out_cv2 <- DataSignerCrc32(vdata)
			defer wg2.Done()
			defer close(out_cv2)
		}(data, cv2, wg)

		wg.Add(1)
		go func(out_cv1 <-chan string,  out_cv2 <-chan string, wg3 *sync.WaitGroup) {
			out <- <- out_cv2 + "~" + <- out_cv1
			defer wg3.Done()
		}(cv1, cv2, wg)
	}

	wg.Wait()
}
*/

func SingleHash(in, out chan interface{}){

}

/*
func MultiHash(in, out chan interface{}){

	type rs struct {
		index int
		crc32 string
	}
	wg := &sync.WaitGroup{}
	for d := range in {
		v, ok := d.(string)
		if !ok {
			fmt.Println("cant convert result data to string")
		}

		cv1 := make(chan rs)
		for i := 0; i<=5;i++ {
			wg.Add(1)
			go func(ci int, vdata string, out_rs chan<- rs,  wg *sync.WaitGroup) {
				out_rs <- rs{ci, DataSignerCrc32(strconv.Itoa(ci) + vdata)}
				defer wg.Done()
			}(i, v, cv1, wg)
		}

		wg.Add(1)
		go func(vdata string, out_rs <-chan rs,) {
			var results [6]string
			for i := 0; i<=5;i++ {
				r := <- out_rs
				results[r.index] = r.crc32
			}
			close(cv1)
			out <- strings.Join(results[:],"")
			defer wg.Done()
		}(v, cv1)
	}
	wg.Wait()
}
*/

func MultiHash(in, out chan interface{}){

}

func CombineResults(in, out chan interface{}){
	var data []string
	for d := range in {
		v, ok := d.(string)
		if !ok {
			fmt.Println("cant convert result data to string")
		}
		data = append(data, v)
	}
	sort.Strings(data)
	out <- strings.Join(data, "_")
}