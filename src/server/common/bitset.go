package common

import "strconv"

type BitSet uint64

func (lang BitSet) Contains(other BitSet) bool {
	return lang&other != 0
}

func (lang BitSet) With(other BitSet) BitSet {
	return BitSet(lang | other)
}

func (lang BitSet) Without(other BitSet) BitSet {
	return BitSet(lang & ^other)
}

func (lang BitSet) Toggle(other BitSet) BitSet {
	if lang.Contains(other) {
		return lang.Without(other)
	}
	return lang.With(other)
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

func (lang BitSet) String() string {
	return strconv.FormatUint(uint64(lang), 10)
}
