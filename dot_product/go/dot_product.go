// package main calculates the dotproduct of two arrays
// see main() godoc for use
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
    "errors"
)

const ARR_SIZE = 1500

// check exits and prints err if exits
func check(e error) {
    if e != nil {
        panic(e)
    }
}

// dotprod calculates the dot product of two array; this is
// multi-thread optimized (start, end).  To use as a single-thread
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

// main takes new-line separated command line args
// file_a (fully qualified path)
// file_b (fully qualified path)
// number of threads (int, must be a factor of ARR_SIZE)
// main prints the elapsed time taken to calculate the dotproduct
// for the given arrays and number of threads.  It writes the result
// as bytes to /tmp/dotprodres
func main() {
    reader := bufio.NewReader(os.Stdin)

    //read in bytes from file_a, file_b
    file1_raw, err := reader.ReadString('\n')
    check(err)
    file1_stripped := strings.TrimSpace(file1_raw)
    f1, err := os.Open(file1_stripped)
    check(err)
    file2_raw, err := reader.ReadString('\n')
    check(err)
    file2_stripped := strings.TrimSpace(file2_raw)
    f2, err := os.Open(file2_stripped)
    check(err)
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

    // read in thread number
    threads_raw, err := reader.ReadString('\n')
    check(err)
    threads_stripped := strings.TrimSpace(threads_raw)
    threads, err := strconv.Atoi(threads_stripped)
    check(err)
    if ARR_SIZE % threads != 0 {
        factor_str := fmt.Sprintf("%d is not a factor of %d\n", threads, ARR_SIZE)
        factor_err := errors.New(factor_str)
        panic(factor_err)
    }

    // start timer here -- all inputs have been read
    var results [ARR_SIZE][ARR_SIZE]uint32
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

    // write output to file tmp/dotprodres
    f, err := os.Create("/tmp/dotprodres")
    check(err)
    defer f.Close()
    for i := 0; i < ARR_SIZE; i++ {
        for j := 0; j < ARR_SIZE; j++ {
            b1 := make([]byte, 4)
            binary.LittleEndian.PutUint32(b1,  results[i][j])
            _, err := f.Write(b1)
            check(err)
        }
    }
}
