package arrowheadfunctions

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type System struct {
	Address            string `json:"address"`
	Port               int    `json:"port"`
	SystemName         string `json:"systemName"`
	AuthenticationInfo string `json:"authenticationInfo"`
}

type Metadata struct {
	Method string `json:"method"`
}

type RegisterSystemRequestDTO struct {
	Address    string `json:"address"`
	Port       int    `json:"port"`
	SystemName string `json:"systemName"`
}

type Service struct {
	Interfaces        []string `json:"interfaces"`
	Metadata          Metadata `json:"metadata"`
	ProviderSystem    System   `json:"providerSystem"`
	Secure            string   `json:"secure"`
	ServiceDefinition string   `json:"serviceDefinition"`
	ServiceUri        string   `json:"serviceUri"`
}

type Orchestrate struct {
	OrchestrationFlags OrchestrationFlag `json:"orchestrationFlags"`
	RequestedService   RequestedService  `json:"requestedService"`
	RequesterSystem    System            `json:"requesterSystem"`
}

type OrchestrateResponse struct {
	Provider   Provider `json:"provider"`
	ServiceUri string   `json:"serviceUri"`
}

type OrchResponse struct {
	Response []OrchestrateResponse `json:"response"`
}

type OrchestrationFlag struct {
	OverrideStore bool `json:"overrideStore"`
}

type RequestedService struct {
	InterfaceRequirements        []string `json:"interfaceRequirements"`
	ServiceDefinitionRequirement string   `json:"serviceDefinitionRequirement"`
}

type Provider struct {
	Address    string `json:"address"`
	Port       int    `json:"port"`
	SystemName string `json:"systemName"`
}

func Hello() {
	fmt.Println("Daniel-sama")
}

func EchoServiceRegistry(address string, port int, certFilePath string, keyFilePath string) ([]byte, error) {
	portSTR := strconv.Itoa(port)
	req, errCreateRequest := http.NewRequest("GET", "https://"+address+":"+portSTR+"/serviceregistry/echo", nil)
	if errCreateRequest != nil {
		return nil, errCreateRequest
	}
	client := GetClient(certFilePath, keyFilePath)
	resp, errDoRequest := client.Do(req)
	if errDoRequest != nil {
		return nil, errDoRequest
	}

	body, errReadingBody := io.ReadAll(resp.Body)
	if errReadingBody != nil {
		return nil, errReadingBody
	}
	return body, nil

}

func RemoveServices(servicesToBeRemoved []Service, address string, port int, certFilePath string, keyFilePath string) {
	for _, service := range servicesToBeRemoved {
		_, err := RemoveService(service, address, port, certFilePath, keyFilePath)
		if err != nil {
			fmt.Println("Error removing service: ", err)
		}
	}
}

func PublishServices(servicesToBeAdded []Service, address string, port int, certFilePath string, keyFilePath string) {
	for _, service := range servicesToBeAdded {
		PublishService(service, address, port, certFilePath, keyFilePath)
	}
}

func GetClient(certFile string, keyFile string) *http.Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panic("Certficate load error. ", err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	client := &http.Client{Transport: tr}
	return client
}

func PublishService(requestBody Service, address string, port int, certFilePath string, keyFilePath string) {
	portSTR := strconv.Itoa(port)
	payload, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://"+address+":"+portSTR+"/serviceregistry/register", bytes.NewBuffer(payload))

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := GetClient(certFilePath, keyFilePath)
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error making HTTP request using client. ", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic("Error reading HTTP response. ", err)
	}

	fmt.Println("## Response Body:\n", string(body))
	fmt.Println("## Response status:\n", resp.Status, resp.StatusCode)
}

func RemoveService(service Service, address string, port int, certFilePath string, keyFilePath string) ([]byte, error) {
	portSTR := strconv.Itoa(port)
	url := fmt.Sprintf("https://"+address+":"+portSTR+"/serviceregistry/unregister?address=%s&port=%s&service_definition=%s&service_uri=%s&system_name=%s", service.ProviderSystem.Address, strconv.Itoa(service.ProviderSystem.Port), service.ServiceDefinition, service.ServiceUri, service.ProviderSystem.SystemName)
	req, errCreateRequest := http.NewRequest("DELETE", url, nil)
	if errCreateRequest != nil {
		return nil, errCreateRequest
	}

	client := GetClient(certFilePath, keyFilePath)
	resp, errDoRequest := client.Do(req)
	if errDoRequest != nil {
		return nil, errCreateRequest
	}
	body, errReadingBody := io.ReadAll(resp.Body)
	if errReadingBody != nil {
		return nil, errReadingBody
	}
	return body, nil

}

func RegisterSystem(rsrDTO RegisterSystemRequestDTO, address string, port int, certFilePath string, keyFilePath string) {
	portSTR := strconv.Itoa(port)
	payload, err := json.Marshal(rsrDTO)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://"+address+":"+portSTR+"/serviceregistry/register-system", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := GetClient(certFilePath, keyFilePath)
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("Error making HTTP request using client. ", err)
	}
	fmt.Println("## Response Body:\n", resp.Body)
	fmt.Println("## Response status:\n", resp.Status, resp.StatusCode)
}

func Orchestration(requestBody Orchestrate, address string, port int, certFilePath string, keyFilePath string) []byte {
	portSTR := strconv.Itoa(port)
	payload, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://"+address+":"+portSTR+"/orchestrator/orchestration", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := GetClient(certFilePath, keyFilePath)
	resp, err := client.Do(req)
	fmt.Println("request sent")
	fmt.Println("requestbody: ", requestBody)
	if err != nil {
		log.Panic("Error making HTTP request using client. ", err)
	}

	body, err := io.ReadAll(resp.Body)
	body2 := string(body[:])

	fmt.Println("## Response Body:\n", body2)
	fmt.Println("## Response status:\n", resp.Status, resp.StatusCode)
	return body
}

func RemoveSystem(system System, address string, port int, certFilePath string, keyFilePath string) {
	portSTR := strconv.Itoa(port)
	url := fmt.Sprintf("https://"+address+":"+portSTR+"/serviceregistry/unregister-system?address=%s&port=%s&system_name=%s", system.Address, strconv.Itoa(system.Port), system.SystemName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := GetClient(certFilePath, keyFilePath)
	resp, err := client.Do(req)
	fmt.Println("request sent")
	if err != nil {
		log.Panic("Error making HTTP request using client. ", err)
	}

	fmt.Println("## Response status:\n", resp.Status, resp.StatusCode)
	fmt.Println("System deleted")
}
