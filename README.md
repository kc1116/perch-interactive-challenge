# perch-interactive-challenge

## Getting started
There are 4 entities involved in our example. 

- Event Aggregator: Listens for events published to the telemetry topic of a specified device registry 

- Device Simulator: Attempts to simulate perch sessions with randomness
    - Session: The time spent at a device from start to finish, there are n number of Interactions within a single session
    - Interaction: Randomly generated event

- Websocket Proxy: simple websocket endpoint that proxies a rethinkdb change set to a React web app 

- React Webapp: bare bones webpage that initiates a ws connection to our proxy endpoint and displays events as they are 
streamed in  

These things are pictured below, in the diagram we can see how they all relate.

<img src= https://www.lucidchart.com/publicSegments/view/02a306a0-1660-4f84-aa24-c6d605d71f7c/image.png />

## Running All in One Docker Compose 
First thing you need to do is update the .env file, add the path to a GOOGLE_APPLICATION_CREDENTIALS json file. ([docs](https://cloud.google.com/docs/authentication/getting-started))

Now run our make target for compose
```bash
make run_network
```

Change any settings you want in the [docker-compose.yaml](./build/docker-compose.yaml) file. 
You can also run each piece individually without docker for more fine grain usage and testing. 
See [Docs](./docs) for manual usage instructions for CLI tool. 

## Start up our UI manually
```bash
yarn global add serve
serve -s ./build/ui
```

