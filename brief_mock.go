/*
Copyright 2018 IBM Corp.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
 package main

import(
  //"errors"
  "encoding/json"
  "fmt"
  "reflect"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  //"strconv"
  pb "github.com/hyperledger/fabric/protos/peer"
  "github.com/syndtr/goleveldb/leveldb"
  //"github.com/golang/protobuf/ptypes/timestamp"
  //"container/list"
  "crypto/sha256"
)

type SimpleChaincode struct{
}
type MockStub struct {
	// arguments the stub was called with
	args [][]byte

	// A pointer back to the chaincode that will invoke this, set by constructor.
	// If a peer calls this stub, the chaincode will be invoked from here.
	cc Chaincode

	// A nice name that can be used for logging
	Name string

	// State keeps name value pairs
	State map[string][]byte

	// Keys stores the list of mapped values in lexical order
	//Keys *list.List

	// registered list of other MockStub chaincodes that can be called from this MockStub
	Invokables map[string]*MockStub

	// stores a transaction uuid while being Invoked / Deployed
	// TODO if a chaincode uses recursion this may need to be a stack of TxIDs or possibly a reference counting map
	TxID string

	//TxTimestamp *timestamp.Timestamp

	// mocked signedProposal
	signedProposal *pb.SignedProposal
}

// ============================================================================================================================
// Asset Definitions - The ledger will store assets and owners
// ============================================================================================================================

// ----- Asset ----- //
type Asset struct{
  ObjectType string `json:"docType"`       //field for couchdb
  Id         string `json:"id"`
}

// ----- Owners ----- //
type Owner struct {
  ObjectType string `json:"docType"`       //field for couchdb
  Id         string `json:"id"`
  Username   string `json:"username"`
}

// ----- pub_kvs ----- //
type pub_kvs struct{
  Asset string `json:"asset"`
  Owner string `json:owner`
  Pub_amt  int `json:"value"`
}

// ----- pvt_kvs ----- //
type pvt_kvs struct{
  Pvt_amt int `json:"pvt_kvs"`
  PvtVal  int `json:"pvtVal"`
}

// ----- enc_kvs ----- //
type enc_kvs struct{
  Owner string `json:"id"`
  EncVal  int `json:"encVal"`
}

type Chaincode interface {
	// Init is called during Instantiate transaction after the chaincode container
	// has been established for the first time, allowing the chaincode to
	// initialize its internal data
	Init(stub shim.ChaincodeStubInterface) pb.Response

	// Invoke is called to update or query the ledger in a proposal transaction.
	// Updated state variables are not committed to the ledger until the
	// transaction is committed.
	Invoke(stub shim.ChaincodeStubInterface) pb.Response
}

func enc(x, key string) string{

  return x+key
}

func hash(x string) string{

  return "h("+x+")"
}
// ============================================================================================================================
// Main
// ===========================================================================================================================
func main(){
  mock := shim.NewMockStub("test", new(SimpleChaincode))
  respA := MockInvoke("proposal.Move", "A", "A", "B", "10")
  respB := MockInvoke("proposal.Move", "B", "A", "B", "10")
  respTP := MockInvoke("proposal.Move", "TP", "A", "B", enc("10", "key"), "ok/err")
  fmt.Println("A:"+valid(respA, respTP))  //validation
  fmt.Println("B:"+valid(respB, respTP))  //validation
  //set
  MockInvoke("set", "A", "B", "hash_enc_amtA", "hash_enc_amtB")
}
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
    if function == "proposal.Move"{
    if (args[0] == "TP"){
      if(args[4] == "ok/err"){
        db, err := leveldb.OpenFile("enc_kvs", nil) //まずenc_kvsに書き込む
        ・・・                 //値の更新
      return shim.Success([]byte(hash(string(new_data1))))
    }}else {
      db, err := leveldb.OpenFile("pvt_kvs", nil)
        ・・・
    return shim.Success([]byte(hash(enc(string(new_data), "key"))))
}}else if(function == "set"){
      stub.PutState(args[1], []byte(args[2]))
      stub.PutState(args[2], []byte(args[3]))
}
func valid(resp, respTP pb.Response) string{
  var amtA = resp
  var amtB = respTP
  if(reflect.DeepEqual(amtA, amtB) == true){
    return "Validation is done"
  }else {
  return "Validation is failed"
 }
}
