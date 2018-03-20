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

let appDir = path.dirname(require.main.filename); // Get where this file is
let daemonAddress = config.controlDaemonAddress + ":" + config.controlDaemonPort;
let pgpKey = fs.readFileSync("./keys/pgpKey.txt","utf8") //pgp key
let pvtKey = fs.readFileSync("./keys/pvtKey.txt","utf8") //ETH wallet private key
let passphrase;

let userData; //data object w user info
let envData; //data object w env info

// Set up prompt
prompt.message = colors.blue("[Gladius-Node]");
prompt.delimiter = " ";

prompt.start();


// Check on userData.json file
if (!fs.existsSync(appDir + "/userData.json")) {
  reset();
}

// Create options for the user where description is the description of the
// argument and toCall is a function.
var options = {
  "init": {
    description: "Gathers information about the user as well as configuration data.",
    toCall: init
  },
  "start": {
    description: "Start the Gladius node, and inform the pool of this.",
    toCall: start
  },
  "stop": {
    description: "Stop the Gladius node, and inform the pool of this.",
    toCall: stop
  },
  "status": {
    description: "Get's the current status of the node daemons",
    toCall: status
  },
  "list-pools": {
    description: "List all available pools from the marketplace",
    toCall: listPools
  },
  "join-pool": {
    description: "Join the beta pool (will have arguments in future to specify pool to join)",
    toCall: joinPool // Eventually replace with arbitrary pool upon launch
  },
  "check-join": {
    description: "Check the status of your applications.",
    toCall: checkJoin
  },
  "config-location": {
    description: "Returns the location of the config.js file",
    toCall: function() {
      console.log(appDir + "/config.js")
    }
  },
  "reset-init": {
    description: "Resets init file (for testing or problem installations)",
    toCall: function() {
      reset();
      console.log(colors.blue("Reset userData.json"));
    }
  },
  "settings": {
    description: "Show settings",
    toCall: getSettings
  },
  "keys": {
    description: "Location of your key files",
    toCall: getKeys
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
* just writes to the userData.json
*/
function init() {
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
    console.log(colors.blue("\nRun gladius-node --help for more actions"));
    userData = {
      email: result.email,
      name: result.name,
      bio: result.bio,
      initialized: true // Set our initialized flag
    };
  });

}

/**
* Create a new Node smart contract (no data set)
*/
function createNode() {
  axios.post(daemonAddress + "/api/node/create", {
    //no data required, data is set AFTER you create the initial node contract
  })
  .then(function(res){
    console.log(res.data);
    userData.nodeAddress = res.data.address
    writeToFile("userData",userData); // Write it to a file
  })
  .catch(function(err){
    console.log(err);

  })
}

/**
* Set the data for the node based on onboarding info
*/
function setNodeData() {
  axios.post(daemonAddress + "/api/node/" + userData.nodeAddress + "/data", {
    //no data required, data is set AFTER you create the initial node contract
  })
  .then(function(res){
    console.log(res.data);
  })
  .catch(function(err){
    console.log(err);
  })
}

/** WIP - should be changed to just join pool
* See status of the node daemon
*/
function joinPool() {
  var initInfo = getInitInfo(); // Grab the information from initialization

  if (initInfo.initialized) {
    axios.post(daemonAddress + "/api/pools/beta", {
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
  } else {
    console.log(colors.red(
      "err: You need to initialize your node first. Run gladius-node init to do this."
    ))
  }
}

/**
* Start accepting connections, right now just posts the settings to start the server
*/
function start() {
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
      "provider": "http://127.0.0.1:9545",
      "privateKey": pvtKey.toString(),
      "pgpKey": pgpKey.toString().replace(/\r?\n|\r/g,"\\n"),
      "passphrase": result.passphrase,
      "marketAddress": "0x345ca3e014aaf5dca488057592ee47305d9b3e10",
      "nodeFactoryAddress": "0xb9a219631aed55ebc3d998f17c3840b7ec39c0cc"
    })
      .then(function(res) {
        console.log(colors.blue("Server is running!"));
      })
      .catch(function(err) {
        console.log(colors.red("Have you set your keys in the ./keys folder yet?"));
        console.log(err);
      });
  });
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

// Get the current status of the node daemons
function status() {
  axios.get(daemonAddress + "/api/status/")
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

function getSettings() {
  axios.get(daemonAddress + "/api/settings")
  .then(function(res) {
    console.log(res);
  })
  .catch(function(err) {
    console.log(err);
  })
}

function getKeys() {
  console.log("private and pgp keys are located in ./keys");
}


/**
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

// Check the status of current applications
function checkJoin() {
  axios.get(daemonAddress + "/api/pools/check/beta")
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

// Reset the userData.json file for testing or for problem installations.
function reset() {
  var json = JSON.stringify({
    email: "",
    name: "",
    bio: "",
    key: "",
    initialized: false
  });
  fs.writeFileSync(appDir + "/userData.json", json);
}

function writeToFile(name, data) {
  var json = JSON.stringify(data);
  fs.writeFileSync(appDir + "/"+name+".json", json);
}

function getInitInfo() {
  return JSON.parse(fs.readFileSync(appDir + "/userData.json"));
}

// Check the up status of the daemon
function checkDaemon() {
  // TODO: Check the status of the daemon from the status method

  return true;
}

// Get the argument that the user provided
var argument = process.argv[2];


// Run the CLI
if (checkDaemon()) {
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
