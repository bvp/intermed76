package intermed76

import (
	"time"
	// "time"
	// "encoding/json"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/kr/pretty"
)

var (
	client *Client
)

func ExampleMain() {
	log.Println("Login to https://intermed76.ru")
	client := NewClient("Имя", "Фамилия", "Отчество", "2000-01-01", "7690299770000000", nil)
	_, err := client.Login()
	if err != nil {
		log.Fatal(err.Error())
	}

	client.GetMos()

	client.GetSession()

	client.GetRecords()

	client.FindSpecs("1.2.643.5.1.13.3.25.76.3", "10304")

	client.FindResources("22")

	client.GetScheduleTable()

	client.GetRecords()
}

func TestClient_Login(t *testing.T) {
	time.Sleep(2 * time.Second)
	t.Log("Testing Login")
	_, err := client.Login()

	if err != nil {
		t.Error(err.Error())
	}
}

func TestClient_GetSession(t *testing.T) {
	time.Sleep(2 * time.Second)
	gsResponse, _ := client.GetSession()
	t.Logf(":: gsResponse body to struct - %# v", pretty.Formatter(&gsResponse))

	t.Logf(":: erz: %s", gsResponse.ErzCode)
	t.Logf(":: rri: %s", gsResponse.Rri)
	t.Logf(":: rri oid: %s", gsResponse.RriOid)
}

func TestClient_FindSpecs(t *testing.T) {
	time.Sleep(2 * time.Second)
	fsResponse, _ := client.FindSpecs("1.2.643.5.1.13.3.25.76.3", "10304")

	specs := fsResponse.GetServiceSpecsInfoResponse.ListServiceSpecs.ServiceSpec
	for _, spec := range specs {
		t.Logf("  :: %d - %s", spec.ServiceSpecID, spec.ServiceSpecName)
	}
}

func TestClient_FindResources(t *testing.T) {
	time.Sleep(2 * time.Second)
	client.GetSession()
	grResponse, _ := client.FindResources("22")

	resources := grResponse.GetResourceInfoResponse.ListResource.Resource

	for _, r := range resources {
		var n string
		if i, ok := r.ResourceID.(string); ok {
			n = string(i)
		} else if s, ok := r.ResourceID.(float64); ok {
			n = fmt.Sprintf("%.0f", s)
		}

		r.ResourceID = n
		t.Logf("  :: %s - %s", r.ResourceID, r.ResourceName)
	}
}

func TestClient_GetSchedules(t *testing.T) {
	time.Sleep(2 * time.Second)
	schedules, _ := client.GetSchedule()
	t.Logf(":: GetSchedule - %# v", pretty.Formatter(&schedules))
}

func TestClient_GetRecords(t *testing.T) {
	time.Sleep(2 * time.Second)
	records, _ := client.GetRecords()

	for _, r := range records {
		t.Logf("  :: %s (%s : %s) - %s %s", r.DoctorName, r.ServiceID, r.ServiceName, r.VisitDate, r.VisitTime)
	}
}

func setup() {
	log.Println("setup")

	csvFile, err := os.Open("./clients.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	row, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	log.Printf("row - %s", row)

	client = NewClient(row[0], row[1], row[2], row[3], row[4], nil)
}

func shutdown() {
	log.Println("shutdown")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
