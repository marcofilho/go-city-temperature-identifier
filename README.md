# Go City and Temperature Identifier

`go-city-temperature-identifier` is a simple web service to find the temperature in a city, informed in Celsius, Fahrenheit, and Kelvin, based on a provided CEP (Brazilian postal code).

## Live Demo

This application is deployed on **Google Cloud Run** and can be accessed at:

```
https://go-city-temperature-identifier-iko4y6cnuq-uc.a.run.app/cep?cep=<CEP_TO_INFORM>
```

Replace `<CEP_TO_INFORM>` with the desired CEP.  
**Example:**  
[https://go-city-temperature-identifier-iko4y6cnuq-uc.a.run.app/cep?cep=22222-222](https://go-city-temperature-identifier-iko4y6cnuq-uc.a.run.app/cep?cep=22222-222)

## Prerequisites (for local development)

- [Docker](https://www.docker.com/) installed on your machine.

## How to Build and Run the Project with Docker (Optional)

### 1. Build the Docker Image

In the root directory of the project, run:

```bash
docker build -t go-city-temperature-identifier .
```

### 2. Run the City Temperature Identifier Locally

After building the image, you can run the service:

```bash
docker run --rm -p 8080:8080 go-city-temperature-identifier
```

Then access:  
[http://localhost:8080/cep?cep=22222-222](http://localhost:8080/cep?cep=22222-222)

## API Usage

- **Endpoint:** `/cep`
- **Method:** `GET`
- **Query Parameter:** `cep` (required)

**Example Request:**
```
GET https://go-city-temperature-identifier-iko4y6cnuq-uc.a.run.app/cep?cep=22222-222
```

**Example Response:**
```json
{
  "city_name": "Rio de Janeiro",
  "city_region": "Rio de Janeiro",
  "city_country": "Brazil",
  "temp_C": 25,
  "temp_F": 77,
  "temp_K": 298.15
}
```

This information is fetched using a weather API and is based on the CEP provided in the request.