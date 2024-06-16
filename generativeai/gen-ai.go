// package generativeai

// import (
// 	"log"
// 	"encoding/json"
// 	"net/http"

// 	"github.com/google/generative-ai-go/genai"
// 	"google.golang.org/api/option"
// )

// 	ctx := content.Background()
// 	client,err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyDB5p_WTTnxZYQfFE-u32qzWb9pHV7RJr8"));
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer client.Close()

// 	model := client.GenerativeModel("gemini-pro")

// func GenerateResponse(w http.ResponseWriter, r *http.Request){
// 	var request struct {
// 		Prompt string `json:"prompt"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		http.Error(w, "Invalid Prompt", http.StatusBadRequest)
// 		return
// 	}

// 	prompt := request.Prompt

// 	result,err := model.GenerateContent(ctx, genai.Text(prompt))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(result)
// }

package generativeai

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateResponse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyDB5p_WTTnxZYQfFE-u32qzWb9pHV7RJr8"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	defer client.Close()

	model := client.GenerativeModel("gemini-pro")

	var request struct {
		Prompt string `json:"prompt"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid Prompt", http.StatusBadRequest)
		return
	}

	prompt := request.Prompt

	result, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
