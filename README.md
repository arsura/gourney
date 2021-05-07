# MoonBase Service

## Schema

### Coin

- id: Int
- name: String
- amount: Float
- total: Float
- rise_rate: Float
- rise_factor: Float
- created_at: DateTime
- updated_at: DateTime

### Exchange Rate

- id: Int
- source: <Coin-id>
- target: <Coin-id>
- price: Float
- created_at: DateTime
- updated_at: DateTime

## Example

THBT Coin

```json
{
  "id": 1,
  "name": "THBT",
  "amount": null,
  "total": null,
  "rise_rate": null,
  "rise_factor": null
}
```

MOON Coin

```json
{
  "id": 2,
  "name": "MOON",
  "amount": 1000.0,
  "total": 1000.0,
  "rise_rate": 0.1,
  "rise_factor": 10.0
}
```

Exchange Rate THBT -> MOON

```json
{
  "id": 1,
  "source": 2,
  "target": 1,
  "price": 50.0
}
```

## Example 1

### Current Stats

```text
50 THBT = 1 MOON
Amount 1000 MOON
```

### Input

```text
100 THBT -> 2 MOON
Slippage Tolerance = 1%

1. Calculate edge of acceptable value
  2 + (2 * 0.01) = 2.02 MOON Upper bound
  2 - (2 * 0.01) = 1.98 MOON Lower bound

2. Calculate amount received token

  source: 1 (THBT),
  target: 2 (MOON),
  price: 50.0

  100.0 / 50.0 = 2.00 (Is in range of acceptable value)

  Get current amount of MOON = 1000.0 to calculate remaining MOON
  1000.0 - 2.00 = 998.0 <--- This value will save to DB
```

## Example 2

### Current Stats

```text
50 THBT = 1 MOON
Amount 991 MOON
```

### Input

```text
1900 THBT -> x MOON
Slippage Tolerance = 1%

1. Calculate edge of acceptable value
  28 + (28 * 0.01) = 28.28 MOON Upper bound
  28 - (28 * 0.01) = 27.72 MOON Lower bound

2. Calculate amount received token

  source: 1 (THBT),
  target: 2 (MOON),
  price: 50.0

  currentPrice = 50.0
  currentAmountCoin = 991.0

  1900 - 50 = 
  991.0




```


Rate
50 THBT = 1 MOON

calcExchangeRates(THBT, MOON, 22) == 0.44 MOON
50 THBT = 1 MOON
22 / 50 = 0.44 MOON

calcExchangeRates(MOON, THBT, 21) == 1050 THBT
1 MOON = 50 THBT
21 MOON = 1050 THBT




THBT coin -> MOON coin



1000 MOON 
50 THBT -> 1 MOON

ถ้าเหรียญถูกซื้อไปถูก ๆ 10 เหรียญ มูลค่าของเหรียญจะเพิ่มขึ้น 10% 

990 MOON
55 THBT -> 1 MOON

980 MOON
60.5 THBT -> 1 MOON


1000 THTB -> x MOON

เริ่ม 50THBT/MOON 500 THTB ->           = 10 MOON ,             55THBT/MOON
    55THBT/MOON 500 THTB -> 500 / 55  = 9.09090909091 MOON ,  60.5THBT/MOON

    = 10 + 9.09090909091 = 19.09090909091 MOON

991 ซื้อ 100
50 => 1
50 / 55 => 1.909091


1000 - 991 => 50
990 - 981 => 55
980 -> 971 => 60.5
970 -> 961 => 66.55
...       => 73.205

967.33554

ซื้อที่ 0.5 บาท

66.55 = 1 
1
