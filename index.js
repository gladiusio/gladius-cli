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
let userData;
let envData = JSON.parse(fs.readFileSync(appDir + "/envData.json")); //for now have this file already made

// Set up prompt
prompt.message = colors.blue("[Gladius-Node]");
prompt.delimiter = " ";

prompt.start();


// Check on userData.json file
if (!fs.existsSync(appDir + "/userData.json")) {
  reset();
}

/**
* Commands for the user to call
* toCall is the function associated with the command
*/
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
  "--help": {
    description: "Show this menu",
    toCall: function() {
      help(options);
    }
  }
}

// Prompt the user for information about themselves
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
      },
      pvKey: {
        description: "Wallet private key location (my/key/here):",
        required: true,
        message: "Not a valid key or key location",
        conform: function(value) {
          return (value == "test"); //replace this to check if valid file path
        }
      },
      pgpKey: {
        description: "PGP key location (my/key/here):",
        required: true,
        message: "Not a valid key or key location",
        conform: function(value) {
          return (value == "test"); //replace this to check if valid file path
        }
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    console.log(colors.blue(
      "\nIf the above information isn't correct, run init again. If it is, you can run gladius-node join-pool."
    ));
    userData = {
      email: result.email,
      name: result.name,
      bio: result.bio,
      pvKey: result.pvKey,
      pgpKey: result.pgpKey,
      initialized: true // Set our initialized flag
    };
    createNode(); //generate a node contract
  });

}

/**
* Create a node based on the onboarding information
*/
function createNode() {
  axios.post(daemonAddress + "/api/node/create", {
    //no data required, data is set AFTER you create the initial node contract
  })
  .then(function(res){
    console.log(res);
    userData.nodeAddress = res.data.address
    writeInitInfo(userData); // Write it to a file
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
    console.log(res);
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
      .then(function(response) {
        console.log(response);
      })
      .catch(function(error) {
        console.log(colors.red(
          "Woah an error! Make sure your daemon is running and can be connected to"
        ));
        console.log(error);
      });
  } else {
    console.log(colors.red(
      "Error: You need to initialize your node first. Run gladius-node init to do this."
    ))
  }
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

/**
* Start accepting connections
*/
function start() {
  axios.put(daemonAddress + "/api/status/", {
      status: true
    })
    .then(function(response) {
      console.log(response);
    })
    .catch(function(error) {
      console.log(colors.red(
        "Woah an error! Make sure your daemon is running and can be connected to"
      ));
      console.log(error);
    });
}

/**
* Stop accepting connections
*/
function stop() {
  axios.put(daemonAddress + "/api/status/", {
      status: false
    })
    .then(function(response) {
      console.log(response);
    })
    .catch(function(error) {
      console.log(colors.red(
        "Woah an error! Make sure your daemon is running and can be connected to"
      ));
      console.log(error);
    });
}

// Get the current status of the node daemons
function status() {
  axios.get(daemonAddress + "/api/status/")
    .then(function(response) {
      console.log(response);
    })
    .catch(function(error) {
      console.log(colors.red(
        "Woah an error! Make sure your daemon is running and can be connected to"
      ));
      console.log(error);
    });
}

/**
* List pools
*/
function listPools() {
  axios.get(daemonAddress + "/api/pools/")
    .then(function(response) {
      console.log(response);
    })
    .catch(function(error) {
      console.log(colors.red(
        "Woah an error! Make sure your daemon is running and can be connected to"
      ));
      console.log(error);
    });
}

// Check the status of current applications
function checkJoin() {
  axios.get(daemonAddress + "/api/pools/check/beta")
    .then(function(response) {
      console.log(response);
    })
    .catch(function(error) {
      console.log(colors.red(
        "Woah an error! Make sure your daemon is running and can be connected to"
      ));
      console.log(error);
    });
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

function writeInitInfo(info) {
  var json = JSON.stringify(info);
  fs.writeFileSync(appDir + "/userData.json", json);
}


function getInitInfo() {
  return JSON.parse(fs.readFileSync(appDir + "/userData.json"));
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
  "--help": {
    description: "Show this menu",
    toCall: function() {
      help(options);
    }
  }
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
