package handler

import (
	"fmt"
	"net/http"
)

type Watchlist struct{}

func (wl *Watchlist) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create route")
}

func (wl *Watchlist) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getbyid route")
}

func (wl *Watchlist) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List route")
}

func (wl *Watchlist) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updatebyid route")
}

func (wl *Watchlist) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deletebyid route")
}
