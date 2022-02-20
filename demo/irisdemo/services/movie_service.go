package services

import (
	"fmt"
	"github.com/l1nkkk/shopSystem/demo/irisdemo/repositories"
)

type MovieService interface {
	ShowMovieName() string
}

type MovieServiceManger struct{
	repo repositories.MovieRepository
}

func NewMovieServiceManger(repo repositories.MovieRepository) MovieService{
	return &MovieServiceManger{repo:repo}
}

func (m *MovieServiceManger) ShowMovieName() string{
	fmt.Println("movie name: " + m.repo.GetMovieName())
	return "movie name: " + m.repo.GetMovieName()
}

