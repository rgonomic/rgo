package camel

import (
	"reflect"
	"testing"
)

var splitTests = []struct {
	str  string
	want []string
}{
	{str: "StatusOK", want: []string{"Status", "OK"}},
	{str: "StatusAccepted", want: []string{"Status", "Accepted"}},
	{str: "StatusNonAuthoritativeInfo", want: []string{"Status", "Non", "Authoritative", "Info"}},
	{str: "StatusNoContent", want: []string{"Status", "No", "Content"}},
	{str: "StatusIMUsed", want: []string{"Status", "IM", "Used"}},
	{str: "X509KeyPair", want: []string{"X", "509", "Key", "Pair"}},
	{str: "keyName", want: []string{"key", "Name"}},
	{str: "__err", want: []string{"err"}},
	{str: "_err", want: []string{"err"}},
	{str: "err", want: []string{"err"}},
	{str: "err_", want: []string{"err"}},
	{str: "err__", want: []string{"err"}},
	{str: "__Err", want: []string{"Err"}},
	{str: "_Err", want: []string{"Err"}},
	{str: "Err", want: []string{"Err"}},
	{str: "Err_", want: []string{"Err"}},
	{str: "Err__", want: []string{"Err"}},
	{str: "AF_X25", want: []string{"AF", "X", "25"}},
	{str: "AF_X2_5", want: []string{"AF", "X", "2", "5"}},
	{str: "AF_X2__5", want: []string{"AF", "X", "2", "5"}},
	{str: "ARPHRD_IEEE80211_RADIOTAP", want: []string{"ARPHRD", "IEEE", "80211", "RADIOTAP"}},
	{str: "_go_hdf5_H5P_DEFAULT", want: []string{"go", "hdf", "5", "H", "5", "P", "DEFAULT"}},
	{str: "IsHDF5", want: []string{"Is", "HDF", "5"}},
	{str: "NMembers", want: []string{"N", "Members"}},
	{str: "IsURLValid", want: []string{"Is", "URL", "Valid"}},
}

func TestSplit(t *testing.T) {
	for _, test := range splitTests {
		got := Split(test.str)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected split from %q: got:%q want:%q", test.str, got, test.want)
		}
	}
}

var splitterTests = []struct {
	str   string
	known []string
	want  []string
}{
	{str: "X509KeyPair", known: []string{"X509"}, want: []string{"X509", "Key", "Pair"}},
	{str: "AF_X25", known: []string{"X25"}, want: []string{"AF", "X25"}},
	{str: "ARPHRD_IEEE80211_RADIOTAP", known: []string{"IEEE80211"}, want: []string{"ARPHRD", "IEEE80211", "RADIOTAP"}},
	{str: "_go_hdf5_H5P_DEFAULT", known: []string{"hdf5", "H5P"}, want: []string{"go", "hdf5", "H5P", "DEFAULT"}},
	{str: "IsHDF5", known: []string{"HDF5"}, want: []string{"Is", "HDF5"}},
	{str: "NMembers", want: []string{"N", "Members"}},
	{str: "IsURLValid", want: []string{"Is", "URL", "Valid"}},
	{str: "ParseWithNaNOrNA", want: []string{"Parse", "With", "Na", "N", "Or", "NA"}},
	{str: "ParseWithNaNOrNA", known: []string{"NaN"}, want: []string{"Parse", "With", "NaN", "Or", "NA"}},
	{str: "ParseWithNaNorNA", want: []string{"Parse", "With", "Na", "Nor", "NA"}},
	{str: "ParseWithNaNorNA", known: []string{"NaN"}, want: []string{"Parse", "With", "NaN", "or", "NA"}},
	{str: "NaNA", known: []string{"NaN", "NA"}, want: []string{"NaN", "A"}},
	{str: "NaNA", known: []string{"NA", "NaN"}, want: []string{"Na", "NA"}},
}

func TestSplitter(t *testing.T) {
	for _, test := range splitterTests {
		s := NewSplitter(test.known)
		got := s.Split(test.str)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected split from %q: got:%q want:%q", test.str, got, test.want)
		}
	}
}
