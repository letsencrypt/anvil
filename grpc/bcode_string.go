// Code generated by "stringer -type=BCode bcodes.go"; DO NOT EDIT

package grpc

import "fmt"

const _BCode_name = "DNSQueryTimeoutDNSError"

var _BCode_index = [...]uint8{0, 15, 23}

func (i BCode) String() string {
	i -= 100
	if i >= BCode(len(_BCode_index)-1) {
		return fmt.Sprintf("BCode(%d)", i+100)
	}
	return _BCode_name[_BCode_index[i]:_BCode_index[i+1]]
}
