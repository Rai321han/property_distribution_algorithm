// package main

// import (
// 	"math"
// )

// type Ratio struct {
// 	Key   string
// 	Value float64
// }

// // propertyDistribution is a function that distributes a limited number of slots among partners based on their specified ratios and priorities, while also considering the available slots for each partner in the database.
// //
// // It takes a ratio map (partner key to ratio), a priority list of partner keys, a db count map (partner key to available slots), and a limit on total slots to distribute.
// //
// // It returns a list of partner keys representing the allocated slots, ordered by priority.
// //
// // example usage:
// //
// //	ratio := map[string]float64{"partnerA": 50, "partnerB": 30, "partnerC": 20}
// //	priority := []string{"partnerA", "partnerB", "partnerC"}
// //	dbcount := map[string]int{"partnerA": 5, "partnerB": 3, "partnerC": 2}
// //	limit := 7
// //	result := propertyDistribution(ratio, priority, dbcount, limit)
// //
// // How it works:
// //
// //  1. Remove partners that have 0 db count and recalculate the ratio
// //     based on the remaining partners' total share.
// //
// //  2. Calculate the initial allocation for each partner using the
// //     available ratio and the given limit.
// //
// //  3. Sum the initial allocations and determine extras
// //     (extras = totalAllocated - limit).
// //
// //  4. If extras exist, subtract from the highest priority partner first, if highest priority partner has 1, then subtract from next higher and so on. If extras still remains or all partner has 1 allocaiton, then continue subtraction from lowest priority to highest partners until the
// //     extras are removed.
// //
// //  5. Check for exhausted partners (where db count < allocated count)
// //     and calculate the remaining slots.
// //
// //  6. Distribute the remaining slots to non-exhausted partners based
// //     on priority until all slots are filled or no partners remain.
// //
// //  7. If total allocated count exceeds the limit, remove in round robin manner from lowest priority to highest until total allocated count is equal to limit.
// //
// //  8. Return the final allocation as a list of partner keys, ordered by priority.
// //
// // Constraints:
// //
// //	minimum allocation is 1, maximum is db count
// //	1 <= limit
// //	dbcount[partner] >= 0
// //	ratio[partner] >= 0

// func propertyDistribution(ratio map[string]float64, priority []string, dbcount map[string]int, limit int) []string {

// 	// length of ratio, priority and dbcount should be same.
// 	ratioLength := len(ratio)
// 	priorityLength := len(priority)
// 	dbcountLength := len(dbcount)

// 	if ratioLength != priorityLength || ratioLength != dbcountLength {
// 		return []string{}
// 	}

// 	// These maps will hold the active partners after filtering out those with 0 db count, and their corresponding ratios and db counts.
// 	activeRatio := make(map[string]float64)
// 	activeDB := make(map[string]int)
// 	removed := make(map[string]bool)
// 	var totalRatio float64

// 	// Remove that shares that has 0 db count, calculate based on remaining sum shares
// 	for _, key := range priority {
// 		if dbcount[key] == 0 {
// 			removed[key] = true
// 			continue
// 		}
// 		activeRatio[key] = ratio[key]
// 		activeDB[key] = dbcount[key]
// 		totalRatio += ratio[key]
// 	}

// 	// remove partner from priority if it is removed in previous step
// 	var updatedPriority []string
// 	for _, key := range priority {
// 		if !removed[key] {
// 			updatedPriority = append(updatedPriority, key)
// 		}
// 	}

// 	// Calculate the initial counts from ratio and limit, and sum the total allocated count for all partners.
// 	initialAllocationCount := 0
// 	partnerWiseCount := make(map[string]int)
// 	for key, value := range activeRatio {
// 		c := int(math.Ceil((value / totalRatio) * float64(limit)))
// 		partnerWiseCount[key] = c
// 		initialAllocationCount += c
// 	}

// 	// calculate extras after initial rounding
// 	extras := initialAllocationCount - limit

// 	// subtract from highest priority if extras are there and highest priority has more than 1 share
// 	if extras > 0 {
// 		// if highest priority partner has only 1 then move toward next higher
// 		for _, key := range updatedPriority {
// 			if partnerWiseCount[key] > 1 {
// 				partnerWiseCount[key]--
// 				extras--
// 				break
// 			}
// 		}
// 	}

// 	// untill extras are removed, subtract from lowest priority if it has more than 1 share, then move to higher priority
// 	exhaustedAfterRemoved := len(partnerWiseCount)
// 	for extras > 0 && exhaustedAfterRemoved > 0 {
// 		for i := len(updatedPriority) - 1; i >= 0 && extras > 0; i-- {
// 			if partnerWiseCount[updatedPriority[i]] > 1 {
// 				partnerWiseCount[updatedPriority[i]]--
// 				extras--
// 			}

