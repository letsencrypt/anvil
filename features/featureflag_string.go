// Code generated by "stringer -type=FeatureFlag"; DO NOT EDIT.

package features

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[unused-0]
	_ = x[WriteIssuedNamesPrecert-1]
	_ = x[HeadNonceStatusOK-2]
	_ = x[RemoveWFE2AccountID-3]
	_ = x[CheckRenewalFirst-4]
	_ = x[ParallelCheckFailedValidation-5]
	_ = x[DeleteUnusedChallenges-6]
	_ = x[BlockedKeyTable-7]
	_ = x[StoreKeyHashes-8]
	_ = x[CAAValidationMethods-9]
	_ = x[CAAAccountURI-10]
	_ = x[EnforceMultiVA-11]
	_ = x[MultiVAFullResults-12]
	_ = x[MandatoryPOSTAsGET-13]
	_ = x[AllowV1Registration-14]
	_ = x[V1DisableNewValidations-15]
	_ = x[PrecertificateRevocation-16]
	_ = x[StripDefaultSchemePort-17]
	_ = x[StoreIssuerInfo-18]
	_ = x[StoreRevokerInfo-19]
	_ = x[RestrictRSAKeySizes-20]
	_ = x[FasterNewOrdersRateLimit-21]
	_ = x[NonCFSSLSigner-22]
	_ = x[DNSServerHealthChecks-23]
}

const _FeatureFlag_name = "unusedWriteIssuedNamesPrecertHeadNonceStatusOKRemoveWFE2AccountIDCheckRenewalFirstParallelCheckFailedValidationDeleteUnusedChallengesBlockedKeyTableStoreKeyHashesCAAValidationMethodsCAAAccountURIEnforceMultiVAMultiVAFullResultsMandatoryPOSTAsGETAllowV1RegistrationV1DisableNewValidationsPrecertificateRevocationStripDefaultSchemePortStoreIssuerInfoStoreRevokerInfoRestrictRSAKeySizesFasterNewOrdersRateLimitNonCFSSLSignerDNSServerHealthChecks"

var _FeatureFlag_index = [...]uint16{0, 6, 29, 46, 65, 82, 111, 133, 148, 162, 182, 195, 209, 227, 245, 264, 287, 311, 333, 348, 364, 383, 407, 421, 442}

func (i FeatureFlag) String() string {
	if i < 0 || i >= FeatureFlag(len(_FeatureFlag_index)-1) {
		return "FeatureFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FeatureFlag_name[_FeatureFlag_index[i]:_FeatureFlag_index[i+1]]
}
