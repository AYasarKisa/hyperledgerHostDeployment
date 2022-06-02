const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')
const helper = require('./helper')
var mqtt = require('mqtt');
var client = mqtt.connect('mqtt://34.230.75.198:1234');
// const createTransactionEventHandler = require('./MyTransactionEventHandler.ts')

var org_name='';
var username='';
var chaincodeName='';
var channelName='';
var fcn='';
var args='';
var args1='';
var args2='';
var args3='';
var args4='';

var topic = 'fabric'
client.on('connect', () => {

    client.subscribe(topic)

})
client.on('message', (topic, message) => {
    console.log("client.on message" +message);
    console.log();
    var data = JSON.parse(message)

    org_name = data.orgname;
    username = data.username;
    chaincodeName = data.chaincodeName;
    channelName = data.channelName;
    fcn = data.fcn;
    args = data.args[0];
    args1 = data.args[1];
    args2= data.args[2];
    args3 = data.args[3];
    console.log("client.on message bitti");
    console.log("fcn: " +fcn);
    console.log("args: " +args);
    console.log("args1: " +args1);
    console.log();    console.log();
    console.log();

    

})

const invokeTransaction = async (channelName, chaincodeName, username, org_name,fcn_,args_) => {
    try {
        logger.debug(util.format('\n============ invoke transaction on channel %s ============\n', channelName));

        const ccp = await helper.getCCP(org_name)
        const walletPath = await helper.getWalletPath(org_name) 
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        let identity = await wallet.get(username);
        if (!identity) {
            console.log(` ${username} adlı kullanıcının kaydı bulunamadı, lütfen kayıt edin.`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            return;
        }
        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: false },
            eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            }
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);
        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        let result
        let message;
        console.log("invokeTransaction     fcn: "+fcn_);
        console.log("args  ",args_);
        if (fcn_=== "createData") {
            console.log("invoke.js createData");
            result = await contract.submitTransaction(fcn_,args_);
            console.log("result: "+result);
            message = `Successfully createData key: ${args_[0]}`
        }
        else {
            return `Invocation require either as function but got ${fcn}`
        }
        await gateway.disconnect();

        let response = {
            message: message
        }

        return response;


    } catch (error) {

        console.log(`Getting error: ${error}`)
        return error.message

    }

}


exports.invokeTransaction = invokeTransaction;
