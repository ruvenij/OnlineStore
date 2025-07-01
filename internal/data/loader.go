package data

import (
	"OnlieStore/internal/model"
	"bufio"
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) LoadUsers(filePath string) ([]*model.User, error) {
	result := make([]*model.User, 0)
	f, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read() // header row
	if err != nil {
		return result, err
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			logrus.WithError(err).Error("Read line failed")
			continue
		}

		u := parseUserRow(line)
		result = append(result, u)
	}

	return result, nil
}

func (l *Loader) LoadProducts(filePath string) ([]*model.ProductDetails, error) {
	result := make([]*model.ProductDetails, 0)
	f, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read() // header row
	if err != nil {
		return result, err
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			logrus.WithError(err).Error("Read line failed")
			continue
		}

		p, err := parseProductRow(line)
		if err != nil {
			logrus.WithError(err).Error("Parse product row failed")
			continue
		}

		result = append(result, p)
	}

	return result, nil
}

func parseUserRow(row []string) *model.User {
	return &model.User{
		ID:       row[0],
		Name:     row[1],
		Password: row[2],
	}
}

func parseProductRow(row []string) (*model.ProductDetails, error) {
	price, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return nil, err
	}

	qty, err := strconv.Atoi(row[3])
	if err != nil {
		return nil, err
	}

	return &model.ProductDetails{
		Name:          row[0],
		Price:         price,
		Category:      row[2],
		AddedQuantity: qty,
	}, nil
}
