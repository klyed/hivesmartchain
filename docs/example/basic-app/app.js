const fs = require('fs')
const hsc = require('@monax/hsc')
const express = require('express')
const bodyParser = require('body-parser')

// HiveSmartChain address
let chainURL = '127.0.0.1:10997'
const abiFile = 'bin/simplestorage.bin'
const deployFile = 'deploy.output.json'
const accountFile = 'account.json'

// Port to run example on locally
const exampleAppPort = 3000

function slurp (file) {
  return JSON.parse(fs.readFileSync(file, 'utf8'))
}

// Grab the account file that is expected to have 'Address' field
let account = slurp(accountFile)
// Connect to running hsc chain using the account address to identify our input account and return values as an object
// using named returns where provided
let chain = hsc.createInstance(chainURL, account.Address, {objectReturn: true})
// The ABI file produced by the solidity compiler (through hsc deploy) that acts as a manifest for our deployed contract
let abi = slurp(abiFile).Abi
// The deployment receipt written to disk by hsc deploy that contains the deployed address of the contract amongst other things
let deploy = slurp(deployFile)
// The contract we will call
let contractAddress = deploy.simplestorage
// A Javascript object that wraps our simplestorage contract and will handle translating Javascript calls to EVM invocations
let store = chain.contracts.new(abi, null, contractAddress)

// For this example we use a simple router based on expressjs
const app = express()
// Apparently this needs to be its own module...
app.use(bodyParser.json())

// Some helpers for parsing/validating input
let asInteger = value => new Promise((resolve, reject) =>
    (i => isNaN(i) ? reject(`${value} is ${typeof value} not integer`) : resolve(i))(parseInt(value)))

let param = (obj, prop) => new Promise((resolve, reject) =>
    prop in obj ? resolve(obj[prop]) : reject(`expected key '${prop}' in ${JSON.stringify(obj)}`))

let handleError = err => {
  console.log(err)
  return err.toString()
}

// We define some method endpoints
// Get the value from the contract by calling the Solidity 'get' method
app.get('/', (req, res) => store.get()
    .then(ret => res.send(ret.values))
    .catch(err => res.send(handleError(err))))

// Sets the value by accepting a value in HTTP POST data and calling the Solidity 'set' method
app.post('/', (req, res) => param(req.body, 'value')
    .then(value => asInteger(value))
    .then(value => store.set(value).then(() => value))
    .then(value => res.send({value: value, success: true}))
    .catch(err => res.send(handleError(err))))

// Sets a value by HTTP POSTing to the value you expect to be stored encoded in the URL - so that the value can be
// updated atomically
app.post('/:test', (req, res) => param(req.body, 'value')
    .then(value => Promise.all([asInteger(req.params.test), asInteger(value)]))
    .then(([test, value]) => store.testAndSet(test, value))
    .then(ret => res.send(ret.values))
    .catch(err => res.send(handleError(err))))

// Send a little value to an account which has the effect of creating that account if it does not exist
app.post('/send/:recipient', (req, res) => param(req.body, 'amount')
    .then(amount =>
        chain.transact.SendTxSync(
          {
            Inputs: [{
              Address: Buffer.from(account.Address, 'hex'),
              Amount: amount
            }],
            Outputs: [{
              Address: Buffer.from(req.params.recipient, 'hex'),
              Amount: amount
            }]
          }))
    .then(txe => res.send({txHash: txe.TxHash.toString('hex'), success: true}))
    .catch(err => res.send(handleError(err))))

const url = `http://127.0.0.1:${exampleAppPort}`

// Listen for requests
app.listen(exampleAppPort, () => console.log(`Example app listening on ${url}...

You may wish to try the following: 
# Inspect current stored value
  $ curl ${url}
  
# Set the value to 2000
  $ curl -d '{"value": 2000}' -H "Content-Type: application/json" -X POST ${url}
  
# Set the value via a testAndSet operation
  $ curl -d '{"value": 30}' -H "Content-Type: application/json" -X POST ${url}/2000
  
# Attempt the same testAndSet which now fails since the value stored is no longer '2000'
  $ curl -d '{"value": 30}' -H "Content-Type: application/json" -X POST ${url}/2000
  $ curl ${url}
  `))
