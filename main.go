// intermed76 project main.go
package intermed76

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	urlBase               = "https://intermed76.ru"
	urlFindPatient        = "https://intermed76.ru/intermed/findPatient"
	urlGetMos             = "https://intermed76.ru/intermed/getAvaliableMos"
	urlGetAuthEsia        = "https://intermed76.ru/intermed/getAuthEsia"
	urlCancelAppointment  = "https://intermed76.ru/intermed/cancelAppointment"
	urlGetSession         = "https://intermed76.ru/intermed/getSession"
	urlFindSpecs          = "https://intermed76.ru/intermed/findSpecs"
	urlFindResources      = "https://intermed76.ru/intermed/findResources"
	urlFindSchedules      = "https://intermed76.ru/intermed/findSchedules"
	urlFindSchedulesTable = "https://intermed76.ru/intermed/findSchedulesTable"
	urlGetRecordsWithErz  = "https://intermed76.ru/intermed/getRecordsWithErz"
)

type (
	Client struct {
		FirstName    string             `json:"firstName"`
		LastName     string             `json:"lastName"`
		MiddleName   string             `json:"middleName"`
		BirthDate    string             `json:"birthDate"`
		RecordSource string             `json:"recordSource"`
		OmsNumber    string             `json:"omsNumber"`
		PatientID    string             `json:"patientID"`
		SessionID    string             `json:"sessionId"`
		http         *http.Client       `xorm:"-"`
		savedCookie  *http.Cookie       `xorm:"-"`
		session      GetSessionResponse `xorm:"-"`
	}

	PatientData struct {
		OmsNumber    string `json:"OMS_Number"`
		FirstName    string `json:"First_Name"`
		LastName     string `json:"Last_Name"`
		MiddleName   string `json:"Middle_Name"`
		BirthDate    string `json:"Birth_Date"`
		RecordSource string `json:"Record_Source"` // "intermed"
	}

	GetPatientInfoRequest struct {
		PatientData `json:"Patient_Data"`
	}

	LoginRequest struct {
		GetPatientInfoRequest `json:"GetPatientInfoRequest"`
	}

	LoginResponse struct {
		GetPatientInfoResponse struct {
			Error struct {
				ErrorDetail struct {
					ErrorCode    int    `json:"errorCode"`
					ErrorMessage string `json:"errorMessage"`
				} `json:"errorDetail"`
			} `json:"Error"`
			PatientID string `json:"Patient_Id"`
			SessionID string `json:"Session_ID"`
		} `json:"GetPatientInfoResponse"`
	}

	MosResponse []struct {
		Address      string `json:"address"`
		AddressTitle string `json:"addressTitle"`
		AddressWsdl  string `json:"address_wsdl"`
		Always       string `json:"always"`
		CodeTfoms    string `json:"codeTfoms"`
		CodeTfoms2   string `json:"code_tfoms"`
		Email        string `json:"email"`
		IDMo         string `json:"id_mo"`
		Irrelevant   string `json:"irrelevant"`
		Name         string `json:"name"`
		Never        string `json:"never"`
		Oid          string `json:"oid"`
		OldOid       string `json:"oldOid"`
		Phone        string `json:"phone"`
		Refferal     string `json:"refferal"`
		Site         string `json:"site"`
		Territory    string `json:"territory"`
	}

	GetSessionResponse struct {
		AttachmentMoName  string `json:"attachmentMoName"`
		DepartOid         string `json:"departOid"`
		Doctor            string `json:"doctor"`
		Email             string `json:"email"`
		EndDateRange      string `json:"endDateRange"`
		EndTimeRange      string `json:"endTimeRange"`
		ErzCode           string `json:"erzCode"`
		IDPost            string `json:"idPost"`
		Inserted          int    `json:"inserted"`
		IP                string `json:"ip"`
		MoID              string `json:"moId"`
		NumberRefferal    string `json:"numberRefferal"`
		PatientBirthdate  int    `json:"patientBirthdate"`
		PatientFirstName  string `json:"patientFirstName"`
		PatientLastName   string `json:"patientLastName"`
		PatientMiddleName string `json:"patientMiddleName"`
		PatientOms        string `json:"patientOms"`
		PatientPassport   string `json:"patientPassport"`
		PatientSnils      string `json:"patientSnils"`
		Phone             string `json:"phone"`
		RecordSource      string `json:"recordSource"`
		ResourceID        string `json:"resourceId"`
		Rri               string `json:"rri"`
		RriOid            string `json:"rriOid"`
		ServiceSpecID     string `json:"serviceSpecId"`
		SessionID         string `json:"sessionId"`
		SlotID            string `json:"slotId"`
		StartDateRange    string `json:"startDateRange"`
		StartTimeRange    string `json:"startTimeRange"`
		VisitDate         string `json:"visitDate"`
		VisitTime         string `json:"visitTime"`
	}

	GetServiceSpecsInfoRequest struct {
		MO_Id  string `json:"MO_Id"`
		Reg_Id string `json:"Reg_Id"`
		// Session_ID string `json:"Session_ID"`
	}

	FindSpecsRequest struct {
		GetServiceSpecsInfoRequest `json:"GetServiceSpecsInfoRequest"`
	}

	FindSpecsResponse struct {
		GetServiceSpecsInfoResponse struct {
			Error struct {
				ErrorDetail struct {
					ErrorCode    int    `json:"errorCode"`
					ErrorMessage string `json:"errorMessage"`
				} `json:"errorDetail"`
			} `json:"Error"`
			ListServiceSpecs struct {
				ServiceSpec []struct {
					ServiceSpecID   int    `json:"ServiceSpec_Id"`
					ServiceSpecName string `json:"ServiceSpec_Name"`
				} `json:"ServiceSpec"`
			} `json:"ListServiceSpecs"`
			SessionID string `json:"Session_ID"`
		} `json:"GetServiceSpecsInfoResponse"`
	}

	GetResourceInfoRequest struct {
		RecordSource  string `json:"Record_Source"`
		ServiceSpecID string `json:"ServiceSpec_Id"`
		SessionID     string `json:"Session_ID"`
	}

	GetResourceRequest struct {
		GetResourceInfoRequest `json:"GetResourceInfoRequest"`
	}

	GetResourceResponse struct {
		GetResourceInfoResponse struct {
			Error struct {
				ErrorDetail struct {
					ErrorCode    int    `json:"errorCode"`
					ErrorMessage string `json:"errorMessage"`
				} `json:"errorDetail"`
			} `json:"Error"`
			ListResource struct {
				Resource []struct {
					ResourceID   interface{} `json:"Resource_Id"`
					ResourceName string      `json:"Resource_Name"`
				} `json:"Resource"`
			} `json:"ListResource"`
			SessionID string `json:"Session_ID"`
		} `json:"GetResourceInfoResponse"`
	}

	GetScheduleInfoRequest struct {
		EndDateRange   string `json:"EndDateRange"`
		EndTimeRange   string `json:"EndTimeRange"`
		RecordSource   string `json:"Record_Source"`
		ResourceID     string `json:"Resource_Id"`
		SessionID      string `json:"Session_ID"`
		StartDateRange string `json:"StartDateRange"`
		StartTimeRange string `json:"StartTimeRange"`
	}

	GetScheduleRequest struct {
		GetScheduleInfoRequest `json:"GetScheduleInfoRequest"`
	}

	GetScheduleResponse struct {
		GetScheduleInfoResponse struct {
			Error struct {
				ErrorDetail struct {
					ErrorCode    int    `json:"errorCode"`
					ErrorMessage string `json:"errorMessage"`
				} `json:"errorDetail"`
			} `json:"Error"`
			Schedule struct {
				Slots []struct {
					SlotID    string `json:"Slot_Id"`
					VisitTime string `json:"VisitTime"`
				} `json:"Slots"`
			} `json:"Schedule"`
			SessionID string `json:"Session_ID"`
		} `json:"GetScheduleInfoResponse"`
	}

	Spec struct {
		Spec string `json:"Spec"`
	}

	GetScheduleTableRequest struct {
		DateFrom     string `json:"DateFrom"`
		DateTo       string `json:"DateTo"`
		ListSpecs    []Spec `json:"ListSpecs"`
		RecordSource string `json:"RecordSource"`
		RegID        string `json:"RegId"`
	}

	GSTRequest struct {
		GetScheduleTableRequest `json:"GetScheduleTableRequest"`
	}

	GSTResponse struct {
		GetScheduleTableResponse struct {
			Error struct {
				ErrorDetail struct {
					ErrorCode    int    `json:"errorCode"`
					ErrorMessage string `json:"errorMessage"`
				} `json:"errorDetail"`
			} `json:"Error"`
			ListScheduleRecord struct {
				ScheduleRecord []struct {
					Area            string      `json:"Area"`
					Cabinet         int         `json:"Cabinet"`
					DoctorID        interface{} `json:"DoctorId"`
					DoctorName      string      `json:"DoctorName"`
					DoctorSpec      int         `json:"DoctorSpec"`
					ListDateRecords struct {
						DateRecords []struct {
							AllRecords  int    `json:"AllRecords"`
							Day         string `json:"Day"`
							FreeRecords int    `json:"FreeRecords"`
							Time        string `json:"Time"`
						} `json:"DateRecords"`
					} `json:"ListDateRecords"`
				} `json:"ScheduleRecord"`
			} `json:"ListScheduleRecord"`
		} `json:"GetScheduleTableResponse"`
	}

	GetRecordsWithErzResponse struct {
		BookMisID      string `json:"bookMisId"`
		Cabinet        string `json:"cabinet"`
		Canceled       bool   `json:"canceled"`
		Creator        string `json:"creator,string"`
		DepartOid      string `json:"departOid"`
		DoctorID       string `json:"doctorId"`
		DoctorName     string `json:"doctorName"`
		Email          string `json:"email"`
		ErzCode        string `json:"erzCode"`
		FirstName      string `json:"firstName"`
		IDPost         string `json:"idPost"`
		Inserted       string `json:"inserted"`
		LastName       string `json:"lastName"`
		MiddleName     string `json:"middleName"`
		MoID           string `json:"moId"`
		MoName         string `json:"moName"`
		NamePost       string `json:"namePost"`
		NumberRefferal string `json:"numberRefferal"`
		Phone          string `json:"phone"`
		RecordSource   string `json:"recordSource"`
		RejectReason   string `json:"rejectReason"`
		Rri            string `json:"rri"`
		ServiceID      string `json:"serviceId"`
		ServiceName    string `json:"serviceName"`
		SlotID         string `json:"slotId"`
		Status         string `json:"status"`
		Updated        string `json:"updated"`
		UploadStatus   string `json:"uploadStatus"`
		VisitDate      string `json:"visitDate"`
		VisitTime      string `json:"visitTime"`
	}
)

