// Configuration for the Gladius CLI

var config = {
  controlDaemonPort: 5544, // Port to access control daemon
  controlDaemonAddress: "http://localhost", // What address to access the daemon
  transferKeyOverHttp: true // Send file path or send actual key (used for docker)
};

module.exports = config;
