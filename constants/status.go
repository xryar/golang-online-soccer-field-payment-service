package constants

type PaymentStatus int
type PaymentStatusString string

const (
	Initial    PaymentStatus = 0
	Pending    PaymentStatus = 100
	Settlement PaymentStatus = 200
	Expire     PaymentStatus = 300

	InitialString    PaymentStatusString = "Initial"
	PendingString    PaymentStatusString = "Pending"
	SettlementString PaymentStatusString = "Settlement"
	ExpireString     PaymentStatusString = "Expire"
)

var mapStringToInt = map[PaymentStatusString]PaymentStatus{
	InitialString:    Initial,
	PendingString:    Pending,
	SettlementString: Settlement,
	ExpireString:     Expire,
}

var mapIntToString = map[PaymentStatus]PaymentStatusString{
	Initial:    InitialString,
	Pending:    PendingString,
	Settlement: SettlementString,
	Expire:     ExpireString,
}

func (p PaymentStatus) GetStatusString() PaymentStatusString {
	return mapIntToString[p]
}

func (ps PaymentStatusString) GetStatusInt() PaymentStatus {
	return mapStringToInt[ps]
}
