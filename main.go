package main

import (
	"api-rest-go/database"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

var databaseConnection *sql.DB

type Product struct {
	ID           int    `json:"id"`
	Product_Code string `json:"product_code"`
	Description  string `json:"description"`
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	databaseConnection = database.InitDB()
	defer databaseConnection.Close()

	r := chi.NewRouter()
	r.Get("/products", AllProductos)
	r.Post("/products", CreateProducto)
	r.Put("/products/{id}", UpdateProducto)
	r.Delete("/products/{id}", DeleteProducto)
	http.ListenAndServe(":3000", r)

}

func DeleteProducto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query, err := databaseConnection.Prepare("delete from products where id=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "successfully deleted"})
}

func UpdateProducto(w http.ResponseWriter, r *http.Request) {
	var product Product
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&product)

	query, err := databaseConnection.Prepare("Update products set product_code=?, description=? where id=?")
	catch(err)
	_, er := query.Exec(product.Product_Code, product.Description, id)
	catch(er)

	defer query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})

}

func CreateProducto(w http.ResponseWriter, r *http.Request) {
	var producto Product
	json.NewDecoder(r.Body).Decode(&producto)

	query, err := databaseConnection.Prepare("Insert products SET product_code=?, description=?")
	catch(err)

	_, er := query.Exec(producto.Product_Code, producto.Description)
	catch(er)
	defer query.Close()

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "successfully created"})
}

func AllProductos(w http.ResponseWriter, r *http.Request) {
	const sql = `SELECT id,product_code,COALESCE(description,'')
				 FROM products`
	results, err := databaseConnection.Query(sql)
	catch(err)
	var products []*Product

	for results.Next() {
		product := &Product{}
		err = results.Scan(&product.ID, &product.Product_Code, &product.Description)

		catch(err)
		products = append(products, product)
	}
	respondwithJSON(w, http.StatusOK, products)
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
