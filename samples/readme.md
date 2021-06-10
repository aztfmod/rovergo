# Using the CAF / Symphony reference app with Rover v2

Clone the symphony and rover v2 (rovergo) repos

```bash
git clone https://github.com/aztfmod/rovergo.git
git clone https://github.com/aztfmod/symphony.git
export roverRepo=$(pwd)/rovergo
```

It's important you run the command from the symphony/caf directory
```
cd symphony/caf
```

```bash
rover cd apply -c $roverRepo/samples/ref-app-symphony.yaml -l level0
```

Then deploy the other levels as you require

```bash
rover cd apply -c $roverRepo/samples/ref-app-symphony.yaml -l level1
```

You can also run init, plan or validate to check things first
```bash
rover cd init -c $roverRepo/samples/ref-app-symphony.yaml -l level1
rover cd plan -c $roverRepo/samples/ref-app-symphony.yaml -l level1
rover cd validate -c $roverRepo/samples/ref-app-symphony.yaml -l level1
```

or if you are feeling bold, deploy ALL levels

```bash
rover cd apply -c $roverRepo/samples/ref-app-symphony.yaml
```