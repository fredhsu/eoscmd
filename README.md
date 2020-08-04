# eoscmd
Command line tool for EOS

## Run command against array of EOS
Goal: take in a list of eos devices in json format, execute cli command via eapi against them all and return the result in a series of text files
Inputs: json file, command, output directory
Outputs: multiple text files for each switch and command run + timestamp?
Note: user may need to create a directory called "output" in their cloned repo, until further development is done

e.g.

    cat devices.json | eoscmd "show tech" -o ./showtech
=======

Device format is based on JSON verison of Ansible inventory list.  Here is an example file: (devices.json)[devices.json].  If not using a pipe, it will optionally look for a `devices.json` file in the local directory 

## Example devices input
```
    {
    "hosts": ["dmz-lf11", "dmz-lf12"],
    "vars": {
        "username": "fredlhsu",
        "password": "arista",
        "transport": "https",
        "port": 443
        }
    }
```

## Providing devices to run against
1. Specify file with -f option
2. Stdin/Pipe
3. devices.json file in local directory

## Credentials
Can be specified in the devices JSON mentioned above.  Or can be specified as environment variables EAPI_USERNAME, EAPI_PASSWORD.  Device file and pipes override the environment variables.

## Output
Output is directed to a path given by the -o parameter
Currently it will need an `output` directory created in the same folder for the output to be sent to
