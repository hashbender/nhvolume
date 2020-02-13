package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"

	"github.com/jasonlvhit/gocron"
)

type AlgoResultResp struct {
	Algos []Algo `json:"miningAlgorithms"`
}

type Algo struct {
	Index string `json:"algorithm"`
	Name  string `json:"title"`
}

type OrderGetResp struct {
	Orders []OrderRep `json:"list"`
}

type ResultResp struct {
	Orders []OrderRep `json:"list"`
}

type OrderRep struct {
	Alive     bool    `json:"alive"`
	Price     float64 `json:"price,string"`
	Speed     float64 `json:"acceptedCurrentSpeed,string"`
	Market    string  `json:"market"`
	Algorithm Algo    `json:"algorithm"`
}

type Stat struct {
	Time      time.Time
	Location  string
	Algorithm string
	Volume    float64
	Orders    int
}

var (
	db *pg.DB
)

type dbLogger struct {
	prefix string
}

func (d dbLogger) BeforeQuery(c context.Context, evt *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, event *pg.QueryEvent) error {
	query, err := event.FormattedQuery()
	log.Printf("%s: ##### (%s) %s", d.prefix, time.Since(event.StartTime), query)
	return err
}

func main() {
	// Connect to Accounts database
	db = initDatabase("nhvolume",
		"accountmanager",
		"password12!")
	defer db.Close()

	db.AddQueryHook(dbLogger{
		prefix: "MASTER",
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Run DB migrations
	createSchema(db)

	checkStats()

	gocron.Every(10).Minutes().Do(checkStats)
	<-gocron.Start()
}

type algoStat struct {
	volume float64
	orders int
}

func checkStats() {
	log.Printf("Checking stats")
	resp, err := http.Get("https://api2.nicehash.com/main/api/v2/public/orders/active/")
	if err != nil {
		log.Printf("Error getting orders.  Error: %s", err)
		return
	}
	var orders OrderGetResp
	err = json.NewDecoder(resp.Body).Decode(&orders)
	if err != nil {
		log.Printf("Error parsing: %s", err)
		return
	}

	totalVolume := 0.0
	totalOrders := 0
	stats := make(map[string]map[string]algoStat)
	for _, order := range orders.Orders {
		if stats[order.Market] == nil {
			stats[order.Market] = make(map[string]algoStat)
		}
		stat := stats[order.Market][order.Algorithm.Name]
		stat.volume += order.Price * order.Speed
		stat.orders++
		stats[order.Market][order.Algorithm.Name] = stat
		totalVolume += order.Price * order.Speed
		totalOrders++
	}

	for region, algoMap := range stats {
		for algo, stat := range algoMap {
			db.Insert(&Stat{
				Time:      time.Now(),
				Location:  region,
				Algorithm: algo,
				Volume:    stat.volume,
				Orders:    stat.orders,
			})
		}
	}

	// for _, algo := range algoResult.Algos {
	// 	resp, err := http.Get(fmt.Sprintf("https://api.nicehash.com/api?method=orders.get&location=0&algo=%d", algo.Index))
	// 	if err != nil {
	// 		glog.Infof("Error getting %s.  Error: %s", algo.Name, err)
	// 		continue
	// 	}
	// 	var resultRep OrderGetResp
	// 	err = json.NewDecoder(resp.Body).Decode(&resultRep)
	// 	if err != nil {
	// 		glog.Infof("Error parsing: %s", err)
	// 		continue
	// 	}
	// 	sum := 0.0
	// 	for _, order := range resultRep.Result.Orders {
	// 		if order.Alive {
	// 			sum += order.Speed * algo.Norm * order.Price
	// 		}
	// 	}
	// 	totalVolume += sum
	// 	totalOrders += len(resultRep.Result.Orders)
	// 	db.Insert(&Stat{
	// 		Time:      time.Now(),
	// 		Location:  "EU",
	// 		Algorithm: algo.Name,
	// 		Volume:    sum,
	// 		Orders:    len(resultRep.Result.Orders),
	// 	})
	// 	glog.Infof("EU Algo: %s. Total BTC: %f.  Num Orders: %d", algo.Name, sum, len(resultRep.Result.Orders))
	// }

	// for _, algo := range algoResult.Algos {
	// 	resp, err := http.Get(fmt.Sprintf("https://api.nicehash.com/api?method=orders.get&location=1&algo=%d", algo.Index))
	// 	if err != nil {
	// 		glog.Infof("Error getting %s.  Error: %s", algo.Name, err)
	// 		continue
	// 	}
	// 	var resultRep OrderGetResp
	// 	err = json.NewDecoder(resp.Body).Decode(&resultRep)
	// 	if err != nil {
	// 		glog.Infof("Error parsing: %s", err)
	// 		continue
	// 	}
	// 	sum := 0.0
	// 	for _, order := range resultRep.Result.Orders {
	// 		if order.Alive {
	// 			sum += order.Speed * algo.Norm * order.Price
	// 		}
	// 	}
	// 	totalVolume += sum
	// 	totalOrders += len(resultRep.Result.Orders)
	// 	db.Insert(&Stat{
	// 		Time:      time.Now(),
	// 		Location:  "US",
	// 		Algorithm: algo.Name,
	// 		Volume:    sum,
	// 		Orders:    len(resultRep.Result.Orders),
	// 	})
	// 	glog.Infof("US Algo: %s. Total BTC: %f.  Num Orders: %d", algo.Name, sum, len(resultRep.Result.Orders))
	// }

	db.Insert(&Stat{
		Time:      time.Now(),
		Location:  "BOTH",
		Algorithm: "total",
		Volume:    totalVolume,
		Orders:    totalOrders,
	})
	log.Printf("Total BTC: %f", totalVolume)
}

func initDatabase(databaseName string, databaseUsername string, databasePassword string) (db *pg.DB) {
	log.Printf("Opening db: " + databaseName)

	var num int
	var err error

	for i := 0; i <= 10; i++ {
		db = pg.Connect(&pg.Options{
			User:     databaseUsername,
			Password: databasePassword,
			Database: databaseName,
		})

		_, err = db.Query(pg.Scan(&num), "SELECT ?", 1)
		if err != nil {
			log.Printf("Error opening database, retrying: %s", err)
		} else {
			break
		}
		time.Sleep(time.Second * 2)
	}

	if db == nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to database: %s", db)
	return db
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*Stat)(nil),
	} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			log.Printf("Could not create tables from schema. Error: %s", err.Error())
		}
	}
	return nil
}
