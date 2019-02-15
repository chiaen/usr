package main

const (
	TableUserRecords     = "user_records"
	TableUserHasInterest = "user_has_interests"
	TableInterests       = "interests"
)

const (
	ColID         = "id"
	ColName       = "name"
	ColUID        = "uid"
	ColEmail      = "email"
	ColPassword   = "password"
	ColUserID     = "user_id"
	ColInterestID = "interest_id"
	ColPublished  = "published"
)

var (
	TableUserRecordsColums = []string{
		ColUID,
		ColEmail,
		ColPassword,
	}

	TableUserHasInterestColumns = []string{
		ColUserID,
		ColInterestID,
		ColPublished,
	}
	TableInterestsColumns = []string{
		ColID,
		ColName,
	}
)
