// Code generated by "stringer -type=FeatureFlag"; DO NOT EDIT.

package features

import "fmt"

const _FeatureFlag_name = "unusedAllowAccountDeactivationAllowKeyRolloverResubmitMissingSCTsOnlyGoogleSafeBrowsingV4UseAIAIssuerURLAllowTLS02ChallengesGenerateOCSPEarlyCountCertificatesExactRandomDirectoryEntryIPv6FirstDirectoryMeta"

var _FeatureFlag_index = [...]uint8{0, 6, 30, 46, 69, 89, 104, 124, 141, 163, 183, 192, 205}

func (i FeatureFlag) String() string {
	if i < 0 || i >= FeatureFlag(len(_FeatureFlag_index)-1) {
		return fmt.Sprintf("FeatureFlag(%d)", i)
	}
	return _FeatureFlag_name[_FeatureFlag_index[i]:_FeatureFlag_index[i+1]]
}
