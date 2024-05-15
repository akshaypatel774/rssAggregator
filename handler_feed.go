package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akshaypatel774/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeed (w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:		   uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:	   params.Name,
		Url:	   params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt create feed: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseFeedToFeed(feed))
}


func (apiCfg *apiConfig) handlerGetFeeds (w http.ResponseWriter, r *http.Request) {	
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldnt get any feeds: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseFeedsToFeeds(feeds))
}