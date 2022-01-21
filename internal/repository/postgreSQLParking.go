package repository

import (
	"EgMeln/CRUDentity/internal/model"
	"context"
	"fmt"
	"log"
)

func (rep *Postgres) CreateRecord(lot *model.ParkingLot) {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	_, err := rep.pool.Exec(context.Background(), "INSERT INTO parking (num,inparking,remark) VALUES ($1,$2,$3)", lot.Num, lot.InParking, lot.Remark)
	if err != nil {
		log.Fatalln(err)
	}
}
func (rep *Postgres) ReadAllRecords() []*model.ParkingLot {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	rows, _ := rep.pool.Query(context.Background(), "SELECT * FROM parking")
	var lots []*model.ParkingLot
	for rows.Next() {
		var lot model.ParkingLot
		values, err := rows.Values()
		if err != nil {
			log.Fatalln(err)
		}
		lot.Num = int(values[0].(int32))
		lot.InParking = values[1].(bool)
		lot.Remark = values[2].(string)
		lots = append(lots, &lot)
		fmt.Println(lots)
		//str = "Car number: " + strconv.Itoa(int(values[0].(int32))) + ". Is the car in the parking: " + strconv.FormatBool(values[1].(bool)) + ". Remark: " + values[2].(string) + "\n"
	}
	return lots
}
func (rep *Postgres) ReadRecordByNum(num int) *model.ParkingLot {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	var lot model.ParkingLot
	err := rep.pool.QueryRow(context.Background(), "SELECT num,inparking, remark from parking where num=$1", num).Scan(&lot.Num, &lot.InParking, &lot.Remark)
	if err != nil {
		log.Fatalln(err)
	}
	return &lot
}
func (rep *Postgres) UpdateRecord(num int, inParking bool, remark string) {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	_, err := rep.pool.Exec(context.Background(), "UPDATE parking SET inparking =$1,remark =$2 WHERE num = $3", inParking, remark, num)
	if err != nil {
		log.Fatalln(err)
	}
}
func (rep *Postgres) DeleteRecord(num int) {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://egormelnikov:54236305@localhost:5432/egormelnikov"))
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close(context.Background())
	row, err := rep.pool.Exec(context.Background(), "DELETE FROM parking where num=$1", num)
	if err != nil {
		log.Fatalln(err)
	}
	if row.RowsAffected() != 1 {
		log.Fatalln("nothing to delete")
	}
}
