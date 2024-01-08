package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Version               = "0.0.1"
	ServiceName           = "HelloWordrApp"
	ServiceDescription    = "Hello Word rApp for testing Non-RT RIC guide development of future rApps and demo purposes"
	ServiceDisplayName     = "Hello Word rApp"
	DefaultHostPMS         = "http://nonrtricgateway.nonrtric.svc.cluster.local:9090"
	BasePathPMS            = "/a1-policy/v2"
	DefaultHostRAppCatalog = "http://rappcatalogues
	ervice.nonrtric.svc.cluster.local:9085"
	BasePathRAppCatalog    = "/services"
	DefaultPolicyBodyPath  = "nonrtric-rapp-helloword/src/pihw_template.json"
	DefaultPolicyTypeID    = "2"
	DefaultPolicyID        = "1"
	DefaultRICID           = "ric4"
)

var (
	baseURLRAppCatalogue  = DefaultHostRAppCatalog + BasePathRAppCatalog
	baseURLPMS            = DefaultHostPMS + BasePathPMS
	typeToUse             = DefaultPolicyTypeID
	ricToUse              = DefaultRICID
	bodyTypeToUse         = DefaultPolicyTypeID
	bodyPathToUse         = DefaultPolicyBodyPath
	policyIDToUse         = DefaultPolicyID
	verbose               bool
)

func registerServiceRAppCatalogue() bool {
	completeURL := fmt.Sprintf("%s/%s", baseURLRAppCatalogue, ServiceName)
	headers := map[string]string{"content-type": "application/json"}
	body := map[string]string{
		"version":       Version,
		"display_name":  ServiceDisplayName,
		"description":   ServiceDescription,
	}

	resp, err := makeRequest(http.MethodPut, completeURL, body, headers)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to register rApp %s\n", resp)
		return false
	}
	return true
}

func getRICsFromAgent() map[string]interface{} {
	resp, err := makeRequest(http.MethodGet, baseURLPMS+"/rics", nil, nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to get Rics %d\n", resp.StatusCode)
		return map[string]interface{}{}
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding RICs response:", err)
		return map[string]interface{}{}
	}
	return result
}

func putPolicy(ricName, policyID string) bool {
	threshold := rand.Intn(10^15 - 1)
	completeURL := baseURLPMS + "/policies"
	headers := map[string]string{"content-type": "application/json"}

	policyObj := make(map[string]interface{})
	err := json.Unmarshal([]byte(strings.ReplaceAll(policyData, "XXX", strconv.Itoa(threshold))), &policyObj)
	if err != nil {
		fmt.Println("Error unmarshaling policy data:", err)
		return false
	}

	body := map[string]interface{}{
		"ric_id":        ricName,
		"policy_id":     policyID,
		"service_id":    ServiceName,
		"policy_data":   policyObj,
		"policytype_id": typeToUse,
	}

	resp, err := makeRequest(http.MethodPut, completeURL, body, headers)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to create policy %s\n", resp)
		return false
	}

	fmt.Printf("Updating policy: %s threshold now: %d\n", policyID, threshold)
	return true
}

func getPolicyInstances() interface{} {
	completeURL := fmt.Sprintf("%s/policy-instances", baseURLPMS)
	resp, err := http.Get(completeURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Unable to get policy %s\n", resp)
		return false
	}

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding policy instances response:", err)
		return false
	}
	return result
}

// Função makeRequest auxiliar para fazer solicitações HTTP genéricas
func makeRequest(method, url string, body map[string]interface{}, headers map[string]string) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}

func main() {
	flag.BoolVar(&verbose, "v", false, "Turn on verbose printing")
	rAppHost := flag.String("rAppHost", "", "The host of the A1 PMS, e.g., http://localhost:9085")
	ricID := flag.String("ricId", "", "The ID of the policy type to use")
	pmsHost := flag.String("pmsHost", "", "The host of the A1 PMS, e.g., http://localhost:8081")
	policyID := flag.String("policyId", "", "The ID of the policy to use")
	policyTypeID := flag.String("policyTypeId", "", "The ID of the policy type to use")
	policyBodyPath := flag.String("policyBodyPath", "", "The path to the JSON body of the policy to create")

	flag.Parse()

	if *rAppHost != "" {
		baseURLRAppCatalogue = *pmsHost + BasePathRAppCatalogue
	}

	if *pmsHost != "" {
		baseURLPMS = *pmsHost + BasePathPMS
	}

	if *ricID != "" {
		ricToUse = *ricID
	}

	if *policyID != "" {
		policyIDToUse = *policyID
	}

	if *policyTypeID != "" {
		typeToUse = *policyTypeID
	}

	if *policyBodyPath != "" {
		bodyTypeToUse = *policyTypeID
	} else {
		bodyTypeToUse = DefaultPolicyBodyPath
	}

	rand.Seed(time.Now().UnixNano())

	// Registra o serviço no catálogo rApp
	fmt.Printf("Registering in rApp catalog %s\n", ServiceName)
	registerServiceRAppCatalogue()

	// Exibe os tipos de política disponíveis
	fmt.Println("Policy Types:")
	fmt.Println(getPolicyTypes())

	// Tenta obter informações do A1 Policy Manager
}