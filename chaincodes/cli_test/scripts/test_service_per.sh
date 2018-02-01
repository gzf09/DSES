#!/usr/bin/env bash

# Detecting whether can import the header file to render colorful cli output
if [ -f ./header.sh ]; then
 source ./header.sh
elif [ -f scripts/header.sh ]; then
 source scripts/header.sh
else
 alias echo_r="echo"
 alias echo_g="echo"
 alias echo_b="echo"
fi

CHANNEL_NAME="$1"
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="60"}
COUNTER=0
MAX_RETRY=5
CC_PATH=github.com/inklabsfoundation/inkchain/examples/chaincode/go/service
ORDERER_CA=/opt/gopath/src/github.com/inklabsfoundation/inkchain/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo_b "Chaincode Path : "$CC_PATH
echo_b "Channel name : "$CHANNEL_NAME

verifyResult () {
    if [ $1 -ne 0 ] ; then
        echo_b "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
        echo_r "================== ERROR !!! FAILED to execute MVE =================="
        echo
        exit 1
    fi
}

issueToken(){
#    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","100","18","4230a12f5b0693dd88bb35c79d7e56a68614b199"]}' >log.txt
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n ascc -c '{"Args":["registerAndIssueToken","'$1'","1000000000","8","07caf88941eafcaaa3370657fccc261acb75dfba"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "Issue a new token using ascc has Failed."
    echo_g "===================== A new token has been successfully issued======================= "
    echo
}

makeTransfer(){
    echo_b "pls wait 5 secs..."
    sleep 5
    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","a5ff00eb44bf19d5dfbde501c90e286badb58df4","INK","30000"]}' -i "1" -z 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4 >log.txt
#    peer chaincode invoke -o orderer.example.com:7050  --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C ${CHANNEL_NAME} -n token -c '{"Args":["transfer","3c97f146e8de9807ef723538521fcecd5f64c79a","INK","10"]}' -i "1" -z bc4bcb06a0793961aec4ee377796e050561b6a84852deccea5ad4583bb31eebe >log.txt

    res=$?
    cat log.txt
    verifyResult $res "Make transfer has Failed."
    echo_g "===================== Make transfer success ======================= "
    echo
}

chaincodeQueryA () {
    echo_b "Attempting to Query account A's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","07caf88941eafcaaa3370657fccc261acb75dfba","INK"]}' >log.txt
#    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","4230a12f5b0693dd88bb35c79d7e56a68614b199","INK"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query account A Failed."
}

chaincodeQueryB () {
    echo_b "Attempting to  query account B's balance on peer "
    sleep 3
    peer chaincode query -C mychannel -n token -c '{"Args":["getBalance","a5ff00eb44bf19d5dfbde501c90e286badb58df4","INK"]}' >log.txt
    res=$?
    cat log.txt
    verifyResult $res "query account B Failed."
   
}

# for init user
serviceInvoke_AddUser(){
    peer chaincode invoke -C mychannel -n service --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["registerUser","'$1'","'"$2"'"]}' -i "10" -z $3 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "service invoke: addUser has Failed."
    echo_g "===================== service invoke successfully ======================= "
    echo
}

serviceQuery_User() {
    echo_b "Attempting to Query user "
    sleep 3
    peer chaincode query -C mychannel -n service -c '{"Args":["queryUser","'$1'"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query user: Failed."
}

# for init service
serviceInvoke_AddService(){
#    peer chaincode invoke -C mychannel -n service --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["registerService","'$1'","'$2'","'"$3"'","'$4'"]}' -i "10" -z $5 >&log.txt

    echo $1
    echo $2
    echo $3
    echo $4
    echo $5

    echo '{"Args":["registerService","'"$1"'","'"$2"'","'"$3"'","'"$4"'"]}'

    peer chaincode invoke -C mychannel -n service --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["registerService","'"$1"'","'"$2"'","'"$3"'","'"$4"'"]}' -i "10" -z $5 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "service invoke: addService has Failed."
    echo_g "===================== service invoke successfully ======================= "
    echo
}

serviceQuery_Service() {
    echo_b "Attempting to Query user "
#    sleep 3
    peer chaincode query -C mychannel -n service -c '{"Args":["queryService","'"$1"'"]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "query service: Failed."
}


serviceRangeQuery_Service() {
    echo_b "Attempting to RangeQuery service "
#    sleep 3
    peer chaincode query -C mychannel -n service -c '{"Args":["queryServiceByRange","",""]}' >log.txt

    res=$?
    cat log.txt
    verifyResult $res "rangeQuery service: Failed."
}

# test publish service
# for init service
serviceInvoke_PublishService(){
    peer chaincode invoke -C mychannel -n service --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -c '{"Args":["invalidateService","'$1'"]}' -i "10" -z $2 >&log.txt
    res=$?
    cat log.txt
    verifyResult $res "service invoke: addService has Failed."
    echo_g "===================== service invoke successfully ======================= "
    echo
}


#echo_b "=====================6.Issue a token using ascc========================"
#issueToken INK
#
#echo_b "=====================7.Transfer 10 amount of INK====================="
#makeTransfer
#
echo_b "=====================8.Query transfer result of From account====================="
chaincodeQueryA
#
#echo_b "=====================9.Query transfer result of To account====================="
#chaincodeQueryB
#
#echo_b "=====================0. Init for test Service====================="
#echo_b "=====================0.1 add 2 user======================="
#
#serviceInvoke_AddUser user1 "An active service developer from Tsinghua University." 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4
#serviceInvoke_AddUser user2 "An active service developer from Inklabs Foundation." 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5
#
#echo_b "=====================0.2 query 2 user======================="
#serviceQuery_User user1
#
#echo_b "=====================0.3 register 2 service======================="
#sleep 2
#serviceInvoke_AddService S1 Map "A service about map APIs." user1 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4
#sleep 2
#serviceInvoke_AddService S2 2ap "A service about 2map APIs." user2 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5
#
#
#echo_b "=====================0.3 register service======================="
#file=servicelist.csv
#IFS=","
#index=0
#beginNumber=12820
#endNumber=12900
#while read name type des
#do
#    index=`expr $index + 1`
#
#    tem=${des:0:`expr ${#des} - 1`}
#
#    if [ $index -gt $beginNumber ]
#    then
#        sleep 2
#        time serviceInvoke_AddService "$name" "$type" "$tem" user1 70698e364537a106b5aa5332d660e2234b37eebcb3768a2a97ffb8042dfe2fc4
#    fi
#
#    echo $index
#    echo 名称=$name
#    echo 描述=$des
#    echo 类型=$type
#
#    if [ $index -eq $endNumber ]
#    then
#        break
#    fi
#
#done <  $file

#serviceRangeQuery_Service

echo_b "=====================0.4 query service======================="
sleep 3
time serviceQuery_Service "thingk.me"
sleep 3
time serviceQuery_Service "indico"
sleep 3
time serviceQuery_Service "Zimba JavaScript"
sleep 3
time serviceQuery_Service "Zannel Open"
sleep 3
time serviceQuery_Service "Yi.tl"
#echo_b "=====================0.5 publis a service======================="
#serviceInvoke_PublishService S2 344c267e5acb2ac9107465fc85eba24cbb17509e918c3cc3f5098dddf42167e5
#
#
#echo_b "=====================0.6 query service======================="
#time serviceQuery_Service "Zoho"
#time serviceQuery_Service "eBay Shopping"
#time serviceQuery_Service S2


echo
echo_g "=====================All GOOD, MVE Test completed ===================== "
echo
exit 0

