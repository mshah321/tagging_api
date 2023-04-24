# Environment Variables

To get started, you need to copy the `env.sample` file and rename it to `.env`. This file contains environment variables that are needed for the application to function properly. 

## Setting up KEY variable

The `KEY` variable is used as a placeholder for the actual value of the API key that is required by the application. Before running the application, you need to replace `KEY` with the actual API key.

To set the value of `KEY`, open the `.env` file in a text editor and replace `KEY` with the actual API key provided by the service you want to use. Make sure you do not include any quotes or whitespace around the value.

```
# Example .env file
KEY=your_api_key_here
``` 

Once you have updated the `.env` file with the correct value for `KEY`, save the file and the application will be able to access the API using the specified key.

# GET /get_tags

This REST endpoint retrieves tags within a specified radius of a given latitude and longitude. It returns an array of tag names and their count, as well as the input parameters.

## Request

The request must be sent as an HTTP GET request with the following parameters:

- `lat`: The latitude of the center point (float)
- `lng`: The longitude of the center point (float)
- `radius`: The radius from the center point (int)

Example:

```
GET /get_tags?lat=37.7749&lng=-122.4194&radius=500
```

## Response

The response will be a JSON object with the following fields:

- `input_params`: An object containing the input parameters:
  - `lat`: The latitude of the center point (float)
  - `lng`: The longitude of the center point (float)
  - `radius`: The radius from the center point (int)
- `result`: An array of objects representing tags and their count:
  - `tag`: The name of the tag (string)
  - `count`: The number of occurrences of the tag within the specified radius (int)

Example response:

```
{
  "input_params": {
    "lat": 37.7749,
    "lng": -122.4194,
    "radius": 500
  },
  "result": [
    {
      "tag": "Coffee Shop",
      "count": 10
    },
    {
      "tag": "Park",
      "count": 5
    }
  ]
}
```

If no tags are found within the specified radius, the response will be an empty array:

```
{
  "input_params": {
    "lat": 37.7749,
    "lng": -122.4194,
    "radius": 500
  },
  "result": []
}
```

## Errors

If there is an error processing the request, the response will contain an appropriate HTTP status code and an error message in the body:

- `400 Bad Request`: Invalid parameters or lat/lng are not valid coordinates
- `401 Unauthorized`: API token is missing.
