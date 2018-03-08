#!/usr/bin/env node

/*
Gladius node CLI. Install with "npm install gladius-cli"
then run: "gladius-node init" and follow the steps. After you have
successfully completed this you can run: "gladius-node start"

Requires the Gladius daemon to be installed and running to opperate.
*/

var prompt = require("prompt");
var colors = require("colors/safe");
var fs = require("fs");
var path = require("path");

var config = require("./config.js") // Load our config file

var appDir = path.dirname(require.main.filename); // Get where this file is


console.log(appDir);

prompt.message = colors.blue("[Gladius-Node]");
prompt.delimiter = " ";

prompt.start();


// The functions below are called by the various options provided
/******************************************************************************/

// Prompt the user for information about themselves
function init() {
  // Create a schema for the paremeters to be asked
  let schema = {
    properties: {
      email: {
        description: "What's your email?",
        pattern: /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
        message: "Not a valid email",
        required: true
      },
      name: {
        description: "What's your first name?",
        required: true
      },
      privateKeyLocation: {
        description: "Private key location (my/key/here):",
        required: true
      }
    }
  };

  // Prompt and forward data
  prompt.get(schema, function(err, result) {
    // TODO: Save the data to a config.json file
  });
}

// Build the help menu from the options
function help(options) {
  console.log("\n------Available arguments------ \n" + Object.keys(options).map(function(key) {
    return ("\n\n" + colors.blue(key) + ": " + options[key].description);
  }).join(""));
}

// Inform the control daemon that the node is ready
function start() {

}

// Inform the control daemon that the node is no longer ready
function stop() {

}

// Get the current status of the node daemons
function status() {

}

// Join the beta pool
function joinBetaPool() {

}

// List all available pools
function listPools() {

}

// Check the status of current applications
function checkJoin() {

}

/******************************************************************************/

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
    toCall: joinBetaPool // Eventually replace with arbitrary pool upon launch
  },
  "check-join": {
    description: "Check the status of your applications.",
    toCall: checkJoin
  },
  "config-location": {
    description: "Returns the location of the config.js file",
    toCall: function() {
      console.log(appDir + "config.js")
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
  // Check the status of the daemon from the status method

  return true;
}

// Get the argument that the user provided
var argument = process.argv[2];

// Run the CLI
if (checkDaemon()) {
  if (argument in options) {
    options[argument].toCall();
  } else {
    console.log("Unknown (or no) argument, please use --help to see available arguments")
  }
} else {
  console.log(colors.red("Cannot connect to the Gladius daemon. See setup instructions here: https://github.com/gladiusio"));;
}
