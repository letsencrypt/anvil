// Code generated by "stringer -type=FeatureFlag"; DO NOT EDIT.

package features

import "strconv"

const _FeatureFlag_name = "unusedUseAIAIssuerURLReusePendingAuthzCountCertificatesExactIPv6FirstAllowRenewalFirstRLWildcardDomainsForceConsistentStatusEnforceChallengeDisableRPCHeadroomTLSSNIRevalidationEmbedSCTsCancelCTSubmissionsVAChecksGSBEnforceV2ContentTypeEnforceOverlappingWildcardsOrderReadyStatusCAAValidationMethodsAllowTLSALPN01Challenge"

var _FeatureFlag_index = [...]uint16{0, 6, 21, 38, 60, 69, 88, 103, 124, 147, 158, 176, 185, 204, 215, 235, 262, 278, 298, 321}

func (i FeatureFlag) String() string {
	if i < 0 || i >= FeatureFlag(len(_FeatureFlag_index)-1) {
		return "FeatureFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FeatureFlag_name[_FeatureFlag_index[i]:_FeatureFlag_index[i+1]]
}
