# Learn Hyperledger

+ Pull the docker images

```
docker pull hyperledger/fabric-tools:x86_64-1.1.0
docker pull hyperledger/fabric-orderer:x86_64-1.1.0
docker pull hyperledger/fabric-peer:x86_64-1.1.0
docker pull hyperledger/fabric-javaenv:x86_64-1.1.0
docker pull hyperledger/fabric-ccenv:x86_64-1.1.0
docker pull hyperledger/fabric-ca:x86_64-1.1.0
docker pull hyperledger/fabric-couchdb:x86_64-0.4.6
```

+ Tag them as latest

```
docker tag hyperledger/fabric-ca:x86_64-1.1.0 hyperledger/fabric-ca:latest
docker tag hyperledger/fabric-tools:x86_64-1.1.0 hyperledger/fabric-tools:latest
docker tag hyperledger/fabric-orderer:x86_64-1.1.0 hyperledger/fabric-orderer:latest
docker tag hyperledger/fabric-peer:x86_64-1.1.0 hyperledger/fabric-peer:latest
docker tag hyperledger/fabric-javaenv:x86_64-1.1.0 hyperledger/fabric-javaenv:latest
docker tag hyperledger/fabric-ccenv:x86_64-1.1.0 hyperledger/fabric-ccenv:latest
docker tag hyperledger/fabric-couchdb:x86_64-0.4.6 hyperledger/fabric-couchdb:latest
```

+ Verify images

```
docker images | grep hyperledger/fabric
```

+ Create a directory

```
mkdir multichain_network && cd multichain_network
```

+ Pull platform tools

```
curl -sSL https://goo.gl/6wtTN5 | bash -s 1.1.0
```

+ Create network related files
```
touch docker-compose.yml
touch crypto-config.yml
touch configtx.yaml
```

+ Set environmental variables
```
export CHANNEL_ONE_NAME=channelone
export CHANNEL_ONE_PROFILE=ChannelOne
export CHANNEL_TWO_NAME=channeltwo
export CHANNEL_TWO_PROFILE=ChannelTwo
export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

+ Generate certificates
```
./bin/cryptogen generate --config=./crypto-config.yml
```

+ Create `channel-artifacts` folder
```
mkdir channel-artifacts
```

+ Generate genesis block
```
./bin/configtxgen -profile OrdererGenesis -outputBlock ./channel-artifacts/genesis.block
```

+ Generate channel transactions
```
./bin/configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputCreateChannelTx ./channel-artifacts/${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME
./bin/configtxgen -profile ${CHANNEL_TWO_PROFILE} -outputCreateChannelTx ./channel-artifacts/${CHANNEL_TWO_NAME}.tx -channelID $CHANNEL_TWO_NAME
```

+ Generate anchor peers
```
./bin/configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors_${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME -asOrg Org1MSP
./bin/configtxgen -profile ${CHANNEL_ONE_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_${CHANNEL_ONE_NAME}.tx -channelID $CHANNEL_ONE_NAME -asOrg Org2MSP
./bin/configtxgen -profile ${CHANNEL_TWO_PROFILE} -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors_${CHANNEL_TWO_NAME}.tx -channelID $CHANNEL_TWO_NAME -asOrg Org2MSP
```

+ Start the network
```
docker-compose -f docker-compose.yml up -d
```

+ Verify the network is running
```
docker ps -a
```

+ Create channel one
```
docker exec cli peer channel create -o orderer.example.com:7050 -c $CHANNEL_ONE_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/${CHANNEL_ONE_NAME}.tx --tls --cafile $ORDERER_CA
```

+ Join peers to channel
```
docker exec cli peer channel join -b ${CHANNEL_ONE_NAME}.block --tls --cafile $ORDERER_CA
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel join -b ${CHANNEL_ONE_NAME}.block --tls --cafile $ORDERER_CA
```

+ Update anchor peers in channel one
```
docker exec cli peer channel update -o orderer.example.com:7050 -c $CHANNEL_ONE_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/Org1MSPanchors_${CHANNEL_ONE_NAME}.tx --tls --cafile $ORDERER_CA 
docker exec  -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel update -o orderer.example.com:7050 -c $CHANNEL_ONE_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/Org2MSPanchors_${CHANNEL_ONE_NAME}.tx --tls --cafile $ORDERER_CA
```

+ Create channel two
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel create -o orderer.example.com:7050 -c $CHANNEL_TWO_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/${CHANNEL_TWO_NAME}.tx --tls --cafile $ORDERER_CA
```

+ Join peer to channel
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel join -b ${CHANNEL_TWO_NAME}.block --tls --cafile $ORDERER_CA
```

+ Update anchor peer in channel two
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer channel update -o orderer.example.com:7050 -c $CHANNEL_TWO_NAME -f /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/Org2MSPanchors_${CHANNEL_TWO_NAME}.tx --tls --cafile $ORDERER_CA
```

+ Export chaincode location
```
export FIRST_CHAINCODE_NAME="firstchaincode"
export FIRST_CHAINCODE_SRC="github.com/chaincode/one"
export SECOND_CHAINCODE_NAME="secondchaincode"
export SECOND_CHAINCODE_SRC="github.com/chaincode/two"
export CHAINCODE_VERSION="1.0"
```

+ Install chaincode one on peer one
```
docker exec cli peer chaincode install -n $FIRST_CHAINCODE_NAME -p $FIRST_CHAINCODE_SRC -v $CHAINCODE_VERSION
```

+ Instantiate chaincode one in peer of org one
```
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_ONE_NAME -c '{"Args":["1","one"]}' -n $FIRST_CHAINCODE_NAME -v $CHAINCODE_VERSION -P "OR('Org1MSP.member', 'Org2MSP.member')"
```

+ Install chaincode one on peer two
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode install -n $FIRST_CHAINCODE_NAME -p $FIRST_CHAINCODE_SRC -v $CHAINCODE_VERSION
```

+ Invoke chaincode one on peer of org two
```
docker exec -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_ONE_NAME -c '{"Args":["set","two","2"]}' -n $FIRST_CHAINCODE_NAME -v $CHAINCODE_VERSION
```

+ Install chaincode two in peer of org two
```
docker exec  -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode install -n $SECOND_CHAINCODE_NAME -p $SECOND_CHAINCODE_SRC -v $CHAINCODE_VERSION
```

+ Instantiate chaincode two in peer of org two
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_TWO_NAME -c '{"Args":["aye","caramba"]}' -n $SECOND_CHAINCODE_NAME -v $CHAINCODE_VERSION -P "OR('Org1MSP.member', 'Org2MSP.member')"
```

+ Check if data is stored by executing below command
```
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.example.com:7051" cli peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C $CHANNEL_TWO_NAME -c '{"Args":["get","aye"]}' -n $SECOND_CHAINCODE_NAME -v $CHAINCODE_VERSION
```

> Congrats!! You've successfully set up a hyperledger fabric network and interacted with it.