# Weather by Zip Code

A Go HTTP API that receives a Brazilian zip code, finds its city through ViaCEP, and returns the current temperature from WeatherAPI in Celsius, Fahrenheit, and Kelvin.

## API

### Get weather

```http
GET /weather/{zipcode}
```

The zip code must contain exactly eight digits.

Example:

```bash
curl http://localhost:8080/weather/01001000
```

Success response:

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

Error responses:

| Status | Body | Condition |
| --- | --- | --- |
| `422 Unprocessable Entity` | `invalid zipcode` | The zip code does not contain exactly eight digits |
| `404 Not Found` | `can not find zipcode` | ViaCEP cannot find the zip code |

## Run locally

Create a free WeatherAPI account and copy your API key. Then run:

```bash
export WEATHER_API_KEY="your-api-key"
go run ./cmd/server
```

The application listens on port `8080` by default. Set `PORT` to use another port.

## Run with Docker

Build the image:

```bash
docker build -t weather-by-zipcode .
```

Start the container:

```bash
docker run --rm \
  -p 8080:8080 \
  -e WEATHER_API_KEY="your-api-key" \
  weather-by-zipcode
```

## Tests

Run the automated tests:

```bash
go test ./...
```

Run the tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Cloud Run

The Cloud Run URL will be added after deployment.
