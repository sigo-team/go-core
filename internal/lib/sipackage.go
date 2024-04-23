package lib

import (
	"archive/zip"
	"gopkg.in/loremipsum.v1"
	"io"
	"math/rand"
	"os"
)

type (
	Slide struct {
		Content     *string
		ContentType *string
		Comment     *string
	}
	Question struct {
		Type *string `json:"Type,omitempty"`

		PriceMin  *int `json:"PriceMin,omitempty"`
		PriceMax  *int `json:"PriceMax,omitempty"`
		PriceStep *int `json:"PriceStep,omitempty"`

		QuestionSlides []*Slide `json:"QuestionSlides,omitempty"`
		AnswerSlides   []*Slide `json:"AnswerSlides,omitempty"`
	}
	Theme struct {
		Name      *string
		Questions []*Question
	}
	Round struct {
		Name   *string
		Themes []*Theme
	}
	Package struct {
		Name   string
		Rounds []*Round
	}
)

func Unzip(path string) error {
	err := os.MkdirAll("./unzipSiPackage/media/", os.ModePerm)
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)

	for _, file := range reader.File {
		if !file.FileInfo().IsDir() {
			openedFile, err := os.OpenFile("./unzipSiPackage/"+file.Name, os.O_CREATE, os.ModePerm)
			if err != nil {
				return err
			}

			readCloser, err := file.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(openedFile, readCloser)
			if err != nil && err != io.EOF {
				return err
			}

			err = openedFile.Close()
			if err != nil {
				return err
			}

			err = readCloser.Close()
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func RemovePackage() error {
	if err := os.RemoveAll("./unzipSiPackage"); err != nil {
		return err
	}
	return nil
}

func GenerateRandomPackage() Package {
	var pck Package
	ipsum := loremipsum.New()

	pck.Name = ipsum.Word()
	roundCount := rand.Intn(5) + 1
	for i := 0; i < roundCount; i++ {
		var roundTemp Round
		*roundTemp.Name = ipsum.Word()
		themeCount := rand.Intn(5) + 1
		for j := 0; j < themeCount; j++ {
			var themeTemp Theme
			*themeTemp.Name = ipsum.Words(2)
			questionCount := rand.Intn(5) + 1
			for k := 0; k < questionCount; k++ {
				var questionTemp Question
				*questionTemp.Type = "secret"
				*questionTemp.PriceMin = rand.Intn(200)
				*questionTemp.PriceMax = rand.Intn(800) + 200
				*questionTemp.PriceStep = rand.Intn(200)
				QuestionSlidesCount := rand.Intn(3) + 1
				AnswerSlidesCount := rand.Intn(3) + 1
				for l := 0; l < QuestionSlidesCount; l++ {
					var tempSlide Slide
					*tempSlide.ContentType = "text"
					tmp := ipsum.SentenceList(2)
					*tempSlide.Comment, *tempSlide.Content = tmp[0], tmp[1]
					questionTemp.QuestionSlides = append(questionTemp.QuestionSlides, &tempSlide)
				}
				for l := 0; l < AnswerSlidesCount; l++ {
					var tempSlide Slide
					*tempSlide.ContentType = "text"
					tmp := ipsum.SentenceList(2)
					*tempSlide.Comment, *tempSlide.Content = tmp[0], tmp[1]
					questionTemp.AnswerSlides = append(questionTemp.AnswerSlides, &tempSlide)
				}
				themeTemp.Questions = append(themeTemp.Questions, &questionTemp)
			}
			roundTemp.Themes = append(roundTemp.Themes, &themeTemp)
		}
		pck.Rounds = append(pck.Rounds, &roundTemp)
	}

	return pck
}
