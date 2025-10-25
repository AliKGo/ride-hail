package types

var (
	EntityRolePassenger = "passenger"
	EntityRoleDriver    = "driver"
)

var (
	RideStatusREQUESTED   = "REQUESTED"
	RideStatusMATCHED     = "MATCHED"
	RideStatusEN_ROUTE    = "EN_ROUTE"
	RideStatusARRIVED     = "ARRIVED"
	RideStatusIN_PROGRESS = "IN_PROGRESS"
	RideStatusCOMPLETED   = "COMPLETED"
	RideStatusCANCELLED   = "CANCELLED"
)

var (
	DriverStatusOffline   = "OFFLINE"
	DriverStatusAvailable = "AVAILABLE"
	DriverStatusBusy      = "BUSY"
	DriverStatusEnRoute   = "EN_ROUTE"
)