func NewClient(firstName, lastName, middleName, birthDate, oms string, httpClient *http.Client) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	cli := &Client{
		FirstName:    firstName,
		LastName:     lastName,
		MiddleName:   middleName,
		BirthDate:    birthDate,
		OmsNumber:    oms,
		RecordSource: "intermed",
		http:         httpClient,
	}
	return cli
}

func (cli *Client) Login() (loginResponse *LoginResponse, err error) {
	u, err := url.Parse(urlBase)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			cli.savedCookie = cookie
		}
	}

	loginPayload := LoginRequest{
		GetPatientInfoRequest: GetPatientInfoRequest{
			PatientData: PatientData{
				OmsNumber:    cli.OmsNumber,
				FirstName:    cli.FirstName,
				LastName:     cli.LastName,
				MiddleName:   cli.MiddleName,
				BirthDate:    cli.BirthDate,
				RecordSource: "intermed",
			},
		},
	}

	respBody, err := cli.DoRequest("POST", urlFindPatient, toJSON(loginPayload), "")
	log.Printf("respBody - %s", respBody)

	loginResponse = &LoginResponse{}
	err = json.Unmarshal([]byte(respBody), loginResponse)

	cli.SessionID = loginResponse.GetPatientInfoResponse.SessionID
	cli.PatientID = loginResponse.GetPatientInfoResponse.PatientID

	return
}

