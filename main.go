package main

import (
	"fmt"
	"math"
)

type Ratio struct {
	Key   string
	Value float64
}

// After rounding
//1.1 If +1 -> subtract from highest priority
//1.2 If +2 or mote -> subtract 1 from highest priority, 1 from from lowest priority, then lower -> goes on

// DB check
//2.1 If missing, addition -> lower priority to higher priority

// Solution
// Rounding -> use ceil -> always sum >= limit -> rule 1.1 or rule 1.2

// vda
// How Vendor Distribution Algorithm works:
// Step 1: Remove that shares that has 0 db count, calculate based on remaining sum shares.
// Step 2: Calculate the initial counts from available ratio and limit using ceiling to ensure we meet the limit.
// Step 3: Calculate the total of initial counts and determine if there are extras (total - limit).
// Step 4: If there are extras, subtract from the highest priority vendor first, then move to lower priorities until extras are removed.
// Step 5: Check for exhausted vendors (where db count is less than allocated count) and calculate remaining slots.
// Step 6: Distribute remaining slots to non-exhausted vendors based on priority until all slots are filled or no more vendors are available.
func vda(ratio map[string]float64, priority []string, dbcount map[string]int, limit int) []string {

	// Step: DB count check
	var availableRatio float64

	// Remove that shares that has 0 db count, calculate based on remaining sum shares
	for key, vendor := range dbcount {
		if vendor == 0 {
			delete(dbcount, key)
			delete(ratio, key)
			continue
		}
		availableRatio += ratio[key]
	}

	// Step: Vendor Wise Initial Count
	vendorCount := make(map[string]int)

	// Calculate the initial counts from ratio and limit
	for key, value := range ratio {
		vendorCount[key] = int(math.Ceil((value / availableRatio) * float64(limit)))
	}

	// sum of initial counts
	vendorCountTotal := 0
	for _, count := range vendorCount {
		vendorCountTotal += count
	}

	// calculate extras after initial rounding
	extras :=  vendorCountTotal - limit

	// subtract from highest priority if extras are there and highest priority has more than 1 share
	if extras > 0 && vendorCount[priority[0]] > 1 {
		vendorCount[priority[0]]--
		extras--
	}

	// untill extras are removed, subtract from lowest priority if it has more than 1 share, then move to higher priority
	for extras > 0 {
		for i := len(priority) - 1; i >= 0 && extras > 0; i-- {
			if vendorCount[priority[i]] > 1 {
				vendorCount[priority[i]]--
				extras--
			}
		}
	}

	// Calculate total spots remains and total allocations
	remainingSlots := 0
	exhaustedVendors := make(map[string]bool)
	for key, count := range vendorCount {
		if dbcount[key] < count {
			remainingSlots += count - dbcount[key]
			vendorCount[key] = dbcount[key]
			exhaustedVendors[key] = true
		}
	}

	var queue []string

	for i := len(priority)-1; i >=0; i-- {
		key := priority[i]
		if !exhaustedVendors[key] {
			queue = append(queue, key)
		}
	}
	
	for remainingSlots > 0 && len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]

		vendorCount[key]++
		remainingSlots--

		if vendorCount[key] < dbcount[key] {
			queue = append(queue, key)
		}
	}
	return result(vendorCount, priority)
}

func result(vendorCount map[string]int, priority []string) []string {
	fmt.Printf("Vendor Count: %v\n", vendorCount)
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
	ratio := map[string]float64{
		"12": 30,
		"24": 20,
		"11": 50,
	}

	priority := []string{"24", "12", "11"}

	dbCount := map[string]int{
		"11": 1,
		"12": 1,
		"24": 2,
	}

	limit := 8

	result := vda(ratio, priority, dbCount, limit)
	fmt.Printf("Result: %v\n", result)
}
