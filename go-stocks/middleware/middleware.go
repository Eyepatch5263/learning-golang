package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/eyepatch5263/go-postgress/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	
)

type response struct {
	ID 	int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`

}

func CreateDbConnection() *sql.DB{
	err:=godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	db,err:=sql.Open("postgres",os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Println("error connecting to the database:", err)
		panic("Error connecting to the database")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("error pinging the database:", err)
		panic("Error pinging the database")
	}
	fmt.Println("Database connection established successfully")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request){
	var stock models.Stock
	//decoding the request body into the stock struct
	err:=json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	insertId := InsertStock(stock)

	res:=response{
		ID: insertId,
		Message: "Stock created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetStockById(w http.ResponseWriter, r *http.Request){
	params:=mux.Vars(r)
	id,err:=strconv.Atoi(params["stockId"])
	if err != nil {
		fmt.Println("Error converting stockId to int:", err)
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}
	stock,err:=getStock(int64(id))
	if err != nil {
		fmt.Println("Error getting stock by ID:", err)
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stock)


}

func GetAllStock(w http.ResponseWriter, r *http.Request){
	stocks,err:=getAllStocks()

	if err != nil {
		fmt.Println("Error getting all stocks:", err)
		http.Error(w, "Error retrieving stocks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStockById(w http.ResponseWriter, r *http.Request){
    params := mux.Vars(r)
    id, err := strconv.Atoi(params["stockId"])
    if err != nil {
        fmt.Println("Error converting stockId to int:", err)
        http.Error(w, "Invalid stock ID", http.StatusBadRequest)
        return
    }
    var stock models.Stock
    err = json.NewDecoder(r.Body).Decode(&stock)
    if err != nil {
        fmt.Println("Error decoding stock:", err)
        http.Error(w, "Invalid stock data", http.StatusBadRequest)
        return
    }
    updatedStock, err := updateStock(int64(id), stock)
    if err != nil {
        fmt.Println("Error updating stock:", err)
        http.Error(w, "Error updating stock", http.StatusInternalServerError)
        return
    }
    msg:=fmt.Sprintf("Stock with ID %v updated successfully", updatedStock)
    res:=response{
        ID:      int64(id),
        Message: msg,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func DeleteStockById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["stockId"])
	if err != nil {
		fmt.Println("Error converting stockId to int:", err)
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}
	deletedStock, err := deleteStock(int64(id))
	if err != nil {
		fmt.Println("Error deleting stock:", err)
		http.Error(w, "Error deleting stock", http.StatusInternalServerError)
		return
	}
	msg:= fmt.Sprintf("Stock with ID %v deleted successfully", deletedStock)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func InsertStock(stock models.Stock) int64 {
	db:= CreateDbConnection()
	defer db.Close()

	sqlStatement:=`INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64
	err:=db.QueryRow(sqlStatement, stock.StockName, stock.StockPrice, stock.StockCompany).Scan(&id)
	if err != nil {
		fmt.Println("Error inserting stock:", err)
		panic("Error inserting stock")
	}
	fmt.Println("Stock inserted successfully with ID:", id)
	return id
}


func getStock(id int64) (models.Stock,error){
	db := CreateDbConnection()
	defer db.Close()
	
	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`
	err := db.QueryRow(sqlStatement, id).Scan(&stock.StockID, &stock.StockName, &stock.StockPrice, &stock.StockCompany)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No stock found with ID:", id)
			return models.Stock{}, nil
		}
		fmt.Println("Error getting stock by ID:", err)
		return models.Stock{}, nil
	}
	return stock, err
}

func getAllStocks() ([]models.Stock,error) {
	db := CreateDbConnection()
	defer db.Close()
	
	var stock []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Error getting all stocks:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Stock
		err := rows.Scan(&s.StockID, &s.StockName, &s.StockPrice, &s.StockCompany)
		if err != nil {
			fmt.Println("Error scanning stock row:", err)
			return nil, err
		}
		stock = append(stock, s)
	}
	return stock, err
}

func updateStock(id int64, stock models.Stock) (models.Stock,error) {
	db := CreateDbConnection()
	defer db.Close()
	
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	_, err := db.Exec(sqlStatement,id, stock.StockName, stock.StockPrice, stock.StockCompany)
	if err != nil {
		fmt.Println("Error updating stock:", err)
		return models.Stock{}, err
	}
	fmt.Println("Stock with ID", id, "updated successfully")
	return stock, err
}

func deleteStock(id int64) (models.Stock, error) {
	db:= CreateDbConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		fmt.Println("Error deleting stock:", err)
		return models.Stock{}, err
	}
	fmt.Println("Stock with ID", id, "deleted successfully")
	return models.Stock{StockID: id}, nil
}