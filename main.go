package main

//local blockchain with peers and a server for genesis block that can be used on cloud to connect distant peers
import (
	"log"
	"net/http"
	"os"
	"project/connection"
	controller "project/controller"
	model "project/model"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//creating the webserver using gorrila/mux

func run() error {
	mux := controller.MakeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {

		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var genesisBlock model.Blockd

		t := time.Now()

		genesisBlock.Index = 0
		genesisBlock.Timestamp = t.String()
		genesisBlock.BPM = 12
		genesisBlock.PrevHash = " "
		genesisBlock.Hash = model.CalculateHash(genesisBlock)
		genesisBlock = model.Blockd{0, t.String(), genesisBlock.BPM, genesisBlock.Hash, ""}
		spew.Dump(genesisBlock)
		model.Blockchain = append(model.Blockchain, genesisBlock)

		stmt, err := connection.Db.Prepare("INSERT INTO Block(data,hash,phash) VALUES(?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(genesisBlock.BPM, genesisBlock.Hash, genesisBlock.PrevHash)
		if err != nil {
			log.Fatal(err)
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		rowCnt, err := res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	}()
	log.Fatal(run())

}
