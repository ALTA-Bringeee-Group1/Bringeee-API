package distance_matrix

import (
	"bringeee-capstone/entities"

	"github.com/stretchr/testify/mock"
)

type DistanceMatrixRepositoryMock struct {
	Mock *mock.Mock
}

func NewDistanceMatrixRepositoryMock(mock *mock.Mock) *DistanceMatrixRepositoryMock {
	return &DistanceMatrixRepositoryMock{
		Mock: mock,
	}
}

var DistanceMatrixCollection = []entities.DistanceMatrix {
	{
		Distance: "507 km",
		DistanceValue: 507337,
		EstimateDuration: "5 Hrs, 23 Min",
		EstimateDurationValue: 323,
		DestinationAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
		OriginAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
	},
	{
		Distance: "338 km",
		DistanceValue: 338221,
		EstimateDuration: "5 Hrs, 23 Min",
		EstimateDurationValue: 323,
		DestinationAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
		OriginAddress: "Suryatmajan, Kec. Danurejan, Kota Yogyakarta, Daerah Istimewa Yogyakarta",
	},
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
func (repo DistanceMatrixRepositoryMock) EstimateShortest(originLat, originLong, destinationLat, destinationLong string) (entities.DistanceMatrix, error) {
	param := repo.Mock.Called(originLat, originLong, destinationLat, destinationLong)
	return param.Get(0).(entities.DistanceMatrix), param.Error(1)
}