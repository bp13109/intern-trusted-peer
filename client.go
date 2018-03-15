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

type MyType struct{
  function string `json:"function"`
  who string `json:"who"`
  from string `json:"from"`
  to string `json:"to"`
  price string  `json:"price"`
}
// ============================================================================================================================
// Main
// ===========================================================================================================================
func main(){

  mock := shim.NewMockStub("test", new(SimpleChaincode))

  //本当は1つの関数で実行したいが、今回は別々に。

  var mt MyType
  mt = MyType{"proposal.Move", "TP", "A", "B", enc("10", "key")}
  respA := mock.MockInvoke("proposal.Move", [][]byte {[]byte("proposal.Move"), []byte("A"), []byte("A"), []byte("B"), []byte("10")}) //, "A", from, to, price ->10
  fmt.Println("sent proposal to A")
  respB := mock.MockInvoke("proposal.Move", [][]byte {[]byte("proposal.Move"), []byte("B"), []byte("A"), []byte("B"), []byte("10")}) //, "B", from, to, price
  fmt.Println("sent proposal to B")
  b, _ := json.Marshal(&mt)
  respT := mock.MockInvoke("proposal.Move", [][]byte{b})
  //respT := mock.MockInvoke("proposal.Move", [][]byte{json.Marshal(&mt)})
  fmt.Println(string(b))
  fmt.Println(b)
  fmt.Println(respT)
  fmt.Println(&respA)
  fmt.Println(sha256.Sum256([]byte("a")))
  respTP := mock.MockInvoke("proposal.Move", [][]byte{[]byte("proposal.Move"), []byte("TP"), []byte("A"), []byte("B"), []byte(enc("10", "key")), []byte("ok/err")})
  fmt.Println("sent proposal to TP")

    //validation
    fmt.Println("A:"+valid(respA, respTP))
    fmt.Println("B:"+valid(respB, respTP))

    //set or broadcast
    var set string = "set"
    mock.MockInvoke(set, [][]byte {[]byte(set), []byte("A"), []byte("B"), []byte("hash_enc_amtA"), []byte("hash_enc_amtB")})
    fmt.Println("broadcast proposal to all as usual")
}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
  return shim.Success(nil)
}
// ============================================================================================================================
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
  function, args := stub.GetFunctionAndParameters()
  fmt.Println("starting invoke, for - " + function)
  // ============================================================================================================================
  // proposal.Move
  //queryと一緒にして、ローカルkvsに書き込む
    if function == "proposal.Move"{               //proposal.Moveみたいな
    //args[0]じゃなくてMockInitの返り値
    if (args[0] == "TP"){
      if(args[4] == "ok/err"){

        //まずenc_kvsに書き込む
        db, err := leveldb.OpenFile("enc_kvs", nil)
        if err != nil{
          fmt.Println("error", err)
        }
        defer db.Close()
        //fromの値の更新
        data1, err := db.Get([]byte(args[1]), nil)
        temp1 := string(data1) + "+" + args[3]
        err = db.Put([]byte(args[1]), []byte(temp1), nil)
        new_data1, err := db.Get([]byte(args[1]), nil)
        fmt.Println("Update encryption kvs --> "+ args[1]+":"+string(new_data1))
        //toの値の更新
        data2, err := db.Get([]byte(args[2]), nil)
        temp2 := string(data2) + "+" + args[3]
        err = db.Put([]byte(args[2]), []byte(temp2), nil)
        new_data2, err := db.Get([]byte(args[2]), nil)
        fmt.Println("Update encryption kvs --> "+ args[2]+":"+string(new_data2))

      return shim.Success([]byte(hash(string(new_data1))))
    }
    }else {                           //AかBだった場合
    //pvt_kvsに書き込む
    db, err := leveldb.OpenFile("pvt_kvs", nil)
    if err != nil{
      fmt.Println("error", err)
    }
    defer db.Close()
    data, err := db.Get([]byte("pv_amt"+args[0]), nil)
    temp := string(data) + "+" + args[3]
    err = db.Put([]byte("pv_amt"+args[0]), []byte(temp), nil)
    new_data, err := db.Get([]byte("pv_amt"+args[0]), nil)
    fmt.Println("Update "+args[0]+"'s private kvs --> "+"pv_amt"+args[0]+":"+string(new_data))
    return shim.Success([]byte(hash(enc(string(new_data), "key"))))
}
  // ============================================================================================================================
  // set
    }else if(function == "set"){
      stub.PutState(args[1], []byte(args[2])) //from, hash_enc_amtA
      stub.PutState(args[2], []byte(args[3]))
}
  return shim.Success(nil)
}
func valid(resp, respTP pb.Response) string{
  var amtA = resp
  var amtB = respTP
  //fmt.Println(args[0])
  if(reflect.DeepEqual(amtA, amtB) == true){
    return "Validation is done"
  }else {
  return "Validation is failed"
 }
}

