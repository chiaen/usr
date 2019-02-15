package main

const (
	TableUserRecords  = "user_records"
	TableRefreshToken = "refresh_tokens"
)

const (
	ColUID        = "uid"
	ColEmail      = "email"
	ColPassword   = "password"
	ColToken      = "token"
	ColIssuedAt   = "issued_at"
	ColValidUntil = "valid_until"
)

var (
	TableUserRecordsColums = []string{
		ColUID,
		ColEmail,
		ColPassword,
	}

	TableRefreshTokens = []string{
		ColUID,
		ColToken,
		ColIssuedAt,
		ColValidUntil,
	}
)
