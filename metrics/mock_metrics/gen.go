package mock_metrics

//go:generate mockgen -package mock_metrics -destination ./mock_scope.go github.com/letsencrypt/boulder/metrics Scope && sed -i 's,github.com/letsencrypt/boulder/vendor/,,' mock_scope.go
