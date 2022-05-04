package entities

type DistanceMatrix struct {
	Distance string
	DistanceValue int
	EstimateDuration string
	EstimateDurationValue int
	DestinationAddress string
	OriginAddress string
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