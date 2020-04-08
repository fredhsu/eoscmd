# eoscmd
Command line tool for EOS

## Run command against array of EOS
Goal: take in a list of eos devices in json format, execute cli command via eapi against them all and return the result in a series of text files
Inputs: json file, command, output directory
Outputs: multiple text files for each switch and command run + timestamp?
Note: user may need to create a directory called "output" in their cloned repo, until further development is done

e.g.

    cat devices.json | eoscmd "show tech" -o ./showtech
