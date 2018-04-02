// Configuration for the Gladius CLI

// address of the control and edge daemon
var address = process.env.ADDRESS || 'localhost';

var network = {
  controlDaemonPort: 3000, // Port to access control daemon
  controlDaemonAddress: "http://" + address, // What address to access the daemon
  transferKeyOverHttp: true // Send file path or send actual key (used for docker)
};

module.exports = {network};