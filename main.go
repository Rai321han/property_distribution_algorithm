package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"os"
)

type Ratio struct {
	Key   string
	Value float64
}

// vda implements the Vendor Distribution Algorithm.
//
// How it works:
//
//  1. Remove vendors that have 0 db count and recalculate the ratio
//     based on the remaining vendors' total share.
//
//  2. Calculate the initial allocation for each vendor using the
//     available ratio and the given limit. Use ceiling to ensure
//     the total can reach the limit.
//
//  3. Sum the initial allocations and determine extras
//     (extras = totalAllocated - limit).
//
//  4. If extras exist, subtract from the highest priority vendor
//     first, then continue to lower priority vendors until the
//     extras are removed.
//
//  5. Check for exhausted vendors (where db count < allocated count)
//     and calculate the remaining slots.
//
//  6. Distribute the remaining slots to non-exhausted vendors based
//     on priority until all slots are filled or no vendors remain.
//
//  7. Return the final allocation as a list of vendor keys, ordered by priority.
//
// Constraints:
//
//	vendors <= limit
func vda(ratio map[string]float64, priority []string, dbcount map[string]int, limit int) []string {

	// Step: DB count check
	var availableRatio float64

	// Create copies of the input maps to avoid modifying the original data
	ratioCopy := make(map[string]float64)
	maps.Copy(ratioCopy, ratio)

	dbCountCopy := make(map[string]int)
	maps.Copy(dbCountCopy, dbcount)

	priorityCopy := make([]string, len(priority))
	copy(priorityCopy, priority)

	// Remove that shares that has 0 db count, calculate based on remaining sum shares
	removedPriority := make(map[string]bool)
	for key, vendor := range dbCountCopy {
		if vendor == 0 {
			delete(dbCountCopy, key)
			delete(ratioCopy, key)
			removedPriority[key] = true
			continue
		}
		availableRatio += ratioCopy[key]
	}

	if len(ratioCopy) > limit {
		fmt.Println("Error: Total vendors exceed the limit.")
		return nil
	}

	// remove vendor from priority if it is removed in previous step
	var updatedPriority []string
	for _, key := range priorityCopy {
		if !removedPriority[key] {
			updatedPriority = append(updatedPriority, key)
		}
	}

	// Step: Vendor Wise Initial Count
	vendorCount := make(map[string]int)

	// Calculate the initial counts from ratio and limit
	for key, value := range ratioCopy {
		vendorCount[key] = int(math.Ceil((value / availableRatio) * float64(limit)))
	}

	// sum of initial counts
	vendorCountTotal := 0
	for _, count := range vendorCount {
		vendorCountTotal += count
	}

	// calculate extras after initial rounding
	extras := vendorCountTotal - limit

	// subtract from highest priority if extras are there and highest priority has more than 1 share
	if extras > 0 {
		for _, key := range updatedPriority {
			if vendorCount[key] > 1 {
				vendorCount[key]--
				extras--
				break
			}
		}
	}

	// untill extras are removed, subtract from lowest priority if it has more than 1 share, then move to higher priority
	for extras > 0 {
		for i := len(updatedPriority) - 1; i >= 0 && extras > 0; i-- {
			if vendorCount[updatedPriority[i]] > 1 {
				vendorCount[updatedPriority[i]]--
				extras--
			}
		}
	}

	// Calculate total spots remains and total allocations
	remainingSlots := 0
	exhaustedVendors := make(map[string]bool)
	for key, count := range vendorCount {
		if dbCountCopy[key] <= count {
			remainingSlots += count - dbCountCopy[key]
			vendorCount[key] = dbCountCopy[key]
			exhaustedVendors[key] = true
		}
	}

	var queue []string

	for i := len(updatedPriority) - 1; i >= 0; i-- {
		key := updatedPriority[i]
		if !exhaustedVendors[key] {
			queue = append(queue, key)
		}
	}

	for remainingSlots > 0 && len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]

		vendorCount[key]++
		remainingSlots--

		if vendorCount[key] < dbCountCopy[key] {
			queue = append(queue, key)
		}
	}
	return result(vendorCount, updatedPriority)
}

func result(vendorCount map[string]int, priority []string) []string {
	var result []string

	for len(vendorCount) > 0 {
		for _, key := range priority {
			_, exists := vendorCount[key]
			if exists && vendorCount[key] > 0 {
				result = append(result, key)
				vendorCount[key]--
			}

			if vendorCount[key] == 0 {
				delete(vendorCount, key)
			}
		}
	}
	return result
}

func main() {
	testCases, err := loadTestCases("test_cases.json")
	if err != nil {
		fmt.Println("Error loading test cases:", err)
		return
	}

	for _, testCase := range testCases {
		fmt.Printf("Test Case: %s\n", testCase.Description)
		result := vda(testCase.Ratio, testCase.Priority, testCase.DBCount, testCase.Limit)
		fmt.Println("Limit:", testCase.Limit)
		fmt.Printf("Ratio: %v\n", testCase.Ratio)
		fmt.Printf("Priority: %v\n", testCase.Priority)
		fmt.Printf("DB Count: %v\n", testCase.DBCount)
		fmt.Printf("Result: %v\n", result)
		fmt.Println("--------------------------------------------------")
	}
}

type TestCase struct {
	Description string
	Ratio       map[string]float64
	Priority    []string
	DBCount     map[string]int
	Limit       int
}

func loadTestCases(path string) ([]TestCase, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tests []TestCase
	if err := json.Unmarshal(file, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}
