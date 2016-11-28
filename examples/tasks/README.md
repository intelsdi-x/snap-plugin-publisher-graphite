# Example tasks

[This](psutil-statistics-graphite.json) example task will publish metrics to **graphite** 
from the psutil plugin.  

## Running the example

### Requirements
 * `docker` and `docker-compose` are **installed** and **configured** 

Running the sample is as *easy* as running the script `./run-example.sh`.  

## Files

- [run-example.sh](run-example.sh) 
    - The example is launched with this script.
- [psutil-statistics-graphite.json](psutil-statistics-graphite.json)
    - Snap task used by large test and example script.
- [psutil-statistics-graphite-simple.json](psutil-statistics-graphite-simple.json)
    - Snap task showed in the [main README.md example](../README.md#Examples).
- [docker-compose.yml](docker-compose.yml)
    - A docker compose file which defines two linked containers
        - "runner" is the container where snapteld is run from. You will be dumped 
        into a shell in this container after running 
        [run-example.sh](run-example.sh). Exiting the shell will 
        trigger cleaning up the containers used in the example.
        - "graphite" is the container running graphite.
- [psutil-statistics-graphite.sh](psutil-statistics-graphite.sh)
    - Downloads `snapteld`, `snaptel`, `snap-plugin-publisher-graphite`,
    `snap-plugin-collector-psutil` `snap-plugin-processor-passthru-grpc` and
    starts the task [psutil-statistics-graphite.json](psutil-statistics-graphite.json).
- [.setup.sh](.setup.sh)
    - Verifies dependencies and starts the containers.  It's called 
    by [run-example.sh](run-example.sh).