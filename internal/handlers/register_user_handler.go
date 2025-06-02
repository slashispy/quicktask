package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slashispy/quicktask/internal/models"
	"github.com/slashispy/quicktask/internal/usecases"
)

func RegisterUserHandler(uc *usecases.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.UserDTO

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := uc.RegisterUser(input); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

}
