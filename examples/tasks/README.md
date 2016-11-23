# Example tasks

[This](mock-passthru-graphite.json) example task will publish metrics to **graphite** 
from the mock plugin.  

## Running the example

### Requirements
 * `docker` and `docker-compose` are **installed** and **configured** 

Running the sample is as *easy* as running the script `./run-example.sh`.  

## Files

- [run-example.sh](run-example.sh) 
    - The example is launched with this script     
- [mock-passthru-graphite.json](mock-passthru-graphite.json)
    - Snap task definition
- [docker-compose.yml](docker-compose.yml)
    - A docker compose file which defines two linked containers
        - "runner" is the container where snapteld is run from. You will be dumped 
        into a shell in this container after running 
        [run-example.sh](run-example.sh). Exiting the shell will 
        trigger cleaning up the containers used in the example.
        - "graphite" is the container running graphite.
- [mock-passthru-graphite.sh](mock-passthru-graphite.sh)
    - Downloads `snapteld`, `snaptel`, `snap-plugin-publisher-graphite`,
    `snap-plugin-collector-mock2-grpc` `snap-plugin-processor-passthru-grpc` and
    starts the task [mock-passthru-graphite.json](mock-passthru-graphite.json).
- [.setup.sh](.setup.sh)
    - Verifies dependencies and starts the containers.  It's called 
    by [run-example.sh](run-example.sh).