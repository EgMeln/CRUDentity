package postgreSQL

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"strconv"
)

func CreateRecord(num int, inPark bool, remark string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "INSERT INTO parking (num,inparking,remark) VALUES ($1,$2,$3)", num, inPark, remark)
	if err != nil {
		log.Fatalln(err)
	}
}
func ReadAllRecords() string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), "SELECT * FROM parking")
	var str string
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatalln(err)
		}
		str += "Car number: " + strconv.Itoa(int(values[0].(int32))) + ". Is the car in the parking: " + strconv.FormatBool(values[1].(bool)) + ". Remark: " + values[2].(string) + "\n"
	}
	return str
}
func ReadRecordByNum(num int) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var inPark bool
	var remark string
	err = conn.QueryRow(context.Background(), "SELECT inparking, remark from parking where num=$1", num).Scan(&inPark, &remark)
	if err != nil {
		log.Fatalln(err)
	}
	return "Car number: " + strconv.Itoa(num) + ". Is the car in the parking: " + strconv.FormatBool(inPark) + ". Remark: " + remark + "\n"
}
func UpdateRecord(num int, inParking bool, remark string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE parking SET inparking =$1,remark =$2 WHERE num = $3", inParking, remark, num)
	if err != nil {
		log.Fatalln(err)
	}
}
func DeleteRecord(num int) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	row, err := conn.Exec(context.Background(), "DELETE FROM parking where num=$1", num)
	if err != nil {
		log.Fatalln(err)
	}
	if row.RowsAffected() != 1 {
		log.Fatalln("nothing to delete")
	}
}