// 			if partnerWiseCount[updatedPriority[i]] == 1 {
// 				exhaustedAfterRemoved--
// 			}
// 		}
// 	}

// 	// Calculate total spots remains and total allocations
// 	remainingSlots := 0
// 	exhaustedpartners := make(map[string]bool)
// 	for key, count := range partnerWiseCount {
// 		if activeDB[key] <= count {
// 			remainingSlots += count - activeDB[key]
// 			partnerWiseCount[key] = activeDB[key]
// 			exhaustedpartners[key] = true
// 		}
// 	}

// 	var queue []string

// 	// Distribute remaining slots to non-exhausted partners based on priority until all slots are filled or no partners remain.
// 	if remainingSlots > 0 {
// 		for i := len(updatedPriority) - 1; i >= 0; i-- {
// 			key := updatedPriority[i]
// 			if !exhaustedpartners[key] {
// 				queue = append(queue, key)
// 			}
// 		}

// 		for remainingSlots > 0 && len(queue) > 0 {
// 			key := queue[0]
// 			queue = queue[1:]

// 			partnerWiseCount[key]++
// 			remainingSlots--

// 			if partnerWiseCount[key] < activeDB[key] {
// 				queue = append(queue, key)
// 			}
// 		}
// 	}

// 	// limit < partners
// 	// only possible when all partners have at most 1 allocaiton
// 	// in that case we will remove in round robin manner from lowest priority to highest until total allocated count is equal to limit.
// 	// that means we need to remove (partners - limit) partners from the parterWise map
// 	if len(partnerWiseCount) > limit {
// 		toRemove := len(partnerWiseCount) - limit
// 		// remove from lowest priority to highest until we removed required number of partners
// 		for i := len(updatedPriority) - 1; i >= 0 && toRemove > 0; i-- {
// 			key := updatedPriority[i]
// 			if _, exists := partnerWiseCount[key]; exists {
// 				delete(partnerWiseCount, key)
// 				toRemove--
// 			}
// 		}
// 	}

// 	return generateSequence(partnerWiseCount, updatedPriority)
// }

// func generateSequence(partnerWiseCount map[string]int, priority []string) []string {
// 	var result []string

// 	for len(partnerWiseCount) > 0 {
// 		for _, key := range priority {
// 			_, exists := partnerWiseCount[key]
// 			if exists && partnerWiseCount[key] > 0 {
// 				result = append(result, key)
// 				partnerWiseCount[key]--
// 			}

// 			if partnerWiseCount[key] == 0 {
// 				delete(partnerWiseCount, key)
// 			}
// 		}
// 	}
// 	return result
// }

package main

import (
	"math"
)

// propertyDistribution distributes a limited number of slots among partners
// based on their specified ratios and priorities, capped by each partner's
// available db count.
//
// Parameters:
//   - ratio:    partner key → share weight (ignored when dbcount is 0)
//   - priority: partner keys ordered highest → lowest priority
//   - dbcount:  partner key → available slots in the database
//   - limit:    total slots to distribute
//
// Returns a slot sequence ordered by priority (round-robin across partners).
//
// Algorithm:
//  1. Filter out partners with dbcount == 0 and normalise ratios.
//  2. Ceiling-allocate each partner's initial share from `limit`.
//  3. Remove over-allocation (extras) from ceiling rounding:
//     subtract from the highest-priority partner first (if > 1),
//     then sweep lowest→highest until extras reach zero.
//  4. Cap each partner at its dbcount; collect freed slots.
//  5. Redistribute freed slots to non-exhausted partners, highest priority first.
//  6. If active partners > limit (all at count 1), drop the lowest-priority ones.
//  7. Emit the final sequence via round-robin in priority order.
//
// Constraints:
//
//	limit >= 1
//	dbcount[p] >= 0, ratio[p] >= 0
//	minimum allocation per active partner: 1
func propertyDistribution(
	ratio map[string]float64,
	priority []string,
	dbcount map[string]int,
	limit int,
) []string {
	if len(ratio) != len(priority) || len(ratio) != len(dbcount) {
		return []string{}
	}

	// Step 1 – drop partners with no available slots.
	activePriority, activeRatio, activeDB := filterActive(priority, ratio, dbcount)
	if len(activePriority) == 0 {
		return []string{}
	}

	// Step 2 – ceiling-allocate initial shares.
	allocated := ceilAllocate(activePriority, activeRatio, limit)

	// Step 3 – remove over-allocation caused by ceiling rounding.
	removeExtras(allocated, activePriority, limit)

	// Step 4 – cap at dbcount and collect freed slots.
	exhausted, freed := capAtDBCount(allocated, activeDB)

	// Step 5 – give freed slots back to non-exhausted partners.
	if freed > 0 {
		redistributeFreed(allocated, activePriority, activeDB, exhausted, freed)
	}

	// Step 6 – if partners outnumber limit (all at 1), drop lowest-priority.
	if len(allocated) > limit {
		dropLowestPriority(allocated, activePriority, len(allocated)-limit)
	}

	// Step 7 – emit round-robin sequence.
	return buildSequence(allocated, activePriority)
}

