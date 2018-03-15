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
  "encoding/json"
  "fmt"
  "reflect"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  //"strconv"
  pb "github.com/hyperledger/fabric/protos/peer"
  "github.com/syndtr/goleveldb/leveldb"
  //"container/list"
  "crypto/sha256"
)

type SimpleChaincode struct{
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

    //set
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
    if function == "proposal.Move"{
    //args[0]じゃなくてMockInitの返り値
    if (args[0] == "TP"){ //executorがTPの場合
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
    }else {           //executorがTP以外の場合
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
      stub.PutState(args[2], []byte(args[3])) //to, hash_enc_amtB
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
