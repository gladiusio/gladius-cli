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
var os = require("os");
var {promisify} = require('util'); //<-- Require promisify
var getIP = promisify(require('external-ip')()); // <-- And then wrap the library
var rpc = require('node-json-rpc');
var kbpgp = require('kbpgp');
let settings = require("./settings.json")

let appDir = path.dirname(require.main.filename); // Get where this file is
let daemonAddress = config.network.controlDaemonAddress + ":" + config.network.controlDaemonPort;
let pgpKey //pgp key
let pvtKey //ETH wallet private key

// Set up prompt
prompt.message = colors.blue("[Gladius-Node]");
prompt.delimiter = " ";

prompt.start();

if(!fs.existsSync(appDir + "/nodeFile.json")) {
  reset();
}

let nodeFile = require("./nodeFile.json")

// RPC Stuff ==
var rpcOptions = {
  // int port of rpc server, default 5080 for http or 5433 for https
  port: 5000,
  // string domain name or ip of rpc server, default '127.0.0.1'
  host: 'localhost',
  // string with default path, default '/'
  path: '/rpc',
  // boolean false to turn rpc checks off, default true
  strict: true
};
var rpcClient = new rpc.Client(rpcOptions);
// RPC Stuff ==

/**
* Create options for the user where description is the description of the
* argument and toCall is a function.
*/
var options = {
  "init": {
    description: "Gathers information about the user as well as configuration data.",
    toCall: function(){reset(init)}
  },
  "create": {
    description: "Create a node",
    toCall: function(){postSettings(create)}
  },
  "apply": {
    description: "Apply to a pool",
    toCall: function(){postSettings(apply)} // Eventually replace with arbitrary pool upon launch
  },
  "check": {
    description: "Status of your application to a pool",
    toCall: function(){postSettings(checkPoolStatus)}
  },
  "status": {
    description: "Get's the current status of the node daemons",
    toCall: function(){postSettings(status)}
  },
  "start": {
    description: "Start accepting connections and acting as an edge node",
    toCall: function(){postSettings(startNetworking)}
  },
  "stop": {
    description: "Stop accepting connections",
    toCall: function(){postSettings(stopNetworking)}
  },
  "gen-keys": {
    description: "Generate new PGP keys",
    toCall: function() {genPGPKey(function(){console.log(colors.green("[Gladius-Node] ") + "New PGP keys generated");})}
  },
  "settings": {
    description: "Show settings",
    toCall: function(){postSettings(getSettings)}
  },
  "dirs": {
    description: "Returns the location of the config.js file",
    toCall: function(){postSettings(locations)}
  },
  "reset": {
    description: "Resets init file (for testing or problem installations)",
    toCall: function() {
      reset();
      console.log(colors.blue("nodeFile.json has been reset"));
    }
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
* just writes to the nodeFile.json
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
      },
      pvtKey: {
        description: "Please make a new ETH wallet and paste your private key (include 0x at the beginning): ",
        required: true,
        hidden: true,
        pattern: /0[xX][0-9a-zA-Z]+/,
        message: "Please include \'0x\' at the beginning of your private key"
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

    fs.writeFileSync(appDir+"/keys/ethPvtKey.txt", result.pvtKey)

    genPGPKey(function() {
      getIP()
      .then(function(ip) {
        nodeFile.userData.ip = ip;
        writeToFile("nodeFile", nodeFile); // Write it to a file
        // console.log(colors.green("[Gladius-Node]") + " User profile created! Please paste your ETH and PGP private keys in the " + colors.blue("./keys") + " directory");
        console.log(colors.green("[Gladius-Node]") + " User profile created! You may create a node with " + colors.blue("gladius-node create"));
        console.log(colors.green("[Gladius-Node]") + " If you'd like to change your information run " + colors.blue("gladius-node init") + " again");
      })
      .catch(function(error){
        console.error(error);
      });
    })
  });
}

