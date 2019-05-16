package model

// Privilege is a bitmask containing one or more privileges assigned to a user
// or a session or required by an operation.
type Privilege uint8

const (
	// PrivSetup is the privilege needed to create, modify, or delete
	// products, SKUs, and events.
	PrivSetup Privilege = 1 << iota

	// PrivAnalyze allows read-only access to all data in the system.
	PrivAnalyze

	// PrivHandleOrders allows recording offline orders and making notes on
	// orders.
	PrivHandleOrders

	// PrivSell allows making and recording in-person sales.  Note that this
	// includes refunding those sales within a few minutes of recording
	// them.
	PrivSell

	// PrivAdmit allows recording ticket usage.
	PrivAdmit
)
