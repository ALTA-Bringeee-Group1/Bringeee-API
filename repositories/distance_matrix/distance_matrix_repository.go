package distanceMatrix

import (
	"bringeee-capstone/configs"
	"bringeee-capstone/entities"
	"bringeee-capstone/entities/web"
	"encoding/json"
	"fmt"
	"net/http"
)

type DistanceMatrixRepository struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewDistanceMatrixRepository() *DistanceMatrixRepository {
	return &DistanceMatrixRepository{
		client:  &http.Client{},
		baseURL: configs.Get().DistanceMatrix.DistanceMatrixBaseURL,
		apiKey:  configs.Get().DistanceMatrix.DistanceMatrixAPIKey,
	}
}

/*
 * Estimate Shortest Distance
 * -------------------------------
 * Mendapatkan estimasi jarak berdasarkan berdasarkan
 * koordinat lintang dan bujur di 2 titik
 *
 * @var originLat 		string 		koordinat lintang titik asal
 * @var originLong 		string		koordinat bujur titik asal
 * @var destinationLat 	string		koordinat lintang titik tujuan
 * @var destinationLat 	string		koordinat bujur titik tujuan
 */
func (repository DistanceMatrixRepository) EstimateShortest(originLat, originLong, destinationLat, destinationLong string) (entities.DistanceMatrix, error) {
	// request
	request, err := http.NewRequest(http.MethodGet, repository.baseURL+"/json", nil)
	if err != nil {
		return entities.DistanceMatrix{}, web.WebError{
			Code:               500,
			ProductionMessage:  "Cannot estimate and calculate distance",
			DevelopmentMessage: "DistanceMatrix Making Request err: " + err.Error(),
		}
	}
	q := request.URL.Query()
	q.Add("origins", fmt.Sprintf("%s,%s", originLat, originLong))
	q.Add("destinations", fmt.Sprintf("%s,%s", destinationLat, destinationLong))
	q.Add("key", repository.apiKey)
	request.URL.RawQuery = q.Encode()
	request.Header.Set("Accept", "application/json")

	// do request
	response, err := repository.client.Do(request)
	if err != nil {
		return entities.DistanceMatrix{}, web.WebError{
			Code:               500,
			ProductionMessage:  "Cannot estimate and calculate distance",
			DevelopmentMessage: "DistanceMatrix Request err: " + err.Error(),
		}
	}
	defer response.Body.Close()

	// parse response
	var data entities.DistanceMatrixResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return entities.DistanceMatrix{}, web.WebError{
			Code:               500,
			ProductionMessage:  "Cannot estimate and calculate distance",
			DevelopmentMessage: "DistanceMatrix parsing response err: " + err.Error(),
		}
	}

	// Translate to universal entity
	return entities.DistanceMatrix{
		Distance:              data.Rows[0].Elements[0].Distance.Text,
		DistanceValue:         data.Rows[0].Elements[0].Distance.Value,
		EstimateDuration:      data.Rows[0].Elements[0].Duration.Text,
		EstimateDurationValue: data.Rows[0].Elements[0].Duration.Value,
		DestinationAddress:    data.DestinationAddresses[0],
		OriginAddress:         data.OriginAddresses[0],
	}, nil
}