func (cli *Client) DoRequest(method string, url string, body string, query string) (respBody string, err error) {
	br := strings.NewReader(body)
	req, err := http.NewRequest(method, url, br)
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "intermed76.ru")
	req.Header.Add("Origin", "https://intermed76.ru")
	req.Header.Add("Referer", "https://intermed76.ru/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:81.0) Gecko/20100101 Firefox/81.0")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	if query != "" {
		req.URL.RawQuery = query
	}

	if cli.savedCookie != nil {
		req.AddCookie(cli.savedCookie)
	}
	reqBody, _ := httputil.DumpRequest(req, true)
	log.Printf("req - %s", string(reqBody))

	resp, err := cli.http.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respdBody, _ := httputil.DumpResponse(resp, true)
	log.Printf("resp - %s", string(respdBody))

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	respBody = doc.Find("body").Text()

	return
}

func (cli *Client) GetMos() (mosResponse *MosResponse, err error) {
	q := url.Values{}
	q.Add("sessionID", cli.SessionID)
	respBody, err := cli.DoRequest("GET", urlGetMos, "", q.Encode())
	mosResponse = &MosResponse{}
	err = json.Unmarshal([]byte(respBody), &mosResponse)

	return
}

func (cli *Client) GetSession() (gsResponse GetSessionResponse, err error) {
	q := url.Values{}
	q.Add("sessionId", cli.SessionID)
	respBody, err := cli.DoRequest("GET", urlGetSession, "", q.Encode())
	err = json.Unmarshal([]byte(respBody), &gsResponse)
	cli.session = gsResponse
	cli.SessionID = gsResponse.SessionID

	return
}