/**
* Create a new Node smart contract (no data set)
*/
function create() {
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
        console.log();
        console.log(colors.green("[Gladius-Node] Application sent to Pool!"));
        console.log(colors.green("[Gladius-Node]") + " Use " + colors.blue("gladius-node check") + " to check your application status");
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
  checkKeys(function() {
    let schema = {
      properties: {
        passphrase: {
          description: "Please enter the passphrase for your PGP private key:",
          required: true,
          hidden: true
        }
      }
    };

    pvtKey = fs.readFileSync(appDir + "./keys/ethPvtKey.txt","utf8")
    pgpKey = fs.readFileSync(appDir + "./keys/pgpPvtKey.txt","utf8")
    // console.log(pvtKey);
    // console.log(pgpKey);

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
          // console.log(err);
          console.log(colors.red("There was a problem posting your settings"));
        });
    });
  })
}

/**
* Check status of the BC control daemon
*/
function status(callback) {
  axios.get(daemonAddress + "/api/status/")
    .then(function(res) {
      if(callback == null) {
        console.log(colors.green("[Gladius-Node]") + " gladius-control-daemon server is running!");
      }
      else {
        callback()
      }
    })
    .catch(function(err) {
      console.log(err);
      console.log(colors.red("[Gladius-Node]") + " gladius-control-daemon server is down at "+ daemonAddress +" Run " + colors.blue("node index.js") + " in the gladius-control-daemon directory");
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
        process.stdout.write(colors.red("[Gladius-Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      case 1:
        _status = colors.green("[Success]")
        process.stdout.write(colors.green("[Gladius-Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      case 2:
        _status = colors.yellow("[Pending]"+"\r")
        process.stdout.write(colors.yellow("[Gladius-Node]")+ " Transaction: " + tx + "\t" + _status)
        break;
      default:
        _status = "[Unknown]"
        process.stdout.write(colors.blue("[Gladius-Node]")+ " Transaction: " + tx + "\t" + _status)
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

/**
* Where settings are stored
*/
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

/**
* Check if user keys are there, if not end proccess
*/
function checkKeys(callback) {
  if (!fs.existsSync(appDir+"/keys/ethPvtKey.txt") || !fs.existsSync(appDir+"/keys/pgpPubKey.txt") || !fs.existsSync(appDir+"/keys/pgpPvtKey.txt")) {
    console.log(colors.red("[Gladius-Node] ") + "You do not have any key files. Run " + colors.blue("gladius-node init") + " to set up your information");
    process.exit(1);
  }
  else {
    callback()
  }
}

/**
* Application status for this node's pools
*/
function checkPoolStatus() {
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

      if(res != null) {
        let poolStatus;
        let message;

        switch(res.data.code) {
          case 1:
            poolStatus = colors.green("[Gladius-Node] ") + "Pool: " + result.poolAddress + "\t" + colors.green("[Application Status: Green]")
            message = colors.green("[Gladius-Node] ") + "You've been accepted! Use " + colors.blue("gladius-node start") + " to start accepting connections!"
            break;
          case 2:
            poolStatus = colors.red("[Gladius-Node] ") + "Pool: " + result.poolAddress + "\t" + colors.red("[Application Status: Rejected]");
            message = colors.red("[Gladius-Node] ") + "Consider applying to a different pool"
            break;
          case 3:
            poolStatus = colors.yellow("[Gladius-Node] ") + "Pool: " + result.poolAddress + "\t" + colors.yellow("[Application Status: Pending]")
            message = colors.yellow("[Gladius-Node] ") + "Wait until the pool manager accepts your application in order to become an edge node"
            break;
          default:
          poolStatus = colors.magenta("[Gladius-Node] ") + "Pool: " + result.poolAddress + "\t" + colors.red("[Application Status: Not Sent]");
          message = colors.magenta("[Gladius-Node] ") + "You have not sent an application to this pool. Use " + colors.blue("gladius-node apply") + " to apply"
          break;
        }
        console.log(poolStatus);
        console.log(message);
      }
      else {
        console.log(colors.red("[Gladius-Node] ") + "You've entered the wrong pool address");
      }

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
* Generate pgp key pairs
*/
function genPGPKey(callback) {
  let schema = {
    properties: {
      passphrase: {
        description: "Please enter a passphrase for your new PGP keys: ",
        required: true,
        hidden: true
      }
    }
  };

  // Prompt and store the data
  prompt.get(schema, function(err, result) {
    var F = kbpgp["const"].openpgp;

    var opts = {
      userid: "User McTester (Born 1979) <user@example.com>",
      primary: {
        nbits: 1024,
        flags: F.certify_keys | F.sign_data | F.auth | F.encrypt_comm,
        expire_in: 0  // never expire
      },
      subkeys: []
    };

    kbpgp.KeyManager.generate(opts, function(err, alice) {
      if (!err) {
        // sign alice's subkeys
        alice.sign({}, function(err) {
          // console.log(alice);
          // export demo; dump the private with a passphrase
          alice.export_pgp_private ({
            passphrase: result.passphrase
          }, function(err, pgp_private) {
            fs.writeFileSync(appDir+"/keys/pgpPvtKey.txt", pgp_private)
            callback()
          });
          alice.export_pgp_public({}, function(err, pgp_public) {
            fs.writeFileSync(appDir+"/keys/pgpPubKey.txt", pgp_public)
          });
        });
      }
    });
  });
}


/*
* For testing rando functions
*/
function test() {
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
    var F = kbpgp["const"].openpgp;

    var opts = {
      userid: "User McTester (Born 1979) <user@example.com>",
      primary: {
        nbits: 1024,
        flags: F.certify_keys | F.sign_data | F.auth | F.encrypt_comm | F.encrypt_storage,
        expire_in: 0  // never expire
      },
      subkeys: [
        {
          nbits: 1024,
          flags: F.sign_data,
          expire_in: 86400 * 365 * 8 // 8 years
        }, {
          nbits: 1024,
          flags: F.encrypt_comm | F.encrypt_storage,
          expire_in: 86400 * 365 * 8
        }
      ]
    };

    kbpgp.KeyManager.generate(opts, function(err, alice) {
      if (!err) {
        // sign alice's subkeys
        alice.sign({}, function(err) {
          // console.log(alice);
          // export demo; dump the private with a passphrase
          alice.export_pgp_private ({
            passphrase: result.passphrase
          }, function(err, pgp_private) {
            fs.writeFileSync(appDir+"/keys/pgpPvtKey.txt", pgp_private)
            console.log("Done");
          });
          alice.export_pgp_public({}, function(err, pgp_public) {
            fs.writeFileSync(appDir+"/keys/pgpPubKey.txt", pgp_public)
          });
        });
      }
    });
  });
}

/*
* start the RPC server
*/
function startNetworking() {

  rpcClient.call( {"jsonrpc": "2.0", "method": "start", "id": 1}, function (err, res) {
    if(err) {
      console.log(err);
    }
    else {
      console.log(res);
    }
  })
}

/*
* stop the RPC server
*/
function stopNetworking() {
  rpcClient.call( {"jsonrpc": "2.0", "method": "stop", "id": 1}, function (err, res) {
    if(err) {
      console.log(err);
    }
    else {
      console.log(res);
    }
  })
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
  fs.writeFileSync(appDir+"/"+name+".json", JSON.stringify(data, null, 2))
}

/*
* Reset the nodeFile
*/
function reset(callback) {
  let data = {
    "userData":
      {
        "email":"",
        "name":"",
        "bio":"",
        "ip":"",
        "initialized":false
      },
    "address":""
  }
  fs.writeFileSync(appDir+"/nodeFile.json", JSON.stringify(data, null, 2))
  if (callback != null) {
    callback()
  }
}

// Get the argument that the user provided
var argument = process.argv[2];

status(function() {
  if (argument in options) {
    options[argument].toCall();
  } else {
    console.log(colors.red("[Gladius-Node]") + " Invalid arguments. See " + colors.blue("gladius-node --help") + " for a list of commands");
  }
})

// /** WIP - need to add a stop/kill endpoint
// * Stop accepting connections
// */
// function stop() {
//   axios.put(daemonAddress + "/api/status/", {
//       status: false
//     })
//     .then(function(res) {
//       console.log(res);
//     })
//     .catch(function(err) {
//       console.log(colors.red(
//         "Woah an err! Make sure your daemon is running and can be connected to"
//       ));
//       console.log(err);
//     });
// }
