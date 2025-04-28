# Go City and Temperature Identifer

`go-city-temperature-identifier` is a simple tool to find a temperature in a city, informed in Celsius, Fahrenheit and Kelvin.

## Prerequisites

- [Docker](https://www.docker.com/) installed on your machine.

## How to Build and Run the Project with Docker

### 1. Build the Docker Image

In the root directory of the project, run the following command to build the Docker image:

```bash
docker build -t go-city-temperature-identifier .
```

### 2. Run the City Temperature Identifier

After building the image, you can run the city temperature identifier by passing the required parameter via the command line:

```bash
docker run --rm go-city-temperature-identifier -cep=<CEP>
```

#### Required Parameters:
- `-cep`: CEP to be searched.

#### Example:
```bash
docker run --rm go-city-temperature-identifier -cep=22222-222
```

### 3. Final Report
After running the tool, the output will display the temperature of the city corresponding to the provided CEP in Celsius, Fahrenheit, and Kelvin. The report will look like this:

```
City: Rio de Janeiro
Temperature:
    - Celsius: 25°C
    - Fahrenheit: 77°F
    - Kelvin: 298.15K
```

This information is fetched using a weather API and is based on the CEP provided in the command.