func (cli *Client) FindSpecs(oid string, idMo string) (fsResponse *FindSpecsResponse, err error) {
	fsPayload := FindSpecsRequest{
		GetServiceSpecsInfoRequest: GetServiceSpecsInfoRequest{
			MO_Id:  oid,  // Oid
			Reg_Id: idMo, // IDMo
			// Session_ID: loginResponse.GetPatientInfoResponse.SessionID,
		},
	}

	fsPayloadStr := toJSON(fsPayload)

	respBody, err := cli.DoRequest("POST", urlFindSpecs, fsPayloadStr, "")

	err = json.Unmarshal([]byte(respBody), &fsResponse)

	return
}

func (cli *Client) FindResourcesWithSession(specId string, sessionID string) (grResponse *GetResourceResponse, err error) {
	grPayload := GetResourceRequest{
		GetResourceInfoRequest: GetResourceInfoRequest{
			RecordSource:  "intermed",
			ServiceSpecID: specId,
			SessionID:     sessionID,
		},
	}

	grPayloadStr := toJSON(grPayload)

	respBody, err := cli.DoRequest("POST", urlFindResources, grPayloadStr, "")

	err = json.Unmarshal([]byte(respBody), &grResponse)
	log.Printf("ERROR code '%d'", grResponse.GetResourceInfoResponse.Error.ErrorDetail.ErrorCode)
	log.Printf("ERROR message '%s'", grResponse.GetResourceInfoResponse.Error.ErrorDetail.ErrorMessage)

	return
}

func (cli *Client) FindResources(specId string) (grResponse *GetResourceResponse, err error) {
	return cli.FindResourcesWithSession(specId, cli.SessionID)
}

func (cli *Client) GetSchedule() (gsResponse *GetScheduleResponse, err error) {
	gsPayload := GetScheduleRequest{
		GetScheduleInfoRequest: GetScheduleInfoRequest{
			SessionID:      cli.SessionID,
			RecordSource:   "intermed",
			ResourceID:     "04765670894",
			StartDateRange: "2020-03-01",
			EndDateRange:   "2020-04-11",
			StartTimeRange: "00:00:00.000+03:00",
			EndTimeRange:   "23:59:00.000+03:00",
		},
	}

	gsPayloadStr := toJSON(gsPayload)

	respBody, err := cli.DoRequest("POST", urlFindSchedules, gsPayloadStr, "")

	err = json.Unmarshal([]byte(respBody), &gsResponse)

	return
}

func (cli *Client) GetScheduleTable() (gstResponse *GSTResponse, err error) {
	sp := Spec{Spec: "22"}

	gstPayload := GSTRequest{
		GetScheduleTableRequest: GetScheduleTableRequest{
			RecordSource: "epgu",
			RegID:        "10304",
			DateFrom:     "2020-03-01",
			DateTo:       "2020-03-30",
		},
	}
	gstPayload.GetScheduleTableRequest.ListSpecs = append(gstPayload.GetScheduleTableRequest.ListSpecs, sp)

	gstPayloadStr := toJSON(gstPayload)

	respBody, err := cli.DoRequest("POST", urlFindSchedulesTable, gstPayloadStr, "")

	err = json.Unmarshal([]byte(respBody), &gstResponse)

	return
}

func (cli *Client) GetRecords() (records []GetRecordsWithErzResponse, err error) {
	now := time.Now()

	q := url.Values{}
	q.Add("rri", cli.session.Rri)
	q.Add("erzCode", cli.session.ErzCode)
	q.Add("accepted", "true")
	q.Add("declined", "false")                   // for total: true
	q.Add("deleted", "false")                    // for total: true
	q.Add("beginDate", now.Format("2006-01-02")) // for total: empty
	// q.Add("beginDate", "") // for total: empty
	q.Add("endDate", "")

	respBody, err := cli.DoRequest("GET", urlGetRecordsWithErz, "", q.Encode())

	err = json.Unmarshal([]byte(respBody), &records)

	return
}

func (cli *Client) CreateAppointment() {
	// POST https://intermed76.ru/intermed/createAppointment
	// {"CreateAppointmentRequest":{"Patient_Data":{"OMS_Number":"7690299770000000","First_Name":"Имя","Last_Name":"Фамилия","Middle_Name":"Отчество","Birth_Date":"2000-01-01"},"Session_ID":"eb9e6bed-0000-0000-0000-000000000000","Slot_Id":"d8ffda1b-91cb-4917-acca-1adb012f01ad"}}
}

// toJSON convert object to json notation
func toJSON(o interface{}) string {
	oj, _ := json.Marshal(o)
	return string(oj)
}

// fromJSON convert from json notation to object
func fromJSON(s string, r interface{}) (o interface{}) {
	_ = json.Unmarshal([]byte(s), o)
	return
}
