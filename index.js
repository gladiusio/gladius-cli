#!/usr/bin/env node

/*
Gladius node CLI. Install with "npm install gladius-cli"
then run: "gladius-node init" and follow the steps. After you have
successfully completed this you can run: "gladius-node start"

Requires the Gladius daemon to be installed and running to opperate.
*/

//We should have color (yellow probably) that indicates that the result of executing this action will incur gas costs

var prompt = require("prompt");
var colors = require("colors/safe");
var fs = require("fs");
var path = require("path");
var axios = require("axios");
var config = require("./config.js") // Load our config file
var Spinner = require('cli-spinner').Spinner;

let appDir = path.dirname(require.main.filename); // Get where this file is
let daemonAddress = config.network.controlDaemonAddress + ":" + config.network.controlDaemonPort;
let pgpKey = fs.readFileSync("./keys/pgpKey.txt","utf8") //pgp key
let pvtKey = fs.readFileSync("./keys/pvtKey.txt","utf8") //ETH wallet private key
let settings = require("./settings.json")
let passphrase;

// Set up prompt
prompt.message = colors.blue("[Gladius-Node]");
prompt.delimiter = " ";

prompt.start();

if(!fs.existsSync("./nodeFile.json")) {
  reset();
}

let nodeFile = require("./nodeFile.json")

/**
* Create options for the user where description is the description of the
* argument and toCall is a function.
*/

var options = {
  "init": {
    description: "Gathers information about the user as well as configuration data.",
    toCall: init
  },
  "create": {
    description: "Create a node",
    toCall: create
  },
  "apply": {
    description: "Apply to a pool",
    toCall: apply // Eventually replace with arbitrary pool upon launch
  },
  "check": {
    description: "Status of your application to a pool",
    toCall: getPoolStatus
  },
  "status": {
    description: "Get's the current status of the node daemons",
    toCall: status
  },
  "dirs": {
    description: "Returns the location of the config.js file",
    toCall: locations
  },
  "reset": {
    description: "Resets init file (for testing or problem installations)",
    toCall: function() {
      reset();
      console.log(colors.blue("nodeFile.json has been reset"));
    }
  },
  "settings": {
    description: "Show settings",
    toCall: getSettings
  },
  "set-node": {
    description: "Push information to the Node",
    toCall: setNodeData
  },
  "test": {
    description: "Test random functions",
    toCall: test
  },
  "--help": {
    description: "Show this menu",
    toCall: function() {
      help(options);
    }
  }
}

