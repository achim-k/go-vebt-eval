package main

import (	
	"fmt"
	"github.com/achimk1704/go-vebt"
	"time"
	"math"
	"math/rand"
	"strings"
)

type vebEval struct {
	m int
	size int
	insert, delete, isMember, successor, predecessor, min, max int
}

func main() {

	// Default values
	maxM := 16
	runs := 50
	
	fmt.Printf("Maximum tree size m (u = 2^m)?\n")
	fmt.Scanf("%d", &maxM)
	fmt.Printf("How many runs per tree size (reduce effect of randomized keys)?\n")	
	fmt.Scanf("%d", &runs)

	fmt.Printf("m (u=2^m)\t")
	fmt.Printf("#STRUCTS\t")
	fmt.Printf("#TESTKEYS\t")
	fmt.Printf("INSERT (ns)\t")
	fmt.Printf("SUCCESSOR (ns)\t")
	fmt.Printf("PREDECESSOR (ns)\t")
	fmt.Printf("MIN (ns)\t")
	fmt.Printf("MAX (ns)\t")
	fmt.Printf("DELETE (ns)\t")
	fmt.Printf("\n")
	// measure time + space for different universe sizes (u = 2^i)
	for i := 1; i <= maxM; i++ {
		u := int(math.Pow(2, float64(i)))
		eval := vebEval{m: i}
		keyNo := 100 // 5%
		//keyNo = u/2


		// Create tree
		V := vebt.CreateTree(u)
		V_full := vebt.CreateTree(u)

		eval.size = V.Count()

		// run more than once, otherwise random keys have to much influence
		for run := 0; run < runs; run++ {
			V.Clear() // perform measurements on empty tree
			V_full.Fill() // fill tree again (for deletion)

			keys := createRandomKeys(keyNo, u)
			// measure average insert time
			eval.insert += int(insertTime(V, keys).Nanoseconds()/int64(keyNo))
			// measure average successor time
			eval.successor += int(successorTime(*V, keys).Nanoseconds()/int64(keyNo))
			// measure average predecessor time
			eval.predecessor += int(predecessorTime(*V, keys).Nanoseconds()/int64(keyNo))
			// measure average min time
			eval.min += int(minTime(*V, keys).Nanoseconds()/int64(keyNo))
			// measure average max time
			eval.max += int(maxTime(*V, keys).Nanoseconds()/int64(keyNo))
			// measure delete time
			eval.delete += int(deleteTime(V_full, keys).Nanoseconds()/int64(keyNo))
		}

		eval.insert /= runs
		eval.successor /= runs
		eval.predecessor /= runs
		eval.min /= runs
		eval.max /= runs
		eval.delete /= runs

		fmt.Printf("%v\t\t", eval.m)
		fmt.Printf("%v\t\t", eval.size)
		fmt.Printf("%v\t\t", keyNo)
		fmt.Printf("%v\t\t", eval.insert)
		fmt.Printf("%v\t\t", eval.successor)
		fmt.Printf("%v\t\t\t", eval.predecessor)
		fmt.Printf("%v\t\t", eval.min)
		fmt.Printf("%v\t\t", eval.max)
		fmt.Printf("%v\t\t", eval.delete)

		fmt.Printf("\n")
	}


	var operation string
	fmt.Printf("Measuring time for {operation} for different tree fullness rate\n")
	fmt.Printf("operation? Choices: [insert, delete, successor, predecessor, min, max]\n")
	fmt.Scanf("%s", &operation)
	fmt.Printf("Maximum tree size m (u = 2^m)?\n")
	fmt.Scanf("%d", &maxM)
	fmt.Printf("How many runs per tree size (reduce effect of randomized keys)?\n")	
	fmt.Scanf("%d", &runs)

	fmt.Printf("Measuring avg %v time [ns] for different tree fullness ratios\n", strings.ToUpper(operation))
	fmt.Printf("m\t")
	for fillRate := 0; fillRate <= 100; fillRate += 10 {
		fmt.Printf("%v%%\t", fillRate)
	}
	fmt.Printf("\n")


	for i := 1; i <= maxM; i++ {
		
		u := int(math.Pow(2, float64(i)))
		V := vebt.CreateTree(u)
		fmt.Printf("%v\t", i)

		for fillRate := 0; fillRate <= 100; fillRate += 10 {
			V.Clear()
			// Create number of random keys to fill tree (depending on fillRate)
			insertKeys := createRandomKeys(int(u * fillRate / 100), u)
			// Fill tree
			for i := 0; i < len(insertKeys); i++ {
				V.Insert(insertKeys[i])
			}


			timeSum := 0

			// Measure average time it takes to insert 1 random key
			for r := 0; r < runs; r++ {
				keys := createRandomKeys(100, u)

				switch operation {
				case "insert":
					timeSum += int(insertTime(V, keys).Nanoseconds()/int64(len(keys)))
				case "delete":
					timeSum += int(deleteTime(V, keys).Nanoseconds()/int64(len(keys)))	
				case "successor":
					timeSum += int(successorTime(*V, keys).Nanoseconds()/int64(len(keys)))	
				case "predecessor":
					timeSum += int(predecessorTime(*V, keys).Nanoseconds()/int64(len(keys)))	
				case "min":
					timeSum += int(minTime(*V, keys).Nanoseconds()/int64(len(keys)))	
				case "max":
					timeSum += int(maxTime(*V, keys).Nanoseconds()/int64(len(keys)))	
				}				
			}

			timeSum /= runs
			fmt.Printf("%v\t", timeSum)
		}
		fmt.Printf("\n")
	}

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

func isMemberTime(V vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.IsMember(keys[i])
	}
	return time.Since(start)
}

func successorTime(V vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Successor(keys[i])
	}
	return time.Since(start)
}

func predecessorTime(V vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Predecessor(keys[i])
	}
	return time.Since(start)
}

func minTime(V vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Min()
	}
	return time.Since(start)
}

func maxTime(V vebt.VEB, keys []int) time.Duration {
	start := time.Now()
	for i := 0; i < len(keys); i++ {
		V.Max()
	}
	return time.Since(start)
}

func createRandomKeys(count, max int) []int {
	keys := []int{}

	for i := 0; i < max; i++ {
		keys = append(keys, i)
	}

	//shuffle order of keys
	rand.Seed(time.Now().UnixNano())
	for i := range keys {
        j := rand.Intn(i + 1)
        keys[i], keys[j] = keys[j], keys[i]
    }

	for count > len(keys) {
		rand.Seed(time.Now().UnixNano())
		keys = append(keys, rand.Intn(max))
	}

    return keys[0:count]
}
