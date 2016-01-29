# go-trueskilld

A simple http API for [go-trueskill](https://github.com/mafredri/go-trueskill).

If the game configuration is not provided in the JSON request, the [default configuration](https://github.com/mafredri/go-trueskill/blob/4dbcbb9/trueskill.go#L15-L19) for TrueSkill™ is used.

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

Response

```json
{
  "players": [
    {
      "mu": 29.396,
      "sigma": 7.171,
      "trueskill": 7.882
    },
    {
      "mu": 20.604,
      "sigma": 7.171,
      "trueskill": 0
    }
  ],
  "probability_of_outcome": 47.759
}
```

**NOTE:** If you provide the game configuration in your request, you might encounter rounding errors without sufficient precision. Example (compare probability of outcome to above):

POST

```json
{
  "mu": 25.0,
  "sigma": 8.333,
  "beta": 4.166,
  "tau": 0.083,
  "draw_probability": 10,
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

Response

```json
{
  "players": [
    {
      "mu": 29.396,
      "sigma": 7.171,
      "trueskill": 7.882
    },
    {
      "mu": 20.604,
      "sigma": 7.171,
      "trueskill": 0
    }
  ],
  "probability_of_outcome": 47.76
}
```