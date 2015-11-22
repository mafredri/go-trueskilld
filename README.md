# go-trueskilld

A simple http API for [go-trueskill](https://github.com/mafredri/go-trueskill).

## Example

POST http://localhost:8495/rate

```json
{
  "players": [
    {
      "mu": 25.000,
      "sigma": 8.333
    },
	{
      "mu": 25.000,
      "sigma": 8.333
    }
  ]
}
```

Response:

```json
{
  "players": [
    {
      "mu": 29.205,
      "sigma": 7.195,
      "trueskill": 8
    },
    {
      "mu": 20.795,
      "sigma": 7.195,
      "trueskill": 0
    }
  ],
  "probability_of_outcome": 50
}
```