/**
* Prompt the user for information about themselves
* just writes to the nodeFile.json TODO - pull in ip address
*/
function init() {
  // checkKeys() //make sure the keys are there before doing this

  // Create a schema for the paremeters to be asked
  let schema = {
    properties: {
      email: {
        description: "What's your email? (So we can contact you about the beta):",
        pattern: /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
        message: "Not a valid email",
        required: true
      },
      name: {
        description: "What's your first name?",
        required: true
      },
      bio: {
        description: "Short bio about why you're interested in Gladius:",
        required: true
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    nodeFile.userData = {
      email: result.email,
      name: result.name,
      bio: result.bio,
      initialized: true // Set our initialized flag
    };
    writeToFile("nodeFile", nodeFile); // Write it to a file
    console.log(colors.green("User profile created! You may now create a node with gladius-node create"));
  });
}

/**
* Create a new Node smart contract (no data set)
*/
function create() {
  postSettings(function() {
    axios.post(daemonAddress + "/api/node/create", {
      //no data required, data is set AFTER you create the initial node contract
    })
    .then(function(res){
      console.log(colors.blue("[Gladius-Node] ") + "Creating Node contract, please wait for tx to complete (this might take a couple of minutes) ");

      creationStatus(res.data.txHash, function(err, res) {
        if(res == colors.green("[Success]")) {
          console.log();
          console.log(colors.blue("[Gladius-Node] ") + "Setting Node data, please wait for tx to complete (this might take a couple of minutes) ");

          getNodeAddress(function() {
            setNodeData(function(tx) {
              creationStatus(tx, function(err, res) {
                if(res == colors.green("[Success]")) {
                  console.log();
                  console.log(colors.green("[Gladius-Node] " + "Node successfully created and ready to use"));
                  console.log(colors.blue("[Gladius-Node] ") + "Use " + colors.blue("gladius-node apply") + " to apply to a pool");
                }
                else {
                  console.log(colors.red("[Gladius-Node] ") + "There was a problem accessing your Node Contract");
                }
              })
            })
          })

        }
        else{
          console.log(colors.red("[Gladius-Node] ") + "There was a problem creating your Node Contract");
        }

      })
    })
    .catch(function(err){
      console.log(err);
      console.log(colors.red("There was a problem creating a node"));
    })
  })
}

/**
* See status of the node daemon
*/
function apply() {
  let schema = {
    properties: {
      poolAddress: {
        description: "Please enter the address of the pool you want to join: ",
        required: true
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    axios.post(daemonAddress + "/api/node/" + nodeFile.address + "/apply/" + result.poolAddress, nodeFile.userData)
    .then(function(res) {
      creationStatus(res.data.tx, function(){
        console.log(colors.green("[Gladius-Node] Application sent to Pool! Use " + colors.blue("gladius-node") + " check to check your application status"));
      })
    })
    .catch(function(err) {
      console.log(err.data);
    })
  });
}

/**
* Start accepting connections, right now just posts the settings to start the server
*/
function postSettings(callback) {
  let schema = {
    properties: {
      passphrase: {
        description: "Please enter the passphrase for your PGP private key:",
        required: true,
        hidden: true
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    axios.post(daemonAddress + "/api/settings/start", {
      "provider": settings.provider,
      "privateKey": pvtKey.toString().replace(/\r?\n|\r/g,""),
      "pgpKey": pgpKey.toString().replace(/\r?\n|\r/g,"\n"),
      "marketAddress": settings.marketAddress,
      "nodeFactoryAddress": settings.nodeFactoryAddress,
      "passphrase": result.passphrase
    })
      .then(function() {
        callback()
      })
      .catch(function(err) {
        console.log(err);
        console.log(colors.red("There was a problem posting your settings"));
      });
  });
}

/**
* Check status of the BC control daemon
*/
function status() {
  axios.get(daemonAddress + "/api/status/")
    .then(function(res) {
      console.log(colors.blue("Server is running!"));
    })
    .catch(function(err) {
      console.log(colors.red("Server is down"));
    });
}

/**
* Check on status of creating a transaction and halt i/o until done
*/
function creationStatus(tx, callback) {
  let status = 0;

  axios.get(daemonAddress + "/api/status/tx/" + tx)
  .then(function(res) {
    if(res.data.receipt) {
      if(res.data.receipt.status == "0x0") {
        status = 0 //Status: Failed
      }
      else if (res.data.receipt.status = "0x1") {
        status = 1 //Status: Success
      }
      else {
        status = 3 //Status: Unknown
      }
    }
    else {
      status = 2 //Status: Pending
    }

    if(status == 2) {
      creationStatus(tx, callback)
    }

    let _status;

    switch(status) {
      case 0:
        _status = colors.red("[Failed]")
        process.stdout.write(colors.red("[Gladius Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      case 1:
        _status = colors.green("[Success]")
        process.stdout.write(colors.green("[Gladius Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      case 2:
        _status = colors.yellow("[Pending]"+"\r")
        process.stdout.write(colors.yellow("[Gladius Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      default:
        _status = "[Unknown]"
        process.stdout.write(colors.blue("[Gladius Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
    }

    if (status == 1) {
      callback(null, _status)
    }
  })
  .catch(function(err) {
    console.log(err);
  })
}

/**
* Set the data for the node based on onboarding info
*/
function setNodeData(callback) {
  axios.post(daemonAddress + "/api/node/" + nodeFile.address + "/data/", nodeFile.userData)
  .then(function(res){
    callback(res.data.txHash)
  })
  .catch(function(err){
    console.log();
    console.log(colors.red("[Gladius-Node]") + " Error setting Node data");
    console.log(err);
  })
}

/**
* Get the node address by reverse looking up the owner
*/
function getNodeAddress(callback) {
  axios.get(daemonAddress + "/api/node")
  .then(function(res) {
    nodeFile.address = res.data.address
    writeToFile("nodeFile", nodeFile)
    callback()
  })
  .catch(function(err) {
    console.log(colors.red("[Gladius-Node] ") + "Couldn't get node address");
  })
}

/** WIP - need to add a stop/kill endpoint
* Stop accepting connections
*/
function stop() {
  axios.put(daemonAddress + "/api/status/", {
      status: false
    })
    .then(function(res) {
      console.log(res);
    })
    .catch(function(err) {
      console.log(colors.red(
        "Woah an err! Make sure your daemon is running and can be connected to"
      ));
      console.log(err);
    });
}

function locations() {
  console.log("Keys: " + appDir + "/keys");
  console.log("UserData: " + appDir + "/nodeFile.json");
  console.log("Settings: " + appDir + "/settings.json");
}
/**
* Get the data for the node env. PVT key, PGP, etc...
*/
function getSettings() {
  axios.get(daemonAddress + "/api/settings")
  .then(function(res) {
    console.log(res.data);
  })
  .catch(function(err) {
    console.log(err);
  })
}

/**
* Key info
*/
function getKeys() {
  console.log("private and pgp keys are located in ./keys");
}

function checkKeys() {
  if (pvtKey.toString() == "INSERT PRIVATE KEY TO THROWAWAY WALLET HERE" || pgpKey.toString() == "INSERT PGP PRIVATE KEY HERE") {
    console.log(colors.red("[Gladius-Node] ") + "You have not pasted your keys in the ./keys folder, do this before proceeding");
    process.exit(1);
  }
}

/**
* Application status for this node's pools
*/
function getPoolStatus() {
  let schema = {
    properties: {
      poolAddress: {
        description: "Please enter the address of the pool you want to check on: ",
        required: true
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    axios.get(daemonAddress + "/api/node/" + nodeFile.address + "/status/" + result.poolAddress)
    .then(function(res) {
      let poolStatus;

      switch(res.data.status) {
        case "Rejected":
          poolStatus = (colors.red("[Application Status: Rejected]"));
          break;
        case "Pending":
          poolStatus = (colors.yellow("[Application Status: Pending]"));
          break;
        case "Approved":
          poolStatus = (colors.green("[Application Status: Approved]"));
          break;
      }
      console.log();
      console.log("Pool: " + result.poolAddress + "\t" + poolStatus);
    })
    .catch(function(err){
      console.log(err);
    })
  });
}

/** WIP - Missing Endpoint
* List pools
*/
function listPools() {
  axios.get(daemonAddress + "/api/pools/")
    .then(function(res) {
      console.log(res.data);
    })
    .catch(function(err) {
      console.log(colors.red(
        "Woah an err! Make sure your daemon is running and can be connected to"
      ));
      console.log(err);
    });
}

/*
* For testing rando functions
*/
function test(reeeee) {
  console.log(nodeFile.address);
}

/**
* Builds and prints help menu
* @param options - commands that the users run
*/
function help(options) {
  console.log(colors.blue(
      "\n--------------Available arguments-------------- \n") +
    Object.keys(options).map(
      function(key) {
        return ("\n\n" + colors.blue(key) + ": " + options[key].description);
      }
    ).join(""));
}

/*
* Helper function for writing to files
*/
function writeToFile(name, data) {
  var json = JSON.stringify(data);
  fs.writeFileSync(appDir + "/"+name+".json", json);
}

/*
* Reset the nodeFile
*/
function reset() {
  let data = {
    "userData":
      {
        "email":"",
        "name":"",
        "bio":"",
        "initialized":false
      },
    "address":""
  }
  fs.writeFileSync("./nodeFile.json", JSON.stringify(data, null, 2))
}

// Get the argument that the user provided
var argument = process.argv[2];

// Run the CLI
if (true) {
  if (argument in options) {
    options[argument].toCall();
  } else {
    console.log(colors.red(
      "Unknown (or no) argument, please use --help to see available arguments"
    ));
  }
} else {
  console.log(colors.red(
    "Cannot connect to the Gladius daemon. See setup instructions here: https://github.com/gladiusio"
  ));
}
