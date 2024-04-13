package common

import "strconv"

type BitSet uint64

const EmptyBitSet = BitSet(0)

func (bs BitSet) Contains(other BitSet) bool {
	return bs&other != 0
}

func (bs BitSet) With(other BitSet) BitSet {
	return BitSet(bs | other)
}

func (bs BitSet) Without(other BitSet) BitSet {
	return BitSet(bs & ^other)
}

func (bs BitSet) Toggle(other BitSet) BitSet {
	if bs.Contains(other) {
		return bs.Without(other)
	}
	return bs.With(other)
}

func BitSetFromString(s string, defaultValue BitSet, zeroValue BitSet) BitSet {
	n, err := strconv.ParseUint(s, 10, 3)
	if err != nil {
		return defaultValue
	}
	if n == 0 {
		return zeroValue
	}
	return BitSet(n)
}

func (bs BitSet) String() string {
	return strconv.FormatUint(uint64(bs), 10)
}

func (bs BitSet) MaskNames(names []string) []string {
	maskedNames := make([]string, 0, len(names))
	for i, name := range names {
		if bs.Contains(BitSet(1 << i)) {
			maskedNames = append(maskedNames, name)
		}
	}
	return maskedNames
}

type UnselectedItem struct {
	Value uint64
	Name  string
}

func (bs BitSet) DivideNamesByMask(names []string) ([]string, []UnselectedItem) {
	maskedNames := make([]string, 0, len(names))
	unmaskedNames := make([]UnselectedItem, 0, len(names))
	for i, name := range names {
		value := BitSet(1 << i)
		if bs.Contains(value) {
			maskedNames = append(maskedNames, name)
		} else {
			unmaskedNames = append(unmaskedNames, UnselectedItem{
				Value: uint64(value),
				Name:  name,
			})
		}
	}
	return maskedNames, unmaskedNames
}
