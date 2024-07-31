package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Item struct {
	ItemID       int     `json:"itemId"`
	ProductName  string  `json:"productName"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
}

type ShippingAddress struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
	Country     string `json:"country"`
}

type ShippingDetails struct {
	ShippingID            int             `json:"shippingId"`
	Carrier               string          `json:"carrier"`
	TrackingNumber        string          `json:"trackingNumber"`
	EstimatedDeliveryDate string          `json:"estimatedDeliveryDate"`
	ShippingAddress       ShippingAddress `json:"shippingAddress"`
}

type Order struct {
	OrderID         int             `json:"orderId"`
	CustomerID      int             `json:"customerId"`
	OrderDate       string          `json:"orderDate"`
	Status          string          `json:"status"`
	Items           []Item          `json:"items"`
	ShippingDetails ShippingDetails `json:"shippingDetails"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}

func loadOrdersFromFile(filename string) ([]Order, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var response OrdersResponse
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Orders, nil
}

func getOrdersByDate(orders []Order, orderDate string, status *string) []Order {
	var filteredOrders []Order

	// แปลงวันที่จากสตริงเป็นรูปแบบวันที่
	dateLayout := "2006-01-02"
	parsedDate, err := time.Parse(dateLayout, orderDate)
	if err != nil {
		fmt.Println("Invalid date format. Use YYYY-MM-DD.")
		return filteredOrders
	}

	for _, order := range orders {
		// แปลงเป็นรูปแบบวันที่ RFC3339
		orderDateTime, err := time.Parse(time.RFC3339, order.OrderDate)
		if err != nil {
			continue // ข้าม order ที่มีวันที่ไม่ถูกต้อง
		}

		// เปรียบเทียบวันที่
		if orderDateTime.Format(dateLayout) == parsedDate.Format(dateLayout) {
			if status == nil || *status == order.Status {
				// เพิ่ม order ที่ status ตรงกัน
				filteredOrders = append(filteredOrders, order)
			}
		}
	}

	return filteredOrders
}

// คำนวณมูลค่ารวม
func calculateTotalValueOrder(order Order) float64 {
	totalValue := 0.0
	for _, item := range order.Items {
		totalValue += float64(item.Quantity) * item.Price
	}
	return totalValue
}

// แสดง order ในรูปแบบตาราง
func displayTableOrders(orders []Order) {
	fmt.Printf("| %-8s | %-11s | %-21s | %-10s | %-11s |\n", "Order ID", "Customer ID", "Order Date", "Status", "Total Value")
	for _, order := range orders {
		totalValue := calculateTotalValueOrder(order)
		fmt.Printf("| %-8d | %-11d | %-21s | %-10s | $%9.2f |\n",
			order.OrderID, order.CustomerID, order.OrderDate, order.Status, totalValue)
	}
}

func main() {
	orders, err := loadOrdersFromFile("data.json")
	if err != nil {
		fmt.Println("Error loading orders:", err)
		return
	}

	// input
	orderDate := "2024-07-29"
	status := "Processing"

	filteredOrders := getOrdersByDate(orders, orderDate, &status)
	// output
	displayTableOrders(filteredOrders)
}
