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

func main(){
  err := shim.Start(new(SimpleChaincode))
}
func  Init() {
  }
func InvokeOnTrustedPeer(executor, enc_kvs, asset,from, to, enc_price) (bool, int, int){
       amountFrom := enc_kvs.GetState(from)
       updated_amtF := enc_minus(amountFrom, enc_price)
       amountTo := enc_kvs.GetState(to)
       updated_amtT := enc_plus(amountTo, enc_price)
       return enc_geq(a, 0), hash(updated_amtF), hash(updated_amtT)
}
func InvokePrivate(executor, pvt_kvs, asset, from, to, price) (int){
	     if executor == from{
         updatedAmount := pvt_kvs.GetState(from) - price
	        new_from := pvt_kvs.PutState(from, updatedAmount)
           return hash(enc(new_from, key))
          }else if(executor == to){
            updatedAmount := pvt_kvs.GetState(to) + price
            new_to := pvt_kvs.PutState(to, updatedAmount)
            return hash(enc(new_to, key))
            }else{ 				                        //do nothing
              }}
func Invoke(function, executor, pub_kvs, enc_kvs, pvt_kvs, asset, from, to, price, enc_price) (int){
  if function == "move"{
    if executor == "TP"{
      return t.InvokeOnTrustedPeer(stub, args)
    }else{
      return t.InvokePrivate(stub, args)
    }
    if pub_kvs.GetState(asset) == to{
          pub_kvs.PutState(asset, from)
    }else{    		//do nothing
    }
}
}
func enc_minus(b int, c int) int{
    return b-c
}
func enc_geq(d int, e int) bool{
    if d>e {
      return true
      }
      return false
}
