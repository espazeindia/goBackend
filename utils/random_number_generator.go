package utils

import (
	"crypto/rand"
	"espazeBackend/domain/entities"
	"fmt"
	"math/big"
)

func GenerateRandomIndex(arr []*entities.GetProductsForStoreSubcategory) (*entities.GetProductsForStoreSubcategory, error) {
	if len(arr) == 0 {
		return nil, fmt.Errorf("array is empty")
	}

	digit, err := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
	if err != nil {
		return nil, err
	}

	index := int(digit.Int64())
	return arr[index], nil
}
