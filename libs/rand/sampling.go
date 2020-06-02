package rand

import (
	"fmt"
	s "sort"
)

// Interface for performing weighted deterministic random selection.
type Candidate interface {
	Priority() uint64
	LessThan(other Candidate) bool
	SetWinPoint(winPoint uint64)
}

const (
	// Casting a larger number than this as float64 can result in that lower bytes can be truncated
	MaxFloat64Significant = uint64(0x1FFFFFFFFFFFFF)
)

// Select a specified number of candidates randomly from the candidate set based on each priority. This function is
// deterministic and will produce the same result for the same input.
//
// Inputs:
// seed - 64bit integer used for random selection.
// candidates - A set of candidates. The order is disregarded.
// sampleSize - The number of candidates to select at random.
// totalPriority - The exact sum of the priorities of each candidate.
//
// Returns:
// samples - A randomly selected candidate from a set of candidates. NOTE that the same candidate may have been
// selected in duplicate.
func RandomSamplingWithPriority(
	seed uint64, candidates []Candidate, sampleSize int, totalPriority uint64) (samples []Candidate) {

	// generates a random selection threshold for candidates' cumulative priority
	thresholds := make([]uint64, sampleSize)
	for i := 0; i < sampleSize; i++ {
		// calculating [gross weights] × [(0,1] random number]
		thresholds[i] = randomThreshold(&seed, totalPriority)
	}
	s.Slice(thresholds, func(i, j int) bool { return thresholds[i] < thresholds[j] })

	// generates a copy of the set to keep the given array order
	candidates = sort(candidates)

	// extract candidates with a cumulative priority threshold
	samples = make([]Candidate, sampleSize)
	cumulativePriority := uint64(0)
	undrawn := 0
	for _, candidate := range candidates {
		for thresholds[undrawn] < cumulativePriority+candidate.Priority() {
			samples[undrawn] = candidate
			undrawn++
			if undrawn == len(samples) {
				return
			}
		}
		cumulativePriority += candidate.Priority()
	}

	// This step is performed if and only if the parameter is invalid. The reasons are as stated in the message:
	actualTotalPriority := uint64(0)
	for i := 0; i < len(candidates); i++ {
		actualTotalPriority += candidates[i].Priority()
	}
	panic(fmt.Sprintf("Either the given candidate is an empty set, the actual cumulative priority is zero,"+
		" or the total priority is less than the actual one; totalPriority=%d, actualTotalPriority=%d,"+
		" seed=%d, sampleSize=%d, undrawn=%d, threshold[%d]=%d, len(candidates)=%d",
		totalPriority, actualTotalPriority, seed, sampleSize, undrawn, undrawn, thresholds[undrawn], len(candidates)))
}

func moveWinnerToLast(candidates []Candidate, winner int) {
	winnerCandidate := candidates[winner]
	copy(candidates[winner:], candidates[winner+1:])
	candidates[len(candidates)-1] = winnerCandidate
}

func randomThreshold(seed *uint64, total uint64) uint64 {
	return uint64(float64(nextRandom(seed)&MaxFloat64Significant) /
		float64(MaxFloat64Significant+1) * float64(total))
}

// `RandomSamplingWithoutReplacement` elects winners among candidates without replacement
// so it updates rewards of winners. This function continues to elect winners until the both of two
// conditions(minSamplingCount, minPriorityPercent) are met.
func RandomSamplingWithoutReplacement(
	seed uint64, candidates []Candidate, minSamplingCount int, minPriorityPercent uint, winPointUnit uint64) (
	winners []Candidate) {

	if len(candidates) < minSamplingCount {
		panic(fmt.Sprintf("The number of candidates(%d) cannot be less minSamplingCount %d",
			len(candidates), minSamplingCount))
	}

	if minPriorityPercent > 100 {
		panic(fmt.Sprintf("minPriorityPercent must be equal or less than 100: %d", minPriorityPercent))
	}

	totalPriority := sumTotalPriority(candidates)
	if totalPriority > MaxFloat64Significant {
		// totalPriority will be casting to float64, so it must be less than 0x1FFFFFFFFFFFFF(53bits)
		panic(fmt.Sprintf("total priority cannot exceed %d: %d", MaxFloat64Significant, totalPriority))
	}

	priorityThreshold := totalPriority * uint64(minPriorityPercent) / 100
	candidates = sort(candidates)
	winnersPriority := uint64(0)
	losersPriorities := make([]uint64, len(candidates))
	winnerNum := 0
	for winnerNum < minSamplingCount || winnersPriority < priorityThreshold {
		threshold := randomThreshold(&seed, totalPriority-winnersPriority)
		cumulativePriority := uint64(0)
		found := false
		for i, candidate := range candidates[:len(candidates)-winnerNum] {
			if threshold < cumulativePriority+candidate.Priority() {
				moveWinnerToLast(candidates, i)
				winnersPriority += candidate.Priority()
				losersPriorities[winnerNum] = totalPriority - winnersPriority
				winnerNum++
				found = true
				break
			}
			cumulativePriority += candidate.Priority()
		}

		if !found {
			panic(fmt.Sprintf("Cannot find random sample. winnerNum=%d, minSamplingCount=%d, "+
				"winnersPriority=%d, priorityThreshold=%d, totalPriority=%d, threshold=%d",
				winnerNum, minSamplingCount, winnersPriority, priorityThreshold, totalPriority, threshold))
		}
	}
	compensationProportions := make([]float64, winnerNum)
	for i := winnerNum - 2; i >= 0; i-- { // last winner doesn't get compensation reward
		compensationProportions[i] = compensationProportions[i+1] + 1/float64(losersPriorities[i])
	}
	winners = candidates[len(candidates)-winnerNum:]
	for i, winner := range winners {
		winner.SetWinPoint(winPointUnit +
			uint64(float64(winner.Priority())*compensationProportions[i]*float64(winPointUnit)))
	}
	return winners
}

func sumTotalPriority(candidates []Candidate) (sum uint64) {
	for _, candi := range candidates {
		if candi.Priority() == 0 {
			panic("candidate(%d) priority must not be 0")
		}
		sum += candi.Priority()
	}
	return sum
}

// SplitMix64
// http://xoshiro.di.unimi.it/splitmix64.c
//
// The PRNG used for this random selection:
//   1. must be deterministic.
//   2. should easily portable, independent of language or library
//   3. is not necessary to keep a long period like MT, since there aren't many random numbers to generate and
//      we expect a certain amount of randomness in the seed.
func nextRandom(rand *uint64) uint64 {
	*rand += uint64(0x9e3779b97f4a7c15)
	var z = *rand
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// sort candidates in descending priority and ascending nature order
func sort(candidates []Candidate) []Candidate {
	temp := make([]Candidate, len(candidates))
	copy(temp, candidates)
	s.Slice(temp, func(i, j int) bool {
		if temp[i].Priority() != temp[j].Priority() {
			return temp[i].Priority() > temp[j].Priority()
		}
		return temp[i].LessThan(temp[j])
	})
	return temp
}
