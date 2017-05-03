package mml

type intSet struct {
	buckets []nodeType
}

func (s *intSet) set(i nodeType) {
	bucket, flag := i/64, i%64
	need := int(bucket + 1)
	if len(s.buckets) < need {
		if cap(s.buckets) >= need {
			s.buckets = s.buckets[:need]
		} else {
			s.buckets = s.buckets[:cap(s.buckets)]
			for len(s.buckets) < need {
				s.buckets = append(s.buckets, 0)
			}
		}
	}

	s.buckets[bucket] |= 1 << flag
}

func (s *intSet) has(i nodeType) bool {
	bucket, flag := i/64, i%64
	if nodeType(len(s.buckets)) <= bucket {
		return false
	}

	return s.buckets[bucket]&(1<<flag) != 0
}
