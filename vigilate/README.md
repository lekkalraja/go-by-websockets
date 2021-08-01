## Vigilate

#### Migrate Database
* Install Soda : https://gobuffalo.io/en/docs/db/toolbox/ `(On Mac/Unix : brew install gobuffalo/tap/pop)`
    ```shellscript
        brew install gobuffalo/tap/pop
        go get github.com/gobuffalo/pop/...
        go install github.com/gobuffalo/pop/soda
    ```
* Install Postgres
* Then Execute `soda migrate` -> To Create Tables and populate data
* `./run.sh` <- To Run the Vigilate application

#### Start Pusher Service
* For Configs look @ `./ipe/linux/config.json`
    ```shellscript
        > cd ./ipe/linux
        > ./ipe
    ```

##### Create Database Table Creation Scripts (migrtion (Buffaloo fizz) scripts) using `soda`
* To generate Scripts : `soda generate fizz CreateHostsTable`
    ```shellscript
        raja@raja-Latitude-3460:~/go-by-websockets/vigilate$ soda generate fizz CreateHostsTable
        v4.13.1

        DEBU[2021-08-01T12:36:28+08:00] Step: d4aeccff
        DEBU[2021-08-01T12:36:28+08:00] Chdir: /home/raja/go-by-websockets/vigilate
        DEBU[2021-08-01T12:36:28+08:00] File: /home/raja/go-by-websockets/vigilate/migrations/20210801043628_create_create_hosts_tables.up.fizz
        DEBU[2021-08-01T12:36:28+08:00] File: /home/raja/go-by-websockets/vigilate/migrations/20210801043628_create_create_hosts_tables.down.fizz

        raja@raja-Latitude-3460:~/Documents/coding/golang/go-by-websockets/vigilate$ soda migrate
        v4.13.1

        [POP] 2021/08/01 12:39:43 info - > create_create_hosts_tables
        [POP] 2021/08/01 12:39:43 info - 2.8215 seconds
        [POP] 2021/08/01 12:39:48 info - dumped schema for vigilate
    ```


Source : https://github.com/tsawler/vigilate