// filterActive returns only partners whose dbcount > 0, preserving order.
func filterActive(
	priority []string,
	ratio map[string]float64,
	dbcount map[string]int,
) (activePriority []string, activeRatio map[string]float64, activeDB map[string]int) {
	activeRatio = make(map[string]float64, len(priority))
	activeDB = make(map[string]int, len(priority))

	for _, key := range priority {
		if dbcount[key] == 0 {
			continue
		}
		activePriority = append(activePriority, key)
		activeRatio[key] = ratio[key]
		activeDB[key] = dbcount[key]
	}
	return
}

// ceilAllocate assigns each partner at least 1 slot using ceiling division so
// that the sum of initial allocations is >= limit.
func ceilAllocate(priority []string, activeRatio map[string]float64, limit int) map[string]int {
	var totalRatio float64
	for _, key := range priority {
		totalRatio += activeRatio[key]
	}

	allocated := make(map[string]int, len(priority))
	for _, key := range priority {
		share := activeRatio[key] / totalRatio
		count := int(math.Ceil(share * float64(limit)))
		if count < 1 {
			count = 1
		}
		allocated[key] = count
	}
	return allocated
}

// removeExtras reduces the over-allocation caused by ceiling rounding.
//
// It first tries the highest-priority partner (preserving high-priority slots),
// then sweeps lowest→highest for any remainder.
// No partner is reduced below 1.
func removeExtras(allocated map[string]int, priority []string, limit int) {
	extras := sumValues(allocated) - limit
	if extras == 0 {
		return
	}

	// Pass 1: try the single highest-priority partner.
	for _, key := range priority {
		if allocated[key] > 1 {
			allocated[key]--
			extras--
			break
		}
	}

	// Pass 2: sweep lowest→highest until extras are gone.
	exhausted := 0
	for extras > 0 && exhausted < len(allocated) {
		for i := len(priority) - 1; i >= 0 && extras > 0; i-- {
			if allocated[priority[i]] > 1 {
				allocated[priority[i]]--
				extras--
			}
			if allocated[priority[i]] == 1 {
				exhausted++
			}
		}
	}
}

// capAtDBCount ensures no partner is allocated more than its db count.
// It returns which partners are now exhausted and the total freed slots.
func capAtDBCount(allocated, activeDB map[string]int) (exhausted map[string]bool, freed int) {
	exhausted = make(map[string]bool)
	for key, count := range allocated {
		if count >= activeDB[key] {
			freed += count - activeDB[key]
			allocated[key] = activeDB[key]
			exhausted[key] = true
		}
	}
	return
}

// redistributeFreed hands freed slots one at a time to non-exhausted partners,
// highest priority first. A partner leaves the rotation once it reaches its
// db count.
func redistributeFreed(
	allocated map[string]int,
	priority []string,
	activeDB map[string]int,
	exhausted map[string]bool,
	freed int,
) {
	// Seed the queue lowest->highest priority.
	queue := make([]string, 0, len(priority))
	for i := len(priority) - 1; i >= 0; i-- {
		key := priority[i]
		if !exhausted[key] {
			queue = append(queue, key)
		}
	}

	for freed > 0 && len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]

		allocated[key]++
		freed--

		if allocated[key] < activeDB[key] {
			queue = append(queue, key) // still has room; re-enqueue
		}
	}
}

// dropLowestPriority removes n lowest-priority partners from the allocation
// map. Used only when every active partner has exactly 1 slot but there are
// more partners than the limit allows.
func dropLowestPriority(allocated map[string]int, priority []string, n int) {
	for i := len(priority) - 1; i >= 0 && n > 0; i-- {
		if _, ok := allocated[priority[i]]; ok {
			delete(allocated, priority[i])
			n--
		}
	}
}

// buildSequence emits all allocated partner slots in round-robin priority
// order (highest priority partner first in each round).
func buildSequence(allocated map[string]int, priority []string) []string {
	// Work on a copy so we don't mutate the caller's map.
	counts := make(map[string]int, len(allocated))
	for k, v := range allocated {
		counts[k] = v
	}

	var result []string
	for len(counts) > 0 {
		for _, key := range priority {
			if counts[key] <= 0 {
				continue
			}
			result = append(result, key)
			counts[key]--
			if counts[key] == 0 {
				delete(counts, key)
			}
		}
	}
	return result
}

// sumValues returns the sum of all values in an int map.
func sumValues(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}
