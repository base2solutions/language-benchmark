// package main calculates the dotproduct of two arrays
// see godoc for main() for use
package main

import (
    "bufio"
    "os"
    "fmt"
    "strings"
    "strconv"
    "encoding/binary"
    "sync"
    "time"
)

const ARR_SIZE = 1500

// dotprod calculates the dot product of two arrays; this is
// multi-thread optimized (start, end).  to use as single thread
// set start to 0 and end to ARR_SIZE
func dotprod(input1, input2, output *[ARR_SIZE][ARR_SIZE]uint32,
             start, end int,
             wg *sync.WaitGroup) {
    for i := start; i < end; i++ {
        for j := 0; j < ARR_SIZE; j++ {
            var sum uint32
            for k := 0; k < ARR_SIZE; k++ {
                sum += (input1[i][j] * input2[j][k])
            }
            output[i][j] = sum
        }
    }
    wg.Done()
}

// dot_product.go takes new-line separated command line arguments
// file_a (fully qualified path)
// file_b (fully qualified path)
// number of threads (int, should be a factor of ARR_SIZE)
// it returns the elapsed time taken to calculate the dotproduct
// for the given arrays and number of threads
func main() {
    reader := bufio.NewReader(os.Stdin)
    file1_raw, _ := reader.ReadString('\n')
    file1_stripped := strings.TrimSpace(file1_raw)
    f1, _ := os.Open(file1_stripped)
    file2_raw, _ := reader.ReadString('\n')
    file2_stripped := strings.TrimSpace(file2_raw)
    f2, _ := os.Open(file2_stripped)
    var ints1 [ARR_SIZE][ARR_SIZE]uint32
    var ints2 [ARR_SIZE][ARR_SIZE]uint32
    for i := 0; i < ARR_SIZE; i++ {
        for j := 0; j < ARR_SIZE; j++ {
            b1 := make([]byte, 4)
            b2 := make([]byte, 4)
            f1.Read(b1)
            f2.Read(b2)
            this_int1 := binary.LittleEndian.Uint32(b1)
            this_int2 := binary.LittleEndian.Uint32(b2)
            ints1[i][j] = this_int1
            ints2[i][j] = this_int2
        }
    }
    var results [ARR_SIZE][ARR_SIZE]uint32
    threads_raw, _ := reader.ReadString('\n')
    threads_stripped := strings.TrimSpace(threads_raw)
    threads, _ := strconv.Atoi(threads_stripped)
    //start timer here -- all inputs have been read
    time_start := time.Now()
    var wg sync.WaitGroup
    wg.Add(threads)
    for t := 0; t < threads; t ++ {
        go_start := t * (ARR_SIZE/threads)
        go_end := (t+1) * (ARR_SIZE/threads)
        go dotprod(&ints1, &ints2, &results, go_start, go_end, &wg)
    }
    wg.Wait()
    duration := time.Since(time_start)
    fmt.Printf("sec: %f\n", duration.Seconds())
}
