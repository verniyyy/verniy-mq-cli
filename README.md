# verniy-mq-cli
verniy-mq-cli is a simple Command-Line Interface (CLI) tool designed for interacting with [verniy-mq](https://github.com/verniyyy/verniy-mq), a lightweight message queue system developed by the creator of this tool.

## Installation
### Using `curl`:
Make sure you have `curl` on your machine.

Then, you can install verniy-mq-cli using the following command:
```bash
curl -LJO https://github.com/verniyyy/verniy-mq-cli/releases/download/v1.0.0/vmq-cli
chmod +x vmq-cli
sudo cp -p vmq-cli /usr/local/bin
```

### Using `git clone`` and manual build
Alternatively, you can clone the repository and build verniy-mq-cli manually.

Here's how:
```
git clone https://github.com/verniyyy/verniy-mq-cli.git
cd verniy-mq-cli
sudo make install
```

This will clone the repository and use the provided Makefile to build and install verniy-mq-cli.
The binary will be installed in /usr/local/bin.

## Usage

### Connect to the server
```bash
VERNIY_MQ_PASSWORD=password vmq-cli -H localhost -u root -p 5672
```

###  Create a new queue
```
> create ${new queue name}
```

### Get list of queues
```bash
> list
```

### Delete a queues
```bash
> deleteq ${queue name}
```

### Publish a Message
```bash
> ${queue name} publish
```

### Consume a Message
```bash
> ${queue name} consume
```

### Delete a Message
```bash
> ${queue name} ${message id}
```

### Set a queue
The queue specification can be omitted.
```bash
> use ${queue name}
${queue name}> consume
```

### Close connection
```bash
> quit
```
or
```bash
> exit
```

## Contributing
If you have suggestions, bug reports, or want to contribute to verniy-mq-cli, feel free to open an issue or submit a pull request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.