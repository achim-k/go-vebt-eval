package main

import (	
	"fmt"
	"github.com/achimk1704/go-vebt"
	"time"
	"math"
	"math/rand"
	"strings"
	"strconv"
	"unsafe"
)

type vebEval struct {
	m int
	size int
	insert, delete, isMember, successor, predecessor, min, max int
}

type VEB struct {
	u, min, max int    //universe size, min-, max value
	summary     *VEB   //pointer to summary
	cluster     []*VEB // array of pointers to each child cluster
}

func main() {

	// Default values
	maxM := 16
	runs := 50
	testkeys := ""
	treeFullness := ""
	var operation string
	randomTestKeys := 0

	var c VEB
	var V VEB
	V.u, V.min, V.max = 1,1,1
	V.summary = &c
	V.cluster = append(V.cluster, &c)

	fmt.Println("Sizeof VEB struct in Bytes:", unsafe.Sizeof(c), unsafe.Sizeof(V.cluster))
	
	
	fmt.Printf("Maximum tree size m (u = 2^m)?\n")
	fmt.Scanf("%d", &maxM)
	fmt.Printf("operation? Choices: \n[count, insert, delete, member, successor, predecessor, min, max]\n")
	fmt.Scanf("%s", &operation)
	fmt.Printf("Number of test keys for each operation? [e.g. 100 or 50%%]\n")	
	fmt.Scanf("%s", &testkeys)
	fmt.Printf("Create keys randomly=1 or sequentially=0?\n")	
	fmt.Scanf("%d", &randomTestKeys)

	fmt.Printf("How many runs per tree size (reduce effect of randomized keys)?\n")	
	fmt.Scanf("%d", &runs)

		fmt.Printf("Fullness ratio of tree? [e.g. 100 or 50%%]\n")
	fmt.Scanf("%s", &treeFullness)
	
	fmt.Printf("Avg time [ns] for each operation with random generated testkeys for %v runs\n", runs)
	fmt.Printf("m (u=2^m)\t%v\n", operation)
	// measure time + space for different universe sizes (u = 2^i)
	for i := 1; i <= maxM; i++ {
		u := int(math.Pow(2, float64(i)))
		keyNo := 100
		treeInitKeyNo := 100

		if strings.Contains(testkeys, "%") {
			//relative to u
			ratio, _ := strconv.ParseInt(testkeys[:strings.Index(testkeys, "%")], 0, 0)
			keyNo = int(u * int(ratio)/100)
			if keyNo <= 0 {
				keyNo = 1
			}
		} else {
			ratio, _ := strconv.ParseInt(testkeys, 0, 0)
			keyNo = int(ratio)
		}
		if strings.Contains(treeFullness, "%") {
			//relative to u
			treeFullnessRatio, _ := strconv.ParseInt(treeFullness[:strings.Index(treeFullness, "%")], 0, 0)
			treeInitKeyNo = int(u * int(treeFullnessRatio)/100)
			if treeInitKeyNo <= 0 {
				treeInitKeyNo = 1
			}
		} else {
			treeFullnessRatio, _ := strconv.ParseInt(treeFullness, 0, 0)
			treeInitKeyNo = int(treeFullnessRatio)
		}

		// Create tree
		V := vebt.CreateTree(u)
		timeSum := 0

		// run more than once, otherwise random keys have to much influence
		for run := 0; run < runs; run++ {
			keys := []int{}

			// Fill tree with random keys to given number before performing operation
			if operation != "count" {
				V.Clear() // clear keys inserted before
				initKeys := createKeys(treeInitKeyNo, u, 1)
				for l := 0; l < len(initKeys); l++ {
					V.Insert(initKeys[l])
				}

				// Create keys for testing operation with
				keys = createKeys(keyNo, u, randomTestKeys)
			}
		
			switch operation {
			case "count":
				timeSum = V.Count()
			case "insert":
				timeSum += int(insertTime(V, keys).Nanoseconds())
			case "delete":
				timeSum += int(deleteTime(V, keys).Nanoseconds())	
			case "member":
				timeSum += int(isMemberTime(*V, keys).Nanoseconds())
			case "successor":
				timeSum += int(successorTime(*V, keys).Nanoseconds())	
			case "predecessor":
				timeSum += int(predecessorTime(*V, keys).Nanoseconds())	
			case "min":
				timeSum += int(minTime(*V, keys).Nanoseconds())	
			case "max":
				timeSum += int(maxTime(*V, keys).Nanoseconds())	
			}	

		}

		timeSum /= runs * keyNo
		
		fmt.Printf("%v\t\t%v\n", i, timeSum)
	}


	
	fmt.Printf("Measuring time for {operation} for different tree fullness rate\n")
	fmt.Printf("operation? Choices: \n[insert, delete, member, successor, predecessor, min, max]\n")
	fmt.Scanf("%s", &operation)
	fmt.Printf("Maximum tree size m (u = 2^m)?\n")
	fmt.Scanf("%d", &maxM)
	fmt.Printf("How many runs per tree size (reduce effect of randomized keys)?\n")	
	fmt.Scanf("%d", &runs)

	fmt.Printf("Measuring avg %v time [ns] for different tree fullness ratios for %v testkeys\n", strings.ToUpper(operation), testkeys)
	fmt.Printf("m\t#keys\t")
	for fillRate := 0; fillRate <= 100; fillRate += 10 {
		fmt.Printf("%v%%\t", fillRate)
	}
	fmt.Printf("\n")


	for i := 1; i <= maxM; i++ {
		
		u := int(math.Pow(2, float64(i)))
		V := vebt.CreateTree(u)
		fmt.Printf("%v\t", i)
		keyNo := 100

		if strings.Contains(testkeys, "%") {
			//relative to u
			ratio, _ := strconv.ParseInt(testkeys[:strings.Index(testkeys, "%")], 0, 0)
			keyNo = int(u * int(ratio)/100)
			if keyNo <= 0 {
				keyNo = 1
			}
		} else {
			ratio, _ := strconv.ParseInt(testkeys, 0, 0)
			keyNo = int(ratio)
		}

		fmt.Printf("%v\t", keyNo)

		for fillRate := 0; fillRate <= 100; fillRate += 10 {
			V.Clear()
			// Create number of random keys to fill tree (depending on fillRate)
			insertKeys := createKeys(int(u * fillRate / 100), u, 1)
			// Fill tree
			for i := 0; i < len(insertKeys); i++ {
				V.Insert(insertKeys[i])
			}
			timeSum := 0

			// Measure average time it takes to insert 1 random key
			for r := 0; r < runs; r++ {
				keys := createKeys(keyNo, u, randomTestKeys)
				switch operation {
				case "insert":
					timeSum += int(insertTime(V, keys).Nanoseconds())
				case "delete":
					timeSum += int(deleteTime(V, keys).Nanoseconds())	
				case "member":
					timeSum += int(isMemberTime(*V, keys).Nanoseconds())
				case "successor":
					timeSum += int(successorTime(*V, keys).Nanoseconds())	
				case "predecessor":
					timeSum += int(predecessorTime(*V, keys).Nanoseconds())	
				case "min":
					timeSum += int(minTime(*V, keys).Nanoseconds())	
				case "max":
					timeSum += int(maxTime(*V, keys).Nanoseconds())	
				}				
			}

			timeSum /= runs * keyNo
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

func createKeys(count, max int, randomized int) []int {
	keys := []int{}

	for i := 0; i < max; i++ {
		keys = append(keys, i)
	}

	//shuffle order of keys
	rand.Seed(time.Now().UnixNano())

	if randomized != 0 {
		for i := range keys {
        j := rand.Intn(i + 1)
        keys[i], keys[j] = keys[j], keys[i]
    	}
	}

	for count > len(keys) {
		rand.Seed(time.Now().UnixNano())
		keys = append(keys, rand.Intn(max))
	}

    return keys[0:count]
}
