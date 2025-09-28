package handlers

import "fmt"
import "net/http"
import "github.com/rakshitg600/notakto-solo/fxns"
import "github.com/rakshitg600/notakto-solo/types"
import "encoding/json"


// HelloHandler handles the "/" route
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

//CreateHandler handles the /create route
func CreateHandler (w http.ResponseWriter, r *http.Request){

	//allow POST req only
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	//this is how we decode the json req body to golang variables
	var req types.CreateGameRequest
	decoder:=json.NewDecoder(r.Body)
	err:=decoder.Decode(&req)
	if err!=nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return 
	}
	//if return early since, go might keep reading the req, to prevent it we use defer, anything after defer acts like a cleanup function of useEffect in react
	defer r.Body.Close()

	//actual fxn call
	res := fxns.CreateGame(
		types.BoardNumber(req.NumberOfBoards),
		types.BoardSize(req.BoardSize),
		types.DifficultyLevel(req.Difficulty),
	)

	//respond
	w.Header().Set("Content-Type", "application/json") //in header, telling that data in body is in json format, not html or protobuf or xml, etc
	//setting the body of response
	resp := types.GameResponse{
    	Success:   true,
    	SessionId: res.SessionId,
    	GameState: res.GameState,
	}
	//encoding the response in json format
	json.NewEncoder(w).Encode(resp)
}