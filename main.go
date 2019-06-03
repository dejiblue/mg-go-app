package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/matryer/way"
	_ "github.com/motemen/go-loghttp/global"
	"github.com/antchfx/xquery/xml"
)

// The Function Run When The Server Starts
func main() {
	log.Println("We Started Windows MDM GoLang Server On Port 8000")
	router := way.NewRouter()

	// The HTTP Routes
	router.HandleFunc("GET", "/", indexHandler)                                                       // main.go
	router.HandleFunc("GET", "/EnrollmentServer/Discovery.svc", discoveryGETHandler)                  // discovery.go
	router.HandleFunc("POST", "/EnrollmentServer/Discovery.svc", discoveryPOSTHandler)                // discovery.go
	//router.HandleFunc("POST", "/EnrollmentServer/PolicyService.svc", enrollmentPolicyHandler)         // enrollment.go
	//router.HandleFunc("POST", "/EnrollmentServer/EnrollmentService.svc", enrollmentWebServiceHandler) // enrollment.go
	//router.HandleFunc("POST", "/EnrollmentServer/DeviceEnrollment.svc", ) // enrollment.go

	router.NotFound = http.HandlerFunc(notFoundHandler) // main.go

	// Start The HTTP Server Listening
	log.Fatalln(http.ListenAndServe(":8000", logRequest(router)))
}

// The Response To Access The Index Page (Just A Placeholder)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Demo Windows MDM Server!"))
}

// The Response To Known Web Routes
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // Get The Body From The Request
	log.Println(string(body))         // Print The Body To The Console

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found: 404"))
}

// The HTTP Request Logger
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Return a 200 status to show the device a MDM server exists
func discoveryGETHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Discovery GET Request")
	w.Write([]byte(""))
}

// Return the locations of the MDM server
func discoveryPOSTHandler(w http.ResponseWriter, r *http.Request) { // TODO: Handle The Device Trying To Join - Valid Windows Version, Authentication, etc
  log.Println("The Discovery POST Request")
	soapBody, err := xmlquery.Parse(r.Body)
	if err != nil {
		panic(err)
	}
	MessageID := xmlquery.FindOne(soapBody, "//s:Header/a:MessageID").InnerText()

	// Send The Response
	w.Write([]byte(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
       xmlns:a="http://www.w3.org/2005/08/addressing">
      <s:Header>
        <a:Action s:mustUnderstand="1">
          http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse
        </a:Action>
        <ActivityId>
          d9eb2fdd-e38a-46ee-bd93-aea9dc86a3b8
        </ActivityId>
        <a:RelatesTo>` + MessageID + `</a:RelatesTo>
      </s:Header>
      <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xmlns:xsd="http://www.w3.org/2001/XMLSchema">
        <DiscoverResponse
           xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment">
          <DiscoverResult>
            <AuthPolicy>OnPremise</AuthPolicy>
            <EnrollmentVersion>4.0</EnrollmentVersion>
            <EnrollmentPolicyServiceUrl>
              https://winmdm.stg.mobileguardian.com/winmdm/EnrollmentServer/PolicyService.svc
            </EnrollmentPolicyServiceUrl>
            <EnrollmentServiceUrl>
              https://winmdm.stg.mobileguardian.com/winmdm/EnrollmentServer/EnrollmentService.svc
            </EnrollmentServiceUrl>
          </DiscoverResult>
        </DiscoverResponse>
      </s:Body>
    </s:Envelope>`))

		// Print out response
		log.Println(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
	       xmlns:a="http://www.w3.org/2005/08/addressing">
	      <s:Header>
	        <a:Action s:mustUnderstand="1">
	          http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse
	        </a:Action>
	        <ActivityId>
	          d9eb2fdd-e38a-46ee-bd93-aea9dc86a3b8
	        </ActivityId>
	        <a:RelatesTo>` + MessageID + `</a:RelatesTo>
	      </s:Header>
	      <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	         xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	        <DiscoverResponse
	           xmlns="http://schemas.microsoft.com/windows/management/2012/01/enrollment">
	          <DiscoverResult>
	            <AuthPolicy>OnPremise</AuthPolicy>
	            <EnrollmentVersion>4.0</EnrollmentVersion>
	            <EnrollmentPolicyServiceUrl>
	              https://winmdm.stg.mobileguardian.com/winmdm/EnrollmentServer/PolicyService.svc
	            </EnrollmentPolicyServiceUrl>
	            <EnrollmentServiceUrl>
	              https://winmdm.stg.mobileguardian.com/winmdm/EnrollmentServer/EnrollmentService.svc
	            </EnrollmentServiceUrl>
	          </DiscoverResult>
	        </DiscoverResponse>
	      </s:Body>
	    </s:Envelope>`)
}