/*
var query string = "query"
//hash_enc_amountA, B及びOk/errをそれぞれのPeerから受け取る
  respA := mock.MockInvoke(query, [][]byte {[]byte(query), []byte("A"), []byte("TP"), []byte("hash_enc_amtA")})
  fmt.Println("received proposal response from A")
  respB := mock.MockInvoke(query, [][]byte {[]byte(query), []byte("B"), []byte("TP"), []byte("hash_enc_amtB")})
  fmt.Println("received proposal response from B")
  respTP_A := mock.MockInvoke(query, [][]byte {[]byte(query), []byte("TP_A"), []byte("hash_enc_amtA"), []byte("hash_enc_amtB"), []byte("ok/err")})
  fmt.Println("received proposal response from TP")
  respTP_B := mock.MockInvoke(query, [][]byte {[]byte(query), []byte("TP_B"), []byte("hash_enc_amtA"), []byte("hash_enc_amtB"), []byte("ok/err")})
  fmt.Println("received proposal response from TP")
*/
/*
//make hogeDB
db, err := leveldb.OpenFile("pub_kvs", nil)
if err != nil{
  fmt.Println("error", err)
}
defer db.Close()

err = db.Put([]byte("key"), []byte("value"), nil)
data, err := db.Get([]byte("key"), nil)
fmt.Println(string(data))
*/
// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
/*
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
  funcName, args := stub.GetFunctionAndParameters()
  var number int
	var err error
	txId := stub.GetTxID()

  fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)

  // expecting 1 arg for instantiate or upgrade
	if len(args) == 1 {
		fmt.Println("  GetFunctionAndParameters() arg[0] length", len(args[0]))

    // expecting arg[0] to be length 0 for upgrade
		if len(args[0]) == 0 {
			fmt.Println("  Uh oh, args[0] is empty...")
		} else {
			fmt.Println("  Great news everyone, args[0] is not empty")

      // this is a very simple test. let's write to the ledger and error out on any errors
			// it's handy to read this right away to verify network is healthy if it wrote the correct value
			err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
			if err != nil {
				return shim.Error(err.Error())                  //self-test fail
			}
		}
	}
  // showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

  fmt.Println("Ready for action")                          //self-test pass
  return shim.Success(nil)
  }
/*
// ============================================================================================================================
// InvokeOnTrustedPeer - Trusted peer confirms balance of receiver
// ============================================================================================================================

//Process of verification on trusted peer
func (t *SimpleChaincode) InvokeOnTrustedPeer(stub shim.ChaincodeStubInterface, args []string) pb.Response{
  function, args := stub.GetFunctionAndParameters()
  fmt.Println(" ")
	fmt.Println("starting invoke on trusted peer, for - " + function)
    //function, executor, enc_kvs, asset,| from, to, enc_price
    //ここのkvsはTrusted Peerのencryptedなkvs
       var amountFrom, a int
       amountFrom = enc_kvs.GetState(from) //enc_kvs.
       a = enc_minus(amountFrom, enc_price)
       return enc_geq(a, 0)  //boolean
}

// ============================================================================================================================
// InvokePrivate - Transactors simulates their transaction proposal
// ============================================================================================================================
func (t *SimpleChaincode) InvokePrivate(stub shim.ChaincodeStubInterface) pb.Response{
  //実際の引数は　function, executor, pvt_kvsまでがshim, asset, from, to, price
  //ここのkvsは各Peerのpvtなkvs
    function, args := stub.GetFunctionAndParameters()
    fmt.Println(" ")
  	fmt.Println("starting InvokePrivate, for - " + function)

	     if executor == from{
         updatedAmount = pvt_kvs.GetState(from) - price
	         pvt_kvs.PutState(from, updatedAmount)
          }else if(executor == to){
	     updatedAmount  = pvt_kvs.GetState(to) + price
       pvt_kvs.PutState(to, updatedAmount)
       }	else{ 				//do nothing
       }
}

// ============================================================================================================================
// Invoke - Public transaction
// ============================================================================================================================
//ownerの移動
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
  //引数はfunction, executor, pub_kvs, enc_kvs, pvt_kvs, asset, from, to, price, enc_price
  function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

  if function == "move"{
    if executor == "TP"{
      return t.InvokeOnTrustedPeer(stub, args)
    }else{
      return t.InvokePrivate(stub, args)
    }
    if pub_kvs.GetState(asset) == to{              //もし資産のオーナーがお金を受け取る人なら
          pub_kvs.PutState(asset, from)              //資産を支払う人に移動
        }		//do nothing
    }else if function == "init_asset"{               //資産の初期化
      return init_asset(stub, args)
    }else if function == "init_owner"{               //ownerの初期化
      return init_owner(stub, args)
    }else if function == "set_owner"{                //ownerをセット
      return set_owner(stub, args)
    }
}
*/
