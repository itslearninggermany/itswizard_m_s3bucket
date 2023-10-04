package itswizard_m_s3bucket

import (
	"time"
)

type timeSlice []time.Time

// Functions for sort
func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }
