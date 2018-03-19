// Configuration for the Gladius CLI

var config = {
  controlDaemonPort: 9545, // Port to access control daemon
  controlDaemonAddress: "http://127.0.0.1", // What address to access the daemon
  transferKeyOverHttp: true // Send file path or send actual key (used for docker)
};

module.exports = config;
