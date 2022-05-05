package entities

type DistanceMatrix struct {
	Distance string
	DistanceValue int
	EstimateDuration string
	EstimateDurationValue int
	DestinationAddress string
	OriginAddress string
}

type DistanceAPIResponse struct {
	Distance string 			`json:"distance"`
	DistanceValue int 			`json:"distance_value"`
	EstimateDuration string 	`json:"estimate_duration"`
	EstimateDurationValue int 	`json:"estimate_duration_value"`
	DestinationAddress string 	`json:"destination_address"`
	OriginAddress string 		`json:"origin_address"`
	EstimatedPrice int64		`json:"estimated_price"`
}


type DistanceMatrixResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses []string `json:"origin_addresses"`
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text string `json:"text"`
				Value int `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text string `json:"text"`
				Value int `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

type EstimateOrderPriceRequest struct {
	DestinationStartLat string `form:"destination_start_lat" validate:"required"`
	DestinationStartLong string `form:"destination_start_long" validate:"required"`
	DestinationEndLat string `form:"destination_end_lat" validate:"required"`
	DestinationEndLong string `form:"destination_end_long" validate:"required"`
	TruckTypeID string `form:"truck_type" validate:"required"`
}