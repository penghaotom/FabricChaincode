package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
//Smartcontract provides the function for managing
type SmartContract struct {
	contractapi.Contract
}
//发起人身份信息
type Voter struct{
	ID         string  `json:"id"`   //投票人ID
	Weight     int  `json:"weight"`   //投票人权重
	Voted      bool `json:"voted"`    //是否已经投票标记
	// vote       int  `json:"vote"`     //当前投票索引
}
//提议模板
type Proposal struct{
	Name       string `json:"name"`
	VoteCount  int    `json:"votecount"`    //已投票数
	Mold       string `json:"mold"`
}


//提案初始化
func (s *SmartContract) InitProposalLedger(ctx contractapi.TransactionContextInterface) error {
	proposals := []Proposal{
		{Name: "proposal1", Mold: "education", VoteCount: 0},
		{Name: "proposal2", Mold: "sport", VoteCount: 0},
		{Name: "proposal3", Mold: "military", VoteCount: 0},
		{Name: "proposal4", Mold: "entertainment", VoteCount: 0},
		{Name: "proposal5", Mold: "government",VoteCount: 0},
		{Name: "proposal6", Mold: "public", VoteCount: 0},
	}

	for _, proposal := range proposals {
		proposalJSON, err := json.Marshal(proposal)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(proposal.Name, proposalJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}
//投票人初始化
func (s *SmartContract) InitVoterLedger(ctx contractapi.TransactionContextInterface) error {
	voters := []Voter{
		{ID: "car1", Weight: 1, Voted: false},
		{ID: "car2", Weight: 1, Voted: false},
		{ID: "car3", Weight: 1, Voted: false},
		{ID: "car4", Weight: 1, Voted: false},
		{ID: "car5", Weight: 1, Voted: false},
		{ID: "car6", Weight: 1, Voted: false},
	}

	for _, voter := range voters {
		voterJSON, err := json.Marshal(voter)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(voter.ID, voterJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

//检测账本中是否存在输入的键
func (s *SmartContract) ExistInquire(ctx contractapi.TransactionContextInterface, everything string) (bool,error){
	existsJSON, err :=ctx.GetStub().GetState(everything)
	if err != nil{
		return false, fmt.Errorf("can't read from world state: %v" , err)
	}
	return existsJSON != nil, nil
}


//投票委托
func (s *SmartContract) DelegateVote(ctx contractapi.TransactionContextInterface, ownerName string, delegatedName string) (error) {
	owner, err := s.GetVoter(ctx, ownerName) //获得owner的结构
	if err != nil {
		return nil
	}
	delegatedvoter, err := s.GetVoter(ctx, delegatedName) //获得delegatedvoter结构
	if err != nil {
		return  nil
	}

	//检测是否已投票
	if owner.Voted == true || delegatedvoter.Voted == true {
		return fmt.Errorf("the delegated voter %s has voted", delegatedName)
	}
	//检测是否投给本人
	if owner.ID == delegatedvoter.ID {
		return fmt.Errorf("cannot delegate the vote to yourself")
	}
	delegatedvoter.Weight = owner.Weight + delegatedvoter.Weight
	owner.Weight = 0
	owner.Voted = true

	ownerJSON, errOwner := json.Marshal(owner)
	if errOwner != nil {
		return  nil
	}
	delegatedvoterJSON, errdelegatedvoter := json.Marshal(delegatedvoter)
	if errdelegatedvoter != nil {
		return  nil
	}
	err = ctx.GetStub().PutState(ownerName, ownerJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	err = ctx.GetStub().PutState(delegatedName, delegatedvoterJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	return nil
	}

	//投票
func (s *SmartContract) Vote(ctx contractapi.TransactionContextInterface, id string, name string) error{
	voter, err := s.GetVoter(ctx, id)
	if err != nil {
	return err
}
	proposal, err := s.GetProposal(ctx, name)
	if err != nil {
	return err
}
	proposal.VoteCount = proposal.VoteCount + voter.Weight
	voter.Weight = 0
	voter.Voted = true
	voterJSON, err := json.Marshal(voter)
	if err != nil {
	return err
}
	proposalJSON, err := json.Marshal(proposal)
	if err != nil {
	return err
}
	err = ctx.GetStub().PutState(id,voterJSON)
	if err != nil {
	return fmt.Errorf("failed to put to world state. %v", err)
}
	err = ctx.GetStub().PutState(name,proposalJSON)
	if err != nil {
	return fmt.Errorf("failed to put to world state. %v", err)
}
	return nil
}
// GetAsset returns the basic asset with id given from the world state
func (s *SmartContract) GetVoter(ctx contractapi.TransactionContextInterface, id string) (*Voter, error) {
	existing, err:= ctx.GetStub().GetState(id)

	if existing == nil {
		return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", id)
	}

	ba := new(Voter)

	err = json.Unmarshal(existing, ba)

	if err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type Voter", id)
	}
	return ba, nil
}//GetVoter和 GetProposal应该可以写成一个函数
//从世界状态中返回所输入name的proposal
func (s *SmartContract) GetProposal(ctx contractapi.TransactionContextInterface, name string) (*Proposal, error) {
		existing, err:= ctx.GetStub().GetState(name)

		if existing == nil {
			return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", name)
		}

		ma := new(Proposal)

		err = json.Unmarshal(existing, ma)

		if err != nil {
			return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type Proposal", name)
		}
		return ma, nil
}
//计票,返回最大票数
func (s *SmartContract) CalculateVoteNumber(ctx contractapi.TransactionContextInterface) (int, error){
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
			return 0, err
		}
		defer resultsIterator.Close()
	//var proposals []*Proposal
	//var maxvote int
	maxvote := 0
	//var votedproposal string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		var proposal Proposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			return 0, err
		}
		if maxvote < proposal.VoteCount{
			maxvote = proposal.VoteCount
			//votedproposal := proposal.Name
		}

	}
	return maxvote, nil
}

//计票，返回得票最多的提案
/*func (s *SmartContract) CalculateVoteName(ctx contractapi.TransactionContextInterface, votedname string) (string, error){
	resultsIterator, err := ctx.GetStub().GetStateByRange("proposal1", "proposal6")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()
	//var proposals []*Proposal
	maxvote := 0
	var votedproposal string
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var proposal Proposal
		err = json.Unmarshal(queryResponse.Value, &proposal)
		if err != nil {
			return "", err
		}
		if maxvote < proposal.VoteCount{
			maxvote = proposal.VoteCount
			votedproposal = proposal.Name
		}
	}
	 votedJSON, err := json.Marshal(votedproposal)
	err = ctx.GetStub().PutState(votedname,votedJSON)
	if err != nil{
		return "",nil
	}
	return votedproposal, nil
}*/
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}


