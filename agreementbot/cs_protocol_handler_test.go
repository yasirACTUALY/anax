// +build unit

package agreementbot

import (
	"encoding/json"
	"github.com/open-horizon/anax/citizenscientist"
	"github.com/open-horizon/anax/policy"
	"os"
	"sync"
	"testing"
)

func Test_agreement_with_override(t *testing.T) {

	expectedType := policy.Ethereum_bc
	expectedName := "ethel"

	testProposal := `{"address":"123456","producerPolicy":"policy","consumerId":"ag12345","type":"proposal","protocol":"Citizen Scientist","version":1,"agreementId":"deadbeef"}`
	testPolicy := `{"header":{"name":"testpolicy","version":"1.0"},"agreementProtocols":[{"name":"Citizen Scientist"}]}`

	os.Setenv("CMTN_BLOCKCHAIN", "ethel")
	if ag, err := createAgreement(testProposal, testPolicy, 0, "", ""); err != nil {
		t.Errorf("Error creating mock agreement, %v", err)
	} else if bcType, bcName := createEmptyPH().GetKnownBlockchain(ag); bcType != expectedType || bcName != expectedName {
		t.Errorf("Wrong BC type %v and name %v returned, expecting %v %v", bcType, bcName, expectedType, expectedName)
	}
	os.Setenv("CMTN_BLOCKCHAIN", "")

}

func Test_old_agreement_success(t *testing.T) {

	expectedType := policy.Ethereum_bc
	expectedName := policy.Default_Blockchain_name

	testProposal := `{"address":"123456","producerPolicy":"policy","consumerId":"ag12345","type":"proposal","protocol":"Citizen Scientist","version":1,"agreementId":"deadbeef"}`
	testPolicy := `{"header":{"name":"testpolicy","version":"1.0"},"agreementProtocols":[{"name":"Citizen Scientist"}]}`

	if ag, err := createAgreement(testProposal, testPolicy, 0, "", ""); err != nil {
		t.Errorf("Error creating mock agreement, %v", err)
	} else if bcType, bcName := createEmptyPH().GetKnownBlockchain(ag); bcType != expectedType || bcName != expectedName {
		t.Errorf("Wrong BC type %v and name %v returned, expecting %v %v", bcType, bcName, expectedType, expectedName)
	}

}

func Test_new_agreement_success1(t *testing.T) {

	expectedType := "fred"
	expectedName := "ethel"

	testProposal := `{"address":"123456","producerPolicy":"policy","consumerId":"ag12345","type":"proposal","protocol":"Citizen Scientist","version":1,"agreementId":"deadbeef"}`
	testPolicy := `{"header":{"name":"testpolicy","version":"1.0"},"agreementProtocols":[{"name":"Citizen Scientist","blockchains":[{"type":"fred","name":"ethel"}]}]}`

	if ag, err := createAgreement(testProposal, testPolicy, 0, "", ""); err != nil {
		t.Errorf("Error creating mock agreement, %v", err)
	} else if bcType, bcName := createEmptyPH().GetKnownBlockchain(ag); bcType != expectedType || bcName != expectedName {
		t.Errorf("Wrong BC type %v and name %v returned, expecting %v %v", bcType, bcName, expectedType, expectedName)
	}

}

func Test_new_agreement_success2(t *testing.T) {

	expectedType := "fred"
	expectedName := "ethel"

	testProposal := `{"address":"123456","producerPolicy":"policy","consumerId":"ag12345","type":"proposal","protocol":"Citizen Scientist","version":1,"agreementId":"deadbeef"}`
	testPolicy := `{"header":{"name":"testpolicy","version":"1.0"},"agreementProtocols":[{"name":"Citizen Scientist","blockchains":[{"type":"fred","name":"ethel"},{"type":"lucy","name":"dezi"}]}]}`

	if ag, err := createAgreement(testProposal, testPolicy, 0, "", ""); err != nil {
		t.Errorf("Error creating mock agreement, %v", err)
	} else if bcType, bcName := createEmptyPH().GetKnownBlockchain(ag); bcType != expectedType || bcName != expectedName {
		t.Errorf("Wrong BC type %v and name %v returned, expecting %v %v", bcType, bcName, expectedType, expectedName)
	}

}

func Test_new_agreement_success3(t *testing.T) {

	expectedType := "fred"
	expectedName := "ethel"

	testProposal := `{"address":"123456","producerPolicy":"policy","consumerId":"ag12345","type":"proposal","protocol":"Citizen Scientist","version":1,"agreementId":"deadbeef"}`
	testPolicy := `{"header":{"name":"testpolicy","version":"1.0"},"agreementProtocols":[{"name":"Citizen Scientist","blockchains":[{"type":"fred","name":"ethel"},{"type":"lucy","name":"dezi"}]}]}`

	if ag, err := createAgreement(testProposal, testPolicy, 2, expectedType, expectedName); err != nil {
		t.Errorf("Error creating mock agreement, %v", err)
	} else if bcType, bcName := createEmptyPH().GetKnownBlockchain(ag); bcType != expectedType || bcName != expectedName {
		t.Errorf("Wrong BC type %v and name %v returned, expecting %v %v", bcType, bcName, expectedType, expectedName)
	}

}

// Utility to help create the testing context
func createEmptyPH() *CSProtocolHandler {
	return &CSProtocolHandler{
		BaseConsumerProtocolHandler: &BaseConsumerProtocolHandler{
			name:             "test",
			pm:               nil,
			db:               nil,
			config:           nil,
			httpClient:       nil,
			agbotId:          "ag12345",
			token:            "abcdefg",
			deferredCommands: nil,
			messages:         nil,
		},
		genericAgreementPH: nil,
		Work:               nil,
		bcState:            nil,
		bcStateLock:        sync.Mutex{},
	}
}

func createAgreement(proposal string, pol string, agpVersion int, bcType string, bcName string) (*Agreement, error) {
	if ag, err := agreement("testagid", "deviceid", "testpolicy", bcType, bcName, "Citizen Scientist"); err != nil {
		return nil, err
	} else {
		prop := new(citizenscientist.CSProposal)
		if err := json.Unmarshal([]byte(proposal), prop); err != nil {
			return nil, err
		} else {
			prop.TsandCs = pol
			if propString, err := json.Marshal(prop); err != nil {
				return nil, err
			} else {
				ag.Proposal = string(propString)
				ag.AgreementProtocolVersion = agpVersion
				return ag, nil
			}
		}
	}
}
