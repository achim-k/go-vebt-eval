package main

import (
	"fmt"
	"github.com/achimk1704/go-vebt"
	"time"
	"math"
	"math/rand"
)

type vebEval struct {
	m int
	insert, delete time.Duration
}

func main() {
	rand.Seed(time.Now().UnixNano())

	evalTimes := []vebEval{}


	// measure time + space for different universe sizes (u = 2^i)
	for i := 1; i <= 10; i++ {
		u := int(math.Pow(2, float64(i)))
		times := vebEval{m: i}

		// Create keys to be used for operations
		keys := []int{}
		for k := 0; k < u; k++ {
			keys = append(keys, k)
		}
		// shuffle keys
		shuffle(keys)

		// Create tree
		V := vebt.CreateTree(u)

		// Measure insert time (averaged for all keys)
		times.insert = insertTime(V, keys)
		times.delete = deleteTime(V, keys)
		evalTimes = append(evalTimes, times)
	}


	fmt.Printf("u = 2^x\tINSERT (ns)\tDELETE (ns)\n")
	for i := 0; i < len(evalTimes); i++ {
		t := evalTimes[i]
		fmt.Printf("%v\t%v\t\t%v\t\t\n", t.m, t.insert.Nanoseconds(), t.delete.Nanoseconds())
	}



	fmt.Println("test done")
	
}

func insertTime(V *vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Insert(keys[i])
	}
	return time.Since(start)
}

func deleteTime(V *vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Delete(keys[i])
	}
	return time.Since(start)
} 



func shuffle(a []int) {
    for i := range a {
        j := rand.Intn(i + 1)
        a[i], a[j] = a[j], a[i]
    }
}