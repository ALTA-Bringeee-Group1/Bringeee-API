package distance_matrix

import "bringeee-capstone/entities"

type DistanceMatrixRepositoryInterface interface {
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
	EstimateShortest(originLat, originLong, destinationLat, destinationLong string) (entities.DistanceMatrix, error)
}